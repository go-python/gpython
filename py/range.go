// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Range object

package py

// A python Range object
// FIXME one day support BigInts too!
type Range struct {
	Start  Int
	Stop   Int
	Step   Int
	Length Int
}

// A python Range iterator
type RangeIterator struct {
	Range
	Index Int
}

var RangeType = NewTypeX("range", `range(stop) -> range object
range(start, stop[, step]) -> range object

Return a virtual sequence of numbers from start to stop by step.`,
	RangeNew, nil)

var RangeIteratorType = NewType("range_iterator", `range_iterator object`)

// Type of this object
func (o *Range) Type() *Type {
	return RangeType
}

// Type of this object
func (o *RangeIterator) Type() *Type {
	return RangeIteratorType
}

// RangeNew
func RangeNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var start Object
	var stop Object
	var step Object = Int(1)
	err := UnpackTuple(args, kwargs, "range", 1, 3, &start, &stop, &step)
	if err != nil {
		return nil, err
	}
	startIndex, err := Index(start)
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		length := computeRangeLength(0, startIndex, 1)
		return &Range{
			Start:  Int(0),
			Stop:   startIndex,
			Step:   Int(1),
			Length: length,
		}, nil
	}
	stopIndex, err := Index(stop)
	if err != nil {
		return nil, err
	}
	stepIndex, err := Index(step)
	if err != nil {
		return nil, err
	}
	length := computeRangeLength(startIndex, stopIndex, stepIndex)
	return &Range{
		Start:  startIndex,
		Stop:   stopIndex,
		Step:   stepIndex,
		Length: length,
	}, nil
}

func (r *Range) M__getitem__(key Object) (Object, error) {
	index, err := Index(key)
	if err != nil {
		return nil, err
	}
	// TODO(corona10): Support slice case
	length := computeRangeLength(r.Start, r.Stop, r.Step)
	if index < 0 {
		index += length
	}

	if index < 0 || index >= length {
		return nil, ExceptionNewf(TypeError, "range object index out of range")
	}
	result := computeItem(r, index)
	return result, nil
}

// Make a range iterator from a range
func (r *Range) M__iter__() (Object, error) {
	return &RangeIterator{
		Range: *r,
		Index: r.Start,
	}, nil
}

func (r *Range) M__len__() (Object, error) {
	return r.Length, nil
}

// Range iterator
func (it *RangeIterator) M__iter__() (Object, error) {
	return it, nil
}

// Range iterator next
func (it *RangeIterator) M__next__() (Object, error) {
	r := it.Index
	if it.Step >= 0 && r >= it.Stop {
		return nil, StopIteration
	}

	if it.Step < 0 && r <= it.Stop {
		return nil, StopIteration
	}
	it.Index += it.Step
	return r, nil
}

func computeItem(r *Range, item Int) Int {
	incr := item * r.Step
	res := r.Start + incr
	return res
}

func computeRangeLength(start, stop, step Int) Int {
	var lo, hi Int
	if step > 0 {
		lo = start
		hi = stop
		step = step
	} else {
		lo = stop
		hi = start
		step = (-step)
	}

	if lo >= hi {
		return Int(0)
	}
	res := (hi-lo-1)/step + 1
	return res
}

// Check interface is satisfied
var _ I__getitem__ = (*Range)(nil)
var _ I__iter__ = (*Range)(nil)
var _ I_iterator = (*RangeIterator)(nil)
