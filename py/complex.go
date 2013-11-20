// Complex objects

package py

import (
	"math"
)

var ComplexType = NewType("complex64")

type Complex complex128

// Type of this Complex object
func (o Complex) Type() *Type {
	return ComplexType
}

// Convert an Object to an Complex
//
// Retrurns ok as to whether the conversion worked or not
func convertToComplex(other Object) (Complex, bool) {
	switch b := other.(type) {
	case Complex:
		return b, true
	case Float:
		return Complex(complex(b, 0)), true
	case Int:
		return Complex(complex(float64(b), 0)), true
	case Bool:
		if b {
			return Complex(1), true
		} else {
			return Complex(0), true
		}
	}
	return 0, false
}

func (a Complex) M__add__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(a + b)
	}
	return NotImplemented
}

func (a Complex) M__radd__(other Object) Object {
	return a.M__add__(other)
}

func (a Complex) M__iadd__(other Object) Object {
	return a.M__add__(other)
}

func (a Complex) M__sub__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(a - b)
	}
	return NotImplemented
}

func (a Complex) M__rsub__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(b - a)
	}
	return NotImplemented
}

func (a Complex) M__isub(other Object) Object {
	return a.M__sub__(other)
}

func (a Complex) M__mul__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(a * b)
	}
	return NotImplemented
}

func (a Complex) M__rmul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Complex) M__imul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Complex) M__truediv__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(a / b)
	}
	return NotImplemented
}

func (a Complex) M__rtruediv__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(b / a)
	}
	return NotImplemented
}

func (a Complex) M__itruediv(other Object) Object {
	return Complex(a).M__truediv__(other)
}

// Floor a complex number
func complexFloor(a Complex) Complex {
	return Complex(complex(math.Floor(real(a)), math.Floor(imag(a))))
}

// Floor divide two complex numbers
func complexFloorDiv(a, b Complex) Complex {
	q := complexFloor(a / b)
	r := a - q*b
	return Complex(r)
}

func (a Complex) M__floordiv__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return complexFloor(a / b)
	}
	return NotImplemented
}

func (a Complex) M__rfloordiv__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return complexFloor(b / a)
	}
	return NotImplemented
}

func (a Complex) M__ifloordiv(other Object) Object {
	return a.M__floordiv__(other)
}

// Does Mod of two floating point numbers
func complexMod(a, b Complex) Complex {
	q := complexFloor(a / b)
	r := a - Complex(q)*b
	return Complex(r)
}

func (a Complex) M__mod__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return complexMod(a, b)
	}
	return NotImplemented
}

func (a Complex) M__rmod__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return complexMod(b, a)
	}
	return NotImplemented
}

func (a Complex) M__imod(other Object) Object {
	return a.M__mod__(other)
}
