// Sequence operations

package py

// Converts a sequence object v into a Tuple
func SequenceTuple(v Object) Tuple {
	// FIXME need to support iterable objects etc!
	switch x := v.(type) {
	case Tuple:
		return x
	case *List:
		return Tuple(x.Items).Copy()
	}
	panic(ExceptionNewf(TypeError, "SequenceTuple not fully implemented, can't convert %s", v.Type().Name))
}

// Converts a sequence object v into a List
func SequenceList(v Object) *List {
	// FIXME need to support iterable objects etc!
	switch x := v.(type) {
	case Tuple:
		return NewListFromItems(x)
	case *List:
		return x.Copy()
	}
	panic(ExceptionNewf(TypeError, "SequenceList not fully implemented, can't convert %s", v.Type().Name))
}
