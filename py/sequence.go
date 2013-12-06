// Sequence operations

package py

// Converts a sequence object v into a Tuple
func SequenceTuple(v Object) Tuple {
	switch x := v.(type) {
	case Tuple:
		return x
	case *List:
		return Tuple(x.Items).Copy()
	default:
		t := Tuple{}
		Iterate(Iter(v), func(item Object) {
			t = append(t, item)
		})
		return t
	}
}

// Converts a sequence object v into a List
func SequenceList(v Object) *List {
	switch x := v.(type) {
	case Tuple:
		return NewListFromItems(x)
	case *List:
		return x.Copy()
	default:
		l := NewList()
		Iterate(Iter(v), func(item Object) {
			l.Append(item)
		})
		return l
	}
}
