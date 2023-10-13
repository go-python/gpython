// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compile python code
package compile

// FIXME name mangling
// FIXME kill ast.Identifier and turn into string?

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/parser"
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/symtable"
	"github.com/go-python/gpython/vm"
)

type loopType byte

// type of loop
const (
	loopLoop loopType = iota
	exceptLoop
	finallyTryLoop
	finallyEndLoop
)

// Loop - used to track loops, try/except and try/finally
type loop struct {
	Start *Label
	End   *Label
	Type  loopType
}

// Loopstack
type loopstack []loop

// Push a loop
func (ls *loopstack) Push(l loop) {
	*ls = append(*ls, l)
}

// Pop a loop
func (ls *loopstack) Pop() {
	*ls = (*ls)[:len(*ls)-1]
}

// Return current loop or nil for none
func (ls loopstack) Top() *loop {
	if len(ls) == 0 {
		return nil
	}
	return &ls[len(ls)-1]
}

type compilerScopeType uint8

const (
	compilerScopeModule compilerScopeType = iota
	compilerScopeClass
	compilerScopeFunction
	compilerScopeLambda
	compilerScopeComprehension
)

// State for the compiler
type compiler struct {
	Code        *py.Code // code being built up
	Filename    string
	Lineno      int // current line number
	OpCodes     Instructions
	loops       loopstack
	SymTable    *symtable.SymTable
	scopeType   compilerScopeType
	qualname    string
	private     string
	parent      *compiler
	depth       int
	interactive bool
}

// Set in py to avoid circular import
func init() {
	py.Compile = Compile
}

// Compile(src, srcDesc, compileMode, flags, dont_inherit) -> code object
//
// Compile the source string (a Python module, statement or expression)
// into a code object that can be executed.
//
// srcDesc is used for run-time error messages and is typically a file system pathname,
//
// See py.CompileMode for compile mode options.
//
// The flags argument, if present, controls which future statements influence
// the compilation of the code.
//
// The dont_inherit argument, if non-zero, stops the compilation inheriting
// the effects of any future statements in effect in the code calling
// compile; if absent or zero these statements do influence the compilation,
// in addition to any features explicitly specified.
func Compile(src, srcDesc string, mode py.CompileMode, futureFlags int, dont_inherit bool) (*py.Code, error) {
	// Parse Ast
	Ast, err := parser.Parse(bytes.NewBufferString(src), srcDesc, mode)
	if err != nil {
		return nil, err
	}
	// Make symbol table
	SymTable, err := symtable.NewSymTable(Ast, srcDesc)
	if err != nil {
		return nil, err
	}
	c := newCompiler(nil, compilerScopeModule)
	c.Filename = srcDesc
	err = c.compileAst(Ast, srcDesc, futureFlags, dont_inherit, SymTable)
	if err != nil {
		return nil, err
	}
	return c.Code, nil
}

// Make a new compiler object with empty code object
func newCompiler(parent *compiler, scopeType compilerScopeType) *compiler {
	code := &py.Code{
		Firstlineno: 1,          // FIXME
		Name:        "<module>", // FIXME
	}
	c := &compiler{
		Code:        code,
		parent:      parent,
		scopeType:   scopeType,
		depth:       1,
		interactive: false,
	}
	if parent != nil {
		c.depth = parent.depth + 1
		c.Filename = parent.Filename
	}
	return c
}

// Panics abount a syntax error on this ast node
func (c *compiler) panicSyntaxErrorf(Ast ast.Ast, format string, a ...interface{}) {
	err := py.ExceptionNewf(py.SyntaxError, format, a...)
	err = py.MakeSyntaxError(err, c.Filename, Ast.GetLineno(), Ast.GetColOffset(), "")
	panic(err)
}

// Sets Lineno from an ast node
func (c *compiler) SetLineno(node ast.Ast) {
	c.Lineno = node.GetLineno()
}

// Create a new compiler object at Ast, using private for name mangling
func (c *compiler) newCompilerScope(compilerScope compilerScopeType, Ast ast.Ast, private string) (newC *compiler) {
	newSymTable := c.SymTable.FindChild(Ast)
	if newSymTable == nil {
		panic(fmt.Sprintf("No symtable found for scope type %v", compilerScope))
	}

	newC = newCompiler(c, compilerScope)

	/* use the class name for name mangling */
	newC.private = private

	if newSymTable.NeedsClassClosure {
		// Cook up a implicit __class__ cell.
		if compilerScope != compilerScopeClass {
			panic("class closure not in class")
		}
		newSymTable.Symbols["__class__"] = symtable.Symbol{Scope: symtable.ScopeCell}
	}

	err := newC.compileAst(Ast, c.Code.Filename, 0, false, newSymTable)
	if err != nil {
		panic(err)
	}
	return newC
}

// Compile an Ast with the current compiler
func (c *compiler) compileAst(Ast ast.Ast, filename string, futureFlags int, dont_inherit bool, SymTable *symtable.SymTable) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = py.MakeException(r)
		}
	}()
	c.SymTable = SymTable
	code := c.Code
	code.Filename = filename
	code.Varnames = append(code.Varnames, SymTable.Varnames...)
	code.Cellvars = append(code.Cellvars, SymTable.Find(symtable.ScopeCell, 0)...)
	code.Freevars = append(code.Freevars, SymTable.Find(symtable.ScopeFree, symtable.DefFreeClass)...)
	code.Flags = c.codeFlags(SymTable) | int32(futureFlags&py.CO_COMPILER_FLAGS_MASK)
	valueOnStack := false
	c.SetLineno(Ast)
	switch node := Ast.(type) {
	case *ast.Module:
		c.Stmts(c.docString(node.Body, false))
	case *ast.Interactive:
		c.interactive = true
		c.Stmts(node.Body)
	case *ast.Expression:
		c.Expr(node.Body)
		valueOnStack = true
	case *ast.Suite:
		panic("suite should not be possible")
	case *ast.Lambda:
		code.Argcount = int32(len(node.Args.Args))
		code.Kwonlyargcount = int32(len(node.Args.Kwonlyargs))
		// Make None the first constant as lambda can't have a docstring
		c.Const(py.None)
		code.Name = "<lambda>"
		c.setQualname()
		c.Expr(node.Body)
		valueOnStack = true
	case *ast.FunctionDef:
		code.Argcount = int32(len(node.Args.Args))
		code.Kwonlyargcount = int32(len(node.Args.Kwonlyargs))
		code.Name = string(node.Name)
		c.setQualname()
		c.Stmts(c.docString(node.Body, true))
	case *ast.ClassDef:
		code.Name = string(node.Name)
		/* load (global) __name__ ... */
		c.NameOp("__name__", ast.Load)
		/* ... and store it as __module__ */
		c.NameOp("__module__", ast.Store)
		c.setQualname()
		if c.qualname == "" {
			panic("Need qualname")
		}

		c.LoadConst(py.String(c.qualname))
		c.NameOp("__qualname__", ast.Store)

		/* compile the body proper */
		c.Stmts(c.docString(node.Body, false))

		if SymTable.NeedsClassClosure {
			/* return the (empty) __class__ cell */
			i := c.FindId("__class__", code.Cellvars)
			if i != 0 {
				panic("__class__ must be first constant")
			}
			/* Return the cell where to store __class__ */
			c.OpArg(vm.LOAD_CLOSURE, uint32(i))
		} else {
			if len(code.Cellvars) != 0 {
				panic("Can't have cellvars without closure")
			}
			/* This happens when nobody references the cell. Return None. */
			c.LoadConst(py.None)
		}
		c.Op(vm.RETURN_VALUE)
	case *ast.ListComp:
		// Elt        Expr
		// Generators []Comprehension
		valueOnStack = true
		code.Name = "<listcomp>"
		c.OpArg(vm.BUILD_LIST, 0)
		c.comprehensionGenerator(node.Generators, 0, node.Elt, nil, Ast)
	case *ast.SetComp:
		// Elt        Expr
		// Generators []Comprehension
		valueOnStack = true
		code.Name = "<setcomp>"
		c.OpArg(vm.BUILD_SET, 0)
		c.comprehensionGenerator(node.Generators, 0, node.Elt, nil, Ast)
	case *ast.DictComp:
		// Key        Expr
		// Value      Expr
		// Generators []Comprehension
		valueOnStack = true
		code.Name = "<dictcomp>"
		c.OpArg(vm.BUILD_MAP, 0)
		c.comprehensionGenerator(node.Generators, 0, node.Key, node.Value, Ast)
	case *ast.GeneratorExp:
		// Elt        Expr
		// Generators []Comprehension
		code.Name = "<genexpr>"
		c.comprehensionGenerator(node.Generators, 0, node.Elt, nil, Ast)

	default:
		panic(fmt.Sprintf("Unknown ModuleBase: %v", Ast))
	}
	if !c.OpCodes.EndsWithReturn() {
		// add a return
		if !valueOnStack {
			// return None if there is nothing on the stack
			c.LoadConst(py.None)
		}
		c.Op(vm.RETURN_VALUE)
	}
	code.Code = c.OpCodes.Assemble()
	code.Stacksize = int32(c.OpCodes.StackDepth())
	code.Nlocals = int32(len(code.Varnames))
	code.Lnotab = string(c.OpCodes.Lnotab())
	code.InitCell2arg()
	return nil
}

