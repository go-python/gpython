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
func CompileCheat(str, filename, mode string, flags int, dont_inherit bool) py.Object {
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
		panic(err)
	}
	fmt.Println(ast.Dump(Ast))
	code := &py.Code{
		Filename:    filename,
		Stacksize:   1,          // FIXME
		Firstlineno: 1,          // FIXME
		Name:        "<module>", // FIXME
		Flags:       64,         // FIXME
	}
	c := &compiler{
		Code: code,
	}
	switch node := Ast.(type) {
	case *ast.Module:
		c.compileStmts(node.Body)
	case *ast.Interactive:
		c.compileStmts(node.Body)
	case *ast.Expression:
		c.compileExpr(node.Body)
	case *ast.Suite:
		c.compileStmts(node.Body)
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown ModuleBase: %v", Ast))
	}
	c.Op(vm.RETURN_VALUE)
	code.Code = c.OpCodes.Assemble()
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
	// FIXME find an existing one
	c.Code.Consts = append(c.Code.Consts, obj)
	return uint32(len(c.Code.Consts) - 1)
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
		panic("FIXME JUMP_FORWARD NOT implemented")
	default:
		panic("Jump called with non jump instruction")
	}
}

// Compile statements
func (c *compiler) compileStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		c.compileStmt(stmt)
	}
}

// Compile statement
func (c *compiler) compileStmt(stmt ast.Stmt) {
	switch node := stmt.(type) {
	case *ast.FunctionDef:
		// Name          Identifier
		// Args          *Arguments
		// Body          []Stmt
		// DecoratorList []Expr
		// Returns       Expr
		panic("FIXME not implemented")
		_ = node
	case *ast.ClassDef:
		// Name          Identifier
		// Bases         []Expr
		// Keywords      []*Keyword
		// Starargs      Expr
		// Kwargs        Expr
		// Body          []Stmt
		// DecoratorList []Expr
		panic("FIXME not implemented")
	case *ast.Return:
		// Value Expr
		panic("FIXME not implemented")
	case *ast.Delete:
		// Targets []Expr
		panic("FIXME not implemented")
	case *ast.Assign:
		// Targets []Expr
		// Value   Expr
		panic("FIXME not implemented")
	case *ast.AugAssign:
		// Target Expr
		// Op     OperatorNumber
		// Value  Expr
		panic("FIXME not implemented")
	case *ast.For:
		// Target Expr
		// Iter   Expr
		// Body   []Stmt
		// Orelse []Stmt
		panic("FIXME not implemented")
	case *ast.While:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		panic("FIXME not implemented")
	case *ast.If:
		// Test   Expr
		// Body   []Stmt
		// Orelse []Stmt
		panic("FIXME not implemented")
	case *ast.With:
		// Items []*WithItem
		// Body  []Stmt
		panic("FIXME not implemented")
	case *ast.Raise:
		// Exc   Expr
		// Cause Expr
		panic("FIXME not implemented")
	case *ast.Try:
		// Body      []Stmt
		// Handlers  []*ExceptHandler
		// Orelse    []Stmt
		// Finalbody []Stmt
		panic("FIXME not implemented")
	case *ast.Assert:
		// Test Expr
		// Msg  Expr
		panic("FIXME not implemented")
	case *ast.Import:
		// Names []*Alias
		panic("FIXME not implemented")
	case *ast.ImportFrom:
		// Module Identifier
		// Names  []*Alias
		// Level  int
		panic("FIXME not implemented")
	case *ast.Global:
		// Names []Identifier
		panic("FIXME not implemented")
	case *ast.Nonlocal:
		// Names []Identifier
		panic("FIXME not implemented")
	case *ast.ExprStmt:
		// Value Expr
		panic("FIXME not implemented")
	case *ast.Pass:
		// No nothing
	case *ast.Break:
		panic("FIXME not implemented")
	case *ast.Continue:
		panic("FIXME not implemented")
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown StmtBase: %v", stmt))
	}
}

