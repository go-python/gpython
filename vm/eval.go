// Evaluate opcodes
package vm

import (
	"errors"
	"fmt"
	"github.com/ncw/gpython/py"
)

// Stack operations
func (vm *Vm) STACK_LEVEL() int             { return len(vm.stack) }
func (vm *Vm) EMPTY() bool                  { return len(vm.stack) == 0 }
func (vm *Vm) TOP() py.Object               { return vm.stack[len(vm.stack)-1] }
func (vm *Vm) SECOND() py.Object            { return vm.stack[len(vm.stack)-2] }
func (vm *Vm) THIRD() py.Object             { return vm.stack[len(vm.stack)-3] }
func (vm *Vm) FOURTH() py.Object            { return vm.stack[len(vm.stack)-4] }
func (vm *Vm) PEEK(n int) py.Object         { return vm.stack[len(vm.stack)-n] }
func (vm *Vm) SET_TOP(v py.Object)          { vm.stack[len(vm.stack)-1] = v }
func (vm *Vm) SET_SECOND(v py.Object)       { vm.stack[len(vm.stack)-2] = v }
func (vm *Vm) SET_THIRD(v py.Object)        { vm.stack[len(vm.stack)-3] = v }
func (vm *Vm) SET_FOURTH(v py.Object)       { vm.stack[len(vm.stack)-4] = v }
func (vm *Vm) SET_VALUE(n int, v py.Object) { vm.stack[len(vm.stack)-(n)] = (v) }
func (vm *Vm) DROPN(n int)                  { vm.stack = vm.stack[:len(vm.stack)-n] }

// Pop from top of vm stack
func (vm *Vm) POP() py.Object {
	// FIXME what if empty?
	out := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return out
}

// Push to top of vm stack
func (vm *Vm) PUSH(obj py.Object) {
	vm.stack = append(vm.stack, obj)
}

// Illegal instruction
func do_ILLEGAL(vm *Vm, arg int32) {
	panic("Illegal opcode")
}

// Do nothing code. Used as a placeholder by the bytecode optimizer.
func do_NOP(vm *Vm, arg int32) {
}

// Removes the top-of-stack (TOS) item.
func do_POP_TOP(vm *Vm, arg int32) {
	vm.DROPN(1)
}

// Swaps the two top-most stack items.
func do_ROT_TWO(vm *Vm, arg int32) {
	top := vm.TOP()
	second := vm.SECOND()
	vm.SET_TOP(second)
	vm.SET_SECOND(top)
}

// Lifts second and third stack item one position up, moves top down
// to position three.
func do_ROT_THREE(vm *Vm, arg int32) {
	top := vm.TOP()
	second := vm.SECOND()
	third := vm.THIRD()
	vm.SET_TOP(second)
	vm.SET_SECOND(third)
	vm.SET_THIRD(top)
}

// Duplicates the reference on top of the stack.
func do_DUP_TOP(vm *Vm, arg int32) {
	vm.PUSH(vm.TOP())
}

// Duplicates the top two reference on top of the stack.
func do_DUP_TOP_TWO(vm *Vm, arg int32) {
	top := vm.TOP()
	second := vm.SECOND()
	vm.PUSH(second)
	vm.PUSH(top)
}

// Unary Operations take the top of the stack, apply the operation,
// and push the result back on the stack.

// Implements TOS = +TOS.
func do_UNARY_POSITIVE(vm *Vm, arg int32) {
	vm.NotImplemented("UNARY_POSITIVE", arg)
}

// Implements TOS = -TOS.
func do_UNARY_NEGATIVE(vm *Vm, arg int32) {
	vm.NotImplemented("UNARY_NEGATIVE", arg)
}

// Implements TOS = not TOS.
func do_UNARY_NOT(vm *Vm, arg int32) {
	vm.NotImplemented("UNARY_NOT", arg)
}

// Implements TOS = ~TOS.
func do_UNARY_INVERT(vm *Vm, arg int32) {
	vm.NotImplemented("UNARY_INVERT", arg)
}

