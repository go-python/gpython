package compile

import (
	"github.com/ncw/gpython/vm"
)

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
