// Automatically generated - DO NOT EDIT
// Regenerate with: go run gen.go | gofmt >arithmetic.go

// Arithmetic operations

package py

import (
	"fmt"
)

// Neg the python Object returning an Object
//
// Will raise TypeError if Neg can't be run on this object
func Neg(a Object) Object {

	if A, ok := a.(I__neg__); ok {
		res := A.M__neg__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for -: '%s'", a.Type().Name))
}

// Pos the python Object returning an Object
//
// Will raise TypeError if Pos can't be run on this object
func Pos(a Object) Object {

	if A, ok := a.(I__pos__); ok {
		res := A.M__pos__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for +: '%s'", a.Type().Name))
}

// Abs the python Object returning an Object
//
// Will raise TypeError if Abs can't be run on this object
func Abs(a Object) Object {

	if A, ok := a.(I__abs__); ok {
		res := A.M__abs__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for abs: '%s'", a.Type().Name))
}

// Invert the python Object returning an Object
//
// Will raise TypeError if Invert can't be run on this object
func Invert(a Object) Object {

	if A, ok := a.(I__invert__); ok {
		res := A.M__invert__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for ~: '%s'", a.Type().Name))
}

// MakeComplex the python Object returning an Object
//
// Will raise TypeError if MakeComplex can't be run on this object
func MakeComplex(a Object) Object {

	if _, ok := a.(Complex); ok {
		return a
	}

	if A, ok := a.(I__complex__); ok {
		res := A.M__complex__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for complex: '%s'", a.Type().Name))
}

// MakeInt the python Object returning an Object
//
// Will raise TypeError if MakeInt can't be run on this object
func MakeInt(a Object) Object {

	if _, ok := a.(Int); ok {
		return a
	}

	if A, ok := a.(I__int__); ok {
		res := A.M__int__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for int: '%s'", a.Type().Name))
}

// MakeFloat the python Object returning an Object
//
// Will raise TypeError if MakeFloat can't be run on this object
func MakeFloat(a Object) Object {

	if _, ok := a.(Float); ok {
		return a
	}

	if A, ok := a.(I__float__); ok {
		res := A.M__float__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for float: '%s'", a.Type().Name))
}

// Add two python objects together returning an Object
//
// Will raise TypeError if can't be add can't be run on these objects
func Add(a, b Object) Object {
	// Try using a to add
	if A, ok := a.(I__add__); ok {
		res := A.M__add__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to radd if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__radd__); ok {
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
	if A, ok := a.(I__iadd__); ok {
		res := A.M__iadd__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Add(a, b)
}

// Sub two python objects together returning an Object
//
// Will raise TypeError if can't be sub can't be run on these objects
func Sub(a, b Object) Object {
	// Try using a to sub
	if A, ok := a.(I__sub__); ok {
		res := A.M__sub__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rsub if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rsub__); ok {
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
	if A, ok := a.(I__isub__); ok {
		res := A.M__isub__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Sub(a, b)
}

// Mul two python objects together returning an Object
//
// Will raise TypeError if can't be mul can't be run on these objects
func Mul(a, b Object) Object {
	// Try using a to mul
	if A, ok := a.(I__mul__); ok {
		res := A.M__mul__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rmul if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rmul__); ok {
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
	if A, ok := a.(I__imul__); ok {
		res := A.M__imul__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Mul(a, b)
}

// TrueDiv two python objects together returning an Object
//
// Will raise TypeError if can't be truediv can't be run on these objects
func TrueDiv(a, b Object) Object {
	// Try using a to truediv
	if A, ok := a.(I__truediv__); ok {
		res := A.M__truediv__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rtruediv if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rtruediv__); ok {
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
	if A, ok := a.(I__itruediv__); ok {
		res := A.M__itruediv__(b)
		if res != NotImplemented {
			return res
		}
	}
	return TrueDiv(a, b)
}

// FloorDiv two python objects together returning an Object
//
// Will raise TypeError if can't be floordiv can't be run on these objects
func FloorDiv(a, b Object) Object {
	// Try using a to floordiv
	if A, ok := a.(I__floordiv__); ok {
		res := A.M__floordiv__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rfloordiv if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rfloordiv__); ok {
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
	if A, ok := a.(I__ifloordiv__); ok {
		res := A.M__ifloordiv__(b)
		if res != NotImplemented {
			return res
		}
	}
	return FloorDiv(a, b)
}

// Mod two python objects together returning an Object
//
// Will raise TypeError if can't be mod can't be run on these objects
func Mod(a, b Object) Object {
	// Try using a to mod
	if A, ok := a.(I__mod__); ok {
		res := A.M__mod__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rmod if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rmod__); ok {
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
	if A, ok := a.(I__imod__); ok {
		res := A.M__imod__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Mod(a, b)
}

// DivMod two python objects together returning an Object
//
// Will raise TypeError if can't be divmod can't be run on these objects
func DivMod(a, b Object) (Object, Object) {
	// Try using a to divmod
	if A, ok := a.(I__divmod__); ok {
		res, res2 := A.M__divmod__(b)
		if res != NotImplemented {
			return res, res2
		}
	}

	// Now using b to rdivmod if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rdivmod__); ok {
			res, res2 := B.M__rdivmod__(a)
			if res != NotImplemented {
				return res, res2
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for divmod: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Lshift two python objects together returning an Object
//
// Will raise TypeError if can't be lshift can't be run on these objects
func Lshift(a, b Object) Object {
	// Try using a to lshift
	if A, ok := a.(I__lshift__); ok {
		res := A.M__lshift__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rlshift if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rlshift__); ok {
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
	if A, ok := a.(I__ilshift__); ok {
		res := A.M__ilshift__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Lshift(a, b)
}

// Rshift two python objects together returning an Object
//
// Will raise TypeError if can't be rshift can't be run on these objects
func Rshift(a, b Object) Object {
	// Try using a to rshift
	if A, ok := a.(I__rshift__); ok {
		res := A.M__rshift__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rrshift if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rrshift__); ok {
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
	if A, ok := a.(I__irshift__); ok {
		res := A.M__irshift__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Rshift(a, b)
}

// And two python objects together returning an Object
//
// Will raise TypeError if can't be and can't be run on these objects
func And(a, b Object) Object {
	// Try using a to and
	if A, ok := a.(I__and__); ok {
		res := A.M__and__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rand if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rand__); ok {
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
	if A, ok := a.(I__iand__); ok {
		res := A.M__iand__(b)
		if res != NotImplemented {
			return res
		}
	}
	return And(a, b)
}

// Xor two python objects together returning an Object
//
// Will raise TypeError if can't be xor can't be run on these objects
func Xor(a, b Object) Object {
	// Try using a to xor
	if A, ok := a.(I__xor__); ok {
		res := A.M__xor__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rxor if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__rxor__); ok {
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
	if A, ok := a.(I__ixor__); ok {
		res := A.M__ixor__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Xor(a, b)
}

// Or two python objects together returning an Object
//
// Will raise TypeError if can't be or can't be run on these objects
func Or(a, b Object) Object {
	// Try using a to or
	if A, ok := a.(I__or__); ok {
		res := A.M__or__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to ror if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__ror__); ok {
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
	if A, ok := a.(I__ior__); ok {
		res := A.M__ior__(b)
		if res != NotImplemented {
			return res
		}
	}
	return Or(a, b)
}

// Pow three python objects together returning an Object
//
// If c != None then it won't attempt to call __rpow__
//
// Will raise TypeError if can't be pow can't be run on these objects
func Pow(a, b, c Object) Object {
	// Try using a to pow
	if A, ok := a.(I__pow__); ok {
		res := A.M__pow__(b, c)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to rpow if different in type to a
	if c == None && a.Type() != b.Type() {
		if B, ok := b.(I__rpow__); ok {
			res := B.M__rpow__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for ** or pow(): '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace pow
func IPow(a, b, c Object) Object {
	if A, ok := a.(I__ipow__); ok {
		res := A.M__ipow__(b, c)
		if res != NotImplemented {
			return res
		}
	}
	return Pow(a, b, c)
}

// Gt two python objects returning a boolean result
//
// Will raise TypeError if Gt can't be run on this object
func Gt(a Object, b Object) Object {
	// Try using a to gt
	if A, ok := a.(I__gt__); ok {
		res := A.M__gt__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to le with reversed parameters
	if B, ok := a.(I__le__); ok {
		res := B.M__le__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for >: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Ge two python objects returning a boolean result
//
// Will raise TypeError if Ge can't be run on this object
func Ge(a Object, b Object) Object {
	// Try using a to ge
	if A, ok := a.(I__ge__); ok {
		res := A.M__ge__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to lt with reversed parameters
	if B, ok := a.(I__lt__); ok {
		res := B.M__lt__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for >=: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Lt two python objects returning a boolean result
//
// Will raise TypeError if Lt can't be run on this object
func Lt(a Object, b Object) Object {
	// Try using a to lt
	if A, ok := a.(I__lt__); ok {
		res := A.M__lt__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to ge with reversed parameters
	if B, ok := a.(I__ge__); ok {
		res := B.M__ge__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for <: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Le two python objects returning a boolean result
//
// Will raise TypeError if Le can't be run on this object
func Le(a Object, b Object) Object {
	// Try using a to le
	if A, ok := a.(I__le__); ok {
		res := A.M__le__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to gt with reversed parameters
	if B, ok := a.(I__gt__); ok {
		res := B.M__gt__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for <=: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Eq two python objects returning a boolean result
//
// Will raise TypeError if Eq can't be run on this object
func Eq(a Object, b Object) Object {
	// Try using a to eq
	if A, ok := a.(I__eq__); ok {
		res := A.M__eq__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to ne with reversed parameters
	if B, ok := a.(I__ne__); ok {
		res := B.M__ne__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for ==: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Ne two python objects returning a boolean result
//
// Will raise TypeError if Ne can't be run on this object
func Ne(a Object, b Object) Object {
	// Try using a to ne
	if A, ok := a.(I__ne__); ok {
		res := A.M__ne__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to eq with reversed parameters
	if B, ok := a.(I__eq__); ok {
		res := B.M__eq__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for !=: '%s' and '%s'", a.Type().Name, b.Type().Name))
}
