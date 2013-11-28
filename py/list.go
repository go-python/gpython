// List objects

package py

var ListType = NewType("list", "list() -> new empty list\nlist(iterable) -> new list initialized from iterable's items")

// FIXME lists are mutable so this should probably be struct { Tuple } then can use the sub methods on Tuple
type List []Object

// Type of this List object
func (o List) Type() *Type {
	return ListType
}

// Copy a list object
func (l List) Copy() List {
	newL := make(List, len(l))
	for i := range l {
		newL[i] = l[i]
	}
	return newL
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

func (t List) M__getitem__(key Object) Object {
	i := IndexIntCheck(key, len(t))
	return t[i]
}

func (t List) M__setitem__(key, value Object) Object {
	i := IndexIntCheck(key, len(t))
	t[i] = value
	return None
}

// Check interface is satisfied
var _ I__len__ = List(nil)
var _ I__bool__ = List(nil)
var _ I__iter__ = List(nil)
var _ I__getitem__ = List(nil)
var _ I__setitem__ = List(nil)

// var _ richComparison = List(nil)