// Check for docstring as first Expr in body and remove it and set the
// first constant if found if fn is set, or set __doc__ if it isn't
func (c *compiler) docString(body []ast.Stmt, fn bool) []ast.Stmt {
	var docstring *ast.Str
	if len(body) > 0 {
		if expr, ok := body[0].(*ast.ExprStmt); ok {
			if docstring, ok = expr.Value.(*ast.Str); ok {
				body = body[1:]
			}
		}
	}
	if fn {
		if docstring != nil {
			c.Const(docstring.S)
		} else {
			// If no docstring put None in
			c.Const(py.None)
		}
	} else {
		if docstring != nil {
			c.LoadConst(docstring.S)
			c.NameOp("__doc__", ast.Store)
		}
	}
	return body
}

// Compiles a python constant
//
// Returns the index into the Consts tuple
func (c *compiler) Const(obj py.Object) uint32 {
	// FIXME back this with a dict to stop O(N**2) behaviour on lots of consts
	for i, c := range c.Code.Consts {
		if obj.Type() == c.Type() {
			eq, err := py.Eq(obj, c)
			if err != nil {
				log.Printf("compiler: Const: error %v", err) // FIXME
			} else if eq == py.True {
				return uint32(i)
			}
		}
	}
	c.Code.Consts = append(c.Code.Consts, obj)
	return uint32(len(c.Code.Consts) - 1)
}

// Loads a constant
func (c *compiler) LoadConst(obj py.Object) {
	c.OpArg(vm.LOAD_CONST, c.Const(obj))
}

// Finds the Id in the slice provided, returning -1 if not found
func (c *compiler) FindId(Id string, Names []string) int {
	// FIXME back this with a dict to stop O(N**2) behaviour on lots of vars
	for i, s := range Names {
		if Id == s {
			return i
		}
	}
	return -1
}

// Returns the index into the slice provided, updating the slice if necessary
func (c *compiler) Index(Id string, Names *[]string) uint32 {
	i := c.FindId(Id, *Names)
	if i >= 0 {
		return uint32(i)
	}
	*Names = append(*Names, Id)
	return uint32(len(*Names) - 1)
}

// Compiles a python name
//
// Returns the index into the Name tuple
func (c *compiler) Name(Id ast.Identifier) uint32 {
	return c.Index(string(Id), &c.Code.Names)
}

// Adds this opcode with mangled name as an argument
func (c *compiler) OpName(opcode vm.OpCode, name ast.Identifier) {
	// FIXME mangled := _Py_Mangle(c->u->u_private, o);
	mangled := name
	c.OpArg(opcode, c.Name(mangled))
}

// Compiles an instruction with an argument
func (c *compiler) OpArg(Op vm.OpCode, Arg uint32) {
	if !Op.HAS_ARG() {
		panic("OpArg called with an instruction which doesn't take an Arg")
	}
	instr := &OpArg{Op: Op, Arg: Arg}
	instr.SetLineno(c.Lineno)
	c.OpCodes.Add(instr)
}

// Compiles an instruction without an argument
func (c *compiler) Op(op vm.OpCode) {
	if op.HAS_ARG() {
		panic("Op called with an instruction which takes an Arg")
	}
	instr := &Op{Op: op}
	instr.SetLineno(c.Lineno)
	c.OpCodes.Add(instr)
}

// Inserts an existing label
func (c *compiler) Label(Dest *Label) {
	c.OpCodes.Add(Dest)
}

// Inserts and creates a label
func (c *compiler) NewLabel() *Label {
	Dest := new(Label)
	c.OpCodes.Add(Dest)
	return Dest
}

// Compiles a jump instruction
func (c *compiler) Jump(Op vm.OpCode, Dest *Label) {
	var instr Instruction
	switch Op {
	case vm.JUMP_IF_FALSE_OR_POP, vm.JUMP_IF_TRUE_OR_POP, vm.JUMP_ABSOLUTE, vm.POP_JUMP_IF_FALSE, vm.POP_JUMP_IF_TRUE, vm.CONTINUE_LOOP: // Absolute
		instr = &JumpAbs{OpArg: OpArg{Op: Op}, Dest: Dest}
	case vm.JUMP_FORWARD, vm.SETUP_WITH, vm.FOR_ITER, vm.SETUP_LOOP, vm.SETUP_EXCEPT, vm.SETUP_FINALLY:
		instr = &JumpRel{OpArg: OpArg{Op: Op}, Dest: Dest}
	default:
		panic("Jump called with non jump instruction")
	}
	instr.SetLineno(c.Lineno)
	c.OpCodes.Add(instr)
}

/*
The test for LOCAL must come before the test for FREE in order to

	handle classes where name is both local and free.  The local var is
	a method and the free var is a free var referenced within a method.
*/
func (c *compiler) getRefType(name string) symtable.Scope {
	if c.scopeType == compilerScopeClass && name == "__class__" {
		return symtable.ScopeCell
	}
	scope := c.SymTable.GetScope(name)
	if scope == symtable.ScopeInvalid {
		panic(fmt.Sprintf("compile: getRefType: unknown scope for %s in %s\nsymbols: %v\nlocals: %s\nglobals: %s", name, c.Code.Name, c.SymTable.Symbols, c.Code.Varnames, c.Code.Names))
	}
	return scope
}

// makeClosure constructs the function or closure for a func/class/lambda etc
func (c *compiler) makeClosure(code *py.Code, args uint32, child *compiler, qualname string) {
	free := uint32(len(code.Freevars))

	if free == 0 {
		c.LoadConst(code)
		c.LoadConst(py.String(qualname))
		c.OpArg(vm.MAKE_FUNCTION, args)
		return
	}
	for i := range code.Freevars {
		/* Bypass com_addop_varname because it will generate
		   LOAD_DEREF but LOAD_CLOSURE is needed.
		*/
		name := code.Freevars[i]

		/* Special case: If a class contains a method with a
		   free variable that has the same name as a method,
		   the name will be considered free *and* local in the
		   class.  It should be handled by the closure, as
		   well as by the normal name loookup logic.
		*/
		reftype := c.getRefType(name)
		arg := 0
		if reftype == symtable.ScopeCell {
			arg = c.FindId(name, c.Code.Cellvars)
		} else { /* (reftype == FREE) */
			// using CellAndFreeVars in closures requires skipping Cellvars
			arg = len(c.Code.Cellvars) + c.FindId(name, c.Code.Freevars)
		}
		if arg < 0 {
			panic(fmt.Sprintf("compile: makeClosure: lookup %q in %q %v %v\nfreevars of %q: %v\n", name, c.SymTable.Name, reftype, arg, code.Name, code.Freevars))
		}
		c.OpArg(vm.LOAD_CLOSURE, uint32(arg))
	}
	c.OpArg(vm.BUILD_TUPLE, free)
	c.LoadConst(code)
	c.LoadConst(py.String(qualname))
	c.OpArg(vm.MAKE_CLOSURE, args)
}

