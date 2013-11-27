// Tuple objects

package py

var TupleType = NewType("tuple", "tuple() -> empty tuple\ntuple(iterable) -> tuple initialized from iterable's items\n\nIf the argument is a tuple, the return value is the same object.")

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
}

// Copy a tuple object
func (t Tuple) Copy() Tuple {
	newT := make(Tuple, len(t))
	for i := range t {
		newT[i] = t[i]
	}
	return newT
}
func (t Tuple) M__len__() Object {
	return Int(len(t))
}

func (t Tuple) M__bool__() Object {
	if len(t) > 0 {
		return True
	}
	return False
}

func (t Tuple) M__iter__() Object {
	return NewIterator(t)
}

func (t Tuple) M__getitem__(key Object) Object {
	i := IndexIntCheck(key, len(t))
	return t[i]
}

// Check interface is satisfied
var _ I__len__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__iter__ = Tuple(nil)
var _ I__getitem__ = Tuple(nil)

// var _ richComparison = Tuple(nil)
