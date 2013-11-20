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

func (a Float) M__isub(other Object) Object {
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

func (a Float) M__itruediv(other Object) Object {
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

func (a Float) M__ifloordiv(other Object) Object {
	return a.M__floordiv__(other)
}

// Does Mod of two floating point numbers
func floatMod(a, b Float) Float {
	q := Float(math.Floor(float64(a / b)))
	r := a - q*b
	return Float(r)
}

func (a Float) M__mod__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return floatMod(a, b)
	}
	return NotImplemented
}

func (a Float) M__rmod__(other Object) Object {
	if b, ok := convertToFloat(other); ok {
		return floatMod(b, a)
	}
	return NotImplemented
}

func (a Float) M__imod(other Object) Object {
	return a.M__mod__(other)
}
