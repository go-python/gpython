// Set and FrozenSet types
//
// FIXME preliminary implementation only - doesn't work properly!

package py

var SetType = NewType("set", "set() -> new empty set object\nset(iterable) -> new set object\n\nBuild an unordered collection of unique elements.")

type SetValue struct{}

type Set struct {
	items map[Object]SetValue
}

// Type of this Set object
func (o *Set) Type() *Type {
	return SetType
}

// Make a new empty set
func NewSet() *Set {
	return &Set{
		items: make(map[Object]SetValue),
	}
}

// Make a new empty set with capacity for n items
func NewSetWithCapacity(n int) *Set {
	return &Set{
		items: make(map[Object]SetValue, n),
	}
}

// Make a new set with the items passed in
func NewSetFromItems(items []Object) *Set {
	s := NewSetWithCapacity(len(items))
	for _, item := range items {
		s.items[item] = SetValue{}
	}
	return s
}

// Add an item to the set
func (s *Set) Add(item Object) {
	s.items[item] = SetValue{}
}

var FrozenSetType = NewType("frozenset", "frozenset() -> empty frozenset object\nfrozenset(iterable) -> frozenset object\n\nBuild an immutable unordered collection of unique elements.")

type FrozenSet struct {
	Set
}

// Type of this FrozenSet object
func (o *FrozenSet) Type() *Type {
	return FrozenSetType
}

// Make a new empty frozen set
func NewFrozenSet() *FrozenSet {
	return &FrozenSet{
		Set: *NewSet(),
	}
}

// Make a new set with the items passed in
func NewFrozenSetFromItems(items []Object) *FrozenSet {
	return &FrozenSet{
		Set: *NewSetFromItems(items),
	}
}
