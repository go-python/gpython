// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Float objects

package py

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var FloatType = ObjectType.NewType("float", "float(x) -> floating point number\n\nConvert a string or number to a floating point number, if possible.", FloatNew, nil)

// Bits of precision in a float64
const (
	float64precision   = 53
	float64MaxExponent = 1023
)

type Float float64

// Type of this Float64 object
func (o Float) Type() *Type {
	return FloatType
}

// FloatNew
func FloatNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var xObj Object = Float(0)
	err := ParseTupleAndKeywords(args, kwargs, "|O", []string{"x"}, &xObj)
	if err != nil {
		return nil, err
	}
	// Special case converting string types
	switch x := xObj.(type) {
	// FIXME Bytearray
	case Bytes:
		return FloatFromString(string(x))
	case String:
		return FloatFromString(string(x))
	}
	return MakeFloat(xObj)
}

func (a Float) M__str__() (Object, error) {
	if i := int64(a); Float(i) == a {
		return String(fmt.Sprintf("%d.0", i)), nil
	}
	return String(fmt.Sprintf("%g", a)), nil
}

func (a Float) M__repr__() (Object, error) {
	return a.M__str__()
}

// FloatFromString turns a string into a Float
func FloatFromString(str string) (Object, error) {
	str = strings.TrimSpace(str)
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			if numErr.Err == strconv.ErrRange {
				if str[0] == '-' {
					return Float(math.Inf(-1)), nil
				} else {
					return Float(math.Inf(1)), nil
				}
			}
		}
		return nil, ExceptionNewf(ValueError, "invalid literal for float: '%s'", str)
	}
	return Float(f), nil
}

var expectingFloat = ExceptionNewf(TypeError, "a float is required")

// Returns the float value of obj if it is exactly a float
func FloatCheckExact(obj Object) (Float, error) {
	f, ok := obj.(Float)
	if !ok {
		return 0, expectingFloat
	}
	return f, nil
}

// Returns the float value of obj if it is a float subclass
func FloatCheck(obj Object) (Float, error) {
	// FIXME should be checking subclasses
	return FloatCheckExact(obj)
}

// PyFloat_AsDouble
func FloatAsFloat64(obj Object) (float64, error) {
	f, err := FloatCheck(obj)
	if err == nil {
		return float64(f), nil
	}
	fObj, err := MakeFloat(obj)
	if err != nil {
		return 0, err
	}
	f, err = FloatCheck(fObj)
	if err == nil {
		return float64(f), nil
	}
	return float64(f), err
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
	case *BigInt:
		x, err := b.Float()
		return x, err == nil
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
	if a >= IntMin && a <= IntMax {
		return Int(a), nil
	}
	frac, exp := math.Frexp(float64(a))              // x = frac << exp; 0.5 <= abs(x) < 1
	fracInt := int64(frac * (1 << float64precision)) // x = frac << (exp - float64precision)
	res := big.NewInt(fracInt)
	shift := exp - float64precision
	switch {
	case shift > 0:
		res.Lsh(res, uint(shift))
	case shift < 0:
		res.Rsh(res, uint(-shift))
	}
	return (*BigInt)(res), nil
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
		digits, err = MakeGoInt(digitsObj)
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

// Properties
func init() {
	FloatType.Dict["is_integer"] = MustNewMethod("is_integer", func(self Object) (Object, error) {
		if a, ok := convertToFloat(self); ok {
			f, err := FloatAsFloat64(a)
			if err != nil {
				return nil, err
			}
			return NewBool(math.Floor(f) == f), nil
		}
		return cantConvert(self, "float")
	}, 0, "is_integer() -> Return True if the float instance is finite with integral value, and False otherwise.")
}

// Check interface is satisfied
var _ floatArithmetic = Float(0)
var _ conversionBetweenTypes = Float(0)
var _ I__bool__ = Float(0)
var _ richComparison = Float(0)