// Compile expressions
func (c *compiler) compileExpr(expr ast.Expr) {
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
			c.compileExpr(e)
			if i != len(node.Values)-1 {
				c.Jump(op, label)
			}
		}
		c.Label(label)
	case *ast.BinOp:
		// Left  Expr
		// Op    OperatorNumber
		// Right Expr
		c.compileExpr(node.Left)
		c.compileExpr(node.Right)
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
		c.compileExpr(node.Operand)
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
		panic("FIXME not implemented")
	case *ast.IfExp:
		// Test   Expr
		// Body   Expr
		// Orelse Expr
		panic("FIXME not implemented")
	case *ast.Dict:
		// Keys   []Expr
		// Values []Expr
		panic("FIXME not implemented")
	case *ast.Set:
		// Elts []Expr
		panic("FIXME not implemented")
	case *ast.ListComp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME not implemented")
	case *ast.SetComp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME not implemented")
	case *ast.DictComp:
		// Key        Expr
		// Value      Expr
		// Generators []Comprehension
		panic("FIXME not implemented")
	case *ast.GeneratorExp:
		// Elt        Expr
		// Generators []Comprehension
		panic("FIXME not implemented")
	case *ast.Yield:
		// Value Expr
		panic("FIXME not implemented")
	case *ast.YieldFrom:
		// Value Expr
		panic("FIXME not implemented")
	case *ast.Compare:
		// Left        Expr
		// Ops         []CmpOp
		// Comparators []Expr
		panic("FIXME not implemented")
	case *ast.Call:
		// Func     Expr
		// Args     []Expr
		// Keywords []*Keyword
		// Starargs Expr
		// Kwargs   Expr
		panic("FIXME not implemented")
	case *ast.Num:
		// N Object
		c.OpArg(vm.LOAD_CONST, c.Const(node.N))
	case *ast.Str:
		// S py.String
		c.OpArg(vm.LOAD_CONST, c.Const(node.S))
	case *ast.Bytes:
		// S py.Bytes
		panic("FIXME not implemented")
	case *ast.NameConstant:
		// Value Singleton
		panic("FIXME not implemented")
	case *ast.Ellipsis:
		panic("FIXME not implemented")
	case *ast.Attribute:
		// Value Expr
		// Attr  Identifier
		// Ctx   ExprContext
		panic("FIXME not implemented")
	case *ast.Subscript:
		// Value Expr
		// Slice Slicer
		// Ctx   ExprContext
		panic("FIXME not implemented")
	case *ast.Starred:
		// Value Expr
		// Ctx   ExprContext
		panic("FIXME not implemented")
	case *ast.Name:
		// Id  Identifier
		// Ctx ExprContext
		panic("FIXME not implemented")
	case *ast.List:
		// Elts []Expr
		// Ctx  ExprContext
		panic("FIXME not implemented")
	case *ast.Tuple:
		// Elts []Expr
		// Ctx  ExprContext
		panic("FIXME not implemented")
	default:
		panic(py.ExceptionNewf(py.SyntaxError, "Unknown ExprBase: %v", expr))
	}
}

// Resolved or unresolved instruction stream
type Instructions []Instruction

// Add an instruction to the instructions
func (is *Instructions) Add(i Instruction) {
	*is = append(*is, i)
}

// Do a pass of assembly
//
// Returns a boolean as to whether the stream changed
func (is Instructions) Pass(pass int) bool {
	addr := uint32(0)
	changed := false
	for _, i := range is {
		changed = changed || i.SetPos(addr)
		if pass > 0 {
			// Only resolve addresses on 2nd pass
			if resolver, ok := i.(Resolver); ok {
				resolver.Resolve()
			}
		}
		addr += i.Size()
	}
	return changed
}

// Assemble the instructions into an Opcode string
func (is Instructions) Assemble() string {
	for i := 0; i < 10; i++ {
		changed := is.Pass(i)
		if !changed {
			goto done
		}
	}
	panic("Failed to assemble after 10 passes")
done:
	out := make([]byte, 0, 3*len(is))
	for _, i := range is {
		out = append(out, i.Output()...)
	}
	return string(out)
}

type Instruction interface {
	Pos() uint32
	SetPos(uint32) bool
	Size() uint32
	Output() []byte
}

type Resolver interface {
	Resolve()
}

// Position
type pos uint32

// Read position
func (p *pos) Pos() uint32 {
	return uint32(*p)
}

// Set Position - returns changed
func (p *pos) SetPos(newPos uint32) bool {
	oldP := *p
	newP := pos(newPos)
	*p = newP
	return oldP != newP
}

// A plain opcode
type Op struct {
	pos
	Op byte
}

// Uses 1 byte in the output stream
func (o *Op) Size() uint32 {
	return 1
}

// Output
func (o *Op) Output() []byte {
	return []byte{byte(o.Op)}
}

// An opcode with argument
type OpArg struct {
	pos
	Op  byte
	Arg uint32
}

// Uses 1 byte in the output stream
func (o *OpArg) Size() uint32 {
	if o.Arg <= 0xFFFF {
		return 3 // Op Arg1 Arg2
	} else {
		return 6 // Extend Arg1 Arg2 Op Arg3 Arg4
	}
}

// Output
func (o *OpArg) Output() []byte {
	out := []byte{o.Op, byte(o.Arg), byte(o.Arg >> 8)}
	if o.Arg > 0xFFFF {
		out = append([]byte{vm.EXTENDED_ARG, byte(o.Arg >> 16), byte(o.Arg >> 24)}, out...)
	}
	return out
}

// A label
type Label struct {
	pos
}

// Uses 0 bytes in the output stream
func (o *Label) Size() uint32 {
	return 0
}

// Output
func (o Label) Output() []byte {
	return []byte{}
}

// An absolute JUMP with destination label
type JumpAbs struct {
	pos
	OpArg
	Dest *Label
}

// Set the Arg from the Jump Label
func (o *JumpAbs) Resolve() {
	o.OpArg.Arg = o.Dest.Pos()
}

// FIXME Jump Relative
