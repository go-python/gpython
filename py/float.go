// Float objects

package py

import (
	"math"
)

var FloatType = NewType("float", "float(x) -> floating point number\n\nConvert a string or number to a floating point number, if possible.")

type Float float64

// Type of this Float64 object
func (o Float) Type() *Type {
	return FloatType
}

// Arithmetic

// Errors
var floatDivisionByZero = ExceptionNewf(ZeroDivisionError, "float division by zero")

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

func (a Float) M__neg__() (Object, error) {
	return -a, nil
}

func (a Float) M__pos__() (Object, error) {
	return a, nil
}

func (a Float) M__abs__() (Object, error) {
	return Float(math.Abs(float64(a))), nil
}

func (a Float) M__add__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(a + b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__radd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Float) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Float) M__sub__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(a - b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__rsub__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(b - a), nil
	}
	return NotImplemented, nil
}

func (a Float) M__isub__(other Object) (Object, error) {
	return a.M__sub__(other)
}

func (a Float) M__mul__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(a * b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Float) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Float) M__truediv__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		if b == 0 {
			return nil, floatDivisionByZero
		}
		return Float(a / b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__rtruediv__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		if a == 0 {
			return nil, floatDivisionByZero
		}
		return Float(b / a), nil
	}
	return NotImplemented, nil
}

func (a Float) M__itruediv__(other Object) (Object, error) {
	return Float(a).M__truediv__(other)
}

func (a Float) M__floordiv__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Floor(float64(a / b))), nil
	}
	return NotImplemented, nil
}

func (a Float) M__rfloordiv__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Floor(float64(b / a))), nil
	}
	return NotImplemented, nil
}

func (a Float) M__ifloordiv__(other Object) (Object, error) {
	return a.M__floordiv__(other)
}

// Does DivMod of two floating point numbers
func floatDivMod(a, b Float) (Float, Float, error) {
	if b == 0 {
		return 0, 0, floatDivisionByZero
	}
	q := Float(math.Floor(float64(a / b)))
	r := a - q*b
	return q, Float(r), nil
}

func (a Float) M__mod__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		_, r, err := floatDivMod(a, b)
		return r, err
	}
	return NotImplemented, nil
}

func (a Float) M__rmod__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		_, r, err := floatDivMod(b, a)
		return r, err
	}
	return NotImplemented, nil
}

func (a Float) M__imod__(other Object) (Object, error) {
	return a.M__mod__(other)
}

func (a Float) M__divmod__(other Object) (Object, Object, error) {
	if b, ok := convertToFloat(other); ok {
		return floatDivMod(a, b)
	}
	return NotImplemented, None, nil
}

func (a Float) M__rdivmod__(other Object) (Object, Object, error) {
	if b, ok := convertToFloat(other); ok {
		return floatDivMod(b, a)
	}
	return NotImplemented, None, nil
}

func (a Float) M__pow__(other, modulus Object) (Object, error) {
	if modulus != None {
		return NotImplemented, nil
	}
	if b, ok := convertToFloat(other); ok {
		return Float(math.Pow(float64(a), float64(b))), nil
	}
	return NotImplemented, nil
}

func (a Float) M__rpow__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return Float(math.Pow(float64(b), float64(a))), nil
	}
	return NotImplemented, nil
}

func (a Float) M__ipow__(other, modulus Object) (Object, error) {
	return a.M__pow__(other, modulus)
}

func (a Float) M__bool__() (Object, error) {
	return NewBool(a != 0), nil
}

func (a Float) M__int__() (Object, error) {
	return Int(a), nil
}

func (a Float) M__float__() (Object, error) {
	return a, nil
}

func (a Float) M__complex__() (Object, error) {
	if r, ok := convertToComplex(a); ok {
		return r, nil
	}
	return cantConvert(a, "complex")
}

func (a Float) M__round__(digitsObj Object) (Object, error) {
	digits := 0
	if digitsObj != None {
		var err error
		digits, err = IndexInt(digitsObj)
		if err != nil {
			return nil, err
		}
	}
	scale := Float(math.Pow(10, float64(digits)))
	return scale * Float(math.Floor(float64(a)/float64(scale))), nil
}

// Rich comparison

func (a Float) M__lt__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a < b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__le__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a <= b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__eq__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a == b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__ne__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a != b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__gt__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a > b), nil
	}
	return NotImplemented, nil
}

func (a Float) M__ge__(other Object) (Object, error) {
	if b, ok := convertToFloat(other); ok {
		return NewBool(a >= b), nil
	}
	return NotImplemented, nil
}

// Check interface is satisfied
var _ floatArithmetic = Float(0)
var _ conversionBetweenTypes = Float(0)
var _ I__bool__ = Float(0)
var _ richComparison = Float(0)
