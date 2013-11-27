// Range object

package py

// A python Range object
// FIXME one day support BigInts too!
type Range struct {
	Start Int
	Stop  Int
	Step  Int
	//Length Object
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
func RangeNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	var start Object
	var stop Object
	var step Object = Int(1)
	UnpackTuple(args, kwargs, "range", 1, 3, &start, &stop, &step)
	if len(args) == 1 {
		return &Range{
			Start: Int(0),
			Stop:  Index(start),
			Step:  Int(1),
		}
	}
	return &Range{
		Start: Index(start),
		Stop:  Index(stop),
		Step:  Index(step),
	}
}

// Make a range iterator from a range
func (r *Range) M__iter__() Object {
	return &RangeIterator{
		Range: *r,
		Index: r.Start,
	}
}

// Range iterator
func (it *RangeIterator) M__iter__() Object {
	return it
}

// Range iterator next
func (it *RangeIterator) M__next__() Object {
	r := it.Index
	if r >= it.Stop {
		panic(StopIteration)
	}
	it.Index += it.Step
	return r
}

// Check interface is satisfied
var _ I__iter__ = (*Range)(nil)
var _ I_iterator = (*RangeIterator)(nil)
