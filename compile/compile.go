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
		panic("FIXME compile: FunctionDef not implemented")
		_ = node
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
		panic("FIXME compile: Assign not implemented")
	case *ast.AugAssign:
		// Target Expr
		// Op     OperatorNumber
		// Value  Expr
		panic("FIXME compile: AugAssign not implemented")
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
		panic("FIXME compile: Assert not implemented")
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
		panic("FIXME compile: ExprStmt not implemented")
	case *ast.Pass:
		// No nothing
	case *ast.Break:
		panic("FIXME compile: Break not implemented")
	case *ast.Continue:
		panic("FIXME compile: Continue not implemented")
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
		panic("FIXME compile: Lambda not implemented")
	case *ast.IfExp:
		// Test   Expr
		// Body   Expr
		// Orelse Expr
		elseBranch := new(Label)
		endifBranch := new(Label)
		c.compileExpr(node.Test)
		c.Jump(vm.POP_JUMP_IF_FALSE, elseBranch)
		c.compileExpr(node.Body)
		c.Jump(vm.JUMP_FORWARD, endifBranch)
		c.Label(elseBranch)
		c.compileExpr(node.Orelse)
		c.Label(endifBranch)
	case *ast.Dict:
		// Keys   []Expr
		// Values []Expr
		panic("FIXME compile: Dict not implemented")
	case *ast.Set:
		// Elts []Expr
		panic("FIXME compile: Set not implemented")
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
		panic("FIXME compile: Compare not implemented")
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
		panic("FIXME compile: Bytes not implemented")
	case *ast.NameConstant:
		// Value Singleton
		panic("FIXME compile: NameConstant not implemented")
	case *ast.Ellipsis:
		panic("FIXME compile: Ellipsis not implemented")
	case *ast.Attribute:
		// Value Expr
		// Attr  Identifier
		// Ctx   ExprContext
		panic("FIXME compile: Attribute not implemented")
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
		// FIXME do something with Ctx
		c.OpArg(vm.LOAD_NAME, c.Name(node.Id))
	case *ast.List:
		// Elts []Expr
		// Ctx  ExprContext
		panic("FIXME compile: List not implemented")
	case *ast.Tuple:
		// Elts []Expr
		// Ctx  ExprContext
		panic("FIXME compile: Tuple not implemented")
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
	for i, instr := range is {
		posChanged := instr.SetPos(i, addr)
		changed = changed || posChanged
		if pass > 0 {
			// Only resolve addresses on 2nd pass
			if resolver, ok := instr.(Resolver); ok {
				resolver.Resolve()
			}
		}
		addr += instr.Size()
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

// Calculate number of arguments for CALL_FUNCTION etc
func nArgs(o uint32) int {
	return (int(o) & 0xFF) + 2*((int(o)>>8)&0xFF)
}

// Effect the opcode has on the stack
func opcodeStackEffect(opcode byte, oparg uint32) int {
	switch opcode {
	case vm.POP_TOP:
		return -1
	case vm.ROT_TWO, vm.ROT_THREE:
		return 0
	case vm.DUP_TOP:
		return 1
	case vm.DUP_TOP_TWO:
		return 2
	case vm.UNARY_POSITIVE, vm.UNARY_NEGATIVE, vm.UNARY_NOT, vm.UNARY_INVERT:
		return 0
	case vm.SET_ADD, vm.LIST_APPEND:
		return -1
	case vm.MAP_ADD:
		return -2
	case vm.BINARY_POWER, vm.BINARY_MULTIPLY, vm.BINARY_MODULO, vm.BINARY_ADD, vm.BINARY_SUBTRACT, vm.BINARY_SUBSCR, vm.BINARY_FLOOR_DIVIDE, vm.BINARY_TRUE_DIVIDE:
		return -1
	case vm.INPLACE_FLOOR_DIVIDE, vm.INPLACE_TRUE_DIVIDE:
		return -1
	case vm.INPLACE_ADD, vm.INPLACE_SUBTRACT, vm.INPLACE_MULTIPLY, vm.INPLACE_MODULO:
		return -1
	case vm.STORE_SUBSCR:
		return -3
	case vm.STORE_MAP:
		return -2
	case vm.DELETE_SUBSCR:
		return -2
	case vm.BINARY_LSHIFT, vm.BINARY_RSHIFT, vm.BINARY_AND, vm.BINARY_XOR, vm.BINARY_OR:
		return -1
	case vm.INPLACE_POWER:
		return -1
	case vm.GET_ITER:
		return 0
	case vm.PRINT_EXPR:
		return -1
	case vm.LOAD_BUILD_CLASS:
		return 1
	case vm.INPLACE_LSHIFT, vm.INPLACE_RSHIFT, vm.INPLACE_AND, vm.INPLACE_XOR, vm.INPLACE_OR:
		return -1
	case vm.BREAK_LOOP:
		return 0
	case vm.SETUP_WITH:
		return 7
	case vm.WITH_CLEANUP:
		return -1 /* XXX Sometimes more */
	case vm.RETURN_VALUE:
		return -1
	case vm.IMPORT_STAR:
		return -1
	case vm.YIELD_VALUE:
		return 0
	case vm.YIELD_FROM:
		return -1
	case vm.POP_BLOCK:
		return 0
	case vm.POP_EXCEPT:
		return 0 /* -3 except if bad bytecode */
	case vm.END_FINALLY:
		return -1 /* or -2 or -3 if exception occurred */
	case vm.STORE_NAME:
		return -1
	case vm.DELETE_NAME:
		return 0
	case vm.UNPACK_SEQUENCE:
		return int(oparg) - 1
	case vm.UNPACK_EX:
		return (int(oparg) & 0xFF) + (int(oparg) >> 8)
	case vm.FOR_ITER:
		return 1 /* or -1, at end of iterator */
	case vm.STORE_ATTR:
		return -2
	case vm.DELETE_ATTR:
		return -1
	case vm.STORE_GLOBAL:
		return -1
	case vm.DELETE_GLOBAL:
		return 0
	case vm.LOAD_CONST:
		return 1
	case vm.LOAD_NAME:
		return 1
	case vm.BUILD_TUPLE, vm.BUILD_LIST, vm.BUILD_SET:
		return 1 - int(oparg)
	case vm.BUILD_MAP:
		return 1
	case vm.LOAD_ATTR:
		return 0
	case vm.COMPARE_OP:
		return -1
	case vm.IMPORT_NAME:
		return -1
	case vm.IMPORT_FROM:
		return 1
	case vm.JUMP_FORWARD, vm.JUMP_ABSOLUTE:
		return 0
	case vm.JUMP_IF_TRUE_OR_POP: /* -1 if jump not taken */
		return 0
	case vm.JUMP_IF_FALSE_OR_POP: /*  "" */
		return 0
	case vm.POP_JUMP_IF_FALSE, vm.POP_JUMP_IF_TRUE:
		return -1
	case vm.LOAD_GLOBAL:
		return 1
	case vm.CONTINUE_LOOP:
		return 0
	case vm.SETUP_LOOP:
		return 0
	case vm.SETUP_EXCEPT, vm.SETUP_FINALLY:
		// can push 3 values for the new exception
		// + 3 others for the previous exception state
		return 6
	case vm.LOAD_FAST:
		return 1
	case vm.STORE_FAST:
		return -1
	case vm.DELETE_FAST:
		return 0

	case vm.RAISE_VARARGS:
		return -int(oparg)
	case vm.CALL_FUNCTION:
		return -nArgs(oparg)
	case vm.CALL_FUNCTION_VAR, vm.CALL_FUNCTION_KW:
		return -nArgs(oparg) - 1
	case vm.CALL_FUNCTION_VAR_KW:
		return -nArgs(oparg) - 2
	case vm.MAKE_FUNCTION:
		return -1 - nArgs(oparg) - ((int(oparg) >> 16) & 0xffff)
	case vm.MAKE_CLOSURE:
		return -2 - nArgs(oparg) - ((int(oparg) >> 16) & 0xffff)
	case vm.BUILD_SLICE:
		if oparg == 3 {
			return -2
		} else {
			return -1
		}
	case vm.LOAD_CLOSURE:
		return 1
	case vm.LOAD_DEREF, vm.LOAD_CLASSDEREF:
		return 1
	case vm.STORE_DEREF:
		return -1
	case vm.DELETE_DEREF:
		return 0
	default:
		panic("Unknown opcode in StackEffect")
	}
}

// Recursive instruction walker to find max stack depth
func (is Instructions) stackDepthWalk(baseIs Instructions, seen map[int]bool, startDepth map[int]int, depth int, maxdepth int) int {
	// var i, target_depth, effect int
	// var instr *struct instr
	// if b.b_seen || b.b_startdepth >= depth {
	// 	return maxdepth
	// }
	// b.b_seen = 1
	// b.b_startdepth = depth
	if len(is) == 0 {
		return maxdepth
	}
	start := is[0].Number()
	if seen[start] {
		// We are processing this block already
		return maxdepth
	}
	if d, ok := startDepth[start]; ok && d >= depth {
		// We've processed this block with a larger depth already
		return maxdepth
	}
	seen[start] = true
	startDepth[start] = depth
	for _, instr := range is {
		depth += instr.StackEffect()
		if depth > maxdepth {
			maxdepth = depth
		}
		if depth < 0 {
			panic("Stack depth negative")
		}
		jrel, isJrel := instr.(*JumpRel)
		jabs, isJabs := instr.(*JumpAbs)
		if isJrel || isJabs {
			var oparg *OpArg
			var dest *Label
			if isJrel {
				oparg = &jrel.OpArg
				dest = jrel.Dest
			} else {
				oparg = &jabs.OpArg
				dest = jabs.Dest
			}
			opcode := oparg.Op
			target_depth := depth
			if opcode == vm.FOR_ITER {
				target_depth = depth - 2
			} else if opcode == vm.SETUP_FINALLY || opcode == vm.SETUP_EXCEPT {
				target_depth = depth + 3
				if target_depth > maxdepth {
					maxdepth = target_depth
				}
			} else if opcode == vm.JUMP_IF_TRUE_OR_POP || opcode == vm.JUMP_IF_FALSE_OR_POP {
				depth = depth - 1
			}
			isTarget := baseIs[dest.Number():]
			maxdepth = isTarget.stackDepthWalk(baseIs, seen, startDepth, target_depth, maxdepth)
			if opcode == vm.JUMP_ABSOLUTE ||
				opcode == vm.JUMP_FORWARD {
				goto out // remaining code is dead
			}
		}
	}
out:
	seen[start] = false
	return maxdepth
}

// Find the flow path that needs the largest stack.  We assume that
// cycles in the flow graph have no net effect on the stack depth.
func (is Instructions) StackDepth() int {
	return is.stackDepthWalk(is, make(map[int]bool), make(map[int]int), 0, 0)
}

type Instruction interface {
	Pos() uint32
	Number() int
	SetPos(int, uint32) bool
	Size() uint32
	Output() []byte
	StackEffect() int
}

type Resolver interface {
	Resolve()
}

// Position
type pos struct {
	n uint32
	p uint32
}

// Read instruction number
func (p *pos) Number() int {
	return int(p.n)
}

// Read position
func (p *pos) Pos() uint32 {
	return p.p
}

// Set Position - returns changed
func (p *pos) SetPos(number int, newPos uint32) bool {
	p.n = uint32(number)
	oldPos := p.p
	p.p = newPos
	return oldPos != newPos
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

// StackEffect
func (o *Op) StackEffect() int {
	return opcodeStackEffect(o.Op, 0)
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

// StackEffect
func (o *OpArg) StackEffect() int {
	return opcodeStackEffect(o.Op, o.Arg)
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

// StackEffect
func (o *Label) StackEffect() int {
	return 0
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

// A relative JUMP with destination label
type JumpRel struct {
	pos
	OpArg
	Dest *Label
}

// Set the Arg from the Jump Label
func (o *JumpRel) Resolve() {
	currentSize := o.Size()
	currentPos := o.Pos() + currentSize
	if o.Dest.Pos() < currentPos {
		panic("JUMP_FORWARD can't jump backwards")
	}
	o.OpArg.Arg = o.Dest.Pos() - currentPos
	if o.Size() != currentSize {
		// FIXME There is an awkward moment where jump forwards is
		// between 0x1000 and 0x1002 where the Arg oscillates
		// between 2 and 4 bytes
		panic("FIXME compile: JUMP_FOWARDS size changed")
	}
}