// Compute the flags for the current Code
func (c *compiler) codeFlags(st *symtable.SymTable) (flags int32) {
	if st.Type == symtable.FunctionBlock {
		flags |= py.CO_NEWLOCALS
		if st.Unoptimized == 0 {
			flags |= py.CO_OPTIMIZED
		}
		if st.Nested {
			flags |= py.CO_NESTED
		}
		if st.Generator {
			flags |= py.CO_GENERATOR
		}
		if st.Varargs {
			flags |= py.CO_VARARGS
		}
		if st.Varkeywords {
			flags |= py.CO_VARKEYWORDS
		}
	}

	/* (Only) inherit compilerflags in PyCF_MASK */
	flags |= c.Code.Flags & py.CO_COMPILER_FLAGS_MASK

	if len(c.Code.Freevars) == 0 && len(c.Code.Cellvars) == 0 {
		flags |= py.CO_NOFREE
	}

	return flags
}

// Sets the qualname
func (c *compiler) setQualname() {
	var base string
	if c.depth > 1 {
		force_global := false
		parent := c.parent
		if parent == nil {
			panic("compile: setQualname: expecting a parent")
		}
		if c.scopeType == compilerScopeFunction || c.scopeType == compilerScopeClass {
			// FIXME mangled = _Py_Mangle(parent.u_private, u.u_name)
			mangled := c.Code.Name
			scope := parent.SymTable.GetScope(mangled)
			if scope == symtable.ScopeGlobalImplicit {
				panic("compile: setQualname: not expecting scopeGlobalImplicit")
			}
			if scope == symtable.ScopeGlobalExplicit {
				force_global = true
			}
		}
		if !force_global {
			if parent.scopeType == compilerScopeFunction || parent.scopeType == compilerScopeLambda {
				base = parent.qualname + ".<locals>"
			} else {
				base = parent.qualname
			}
		}
	}
	if base != "" {
		c.qualname = base + "." + c.Code.Name
	} else {
		c.qualname = c.Code.Name
	}
}

// Compile a function
func (c *compiler) compileFunc(compilerScope compilerScopeType, Ast ast.Ast, Args *ast.Arguments, DecoratorList []ast.Expr, Returns ast.Expr) {
	newC := c.newCompilerScope(compilerScope, Ast, "")
	newC.Code.Argcount = int32(len(Args.Args))
	newC.Code.Kwonlyargcount = int32(len(Args.Kwonlyargs))

	// Defaults
	c.Exprs(Args.Defaults)

	// KwDefaults
	if len(Args.KwDefaults) > len(Args.Kwonlyargs) {
		panic("compile: more KwDefaults than Kwonlyargs")
	}
	for i := range Args.KwDefaults {
		c.LoadConst(py.String(Args.Kwonlyargs[i].Arg))
		c.Expr(Args.KwDefaults[i])
	}

	// Annotations
	annotations := py.Tuple{}
	addAnnotation := func(args ...*ast.Arg) {
		for _, arg := range args {
			if arg != nil && arg.Annotation != nil {
				c.Expr(arg.Annotation)
				annotations = append(annotations, py.String(arg.Arg))
			}
		}
	}
	addAnnotation(Args.Args...)
	addAnnotation(Args.Vararg)
	addAnnotation(Args.Kwonlyargs...)
	addAnnotation(Args.Kwarg)
	if Returns != nil {
		c.Expr(Returns)
		annotations = append(annotations, py.String("return"))
	}
	num_annotations := uint32(len(annotations))
	if num_annotations > 0 {
		num_annotations++ // include the tuple
		c.LoadConst(annotations)
	}

	// Load decorators onto stack
	c.Exprs(DecoratorList)

	// Make function or closure, leaving it on the stack
	posdefaults := uint32(len(Args.Defaults))
	kwdefaults := uint32(len(Args.KwDefaults))
	args := uint32(posdefaults + (kwdefaults << 8) + (num_annotations << 16))
	c.makeClosure(newC.Code, args, newC, newC.qualname)

	// Call decorators
	for range DecoratorList {
		c.OpArg(vm.CALL_FUNCTION, 1) // 1 positional, 0 keyword pair
	}
}

// Compile class definition
func (c *compiler) class(Ast ast.Ast, class *ast.ClassDef) {
	// Load decorators onto stack
	c.Exprs(class.DecoratorList)

	/* ultimately generate code for:
	     <name> = __build_class__(<func>, <name>, *<bases>, **<keywords>)
	   where:
	     <func> is a function/closure created from the class body;
	        it has a single argument (__locals__) where the dict
	        (or MutableSequence) representing the locals is passed
	     <name> is the class name
	     <bases> is the positional arguments and *varargs argument
	     <keywords> is the keyword arguments and **kwds argument
	   This borrows from compiler_call.
	*/

	/* 1. compile the class body into a code object */
	newC := c.newCompilerScope(compilerScopeClass, Ast, string(class.Name))
	// newSymTable := c.SymTable.FindChild(Ast)
	// if newSymTable == nil {
	// 	panic("No symtable found for class")
	// }
	// newC := newCompiler(c, compilerScopeClass)
	/* use the class name for name mangling */
	newC.private = string(class.Name)
	// code, err := newC.compileAst(Ast, c.Code.Filename, 0, false, newSymTable)
	// if err != nil {
	// 	panic(err)
	// }

	/* 2. load the 'build_class' function */
	c.Op(vm.LOAD_BUILD_CLASS)

	/* 3. load a function (or closure) made from the code object */
	c.makeClosure(newC.Code, 0, newC, string(class.Name))

	/* 4. load class name */
	c.LoadConst(py.String(class.Name))

	/* 5. generate the rest of the code for the call */
	c.callHelper(2, class.Bases, class.Keywords, class.Starargs, class.Kwargs)

	/* 6. apply decorators */
	for range class.DecoratorList {
		c.OpArg(vm.CALL_FUNCTION, 1) // 1 positional, 0 keyword pair
	}

	/* 7. store into <name> */
	c.NameOp(string(class.Name), ast.Store)
}

/*
Implements the with statement from PEP 343.

The semantics outlined in that PEP are as follows:

with EXPR as VAR:

	BLOCK

It is implemented roughly as:

context = EXPR
exit = context.__exit__  # not calling it
value = context.__enter__()
try:

	VAR = value  # if VAR present in the syntax
	BLOCK

finally:

	if an exception was raised:
	exc = copy of (exception, instance, traceback)
	else:
	exc = (None, None, None)
	exit(*exc)
*/
func (c *compiler) with(node *ast.With, pos int) {
	item := node.Items[pos]
	finally := new(Label)

	/* Evaluate EXPR */
	c.Expr(item.ContextExpr)
	c.Jump(vm.SETUP_WITH, finally)

	/* SETUP_WITH pushes a finally block. */
	c.loops.Push(loop{Type: finallyTryLoop})
	if item.OptionalVars != nil {
		c.Expr(item.OptionalVars)
	} else {
		/* Discard result from context.__enter__() */
		c.Op(vm.POP_TOP)
	}

	pos++
	if pos == len(node.Items) {
		/* BLOCK code */
		c.Stmts(node.Body)
	} else {
		c.with(node, pos)
	}

	/* End of try block; start the finally block */
	c.Op(vm.POP_BLOCK)
	c.loops.Pop()
	c.LoadConst(py.None)

	/* Finally block starts; context.__exit__ is on the stack under
	   the exception or return information. Just issue our magic
	   opcode. */
	c.Label(finally)
	c.Op(vm.WITH_CLEANUP)

	/* Finally block ends. */
	c.Op(vm.END_FINALLY)
}