// Implements TOS = iter(TOS).
func do_GET_ITER(vm *Vm, arg int32) {
	vm.NotImplemented("GET_ITER", arg)
}

// Binary operations remove the top of the stack (TOS) and the second
// top-most stack item (TOS1) from the stack. They perform the
// operation, and put the result back on the stack.

// Implements TOS = TOS1 ** TOS.
func do_BINARY_POWER(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_POWER", arg)
}

// Implements TOS = TOS1 * TOS.
func do_BINARY_MULTIPLY(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_MULTIPLY", arg)
}

// Implements TOS = TOS1 // TOS.
func do_BINARY_FLOOR_DIVIDE(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_FLOOR_DIVIDE", arg)
}

// Implements TOS = TOS1 / TOS when from __future__ import division is
// in effect.
func do_BINARY_TRUE_DIVIDE(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_TRUE_DIVIDE", arg)
}

// Implements TOS = TOS1 % TOS.
func do_BINARY_MODULO(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_MODULO", arg)
}

// Implements TOS = TOS1 + TOS.
func do_BINARY_ADD(vm *Vm, arg int32) {
	b := vm.POP()
	a := vm.POP()
	vm.PUSH(py.Add(a, b))
}

// Implements TOS = TOS1 - TOS.
func do_BINARY_SUBTRACT(vm *Vm, arg int32) {
	b := vm.POP()
	a := vm.POP()
	vm.PUSH(py.Sub(a, b))
}

// Implements TOS = TOS1[TOS].
func do_BINARY_SUBSCR(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_SUBSCR", arg)
}

// Implements TOS = TOS1 << TOS.
func do_BINARY_LSHIFT(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_LSHIFT", arg)
}

// Implements TOS = TOS1 >> TOS.
func do_BINARY_RSHIFT(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_RSHIFT", arg)
}

// Implements TOS = TOS1 & TOS.
func do_BINARY_AND(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_AND", arg)
}

// Implements TOS = TOS1 ^ TOS.
func do_BINARY_XOR(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_XOR", arg)
}

// Implements TOS = TOS1 | TOS.
func do_BINARY_OR(vm *Vm, arg int32) {
	vm.NotImplemented("BINARY_OR", arg)
}

// In-place operations are like binary operations, in that they remove
// TOS and TOS1, and push the result back on the stack, but the
// operation is done in-place when TOS1 supports it, and the resulting
// TOS may be (but does not have to be) the original TOS1.

// Implements in-place TOS = TOS1 ** TOS.
func do_INPLACE_POWER(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_POWER", arg)
}

// Implements in-place TOS = TOS1 * TOS.
func do_INPLACE_MULTIPLY(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_MULTIPLY", arg)
}

// Implements in-place TOS = TOS1 // TOS.
func do_INPLACE_FLOOR_DIVIDE(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_FLOOR_DIVIDE", arg)
}

// Implements in-place TOS = TOS1 / TOS when from __future__ import
// division is in effect.
func do_INPLACE_TRUE_DIVIDE(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_TRUE_DIVIDE", arg)
}

// Implements in-place TOS = TOS1 % TOS.
func do_INPLACE_MODULO(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_MODULO", arg)
}

// Implements in-place TOS = TOS1 + TOS.
func do_INPLACE_ADD(vm *Vm, arg int32) {
	b := vm.POP()
	a := vm.POP()
	vm.PUSH(py.IAdd(a, b))
}

// Implements in-place TOS = TOS1 - TOS.
func do_INPLACE_SUBTRACT(vm *Vm, arg int32) {
	b := vm.POP()
	a := vm.POP()
	vm.PUSH(py.ISub(a, b))
}

// Implements in-place TOS = TOS1 << TOS.
func do_INPLACE_LSHIFT(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_LSHIFT", arg)
}

// Implements in-place TOS = TOS1 >> TOS.
func do_INPLACE_RSHIFT(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_RSHIFT", arg)
}

// Implements in-place TOS = TOS1 & TOS.
func do_INPLACE_AND(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_AND", arg)
}

