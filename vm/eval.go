// Evaluate opcodes
package vm

import (
	"errors"
	"fmt"
	"github.com/ncw/gpython/py"
)

// Globals
var (
	jumpTable [256]func(*Vm, int32)
)

// Initialise jump table
func init() {
	for i := range jumpTable {
		jumpTable[i] = do_ILLEGAL
	}
	jumpTable[POP_TOP] = do_POP_TOP
	jumpTable[ROT_TWO] = do_ROT_TWO
	jumpTable[ROT_THREE] = do_ROT_THREE
	jumpTable[DUP_TOP] = do_DUP_TOP
	jumpTable[DUP_TOP_TWO] = do_DUP_TOP_TWO
	jumpTable[NOP] = do_NOP

	jumpTable[UNARY_POSITIVE] = do_UNARY_POSITIVE
	jumpTable[UNARY_NEGATIVE] = do_UNARY_NEGATIVE
	jumpTable[UNARY_NOT] = do_UNARY_NOT

	jumpTable[UNARY_INVERT] = do_UNARY_INVERT

	jumpTable[BINARY_POWER] = do_BINARY_POWER

	jumpTable[BINARY_MULTIPLY] = do_BINARY_MULTIPLY

	jumpTable[BINARY_MODULO] = do_BINARY_MODULO
	jumpTable[BINARY_ADD] = do_BINARY_ADD
	jumpTable[BINARY_SUBTRACT] = do_BINARY_SUBTRACT
	jumpTable[BINARY_SUBSCR] = do_BINARY_SUBSCR
	jumpTable[BINARY_FLOOR_DIVIDE] = do_BINARY_FLOOR_DIVIDE
	jumpTable[BINARY_TRUE_DIVIDE] = do_BINARY_TRUE_DIVIDE
	jumpTable[INPLACE_FLOOR_DIVIDE] = do_INPLACE_FLOOR_DIVIDE
	jumpTable[INPLACE_TRUE_DIVIDE] = do_INPLACE_TRUE_DIVIDE

	jumpTable[STORE_MAP] = do_STORE_MAP
	jumpTable[INPLACE_ADD] = do_INPLACE_ADD
	jumpTable[INPLACE_SUBTRACT] = do_INPLACE_SUBTRACT
	jumpTable[INPLACE_MULTIPLY] = do_INPLACE_MULTIPLY

	jumpTable[INPLACE_MODULO] = do_INPLACE_MODULO
	jumpTable[STORE_SUBSCR] = do_STORE_SUBSCR
	jumpTable[DELETE_SUBSCR] = do_DELETE_SUBSCR

	jumpTable[BINARY_LSHIFT] = do_BINARY_LSHIFT
	jumpTable[BINARY_RSHIFT] = do_BINARY_RSHIFT
	jumpTable[BINARY_AND] = do_BINARY_AND
	jumpTable[BINARY_XOR] = do_BINARY_XOR
	jumpTable[BINARY_OR] = do_BINARY_OR
	jumpTable[INPLACE_POWER] = do_INPLACE_POWER
	jumpTable[GET_ITER] = do_GET_ITER
	jumpTable[PRINT_EXPR] = do_PRINT_EXPR
	jumpTable[LOAD_BUILD_CLASS] = do_LOAD_BUILD_CLASS
	jumpTable[YIELD_FROM] = do_YIELD_FROM

	jumpTable[INPLACE_LSHIFT] = do_INPLACE_LSHIFT
	jumpTable[INPLACE_RSHIFT] = do_INPLACE_RSHIFT
	jumpTable[INPLACE_AND] = do_INPLACE_AND
	jumpTable[INPLACE_XOR] = do_INPLACE_XOR
	jumpTable[INPLACE_OR] = do_INPLACE_OR
	jumpTable[BREAK_LOOP] = do_BREAK_LOOP
	jumpTable[WITH_CLEANUP] = do_WITH_CLEANUP

	jumpTable[RETURN_VALUE] = do_RETURN_VALUE
	jumpTable[IMPORT_STAR] = do_IMPORT_STAR

	jumpTable[YIELD_VALUE] = do_YIELD_VALUE
	jumpTable[POP_BLOCK] = do_POP_BLOCK
	jumpTable[END_FINALLY] = do_END_FINALLY
	jumpTable[POP_EXCEPT] = do_POP_EXCEPT

	jumpTable[STORE_NAME] = do_STORE_NAME
	jumpTable[DELETE_NAME] = do_DELETE_NAME
	jumpTable[UNPACK_SEQUENCE] = do_UNPACK_SEQUENCE
	jumpTable[FOR_ITER] = do_FOR_ITER
	jumpTable[UNPACK_EX] = do_UNPACK_EX

	jumpTable[STORE_ATTR] = do_STORE_ATTR
	jumpTable[DELETE_ATTR] = do_DELETE_ATTR
	jumpTable[STORE_GLOBAL] = do_STORE_GLOBAL
	jumpTable[DELETE_GLOBAL] = do_DELETE_GLOBAL

	jumpTable[LOAD_CONST] = do_LOAD_CONST
	jumpTable[LOAD_NAME] = do_LOAD_NAME
	jumpTable[BUILD_TUPLE] = do_BUILD_TUPLE
	jumpTable[BUILD_LIST] = do_BUILD_LIST
	jumpTable[BUILD_SET] = do_BUILD_SET
	jumpTable[BUILD_MAP] = do_BUILD_MAP
	jumpTable[LOAD_ATTR] = do_LOAD_ATTR
	jumpTable[COMPARE_OP] = do_COMPARE_OP
	jumpTable[IMPORT_NAME] = do_IMPORT_NAME
	jumpTable[IMPORT_FROM] = do_IMPORT_FROM

	jumpTable[JUMP_FORWARD] = do_JUMP_FORWARD
	jumpTable[JUMP_IF_FALSE_OR_POP] = do_JUMP_IF_FALSE_OR_POP
	jumpTable[JUMP_IF_TRUE_OR_POP] = do_JUMP_IF_TRUE_OR_POP
	jumpTable[JUMP_ABSOLUTE] = do_JUMP_ABSOLUTE
	jumpTable[POP_JUMP_IF_FALSE] = do_POP_JUMP_IF_FALSE
	jumpTable[POP_JUMP_IF_TRUE] = do_POP_JUMP_IF_TRUE

	jumpTable[LOAD_GLOBAL] = do_LOAD_GLOBAL

	jumpTable[CONTINUE_LOOP] = do_CONTINUE_LOOP
	jumpTable[SETUP_LOOP] = do_SETUP_LOOP
	jumpTable[SETUP_EXCEPT] = do_SETUP_EXCEPT
	jumpTable[SETUP_FINALLY] = do_SETUP_FINALLY

	jumpTable[LOAD_FAST] = do_LOAD_FAST
	jumpTable[STORE_FAST] = do_STORE_FAST
	jumpTable[DELETE_FAST] = do_DELETE_FAST

	jumpTable[RAISE_VARARGS] = do_RAISE_VARARGS
	jumpTable[CALL_FUNCTION] = do_CALL_FUNCTION
	jumpTable[MAKE_FUNCTION] = do_MAKE_FUNCTION
	jumpTable[BUILD_SLICE] = do_BUILD_SLICE

	jumpTable[MAKE_CLOSURE] = do_MAKE_CLOSURE
	jumpTable[LOAD_CLOSURE] = do_LOAD_CLOSURE
	jumpTable[LOAD_DEREF] = do_LOAD_DEREF
	jumpTable[STORE_DEREF] = do_STORE_DEREF
	jumpTable[DELETE_DEREF] = do_DELETE_DEREF

	jumpTable[CALL_FUNCTION_VAR] = do_CALL_FUNCTION_VAR
	jumpTable[CALL_FUNCTION_KW] = do_CALL_FUNCTION_KW
	jumpTable[CALL_FUNCTION_VAR_KW] = do_CALL_FUNCTION_VAR_KW

	jumpTable[SETUP_WITH] = do_SETUP_WITH

	jumpTable[EXTENDED_ARG] = do_EXTENDED_ARG

	jumpTable[LIST_APPEND] = do_LIST_APPEND
	jumpTable[SET_ADD] = do_SET_ADD
	jumpTable[MAP_ADD] = do_MAP_ADD

	jumpTable[LOAD_CLASSDEREF] = do_LOAD_CLASSDEREF
}

