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
		Iterate(v, func(item Object) bool {
			t = append(t, item)
			return false
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
		l.ExtendSequence(v)
		return l
	}
}

// Call __next__ for the python object
//
// Returns the next object
//
// finished == StopIteration or subclass when finished
func Next(self Object) (obj Object, finished Object) {
	defer func() {
		if r := recover(); r != nil {
			if IsException(StopIteration, r) {
				// StopIteration or subclass raised
				finished = r.(Object)
			} else {
				panic(r)
			}
		}
	}()
	if I, ok := self.(I__next__); ok {
		obj = I.M__next__()
		return
	} else if obj, ok = TypeCall0(self, "__next__"); ok {
		return
	}

	panic(ExceptionNewf(TypeError, "'%s' object is not iterable", self.Type().Name))
}

// Create an iterator from obj and iterate the iterator until finished
// calling the function passed in on each object.  The iteration is
// finished if the function returns true
func Iterate(obj Object, fn func(Object) bool) {
	// Some easy cases
	switch x := obj.(type) {
	case Tuple:
		for _, item := range x {
			if fn(item) {
				break
			}
		}
	case *List:
		for _, item := range x.Items {
			if fn(item) {
				break
			}
		}
	case String:
		for _, item := range x {
			if fn(String(item)) {
				break
			}
		}
	case Bytes:
		for _, item := range x {
			if fn(Int(item)) {
				break
			}
		}
	default:
		iterator := Iter(obj)
		for {
			item, finished := Next(iterator)
			if finished != nil {
				break
			}
			if fn(item) {
				break
			}
		}
	}
}

// Call send for the python object
func Send(self, value Object) Object {
	if I, ok := self.(I_send); ok {
		return I.Send(value)
	} else if res, ok := TypeCall1(self, "send", value); ok {
		return res
	}

	panic(ExceptionNewf(TypeError, "'%s' object doesn't have send method", self.Type().Name))
}

// SequenceContains returns True if obj is in seq
func SequenceContains(seq, obj Object) (found bool) {
	Iterate(seq, func(item Object) bool {
		if Eq(item, obj) == True {
			found = true
			return true
		}
		return false
	})
	return
}
