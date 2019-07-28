// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tuple objects

package py

import "bytes"

var TupleType = ObjectType.NewType("tuple", "tuple() -> empty tuple\ntuple(iterable) -> tuple initialized from iterable's items\n\nIf the argument is a tuple, the return value is the same object.", TupleNew, nil)

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
}

// TupleNew
func TupleNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var iterable Object
	err = UnpackTuple(args, kwargs, "tuple", 0, 1, &iterable)
	if err != nil {
		return nil, err
	}
	if iterable != nil {
		return SequenceTuple(iterable)
	}
	return Tuple{}, nil
}

// Copy a tuple object
func (t Tuple) Copy() Tuple {
	newT := make(Tuple, len(t))
	copy(newT, t)
	return newT
}

// Reverses a tuple (in-place)
func (t Tuple) Reverse() {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
}

// output the tuple to out, using fn to transform the tuple to out
// start and end brackets
func (t Tuple) repr(start, end string) (Object, error) {
	var out bytes.Buffer
	out.WriteString(start)
	for i, obj := range t {
		if i != 0 {
			out.WriteString(", ")
		}
		str, err := ReprAsString(obj)
		if err != nil {
			return nil, err
		}
		out.WriteString(str)
	}
	out.WriteString(end)
	return String(out.String()), nil
}

func (t Tuple) M__str__() (Object, error) {
	return t.M__repr__()
}

func (t Tuple) M__repr__() (Object, error) {
	return t.repr("(", ")")
}

func (t Tuple) M__len__() (Object, error) {
	return Int(len(t)), nil
}

func (t Tuple) M__bool__() (Object, error) {
	return NewBool(len(t) > 0), nil
}

func (t Tuple) M__iter__() (Object, error) {
	return NewIterator(t), nil
}

func (t Tuple) M__getitem__(key Object) (Object, error) {
	if slice, ok := key.(*Slice); ok {
		start, stop, step, slicelength, err := slice.GetIndices(len(t))
		if err != nil {
			return nil, err
		}
		if step == 1 {
			// Return a subslice since tuples are immutable
			return t[start:stop], nil
		}
		newTuple := make(Tuple, slicelength)
		for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
			newTuple[j] = t[i]
		}
		return newTuple, nil
	}
	i, err := IndexIntCheck(key, len(t))
	if err != nil {
		return nil, err
	}
	return t[i], nil
}

func (a Tuple) M__add__(other Object) (Object, error) {
	if b, ok := other.(Tuple); ok {
		newTuple := make(Tuple, len(a)+len(b))
		copy(newTuple, a)
		copy(newTuple[len(b):], b)
		return newTuple, nil
	}

	return NotImplemented, nil
}

func (a Tuple) M__radd__(other Object) (Object, error) {
	if b, ok := other.(Tuple); ok {
		return b.M__add__(a)
	}
	return NotImplemented, nil
}

func (a Tuple) M__iadd__(other Object) (Object, error) {
	return a.M__add__(other)
}

func (l Tuple) M__mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		m := len(l)
		n := int(b) * m
		if n < 0 {
			n = 0
		}
		newTuple := make(Tuple, n)
		for i := 0; i < n; i += m {
			copy(newTuple[i:i+m], l)
		}
		return newTuple, nil
	}
	return NotImplemented, nil
}

func (a Tuple) M__rmul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Tuple) M__imul__(other Object) (Object, error) {
	return a.M__mul__(other)
}

func (a Tuple) M__eq__(other Object) (Object, error) {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented, nil
	}
	if len(a) != len(b) {
		return False, nil
	}
	for i := range a {
		eq, err := Eq(a[i], b[i])
		if err != nil {
			return nil, err
		}
		if eq == False {
			return False, nil
		}
	}
	return True, nil
}

func (a Tuple) M__ne__(other Object) (Object, error) {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented, nil
	}
	if len(a) != len(b) {
		return True, nil
	}
	for i := range a {
		eq, err := Eq(a[i], b[i])
		if err != nil {
			return nil, err
		}
		if eq == False {
			return True, nil
		}
	}
	return False, nil
}

// Check interface is satisfied
var _ sequenceArithmetic = Tuple(nil)
var _ I__str__ = Tuple(nil)
var _ I__repr__ = Tuple(nil)
var _ I__len__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__iter__ = Tuple(nil)
var _ I__getitem__ = Tuple(nil)
var _ I__eq__ = Tuple(nil)
var _ I__ne__ = Tuple(nil)

// var _ richComparison = Tuple(nil)
