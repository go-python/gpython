// Arithmetic operations

package py

import (
	"fmt"
)

// Add two python objects together returning an Object
//
// Will raise TypeError if can't be added
func Add(a, b Object) Object {
	// Try using a to add
	A, ok := a.(I__add__)
	if ok {
		res := A.M__add__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to radd if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__radd__)
		if ok {
			res := B.M__radd__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for +: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace add
func Iadd(a, b Object) Object {
	A, ok := a.(I__iadd__)
	if ok {
		res := A.M__iadd__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Add(a, b)
}

// Subtract two python objects together returning an Object
//
// Will raise TypeError if can't be subtracted
func Sub(a, b Object) Object {
	// Try using a to sub
	A, ok := a.(I__sub__)
	if ok {
		res := A.M__sub__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rsub if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rsub__)
		if ok {
			res := B.M__rsub__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for -: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace sub
func Isub(a, b Object) Object {
	A, ok := a.(I__isub__)
	if ok {
		res := A.M__isub__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Sub(a, b)
}