// Implements in-place TOS = TOS1 ^ TOS.
func do_INPLACE_XOR(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_XOR", arg)
}

// Implements in-place TOS = TOS1 | TOS.
func do_INPLACE_OR(vm *Vm, arg int32) {
	vm.NotImplemented("INPLACE_OR", arg)
}

// Implements TOS1[TOS] = TOS2.
func do_STORE_SUBSCR(vm *Vm, arg int32) {
	vm.NotImplemented("STORE_SUBSCR", arg)
}

// Implements del TOS1[TOS].
func do_DELETE_SUBSCR(vm *Vm, arg int32) {
	vm.NotImplemented("DELETE_SUBSCR", arg)
}

// Miscellaneous opcodes.

// Implements the expression statement for the interactive mode. TOS
// is removed from the stack and printed. In non-interactive mode, an
// expression statement is terminated with POP_STACK.
func do_PRINT_EXPR(vm *Vm, arg int32) {
	vm.NotImplemented("PRINT_EXPR", arg)
}

// Terminates a loop due to a break statement.
func do_BREAK_LOOP(vm *Vm, arg int32) {
	vm.NotImplemented("BREAK_LOOP", arg)
}

// Continues a loop due to a continue statement. target is the address
// to jump to (which should be a FOR_ITER instruction).
func do_CONTINUE_LOOP(vm *Vm, target int32) {
	vm.NotImplemented("CONTINUE_LOOP", target)
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
	vm.NotImplemented("UNPACK_EX", counts)
}

// Calls set.add(TOS1[-i], TOS). Used to implement set comprehensions.
func do_SET_ADD(vm *Vm, i int32) {
	vm.NotImplemented("SET_ADD", i)
}

// Calls list.append(TOS[-i], TOS). Used to implement list
// comprehensions. While the appended value is popped off, the list
// object remains on the stack so that it is available for further
// iterations of the loop.
func do_LIST_APPEND(vm *Vm, i int32) {
	vm.NotImplemented("LIST_APPEND", i)
}

// Calls dict.setitem(TOS1[-i], TOS, TOS1). Used to implement dict comprehensions.
func do_MAP_ADD(vm *Vm, i int32) {
	vm.NotImplemented("MAP_ADD", i)
}

// Returns with TOS to the caller of the function.
func do_RETURN_VALUE(vm *Vm, arg int32) {
	vm.PopFrame()
}

// Pops TOS and delegates to it as a subiterator from a generator.
func do_YIELD_FROM(vm *Vm, arg int32) {
	vm.NotImplemented("YIELD_FROM", arg)
}

// Pops TOS and yields it from a generator.
func do_YIELD_VALUE(vm *Vm, arg int32) {
	vm.NotImplemented("YIELD_VALUE", arg)
}

// Loads all symbols not starting with '_' directly from the module
// TOS to the local namespace. The module is popped after loading all
// names. This opcode implements from module import *.
func do_IMPORT_STAR(vm *Vm, arg int32) {
	vm.NotImplemented("IMPORT_STAR", arg)
}

// Removes one block from the block stack. Per frame, there is a stack
// of blocks, denoting nested loops, try statements, and such.
func do_POP_BLOCK(vm *Vm, arg int32) {
	vm.NotImplemented("POP_BLOCK", arg)
}

// Removes one block from the block stack. The popped block must be an
// exception handler block, as implicitly created when entering an
// except handler. In addition to popping extraneous values from the
// frame stack, the last three popped values are used to restore the
// exception state.
func do_POP_EXCEPT(vm *Vm, arg int32) {
	vm.NotImplemented("POP_EXCEPT", arg)
}

// Terminates a finally clause. The interpreter recalls whether the
// exception has to be re-raised, or whether the function returns, and
// continues with the outer-next block.
func do_END_FINALLY(vm *Vm, arg int32) {
	vm.NotImplemented("END_FINALLY", arg)
}

// Creates a new class object. TOS is the methods dictionary, TOS1 the
// tuple of the names of the base classes, and TOS2 the class name.
func do_LOAD_BUILD_CLASS(vm *Vm, arg int32) {
	vm.NotImplemented("LOAD_BUILD_CLASS", arg)
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
	vm.NotImplemented("WITH_CLEANUP", arg)
}

