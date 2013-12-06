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
	t := Tuple{}
	defer func() {
		if r := recover(); r != nil {
			if IsException(StopIteration, r) {
				// StopIteration or subclass raised
				res = t
			} else {
				panic(r)
			}
		}
	}()
	var iterable Object
	UnpackTuple(args, kwargs, "tuple", 0, 1, &iterable)
	if iterable == nil {
		return t
	}
	it := Iter(iterable)
	for {
		item := Next(it)
		t = append(t, item)
	}
}

// Copy a tuple object
func (t Tuple) Copy() Tuple {
	newT := make(Tuple, len(t))
	copy(newT, t)
	return newT
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
	i := IndexIntCheck(key, len(t))
	return t[i]
}

// Check interface is satisfied
var _ I__len__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__iter__ = Tuple(nil)
var _ I__getitem__ = Tuple(nil)

// var _ richComparison = Tuple(nil)
