// Evaluate opcodes
package vm

/* FIXME

cpython has one stack per frame, not one stack in total

We know how big each frame needs to be from

code->co_stacksize

The frame then becomes the important thing

cpython keeps a zombie frame on each code object to speed up execution
of a code object so a frame doesn't have to be allocated and
deallocated each time which seems like a good idea.  If we want to
work with go routines then it might have to be more sophisticated.

To implmenent generators need to check Code.Flags & CO_GENERATOR at
the start of vmRum and if so wrap the created frame into a generator
object.

FIXME could make the stack be permanently allocated and just keep a
pointer into it rather than using append etc...

If we are caching the frames need to make sure we clear the stack
objects so they can be GCed

*/

import (
	"fmt"
	"github.com/ncw/gpython/py"
	"runtime/debug"
)

const (
	nameErrorMsg         = "name '%s' is not defined"
	globalNameErrorMsg   = "global name '%s' is not defined"
	unboundLocalErrorMsg = "local variable '%s' referenced before assignment"
	unboundFreeErrorMsg  = "free variable '%s' referenced before assignment in enclosing scope"
	cannotCatchMsg       = "catching '%s' that does not inherit from BaseException is not allowed"
)

// Stack operations
func (vm *Vm) STACK_LEVEL() int             { return len(vm.frame.Stack) }
func (vm *Vm) EMPTY() bool                  { return len(vm.frame.Stack) == 0 }
func (vm *Vm) TOP() py.Object               { return vm.frame.Stack[len(vm.frame.Stack)-1] }
func (vm *Vm) SECOND() py.Object            { return vm.frame.Stack[len(vm.frame.Stack)-2] }
func (vm *Vm) THIRD() py.Object             { return vm.frame.Stack[len(vm.frame.Stack)-3] }
func (vm *Vm) FOURTH() py.Object            { return vm.frame.Stack[len(vm.frame.Stack)-4] }
func (vm *Vm) PEEK(n int) py.Object         { return vm.frame.Stack[len(vm.frame.Stack)-n] }
func (vm *Vm) SET_TOP(v py.Object)          { vm.frame.Stack[len(vm.frame.Stack)-1] = v }
func (vm *Vm) SET_SECOND(v py.Object)       { vm.frame.Stack[len(vm.frame.Stack)-2] = v }
func (vm *Vm) SET_THIRD(v py.Object)        { vm.frame.Stack[len(vm.frame.Stack)-3] = v }
func (vm *Vm) SET_FOURTH(v py.Object)       { vm.frame.Stack[len(vm.frame.Stack)-4] = v }
func (vm *Vm) SET_VALUE(n int, v py.Object) { vm.frame.Stack[len(vm.frame.Stack)-(n)] = (v) }
func (vm *Vm) DROP()                        { vm.frame.Stack = vm.frame.Stack[:len(vm.frame.Stack)-1] }
func (vm *Vm) DROPN(n int)                  { vm.frame.Stack = vm.frame.Stack[:len(vm.frame.Stack)-n] }

// Pop from top of vm stack
func (vm *Vm) POP() py.Object {
	// FIXME what if empty?
	out := vm.frame.Stack[len(vm.frame.Stack)-1]
	vm.frame.Stack = vm.frame.Stack[:len(vm.frame.Stack)-1]
	return out
}

// Push to top of vm stack
func (vm *Vm) PUSH(obj py.Object) {
	vm.frame.Stack = append(vm.frame.Stack, obj)
}

// Adds a traceback to the exc passed in for the current vm state
func (vm *Vm) AddTraceback(exc *py.ExceptionInfo) {
	exc.Traceback = &py.Traceback{
		Next:   exc.Traceback,
		Frame:  vm.frame,
		Lasti:  vm.frame.Lasti,
		Lineno: vm.frame.Code.Addr2Line(vm.frame.Lasti),
	}
}

// Set an exception in the VM
//
// The exception must be a valid exception instance (eg as returned by
// py.MakeException)
//
// It sets vm.exc.* and sets vm.exit to exitException
func (vm *Vm) SetException(exception py.Object) {
	vm.old_exc = vm.exc
	vm.exc.Value = exception
	vm.exc.Type = exception.Type()
	vm.exc.Traceback = nil
	vm.AddTraceback(&vm.exc)
	vm.exit = exitException
}

// Check for an exception (panic)
//
// Should be called with the result of recover
func (vm *Vm) CheckExceptionRecover(r interface{}) {
	// If what was raised was an ExceptionInfo the stuff this into the current vm
	if exc, ok := r.(py.ExceptionInfo); ok {
		vm.old_exc = vm.exc
		vm.exc = exc
		vm.AddTraceback(&vm.exc)
		vm.exit = exitException
		fmt.Printf("*** Propagating exception: %s\n", exc.Error())
	} else {
		// Coerce whatever was raised into a *Exception
		vm.SetException(py.MakeException(r))
		fmt.Printf("*** Exception raised %v\n", r)
		// Dump the goroutine stack
		debug.PrintStack()
	}
}

// Check for an exception (panic)
//
// Must be called as a defer function
func (vm *Vm) CheckException() {
	if r := recover(); r != nil {
		vm.CheckExceptionRecover(r)
	}
}

// Checks if r is StopIteration and if so returns true
//
// Otherwise deals with the as per vm.CheckException and returns false
func (vm *Vm) catchStopIteration(r interface{}) bool {
	if py.IsException(py.StopIteration, r) {
		// StopIteration or subclass raises
		return true
	} else {
		// Deal with the exception as normal
		vm.CheckExceptionRecover(r)
	}
	return false
}

// Illegal instruction
func do_ILLEGAL(vm *Vm, arg int32) {
	defer vm.CheckException()
	panic("Illegal opcode")
}

// Do nothing code. Used as a placeholder by the bytecode optimizer.
func do_NOP(vm *Vm, arg int32) {
}

// Removes the top-of-stack (TOS) item.
func do_POP_TOP(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.DROPN(1)
}

// Swaps the two top-most stack items.
func do_ROT_TWO(vm *Vm, arg int32) {
	defer vm.CheckException()
	top := vm.TOP()
	second := vm.SECOND()
	vm.SET_TOP(second)
	vm.SET_SECOND(top)
}