// All of the following opcodes expect arguments. An argument is two bytes, with the more significant byte last.

// Implements name = TOS. namei is the index of name in the attribute
// co_names of the code object. The compiler tries to use STORE_FAST
// or STORE_GLOBAL if possible.
func do_STORE_NAME(vm *Vm, namei int32) {
	vm.frame.Locals[vm.frame.Code.Names[namei]] = vm.POP()
}

// Implements del name, where namei is the index into co_names
// attribute of the code object.
func do_DELETE_NAME(vm *Vm, namei int32) {
	vm.NotImplemented("DELETE_NAME", namei)
}

// Unpacks TOS into count individual values, which are put onto the
// stack right-to-left.
func do_UNPACK_SEQUENCE(vm *Vm, count int32) {
	vm.NotImplemented("UNPACK_SEQUENCE", count)
}

// Implements TOS.name = TOS1, where namei is the index of name in
// co_names.
func do_STORE_ATTR(vm *Vm, namei int32) {
	vm.NotImplemented("STORE_ATTR", namei)
}

// Implements del TOS.name, using namei as index into co_names.
func do_DELETE_ATTR(vm *Vm, namei int32) {
	vm.NotImplemented("DELETE_ATTR", namei)
}

// Works as STORE_NAME, but stores the name as a global.
func do_STORE_GLOBAL(vm *Vm, namei int32) {
	vm.NotImplemented("STORE_GLOBAL", namei)
}

// Works as DELETE_NAME, but deletes a global name.
func do_DELETE_GLOBAL(vm *Vm, namei int32) {
	vm.NotImplemented("DELETE_GLOBAL", namei)
}

// Pushes co_consts[consti] onto the stack.
func do_LOAD_CONST(vm *Vm, consti int32) {
	vm.PUSH(vm.frame.Code.Consts[consti])
	// fmt.Printf("LOAD_CONST %v\n", vm.TOP())
}

// Pushes the value associated with co_names[namei] onto the stack.
func do_LOAD_NAME(vm *Vm, namei int32) {
	vm.PUSH(py.String(vm.frame.Code.Names[namei]))
}

// Creates a tuple consuming count items from the stack, and pushes
// the resulting tuple onto the stack.
func do_BUILD_TUPLE(vm *Vm, count int32) {
	vm.NotImplemented("BUILD_TUPLE", count)
}

// Works as BUILD_TUPLE, but creates a set.
func do_BUILD_SET(vm *Vm, count int32) {
	vm.NotImplemented("BUILD_SET", count)
}

// Works as BUILD_TUPLE, but creates a list.
func do_BUILD_LIST(vm *Vm, count int32) {
	vm.NotImplemented("BUILD_LIST", count)
}

// Pushes a new dictionary object onto the stack. The dictionary is
// pre-sized to hold count entries.
func do_BUILD_MAP(vm *Vm, count int32) {
	vm.NotImplemented("BUILD_MAP", count)
}

// Replaces TOS with getattr(TOS, co_names[namei]).
func do_LOAD_ATTR(vm *Vm, namei int32) {
	vm.NotImplemented("LOAD_ATTR", namei)
}

// Performs a Boolean operation. The operation name can be found in
// cmp_op[opname].
func do_COMPARE_OP(vm *Vm, opname int32) {
	vm.NotImplemented("COMPARE_OP", opname)
}

// Imports the module co_names[namei]. TOS and TOS1 are popped and
// provide the fromlist and level arguments of __import__( ). The
// module object is pushed onto the stack. The current namespace is
// not affected: for a proper import statement, a subsequent
// STORE_FAST instruction modifies the namespace.
func do_IMPORT_NAME(vm *Vm, namei int32) {
	vm.NotImplemented("IMPORT_NAME", namei)
}

