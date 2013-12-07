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
		Iterate(v, func(item Object) {
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
		Iterate(v, func(item Object) {
			l.Append(item)
		})
		return l
	}
}

// Call __next__ for the python object
func Next(self Object) Object {
	if I, ok := self.(I__next__); ok {
		return I.M__next__()
	} else if res, ok := TypeCall0(self, "__next__"); ok {
		return res
	}

	panic(ExceptionNewf(TypeError, "'%s' object is not iterable", self.Type().Name))
}

// Create an iterator from obj and iterate the iterator until finished
// calling the function passed in on each object
func Iterate(obj Object, fn func(Object)) {
	defer func() {
		if r := recover(); r != nil {
			if IsException(StopIteration, r) {
				// StopIteration or subclass raised
			} else {
				panic(r)
			}
		}
	}()
	// Some easy cases
	switch x := obj.(type) {
	case Tuple:
		for _, item := range x {
			fn(item)
		}
	case *List:
		for _, item := range x.Items {
			fn(item)
		}
	case String:
		for _, item := range x {
			fn(String(item))
		}
	case Bytes:
		for _, item := range x {
			fn(Int(item))
		}
	default:
		iterator := Iter(obj)
		for {
			item := Next(iterator)
			fn(item)
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
