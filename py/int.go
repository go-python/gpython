// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Int objects

package py

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
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
	IntMax = math.MaxInt64
	// Minimum possible Int
	IntMin = math.MinInt64
	// The largest number such that sqrtIntMax**2 < IntMax
	sqrtIntMax = 3037000499
	// Go integer limits
	GoUintMax = ^uint(0)
	GoUintMin = 0
	GoIntMax  = int(GoUintMax >> 1)
	GoIntMin  = -GoIntMax - 1
)

// Type of this Int object
func (o Int) Type() *Type {
	return IntType
}

// IntNew
func IntNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var xObj Object = Int(0)
	var baseObj Object
	base := 0
	err := ParseTupleAndKeywords(args, kwargs, "|OO:int", []string{"x", "base"}, &xObj, &baseObj)
	if err != nil {
		return nil, err
	}
	if baseObj != nil {
		base, err = MakeGoInt(baseObj)
		if err != nil {
			return nil, err
		}
		if base != 0 && (base < 2 || base > 36) {
			return nil, ExceptionNewf(ValueError, "int() base must be >= 2 and <= 36")
		}
	}
	// Special case converting string types
	switch x := xObj.(type) {
	// FIXME Bytearray
	case Bytes:
		return IntFromString(string(x), base)
	case String:
		return IntFromString(string(x), base)
	}
	if baseObj != nil {
		return nil, ExceptionNewf(TypeError, "int() can't convert non-string with explicit base")
	}
	return MakeInt(xObj)
}

// Create an Int (or BigInt) from the string passed in
//
// FIXME check this is 100% python compatible
func IntFromString(str string, base int) (Object, error) {
	var x *big.Int
	var ok bool
	s := str
	negative := false
	convertBase := base

	// Get rid of padding
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		goto error
	}

	// Get rid of sign
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			negative = true
		}
		s = s[1:]
		if len(s) == 0 {
			goto error
		}
	}

	// Get rid of leading sigils and set convertBase
	if len(s) > 1 && s[0] == '0' {
		switch s[1] {
		case 'x', 'X':
			convertBase = 16
		case 'o', 'O':
			convertBase = 8
		case 'b', 'B':
			convertBase = 2
		default:
			goto nosigil
		}
		if base != 0 && base != convertBase {
			// int("0xFF", 10)
			// int("0b", 16)
			convertBase = base // ignore sigil
			goto nosigil
		}
		s = s[2:]
		if len(s) == 0 {
			goto error
		}
	nosigil:
	}
	if convertBase == 0 {
		convertBase = 10
	}

	// Detect leading zeros which Python doesn't allow using base 0
	if base == 0 {
		if len(s) > 1 && s[0] == '0' && (s[1] >= '0' && s[1] <= '9') {
			goto error
		}
	}

	// Use int64 conversion for short strings since 12**36 < IntMax
	// and 10**18 < IntMax
	if len(s) <= 12 || (convertBase <= 10 && len(s) <= 18) {
		i, err := strconv.ParseInt(s, convertBase, 64)
		if err != nil {
			goto error
		}
		if negative {
			i = -i
		}
		return Int(i), nil
	}

	// The base argument must be 0 or a value from 2 through
	// 36. If the base is 0, the string prefix determines the
	// actual conversion base. A prefix of “0x” or “0X” selects
	// base 16; the “0” prefix selects base 8, and a “0b” or “0B”
	// prefix selects base 2. Otherwise the selected base is 10.
	x, ok = new(big.Int).SetString(s, convertBase)
	if !ok {
		goto error
	}
	if negative {
		x.Neg(x)
	}
	return (*BigInt)(x).MaybeInt(), nil
error:
	return nil, ExceptionNewf(ValueError, "invalid literal for int() with base %d: '%s'", convertBase, str)
}

// Truncates to go int
//
// If it is outside the range of an go int it will return an error
func (x Int) GoInt() (int, error) {
	r := int(x)
	if Int(r) != x {
		return 0, overflowErrorGo
	}
	return int(r), nil
}

// Truncates to go int64
//
// If it is outside the range of an go int64 it will return an error
func (x Int) GoInt64() (int64, error) {
	return int64(x), nil
}

func (a Int) M__str__() (Object, error) {
	return String(fmt.Sprintf("%d", a)), nil
}

func (a Int) M__repr__() (Object, error) {
	return a.M__str__()
}

// Arithmetic

// Errors
var (
	divisionByZero     = ExceptionNewf(ZeroDivisionError, "division by zero")
	negativeShiftCount = ExceptionNewf(ValueError, "negative shift count")
)

// Constructs a TypeError
func cantConvert(a Object, to string) (Object, error) {
	return nil, ExceptionNewf(TypeError, "cant convert %s to %s", a.Type().Name, to)
}