// Loads the attribute co_names[namei] from the module found in
// TOS. The resulting object is pushed onto the stack, to be
// subsequently stored by a STORE_FAST instruction.
func do_IMPORT_FROM(vm *Vm, namei int32) {
	vm.NotImplemented("IMPORT_FROM", namei)
}

// Increments bytecode counter by delta.
func do_JUMP_FORWARD(vm *Vm, delta int32) {
	vm.NotImplemented("JUMP_FORWARD", delta)
}

// If TOS is true, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_TRUE(vm *Vm, target int32) {
	vm.NotImplemented("POP_JUMP_IF_TRUE", target)
}

// If TOS is false, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_FALSE(vm *Vm, target int32) {
	vm.NotImplemented("POP_JUMP_IF_FALSE", target)
}

// If TOS is true, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is false), TOS is popped.
func do_JUMP_IF_TRUE_OR_POP(vm *Vm, target int32) {
	vm.NotImplemented("JUMP_IF_TRUE_OR_POP", target)
}

// If TOS is false, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is true), TOS is popped.
func do_JUMP_IF_FALSE_OR_POP(vm *Vm, target int32) {
	vm.NotImplemented("JUMP_IF_FALSE_OR_POP", target)
}

// Set bytecode counter to target.
func do_JUMP_ABSOLUTE(vm *Vm, target int32) {
	vm.NotImplemented("JUMP_ABSOLUTE", target)
}

// TOS is an iterator. Call its next( ) method. If this yields a new
// value, push it on the stack (leaving the iterator below it). If the
// iterator indicates it is exhausted TOS is popped, and the bytecode
// counter is incremented by delta.
func do_FOR_ITER(vm *Vm, delta int32) {
	vm.NotImplemented("FOR_ITER", delta)
}

// Loads the global named co_names[namei] onto the stack.
func do_LOAD_GLOBAL(vm *Vm, namei int32) {
	// FIXME this is looking in local scope too - is that correct?
	vm.PUSH(vm.frame.Lookup(vm.frame.Code.Names[namei]))
}

