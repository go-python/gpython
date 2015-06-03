// BigInt objects

package py

import (
	"math"
	"math/big"
)

type BigInt big.Int

var BigIntType = NewType("bigint", "Holds large integers")

// Type of this BigInt object
func (o *BigInt) Type() *Type {
	return BigIntType
}

// Some common BigInts
var (
	bigInt0   = (*BigInt)(big.NewInt(0))
	bigInt1   = (*BigInt)(big.NewInt(1))
	bigIntMin = (*BigInt)(big.NewInt(IntMin))
	bigIntMax = (*BigInt)(big.NewInt(IntMax))
)

// Errors
var (
	overflowError      = ExceptionNewf(OverflowError, "Python int too large to convert to int64")
	overflowErrorGo    = ExceptionNewf(OverflowError, "Python int too large to convert to a go int")
	overflowErrorFloat = ExceptionNewf(OverflowError, "long int too large to convert to float")
)

// Arithmetic

// Convert an Object to an BigInt
//
// Retrurns ok as to whether the conversion worked or not
func convertToBigInt(other Object) (*BigInt, bool) {
	switch b := other.(type) {
	case Int:
		return (*BigInt)(big.NewInt(int64(b))), true
	case *BigInt:
		return b, true
	case Bool:
		if b {
			return bigInt1, true
		} else {
			return bigInt0, true
		}
	}
	return nil, false
}

// Truncates to Int
//
// If it is outside the range of an Int it will return an error
func (x *BigInt) Int() (Int, error) {
	if (*big.Int)(x).BitLen() <= 63 || ((*big.Int)(x).Cmp((*big.Int)(bigIntMax)) <= 0 && (*big.Int)(x).Cmp((*big.Int)(bigIntMin)) >= 0) {
		return Int((*big.Int)(x).Int64()), nil
	}
	return 0, overflowError
}

// MaybeInt truncates to Int if it can, otherwise returns the original BigInt
func (x *BigInt) MaybeInt() Object {
	i, err := x.Int()
	if err != nil {
		return x
	}
	return i
}

// Truncates to go int
//
// If it is outside the range of an go int it will return an error
func (x *BigInt) GoInt() (int, error) {
	z, err := x.Int()
	if err != nil {
		return 0, overflowErrorGo
	}
	r := int(z)
	if Int(r) != z {
		return 0, overflowErrorGo
	}
	return int(r), nil
}

// Truncates to Float
//
// If it is outside the range of an Float it will return an error
func (a *BigInt) Float() (Float, error) {
	aBig := (*big.Int)(a)
	bits := aBig.BitLen()
	exp := bits - 63
	// FIXME this is a bit approximate but errs on the low side so
	// we won't ever produce +Infs
	if exp > float64MaxExponent-63 {
		return 0, overflowErrorFloat
	}
	t := new(big.Int).Set(aBig)
	switch {
	case exp > 0:
		t.Rsh(t, uint(exp))
	case exp < 0:
		t.Lsh(t, uint(-exp))
	}
	// t should now have 63 bits of the integer in and will fit in
	// an int64
	return Float(math.Ldexp(float64(t.Int64()), exp)), nil
}

func (a *BigInt) M__neg__() (Object, error) {
	return (*BigInt)(new(big.Int).Neg((*big.Int)(a))), nil
}

func (a *BigInt) M__pos__() (Object, error) {
	return a, nil
}

func (a *BigInt) M__abs__() (Object, error) {
	if (*big.Int)(a).Sign() >= 0 {
		return a, nil
	}
	return (*BigInt)(new(big.Int).Abs((*big.Int)(a))), nil
}

func (a *BigInt) M__invert__() (Object, error) {
	return (*BigInt)(new(big.Int).Not((*big.Int)(a))), nil
}

