// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// List objects

package py

import (
	"sort"
)

var ListType = ObjectType.NewType("list", "list() -> new empty list\nlist(iterable) -> new list initialized from iterable's items", ListNew, nil)

// FIXME lists are mutable so this should probably be struct { Tuple } then can use the sub methods on Tuple
type List struct {
	Items []Object
}

func init() {
	// FIXME: all methods should be callable using list.method([], *args, **kwargs) or [].method(*args, **kwargs)
	ListType.Dict["append"] = MustNewMethod("append", func(self Object, args Tuple) (Object, error) {
		listSelf := self.(*List)
		if len(args) != 1 {
			return nil, ExceptionNewf(TypeError, "append() takes exactly one argument (%d given)", len(args))
		}
		listSelf.Items = append(listSelf.Items, args[0])
		return NoneType{}, nil
	}, 0, "append(item)")

	ListType.Dict["extend"] = MustNewMethod("extend", func(self Object, args Tuple) (Object, error) {
		listSelf := self.(*List)
		if len(args) != 1 {
			return nil, ExceptionNewf(TypeError, "append() takes exactly one argument (%d given)", len(args))
		}
		if oList, ok := args[0].(*List); ok {
			listSelf.Items = append(listSelf.Items, oList.Items...)
		}
		return NoneType{}, nil
	}, 0, "extend([item])")

	ListType.Dict["sort"] = MustNewMethod("sort", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
		const funcName = "sort"
		var l *List
		if self == None {
			// method called using `list.sort([], **kwargs)`
			var o Object
			err := UnpackTuple(args, nil, funcName, 1, 1, &o)
			if err != nil {
				return nil, err
			}
			var ok bool
			l, ok = o.(*List)
			if !ok {
				return nil, ExceptionNewf(TypeError, "descriptor 'sort' requires a 'list' object but received a '%s'", o.Type())
			}
		} else {
			// method called using `[].sort(**kargs)`
			err := UnpackTuple(args, nil, funcName, 0, 0)
			if err != nil {
				return nil, err
			}
			l = self.(*List)
		}
		err := SortInPlace(l, kwargs, funcName)
		if err != nil {
			return nil, err
		}
		return NoneType{}, nil
	}, 0, "sort(key=None, reverse=False)")

}

// Type of this List object
func (o *List) Type() *Type {
	return ListType
}

// ListNew
func ListNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var iterable Object
	err = UnpackTuple(args, kwargs, "list", 0, 1, &iterable)
	if err != nil {
		return nil, err
	}
	if iterable != nil {
		return SequenceList(iterable)
	}
	return NewList(), nil
}

// Make a new empty list
func NewList() *List {
	return &List{}
}

// Make a new empty list with given capacity
func NewListWithCapacity(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, 0, n)
	}
	return l
}

// Make a list with n nil elements
func NewListSized(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, n)
	}
	return l
}

// Make a new list from an []Object
//
// The []Object is copied into the list
func NewListFromItems(items []Object) *List {
	l := NewListSized(len(items))
	copy(l.Items, items)
	return l
}

// Copy a list object
func (l *List) Copy() *List {
	return NewListFromItems(l.Items)
}

// Append an item
func (l *List) Append(item Object) {
	l.Items = append(l.Items, item)
}

// Resize the list
func (l *List) Resize(newSize int) {
	l.Items = l.Items[:newSize]
}

// Extend the list with items
func (l *List) Extend(items []Object) {
	l.Items = append(l.Items, items...)
}

// Extends the list with the sequence passed in
func (l *List) ExtendSequence(seq Object) error {
	return Iterate(seq, func(item Object) bool {
		l.Append(item)
		return false
	})
}

// Len of list
func (l *List) Len() int {
	return len(l.Items)
}

func (l *List) M__str__() (Object, error) {
	return l.M__repr__()
}

func (l *List) M__repr__() (Object, error) {
	return Tuple(l.Items).repr("[", "]")
}

func (l *List) M__len__() (Object, error) {
	return Int(len(l.Items)), nil
}

func (l *List) M__bool__() (Object, error) {
	return NewBool(len(l.Items) > 0), nil
}

func (l *List) M__iter__() (Object, error) {
	return NewIterator(l.Items), nil
}

func (l *List) M__getitem__(key Object) (Object, error) {
	if slice, ok := key.(*Slice); ok {
		start, _, step, slicelength, err := slice.GetIndices(len(l.Items))
		if err != nil {
			return nil, err
		}
		newList := NewListSized(slicelength)
		for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
			newList.Items[j] = l.Items[i]
		}
		return newList, nil
	}
	i, err := IndexIntCheck(key, len(l.Items))
	if err != nil {
		return nil, err
	}
	return l.Items[i], nil
}

func (l *List) M__setitem__(key, value Object) (Object, error) {
	if slice, ok := key.(*Slice); ok {
		start, stop, step, slicelength, err := slice.GetIndices(len(l.Items))
		if err != nil {
			return nil, err
		}
		if step == 1 {
			// Make a copy of the tail
			tailSlice := l.Items[stop:]
			tail := make([]Object, len(tailSlice))
			copy(tail, tailSlice)
			l.Items = l.Items[:start]
			err = l.ExtendSequence(value)
			if err != nil {
				return nil, err
			}
			l.Items = append(l.Items, tail...)
		} else {
			newItems, err := SequenceTuple(value)
			if err != nil {
				return nil, err
			}
			if len(newItems) != slicelength {
				return nil, ExceptionNewf(ValueError, "attempt to assign sequence of size %d to extended slice of size %d", len(newItems), slicelength)
			}
			j := 0
			for i := start; i < stop; i += step {
				l.Items[i] = newItems[j]
				j++
			}
		}
	} else {
		i, err := IndexIntCheck(key, len(l.Items))
		if err != nil {
			return nil, err
		}
		l.Items[i] = value
	}
	return None, nil
}

