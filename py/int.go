// Int and BigInt objects

package py

import (
	"math"
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

func (a Int) M__isub__(other Object) Object {
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

func (a Int) M__itruediv__(other Object) Object {
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

func (a Int) M__ifloordiv__(other Object) Object {
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

func (a Int) M__imod__(other Object) Object {
	return a.M__mod__(other)
}

func (a Int) M__divmod__(other Object) (Object, Object) {
	if b, ok := convertToInt(other); ok {
		return Int(a / b), Int(a % b)
	}
	return NotImplemented, None
}

func (a Int) M__rdivmod__(other Object) (Object, Object) {
	if b, ok := convertToInt(other); ok {
		return Int(b / a), Int(b % a)
	}
	return NotImplemented, None
}

func (a Int) M__pow__(other, modulus Object) Object {
	if modulus != None {
		return NotImplemented
	}
	if b, ok := convertToInt(other); ok {
		// FIXME possible loss of precision
		return Int(math.Pow(float64(a), float64(b)))
	}
	return NotImplemented
}

func (a Int) M__rpow__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		// FIXME possible loss of precision
		return Int(math.Pow(float64(b), float64(a)))
	}
	return NotImplemented
}

func (a Int) M__ipow__(other, modulus Object) Object {
	return a.M__pow__(other, modulus)
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

func (a Int) M__ilshift__(other Object) Object {
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

func (a Int) M__irshift__(other Object) Object {
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

func (a Int) M__iand__(other Object) Object {
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

func (a Int) M__ixor__(other Object) Object {
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

func (a Int) M__ior__(other Object) Object {
	return a.M__or__(other)
}

// Check interface is satisfied
var _ floatArithmetic = Int(0)
var _ booleanArithmetic = Int(0)
