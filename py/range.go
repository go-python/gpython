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
	if slice, ok := key.(*Slice); ok {
		return computeRangeSlice(r, slice)
	}

	index, err := Index(key)
	if err != nil {
		return nil, err
	}
	index = computeNegativeIndex(index, r.Length)

	if index < 0 || index >= r.Length {
		return nil, ExceptionNewf(IndexError, "range object index out of range")
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

func (r *Range) M__str__() (Object, error) {
	return r.M__repr__()
}

func (r *Range) M__repr__() (Object, error) {
	return r.repr()
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

func getIndexWithDefault(i Object, d Int) (Int, error) {
	if i == None {
		return d, nil
	} else if res, err := Index(i); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}

func computeNegativeIndex(index, length Int) Int {
	if index < 0 {
		index += length
	}
	return index
}

func computeBoundIndex(index, length Int) Int {
	if index < 0 {
		index = 0
	} else if index > length {
		index = length
	}
	return index
}

func computeRangeSlice(r *Range, s *Slice) (Object, error) {
	start, err := getIndexWithDefault(s.Start, 0)
	if err != nil {
		return nil, err
	}
	stop, err := getIndexWithDefault(s.Stop, r.Length)
	if err != nil {
		return nil, err
	}
	step, err := getIndexWithDefault(s.Step, 1)
	if err != nil {
		return nil, err
	}

	if step == 0 {
		return nil, ExceptionNewf(ValueError, "slice step cannot be zero")
	}
	start = computeNegativeIndex(start, r.Length)
	stop = computeNegativeIndex(stop, r.Length)

	start = computeBoundIndex(start, r.Length)
	stop = computeBoundIndex(stop, r.Length)

	startIndex := computeItem(r, start)
	stopIndex := computeItem(r, stop)
	stepIndex := step * r.Step

	var sliceLength Int
	if start < stop {
		if stepIndex < 0 {
			startIndex, stopIndex = stopIndex-1, startIndex-1
		}
	} else {
		if stepIndex < 0 {
			startIndex, stopIndex = stopIndex+1, startIndex+1
		}
	}
	sliceLength = computeRangeLength(startIndex, stopIndex, stepIndex)

	return &Range{
		Start:  startIndex,
		Stop:   stopIndex,
		Step:   stepIndex,
		Length: sliceLength,
	}, nil
}

// Check interface is satisfied
var _ I__getitem__ = (*Range)(nil)
var _ I__iter__ = (*Range)(nil)
var _ I_iterator = (*RangeIterator)(nil)

func (a *Range) M__eq__(other Object) (Object, error) {
	b, ok := other.(*Range)
	if !ok {
		return NotImplemented, nil
	}

	if a.Length != b.Length {
		return False, nil
	}

	if a.Length == 0 {
		return True, nil
	}
	if a.Start != b.Start {
		return False, nil
	}

	if a.Step == 1 {
		return True, nil
	}
	if a.Step != b.Step {
		return False, nil
	}

	return True, nil
}

func (a *Range) M__ne__(other Object) (Object, error) {
	b, ok := other.(*Range)
	if !ok {
		return NotImplemented, nil
	}

	if a.Length != b.Length {
		return True, nil
	}

	if a.Length == 0 {
		return False, nil
	}
	if a.Start != b.Start {
		return True, nil
	}

	if a.Step == 1 {
		return False, nil
	}
	if a.Step != b.Step {
		return True, nil
	}

	return False, nil
}