// Pushes a block for a loop onto the block stack. The block spans
// from the current instruction with a size of delta bytes.
func do_SETUP_LOOP(vm *Vm, delta int32) {
	vm.NotImplemented("SETUP_LOOP", delta)
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the first except block.
func do_SETUP_EXCEPT(vm *Vm, delta int32) {
	vm.NotImplemented("SETUP_EXCEPT", delta)
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the finally block.
func do_SETUP_FINALLY(vm *Vm, delta int32) {
	vm.NotImplemented("SETUP_FINALLY", delta)
}

// Store a key and value pair in a dictionary. Pops the key and value
// while leaving the dictionary on the stack.
func do_STORE_MAP(vm *Vm, arg int32) {
	vm.NotImplemented("STORE_MAP", arg)
}

// Pushes a reference to the local co_varnames[var_num] onto the stack.
func do_LOAD_FAST(vm *Vm, var_num int32) {
	vm.PUSH(vm.frame.Locals[vm.frame.Code.Varnames[var_num]])
}

// Stores TOS into the local co_varnames[var_num].
func do_STORE_FAST(vm *Vm, var_num int32) {
	vm.frame.Locals[vm.frame.Code.Varnames[var_num]] = vm.POP()
}

// Deletes local co_varnames[var_num].
func do_DELETE_FAST(vm *Vm, var_num int32) {
	vm.NotImplemented("DELETE_FAST", var_num)
}

// Pushes a reference to the cell contained in slot i of the cell and
// free variable storage. The name of the variable is co_cellvars[i]
// if i is less than the length of co_cellvars. Otherwise it is
// co_freevars[i - len(co_cellvars)].
func do_LOAD_CLOSURE(vm *Vm, i int32) {
	vm.NotImplemented("LOAD_CLOSURE", i)
}

// Loads the cell contained in slot i of the cell and free variable
// storage. Pushes a reference to the object the cell contains on the
// stack.
func do_LOAD_DEREF(vm *Vm, i int32) {
	vm.NotImplemented("LOAD_DEREF", i)
}

// Much like LOAD_DEREF but first checks the locals dictionary before
// consulting the cell. This is used for loading free variables in
// class bodies.
func do_LOAD_CLASSDEREF(vm *Vm, i int32) {
	vm.NotImplemented("LOAD_CLASSDEREF", i)
}

// Stores TOS into the cell contained in slot i of the cell and free
// variable storage.
func do_STORE_DEREF(vm *Vm, i int32) {
	vm.NotImplemented("STORE_DEREF", i)
}

// Empties the cell contained in slot i of the cell and free variable
// storage. Used by the del statement.
func do_DELETE_DEREF(vm *Vm, i int32) {
	vm.NotImplemented("DELETE_DEREF", i)
}

// Raises an exception. argc indicates the number of parameters to the
// raise statement, ranging from 0 to 3. The handler will find the
// traceback as TOS2, the parameter as TOS1, and the exception as TOS.
func do_RAISE_VARARGS(vm *Vm, argc int32) {
	vm.NotImplemented("RAISE_VARARGS", argc)
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
	fmt.Printf("Stack: %v\n", vm.stack)
	fmt.Printf("Locals: %v\n", vm.frame.Locals)
	fmt.Printf("Globals: %v\n", vm.frame.Globals)
	nargs := int(argc & 0xFF)
	nkwargs := int((argc >> 8) & 0xFF)
	p, q := len(vm.stack)-2*nkwargs, len(vm.stack)
	kwargs := vm.stack[p:q]
	p, q = p-nargs, p
	args := py.Tuple(vm.stack[p:q])
	p, q = p-1, p
	fn := vm.stack[p]
	// Drop everything off the stack
	vm.stack = vm.stack[:q]
	vm.Call(fn, args, kwargs)
}

// Pushes a new function object on the stack. TOS is the code
// associated with the function. The function object is defined to
// have argc default parameters, which are found below TOS.
//
// FIXME these docs are slightly wrong.
func do_MAKE_FUNCTION(vm *Vm, argc int32) {
	posdefaults := argc & 0xff
	kwdefaults := (argc >> 8) & 0xff
	num_annotations := (argc >> 16) & 0x7fff
	qualname := vm.POP()
	code := vm.POP()
	function := py.NewFunction(code.(*py.Code), vm.frame.Globals, string(qualname.(py.String)))

	// FIXME share code with MAKE_CLOSURE
	// if opcode == MAKE_CLOSURE {
	// 	function.Closure = vm.POP();
	// }

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

// Creates a new function object, sets its func_closure slot, and
// pushes it on the stack. TOS is the code associated with the
// function, TOS1 the tuple containing cells for the closure’s free
// variables. The function also has argc default parameters, which are
// found below the cells.
func do_MAKE_CLOSURE(vm *Vm, argc int32) {
	vm.NotImplemented("MAKE_CLOSURE", argc)
	// see MAKE_FUNCTION
}

// Pushes a slice object on the stack. argc must be 2 or 3. If it is
// 2, slice(TOS1, TOS) is pushed; if it is 3, slice(TOS2, TOS1, TOS)
// is pushed. See the slice( ) built-in function for more information.
func do_BUILD_SLICE(vm *Vm, argc int32) {
	vm.NotImplemented("BUILD_SLICE", argc)
}

// Prefixes any opcode which has an argument too big to fit into the
// default two bytes. ext holds two additional bytes which, taken
// together with the subsequent opcode’s argument, comprise a
// four-byte argument, ext being the two most-significant bytes.
func do_EXTENDED_ARG(vm *Vm, ext int32) {
	vm.ext = ext
	vm.extended = true
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the variable argument list, followed
// by keyword and positional arguments.
func do_CALL_FUNCTION_VAR(vm *Vm, argc int32) {
	vm.NotImplemented("CALL_FUNCTION_VAR", argc)
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by explicit keyword and positional arguments.
func do_CALL_FUNCTION_KW(vm *Vm, argc int32) {
	vm.NotImplemented("CALL_FUNCTION_KW", argc)
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by the variable-arguments tuple, followed by explicit
// keyword and positional arguments.
func do_CALL_FUNCTION_VAR_KW(vm *Vm, argc int32) {
	vm.NotImplemented("CALL_FUNCTION_VAR_KW", argc)
}

// NotImplemented
func (vm *Vm) NotImplemented(name string, arg int32) {
	fmt.Printf("%s %d NOT IMPLEMENTED\n", name, arg)
}

// Calls function fn with args and kwargs
//
// fn can be a string in which case it will be looked up or another callable type
//
// kwargs is a sequence of name, value pairs
//
// The result is put on the stack
func (vm *Vm) Call(fnObj py.Object, args []py.Object, kwargs []py.Object) {
	fmt.Printf("Call %T %v with args = %v, kwargs = %v\n", fnObj, fnObj, args, kwargs)
	var kwargsd py.StringDict
	if len(kwargs) > 0 {
		// Convert kwargs into dictionary
		if len(kwargs)%2 != 0 {
			panic("Odd length kwargs")
		}
		kwargsd = py.NewStringDict()
		for i := 0; i < len(kwargs); i += 2 {
			kwargsd[string(kwargs[i].(py.String))] = kwargs[i+1]
		}
	}
try_again:
	switch fn := fnObj.(type) {
	case py.String:
		fnObj = vm.frame.Lookup(string(fn))
		goto try_again
	case *py.Method:
		self := py.None // FIXME should be the module
		if kwargsd != nil {
			vm.PUSH(fn.CallWithKeywords(self, args, kwargsd))
		} else {
			vm.PUSH(fn.Call(self, args))
		}
	case *py.Function:
		var locals py.StringDict
		if kwargsd != nil {
			locals = fn.LocalsForCallWithKeywords(args, kwargsd)
		} else {
			locals = fn.LocalsForCall(args)
		}
		vm.PushFrame(vm.frame.Globals, locals, fn.Code)
	default:
		// FIXME should be TypeError
		panic(fmt.Sprintf("TypeError: '%s' object is not callable", fnObj.Type().Name))
	}
}

// Make a new Frame with globals, locals and Code on the frames stack
func (vm *Vm) PushFrame(globals, locals py.StringDict, code *py.Code) {
	frame := py.Frame{
		Globals:  globals,
		Locals:   locals,
		Code:     code,
		Builtins: py.Builtins.Globals,
	}
	vm.frames = append(vm.frames, frame)
	vm.frame = &frame
}

// Drop the current frame
func (vm *Vm) PopFrame() {
	vm.frames = vm.frames[:len(vm.frames)-1]
	if len(vm.frames) > 0 {
		vm.frame = &vm.frames[len(vm.frames)-1]
	} else {
		vm.frame = nil
	}
}

// Run the virtual machine on the code object in the module
//
// FIXME figure out how we are going to signal exceptions!
//
// Any parameters are expected to have been decoded into locals
func Run(globals, locals py.StringDict, code *py.Code) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				err = x
			case string:
				err = errors.New(x)
			default:
				err = errors.New(fmt.Sprintf("Unknown error '%s'", x))
			}
		}
	}()
	vm := NewVm()
	vm.PushFrame(globals, locals, code)

	var opcode byte
	var arg int32
	for vm.frame != nil {
		frame := vm.frame
		opcodes := frame.Code.Code
		opcode = opcodes[frame.Lasti]
		frame.Lasti++
		if HAS_ARG(opcode) {
			arg = int32(opcodes[frame.Lasti])
			frame.Lasti++
			arg += int32(opcodes[frame.Lasti] << 8)
			frame.Lasti++
			if vm.extended {
				arg += vm.ext << 16
			}
			fmt.Printf("* %s(%d)\n", OpCodeToName[opcode], arg)
		} else {
			fmt.Printf("* %s\n", OpCodeToName[opcode])
		}
		vm.extended = false
		jumpTable[opcode](vm, arg)
	}
	if len(vm.stack) != 1 {
		fmt.Printf("vmstack = %v\n", vm.stack)
		panic("vm stack should only have 1 entry on at this point")
	}
	// return vm.POP()
	return nil
}
