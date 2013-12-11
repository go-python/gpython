// List objects

package py

var ListType = ObjectType.NewType("list", "list() -> new empty list\nlist(iterable) -> new list initialized from iterable's items", ListNew, nil)

// FIXME lists are mutable so this should probably be struct { Tuple } then can use the sub methods on Tuple
type List struct {
	Items []Object
}

// Type of this List object
func (o *List) Type() *Type {
	return ListType
}

// ListNew
func ListNew(metatype *Type, args Tuple, kwargs StringDict) (res Object) {
	var iterable Object
	UnpackTuple(args, kwargs, "list", 0, 1, &iterable)
	if iterable != nil {
		return SequenceList(iterable)
	}
	return NewList()
}

// Make a new empty list
func NewList() *List {
	return &List{}
}

// Make a new empty list with given capacity
func NewListWithCapacity(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, 0, n)
	}
	return l
}

// Make a list with n nil elements
func NewListSized(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, n)
	}
	return l
}

// Make a new list from an []Object
//
// The []Object is copied into the list
func NewListFromItems(items []Object) *List {
	l := NewListSized(len(items))
	copy(l.Items, items)
	return l
}

// Copy a list object
func (l *List) Copy() *List {
	return NewListFromItems(l.Items)
}

// Append an item
func (l *List) Append(item Object) {
	l.Items = append(l.Items, item)
}

// Extend the list with items
func (l *List) Extend(items []Object) {
	l.Items = append(l.Items, items...)
}

func (l *List) M__len__() Object {
	return Int(len(l.Items))
}

func (l *List) M__bool__() Object {
	return NewBool(len(l.Items) > 0)
}

func (l *List) M__iter__() Object {
	return NewIterator(l.Items)
}

func (l *List) M__getitem__(key Object) Object {
	if slice, ok := key.(*Slice); ok {
		start, _, step, slicelength := slice.GetIndices(len(l.Items))
		newList := NewListSized(slicelength)
		for i, j := start, 0; j < slicelength; i, j = i+step, j+1 {
			newList.Items[j] = l.Items[i]
		}
		return newList
	}
	i := IndexIntCheck(key, len(l.Items))
	return l.Items[i]
}

func (l *List) M__setitem__(key, value Object) Object {
	if slice, ok := key.(*Slice); ok {
		start, stop, step, slicelength := slice.GetIndices(len(l.Items))
		if step != 1 {
			panic("Setting slices with step != 1 not implemented yet")
		}
		if stop == len(l.Items) {
			// tail of the list only
			l.Items = l.Items[:start]
			Iterate(value, func(item Object) {
				l.Append(item)
			})
			return None
		}
		_ = slicelength
		_ = start
		_ = stop
		panic("Set slice not implemented fully yet")
		return None
	}
	i := IndexIntCheck(key, len(l.Items))
	l.Items[i] = value
	return None
}

func (a *List) M__add__(other Object) Object {
	if b, ok := other.(*List); ok {
		newList := NewListSized(len(a.Items) + len(b.Items))
		copy(newList.Items, a.Items)
		copy(newList.Items[len(b.Items):], b.Items)
		return newList
	}
	return NotImplemented
}

func (a *List) M__radd__(other Object) Object {
	if b, ok := other.(*List); ok {
		return b.M__add__(a)
	}
	return NotImplemented
}

func (a *List) M__iadd__(other Object) Object {
	if b, ok := other.(*List); ok {
		a.Extend(b.Items)
		return a
	}
	return NotImplemented
}

func (l *List) M__mul__(other Object) Object {
	if b, ok := convertToInt(other); ok {
		m := len(l.Items)
		n := int(b) * m
		newList := NewListSized(n)
		for i := 0; i < n; i += m {
			copy(newList.Items[i:i+m], l.Items)
		}
		return newList
	}
	return NotImplemented
}

func (a *List) M__rmul__(other Object) Object {
	return a.M__mul__(other)
}

func (a *List) M__imul__(other Object) Object {
	return a.M__mul__(other)
}

// Check interface is satisfied
var _ sequenceArithmetic = (*List)(nil)
var _ I__len__ = (*List)(nil)
var _ I__len__ = (*List)(nil)
var _ I__bool__ = (*List)(nil)
var _ I__iter__ = (*List)(nil)
var _ I__getitem__ = (*List)(nil)
var _ I__setitem__ = (*List)(nil)

// var _ richComparison = (*List)(nil)
