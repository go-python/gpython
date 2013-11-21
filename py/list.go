// List objects

package py

var ListType = NewType("list", "list() -> new empty list\nlist(iterable) -> new list initialized from iterable's items")

type List []Object

// Type of this List object
func (o List) Type() *Type {
	return ListType
}

func (t List) M__len__() Object {
	return Int(len(t))
}

func (t List) M__bool__() Object {
	if len(t) > 0 {
		return True
	}
	return False
}

func (t List) M__iter__() Object {
	return NewIterator(t)
}

// Check interface is satisfied
var _ I__len__ = List(nil)
var _ I__bool__ = List(nil)
var _ I__iter__ = List(nil)

// var _ richComparison = List(nil)
