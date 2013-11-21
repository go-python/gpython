// Set and FrozenSet types

package py

var SetType = NewType("set", "set() -> new empty set object\nset(iterable) -> new set object\n\nBuild an unordered collection of unique elements.")

type SetValue struct{}

type Set map[Object]SetValue

// Type of this Set object
func (o Set) Type() *Type {
	return SetType
}

var FrozenSetType = NewType("frozenset", "frozenset() -> empty frozenset object\nfrozenset(iterable) -> frozenset object\n\nBuild an immutable unordered collection of unique elements.")

type FrozenSet map[Object]SetValue

// Type of this FrozenSet object
func (o FrozenSet) Type() *Type {
	return FrozenSetType
}
