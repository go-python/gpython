// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// A python Enumerate object
type Enumerate struct {
	Iterable Object
	Start    Int
}

// A python Enumerate iterator
type EnumerateIterator struct {
	Enumerate
	Index Int
}

var EnumerateType = NewTypeX("enumerate", `enumerate(iterable, start=0)

Return an enumerate object.`,
	EnumerateNew, nil)

var EnumerateIteratorType = NewType("enumerate_iterator", `enumerate_iterator object`)

// Type of this object
func (e *Enumerate) Type() *Type {
	return EnumerateType
}

// Type of this object
func (ei *EnumerateIterator) Type() *Type {
	return EnumerateIteratorType
}

// EnumerateTypeNew
func EnumerateNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var iterable Object
	var start Object
	err := UnpackTuple(args, kwargs, "enumerate", 1, 2, &iterable, &start)
	if err != nil {
		return nil, err
	}

	if start == nil {
		start = Int(0)
	}
	startIndex, err := Index(start)
	if err != nil {
		return nil, err
	}
	iter, err := Iter(iterable)
	if err != nil {
		return nil, err
	}

	return &Enumerate{Iterable: iter, Start: startIndex}, nil
}

// Enumerate iterator
func (e *Enumerate) M__iter__() (Object, error) {
	return &EnumerateIterator{
		Enumerate: *e,
		Index:     e.Start,
	}, nil
}

// EnumerateIterator iterator
func (ei *EnumerateIterator) M__iter__() (Object, error) {
	return ei, nil
}

// EnumerateIterator iterator next
func (ei *EnumerateIterator) M__next__() (Object, error) {
	value, err := Next(ei.Enumerate.Iterable)
	if err != nil {
		return nil, err
	}
	res := make(Tuple, 2)
	res[0] = ei.Index
	res[1] = value
	ei.Index += 1
	return res, nil
}

// Check interface is satisfied
var _ I__iter__ = (*Enumerate)(nil)
var _ I_iterator = (*EnumerateIterator)(nil)