// Convert an Object to an Int
//
// Returns ok as to whether the conversion worked or not
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
		// case Float:
		// 	ib := Int(b)
		// 	if Float(ib) == b {
		// 		return ib, true
		// 	}
	}
	return 0, false
}

// FIXME overflow should promote to BigInt in all these functions

func (a Int) M__neg__() (Object, error) {
	if a == IntMin {
		// Upconvert overflowing case
		r := big.NewInt(IntMin)
		r.Neg(r)
		return (*BigInt)(r), nil
	}
	return -a, nil
}

func (a Int) M__pos__() (Object, error) {
	return a, nil
}

func (a Int) M__abs__() (Object, error) {
	if a == IntMin {
		return a.M__neg__()
	}
	if a < 0 {
		return -a, nil
	}
	return a, nil
}

func (a Int) M__invert__() (Object, error) {
	return ^a, nil
}

// Integer add with overflow detection
func intAdd(a, b Int) Object {
	if a >= 0 {
		// Overflow when a + b > IntMax
		// b > IntMax - a
		// IntMax - a can't overflow since
		// IntMax = 7FFF, a = 0..7FFF
		if b > IntMax-a {
			goto overflow
		}
	} else {
		// Underflow when a + b < IntMin
		// => b < IntMin-a
		// IntMin-a can't overflow since
		// IntMin=-8000, a = -8000..-1
		if b < IntMin-a {
			goto overflow
		}
	}
	return Int(a + b)

overflow:
	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Add(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Integer subtract with overflow detection
func intSub(a, b Int) Object {
	if b >= 0 {
		// Underflow when a - b < IntMin
		// a < IntMin + b
		// IntMin + b can't overflow since
		// IntMin = -8000, b 0..7FFF
		if a < IntMin+b {
			goto overflow
		}
	} else {
		// Overflow when a - b > IntMax
		// a < IntMax + b
		// IntMax + b can't overflow since
		// IntMax=7FFF, b = -8000..-1, IntMax + b = -1..0x7FFE
		if a < IntMax+b {
			goto overflow
		}
	}
	return Int(a - b)

overflow:
	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Sub(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Integer multiplication with overflow detection
func intMul(a, b Int) Object {
	absA := a
	if a < 0 {
		absA = -a
	}
	absB := b
	if b < 0 {
		absB = -b
	}
	// A crude but effective test!
	if absA <= sqrtIntMax && absB <= sqrtIntMax {
		return Int(a * b)
	}
	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Mul(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Left shift a << b
func intLshift(a, b Int) (Object, error) {
	if b < 0 {
		return nil, negativeShiftCount
	}
	shift := uint(b)
	r := a << shift
	if r>>shift != a {
		aBig := big.NewInt(int64(a))
		aBig.Lsh(aBig, shift)
		return (*BigInt)(aBig), nil
	}
	return Int(r), nil
}

func (a Int) M__add__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intAdd(a, b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__radd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Int) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a Int) M__sub__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intSub(a, b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__rsub__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intSub(b, a), nil
	}
	return NotImplemented, nil
}

func (a Int) M__isub__(other Object) (Object, error) {
	return a.M__sub__(other)
}

func (a Int) M__mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intMul(a, b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Int) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Int) M__truediv__(other Object) (Object, error) {
	b, err := MakeFloat(other)
	if err != nil {
		return nil, err
	}
	fa := Float(a)
	if err != nil {
		return nil, err
	}
	fb := b.(Float)
	if fb == 0 {
		return nil, divisionByZero
	}
	return Float(fa / fb), nil
}

func (a Int) M__rtruediv__(other Object) (Object, error) {
	b, err := MakeFloat(other)
	if err != nil {
		return nil, err
	}
	fa := Float(a)
	if err != nil {
		return nil, err
	}
	fb := b.(Float)
	if fa == 0 {
		return nil, divisionByZero
	}
	return Float(fb / fa), nil
}

func (a Int) M__itruediv__(other Object) (Object, error) {
	return a.M__truediv__(other)
}

func (a Int) M__floordiv__(other Object) (Object, error) {
	result, _, err := a.M__divmod__(other)
	return result, err
}

func (a Int) M__rfloordiv__(other Object) (Object, error) {
	result, _, err := a.M__rdivmod__(other)
	return result, err
}

func (a Int) M__ifloordiv__(other Object) (Object, error) {
	result, _, err := a.M__divmod__(other)
	return result, err
}

func (a Int) M__mod__(other Object) (Object, error) {
	_, result, err := a.M__divmod__(other)
	return result, err
}

func (a Int) M__rmod__(other Object) (Object, error) {
	_, result, err := a.M__rdivmod__(other)
	return result, err
}

func (a Int) M__imod__(other Object) (Object, error) {
	_, result, err := a.M__divmod__(other)
	return result, err
}

func (a Int) divMod(b Int) (Object, Object, error) {
	if b == 0 {
		return nil, nil, divisionByZero
	}
	// Can't overflow
	result, remainder := Int(a/b), Int(a%b)
	// Implement floor division
	negativeResult := (a < 0)
	if b < 0 {
		negativeResult = !negativeResult
	}
	if negativeResult && remainder != 0 {
		result -= 1
		remainder += b
	}
	return result, remainder, nil
}

func (a Int) M__divmod__(other Object) (Object, Object, error) {
	if b, ok := convertToInt(other); ok {
		return a.divMod(b)
	}
	return NotImplemented, NotImplemented, nil
}

func (a Int) M__rdivmod__(other Object) (Object, Object, error) {
	if b, ok := convertToInt(other); ok {
		return b.divMod(a)
	}
	return NotImplemented, NotImplemented, nil
}

func (a Int) M__pow__(other, modulus Object) (Object, error) {
	return (*BigInt)(big.NewInt(int64(a))).M__pow__(other, modulus)
}

func (a Int) M__rpow__(other Object) (Object, error) {
	return (*BigInt)(big.NewInt(int64(a))).M__rpow__(other)
}

func (a Int) M__ipow__(other, modulus Object) (Object, error) {
	return a.M__pow__(other, modulus)
}

func (a Int) M__lshift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intLshift(a, b)
	}
	return NotImplemented, nil
}

func (a Int) M__rlshift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intLshift(b, a)
	}
	return NotImplemented, nil
}

