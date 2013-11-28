// Complex objects

package py

import (
	"fmt"
	"math"
	"math/cmplx"
)

var ComplexType = NewType("complex64", "complex(real[, imag]) -> complex number\n\nCreate a complex number from a real part and an optional imaginary part.\nThis is equivalent to (real + imag*1j) where imag defaults to 0.")

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

func (a Complex) M__neg__() Object {
	return -a
}

func (a Complex) M__pos__() Object {
	return a
}

func (a Complex) M__abs__() Object {
	return Float(cmplx.Abs(complex128(a)))
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

func (a Complex) M__isub__(other Object) Object {
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

func (a Complex) M__itruediv__(other Object) Object {
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

func (a Complex) M__ifloordiv__(other Object) Object {
	return a.M__floordiv__(other)
}

// Does Mod of two floating point numbers
func complexDivMod(a, b Complex) (Complex, Complex) {
	q := complexFloor(a / b)
	r := a - Complex(q)*b
	return q, Complex(r)
}

func (a Complex) M__mod__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		_, r := complexDivMod(a, b)
		return r
	}
	return NotImplemented
}

func (a Complex) M__rmod__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		_, r := complexDivMod(b, a)
		return r
	}
	return NotImplemented
}

func (a Complex) M__imod__(other Object) Object {
	return a.M__mod__(other)
}

func (a Complex) M__divmod__(other Object) (Object, Object) {
	if b, ok := convertToComplex(other); ok {
		return complexDivMod(a, b)
	}
	return NotImplemented, None
}

func (a Complex) M__rdivmod__(other Object) (Object, Object) {
	if b, ok := convertToComplex(other); ok {
		return complexDivMod(b, a)
	}
	return NotImplemented, None
}

func (a Complex) M__pow__(other, modulus Object) Object {
	if modulus != None {
		return NotImplemented
	}
	if b, ok := convertToComplex(other); ok {
		return Complex(cmplx.Pow(complex128(a), complex128(b)))
	}
	return NotImplemented
}

func (a Complex) M__rpow__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return Complex(cmplx.Pow(complex128(b), complex128(a)))
	}
	return NotImplemented
}

func (a Complex) M__ipow__(other, modulus Object) Object {
	return a.M__pow__(other, modulus)
}

// Rich comparison

func (a Complex) M__lt__(other Object) Object {
	if _, ok := convertToComplex(other); ok {
		// FIXME type error
		panic(fmt.Sprintf("TypeError: no ordering relation is defined for complex numbers"))
	}
	return NotImplemented
}

func (a Complex) M__le__(other Object) Object {
	return a.M__lt__(other)
}

func (a Complex) M__eq__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return NewBool(a == b)
	}
	return NotImplemented
}

func (a Complex) M__ne__(other Object) Object {
	if b, ok := convertToComplex(other); ok {
		return NewBool(a != b)
	}
	return NotImplemented
}

func (a Complex) M__gt__(other Object) Object {
	return a.M__lt__(other)
}

func (a Complex) M__ge__(other Object) Object {
	return a.M__lt__(other)
}

// Check interface is satisfied
var _ floatArithmetic = Complex(complex(0, 0))
var _ richComparison = Complex(0)
