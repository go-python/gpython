// Int and BigInt objects

package py

import (
	"math"
	"math/big"
	"strconv"
)

var IntType = ObjectType.NewType("int", `
int(x=0) -> integer
int(x, base=10) -> integer

Convert a number or string to an integer, or return 0 if no arguments
are given.  If x is a number, return x.__int__().  For floating point
numbers, this truncates towards zero.

If x is not a number or if base is given, then x must be a string,
bytes, or bytearray instance representing an integer literal in the
given base.  The literal can be preceded by '+' or '-' and be surrounded
by whitespace.  The base defaults to 10.  Valid bases are 0 and 2-36.
Base 0 means to interpret the base from the string as an integer literal.
>>> int('0b100', base=0)
4`, IntNew, nil)

type Int int64

const (
	// Maximum possible Int
	IntMax = 1<<63 - 1
	// Minimum possible Int
	IntMin = -IntMax - 1
)

// Type of this Int object
func (o Int) Type() *Type {
	return IntType
}

type BigInt big.Int

var BigIntType = NewType("bigint", "Holds large integers")

// Type of this BigInt object
func (o *BigInt) Type() *Type {
	return BigIntType
}

// Make sure it satisfies the interface
var _ Object = (*BigInt)(nil)

// IntNew
func IntNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	var xObj Object = Int(0)
	var baseObj Object = Int(10)
	ParseTupleAndKeywords(args, kwargs, "|OO:int", []string{"x", "base"}, &xObj, &baseObj)
	var res Int
	switch x := xObj.(type) {
	case Int:
		res = x
	case String:
		base := int(baseObj.(Int))
		// FIXME this isn't 100% python compatible but it is close!
		i, err := strconv.ParseInt(string(x), base, 64)
		if err != nil {
			panic(ExceptionNewf(ValueError, "invalid literal for int() with base %d: '%s'", base, string(x)))
		}
		res = Int(i)
	case Float:
		res = Int(x)
	default:
		var ok bool
		res, ok = convertToInt(x)
		if !ok {
			panic(ExceptionNewf(TypeError, "int() argument must be a string or a number, not 'tuple'"))
		}
	}
	return res
}

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

func (a Int) M__neg__() Object {
	return -a
}

func (a Int) M__pos__() Object {
	return a
}

func (a Int) M__abs__() Object {
	if a < 0 {
		return -a
	}
	return a
}

func (a Int) M__invert__() Object {
	return ^a
}

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

// FIXME implement powmod...
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

func (a Int) M__bool__() Object {
	return NewBool(a != 0)
}

func (a Int) M__index__() Int {
	return a
}

func (a Int) M__int__() Object {
	return a
}

func (a Int) M__float__() Object {
	if r, ok := convertToFloat(a); ok {
		return r
	}
	panic("convertToFloat failed")
}

func (a Int) M__complex__() Object {
	if r, ok := convertToComplex(a); ok {
		return r
	}
	panic("convertToComplex failed")
}

func (a Int) M__round__(digits Object) Object {
	return Int(Float(a).M__round__(digits).(Float))
}

// Rich comparison

func (a Int) M__lt__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a < b)
	}
	return NotImplemented
}

func (a Int) M__le__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a <= b)
	}
	return NotImplemented
}

func (a Int) M__eq__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a == b)
	}
	return NotImplemented
}

func (a Int) M__ne__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a != b)
	}
	return NotImplemented
}

func (a Int) M__gt__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a > b)
	}
	return NotImplemented
}

func (a Int) M__ge__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		return NewBool(a >= b)
	}
	return NotImplemented
}

// Check interface is satisfied
var _ floatArithmetic = Int(0)
var _ booleanArithmetic = Int(0)
var _ conversionBetweenTypes = Int(0)
var _ I__bool__ = Int(0)
var _ I__index__ = Int(0)
var _ richComparison = Int(0)