func (a Int) M__ilshift__(other Object) (Object, error) {
	return a.M__lshift__(other)
}

func (a Int) M__rshift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			return nil, negativeShiftCount
		}
		// Can't overflow
		return Int(a >> uint64(b)), nil
	}
	return NotImplemented, nil
}

func (a Int) M__rrshift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			return nil, negativeShiftCount
		}
		// Can't overflow
		return Int(b >> uint64(a)), nil
	}
	return NotImplemented, nil
}

func (a Int) M__irshift__(other Object) (Object, error) {
	return a.M__rshift__(other)
}

func (a Int) M__and__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return Int(a & b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__rand__(other Object) (Object, error) {
	return a.M__and__(other)
}

func (a Int) M__iand__(other Object) (Object, error) {
	return a.M__and__(other)
}

func (a Int) M__xor__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return Int(a ^ b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__rxor__(other Object) (Object, error) {
	return a.M__xor__(other)
}

func (a Int) M__ixor__(other Object) (Object, error) {
	return a.M__xor__(other)
}

func (a Int) M__or__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return Int(a | b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__ror__(other Object) (Object, error) {
	return a.M__or__(other)
}

func (a Int) M__ior__(other Object) (Object, error) {
	return a.M__or__(other)
}

func (a Int) M__bool__() (Object, error) {
	return NewBool(a != 0), nil
}

func (a Int) M__index__() (Int, error) {
	return a, nil
}

func (a Int) M__int__() (Object, error) {
	return a, nil
}

func (a Int) M__float__() (Object, error) {
	if r, ok := convertToFloat(a); ok {
		return r, nil
	}
	return cantConvert(a, "float")
}

func (a Int) M__complex__() (Object, error) {
	if r, ok := convertToComplex(a); ok {
		return r, nil
	}
	return cantConvert(a, "complex")
}

func (a Int) M__round__(digits Object) (Object, error) {
	if b, ok := convertToInt(digits); ok {
		if b >= 0 {
			return a, nil
		}
		// Promote to BigInt if 10**-b > 2**63 or a == IntMin
		if b <= -19 || a == IntMin {
			return (*BigInt)(big.NewInt(int64(a))).M__round__(digits)
		}
		negative := false
		r := a
		if r < 0 {
			r = -r
			negative = true
		}
		scale := Int(math.Pow(10, float64(-b)))
		digits := r % scale
		r -= digits
		// Round
		if 2*digits >= scale {
			r += scale
		}
		if negative {
			r = -r
		}
		return r, nil
	}
	return cantConvert(digits, "int")
}

// Rich comparison

func (a Int) M__lt__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a < b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__le__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a <= b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__eq__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a == b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__ne__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a != b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__gt__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a > b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__ge__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(a >= b), nil
	}
	return NotImplemented, nil
}

func (a Int) M__ceil__() (Object, error) {
	return a, nil
}

func (a Int) M__floor__() (Object, error) {
	return a, nil
}

func (a Int) M__trunc__() (Object, error) {
	return a, nil
}

// Check interface is satisfied
var _ floatArithmetic = Int(0)
var _ booleanArithmetic = Int(0)
var _ conversionBetweenTypes = Int(0)
var _ I__bool__ = Int(0)
var _ I__index__ = Int(0)
var _ richComparison = Int(0)
var _ IGoInt = Int(0)
var _ IGoInt64 = Int(0)