// Virtual machine state
type Vm struct {
	// Object stack
	stack []py.Object
	// Current code object
	co *py.Code
	// Whether ext should be added to the next arg
	extended bool
	// 16 bit extension for argument for next opcode
	ext int32
	// Whether we should exit
	exit bool
}

// Make a new VM
func NewVm() *Vm {
	vm := new(Vm)
	vm.stack = make([]py.Object, 0, 1024)
	return vm
}

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
}

// Implements TOS = -TOS.
func do_UNARY_NEGATIVE(vm *Vm, arg int32) {
}

// Implements TOS = not TOS.
func do_UNARY_NOT(vm *Vm, arg int32) {
}

// Implements TOS = ~TOS.
func do_UNARY_INVERT(vm *Vm, arg int32) {
}

// Implements TOS = iter(TOS).
func do_GET_ITER(vm *Vm, arg int32) {
}

// Binary operations remove the top of the stack (TOS) and the second
// top-most stack item (TOS1) from the stack. They perform the
// operation, and put the result back on the stack.

// Implements TOS = TOS1 ** TOS.
func do_BINARY_POWER(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 * TOS.
func do_BINARY_MULTIPLY(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 // TOS.
func do_BINARY_FLOOR_DIVIDE(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 / TOS when from __future__ import division is
// in effect.
func do_BINARY_TRUE_DIVIDE(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 % TOS.
func do_BINARY_MODULO(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 + TOS.
func do_BINARY_ADD(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 - TOS.
func do_BINARY_SUBTRACT(vm *Vm, arg int32) {
}

// Implements TOS = TOS1[TOS].
func do_BINARY_SUBSCR(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 << TOS.
func do_BINARY_LSHIFT(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 >> TOS.
func do_BINARY_RSHIFT(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 & TOS.
func do_BINARY_AND(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 ^ TOS.
func do_BINARY_XOR(vm *Vm, arg int32) {
}

// Implements TOS = TOS1 | TOS.
func do_BINARY_OR(vm *Vm, arg int32) {
}

// In-place operations are like binary operations, in that they remove
// TOS and TOS1, and push the result back on the stack, but the
// operation is done in-place when TOS1 supports it, and the resulting
// TOS may be (but does not have to be) the original TOS1.

// Implements in-place TOS = TOS1 ** TOS.
func do_INPLACE_POWER(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 * TOS.
func do_INPLACE_MULTIPLY(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 // TOS.
func do_INPLACE_FLOOR_DIVIDE(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 / TOS when from __future__ import
// division is in effect.
func do_INPLACE_TRUE_DIVIDE(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 % TOS.
func do_INPLACE_MODULO(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 + TOS.
func do_INPLACE_ADD(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 - TOS.
func do_INPLACE_SUBTRACT(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 << TOS.
func do_INPLACE_LSHIFT(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 >> TOS.
func do_INPLACE_RSHIFT(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 & TOS.
func do_INPLACE_AND(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 ^ TOS.
func do_INPLACE_XOR(vm *Vm, arg int32) {
}

// Implements in-place TOS = TOS1 | TOS.
func do_INPLACE_OR(vm *Vm, arg int32) {
}

// Implements TOS1[TOS] = TOS2.
func do_STORE_SUBSCR(vm *Vm, arg int32) {
}

// Implements del TOS1[TOS].
func do_DELETE_SUBSCR(vm *Vm, arg int32) {
}

// Miscellaneous opcodes.

// Implements the expression statement for the interactive mode. TOS
// is removed from the stack and printed. In non-interactive mode, an
// expression statement is terminated with POP_STACK.
func do_PRINT_EXPR(vm *Vm, arg int32) {
}

// Terminates a loop due to a break statement.
func do_BREAK_LOOP(vm *Vm, arg int32) {
}

// Continues a loop due to a continue statement. target is the address
// to jump to (which should be a FOR_ITER instruction).
func do_CONTINUE_LOOP(vm *Vm, target int32) {
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
}

// Calls set.add(TOS1[-i], TOS). Used to implement set comprehensions.
func do_SET_ADD(vm *Vm, i int32) {
}

// Calls list.append(TOS[-i], TOS). Used to implement list
// comprehensions. While the appended value is popped off, the list
// object remains on the stack so that it is available for further
// iterations of the loop.
func do_LIST_APPEND(vm *Vm, i int32) {
}

// Calls dict.setitem(TOS1[-i], TOS, TOS1). Used to implement dict comprehensions.
func do_MAP_ADD(vm *Vm, i int32) {
}

// Returns with TOS to the caller of the function.
func do_RETURN_VALUE(vm *Vm, arg int32) {
	vm.exit = true
}

// Pops TOS and delegates to it as a subiterator from a generator.
func do_YIELD_FROM(vm *Vm, arg int32) {
}

// Pops TOS and yields it from a generator.
func do_YIELD_VALUE(vm *Vm, arg int32) {
}

// Loads all symbols not starting with '_' directly from the module
// TOS to the local namespace. The module is popped after loading all
// names. This opcode implements from module import *.
func do_IMPORT_STAR(vm *Vm, arg int32) {
}

// Removes one block from the block stack. Per frame, there is a stack
// of blocks, denoting nested loops, try statements, and such.
func do_POP_BLOCK(vm *Vm, arg int32) {
}

// Removes one block from the block stack. The popped block must be an
// exception handler block, as implicitly created when entering an
// except handler. In addition to popping extraneous values from the
// frame stack, the last three popped values are used to restore the
// exception state.
func do_POP_EXCEPT(vm *Vm, arg int32) {
}

// Terminates a finally clause. The interpreter recalls whether the
// exception has to be re-raised, or whether the function returns, and
// continues with the outer-next block.
func do_END_FINALLY(vm *Vm, arg int32) {
}

// Creates a new class object. TOS is the methods dictionary, TOS1 the
// tuple of the names of the base classes, and TOS2 the class name.
func do_LOAD_BUILD_CLASS(vm *Vm, arg int32) {
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
}

// All of the following opcodes expect arguments. An argument is two bytes, with the more significant byte last.

// Implements name = TOS. namei is the index of name in the attribute
// co_names of the code object. The compiler tries to use STORE_FAST
// or STORE_GLOBAL if possible.
func do_STORE_NAME(vm *Vm, namei int32) {
}

// Implements del name, where namei is the index into co_names
// attribute of the code object.
func do_DELETE_NAME(vm *Vm, namei int32) {
}

// Unpacks TOS into count individual values, which are put onto the
// stack right-to-left.
func do_UNPACK_SEQUENCE(vm *Vm, count int32) {
}

// Implements TOS.name = TOS1, where namei is the index of name in
// co_names.
func do_STORE_ATTR(vm *Vm, namei int32) {
}

// Implements del TOS.name, using namei as index into co_names.
func do_DELETE_ATTR(vm *Vm, namei int32) {
}

// Works as STORE_NAME, but stores the name as a global.
func do_STORE_GLOBAL(vm *Vm, namei int32) {
}

// Works as DELETE_NAME, but deletes a global name.
func do_DELETE_GLOBAL(vm *Vm, namei int32) {
}

// Pushes co_consts[consti] onto the stack.
func do_LOAD_CONST(vm *Vm, consti int32) {
	vm.PUSH(vm.co.Consts[consti])
}

// Pushes the value associated with co_names[namei] onto the stack.
func do_LOAD_NAME(vm *Vm, namei int32) {
	vm.PUSH(vm.co.Names[namei])
}

// Creates a tuple consuming count items from the stack, and pushes
// the resulting tuple onto the stack.
func do_BUILD_TUPLE(vm *Vm, count int32) {
}

// Works as BUILD_TUPLE, but creates a set.
func do_BUILD_SET(vm *Vm, count int32) {
}

// Works as BUILD_TUPLE, but creates a list.
func do_BUILD_LIST(vm *Vm, count int32) {
}

// Pushes a new dictionary object onto the stack. The dictionary is
// pre-sized to hold count entries.
func do_BUILD_MAP(vm *Vm, count int32) {
}

// Replaces TOS with getattr(TOS, co_names[namei]).
func do_LOAD_ATTR(vm *Vm, namei int32) {
}

// Performs a Boolean operation. The operation name can be found in
// cmp_op[opname].
func do_COMPARE_OP(vm *Vm, opname int32) {
}

// Imports the module co_names[namei]. TOS and TOS1 are popped and
// provide the fromlist and level arguments of __import__( ). The
// module object is pushed onto the stack. The current namespace is
// not affected: for a proper import statement, a subsequent
// STORE_FAST instruction modifies the namespace.
func do_IMPORT_NAME(vm *Vm, namei int32) {
}

// Loads the attribute co_names[namei] from the module found in
// TOS. The resulting object is pushed onto the stack, to be
// subsequently stored by a STORE_FAST instruction.
func do_IMPORT_FROM(vm *Vm, namei int32) {
}

// Increments bytecode counter by delta.
func do_JUMP_FORWARD(vm *Vm, delta int32) {
}

// If TOS is true, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_TRUE(vm *Vm, target int32) {
}

// If TOS is false, sets the bytecode counter to target. TOS is popped.
func do_POP_JUMP_IF_FALSE(vm *Vm, target int32) {
}

// If TOS is true, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is false), TOS is popped.
func do_JUMP_IF_TRUE_OR_POP(vm *Vm, target int32) {
}

// If TOS is false, sets the bytecode counter to target and leaves TOS
// on the stack. Otherwise (TOS is true), TOS is popped.
func do_JUMP_IF_FALSE_OR_POP(vm *Vm, target int32) {
}

// Set bytecode counter to target.
func do_JUMP_ABSOLUTE(vm *Vm, target int32) {
}

// TOS is an iterator. Call its next( ) method. If this yields a new
// value, push it on the stack (leaving the iterator below it). If the
// iterator indicates it is exhausted TOS is popped, and the bytecode
// counter is incremented by delta.
func do_FOR_ITER(vm *Vm, delta int32) {
}

// Loads the global named co_names[namei] onto the stack.
func do_LOAD_GLOBAL(vm *Vm, namei int32) {
}

// Pushes a block for a loop onto the block stack. The block spans
// from the current instruction with a size of delta bytes.
func do_SETUP_LOOP(vm *Vm, delta int32) {
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the first except block.
func do_SETUP_EXCEPT(vm *Vm, delta int32) {
}

// Pushes a try block from a try-except clause onto the block
// stack. delta points to the finally block.
func do_SETUP_FINALLY(vm *Vm, delta int32) {
}

// Store a key and value pair in a dictionary. Pops the key and value
// while leaving the dictionary on the stack.
func do_STORE_MAP(vm *Vm, arg int32) {
}

// Pushes a reference to the local co_varnames[var_num] onto the stack.
func do_LOAD_FAST(vm *Vm, var_num int32) {
}

// Stores TOS into the local co_varnames[var_num].
func do_STORE_FAST(vm *Vm, var_num int32) {
}

// Deletes local co_varnames[var_num].
func do_DELETE_FAST(vm *Vm, var_num int32) {
}

// Pushes a reference to the cell contained in slot i of the cell and
// free variable storage. The name of the variable is co_cellvars[i]
// if i is less than the length of co_cellvars. Otherwise it is
// co_freevars[i - len(co_cellvars)].
func do_LOAD_CLOSURE(vm *Vm, i int32) {
}

// Loads the cell contained in slot i of the cell and free variable
// storage. Pushes a reference to the object the cell contains on the
// stack.
func do_LOAD_DEREF(vm *Vm, i int32) {
}

// Much like LOAD_DEREF but first checks the locals dictionary before
// consulting the cell. This is used for loading free variables in
// class bodies.
func do_LOAD_CLASSDEREF(vm *Vm, i int32) {
}

// Stores TOS into the cell contained in slot i of the cell and free
// variable storage.
func do_STORE_DEREF(vm *Vm, i int32) {
}

// Empties the cell contained in slot i of the cell and free variable
// storage. Used by the del statement.
func do_DELETE_DEREF(vm *Vm, i int32) {
}

// Raises an exception. argc indicates the number of parameters to the
// raise statement, ranging from 0 to 3. The handler will find the
// traceback as TOS2, the parameter as TOS1, and the exception as TOS.
func do_RAISE_VARARGS(vm *Vm, argc int32) {
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
	nargs := int(argc & 0xFF)
	nkwargs := int((argc >> 8) & 0xFF)
	p, q := len(vm.stack)-2*nkwargs, len(vm.stack)
	kwargs := vm.stack[p:q]
	p, q = p-nargs, p
	args := py.Tuple(vm.stack[p:q])
	p, q = p-1, p
	fn := vm.stack[p]
	fmt.Printf("Call %v with args = %v, kwargs = %v\n", fn, args, kwargs)
	// FIXME look the function up
	fn_name := string(fn.(py.String))
	if method, ok := py.Builtins.Methods[fn_name]; ok {
		// FIXME need module as self
		self := py.None
		if len(kwargs) > 0 {
			// FIXME need to convert kwargs to dictionary
			kwargsd := py.NewDict()
			vm.stack[p] = method.CallWithKeywords(self, args, kwargsd)
		} else {
			vm.stack[p] = method.Call(self, args)
		}
	} else {
		panic("Couldn't find method")
	}
	// Drop the args off the stack and put return value in
	vm.stack = vm.stack[:q]
	vm.stack[p] = py.None
}

// Pushes a new function object on the stack. TOS is the code
// associated with the function. The function object is defined to
// have argc default parameters, which are found below TOS.
func do_MAKE_FUNCTION(vm *Vm, argc int32) {
}

// Creates a new function object, sets its func_closure slot, and
// pushes it on the stack. TOS is the code associated with the
// function, TOS1 the tuple containing cells for the closure’s free
// variables. The function also has argc default parameters, which are
// found below the cells.
func do_MAKE_CLOSURE(vm *Vm, argc int32) {
}

// Pushes a slice object on the stack. argc must be 2 or 3. If it is
// 2, slice(TOS1, TOS) is pushed; if it is 3, slice(TOS2, TOS1, TOS)
// is pushed. See the slice( ) built-in function for more information.
func do_BUILD_SLICE(vm *Vm, argc int32) {
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
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by explicit keyword and positional arguments.
func do_CALL_FUNCTION_KW(vm *Vm, argc int32) {
}

// Calls a function. argc is interpreted as in CALL_FUNCTION. The top
// element on the stack contains the keyword arguments dictionary,
// followed by the variable-arguments tuple, followed by explicit
// keyword and positional arguments.
func do_CALL_FUNCTION_VAR_KW(vm *Vm, argc int32) {
}

// Run the virtual machine on the code object
//
// FIXME figure out how we are going to signal exceptions!
func (vm *Vm) Run(co *py.Code) (err error) {
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
	vm.co = co
	ip := 0
	var opcode byte
	var arg int32
	code := co.Code
	for !vm.exit {
		opcode = code[ip]
		ip++
		if HAS_ARG(opcode) {
			arg = int32(code[ip])
			ip++
			arg += int32(code[ip] << 8)
			ip++
			if vm.extended {
				arg += vm.ext << 16
			}
			fmt.Printf("Opcode %d with arg %d\n", opcode, arg)
		} else {
			fmt.Printf("Opcode %d\n", opcode)
		}
		vm.extended = false
		jumpTable[opcode](vm, arg)
	}
	return nil
}