/*
Code generated for "try: <body> finally: <finalbody>" is as follows:

	     SETUP_FINALLY           L
	     <code for body>
	     POP_BLOCK
	     LOAD_CONST              <None>
	 L:          <code for finalbody>
	     END_FINALLY

	The special instructions use the block stack.  Each block
	stack entry contains the instruction that created it (here
	SETUP_FINALLY), the level of the value stack at the time the
	block stack entry was created, and a label (here L).

	SETUP_FINALLY:
	 Pushes the current value stack level and the label
	 onto the block stack.
	POP_BLOCK:
	 Pops en entry from the block stack, and pops the value
	 stack until its level is the same as indicated on the
	 block stack.  (The label is ignored.)
	END_FINALLY:
	 Pops a variable number of entries from the *value* stack
	 and re-raises the exception they specify.  The number of
	 entries popped depends on the (pseudo) exception type.

	The block stack is unwound when an exception is raised:
	when a SETUP_FINALLY entry is found, the exception is pushed
	onto the value stack (and the exception condition is cleared),
	and the interpreter jumps to the label gotten from the block
	stack.
*/
func (c *compiler) tryFinally(node *ast.Try) {
	end := new(Label)
	c.Jump(vm.SETUP_FINALLY, end)
	if len(node.Handlers) > 0 {
		c.tryExcept(node)
	} else {
		c.loops.Push(loop{Type: finallyTryLoop})
		c.Stmts(node.Body)
		c.loops.Pop()
	}
	c.Op(vm.POP_BLOCK)
	c.LoadConst(py.None)
	c.Label(end)
	c.loops.Push(loop{Type: finallyEndLoop})
	c.Stmts(node.Finalbody)
	c.loops.Pop()
	c.Op(vm.END_FINALLY)
}

/*
Code generated for "try: S except E1 as V1: S1 except E2 as V2: S2 ...":
(The contents of the value stack is shown in [], with the top
at the right; 'tb' is trace-back info, 'val' the exception's
associated value, and 'exc' the exception.)

Value stack          Label   Instruction     Argument
[]                           SETUP_EXCEPT    L1
[]                           <code for S>
[]                           POP_BLOCK
[]                           JUMP_FORWARD    L0

[tb, val, exc]       L1:     DUP                             )
[tb, val, exc, exc]          <evaluate E1>                   )
[tb, val, exc, exc, E1]      COMPARE_OP      EXC_MATCH       ) only if E1
[tb, val, exc, 1-or-0]       POP_JUMP_IF_FALSE       L2      )
[tb, val, exc]               POP
[tb, val]                    <assign to V1>  (or POP if no V1)
[tb]                         POP
[]                           <code for S1>

	JUMP_FORWARD    L0

[tb, val, exc]       L2:     DUP
.............................etc.......................

[tb, val, exc]       Ln+1:   END_FINALLY     # re-raise exception

[]                   L0:     <next statement>

Of course, parts are not generated if Vi or Ei is not present.
*/
func (c *compiler) tryExcept(node *ast.Try) {
	c.loops.Push(loop{Type: exceptLoop})
	except := new(Label)
	orelse := new(Label)
	end := new(Label)
	c.Jump(vm.SETUP_EXCEPT, except)
	c.Stmts(node.Body)
	c.Op(vm.POP_BLOCK)
	c.Jump(vm.JUMP_FORWARD, orelse)
	n := len(node.Handlers)
	c.Label(except)
	for i, handler := range node.Handlers {
		if handler.ExprType == nil && i < n-1 {
			c.panicSyntaxErrorf(handler, "default 'except:' must be last")
		}
		// FIXME c.u.u_lineno_set = 0
		// c.u.u_lineno = handler.lineno
		// c.u.u_col_offset = handler.col_offset
		except := new(Label)
		if handler.ExprType != nil {
			c.Op(vm.DUP_TOP)
			c.Expr(handler.ExprType)
			c.OpArg(vm.COMPARE_OP, vm.PyCmp_EXC_MATCH)
			c.Jump(vm.POP_JUMP_IF_FALSE, except)
		}
		c.Op(vm.POP_TOP)
		if handler.Name != "" {
			cleanup_end := new(Label)
			c.NameOp(string(handler.Name), ast.Store)
			c.Op(vm.POP_TOP)

			/*
			   try:
			       # body
			   except type as name:
			       try:
			           # body
			       finally:
			           name = None
			           del name
			*/

			/* second try: */
			c.Jump(vm.SETUP_FINALLY, cleanup_end)

			/* second # body */
			c.Stmts(handler.Body)
			c.Op(vm.POP_BLOCK)
			c.Op(vm.POP_EXCEPT)

			/* finally: */
			c.LoadConst(py.None)
			c.Label(cleanup_end)

			/* name = None */
			c.LoadConst(py.None)
			c.NameOp(string(handler.Name), ast.Store)

			/* del name */
			c.NameOp(string(handler.Name), ast.Del)

			c.Op(vm.END_FINALLY)
		} else {
			c.Op(vm.POP_TOP)
			c.Op(vm.POP_TOP)
			c.Stmts(handler.Body)
			c.Op(vm.POP_EXCEPT)
		}
		c.Jump(vm.JUMP_FORWARD, end)
		c.Label(except)
	}
	c.Op(vm.END_FINALLY)
	c.Label(orelse)
	c.Stmts(node.Orelse)
	c.Label(end)
	c.loops.Pop()
}

// Compile a try statement
func (c *compiler) try(node *ast.Try) {
	if len(node.Finalbody) > 0 {
		c.tryFinally(node)
	} else {
		c.tryExcept(node)
	}
}

/*
The IMPORT_NAME opcode was already generated.  This function

	merely needs to bind the result to a name.

	If there is a dot in name, we need to split it and emit a
	LOAD_ATTR for each name.
*/
func (c *compiler) importAs(name ast.Identifier, asname ast.Identifier) {
	attrs := strings.Split(string(name), ".")
	if len(attrs) > 1 {
		for _, attr := range attrs[1:] {
			c.OpArg(vm.LOAD_ATTR, c.Name(ast.Identifier(attr)))
		}
	}
	c.NameOp(string(asname), ast.Store)
}

/*
The Import node stores a module name like a.b.c as a single

	string.  This is convenient for all cases except
	  import a.b.c as d
	where we need to parse that string to extract the individual
	module names.
	XXX Perhaps change the representation to make this case simpler?
*/
func (c *compiler) import_(node *ast.Import) {
	//n = asdl_seq_LEN(s.v.Import.names);

	for _, alias := range node.Names {
		c.LoadConst(py.Int(0))
		c.LoadConst(py.None)
		c.OpName(vm.IMPORT_NAME, alias.Name)

		if alias.AsName != "" {
			c.importAs(alias.Name, alias.AsName)
		} else {
			tmp := alias.Name
			dot := strings.IndexByte(string(alias.Name), '.')
			if dot >= 0 {
				tmp = alias.Name[:dot]
			}
			c.NameOp(string(tmp), ast.Store)
		}
	}
}

