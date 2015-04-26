// compile python code
//
// Need to port the 10,000 lines of compiling machinery, into a
// different module probably.
//
// In the mean time, cheat horrendously by calling python3.4 to do our
// dirty work under the hood!

package compile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ncw/gpython/ast"
	"github.com/ncw/gpython/marshal"
	"github.com/ncw/gpython/parser"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/symtable"
	"github.com/ncw/gpython/vm"
)

// Set in py to avoid circular import
func init() {
	py.Compile = Compile
}

// Compile(source, filename, mode, flags, dont_inherit) -> code object
//
// Compile the source string (a Python module, statement or expression)
// into a code object that can be executed by exec() or eval().
// The filename will be used for run-time error messages.
// The mode must be 'exec' to compile a module, 'single' to compile a
// single (interactive) statement, or 'eval' to compile an expression.
// The flags argument, if present, controls which future statements influence
// the compilation of the code.
// The dont_inherit argument, if non-zero, stops the compilation inheriting
// the effects of any future statements in effect in the code calling
// compile; if absent or zero these statements do influence the compilation,
// in addition to any features explicitly specified.
func LegacyCompile(str, filename, mode string, flags int, dont_inherit bool) py.Object {
	dont_inherit_str := "False"
	if dont_inherit {
		dont_inherit_str = "True"
	}
	// FIXME escaping in filename
	code := fmt.Sprintf(`import sys, marshal
str = sys.stdin.buffer.read().decode("utf-8")
code = compile(str, "%s", "%s", %d, %s)
marshalled_code = marshal.dumps(code)
sys.stdout.buffer.write(marshalled_code)
sys.stdout.close()`,
		filename,
		mode,
		flags,
		dont_inherit_str,
	)
	cmd := exec.Command("python3.4", "-c", code)
	cmd.Stdin = strings.NewReader(str)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "--- Failed to run python3.4 compile ---\n")
		fmt.Fprintf(os.Stderr, "--------------------\n")
		os.Stderr.Write(stderr.Bytes())
		fmt.Fprintf(os.Stderr, "--------------------\n")
		panic(err)
	}
	obj, err := marshal.ReadObject(bytes.NewBuffer(out.Bytes()))
	if err != nil {
		panic(err)
	}
	return obj
}

// Compile(source, filename, mode, flags, dont_inherit) -> code object
//
// Compile the source string (a Python module, statement or expression)
// into a code object that can be executed by exec() or eval().
// The filename will be used for run-time error messages.
// The mode must be 'exec' to compile a module, 'single' to compile a
// single (interactive) statement, or 'eval' to compile an expression.
// The flags argument, if present, controls which future statements influence
// the compilation of the code.
// The dont_inherit argument, if non-zero, stops the compilation inheriting
// the effects of any future statements in effect in the code calling
// compile; if absent or zero these statements do influence the compilation,
// in addition to any features explicitly specified.
func Compile(str, filename, mode string, flags int, dont_inherit bool) (py.Object, error) {
	// Parse Ast
	Ast, err := parser.ParseString(str, mode)
	if err != nil {
		return nil, err
	}
	// Make symbol table
	SymTable, err := symtable.NewSymTable(Ast)
	if err != nil {
		return nil, err
	}
	return CompileAst(Ast, filename, flags, dont_inherit, SymTable)
}

// As Compile but takes an Ast
func CompileAst(Ast ast.Ast, filename string, flags int, dont_inherit bool, SymTable *symtable.SymTable) (code *py.Code, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = py.MakeException(r)
		}
	}()
	//fmt.Println(ast.Dump(Ast))
	code = &py.Code{
		Filename:    filename,
		Firstlineno: 1,                           // FIXME
		Name:        "<module>",                  // FIXME
		Flags:       int32(flags | py.CO_NOFREE), // FIXME
	}
	c := &compiler{
		Code:     code,
		SymTable: SymTable,
	}
	valueOnStack := false
	switch node := Ast.(type) {
	case *ast.Module:
		c.Stmts(node.Body)
	case *ast.Interactive:
		c.Stmts(node.Body)
	case *ast.Expression:
		c.Expr(node.Body)
		valueOnStack = true
	case *ast.Suite:
		c.Stmts(node.Body)
	case ast.Expr:
		// Make None the first constant as lambda can't have a docstring
		c.Const(py.None)
		c.Code.Name = "<lambda>"
		c.Expr(node)
		valueOnStack = true
	case *ast.FunctionDef:
		c.Stmts(c.docString(node.Body))
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown ModuleBase: %v", Ast))
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
	return code, nil
}