// Lifts second and third stack item one position up, moves top down
// to position three.
func do_ROT_THREE(vm *Vm, arg int32) {
	defer vm.CheckException()
	top := vm.TOP()
	second := vm.SECOND()
	third := vm.THIRD()
	vm.SET_TOP(second)
	vm.SET_SECOND(third)
	vm.SET_THIRD(top)
}

// Duplicates the reference on top of the stack.
func do_DUP_TOP(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.PUSH(vm.TOP())
}

// Duplicates the top two reference on top of the stack.
func do_DUP_TOP_TWO(vm *Vm, arg int32) {
	defer vm.CheckException()
	top := vm.TOP()
	second := vm.SECOND()
	vm.PUSH(second)
	vm.PUSH(top)
}

// Unary Operations take the top of the stack, apply the operation,
// and push the result back on the stack.

// Implements TOS = +TOS.
func do_UNARY_POSITIVE(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.Pos(vm.TOP()))
}

// Implements TOS = -TOS.
func do_UNARY_NEGATIVE(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.Neg(vm.TOP()))
}

// Implements TOS = not TOS.
func do_UNARY_NOT(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.Not(vm.TOP()))
}

// Implements TOS = ~TOS.
func do_UNARY_INVERT(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.Invert(vm.TOP()))
}

// Implements TOS = iter(TOS).
func do_GET_ITER(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.Iter(vm.TOP()))
}

// Pops TOS from the stack and stores it as the current frame’s
// f_locals. This is used in class construction.
func do_STORE_LOCALS(vm *Vm, arg int32) {
	defer vm.CheckException()
	locals := vm.POP()
	vm.frame.Locals = locals.(py.StringDict)
}

// Binary operations remove the top of the stack (TOS) and the second
// top-most stack item (TOS1) from the stack. They perform the
// operation, and put the result back on the stack.

// Implements TOS = TOS1 ** TOS.
func do_BINARY_POWER(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Pow(a, b, py.None))
}

// Implements TOS = TOS1 * TOS.
func do_BINARY_MULTIPLY(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Mul(a, b))
}

// Implements TOS = TOS1 // TOS.
func do_BINARY_FLOOR_DIVIDE(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.FloorDiv(a, b))
}

// Implements TOS = TOS1 / TOS when from __future__ import division is
// in effect.
func do_BINARY_TRUE_DIVIDE(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.TrueDiv(a, b))
}

// Implements TOS = TOS1 % TOS.
func do_BINARY_MODULO(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Mod(a, b))
}

// Implements TOS = TOS1 + TOS.
func do_BINARY_ADD(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Add(a, b))
}

// Implements TOS = TOS1 - TOS.
func do_BINARY_SUBTRACT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Sub(a, b))
}

// Implements TOS = TOS1[TOS].
func do_BINARY_SUBSCR(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.GetItem(a, b))
}

// Implements TOS = TOS1 << TOS.
func do_BINARY_LSHIFT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Lshift(a, b))
}

// Implements TOS = TOS1 >> TOS.
func do_BINARY_RSHIFT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Rshift(a, b))
}

// Implements TOS = TOS1 & TOS.
func do_BINARY_AND(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.And(a, b))
}

// Implements TOS = TOS1 ^ TOS.
func do_BINARY_XOR(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Xor(a, b))
}

// Implements TOS = TOS1 | TOS.
func do_BINARY_OR(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Or(a, b))
}

// In-place operations are like binary operations, in that they remove
// TOS and TOS1, and push the result back on the stack, but the
// operation is done in-place when TOS1 supports it, and the resulting
// TOS may be (but does not have to be) the original TOS1.

// Implements in-place TOS = TOS1 ** TOS.
func do_INPLACE_POWER(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IPow(a, b, py.None))
}

// Implements in-place TOS = TOS1 * TOS.
func do_INPLACE_MULTIPLY(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IMul(a, b))
}

// Implements in-place TOS = TOS1 // TOS.
func do_INPLACE_FLOOR_DIVIDE(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IFloorDiv(a, b))
}

// Implements in-place TOS = TOS1 / TOS when from __future__ import
// division is in effect.
func do_INPLACE_TRUE_DIVIDE(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.ITrueDiv(a, b))
}

// Implements in-place TOS = TOS1 % TOS.
func do_INPLACE_MODULO(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.Mod(a, b))
}

// Implements in-place TOS = TOS1 + TOS.
func do_INPLACE_ADD(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IAdd(a, b))
}

// Implements in-place TOS = TOS1 - TOS.
func do_INPLACE_SUBTRACT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.ISub(a, b))
}

// Implements in-place TOS = TOS1 << TOS.
func do_INPLACE_LSHIFT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.ILshift(a, b))
}

// Implements in-place TOS = TOS1 >> TOS.
func do_INPLACE_RSHIFT(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IRshift(a, b))
}

// Implements in-place TOS = TOS1 & TOS.
func do_INPLACE_AND(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IAnd(a, b))
}

// Implements in-place TOS = TOS1 ^ TOS.
func do_INPLACE_XOR(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IXor(a, b))
}

// Implements in-place TOS = TOS1 | TOS.
func do_INPLACE_OR(vm *Vm, arg int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	vm.SET_TOP(py.IOr(a, b))
}

// Implements TOS1[TOS] = TOS2.
func do_STORE_SUBSCR(vm *Vm, arg int32) {
	defer vm.CheckException()
	w := vm.TOP()
	v := vm.SECOND()
	u := vm.THIRD()
	vm.DROPN(3)
	// v[w] = u
	py.SetItem(v, w, u)
}

// Implements del TOS1[TOS].
func do_DELETE_SUBSCR(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.NotImplemented("DELETE_SUBSCR", arg)
}

// Miscellaneous opcodes.

// Implements the expression statement for the interactive mode. TOS
// is removed from the stack and printed. In non-interactive mode, an
// expression statement is terminated with POP_STACK.
func do_PRINT_EXPR(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.NotImplemented("PRINT_EXPR", arg)
}

