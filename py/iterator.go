// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Iterator objects

package py

// A python Iterator object
type Iterator struct {
	Pos int
	Seq Object
}

var IteratorType = NewType("iterator", "iterator type")

// Type of this object
func (o *Iterator) Type() *Type {
	return IteratorType
}

// Define a new iterator
func NewIterator(Seq Object) *Iterator {
	m := &Iterator{
		Pos: 0,
		Seq: Seq,
	}
	return m
}

func (it *Iterator) M__iter__() (Object, error) {
	return it, nil
}

// Get next one from the iteration
func (it *Iterator) M__next__() (res Object, err error) {
	if tuple, ok := it.Seq.(Tuple); ok {
		if it.Pos >= len(tuple) {
			return nil, StopIteration
		}
		res = tuple[it.Pos]
		it.Pos++
		return res, nil
	}
	index := Int(it.Pos)
	if I, ok := it.Seq.(I__getitem__); ok {
		res, err = I.M__getitem__(index)
	} else if res, ok, err = TypeCall1(it.Seq, "__getitem__", index); !ok {
		return nil, ExceptionNewf(TypeError, "'%s' object is not iterable", it.Type().Name)
	}
	if err != nil {
		if IsException(IndexError, err) {
			return nil, StopIteration
		}
		return nil, err
	}
	it.Pos++
	return res, nil
}

// Check interface is satisfied
var _ I_iterator = (*Iterator)(nil)
