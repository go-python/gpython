// Copyright 2019 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CallIterator objects

package py

// A python CallIterator object
type CallIterator struct {
	callable Object
	sentinel Object
}

var CallIteratorType = NewType("callable_iterator", "callable_iterator type")

// Type of this object
func (o *CallIterator) Type() *Type {
	return CallIteratorType
}

func (cit *CallIterator) M__iter__() (Object, error) {
	return cit, nil
}

// Get next one from the iteration
func (cit *CallIterator) M__next__() (Object, error) {
	value, err := Call(cit.callable, nil, nil)

	if err != nil {
		return nil, err
	}

	if value == cit.sentinel {
		return nil, StopIteration
	}

	return value, nil
}

// Define a new CallIterator
func NewCallIterator(callable Object, sentinel Object) *CallIterator {
	c := &CallIterator{
		callable: callable,
		sentinel: sentinel,
	}
	return c
}

// Check interface is satisfied
var _ I_iterator = (*CallIterator)(nil)
