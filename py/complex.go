// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Complex objects

package py

import (
	"fmt"
	"math"
	"math/cmplx"
)

var ComplexType = ObjectType.NewType("complex64", "complex(real[, imag]) -> complex number\n\nCreate a complex number from a real part and an optional imaginary part.\nThis is equivalent to (real + imag*1j) where imag defaults to 0.", ComplexNew, nil)

type Complex complex128

// Type of this Complex object
func (o Complex) Type() *Type {
	return ComplexType
}

// ComplexNew
func ComplexNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var realObj Object = Float(0)
	var imagObj Object = Float(0)
	err := ParseTupleAndKeywords(args, kwargs, "|OO", []string{"real", "imag"}, &realObj, &imagObj)
	if err != nil {
		return nil, err
	}
	real, err := MakeFloat(realObj)
	if err != nil {
		return nil, err
	}
	imag, err := MakeFloat(imagObj)
	if err != nil {
		return nil, err
	}
	return Complex(complex(real.(Float), imag.(Float))), nil
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

func (a Complex) M__str__() (Object, error) {
	return String(fmt.Sprintf("(%g%+gj)", real(complex128(a)), imag(complex128(a)))), nil
}

func (a Complex) M__repr__() (Object, error) {
	return a.M__str__()
}

func (a Complex) M__neg__() (Object, error) {
	return -a, nil
}

func (a Complex) M__pos__() (Object, error) {
	return a, nil
}

func (a Complex) M__abs__() (Object, error) {
	return Float(cmplx.Abs(complex128(a))), nil
}

func (a Complex) M__add__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(a + b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__radd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Complex) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Complex) M__sub__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(a - b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rsub__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(b - a), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__isub__(other Object) (Object, error) {
	return a.M__sub__(other)
}

func (a Complex) M__mul__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(a * b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Complex) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Complex) M__truediv__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(a / b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rtruediv__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(b / a), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__itruediv__(other Object) (Object, error) {
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

func (a Complex) M__floordiv__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return complexFloor(a / b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rfloordiv__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return complexFloor(b / a), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__ifloordiv__(other Object) (Object, error) {
	return a.M__floordiv__(other)
}

// Does Mod of two floating point numbers
func complexDivMod(a, b Complex) (Complex, Complex) {
	q := complexFloor(a / b)
	r := a - Complex(q)*b
	return q, Complex(r)
}

func (a Complex) M__mod__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		_, r := complexDivMod(a, b)
		return r, nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rmod__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		_, r := complexDivMod(b, a)
		return r, nil
	}
	return NotImplemented, nil
}

func (a Complex) M__imod__(other Object) (Object, error) {
	return a.M__mod__(other)
}

func (a Complex) M__divmod__(other Object) (Object, Object, error) {
	if b, ok := convertToComplex(other); ok {
		x, y := complexDivMod(a, b)
		return x, y, nil
	}
	return NotImplemented, None, nil
}

func (a Complex) M__rdivmod__(other Object) (Object, Object, error) {
	if b, ok := convertToComplex(other); ok {
		x, y := complexDivMod(b, a)
		return x, y, nil
	}
	return NotImplemented, None, nil
}

func (a Complex) M__pow__(other, modulus Object) (Object, error) {
	if modulus != None {
		return NotImplemented, nil
	}
	if b, ok := convertToComplex(other); ok {
		return Complex(cmplx.Pow(complex128(a), complex128(b))), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__rpow__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return Complex(cmplx.Pow(complex128(b), complex128(a))), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__ipow__(other, modulus Object) (Object, error) {
	return a.M__pow__(other, modulus)
}

func (a Complex) M__int__() (Object, error) {
	if r, ok := convertToInt(a); ok {
		return r, nil
	}
	return cantConvert(a, "int")

}

func (a Complex) M__float__() (Object, error) {
	if r, ok := convertToFloat(a); ok {
		return r, nil
	}
	return cantConvert(a, "float")
}

func (a Complex) M__complex__() (Object, error) {
	return a, nil
}

// Rich comparison

func (a Complex) M__lt__(other Object) (Object, error) {
	if _, ok := convertToComplex(other); ok {
		return nil, ExceptionNewf(TypeError, "no ordering relation is defined for complex numbers")
	}
	return NotImplemented, nil
}

func (a Complex) M__le__(other Object) (Object, error) {
	return a.M__lt__(other)
}

func (a Complex) M__eq__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return NewBool(a == b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__ne__(other Object) (Object, error) {
	if b, ok := convertToComplex(other); ok {
		return NewBool(a != b), nil
	}
	return NotImplemented, nil
}

func (a Complex) M__gt__(other Object) (Object, error) {
	return a.M__lt__(other)
}

func (a Complex) M__ge__(other Object) (Object, error) {
	return a.M__lt__(other)
}

// Properties
func init() {
	ComplexType.Dict["real"] = &Property{
		Fget: func(self Object) (Object, error) {
			return Float(real(self.(Complex))), nil
		},
	}
	ComplexType.Dict["imag"] = &Property{
		Fget: func(self Object) (Object, error) {
			return Float(imag(self.(Complex))), nil
		},
	}
	ComplexType.Dict["conjugate"] = MustNewMethod("conjugate", func(self Object) (Object, error) {
		cnj := cmplx.Conj(complex128(self.(Complex)))
		return Complex(cnj), nil
	}, 0, "conjugate() -> Returns the complex conjugate.")
}

// Check interface is satisfied
var _ floatArithmetic = Complex(complex(0, 0))
var _ richComparison = Complex(0)
