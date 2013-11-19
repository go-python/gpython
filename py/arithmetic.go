// Automatically generated - DO NOT EDIT
// Regenerate with: go run gen.go | gofmt >arithmetic.go

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
func IAdd(a, b Object) Object {
	A, ok := a.(I__iadd__)
	if ok {
		res := A.M__iadd__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Add(a, b)
}

// Sub two python objects together returning an Object
//
// Will raise TypeError if can't be subed
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
func ISub(a, b Object) Object {
	A, ok := a.(I__isub__)
	if ok {
		res := A.M__isub__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Sub(a, b)
}

// Mul two python objects together returning an Object
//
// Will raise TypeError if can't be muled
func Mul(a, b Object) Object {
	// Try using a to mul
	A, ok := a.(I__mul__)
	if ok {
		res := A.M__mul__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rmul if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rmul__)
		if ok {
			res := B.M__rmul__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for *: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace mul
func IMul(a, b Object) Object {
	A, ok := a.(I__imul__)
	if ok {
		res := A.M__imul__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Mul(a, b)
}

// TrueDiv two python objects together returning an Object
//
// Will raise TypeError if can't be truedived
func TrueDiv(a, b Object) Object {
	// Try using a to truediv
	A, ok := a.(I__truediv__)
	if ok {
		res := A.M__truediv__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rtruediv if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rtruediv__)
		if ok {
			res := B.M__rtruediv__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for /: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace truediv
func ITrueDiv(a, b Object) Object {
	A, ok := a.(I__itruediv__)
	if ok {
		res := A.M__itruediv__(b)
		if res != NotImplemented {
			return res
		}
	}
	return TrueDiv(a, b)
}

// FloorDiv two python objects together returning an Object
//
// Will raise TypeError if can't be floordived
func FloorDiv(a, b Object) Object {
	// Try using a to floordiv
	A, ok := a.(I__floordiv__)
	if ok {
		res := A.M__floordiv__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rfloordiv if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rfloordiv__)
		if ok {
			res := B.M__rfloordiv__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for //: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace floordiv
func IFloorDiv(a, b Object) Object {
	A, ok := a.(I__ifloordiv__)
	if ok {
		res := A.M__ifloordiv__(b)
		if res != NotImplemented {
			return res
		}
	}
	return FloorDiv(a, b)
}

// Mod two python objects together returning an Object
//
// Will raise TypeError if can't be moded
func Mod(a, b Object) Object {
	// Try using a to mod
	A, ok := a.(I__mod__)
	if ok {
		res := A.M__mod__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rmod if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rmod__)
		if ok {
			res := B.M__rmod__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for %: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace mod
func IMod(a, b Object) Object {
	A, ok := a.(I__imod__)
	if ok {
		res := A.M__imod__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Mod(a, b)
}

// Lshift two python objects together returning an Object
//
// Will raise TypeError if can't be lshifted
func Lshift(a, b Object) Object {
	// Try using a to lshift
	A, ok := a.(I__lshift__)
	if ok {
		res := A.M__lshift__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rlshift if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rlshift__)
		if ok {
			res := B.M__rlshift__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for <<: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace lshift
func ILshift(a, b Object) Object {
	A, ok := a.(I__ilshift__)
	if ok {
		res := A.M__ilshift__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Lshift(a, b)
}

// Rshift two python objects together returning an Object
//
// Will raise TypeError if can't be rshifted
func Rshift(a, b Object) Object {
	// Try using a to rshift
	A, ok := a.(I__rshift__)
	if ok {
		res := A.M__rshift__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rrshift if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rrshift__)
		if ok {
			res := B.M__rrshift__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for >>: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace rshift
func IRshift(a, b Object) Object {
	A, ok := a.(I__irshift__)
	if ok {
		res := A.M__irshift__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Rshift(a, b)
}

// And two python objects together returning an Object
//
// Will raise TypeError if can't be anded
func And(a, b Object) Object {
	// Try using a to and
	A, ok := a.(I__and__)
	if ok {
		res := A.M__and__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rand if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rand__)
		if ok {
			res := B.M__rand__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for &: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace and
func IAnd(a, b Object) Object {
	A, ok := a.(I__iand__)
	if ok {
		res := A.M__iand__(b)
		if res != NotImplemented {
			return res
		}
	}
	return And(a, b)
}

// Xor two python objects together returning an Object
//
// Will raise TypeError if can't be xored
func Xor(a, b Object) Object {
	// Try using a to xor
	A, ok := a.(I__xor__)
	if ok {
		res := A.M__xor__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rxor if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__rxor__)
		if ok {
			res := B.M__rxor__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for ^: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace xor
func IXor(a, b Object) Object {
	A, ok := a.(I__ixor__)
	if ok {
		res := A.M__ixor__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Xor(a, b)
}

// Or two python objects together returning an Object
//
// Will raise TypeError if can't be ored
func Or(a, b Object) Object {
	// Try using a to or
	A, ok := a.(I__or__)
	if ok {
		res := A.M__or__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to ror if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__ror__)
		if ok {
			res := B.M__ror__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for |: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace or
func IOr(a, b Object) Object {
	A, ok := a.(I__ior__)
	if ok {
		res := A.M__ior__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Or(a, b)
}