// Removes the item at i
func (a *List) DelItem(i int) {
	a.Items = append(a.Items[:i], a.Items[i+1:]...)
}

// Removes items from a list
func (a *List) M__delitem__(key Object) (Object, error) {
	if slice, ok := key.(*Slice); ok {
		start, stop, step, _, err := slice.GetIndices(len(a.Items))
		if err != nil {
			return nil, err
		}
		if step == 1 {
			a.Items = append(a.Items[:start], a.Items[stop:]...)
		} else {
			j := 0
			for i := start; i < stop; i += step {
				a.DelItem(i - j)
				j++
			}
		}
	} else {
		i, err := IndexIntCheck(key, len(a.Items))
		if err != nil {
			return nil, err
		}
		a.DelItem(i)
	}
	return None, nil
}

func (a *List) M__add__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		newList := NewListSized(len(a.Items) + len(b.Items))
		copy(newList.Items, a.Items)
		copy(newList.Items[len(a.Items):], b.Items)
		return newList, nil
	}
	return NotImplemented, nil
}

func (a *List) M__radd__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		return b.M__add__(a)
	}
	return NotImplemented, nil
}

func (a *List) M__iadd__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		a.Extend(b.Items)
		return a, nil
	}
	return NotImplemented, nil
}

func (l *List) M__mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		m := len(l.Items)
		n := int(b) * m
		if n < 0 {
			n = 0
		}
		newList := NewListSized(n)
		for i := 0; i < n; i += m {
			copy(newList.Items[i:i+m], l.Items)
		}
		return newList, nil
	}
	return NotImplemented, nil
}

func (a *List) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a *List) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

// Check interface is satisfied
var _ sequenceArithmetic = (*List)(nil)
var _ I__str__ = (*List)(nil)
var _ I__repr__ = (*List)(nil)
var _ I__len__ = (*List)(nil)
var _ I__len__ = (*List)(nil)
var _ I__bool__ = (*List)(nil)
var _ I__iter__ = (*List)(nil)
var _ I__getitem__ = (*List)(nil)
var _ I__setitem__ = (*List)(nil)

// var _ richComparison = (*List)(nil)

func (a *List) M__eq__(other Object) (Object, error) {
	b, ok := other.(*List)
	if !ok {
		return NotImplemented, nil
	}
	if len(a.Items) != len(b.Items) {
		return False, nil
	}
	for i := range a.Items {
		eq, err := Eq(a.Items[i], b.Items[i])
		if err != nil {
			return nil, err
		}
		if eq == False {
			return False, nil
		}
	}
	return True, nil
}

func (a *List) M__ne__(other Object) (Object, error) {
	b, ok := other.(*List)
	if !ok {
		return NotImplemented, nil
	}
	if len(a.Items) != len(b.Items) {
		return True, nil
	}
	for i := range a.Items {
		eq, err := Eq(a.Items[i], b.Items[i])
		if err != nil {
			return nil, err
		}
		if eq == False {
			return True, nil
		}
	}
	return False, nil
}

type sortable struct {
	l        *List
	keyFunc  Object
	reverse  bool
	firstErr error
}

type ptrSortable struct {
	s *sortable
}

func (s ptrSortable) Len() int {
	return s.s.l.Len()
}

func (s ptrSortable) Swap(i, j int) {
	itemI, err := s.s.l.M__getitem__(Int(i))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return
	}
	itemJ, err := s.s.l.M__getitem__(Int(j))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return
	}
	_, err = s.s.l.M__setitem__(Int(i), itemJ)
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
	}
	_, err = s.s.l.M__setitem__(Int(j), itemI)
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
	}
}

func (s ptrSortable) Less(i, j int) bool {
	itemI, err := s.s.l.M__getitem__(Int(i))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}
	itemJ, err := s.s.l.M__getitem__(Int(j))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}

	if s.s.keyFunc != None {
		itemI, err = Call(s.s.keyFunc, Tuple{itemI}, nil)
		if err != nil {
			if s.s.firstErr == nil {
				s.s.firstErr = err
			}
			return false
		}
		itemJ, err = Call(s.s.keyFunc, Tuple{itemJ}, nil)
		if err != nil {
			if s.s.firstErr == nil {
				s.s.firstErr = err
			}
			return false
		}
	}

	var cmpResult Object
	if s.s.reverse {
		cmpResult, err = Lt(itemJ, itemI)
	} else {
		cmpResult, err = Lt(itemI, itemJ)
	}

	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}

	if boolResult, ok := cmpResult.(Bool); ok {
		return bool(boolResult)
	}

	return false
}

// SortInPlace sorts the given List in place using a stable sort.
// kwargs can have the keys "key" and "reverse".
func SortInPlace(l *List, kwargs StringDict, funcName string) error {
	var keyFunc Object
	var reverse Object
	err := ParseTupleAndKeywords(nil, kwargs, "|$OO:"+funcName, []string{"key", "reverse"}, &keyFunc, &reverse)
	if err != nil {
		return err
	}
	if keyFunc == nil {
		keyFunc = None
	}
	if reverse == nil {
		reverse = False
	}
	// FIXME: requires the same bool-check like CPython (or better "|$Op" that doesn't panic on nil).
	s := ptrSortable{&sortable{l, keyFunc, ObjectIsTrue(reverse), nil}}
	sort.Stable(s)
	return s.s.firstErr
}