func (c *compiler) importFrom(node *ast.ImportFrom) {
	names := make(py.Tuple, len(node.Names))

	/* build up the names */
	for i, alias := range node.Names {
		names[i] = py.String(alias.Name)
	}

	// FIXME if s.lineno > c.c_future.ff_lineno && s.v.ImportFrom.module && !PyUnicode_CompareWithASCIIString(s.v.ImportFrom.module, "__future__") {
	// 	return compiler_error(c, "from __future__ imports must occur at the beginning of the file")
	// }

	c.LoadConst(py.Int(node.Level))
	c.LoadConst(names)
	c.OpName(vm.IMPORT_NAME, node.Module)
	for i, alias := range node.Names {
		if i == 0 && alias.Name[0] == '*' {
			if len(alias.Name) != 1 {
				panic("can only import *")
			}
			c.Op(vm.IMPORT_STAR)
			return
		}

		c.OpName(vm.IMPORT_FROM, alias.Name)
		store_name := alias.Name
		if alias.AsName != "" {
			store_name = alias.AsName
		}

		c.NameOp(string(store_name), ast.Store)
	}
	/* remove imported module */
	c.Op(vm.POP_TOP)
}

// Compile statements
func (c *compiler) Stmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		c.Stmt(stmt)
	}
}

// Compile statement
func (c *compiler) Stmt(stmt ast.Stmt) {
	c.SetLineno(stmt)
	switch node := stmt.(type) {
	case *ast.FunctionDef:
		// Name          Identifier
		// Args          *Arguments
		// Body          []Stmt
		// DecoratorList []Expr
		// Returns       Expr
		c.compileFunc(compilerScopeFunction, stmt, node.Args, node.DecoratorList, node.Returns)
		c.NameOp(string(node.Name), ast.Store)

	case *ast.ClassDef:
		// Name          Identifier
		// Bases         []Expr
		// Keywords      []*Keyword
		// Starargs      Expr
		// Kwargs        Expr
		// Body          []Stmt
		// DecoratorList []Expr
		c.class(stmt, node)
	case *ast.Return:
		// Value Expr
		if c.SymTable.Type != symtable.FunctionBlock {
			c.panicSyntaxErrorf(node, "'return' outside function")
		}
		if node.Value != nil {
			c.Expr(node.Value)
		} else {
			c.LoadConst(py.None)
		}
		c.Op(vm.RETURN_VALUE)
	case *ast.Delete:
		// Targets []Expr
		c.Exprs(node.Targets)
	case *ast.Assign:
		// Targets []Expr
		// Value   Expr
		c.Expr(node.Value)
		for i, target := range node.Targets {
			if i != len(node.Targets)-1 {
				c.Op(vm.DUP_TOP)
			}
			c.Expr(target)
		}
	case *ast.AugAssign:
		// Target Expr
		// Op     OperatorNumber
		// Value  Expr
		setctx, ok := node.Target.(ast.SetCtxer)
		if !ok {
			panic("compile: can't set context in AugAssign")
		}
		// FIXME untidy modifying the ast temporarily!
		setctx.SetCtx(ast.AugLoad)
		c.Expr(node.Target)
		c.Expr(node.Value)
		var op vm.OpCode
		switch node.Op {
		case ast.Add:
			op = vm.INPLACE_ADD
		case ast.Sub:
			op = vm.INPLACE_SUBTRACT
		case ast.Mult:
			op = vm.INPLACE_MULTIPLY
		case ast.Div:
			op = vm.INPLACE_TRUE_DIVIDE
		case ast.Modulo:
			op = vm.INPLACE_MODULO
		case ast.Pow:
			op = vm.INPLACE_POWER
		case ast.LShift:
			op = vm.INPLACE_LSHIFT
		case ast.RShift:
			op = vm.INPLACE_RSHIFT
		case ast.BitOr:
			op = vm.INPLACE_OR
		case ast.BitXor:
			op = vm.INPLACE_XOR
		case ast.BitAnd:
			op = vm.INPLACE_AND
		case ast.FloorDiv:
			op = vm.INPLACE_FLOOR_DIVIDE
		default:
			panic("Unknown BinOp")
		}
		c.Op(op)
		setctx.SetCtx(ast.AugStore)
		c.Expr(node.Target)
	case *ast.For:
		// Target Expr
		// Iter   Expr
		// Body   []Stmt
		// Orelse []Stmt
		endfor := new(Label)
		endpopblock := new(Label)
		c.Jump(vm.SETUP_LOOP, endpopblock)
		c.Expr(node.Iter)
		c.Op(vm.GET_ITER)
		forloop := c.NewLabel()
		c.loops.Push(loop{Start: forloop, End: endpopblock, Type: loopLoop})
		c.Jump(vm.FOR_ITER, endfor)
		c.Expr(node.Target)
		c.Stmts(node.Body)
		c.Jump(vm.JUMP_ABSOLUTE, forloop)
		c.Label(endfor)
		c.Op(vm.POP_BLOCK)
		c.loops.Pop()
		c.Stmts(node.Orelse)
		c.Label(endpopblock)
	case *ast.While:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		endwhile := new(Label)
		endpopblock := new(Label)
		c.Jump(vm.SETUP_LOOP, endpopblock)
		while := c.NewLabel()
		c.loops.Push(loop{Start: while, End: endpopblock, Type: loopLoop})
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, endwhile)
		c.Stmts(node.Body)
		c.Jump(vm.JUMP_ABSOLUTE, while)
		c.Label(endwhile)
		c.Op(vm.POP_BLOCK)
		c.loops.Pop()
		c.Stmts(node.Orelse)
		c.Label(endpopblock)
	case *ast.If:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		orelse := new(Label)
		endif := new(Label)
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, orelse)
		c.Stmts(node.Body)
		// FIXME this puts a JUMP_FORWARD in when not
		// necessary (when no Orelse statements) but it
		// matches python3.4 (this is fixed in py3.5)
		c.Jump(vm.JUMP_FORWARD, endif)
		c.Label(orelse)
		c.Stmts(node.Orelse)
		c.Label(endif)
	case *ast.With:
		// Items []*WithItem
		// Body  []Stmt
		c.with(node, 0)
	case *ast.Raise:
		// Exc   Expr
		// Cause Expr
		args := uint32(0)
		if node.Exc != nil {
			args++
			c.Expr(node.Exc)
			if node.Cause != nil {
				args++
				c.Expr(node.Cause)
			}
		}
		c.OpArg(vm.RAISE_VARARGS, args)
	case *ast.Try:
		// Body      []Stmt
		// Handlers  []*ExceptHandler
		// Orelse    []Stmt
		// Finalbody []Stmt
		c.try(node)
	case *ast.Assert:
		// Test Expr
		// Msg  Expr
		label := new(Label)
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_TRUE, label)
		c.OpArg(vm.LOAD_GLOBAL, c.Name("AssertionError"))
		if node.Msg != nil {
			c.Expr(node.Msg)
			c.OpArg(vm.CALL_FUNCTION, 1) // 1 positional, 0 keyword pair
		}
		c.OpArg(vm.RAISE_VARARGS, 1)
		c.Label(label)
	case *ast.Import:
		// Names []*Alias
		c.import_(node)
	case *ast.ImportFrom:
		// Module Identifier
		// Names  []*Alias
		// Level  int
		c.importFrom(node)
	case *ast.Global:
		// Names []Identifier
		// Implemented by symtable
	case *ast.Nonlocal:
		// Names []Identifier
		// Implemented by symtable
	case *ast.ExprStmt:
		// Value Expr
		if c.interactive && c.depth <= 1 {
			c.Expr(node.Value)
			c.Op(vm.PRINT_EXPR)
		} else {
			switch node.Value.(type) {
			case *ast.Str:
			case *ast.Num:
			default:
				c.Expr(node.Value)
				c.Op(vm.POP_TOP)
			}
		}
	case *ast.Pass:
		// Do nothing
	case *ast.Break:
		l := c.loops.Top()
		if l == nil {
			c.panicSyntaxErrorf(node, "'break' outside loop")
		}
		c.Op(vm.BREAK_LOOP)
	case *ast.Continue:
		const loopError = "'continue' not properly in loop"
		const inFinallyError = "'continue' not supported inside 'finally' clause"
		l := c.loops.Top()
		if l == nil {
			c.panicSyntaxErrorf(node, loopError)
			panic("impossible")
		}
		switch l.Type {
		case loopLoop:
			c.Jump(vm.JUMP_ABSOLUTE, l.Start)
		case exceptLoop, finallyTryLoop:
			i := len(c.loops) - 2 // next loop out
			for ; i >= 0; i-- {
				l = &c.loops[i]
				if l.Type == loopLoop {
					break
				}
				// Prevent continue anywhere under a finally even if hidden in a sub-try or except.
				if l.Type == finallyEndLoop {
					c.panicSyntaxErrorf(node, inFinallyError)
				}
			}
			if i == -1 {
				c.panicSyntaxErrorf(node, loopError)
			}
			c.Jump(vm.CONTINUE_LOOP, l.Start)
		case finallyEndLoop:
			c.panicSyntaxErrorf(node, inFinallyError)
		default:
			panic("unknown loop type")
		}
	default:
		panic(fmt.Sprintf("Unknown StmtBase: %v", stmt))
	}
}

