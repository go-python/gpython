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

func (it *Iterator) M__iter__() Object {
	return it
}

// Get next one from the iteration
func (it *Iterator) M__next__() Object {
	if it.Pos >= len(it.Objs) {
		panic(StopIteration)
	}
	r := it.Objs[it.Pos]
	it.Pos++
	return r
}

// Check interface is satisfied
var _ I_iterator = (*Iterator)(nil)
