// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set and FrozenSet types
//
// FIXME preliminary implementation only - doesn't work properly!

package py

import "bytes"

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
func SetNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var iterable Object
	err := UnpackTuple(args, kwargs, "set", 0, 1, &iterable)
	if err != nil {
		return nil, err
	}
	if iterable != nil {
		return SequenceSet(iterable)
	}
	return NewSet(), nil
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

func (s *Set) M__len__() (Object, error) {
	return Int(len(s.items)), nil
}

func (s *Set) M__bool__() (Object, error) {
	return NewBool(len(s.items) > 0), nil
}

func (s *Set) M__repr__() (Object, error) {
	var out bytes.Buffer
	out.WriteRune('{')
	spacer := false
	for item := range s.items {
		if spacer {
			out.WriteString(", ")
		}
		str, err := ReprAsString(item)
		if err != nil {
			return nil, err
		}
		out.WriteString(str)
		spacer = true
	}
	out.WriteRune('}')
	return String(out.String()), nil
}

func (s *Set) M__iter__() (Object, error) {
	items := make(Tuple, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return NewIterator(items), nil
}

func (s *Set) M__and__(other Object) (Object, error) {
	ret := NewSet()
	b, ok := other.(*Set)
	if !ok {
		return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for &: '%s' and '%s'", s.Type().Name, other.Type().Name)
	}
	for i := range b.items {
		if _, ok := s.items[i]; ok {
			ret.items[i] = SetValue{}
		}
	}
	return ret, nil
}

func (s *Set) M__or__(other Object) (Object, error) {
	ret := NewSet()
	b, ok := other.(*Set)
	if !ok {
		return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for &: '%s' and '%s'", s.Type().Name, other.Type().Name)
	}
	for j := range s.items {
		ret.items[j] = SetValue{}
	}
	for i := range b.items {
		if _, ok := s.items[i]; !ok {
			ret.items[i] = SetValue{}
		}
	}
	return ret, nil
}

func (s *Set) M__sub__(other Object) (Object, error) {
	ret := NewSet()
	b, ok := other.(*Set)
	if !ok {
		return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for &: '%s' and '%s'", s.Type().Name, other.Type().Name)
	}
	for j := range s.items {
		ret.items[j] = SetValue{}
	}
	for i := range b.items {
		if _, ok := s.items[i]; ok {
			delete(ret.items, i)
		}
	}
	return ret, nil
}

func (s *Set) M__xor__(other Object) (Object, error) {
	ret := NewSet()
	b, ok := other.(*Set)
	if !ok {
		return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for &: '%s' and '%s'", s.Type().Name, other.Type().Name)
	}
	for j := range s.items {
		ret.items[j] = SetValue{}
	}
	for i := range b.items {
		_, ok := s.items[i]
		if ok {
			delete(ret.items, i)
		} else {
			ret.items[i] = SetValue{}
		}
	}
	return ret, nil
}

// Check interface is satisfied
var _ I__len__ = (*Set)(nil)
var _ I__bool__ = (*Set)(nil)
var _ I__iter__ = (*Set)(nil)

// var _ richComparison = (*Set)(nil)

func (a *Set) M__eq__(other Object) (Object, error) {
	b, ok := other.(*Set)
	if !ok {
		return NotImplemented, nil
	}
	if len(a.items) != len(b.items) {
		return False, nil
	}
	// FIXME nasty O(n**2) algorithm, waiting for proper hashing!
	for i := range a.items {
		for j := range b.items {
			eq, err := Eq(i, j)
			if err != nil {
				return nil, err
			}
			if eq == True {
				goto found
			}
		}
		return False, nil
	found:
	}
	return True, nil
}

func (a *Set) M__ne__(other Object) (Object, error) {
	eq, err := a.M__eq__(other)
	if err != nil {
		return nil, err
	}
	if eq == NotImplemented {
		return eq, nil
	}
	if eq == True {
		return False, nil
	}
	return True, nil
}
