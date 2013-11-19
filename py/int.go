// Int and BigInt objects

package py

import (
	"math/big"
)

var IntType = NewType("int")

type Int int64

// Type of this Int object
func (o Int) Type() *Type {
	return IntType
}

type BigInt big.Int

var BigIntType = NewType("bigint")

// Type of this BigInt object
func (o *BigInt) Type() *Type {
	return BigIntType
}

// Make sure it satisfies the interface
var _ Object = (*BigInt)(nil)

// Arithmetic

func (a Int) M__add__(other Object) Object {
	switch b := other.(type) {
	case Int:
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
	switch b := other.(type) {
	case Int:
		return Int(a - b)
	}
	return NotImplemented
}

func (b Int) M__rsub__(other Object) Object {
	switch a := other.(type) {
	case Int:
		return Int(b - a)
	}
	return NotImplemented
}

func (a Int) M__isub(other Object) Object {
	return a.M__sub__(other)
}