func (a *BigInt) M__add__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Add((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__radd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a *BigInt) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (a *BigInt) M__sub__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Sub((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rsub__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Sub((*big.Int)(b), (*big.Int)(a))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__isub__(other Object) (Object, error) {
	return a.M__sub__(other)
}

func (a *BigInt) M__mul__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Mul((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a *BigInt) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a *BigInt) M__truediv__(other Object) (Object, error) {
	b, err := MakeFloat(other)
	if err != nil {
		return nil, err
	}
	fa, err := a.Float()
	if err != nil {
		return nil, err
	}
	fb := b.(Float)
	if fb == 0 {
		return nil, divisionByZero
	}
	return Float(fa / fb), nil
}

func (a *BigInt) M__rtruediv__(other Object) (Object, error) {
	b, err := MakeFloat(other)
	if err != nil {
		return nil, err
	}
	fa, err := a.Float()
	if err != nil {
		return nil, err
	}
	fb := b.(Float)
	if fa == 0 {
		return nil, divisionByZero
	}
	return Float(fb / fa), nil
}

func (a *BigInt) M__itruediv__(other Object) (Object, error) {
	return a.M__truediv__(other)
}

func (a *BigInt) M__floordiv__(other Object) (Object, error) {
	result, _, err := a.M__divmod__(other)
	return result, err
}

func (a *BigInt) M__rfloordiv__(other Object) (Object, error) {
	result, _, err := a.M__rdivmod__(other)
	return result, err
}

func (a *BigInt) M__ifloordiv__(other Object) (Object, error) {
	result, _, err := a.M__divmod__(other)
	return result, err
}

func (a *BigInt) M__mod__(other Object) (Object, error) {
	_, result, err := a.M__divmod__(other)
	return result, err
}

func (a *BigInt) M__rmod__(other Object) (Object, error) {
	_, result, err := a.M__rdivmod__(other)
	return result, err
}

func (a *BigInt) M__imod__(other Object) (Object, error) {
	_, result, err := a.M__divmod__(other)
	return result, err
}

func (a *BigInt) divMod(b *BigInt) (Object, Object, error) {
	if (*big.Int)(b).Sign() == 0 {
		return nil, nil, divisionByZero
	}
	r := new(big.Int)
	q := new(big.Int)
	q.QuoRem((*big.Int)(a), (*big.Int)(b), r)
	// Implement floor division
	negativeResult := (*big.Int)(a).Sign() < 0
	if (*big.Int)(b).Sign() < 0 {
		negativeResult = !negativeResult
	}
	if negativeResult && r.Sign() != 0 {
		q.Sub(q, (*big.Int)(bigInt1))
		r.Add(r, (*big.Int)(b))
	}
	return (*BigInt)(q).MaybeInt(), (*BigInt)(r).MaybeInt(), nil
}

func (a *BigInt) M__divmod__(other Object) (Object, Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return a.divMod(b)
	}
	return NotImplemented, NotImplemented, nil
}

func (a *BigInt) M__rdivmod__(other Object) (Object, Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return b.divMod(a)
	}
	return NotImplemented, NotImplemented, nil
}

func (a *BigInt) M__pow__(other, modulus Object) (Object, error) {
	var m *BigInt
	if modulus != None {
		var ok bool
		if m, ok = convertToBigInt(modulus); !ok {
			return NotImplemented, nil
		}
	}
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Exp((*big.Int)(a), (*big.Int)(b), (*big.Int)(m))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rpow__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Exp((*big.Int)(b), (*big.Int)(a), nil)).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__ipow__(other, modulus Object) (Object, error) {
	return a.M__pow__(other, modulus)
}

func (a *BigInt) M__lshift__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		bb, err := b.GoInt()
		if err != nil {
			return nil, err
		}
		if bb < 0 {
			return nil, negativeShiftCount
		}
		return (*BigInt)(new(big.Int).Lsh((*big.Int)(a), uint(bb))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rlshift__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		aa, err := a.GoInt()
		if err != nil {
			return nil, err
		}
		if aa < 0 {
			return nil, negativeShiftCount
		}
		return (*BigInt)(new(big.Int).Lsh((*big.Int)(b), uint(aa))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__ilshift__(other Object) (Object, error) {
	return a.M__lshift__(other)
}

func (a *BigInt) M__rshift__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		bb, err := b.GoInt()
		if err != nil {
			return nil, err
		}
		if bb < 0 {
			return nil, negativeShiftCount
		}
		return (*BigInt)(new(big.Int).Rsh((*big.Int)(a), uint(bb))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rrshift__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		aa, err := a.GoInt()
		if err != nil {
			return nil, err
		}
		if aa < 0 {
			return nil, negativeShiftCount
		}
		return (*BigInt)(new(big.Int).Rsh((*big.Int)(b), uint(aa))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__irshift__(other Object) (Object, error) {
	return a.M__rshift__(other)
}

func (a *BigInt) M__and__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).And((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rand__(other Object) (Object, error) {
	return a.M__and__(other)
}

func (a *BigInt) M__iand__(other Object) (Object, error) {
	return a.M__and__(other)
}

func (a *BigInt) M__xor__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Xor((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__rxor__(other Object) (Object, error) {
	return a.M__xor__(other)
}

func (a *BigInt) M__ixor__(other Object) (Object, error) {
	return a.M__xor__(other)
}

func (a *BigInt) M__or__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Or((*big.Int)(a), (*big.Int)(b))).MaybeInt(), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__ror__(other Object) (Object, error) {
	return a.M__or__(other)
}

func (a *BigInt) M__ior__(other Object) (Object, error) {
	return a.M__or__(other)
}

func (a *BigInt) M__bool__() (Object, error) {
	return NewBool((*big.Int)(a).Sign() != 0), nil
}

func (a *BigInt) M__index__() (Int, error) {
	return a.Int()
}

func (a *BigInt) M__int__() (Object, error) {
	return a, nil
}

func (a *BigInt) M__float__() (Object, error) {
	return a.Float()
}

func (a *BigInt) M__complex__() (Object, error) {
	// FIXME this is broken
	if r, ok := convertToComplex(a); ok {
		return r, nil
	}
	return cantConvert(a, "complex")
}

func (a *BigInt) M__round__(digits Object) (Object, error) {
	if b, ok := convertToBigInt(digits); ok {
		bb, err := b.GoInt()
		if err != nil {
			return nil, err
		}
		if bb >= 0 {
			return a, nil
		}
		// FIXME return a - (a % 10**(-bb))
		return nil, NotImplementedError
	}
	return cantConvert(digits, "int")
}

// Rich comparison

func (a *BigInt) M__lt__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) < 0), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__le__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) <= 0), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__eq__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) == 0), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__ne__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) != 0), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__gt__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) > 0), nil
	}
	return NotImplemented, nil
}

func (a *BigInt) M__ge__(other Object) (Object, error) {
	if b, ok := convertToBigInt(other); ok {
		return NewBool((*big.Int)(a).Cmp((*big.Int)(b)) >= 0), nil
	}
	return NotImplemented, nil
}

// Check interface is satisfied
var _ Object = (*BigInt)(nil)
var _ floatArithmetic = (*BigInt)(nil)
var _ booleanArithmetic = (*BigInt)(nil)
var _ conversionBetweenTypes = (*BigInt)(nil)
var _ I__bool__ = (*BigInt)(nil)
var _ I__index__ = (*BigInt)(nil)
var _ richComparison = (*BigInt)(nil)