// Terminates a loop due to a break statement.
func do_BREAK_LOOP(vm *Vm, arg int32) {
	defer vm.CheckException()
	// Jump
	vm.frame.Lasti = vm.frame.Block.Handler
	// Reset the stack (FIXME?)
	vm.frame.Stack = vm.frame.Stack[:vm.frame.Block.Level]
	vm.frame.PopBlock()
}

// Continues a loop due to a continue statement. target is the address
// to jump to (which should be a FOR_ITER instruction).
func do_CONTINUE_LOOP(vm *Vm, target int32) {
	defer vm.CheckException()
	switch vm.frame.Block.Type {
	case SETUP_LOOP:
	case SETUP_WITH:
		vm.NotImplemented("CONTINUE_LOOP WITH", target)
	case SETUP_EXCEPT:
		vm.NotImplemented("CONTINUE_LOOP EXCEPT", target)
	case SETUP_FINALLY:
		vm.NotImplemented("CONTINUE_LOOP FINALLY", target)
	default:
	}
	vm.frame.Lasti = vm.frame.Block.Handler
}

// Implements assignment with a starred target: Unpacks an iterable in
// TOS into individual values, where the total number of values can be
// smaller than the number of items in the iterable: one the new
// values will be a list of all leftover items.
//
// The low byte of counts is the number of values before the list
// value, the high byte of counts the number of values after it. The
// resulting values are put onto the stack right-to-left.
func do_UNPACK_EX(vm *Vm, counts int32) {
	defer vm.CheckException()
	vm.NotImplemented("UNPACK_EX", counts)
}

// Calls set.add(TOS1[-i], TOS). Used to implement set comprehensions.
func do_SET_ADD(vm *Vm, i int32) {
	defer vm.CheckException()
	vm.NotImplemented("SET_ADD", i)
}

// Calls list.append(TOS[-i], TOS). Used to implement list
// comprehensions. While the appended value is popped off, the list
// object remains on the stack so that it is available for further
// iterations of the loop.
func do_LIST_APPEND(vm *Vm, i int32) {
	defer vm.CheckException()
	w := vm.POP()
	v := vm.PEEK(int(i))
	v.(*py.List).Append(w)
}

// Calls dict.setitem(TOS1[-i], TOS, TOS1). Used to implement dict comprehensions.
func do_MAP_ADD(vm *Vm, i int32) {
	defer vm.CheckException()
	vm.NotImplemented("MAP_ADD", i)
}

// Returns with TOS to the caller of the function.
func do_RETURN_VALUE(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.result = vm.POP()
	if len(vm.frame.Stack) != 0 {
		fmt.Printf("vmstack = %#v\n", vm.frame.Stack)
		panic("vm stack should be empty at this point")
	}
	vm.frame.Yielded = false
	vm.exit = exitReturn
}

// Pops TOS and delegates to it as a subiterator from a generator.
func do_YIELD_FROM(vm *Vm, arg int32) {
	defer func() {
		if r := recover(); r != nil {
			if vm.catchStopIteration(r) {
				// No extra action needed
			}
		}
	}()

	var retval py.Object
	u := vm.POP()
	x := vm.TOP()
	// send u to x
	if u == py.None {
		retval = py.Next(x)
	} else {
		retval = py.Send(x, u)
	}
	// x remains on stack, retval is value to be yielded
	// FIXME vm.frame.Stacktop = stack_pointer
	//why = exitYield
	// and repeat...
	vm.frame.Lasti--

	vm.result = retval
	vm.frame.Yielded = true
	vm.exit = exitYield
}

// Pops TOS and yields it from a generator.
func do_YIELD_VALUE(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.result = vm.POP()
	vm.frame.Yielded = true
	vm.exit = exitYield
}

// Loads all symbols not starting with '_' directly from the module
// TOS to the local namespace. The module is popped after loading all
// names. This opcode implements from module import *.
func do_IMPORT_STAR(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.NotImplemented("IMPORT_STAR", arg)
}

// Removes one block from the block stack. Per frame, there is a stack
// of blocks, denoting nested loops, try statements, and such.
func do_POP_BLOCK(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.frame.PopBlock()
}

// Removes one block from the block stack. The popped block must be an
// exception handler block, as implicitly created when entering an
// except handler. In addition to popping extraneous values from the
// frame stack, the last three popped values are used to restore the
// exception state.
func do_POP_EXCEPT(vm *Vm, arg int32) {
	defer vm.CheckException()
	frame := vm.frame
	b := vm.frame.Block
	frame.PopBlock()
	if b.Type != EXCEPT_HANDLER {
		vm.SetException(py.ExceptionNewf(py.SystemError, "popped block is not an except handler"))
	} else {
		vm.UnwindExceptHandler(frame, b)
	}
}

// Terminates a finally clause. The interpreter recalls whether the
// exception has to be re-raised, or whether the function returns, and
// continues with the outer-next block.
func do_END_FINALLY(vm *Vm, arg int32) {
	defer vm.CheckException()
	v := vm.POP()
	if vInt, ok := v.(py.Int); ok {
		vm.exit = vmExit(vInt)
		if vm.exit == exitYield {
			panic("Unexpected exitYield in END_FINALLY")
		}
		if vm.exit == exitReturn || vm.exit == exitContinue {
			// Leave return value on the stack
			// retval = vm.POP()
		}
		if vm.exit == exitSilenced {
			// An exception was silenced by 'with', we must
			// manually unwind the EXCEPT_HANDLER block which was
			// created when the exception was caught, otherwise
			// the stack will be in an inconsistent state.
			frame := vm.frame
			b := vm.frame.Block
			frame.PopBlock()
			if b.Type != EXCEPT_HANDLER {
				panic("Expecting EXCEPT_HANDLER in END_FINALLY")
			}
			vm.UnwindExceptHandler(frame, b)
			vm.exit = exitNot
		}
	} else if py.ExceptionClassCheck(v) {
		w := vm.POP()
		u := vm.POP()
		// FIXME PyErr_Restore(v, w, u)
		vm.exc.Type = v.(*py.Type)
		vm.exc.Value = w
		vm.exc.Traceback = u.(*py.Traceback)
		vm.exit = exitReraise
	} else if v != py.None {
		vm.SetException(py.ExceptionNewf(py.SystemError, "'finally' pops bad exception %#v", v))
	}
}

