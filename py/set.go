// Set and FrozenSet types
//
// FIXME preliminary implementation only - doesn't work properly!

package py

var SetType = NewTypeX("set", "set() -> new empty set object\nset(iterable) -> new set object\n\nBuild an unordered collection of unique elements.", SetNew, nil)

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

// SetNew
func SetNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	var iterable Object
	UnpackTuple(args, kwargs, "set", 0, 1, &iterable)
	if iterable == nil {
		return NewSet()
	}
	// FIXME should be able to initialise from an iterable!
	return NewSetFromItems(iterable.(Tuple))
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

// Extend the set with items
func (s *Set) Update(items []Object) {
	for _, item := range items {
		s.items[item] = SetValue{}
	}
}

func (s *Set) M__len__() Object {
	return Int(len(s.items))
}

func (s *Set) M__bool__() Object {
	return NewBool(len(s.items) > 0)
}

func (s *Set) M__iter__() Object {
	items := make(Tuple, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return NewIterator(items)
}

// Check interface is satisfied
var _ I__len__ = (*Set)(nil)
var _ I__bool__ = (*Set)(nil)
var _ I__iter__ = (*Set)(nil)

// var _ richComparison = (*Set)(nil)

func (a *Set) M__eq__(other Object) Object {
	b, ok := other.(*Set)
	if !ok {
		return NotImplemented
	}
	if len(a.items) != len(b.items) {
		return False
	}
	// FIXME nasty O(n**2) algorithm, waiting for proper hashing!
	for i := range a.items {
		for j := range b.items {
			if Eq(i, j) == True {
				goto found
			}
		}
		return False
	found:
	}
	return True
}

func (a *Set) M__ne__(other Object) Object {
	if a.M__eq__(other) == True {
		return False
	}
	return True
}
