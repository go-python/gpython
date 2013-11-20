// Tuple objects

package py

var TupleType = NewType("tuple")

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
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

// Check interface is satisfied
var _ I__len__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__iter__ = Tuple(nil)

// var _ richComparison = Tuple(nil)