// Loads the __build_class__ helper function to the stack which
// creates a new class object.
func do_LOAD_BUILD_CLASS(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.PUSH(py.Builtins.Globals["__build_class__"])
}

// This opcode performs several operations before a with block
// starts. First, it loads __exit__( ) from the context manager and
// pushes it onto the stack for later use by WITH_CLEANUP. Then,
// __enter__( ) is called, and a finally block pointing to delta is
// pushed. Finally, the result of calling the enter method is pushed
// onto the stack. The next opcode will either ignore it (POP_TOP), or
// store it in (a) variable(s) (STORE_FAST, STORE_NAME, or
// UNPACK_SEQUENCE).
func do_SETUP_WITH(vm *Vm, delta int32) {
	defer vm.CheckException()
	vm.NotImplemented("SETUP_WITH", delta)
}

// Cleans up the stack when a with statement block exits. On top of
// the stack are 1–3 values indicating how/why the finally clause was
// entered:
//
// TOP = None
// (TOP, SECOND) = (WHY_{RETURN,CONTINUE}), retval
// TOP = WHY_*; no retval below it
// (TOP, SECOND, THIRD) = exc_info( )
// Under them is EXIT, the context manager’s __exit__( ) bound method.
//
// In the last case, EXIT(TOP, SECOND, THIRD) is called, otherwise
// EXIT(None, None, None).
//
// EXIT is removed from the stack, leaving the values above it in the
// same order. In addition, if the stack represents an exception, and
// the function call returns a ‘true’ value, this information is
// “zapped”, to prevent END_FINALLY from re-raising the
// exception. (But non-local gotos should still be resumed.)
func do_WITH_CLEANUP(vm *Vm, arg int32) {
	defer vm.CheckException()
	vm.NotImplemented("WITH_CLEANUP", arg)
}

// All of the following opcodes expect arguments. An argument is two bytes, with the more significant byte last.

// Implements name = TOS. namei is the index of name in the attribute
// co_names of the code object. The compiler tries to use STORE_FAST
// or STORE_GLOBAL if possible.
func do_STORE_NAME(vm *Vm, namei int32) {
	defer vm.CheckException()
	fmt.Printf("STORE_NAME %v\n", vm.frame.Code.Names[namei])
	vm.frame.Locals[vm.frame.Code.Names[namei]] = vm.POP()
}

// Implements del name, where namei is the index into co_names
// attribute of the code object.
func do_DELETE_NAME(vm *Vm, namei int32) {
	defer vm.CheckException()
	vm.NotImplemented("DELETE_NAME", namei)
}

// Unpacks TOS into count individual values, which are put onto the
// stack right-to-left.
func do_UNPACK_SEQUENCE(vm *Vm, count int32) {
	defer vm.CheckException()
	vm.NotImplemented("UNPACK_SEQUENCE", count)
}

// Implements TOS.name = TOS1, where namei is the index of name in
// co_names.
func do_STORE_ATTR(vm *Vm, namei int32) {
	defer vm.CheckException()
	w := vm.frame.Code.Names[namei]
	v := vm.TOP()
	u := vm.SECOND()
	vm.DROPN(2)
	py.SetAttrString(v, w, u) /* v.w = u */
}

// Implements del TOS.name, using namei as index into co_names.
func do_DELETE_ATTR(vm *Vm, namei int32) {
	defer vm.CheckException()
	py.DeleteAttrString(vm.POP(), vm.frame.Code.Names[namei])
}

// Works as STORE_NAME, but stores the name as a global.
func do_STORE_GLOBAL(vm *Vm, namei int32) {
	defer vm.CheckException()
	vm.NotImplemented("STORE_GLOBAL", namei)
}

// Works as DELETE_NAME, but deletes a global name.
func do_DELETE_GLOBAL(vm *Vm, namei int32) {
	defer vm.CheckException()
	vm.NotImplemented("DELETE_GLOBAL", namei)
}

// Pushes co_consts[consti] onto the stack.
func do_LOAD_CONST(vm *Vm, consti int32) {
	defer vm.CheckException()
	vm.PUSH(vm.frame.Code.Consts[consti])
	// fmt.Printf("LOAD_CONST %v\n", vm.TOP())
}

// Pushes the value associated with co_names[namei] onto the stack.
func do_LOAD_NAME(vm *Vm, namei int32) {
	defer vm.CheckException()
	fmt.Printf("LOAD_NAME %v\n", vm.frame.Code.Names[namei])
	vm.PUSH(vm.frame.Lookup(vm.frame.Code.Names[namei]))
}

// Creates a tuple consuming count items from the stack, and pushes
// the resulting tuple onto the stack.
func do_BUILD_TUPLE(vm *Vm, count int32) {
	defer vm.CheckException()
	tuple := make(py.Tuple, count)
	copy(tuple, vm.frame.Stack[len(vm.frame.Stack)-int(count):])
	vm.DROPN(int(count))
	vm.PUSH(tuple)
}

// Works as BUILD_TUPLE, but creates a set.
func do_BUILD_SET(vm *Vm, count int32) {
	defer vm.CheckException()
	set := py.NewSetFromItems(vm.frame.Stack[len(vm.frame.Stack)-int(count):])
	vm.DROPN(int(count))
	vm.PUSH(set)
}

// Works as BUILD_TUPLE, but creates a list.
func do_BUILD_LIST(vm *Vm, count int32) {
	defer vm.CheckException()
	list := py.NewListFromItems(vm.frame.Stack[len(vm.frame.Stack)-int(count):])
	vm.DROPN(int(count))
	vm.PUSH(list)
}

// Pushes a new dictionary object onto the stack. The dictionary is
// pre-sized to hold count entries.
func do_BUILD_MAP(vm *Vm, count int32) {
	defer vm.CheckException()
	vm.PUSH(py.NewStringDictSized(int(count)))
}

