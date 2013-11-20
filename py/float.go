// Float objects

package py

import (
	"math"
)

var FloatType = NewType("float")

type Float float64

// Type of this Float64 object
func (o Float) Type() *Type {
	return FloatType
}

// Arithmetic

// Convert an Object to an Float
//
// Retrurns ok as to whether the conversion worked or not
func convertToFloat(other Object) (Float, bool) {
	switch b := other.(type) {
	case Float:
		return b, true
	case Int:
		return Float(b), true
	case Bool:
		if b {
			return Float(1), true
		} else {
			return Float(0), true
		}
	}
	return 0, false
}

func (a Float) M__neg__() Object {
	return -a
}

func (a Float) M__pos__() Object {
	return a
}

func (a Float) M__abs__() Object {
	return Float(math.Abs(float64(a)))
}

func (a Float) M__add__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(a + b)
	}
	return NotImplemented
}

func (a Float) M__radd__(other Object) Object {
	return a.M__add__(other)
}

func (a Float) M__iadd__(other Object) Object {
	return a.M__add__(other)
}

func (a Float) M__sub__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(a - b)
	}
	return NotImplemented
}

func (a Float) M__rsub__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(b - a)
	}
	return NotImplemented
}

func (a Float) M__isub__(other Object) Object {
	return a.M__sub__(other)
}

func (a Float) M__mul__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(a * b)
	}
	return NotImplemented
}

func (a Float) M__rmul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Float) M__imul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Float) M__truediv__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(a / b)
	}
	return NotImplemented
}

func (a Float) M__rtruediv__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(b / a)
	}
	return NotImplemented
}

func (a Float) M__itruediv__(other Object) Object {
	return Float(a).M__truediv__(other)
}

func (a Float) M__floordiv__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Floor(float64(a / b)))
	}
	return NotImplemented
}

func (a Float) M__rfloordiv__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Floor(float64(b / a)))
	}
	return NotImplemented
}

func (a Float) M__ifloordiv__(other Object) Object {
	return a.M__floordiv__(other)
}

// Does DivMod of two floating point numbers
func floatDivMod(a, b Float) (Float, Float) {
	q := Float(math.Floor(float64(a / b)))
	r := a - q*b
	return q, Float(r)
}

func (a Float) M__mod__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		_, r := floatDivMod(a, b)
		return r
	}
	return NotImplemented
}

func (a Float) M__rmod__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		_, r := floatDivMod(b, a)
		return r
	}
	return NotImplemented
}

func (a Float) M__imod__(other Object) Object {
	return a.M__mod__(other)
}

func (a Float) M__divmod__(other Object) (Object, Object) {
	if b, ok := convertToFloat(other); ok {
		return floatDivMod(a, b)
	}
	return NotImplemented, None
}

func (a Float) M__rdivmod__(other Object) (Object, Object) {
	if b, ok := convertToFloat(other); ok {
		return floatDivMod(b, a)
	}
	return NotImplemented, None
}

func (a Float) M__pow__(other, modulus Object) Object {
	if modulus != None {
		return NotImplemented
	}
	if b, ok := convertToFloat(other); ok {
		return Float(math.Pow(float64(a), float64(b)))
	}
	return NotImplemented
}

func (a Float) M__rpow__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Pow(float64(b), float64(a)))
	}
	return NotImplemented
}

func (a Float) M__ipow__(other, modulus Object) Object {
	return a.M__pow__(other, modulus)
}

func (a Float) M__bool__() Object {
	if a == 0 {
		return False
	}
	return True
}

func (a Float) M__int__() Object {
	return Int(a)
}

func (a Float) M__float__() Object {
	return a
}

func (a Float) M__complex__() Object {
	if r, ok := convertToComplex(a); ok {
		return r
	}
	panic("convertToComplex failed")
}

func (a Float) M__round__(digitsObj Object) Object {
	digits := 0
	if digitsObj != None {
		digits = Index(digitsObj)
	}
	scale := Float(math.Pow(10, float64(digits)))
	return scale * Float(math.Floor(float64(a)/float64(scale)))
}

// Check interface is satisfied
var _ floatArithmetic = Float(0)
var _ conversionBetweenTypes = Int(0)
var _ I__bool__ = Int(0)
