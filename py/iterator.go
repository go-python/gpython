// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Iterator objects

package py

// A python Iterator object
type Iterator struct {
	Pos  int
	Objs []Object
}

var IteratorType = NewType("iterator", "iterator type")

// Type of this object
func (o *Iterator) Type() *Type {
	return IteratorType
}

// Define a new iterator
func NewIterator(Objs []Object) *Iterator {
	m := &Iterator{
		Pos:  0,
		Objs: Objs,
	}
	return m
}

func (it *Iterator) M__iter__() (Object, error) {
	return it, nil
}

// Get next one from the iteration
func (it *Iterator) M__next__() (Object, error) {
	if it.Pos >= len(it.Objs) {
		return nil, StopIteration
	}
	r := it.Objs[it.Pos]
	it.Pos++
	return r, nil
}

// Check interface is satisfied
var _ I_iterator = (*Iterator)(nil)
