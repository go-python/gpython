// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Build the jump table

package vm

// Globals
var (
	jumpTable [256]func(*Vm, int32) error
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