// Loop
type loop struct {
	Start     *Label
	End       *Label
	IsForLoop bool
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

// State for the compiler
type compiler struct {
	Code     *py.Code // code being built up
	OpCodes  Instructions
	loops    loopstack
	SymTable *symtable.SymTable
}

// Check for docstring as first Expr in body and remove it and set the
// first constant if found.
func (c *compiler) docString(body []ast.Stmt) []ast.Stmt {
	if len(body) > 0 {
		if expr, ok := body[0].(*ast.ExprStmt); ok {
			if docstring, ok := expr.Value.(*ast.Str); ok {
				c.Const(docstring.S)
				return body[1:]
			}
		}
	}
	// If no docstring put None in
	c.Const(py.None)
	return body
}

// Compiles a python constant
//
// Returns the index into the Consts tuple
func (c *compiler) Const(obj py.Object) uint32 {
	// FIXME back this with a dict to stop O(N**2) behaviour on lots of consts
	for i, c := range c.Code.Consts {
		if obj.Type() == c.Type() && py.Eq(obj, c) == py.True {
			return uint32(i)
		}
	}
	c.Code.Consts = append(c.Code.Consts, obj)
	return uint32(len(c.Code.Consts) - 1)
}

// Loads a constant
func (c *compiler) LoadConst(obj py.Object) {
	c.OpArg(vm.LOAD_CONST, c.Const(obj))
}

// Returns the index into the slice provided, updating the slice if necessary
func (c *compiler) Index(Id string, Names *[]string) uint32 {
	// FIXME back this with a dict to stop O(N**2) behaviour on lots of vars
	for i, s := range *Names {
		if Id == s {
			return uint32(i)
		}
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

// Compiles an instruction with an argument
func (c *compiler) OpArg(Op byte, Arg uint32) {
	if !vm.HAS_ARG(Op) {
		panic("OpArg called with an instruction which doesn't take an Arg")
	}
	c.OpCodes.Add(&OpArg{Op: Op, Arg: Arg})
}

// Compiles an instruction without an argument
func (c *compiler) Op(op byte) {
	if vm.HAS_ARG(op) {
		panic("Op called with an instruction which takes an Arg")
	}
	c.OpCodes.Add(&Op{Op: op})
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
func (c *compiler) Jump(Op byte, Dest *Label) {
	switch Op {
	case vm.JUMP_IF_FALSE_OR_POP, vm.JUMP_IF_TRUE_OR_POP, vm.JUMP_ABSOLUTE, vm.POP_JUMP_IF_FALSE, vm.POP_JUMP_IF_TRUE, vm.CONTINUE_LOOP: // Absolute
		c.OpCodes.Add(&JumpAbs{OpArg: OpArg{Op: Op}, Dest: Dest})
	case vm.JUMP_FORWARD, vm.SETUP_WITH, vm.FOR_ITER, vm.SETUP_LOOP, vm.SETUP_EXCEPT, vm.SETUP_FINALLY:
		c.OpCodes.Add(&JumpRel{OpArg: OpArg{Op: Op}, Dest: Dest})
	default:
		panic("Jump called with non jump instruction")
	}
}

// Compile statements
func (c *compiler) Stmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		c.Stmt(stmt)
	}
}

// Compile statement
func (c *compiler) Stmt(stmt ast.Stmt) {
	switch node := stmt.(type) {
	case *ast.FunctionDef:
		// Name          Identifier
		// Args          *Arguments
		// Body          []Stmt
		// DecoratorList []Expr
		// Returns       Expr
		newSymTable := c.SymTable.FindChild(stmt)
		if newSymTable == nil {
			panic("No symtable found for function")
		}
		code, err := CompileAst(node, c.Code.Filename, int(c.Code.Flags)|py.CO_OPTIMIZED|py.CO_NEWLOCALS, false, newSymTable) // FIXME pass on compile args
		if err != nil {
			panic(err)
		}
		code.Argcount = int32(len(node.Args.Args))
		code.Name = string(node.Name)
		code.Kwonlyargcount = int32(len(node.Args.Kwonlyargs))
		code.Nlocals = code.Kwonlyargcount + int32(len(node.Args.Args))
		if code.Kwonlyargcount > 0 {
			code.Flags |= py.CO_VARARGS
		}

		// Arguments
		for _, arg := range node.Args.Args {
			c.Index(string(arg.Arg), &code.Varnames)
		}
		for _, arg := range node.Args.Kwonlyargs {
			c.Index(string(arg.Arg), &code.Varnames)
		}
		if node.Args.Vararg != nil {
			code.Nlocals++
			c.Index(string(node.Args.Vararg.Arg), &code.Varnames)
		}
		if node.Args.Kwarg != nil {
			code.Nlocals++
			c.Index(string(node.Args.Kwarg.Arg), &code.Varnames)
			code.Flags |= py.CO_VARKEYWORDS
		}

		// Defaults
		posdefaults := uint32(len(node.Args.Defaults))
		for _, expr := range node.Args.Defaults {
			c.Expr(expr)
		}

		// KwDefaults
		if len(node.Args.Kwonlyargs) != len(node.Args.KwDefaults) {
			panic("differing number of Kwonlyargs to KwDefaults")
		}
		kwdefaults := uint32(len(node.Args.KwDefaults))
		for i := range node.Args.KwDefaults {
			c.LoadConst(py.String(node.Args.Kwonlyargs[i].Arg))
			c.Expr(node.Args.KwDefaults[i])
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
		addAnnotation(node.Args.Args...)
		addAnnotation(node.Args.Vararg)
		addAnnotation(node.Args.Kwonlyargs...)
		addAnnotation(node.Args.Kwarg)
		if node.Returns != nil {
			c.Expr(node.Returns)
			annotations = append(annotations, py.String("return"))
		}
		num_annotations := uint32(len(annotations))
		if num_annotations > 0 {
			num_annotations++ // include the tuple
			c.LoadConst(annotations)
		}

		c.LoadConst(code)
		c.LoadConst(py.String(node.Name))
		c.OpArg(vm.MAKE_FUNCTION, posdefaults+(kwdefaults<<8)+(num_annotations<<16))
		c.OpArg(vm.STORE_NAME, c.Name(node.Name))
	case *ast.ClassDef:
		// Name          Identifier
		// Bases         []Expr
		// Keywords      []*Keyword
		// Starargs      Expr
		// Kwargs        Expr
		// Body          []Stmt
		// DecoratorList []Expr
		panic("FIXME compile: ClassDef not implemented")
	case *ast.Return:
		// Value Expr
		if node.Value != nil {
			c.Expr(node.Value)
		} else {
			c.LoadConst(py.None)
		}
		c.Op(vm.RETURN_VALUE)
	case *ast.Delete:
		// Targets []Expr
		for _, target := range node.Targets {
			c.Expr(target)
		}
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
		setctx.SetCtx(ast.Load)
		c.Expr(node.Target)
		c.Expr(node.Value)
		var op byte
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
		setctx.SetCtx(ast.Store)
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
		c.loops.Push(loop{Start: forloop, End: endpopblock, IsForLoop: true})
		c.Jump(vm.FOR_ITER, endfor)
		c.Expr(node.Target)
		for _, stmt := range node.Body {
			c.Stmt(stmt)
		}
		c.Jump(vm.JUMP_ABSOLUTE, forloop)
		c.Label(endfor)
		c.Op(vm.POP_BLOCK)
		c.loops.Pop()
		for _, stmt := range node.Orelse {
			c.Stmt(stmt)
		}
		c.Label(endpopblock)
	case *ast.While:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		endwhile := new(Label)
		endpopblock := new(Label)
		c.Jump(vm.SETUP_LOOP, endpopblock)
		while := c.NewLabel()
		c.loops.Push(loop{Start: while, End: endpopblock})
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, endwhile)
		for _, stmt := range node.Body {
			c.Stmt(stmt)
		}
		c.Jump(vm.JUMP_ABSOLUTE, while)
		c.Label(endwhile)
		c.Op(vm.POP_BLOCK)
		c.loops.Pop()
		for _, stmt := range node.Orelse {
			c.Stmt(stmt)
		}
		c.Label(endpopblock)
	case *ast.If:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		orelse := new(Label)
		endif := new(Label)
		c.Expr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, orelse)
		for _, stmt := range node.Body {
			c.Stmt(stmt)
		}
		// FIXME this puts a JUMP_FORWARD in when not
		// necessary (when no Orelse statements) but it
		// matches python3.4 (this is fixed in py3.5)
		c.Jump(vm.JUMP_FORWARD, endif)
		c.Label(orelse)
		for _, stmt := range node.Orelse {
			c.Stmt(stmt)
		}
		c.Label(endif)
	case *ast.With:
		// Items []*WithItem
		// Body  []Stmt
		panic("FIXME compile: With not implemented")
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
		panic("FIXME compile: Try not implemented")
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
		panic("FIXME compile: Import not implemented")
	case *ast.ImportFrom:
		// Module Identifier
		// Names  []*Alias
		// Level  int
		panic("FIXME compile: ImportFrom not implemented")
	case *ast.Global:
		// Names []Identifier
		// Implemented by symtable
	case *ast.Nonlocal:
		// Names []Identifier
		// Implemented by symtable
	case *ast.ExprStmt:
		// Value Expr
		c.Expr(node.Value)
		c.Op(vm.POP_TOP)
	case *ast.Pass:
		// Do nothing
	case *ast.Break:
		l := c.loops.Top()
		if l == nil {
			panic(py.ExceptionNewf(py.SyntaxError, "'break' outside loop"))
		}
		c.Op(vm.BREAK_LOOP)
	case *ast.Continue:
		l := c.loops.Top()
		if l == nil {
			panic(py.ExceptionNewf(py.SyntaxError, "'continue' not properly in loop"))
		}
		if l.IsForLoop {
			// FIXME when do we use CONTINUE_LOOP?
			c.Jump(vm.JUMP_ABSOLUTE, l.Start)
			//c.Jump(vm.CONTINUE_LOOP, l.Start)
		} else {
			c.Jump(vm.JUMP_ABSOLUTE, l.Start)
		}
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown StmtBase: %v", stmt))
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

	// FIXME mangled = _Py_Mangle(c->u->u_private, name);
	mangled := name

	if name == "None" || name == "True" || name == "False" {
		panic("NameOp: Can't compile None, True or False")
	}

	op := byte(0)
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
	default:
		panic(fmt.Sprintf("NameOp: Invalid scope %v for %q", scope, mangled))
	}

	/* XXX Leave assert here, but handle __doc__ and the like better */
	// FIXME assert(scope || PyUnicode_READ_CHAR(name, 0) == '_')

	switch optype {
	case OP_DEREF:
		switch ctx {
		case ast.Load:
			if c.SymTable.Type == symtable.ClassBlock {
				op = vm.LOAD_CLASSDEREF
			} else {
				op = vm.LOAD_DEREF
			}
		case ast.Store:
			op = vm.STORE_DEREF
		case ast.AugLoad:
		case ast.AugStore:
		case ast.Del:
			op = vm.DELETE_DEREF
		case ast.Param:
			panic("NameOp: param invalid for deref variable")
		default:
			panic("NameOp: ctx invalid for deref variable")
		}
	case OP_FAST:
		switch ctx {
		case ast.Load:
			op = vm.LOAD_FAST
		case ast.Store:
			op = vm.STORE_FAST
		case ast.Del:
			op = vm.DELETE_FAST
		case ast.AugLoad:
		case ast.AugStore:
		case ast.Param:
			panic("NameOp: param invalid for local variable")
		default:
			panic("NameOp: ctx invalid for local variable")
		}
		dict = &c.Code.Varnames
	case OP_GLOBAL:
		switch ctx {
		case ast.Load:
			op = vm.LOAD_GLOBAL
		case ast.Store:
			op = vm.STORE_GLOBAL
		case ast.Del:
			op = vm.DELETE_GLOBAL
		case ast.AugLoad:
		case ast.AugStore:
		case ast.Param:
			panic("NameOp: param invalid for global variable")
		default:
			panic("NameOp: ctx invalid for global variable")
		}
	case OP_NAME:
		switch ctx {
		case ast.Load:
			op = vm.LOAD_NAME
		case ast.Store:
			op = vm.STORE_NAME
		case ast.Del:
			op = vm.DELETE_NAME
		case ast.AugLoad:
		case ast.AugStore:
		case ast.Param:
			panic("NameOp: param invalid for name variable")
		default:
			panic("NameOp: ctx invalid for name variable")
		}
		break
	}
	if op == 0 {
		panic("NameOp: Op not set")
	}
	c.OpArg(op, c.Index(mangled, dict))
}

// Compile expressions
func (c *compiler) Expr(expr ast.Expr) {
	switch node := expr.(type) {
	case *ast.BoolOp:
		// Op     BoolOpNumber
		// Values []Expr
		var op byte
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
		var op byte
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
		var op byte
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
		newSymTable := c.SymTable.FindChild(expr)
		if newSymTable == nil {
			panic("No symtable found for lambda")
		}
		code, err := CompileAst(node.Body, c.Code.Filename, int(c.Code.Flags)|py.CO_OPTIMIZED|py.CO_NEWLOCALS, false, newSymTable) // FIXME pass on compile args
		if err != nil {
			panic(err)
		}

		code.Argcount = int32(len(node.Args.Args))
		c.LoadConst(code)
		c.LoadConst(py.String("<lambda>"))
		// FIXME node.Args
		c.OpArg(vm.MAKE_FUNCTION, 0)
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
		for _, elt := range node.Elts {
			c.Expr(elt)
		}
		c.OpArg(vm.BUILD_SET, uint32(len(node.Elts)))
	case *ast.ListComp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME compile: ListComp not implemented")
	case *ast.SetComp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME compile: SetComp not implemented")
	case *ast.DictComp:
		// Key        Expr
		// Value      Expr
		// Generators []Comprehension
		panic("FIXME compile: DictComp not implemented")
	case *ast.GeneratorExp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME compile: GeneratorExp not implemented")
	case *ast.Yield:
		// Value Expr
		panic("FIXME compile: Yield not implemented")
	case *ast.YieldFrom:
		// Value Expr
		panic("FIXME compile: YieldFrom not implemented")
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
		args := len(node.Args)
		for i := range node.Args {
			c.Expr(node.Args[i])
		}
		kwargs := len(node.Keywords)
		for i := range node.Keywords {
			kw := node.Keywords[i]
			c.LoadConst(py.String(kw.Arg))
			c.Expr(kw.Value)
		}
		op := byte(vm.CALL_FUNCTION)
		if node.Starargs != nil {
			c.Expr(node.Starargs)
			if node.Kwargs != nil {
				c.Expr(node.Kwargs)
				op = vm.CALL_FUNCTION_VAR_KW
			} else {
				op = vm.CALL_FUNCTION_VAR
			}
		} else if node.Kwargs != nil {
			c.Expr(node.Kwargs)
			op = vm.CALL_FUNCTION_KW
		}
		c.OpArg(op, uint32(args+kwargs<<8))
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
		panic("FIXME compile: Ellipsis not implemented")
	case *ast.Attribute:
		// Value Expr
		// Attr  Identifier
		// Ctx   ExprContext
		// FIXME do something with Ctx
		c.Expr(node.Value)
		c.OpArg(vm.LOAD_ATTR, c.Name(node.Attr))
	case *ast.Subscript:
		// Value Expr
		// Slice Slicer
		// Ctx   ExprContext
		panic("FIXME compile: Subscript not implemented")
	case *ast.Starred:
		// Value Expr
		// Ctx   ExprContext
		panic("FIXME compile: Starred not implemented")
	case *ast.Name:
		// Id  Identifier
		// Ctx ExprContext
		c.NameOp(string(node.Id), node.Ctx)
	case *ast.List:
		// Elts []Expr
		// Ctx  ExprContext
		// FIXME do something with Ctx
		for _, elt := range node.Elts {
			c.Expr(elt)
		}
		c.OpArg(vm.BUILD_LIST, uint32(len(node.Elts)))
	case *ast.Tuple:
		// Elts []Expr
		// Ctx  ExprContext
		// FIXME do something with Ctx
		for _, elt := range node.Elts {
			c.Expr(elt)
		}
		c.OpArg(vm.BUILD_TUPLE, uint32(len(node.Elts)))
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown ExprBase: %v", expr))
	}
}
