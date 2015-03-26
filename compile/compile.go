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
func Compile(str, filename, mode string, flags int, dont_inherit bool) py.Object {
	Ast, err := parser.ParseString(str, mode)
	if err != nil {
		panic(err) // FIXME error handling!
	}
	return CompileAst(Ast, filename, flags, dont_inherit)
}

// As Compile but takes an Ast
func CompileAst(Ast ast.Ast, filename string, flags int, dont_inherit bool) *py.Code {
	//fmt.Println(ast.Dump(Ast))
	code := &py.Code{
		Filename:    filename,
		Firstlineno: 1,                           // FIXME
		Name:        "<module>",                  // FIXME
		Flags:       int32(flags | py.CO_NOFREE), // FIXME
	}
	c := &compiler{
		Code: code,
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
		// Make None the first constant so lambda can't have a docstring
		c.Code.Name = "<lambda>"
		c.Const(py.None) // FIXME extra None for some reason in Consts
		c.Expr(node)
		valueOnStack = true
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown ModuleBase: %v", Ast))
	}
	if !c.OpCodes.EndsWithReturn() {
		// add a return
		if !valueOnStack {
			// return None if there is nothing on the stack
			c.OpArg(vm.LOAD_CONST, c.Const(py.None))
		}
		c.Op(vm.RETURN_VALUE)
	}
	code.Code = c.OpCodes.Assemble()
	code.Stacksize = int32(c.OpCodes.StackDepth())
	return code
}

// State for the compiler
type compiler struct {
	Code    *py.Code // code being built up
	OpCodes Instructions
}

// Compiles a python constant
//
// Returns the index into the Consts tuple
func (c *compiler) Const(obj py.Object) uint32 {
	for i, c := range c.Code.Consts {
		if obj.Type() == c.Type() && py.Eq(obj, c) == py.True {
			return uint32(i)
		}
	}
	c.Code.Consts = append(c.Code.Consts, obj)
	return uint32(len(c.Code.Consts) - 1)
}

// Compiles a python name
//
// Returns the index into the Name tuple
func (c *compiler) Name(Id ast.Identifier) uint32 {
	for i, s := range c.Code.Names {
		if string(Id) == s {
			return uint32(i)
		}
	}
	c.Code.Names = append(c.Code.Names, string(Id))
	return uint32(len(c.Code.Names) - 1)
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
	case vm.JUMP_IF_FALSE_OR_POP, vm.JUMP_IF_TRUE_OR_POP, vm.JUMP_ABSOLUTE, vm.POP_JUMP_IF_FALSE, vm.POP_JUMP_IF_TRUE: // Absolute
		c.OpCodes.Add(&JumpAbs{OpArg: OpArg{Op: Op}, Dest: Dest})
	case vm.JUMP_FORWARD: // Relative
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
		panic("FIXME compile: FunctionDef not implemented")
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
		panic("FIXME compile: Return not implemented")
	case *ast.Delete:
		// Targets []Expr
		panic("FIXME compile: Delete not implemented")
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
		panic("FIXME compile: For not implemented")
	case *ast.While:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		panic("FIXME compile: While not implemented")
	case *ast.If:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		panic("FIXME compile: If not implemented")
	case *ast.With:
		// Items []*WithItem
		// Body  []Stmt
		panic("FIXME compile: With not implemented")
	case *ast.Raise:
		// Exc   Expr
		// Cause Expr
		panic("FIXME compile: Raise not implemented")
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
		panic("FIXME compile: Global not implemented")
	case *ast.Nonlocal:
		// Names []Identifier
		panic("FIXME compile: Nonlocal not implemented")
	case *ast.ExprStmt:
		// Value Expr
		c.Expr(node.Value)
		c.Op(vm.POP_TOP)
	case *ast.Pass:
		// Do nothing
	case *ast.Break:
		panic("FIXME compile: Break not implemented")
	case *ast.Continue:
		panic("FIXME compile: Continue not implemented")
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown StmtBase: %v", stmt))
	}
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
		code := CompileAst(node.Body, c.Code.Filename, int(c.Code.Flags)|py.CO_OPTIMIZED|py.CO_NEWLOCALS, false) // FIXME pass on compile args
		code.Argcount = int32(len(node.Args.Args))
		c.OpArg(vm.LOAD_CONST, c.Const(code))
		c.OpArg(vm.LOAD_CONST, c.Const(py.String("<lambda>")))
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
		panic("FIXME compile: Call not implemented")
	case *ast.Num:
		// N Object
		c.OpArg(vm.LOAD_CONST, c.Const(node.N))
	case *ast.Str:
		// S py.String
		c.OpArg(vm.LOAD_CONST, c.Const(node.S))
	case *ast.Bytes:
		// S py.Bytes
		c.OpArg(vm.LOAD_CONST, c.Const(node.S))
	case *ast.NameConstant:
		// Value Singleton
		c.OpArg(vm.LOAD_CONST, c.Const(node.Value))
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
		switch node.Ctx {
		case ast.Load:
			c.OpArg(vm.LOAD_NAME, c.Name(node.Id))
		case ast.Store:
			c.OpArg(vm.STORE_NAME, c.Name(node.Id))
		// case ast.Del:
		// case ast.AugLoad:
		// case ast.AugStore:
		// case ast.Param:
		default:
			panic(fmt.Sprintf("FIXME ast.Name Ctx=%v not implemented", node.Ctx))
		}
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
