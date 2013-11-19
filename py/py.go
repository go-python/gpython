// Python global definitions
package py

import (
	"math/big"
)

// A python object
type Object interface {
	Type() *Type
}

// Some well known objects
var (
	None, False, True, StopIteration, Elipsis Object
)

// Some python types
// FIXME factor into own files probably

var TupleType = NewType("tuple")

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
}

var ListType = NewType("list")

type List []Object

// Type of this List object
func (o List) Type() *Type {
	return ListType
}

var BytesType = NewType("bytes")

type Bytes []byte

// Type of this Bytes object
func (o Bytes) Type() *Type {
	return BytesType
}

var Int64Type = NewType("int64")

type Int64 int64

// Type of this Int64 object
func (o Int64) Type() *Type {
	return Int64Type
}

var Float64Type = NewType("float64")

type Float64 float64

// Type of this Float64 object
func (o Float64) Type() *Type {
	return Float64Type
}

var Complex64Type = NewType("complex64")

type Complex64 complex64

// Type of this Complex64 object
func (o Complex64) Type() *Type {
	return Complex64Type
}

var StringDictType = NewType("dict")

// String to object dictionary
//
// Used for variables etc where the keys can only be strings
type StringDict map[string]Object

// Type of this StringDict object
func (o StringDict) Type() *Type {
	return StringDictType
}

// Make a new dictionary
func NewStringDict() StringDict {
	return make(StringDict)
}

var SetType = NewType("set")

type SetValue struct{}

type Set map[Object]SetValue

// Type of this Set object
func (o Set) Type() *Type {
	return SetType
}

var FrozenSetType = NewType("frozenset")

type FrozenSet map[Object]SetValue

// Type of this FrozenSet object
func (o FrozenSet) Type() *Type {
	return FrozenSetType
}

type BigInt big.Int

var BigIntType = NewType("bigint")

// Type of this BigInt object
func (o *BigInt) Type() *Type {
	return BigIntType
}

// Make sure it satisfies the interface
var _ Object = (*BigInt)(nil)

// Interfaces satisfied by a subset of Objects