// Compile a NameOp
func (c *compiler) NameOp(name string, ctx ast.ExprContext) {
	// int op, scope;
	// Py_ssize_t arg;
	const (
		OP_FAST = iota
		OP_GLOBAL
		OP_DEREF
		OP_NAME
	)

	dict := &c.Code.Names
	// PyObject *mangled;
	/* XXX AugStore isn't used anywhere! */

	// FIXME mangled = _Py_Mangle(c.u.u_private, name);
	mangled := name

	if name == "None" || name == "True" || name == "False" {
		panic("NameOp: Can't compile None, True or False")
	}

	op := vm.OpCode(0)
	optype := OP_NAME
	scope := c.SymTable.GetScope(mangled)
	switch scope {
	case symtable.ScopeFree:
		dict = &c.Code.Freevars
		optype = OP_DEREF
	case symtable.ScopeCell:
		dict = &c.Code.Cellvars
		optype = OP_DEREF
	case symtable.ScopeLocal:
		if c.SymTable.Type == symtable.FunctionBlock {
			optype = OP_FAST
		}
	case symtable.ScopeGlobalImplicit:
		if c.SymTable.Type == symtable.FunctionBlock && c.SymTable.Unoptimized == 0 {
			optype = OP_GLOBAL
		}
	case symtable.ScopeGlobalExplicit:
		optype = OP_GLOBAL
	case symtable.ScopeInvalid:
		// scope can be unset
	default:
		panic(fmt.Sprintf("NameOp: Invalid scope %v for %q", scope, mangled))
	}

	/* XXX Leave assert here, but handle __doc__ and the like better */
	// FIXME assert(scope || PyUnicode_READ_CHAR(name, 0) == '_')

	switch optype {
	case OP_DEREF:
		switch ctx {
		case ast.Load, ast.AugLoad:
			if c.SymTable.Type == symtable.ClassBlock {
				op = vm.LOAD_CLASSDEREF
			} else {
				op = vm.LOAD_DEREF
			}
		case ast.Store, ast.AugStore:
			op = vm.STORE_DEREF
		case ast.Del:
			op = vm.DELETE_DEREF
		case ast.Param:
			panic("NameOp: param invalid for deref variable")
		default:
			panic("NameOp: ctx invalid for deref variable")
		}
	case OP_FAST:
		switch ctx {
		case ast.Load, ast.AugLoad:
			op = vm.LOAD_FAST
		case ast.Store, ast.AugStore:
			op = vm.STORE_FAST
		case ast.Del:
			op = vm.DELETE_FAST
		case ast.Param:
			panic("NameOp: param invalid for local variable")
		default:
			panic("NameOp: ctx invalid for local variable")
		}
		dict = &c.Code.Varnames
	case OP_GLOBAL:
		switch ctx {
		case ast.Load, ast.AugLoad:
			op = vm.LOAD_GLOBAL
		case ast.Store, ast.AugStore:
			op = vm.STORE_GLOBAL
		case ast.Del:
			op = vm.DELETE_GLOBAL
		case ast.Param:
			panic("NameOp: param invalid for global variable")
		default:
			panic("NameOp: ctx invalid for global variable")
		}
	case OP_NAME:
		switch ctx {
		case ast.Load, ast.AugLoad:
			op = vm.LOAD_NAME
		case ast.Store, ast.AugStore:
			op = vm.STORE_NAME
		case ast.Del:
			op = vm.DELETE_NAME
		case ast.Param:
			panic("NameOp: param invalid for name variable")
		default:
			panic("NameOp: ctx invalid for name variable")
		}
	}
	if op == 0 {
		panic("NameOp: Op not set")
	}
	i := c.Index(mangled, dict)
	// using CellAndFreeVars in closures requires skipping Cellvars
	if scope == symtable.ScopeFree {
		i += uint32(len(c.Code.Cellvars))
	}
	c.OpArg(op, i)
}

// Call a function which is already on the stack with n arguments already on the stack
func (c *compiler) callHelper(n int, Args []ast.Expr, Keywords []*ast.Keyword, Starargs ast.Expr, Kwargs ast.Expr) {
	args := len(Args) + n
	for i := range Args {
		c.Expr(Args[i])
	}
	kwargs := len(Keywords)
	duplicateDetector := make(map[ast.Identifier]struct{}, len(Keywords))
	var duplicate *ast.Keyword
	for i := range Keywords {
		kw := Keywords[i]
		if _, found := duplicateDetector[kw.Arg]; found {
			if duplicate == nil {
				duplicate = kw
			}
		} else {
			duplicateDetector[kw.Arg] = struct{}{}
		}
		c.LoadConst(py.String(kw.Arg))
		c.Expr(kw.Value)
	}
	if duplicate != nil {
		c.panicSyntaxErrorf(duplicate, "keyword argument repeated")
	}
	op := vm.CALL_FUNCTION
	if Starargs != nil {
		c.Expr(Starargs)
		if Kwargs != nil {
			c.Expr(Kwargs)
			op = vm.CALL_FUNCTION_VAR_KW
		} else {
			op = vm.CALL_FUNCTION_VAR
		}
	} else if Kwargs != nil {
		c.Expr(Kwargs)
		op = vm.CALL_FUNCTION_KW
	}
	c.OpArg(op, uint32(args+kwargs<<8))
}

