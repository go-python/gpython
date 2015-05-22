// Tuple objects

package py

var TupleType = ObjectType.NewType("tuple", "tuple() -> empty tuple\ntuple(iterable) -> tuple initialized from iterable's items\n\nIf the argument is a tuple, the return value is the same object.", TupleNew, nil)

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
}

// TupleNew
func TupleNew(metatype *Type, args Tuple, kwargs StringDict) (res Object) {
	var iterable Object
	UnpackTuple(args, kwargs, "tuple", 0, 1, &iterable)
	if iterable != nil {
		return SequenceTuple(iterable)
	}
	return Tuple{}
}

// Copy a tuple object
func (t Tuple) Copy() Tuple {
	newT := make(Tuple, len(t))
	copy(newT, t)
	return newT
}

// Reverses a tuple (in-place)
func (t Tuple) Reverse() {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
}

func (t Tuple) M__len__() Object {
	return Int(len(t))
}

func (t Tuple) M__bool__() Object {
	return NewBool(len(t) > 0)
}

func (t Tuple) M__iter__() Object {
	return NewIterator(t)
}

func (t Tuple) M__getitem__(key Object) Object {
	if slice, ok := key.(*Slice); ok {
		start, stop, step, slicelength := slice.GetIndices(len(t))
		if step == 1 {
			// Return a subslice since tuples are immutable
			return t[start:stop]
		}
		newTuple := make(Tuple, slicelength)
		for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
			newTuple[j] = t[i]
		}
		return newTuple
	}
	i := IndexIntCheck(key, len(t))
	return t[i]
}

func (a Tuple) M__add__(other Object) Object {
	if b, ok := other.(Tuple); ok {
		newTuple := make(Tuple, len(a)+len(b))
		copy(newTuple, a)
		copy(newTuple[len(b):], b)
		return newTuple
	}
	return NotImplemented
}

func (a Tuple) M__radd__(other Object) Object {
	if b, ok := other.(Tuple); ok {
		return b.M__add__(a)
	}
	return NotImplemented
}

func (a Tuple) M__iadd__(other Object) Object {
	return a.M__add__(other)
}

func (l Tuple) M__mul__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		m := len(l)
		n := int(b) * m
		newTuple := make(Tuple, n)
		for i := 0; i < n; i += m {
			copy(newTuple[i:i+m], l)
		}
		return newTuple
	}
	return NotImplemented
}

func (a Tuple) M__rmul__(other Object) Object {
	return a.M__mul__(other)
}

func (a Tuple) M__imul__(other Object) Object {
	return a.M__mul__(other)
}

// Check interface is satisfied
var _ sequenceArithmetic = Tuple(nil)
var _ I__len__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__iter__ = Tuple(nil)
var _ I__getitem__ = Tuple(nil)

// var _ richComparison = Tuple(nil)

func (a Tuple) M__eq__(other Object) Object {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented
	}
	if len(a) != len(b) {
		return False
	}
	for i := range a {
		if Eq(a[i], b[i]) == False {
			return False
		}
	}
	return True
}

func (a Tuple) M__ne__(other Object) Object {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented
	}
	if len(a) != len(b) {
		return True
	}
	for i := range a {
		if Eq(a[i], b[i]) == False {
			return True
		}
	}
	return False
}
