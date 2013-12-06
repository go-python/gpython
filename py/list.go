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
	i := IndexIntCheck(key, len(l.Items))
	return l.Items[i]
}

func (l *List) M__setitem__(key, value Object) Object {
	i := IndexIntCheck(key, len(l.Items))
	l.Items[i] = value
	return None
}

// Check interface is satisfied
var _ I__len__ = (*List)(nil)
var _ I__bool__ = (*List)(nil)
var _ I__iter__ = (*List)(nil)
var _ I__getitem__ = (*List)(nil)
var _ I__setitem__ = (*List)(nil)

// var _ richComparison = (*List)(nil)
