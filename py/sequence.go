// Sequence operations

package py

// Converts a sequence object v into a Tuple
func SequenceTuple(v Object) Tuple {
	// FIXME need to support iterable objects etc!
	switch x := v.(type) {
	case Tuple:
		return x
	case List:
		return Tuple(x.Copy())
	}
	panic("SequenceTuple not fully implemented")
}

// Converts a sequence object v into a List
func SequenceList(v Object) List {
	// FIXME need to support iterable objects etc!
	switch x := v.(type) {
	case Tuple:
		return List(x).Copy()
	case List:
		return x.Copy()
	}
	panic("SequenceList not fully implemented")
}