// Replaces TOS with getattr(TOS, co_names[namei]).
func do_LOAD_ATTR(vm *Vm, namei int32) {
	defer vm.CheckException()
	vm.SET_TOP(py.GetAttrString(vm.TOP(), vm.frame.Code.Names[namei]))
}

// Performs a Boolean operation. The operation name can be found in
// cmp_op[opname].
func do_COMPARE_OP(vm *Vm, opname int32) {
	defer vm.CheckException()
	b := vm.POP()
	a := vm.TOP()
	var r py.Object
	switch opname {
	case PyCmp_LT:
		r = py.Lt(a, b)
	case PyCmp_LE:
		r = py.Le(a, b)
	case PyCmp_EQ:
		r = py.Eq(a, b)
	case PyCmp_NE:
		r = py.Ne(a, b)
	case PyCmp_GT:
		r = py.Gt(a, b)
	case PyCmp_GE:
		r = py.Ge(a, b)
	case PyCmp_IN:
		vm.NotImplemented("COMPARE_OP PyCmp_IN", opname)
	case PyCmp_NOT_IN:
		vm.NotImplemented("COMPARE_OP PyCmp_NOT_IN", opname)
	case PyCmp_IS:
		// FIXME not right
		r = py.NewBool(a == b)
		vm.NotImplemented("COMPARE_OP PyCmp_IS", opname)
	case PyCmp_IS_NOT:
		// FIXME not right
		r = py.NewBool(a != b)
		vm.NotImplemented("COMPARE_OP PyCmp_NOT_IS", opname)
	case PyCmp_EXC_MATCH:
		if bTuple, ok := b.(py.Tuple); ok {
			for _, exc := range bTuple {
				if !py.ExceptionClassCheck(exc) {
					vm.SetException(py.ExceptionNewf(py.TypeError, cannotCatchMsg, exc.Type().Name))
					goto finished
				}
			}
		} else {
			if !py.ExceptionClassCheck(b) {
				vm.SetException(py.ExceptionNewf(py.TypeError, cannotCatchMsg, b.Type().Name))
				goto finished
			}
		}
		r = py.NewBool(py.ExceptionGivenMatches(a, b))
	finished:
		;
	case PyCmp_BAD:
		vm.NotImplemented("COMPARE_OP PyCmp_BAD", opname)
	default:
		vm.NotImplemented("COMPARE_OP", opname)
	}
	vm.SET_TOP(r)
}

// Imports the module co_names[namei]. TOS and TOS1 are popped and
// provide the fromlist and level arguments of __import__( ). The
// module object is pushed onto the stack. The current namespace is
// not affected: for a proper import statement, a subsequent
// STORE_FAST instruction modifies the namespace.
func do_IMPORT_NAME(vm *Vm, namei int32) {
	defer vm.CheckException()
	name := py.String(vm.frame.Code.Names[namei])
	__import__, ok := vm.frame.Builtins["__import__"]
	if !ok {
		panic(py.ExceptionNewf(py.ImportError, "__import__ not found"))
	}
	v := vm.POP()
	u := vm.TOP()
	var locals py.Object = vm.frame.Locals
	if locals == nil {
		locals = py.None
	}
	var args py.Tuple
	if _, ok := u.(py.Int); ok {
		args = py.Tuple{name, vm.frame.Globals, locals, v, u}
	} else {
		args = py.Tuple{name, vm.frame.Globals, locals, v}
	}
	x := py.Call(__import__, args, nil)
	vm.SET_TOP(x)
}

// Loads the attribute co_names[namei] from the module found in
// TOS. The resulting object is pushed onto the stack, to be
// subsequently stored by a STORE_FAST instruction.
func do_IMPORT_FROM(vm *Vm, namei int32) {
	defer vm.CheckException()
	vm.NotImplemented("IMPORT_FROM", namei)
}

// Increments bytecode counter by delta.
func do_JUMP_FORWARD(vm *Vm, delta int32) {
	defer vm.CheckException()
	vm.frame.Lasti += delta
}

// If TOS is true, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_TRUE(vm *Vm, target int32) {
	defer vm.CheckException()
	if py.MakeBool(vm.POP()).(py.Bool) {
		vm.frame.Lasti = target
	}
}

// If TOS is false, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_FALSE(vm *Vm, target int32) {
	defer vm.CheckException()
	if !py.MakeBool(vm.POP()).(py.Bool) {
		vm.frame.Lasti = target
	}
}

// If TOS is true, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is false), TOS is popped.
func do_JUMP_IF_TRUE_OR_POP(vm *Vm, target int32) {
	defer vm.CheckException()
	if py.MakeBool(vm.TOP()).(py.Bool) {
		vm.frame.Lasti = target
	} else {
		vm.DROP()
	}
}

// If TOS is false, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is true), TOS is popped.
func do_JUMP_IF_FALSE_OR_POP(vm *Vm, target int32) {
	defer vm.CheckException()
	if !py.MakeBool(vm.TOP()).(py.Bool) {
		vm.frame.Lasti = target
	} else {
		vm.DROP()
	}
}

// Set bytecode counter to target.
func do_JUMP_ABSOLUTE(vm *Vm, target int32) {
	defer vm.CheckException()
	vm.frame.Lasti = target
}

// TOS is an iterator. Call its next( ) method. If this yields a new
// value, push it on the stack (leaving the iterator below it). If the
// iterator indicates it is exhausted TOS is popped, and the bytecode
// counter is incremented by delta.
func do_FOR_ITER(vm *Vm, delta int32) {
	defer func() {
		if r := recover(); r != nil {
			if vm.catchStopIteration(r) {
				vm.DROP()
				vm.frame.Lasti += delta
			}
		}
	}()
	r := py.Next(vm.TOP())
	vm.PUSH(r)
}

// Loads the global named co_names[namei] onto the stack.
func do_LOAD_GLOBAL(vm *Vm, namei int32) {
	defer vm.CheckException()
	// FIXME this is looking in local scope too - is that correct?
	fmt.Printf("LOAD_GLOBAL %v\n", vm.frame.Code.Names[namei])
	vm.PUSH(vm.frame.Lookup(vm.frame.Code.Names[namei]))
}

