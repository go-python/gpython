// Int and BigInt objects

package py

import (
	"math/big"
)

var IntType = NewType("int")

type Int int64

// Type of this Int object
func (o Int) Type() *Type {
	return IntType
}

type BigInt big.Int

var BigIntType = NewType("bigint")

// Type of this BigInt object
func (o *BigInt) Type() *Type {
	return BigIntType
}

// Make sure it satisfies the interface
var _ Object = (*BigInt)(nil)

// Arithmetic

// Convert an Object to an Int
//
// Retrurns ok as to whether the conversion worked or not
func convertToInt(other Object) (Int, bool) {
	switch b := other.(type) {
	case Int:
		return b, true
	case Bool:
		if b {
			return Int(1), true
		} else {
			return Int(0), true
		}
	}
	return 0, false
}

// FIXME overflow should promote to Long in all these functions

func (a Int) M__add__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a + b)
	}
	return NotImplemented
}

func (a Int) M__radd__(other Object) Object {
	return a.M__add__(other)
}

func (a Int) M__iadd__(other Object) Object {
	return a.M__add__(other)
}

func (a Int) M__sub__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a - b)
	}
	return NotImplemented
}

func (a Int) M__rsub__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(b - a)
	}
	return NotImplemented
}

func (a Int) M__isub(other Object) Object {
	return a.M__sub__(other)
}

func (a Int) M__mul__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a * b)
	}
	return NotImplemented
}

func (a Int) M__rmul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Int) M__imul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Int) M__truediv__(other Object) Object {
	return Float(a).M__truediv__(other)
}

func (a Int) M__rtruediv__(other Object) Object {
	return Float(a).M__rtruediv__(other)
}

func (a Int) M__itruediv(other Object) Object {
	return Float(a).M__truediv__(other)
}

func (a Int) M__floordiv__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a / b)
	}
	return NotImplemented
}

func (a Int) M__rfloordiv__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(b / a)
	}
	return NotImplemented
}

func (a Int) M__ifloordiv(other Object) Object {
	return a.M__floordiv__(other)
}

func (a Int) M__mod__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a % b)
	}
	return NotImplemented
}

func (a Int) M__rmod__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(b % a)
	}
	return NotImplemented
}

func (a Int) M__imod(other Object) Object {
	return a.M__mod__(other)
}

func (a Int) M__lshift__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			// FIXME should be ValueError
			panic("ValueError: negative shift count")
		}
		return Int(a << uint64(b))
	}
	return NotImplemented
}

func (a Int) M__rlshift__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			// FIXME should be ValueError
			panic("ValueError: negative shift count")
		}
		return Int(b << uint64(a))
	}
	return NotImplemented
}

func (a Int) M__ilshift(other Object) Object {
	return a.M__floordiv__(other)
}

func (a Int) M__rshift__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			// FIXME should be ValueError
			panic("ValueError: negative shift count")
		}
		return Int(a >> uint64(b))
	}
	return NotImplemented
}

func (a Int) M__rrshift__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			// FIXME should be ValueError
			panic("ValueError: negative shift count")
		}
		return Int(b >> uint64(a))
	}
	return NotImplemented
}

func (a Int) M__irshift(other Object) Object {
	return a.M__rshift__(other)
}

func (a Int) M__and__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a & b)
	}
	return NotImplemented
}

func (a Int) M__rand__(other Object) Object {
	return a.M__and__(other)
}

func (a Int) M__iand(other Object) Object {
	return a.M__and__(other)
}

func (a Int) M__xor__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a ^ b)
	}
	return NotImplemented
}

func (a Int) M__rxor__(other Object) Object {
	return a.M__xor__(other)
}

func (a Int) M__ixor(other Object) Object {
	return a.M__xor__(other)
}

func (a Int) M__or__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return Int(a | b)
	}
	return NotImplemented
}

func (a Int) M__ror__(other Object) Object {
	return a.M__or__(other)
}

func (a Int) M__ior(other Object) Object {
	return a.M__or__(other)
}