/*
	List and set comprehensions and generator expressions work by creating a

nested function to perform the actual iteration. This means that the
iteration variables don't leak into the current scope.
The defined function is called immediately following its definition, with the
result of that call being the result of the expression.
The LC/SC version returns the populated container, while the GE version is
flagged in symtable.c as a generator, so it returns the generator object
when the function is called.
This code *knows* that the loop cannot contain break, continue, or return,
so it cheats and skips the SETUP_LOOP/POP_BLOCK steps used in normal loops.

Possible cleanups:
  - iterate over the generator sequence instead of using recursion
*/
func (c *compiler) comprehensionGenerator(generators []ast.Comprehension, gen_index int, elt ast.Expr, val ast.Expr, Ast ast.Ast) {
	// generate code for the iterator, then each of the ifs,
	// and then write to the element
	start := new(Label)
	skip := new(Label)
	anchor := new(Label)
	gen := generators[gen_index]
	if gen_index == 0 {
		/* Receive outermost iter as an implicit argument */
		c.Code.Argcount = 1
		c.OpArg(vm.LOAD_FAST, 0)
	} else {
		/* Sub-iter - calculate on the fly */
		c.Expr(gen.Iter)
		c.Op(vm.GET_ITER)
	}
	c.Label(start)
	c.Jump(vm.FOR_ITER, anchor)
	c.Expr(gen.Target)

	/* XXX this needs to be cleaned up...a lot! */
	for _, e := range gen.Ifs {
		c.Expr(e)
		c.Jump(vm.POP_JUMP_IF_FALSE, start)
	}

	gen_index++
	if gen_index < len(generators) {
		c.comprehensionGenerator(generators, gen_index, elt, val, Ast)
	}

	/* only append after the last for generator */
	if gen_index >= len(generators) {
		/* comprehension specific code */
		switch Ast.(type) {
		case *ast.GeneratorExp:
			c.Expr(elt)
			c.Op(vm.YIELD_VALUE)
			c.Op(vm.POP_TOP)
		case *ast.ListComp:
			c.Expr(elt)
			c.OpArg(vm.LIST_APPEND, uint32(gen_index+1))
		case *ast.SetComp:
			c.Expr(elt)
			c.OpArg(vm.SET_ADD, uint32(gen_index+1))
		case *ast.DictComp:
			// With 'd[k] = v', v is evaluated before k, so we do the same.
			c.Expr(val)
			c.Expr(elt)
			c.OpArg(vm.MAP_ADD, uint32(gen_index+1))
		default:
			panic(fmt.Sprintf("unknown comprehension %v", Ast))
		}
		c.Label(skip)
	}
	c.Jump(vm.JUMP_ABSOLUTE, start)
	c.Label(anchor)
}

// Compile a comprehension
func (c *compiler) comprehension(expr ast.Expr, generators []ast.Comprehension) {
	newC := c.newCompilerScope(compilerScopeComprehension, expr, "")
	c.makeClosure(newC.Code, 0, newC, newC.Code.Name)
	outermost_iter := generators[0].Iter
	c.Expr(outermost_iter)
	c.Op(vm.GET_ITER)
	c.OpArg(vm.CALL_FUNCTION, 1)
}

// Compile a tuple or a list
func (c *compiler) tupleOrList(op vm.OpCode, ctx ast.ExprContext, elts []ast.Expr) {
	const INT_MAX = 0x7FFFFFFF
	n := len(elts)
	if ctx == ast.Store {
		seen_star := false
		for i, elt := range elts {
			starred, isStarred := elt.(*ast.Starred)
			if isStarred && !seen_star {
				if i >= (1<<8) || n-i-1 >= (INT_MAX>>8) {
					c.panicSyntaxErrorf(elt, "too many expressions in star-unpacking assignment")
				}
				c.OpArg(vm.UNPACK_EX, uint32((i + ((n - i - 1) << 8))))
				seen_star = true
				// FIXME Overwrite the starred element
				elts[i] = starred.Value
			} else if isStarred {
				c.panicSyntaxErrorf(elt, "two starred expressions in assignment")
			}
		}
		if !seen_star {
			c.OpArg(vm.UNPACK_SEQUENCE, uint32(n))
		}
	}
	c.Exprs(elts)
	if ctx == ast.Load {
		c.OpArg(op, uint32(n))
	}
}

// compile a subscript
func (c *compiler) subscript(kind string, ctx ast.ExprContext) {
	switch ctx {
	case ast.AugLoad:
		c.Op(vm.DUP_TOP_TWO)
		c.Op(vm.BINARY_SUBSCR)
	case ast.Load:
		c.Op(vm.BINARY_SUBSCR)
	case ast.AugStore:
		c.Op(vm.ROT_THREE)
		c.Op(vm.STORE_SUBSCR)
	case ast.Store:
		c.Op(vm.STORE_SUBSCR)
	case ast.Del:
		c.Op(vm.DELETE_SUBSCR)
	case ast.Param:
		panic(fmt.Sprintf("invalid %v kind %v in subscript", kind, ctx))
	}
}

// build the slice
func (c *compiler) buildSlice(slice *ast.Slice, ctx ast.ExprContext) {
	n := uint32(2)

	/* only handles the cases where BUILD_SLICE is emitted */
	if slice.Lower != nil {
		c.Expr(slice.Lower)
	} else {
		c.LoadConst(py.None)
	}

	if slice.Upper != nil {
		c.Expr(slice.Upper)
	} else {
		c.LoadConst(py.None)
	}

	if slice.Step != nil {
		n++
		c.Expr(slice.Step)
	}
	c.OpArg(vm.BUILD_SLICE, n)
}

// compile a nested slice
func (c *compiler) nestedSlice(s ast.Slicer, ctx ast.ExprContext) {
	switch node := s.(type) {
	case *ast.Slice:
		c.buildSlice(node, ctx)
	case *ast.Index:
		c.Expr(node.Value)
	case *ast.ExtSlice:
		panic("extended slice invalid in nested slice")
	default:
		panic("nestedSlice: unknown type")
	}
}

// Compile a slice
func (c *compiler) slice(s ast.Slicer, ctx ast.ExprContext) {
	kindname := ""
	switch node := s.(type) {
	case *ast.Index:
		kindname = "index"
		if ctx != ast.AugStore {
			c.Expr(node.Value)
		}
	case *ast.Slice:
		kindname = "slice"
		if ctx != ast.AugStore {
			c.buildSlice(node, ctx)
		}
	case *ast.ExtSlice:
		kindname = "extended slice"
		if ctx != ast.AugStore {
			for _, sub := range node.Dims {
				c.nestedSlice(sub, ctx)
			}
			c.OpArg(vm.BUILD_TUPLE, uint32(len(node.Dims)))
		}
	default:
		panic(fmt.Sprintf("invalid subscript kind %T", s))
	}
	c.subscript(kindname, ctx)
}

// Compile expressions
func (c *compiler) Exprs(exprs []ast.Expr) {
	for _, expr := range exprs {
		c.Expr(expr)
	}
}