// Pushes a block for a loop onto the block stack. The block spans
// from the current instruction with a size of delta bytes.
func do_SETUP_LOOP(vm *Vm, delta int32) {
	defer vm.CheckException()
	vm.frame.PushBlock(SETUP_LOOP, vm.frame.Lasti+delta, len(vm.frame.Stack))
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the first except block.
func do_SETUP_EXCEPT(vm *Vm, delta int32) {
	defer vm.CheckException()
	vm.frame.PushBlock(SETUP_EXCEPT, vm.frame.Lasti+delta, len(vm.frame.Stack))
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the finally block.
func do_SETUP_FINALLY(vm *Vm, delta int32) {
	defer vm.CheckException()
	vm.frame.PushBlock(SETUP_FINALLY, vm.frame.Lasti+delta, len(vm.frame.Stack))
}

// Store a key and value pair in a dictionary. Pops the key and value
// while leaving the dictionary on the stack.
func do_STORE_MAP(vm *Vm, arg int32) {
	defer vm.CheckException()
	key := string(vm.TOP().(py.String)) // FIXME
	value := vm.SECOND()
	dictObj := vm.THIRD()
	vm.DROPN(2)
	dict := dictObj.(py.StringDict)
	dict[key] = value
}

// Pushes a reference to the local co_varnames[var_num] onto the stack.
func do_LOAD_FAST(vm *Vm, var_num int32) {
	defer vm.CheckException()
	varname := vm.frame.Code.Varnames[var_num]
	fmt.Printf("LOAD_FAST %q\n", varname)
	if value, ok := vm.frame.Locals[varname]; ok {
		vm.PUSH(value)
	} else {
		vm.SetException(py.ExceptionNewf(py.UnboundLocalError, unboundLocalErrorMsg, varname))
	}
}

// Stores TOS into the local co_varnames[var_num].
func do_STORE_FAST(vm *Vm, var_num int32) {
	defer vm.CheckException()
	vm.frame.Locals[vm.frame.Code.Varnames[var_num]] = vm.POP()
}

// Deletes local co_varnames[var_num].
func do_DELETE_FAST(vm *Vm, var_num int32) {
	defer vm.CheckException()
	varname := vm.frame.Code.Varnames[var_num]
	if _, ok := vm.frame.Locals[varname]; ok {
		delete(vm.frame.Locals, varname)
	} else {
		vm.SetException(py.ExceptionNewf(py.UnboundLocalError, unboundLocalErrorMsg, varname))
	}
}

// Name of slot for LOAD_CLOSURE / LOAD_DEREF / etc
//
// Returns name of variable and bool, true for free var, false for
// cell var
func _var_name(vm *Vm, i int32) (string, bool) {
	cellvars := vm.frame.Code.Cellvars
	if int(i) < len(cellvars) {
		return cellvars[i], false
	}
	return vm.frame.Code.Freevars[int(i)-len(cellvars)], true
}

// Pushes a reference to the cell contained in slot i of the cell and
// free variable storage. The name of the variable is co_cellvars[i]
// if i is less than the length of co_cellvars. Otherwise it is
// co_freevars[i - len(co_cellvars)].
func do_LOAD_CLOSURE(vm *Vm, i int32) {
	defer vm.CheckException()
	varname, _ := _var_name(vm, i)
	// FIXME this is making a new cell each time rather than
	// returning a reference to an old one
	vm.PUSH(py.NewCell(vm.frame.Locals[varname]))
}

// Loads the cell contained in slot i of the cell and free variable
// storage. Pushes a reference to the object the cell contains on the
// stack.
func do_LOAD_DEREF(vm *Vm, i int32) {
	defer vm.CheckException()
	res := vm.frame.Closure[i].(*py.Cell).Get()
	if res == nil {
		varname, free := _var_name(vm, i)
		if free {
			vm.SetException(py.ExceptionNewf(py.UnboundLocalError, unboundFreeErrorMsg, varname))
		} else {
			vm.SetException(py.ExceptionNewf(py.UnboundLocalError, unboundLocalErrorMsg, varname))
		}
	}
	vm.PUSH(res)
}

// Much like LOAD_DEREF but first checks the locals dictionary before
// consulting the cell. This is used for loading free variables in
// class bodies.
func do_LOAD_CLASSDEREF(vm *Vm, i int32) {
	defer vm.CheckException()
	vm.NotImplemented("LOAD_CLASSDEREF", i)
}

// Stores TOS into the cell contained in slot i of the cell and free
// variable storage.
func do_STORE_DEREF(vm *Vm, i int32) {
	defer vm.CheckException()
	vm.NotImplemented("STORE_DEREF", i)
}

// Empties the cell contained in slot i of the cell and free variable
// storage. Used by the del statement.
func do_DELETE_DEREF(vm *Vm, i int32) {
	defer vm.CheckException()
	vm.NotImplemented("DELETE_DEREF", i)
}

// Logic for the raise statement
func (vm *Vm) raise(exc, cause py.Object) {
	if exc == nil {
		// raise (with no parameters == re-raise)
		if vm.exc.Value == nil {
			vm.SetException(py.ExceptionNewf(py.RuntimeError, "No active exception to reraise"))
		} else {
			// Signal the existing exception again
			vm.exit = exitReraise
		}
	} else {
		// raise <instance>
		// raise <type>
		excException := py.MakeException(exc)
		vm.SetException(excException)
		if cause != nil {
			excException.Cause = py.MakeException(cause)
		}
	}
}

// Raises an exception. argc indicates the number of parameters to the
// raise statement, ranging from 0 to 3. The handler will find the
// traceback as TOS2, the parameter as TOS1, and the exception as TOS.
func do_RAISE_VARARGS(vm *Vm, argc int32) {
	defer vm.CheckException()
	var cause, exc py.Object
	switch argc {
	case 2:
		cause = vm.POP()
		fallthrough
	case 1:
		exc = vm.POP()
	case 0:
	default:
		panic("Bad RAISE_VARARGS argc")
	}
	vm.raise(exc, cause)
}

// Calls a function. The low byte of argc indicates the number of
// positional parameters, the high byte the number of keyword
// parameters. On the stack, the opcode finds the keyword parameters
// first. For each keyword argument, the value is on top of the
// key. Below the keyword parameters, the positional parameters are on
// the stack, with the right-most parameter on top. Below the
// parameters, the function object to call is on the stack. Pops all
// function arguments, and the function itself off the stack, and
// pushes the return value.
func do_CALL_FUNCTION(vm *Vm, argc int32) {
	defer vm.CheckException()
	// fmt.Printf("Stack: %v\n", vm.frame.Stack)
	// fmt.Printf("Locals: %v\n", vm.frame.Locals)
	// fmt.Printf("Globals: %v\n", vm.frame.Globals)
	nargs := int(argc & 0xFF)
	nkwargs := int((argc >> 8) & 0xFF)
	p, q := len(vm.frame.Stack)-2*nkwargs, len(vm.frame.Stack)
	kwargs := vm.frame.Stack[p:q]
	p, q = p-nargs, p
	args := py.Tuple(vm.frame.Stack[p:q])
	p, q = p-1, p
	fn := vm.frame.Stack[p]
	// Drop everything off the stack
	vm.frame.Stack = vm.frame.Stack[:p]
	vm.Call(fn, args, kwargs)
}

// Implementation for MAKE_FUNCTION and MAKE_CLOSURE
func _make_function(vm *Vm, argc int32, opcode byte) {
	posdefaults := argc & 0xff
	kwdefaults := (argc >> 8) & 0xff
	num_annotations := (argc >> 16) & 0x7fff
	qualname := vm.POP()
	code := vm.POP()
	function := py.NewFunction(code.(*py.Code), vm.frame.Globals, string(qualname.(py.String)))

	if opcode == MAKE_CLOSURE {
		function.Closure = vm.POP().(py.Tuple)
	}

	if num_annotations > 0 {
		names := vm.POP().(py.Tuple) // names of args with annotations
		anns := py.NewStringDict()
		name_ix := int32(len(names))
		if num_annotations != name_ix+1 {
			panic("num_annotations wrong - corrupt bytecode?")
		}
		for name_ix > 0 {
			name_ix--
			name := names[name_ix]
			value := vm.POP()
			anns[string(name.(py.String))] = value
		}
		function.Annotations = anns
	}

	if kwdefaults > 0 {
		defs := py.NewStringDict()
		for kwdefaults--; kwdefaults >= 0; kwdefaults-- {
			v := vm.POP()   // default value
			key := vm.POP() // kw only arg name
			defs[string(key.(py.String))] = v
		}
		function.KwDefaults = defs
	}

	if posdefaults > 0 {
		defs := make(py.Tuple, posdefaults)
		for posdefaults--; posdefaults >= 0; posdefaults-- {
			defs[posdefaults] = vm.POP()
		}
		function.Defaults = defs
	}

	vm.PUSH(function)
}

// Pushes a new function object on the stack. TOS is the code
// associated with the function. The function object is defined to
// have argc default parameters, which are found below TOS.
//
// FIXME these docs are slightly wrong.
func do_MAKE_FUNCTION(vm *Vm, argc int32) {
	defer vm.CheckException()
	_make_function(vm, argc, MAKE_FUNCTION)
}

// Creates a new function object, sets its func_closure slot, and
// pushes it on the stack. TOS is the code associated with the
// function, TOS1 the tuple containing cells for the closure’s free
// variables. The function also has argc default parameters, which are
// found below the cells.
func do_MAKE_CLOSURE(vm *Vm, argc int32) {
	defer vm.CheckException()
	_make_function(vm, argc, MAKE_CLOSURE)
}

// Pushes a slice object on the stack. argc must be 2 or 3. If it is
// 2, slice(TOS1, TOS) is pushed; if it is 3, slice(TOS2, TOS1, TOS)
// is pushed. See the slice( ) built-in function for more information.
func do_BUILD_SLICE(vm *Vm, argc int32) {
	defer vm.CheckException()
	vm.NotImplemented("BUILD_SLICE", argc)
}

// Prefixes any opcode which has an argument too big to fit into the
// default two bytes. ext holds two additional bytes which, taken
// together with the subsequent opcode’s argument, comprise a
// four-byte argument, ext being the two most-significant bytes.
func do_EXTENDED_ARG(vm *Vm, ext int32) {
	defer vm.CheckException()
	vm.ext = ext
	vm.extended = true
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the variable argument list, followed
// by keyword and positional arguments.
func do_CALL_FUNCTION_VAR(vm *Vm, argc int32) {
	defer vm.CheckException()
	vm.NotImplemented("CALL_FUNCTION_VAR", argc)
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by explicit keyword and positional arguments.
func do_CALL_FUNCTION_KW(vm *Vm, argc int32) {
	defer vm.CheckException()
	vm.NotImplemented("CALL_FUNCTION_KW", argc)
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by the variable-arguments tuple, followed by explicit
// keyword and positional arguments.
func do_CALL_FUNCTION_VAR_KW(vm *Vm, argc int32) {
	defer vm.CheckException()
	vm.NotImplemented("CALL_FUNCTION_VAR_KW", argc)
}

// NotImplemented
func (vm *Vm) NotImplemented(name string, arg int32) {
	fmt.Printf("%s %d NOT IMPLEMENTED\n", name, arg)
	fmt.Printf("vmstack = %#v\n", vm.frame.Stack)
	panic(py.ExceptionNewf(py.SystemError, "Opcode %s %d NOT IMPLEMENTED", name, arg))
}

// Calls function fn with args and kwargs
//
// fn can be a string in which case it will be looked up or another
// callable type such as *py.Method or *py.Function
//
// kwargsTuple is a sequence of name, value pairs
//
// The result is put on the stack
func (vm *Vm) Call(fnObj py.Object, args []py.Object, kwargsTuple []py.Object) {
	// fmt.Printf("Call %T %v with args = %v, kwargsTuple = %v\n", fnObj, fnObj, args, kwargsTuple)
	var kwargs py.StringDict
	if len(kwargsTuple) > 0 {
		// Convert kwargsTuple into dictionary
		if len(kwargsTuple)%2 != 0 {
			panic("Odd length kwargsTuple")
		}
		kwargs = py.NewStringDict()
		for i := 0; i < len(kwargsTuple); i += 2 {
			kwargs[string(kwargsTuple[i].(py.String))] = kwargsTuple[i+1]
		}
	}

	// Call the function pushing the return on the stack
	vm.PUSH(py.Call(fnObj, args, kwargs))
}

// Unwinds the stack for a block
func (vm *Vm) UnwindBlock(frame *py.Frame, block *py.TryBlock) {
	if vm.STACK_LEVEL() > block.Level {
		frame.Stack = frame.Stack[:block.Level]
	}
}

// Unwinds the stack in the presence of an exception
func (vm *Vm) UnwindExceptHandler(frame *py.Frame, block *py.TryBlock) {
	if vm.STACK_LEVEL() < block.Level+3 {
		panic("Couldn't find traceback on stack")
	} else {
		frame.Stack = frame.Stack[:block.Level+3]
	}
	vm.exc.Type = vm.POP().(*py.Type)
	vm.exc.Value = vm.POP()
	vm.exc.Traceback = vm.POP().(*py.Traceback)
}

// Run the virtual machine on a Frame object
//
// FIXME figure out how we are going to signal exceptions!
//
// Returns an Object and an error.  The error will be a py.ExceptionInfo
//
// This is the equivalent of PyEval_EvalFrame
func RunFrame(frame *py.Frame) (res py.Object, err error) {
	vm := NewVm(frame)
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		switch x := r.(type) {
	// 		case error:
	// 			err = x
	// 		case string:
	// 			err = errors.New(x)
	// 		default:
	// 			err = errors.New(fmt.Sprintf("Unknown error '%s'", x))
	// 		}
	// 		fmt.Printf("*** Exception raised %v\n", r)
	// 		// Dump the goroutine stack
	// 		debug.PrintStack()
	// 	}
	// }()

	var opcode byte
	var arg int32
	for vm.exit == exitNot {
		frame := vm.frame
		fmt.Printf("* %4d:", frame.Lasti)
		opcodes := frame.Code.Code
		opcode = opcodes[frame.Lasti]
		frame.Lasti++
		if HAS_ARG(opcode) {
			arg = int32(opcodes[frame.Lasti])
			frame.Lasti++
			arg += int32(opcodes[frame.Lasti]) << 8
			frame.Lasti++
			if vm.extended {
				arg += vm.ext << 16
			}
			fmt.Printf(" %s(%d)\n", OpCodeToName[opcode], arg)
		} else {
			fmt.Printf(" %s\n", OpCodeToName[opcode])
		}
		vm.extended = false
		jumpTable[opcode](vm, arg)
		if vm.frame != nil {
			fmt.Printf("* Stack = %#v\n", vm.frame.Stack)
			// if len(vm.frame.Stack) > 0 {
			// 	if t, ok := vm.TOP().(*py.Type); ok {
			// 		fmt.Printf(" * TOP = %#v\n", t)
			// 	}
			// }
		}

		// Something exceptional has happened - unwind the block stack
		// and find out what
		for vm.exit != exitNot && vm.frame.Block != nil {
			// Peek at the current block.
			frame := vm.frame
			b := frame.Block
			fmt.Printf("*** Unwinding %#v vm %#v\n", b, vm)

			if vm.exit == exitYield {
				return vm.result, nil
			}

			// Now we have to pop the block.
			frame.PopBlock()

			if b.Type == EXCEPT_HANDLER {
				fmt.Printf("*** EXCEPT_HANDLER\n")
				vm.UnwindExceptHandler(frame, b)
				continue
			}
			vm.UnwindBlock(frame, b)
			if b.Type == SETUP_LOOP && vm.exit == exitBreak {
				fmt.Printf("*** Loop\n")
				vm.exit = exitNot
				frame.Lasti = b.Handler
				break
			}
			if vm.exit&(exitException|exitReraise) != 0 && (b.Type == SETUP_EXCEPT || b.Type == SETUP_FINALLY) {
				fmt.Printf("*** Exception\n")
				handler := b.Handler
				// This invalidates b
				frame.PushBlock(EXCEPT_HANDLER, -1, vm.STACK_LEVEL())
				vm.PUSH(vm.old_exc.Traceback)
				vm.PUSH(vm.old_exc.Value)
				vm.PUSH(vm.exc.Type) // can be nil
				// FIXME PyErr_Fetch(&exc, &val, &tb)
				exc := vm.exc.Type
				val := vm.exc.Value
				tb := vm.exc.Traceback
				// Make the raw exception data
				// available to the handler,
				// so a program can emulate the
				// Python main loop.
				// FIXME PyErr_NormalizeException(exc, &val, &tb)
				// FIXME PyException_SetTraceback(val, tb)
				vm.exc.Type = exc
				vm.exc.Value = val
				vm.exc.Traceback = tb
				vm.PUSH(tb)
				vm.PUSH(val)
				vm.PUSH(exc)
				vm.exit = exitNot
				frame.Lasti = handler
				break
			}
			if b.Type == SETUP_FINALLY {
				if vm.exit&(exitReturn|exitContinue) != 0 {
					vm.PUSH(vm.result)
				}
				vm.PUSH(py.Int(vm.exit))
				vm.exit = exitNot
				frame.Lasti = b.Handler
				break
			}
		}
	}
	if vm.exc.Value != nil {
		return vm.result, vm.exc
	}
	return vm.result, nil
}

// Run the virtual machine on a Code object
//
// Any parameters are expected to have been decoded into locals
//
// Returns an Object and an error.  The error will be a py.ExceptionInfo
//
// This is the equivalent of PyEval_EvalCode with closure support
func Run(globals, locals py.StringDict, code *py.Code, closure py.Tuple) (res py.Object, err error) {
	frame := py.NewFrame(globals, locals, code, closure)

	// If this is a generator then make a generator object from
	// the frame and return that instead
	if code.Flags&py.CO_GENERATOR != 0 {
		return py.NewGenerator(frame), nil
	}

	return RunFrame(frame)
}

// Write the py global to avoid circular import
func init() {
	py.Run = Run
	py.RunFrame = RunFrame
}
