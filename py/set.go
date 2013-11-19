// Set and FrozenSet types

package py

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