// Compile and expression
func (c *compiler) Expr(expr ast.Expr) {
	c.SetLineno(expr)
	switch node := expr.(type) {
	case *ast.BoolOp:
		// Op     BoolOpNumber
		// Values []Expr
		var op vm.OpCode
		switch node.Op {
		case ast.And:
			op = vm.JUMP_IF_FALSE_OR_POP
		case ast.Or:
			op = vm.JUMP_IF_TRUE_OR_POP
		default:
			panic("Unknown BoolOp")
		}
		label := new(Label)
		for i, e := range node.Values {
			c.Expr(e)
			if i != len(node.Values)-1 {
				c.Jump(op, label)
			}
		}
		c.Label(label)
	case *ast.BinOp:
		// Left  Expr
		// Op    OperatorNumber
		// Right Expr
		c.Expr(node.Left)
		c.Expr(node.Right)
		var op vm.OpCode
		switch node.Op {
		case ast.Add:
			op = vm.BINARY_ADD
		case ast.Sub:
			op = vm.BINARY_SUBTRACT
		case ast.Mult:
			op = vm.BINARY_MULTIPLY
		case ast.Div:
			op = vm.BINARY_TRUE_DIVIDE
		case ast.Modulo:
			op = vm.BINARY_MODULO
		case ast.Pow:
			op = vm.BINARY_POWER
		case ast.LShift:
			op = vm.BINARY_LSHIFT
		case ast.RShift:
			op = vm.BINARY_RSHIFT
		case ast.BitOr:
			op = vm.BINARY_OR
		case ast.BitXor:
			op = vm.BINARY_XOR
		case ast.BitAnd:
			op = vm.BINARY_AND
		case ast.FloorDiv:
			op = vm.BINARY_FLOOR_DIVIDE
		default:
			panic("Unknown BinOp")
		}
		c.Op(op)
	case *ast.UnaryOp:
		// Op      UnaryOpNumber
		// Operand Expr
		c.Expr(node.Operand)
		var op vm.OpCode
		switch node.Op {
		case ast.Invert:
			op = vm.UNARY_INVERT
		case ast.Not:
			op = vm.UNARY_NOT
		case ast.UAdd:
			op = vm.UNARY_POSITIVE
		case ast.USub:
			op = vm.UNARY_NEGATIVE
		default:
			panic("Unknown UnaryOp")
		}
		c.Op(op)
	case *ast.Lambda:
		// Args *Arguments
		// Body Expr
		// newC := Compiler
		c.compileFunc(compilerScopeLambda, expr, node.Args, nil, nil)
	case *ast.IfExp:
		// Test   Expr
		// Body   Expr
		// Orelse Expr
		elseBranch := new(Label)
		endifBranch := new(Label)
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, elseBranch)
		c.Expr(node.Body)
		c.Jump(vm.JUMP_FORWARD, endifBranch)
		c.Label(elseBranch)
		c.Expr(node.Orelse)
		c.Label(endifBranch)
	case *ast.Dict:
		// Keys   []Expr
		// Values []Expr
		n := len(node.Keys)
		if n != len(node.Values) {
			panic("compile: Dict keys and values differing sizes")
		}
		c.OpArg(vm.BUILD_MAP, uint32(n))
		for i := range node.Keys {
			c.Expr(node.Values[i])
			c.Expr(node.Keys[i])
			c.Op(vm.STORE_MAP)
		}
	case *ast.Set:
		// Elts []Expr
		c.Exprs(node.Elts)
		c.OpArg(vm.BUILD_SET, uint32(len(node.Elts)))
	case *ast.ListComp:
		// Elt        Expr
		// Generators []Comprehension
		c.comprehension(expr, node.Generators)
	case *ast.SetComp:
		// Elt        Expr
		// Generators []Comprehension
		c.comprehension(expr, node.Generators)
	case *ast.DictComp:
		// Key        Expr
		// Value      Expr
		// Generators []Comprehension
		c.comprehension(expr, node.Generators)
	case *ast.GeneratorExp:
		// Elt        Expr
		// Generators []Comprehension
		c.comprehension(expr, node.Generators)
	case *ast.Yield:
		// Value Expr
		if c.SymTable.Type != symtable.FunctionBlock {
			c.panicSyntaxErrorf(node, "'yield' outside function")
		}
		if node.Value != nil {
			c.Expr(node.Value)
		} else {
			c.LoadConst(py.None)
		}
		c.Op(vm.YIELD_VALUE)
	case *ast.YieldFrom:
		// Value Expr
		if c.SymTable.Type != symtable.FunctionBlock {
			c.panicSyntaxErrorf(node, "'yield' outside function")
		}
		c.Expr(node.Value)
		c.Op(vm.GET_ITER)
		c.LoadConst(py.None)
		c.Op(vm.YIELD_FROM)
	case *ast.Compare:
		// Left        Expr
		// Ops         []CmpOp
		// Comparators []Expr
		if len(node.Ops) != len(node.Comparators) {
			panic("compile: Unequal Ops and Comparators in Compare")
		}
		if len(node.Ops) == 0 {
			panic("compile: No Ops or Comparators in Compare")
		}
		c.Expr(node.Left)
		label := new(Label)
		for i := range node.Ops {
			last := i == len(node.Ops)-1
			c.Expr(node.Comparators[i])
			if !last {
				c.Op(vm.DUP_TOP)
				c.Op(vm.ROT_THREE)
			}
			op := node.Ops[i]
			var arg uint32
			switch op {
			case ast.Eq:
				arg = vm.PyCmp_EQ
			case ast.NotEq:
				arg = vm.PyCmp_NE
			case ast.Lt:
				arg = vm.PyCmp_LT
			case ast.LtE:
				arg = vm.PyCmp_LE
			case ast.Gt:
				arg = vm.PyCmp_GT
			case ast.GtE:
				arg = vm.PyCmp_GE
			case ast.Is:
				arg = vm.PyCmp_IS
			case ast.IsNot:
				arg = vm.PyCmp_IS_NOT
			case ast.In:
				arg = vm.PyCmp_IN
			case ast.NotIn:
				arg = vm.PyCmp_NOT_IN
			default:
				panic("compile: Unknown OpArg")
			}
			c.OpArg(vm.COMPARE_OP, arg)
			if !last {
				c.Jump(vm.JUMP_IF_FALSE_OR_POP, label)
			}
		}
		if len(node.Ops) > 1 {
			endLabel := new(Label)
			c.Jump(vm.JUMP_FORWARD, endLabel)
			c.Label(label)
			c.Op(vm.ROT_TWO)
			c.Op(vm.POP_TOP)
			c.Label(endLabel)
		}
	case *ast.Call:
		// Func     Expr
		// Args     []Expr
		// Keywords []*Keyword
		// Starargs Expr
		// Kwargs   Expr
		c.Expr(node.Func)
		c.callHelper(0, node.Args, node.Keywords, node.Starargs, node.Kwargs)
	case *ast.Num:
		// N Object
		c.LoadConst(node.N)
	case *ast.Str:
		// S py.String
		c.LoadConst(node.S)
	case *ast.Bytes:
		// S py.Bytes
		c.LoadConst(node.S)
	case *ast.NameConstant:
		// Value Singleton
		c.LoadConst(node.Value)
	case *ast.Ellipsis:
		c.LoadConst(py.Ellipsis)
	case *ast.Attribute:
		// Value Expr
		// Attr  Identifier
		// Ctx   ExprContext
		if node.Ctx != ast.AugStore {
			c.Expr(node.Value)
		}
		var op vm.OpCode
		switch node.Ctx {
		case ast.AugLoad:
			c.Op(vm.DUP_TOP)
			op = vm.LOAD_ATTR
		case ast.Load:
			op = vm.LOAD_ATTR
		case ast.AugStore:
			c.Op(vm.ROT_TWO)
			op = vm.STORE_ATTR
		case ast.Store:
			op = vm.STORE_ATTR
		case ast.Del:
			op = vm.DELETE_ATTR
		case ast.Param:
			panic("param invalid in attribute expression")
		default:
			panic("unknown context in attribute expression")
		}
		c.OpArg(op, c.Name(node.Attr))
	case *ast.Subscript:
		// Value Expr
		// Slice Slicer
		// Ctx   ExprContext
		switch node.Ctx {
		case ast.AugLoad:
			c.Expr(node.Value)
			c.slice(node.Slice, ast.AugLoad)
		case ast.Load:
			c.Expr(node.Value)
			c.slice(node.Slice, ast.Load)
		case ast.AugStore:
			c.slice(node.Slice, ast.AugStore)
		case ast.Store:
			c.Expr(node.Value)
			c.slice(node.Slice, ast.Store)
		case ast.Del:
			c.Expr(node.Value)
			c.slice(node.Slice, ast.Del)
		default:
			panic("param invalid in subscript expression")
		}
	case *ast.Starred:
		// Value Expr
		// Ctx   ExprContext
		switch node.Ctx {
		case ast.Store:
			// In all legitimate cases, the Starred node was already replaced
			// by tupleOrList: is that okay?
			c.panicSyntaxErrorf(node, "starred assignment target must be in a list or tuple")
		default:
			c.panicSyntaxErrorf(node, "can use starred expression only as assignment target")
		}
	case *ast.Name:
		// Id  Identifier
		// Ctx ExprContext
		c.NameOp(string(node.Id), node.Ctx)
	case *ast.List:
		// Elts []Expr
		// Ctx  ExprContext
		c.tupleOrList(vm.BUILD_LIST, node.Ctx, node.Elts)
	case *ast.Tuple:
		// Elts []Expr
		// Ctx  ExprContext
		c.tupleOrList(vm.BUILD_TUPLE, node.Ctx, node.Elts)
	default:
		panic(fmt.Sprintf("Unknown ExprBase: %v", expr))
	}
}
