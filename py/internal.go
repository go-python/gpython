// Internal interface for use from Go
//
// See arithmetic.go for the auto generated stuff

package py

import (
	"fmt"
)

// Bool is called to implement truth value testing and the built-in
// operation bool(); should return False or True. When this method is
// not defined, __len__() is called, if it is defined, and the object
// is considered true if its result is nonzero. If a class defines
// neither __len__() nor __bool__(), all its instances are considered
// true.
func MakeBool(a Object) Object {
	if _, ok := a.(Bool); ok {
		return a
	}

	if A, ok := a.(I__bool__); ok {
		res := A.M__bool__()
		if res != NotImplemented {
			return res
		}
	}

	if B, ok := a.(I__len__); ok {
		res := B.M__len__()
		if res != NotImplemented {
			return MakeBool(res)
		}
	}

	return True
}

// Index the python Object returning an Int
//
// Will raise TypeError if Index can't be run on this object
func Index(a Object) Int {
	A, ok := a.(I__index__)
	if ok {
		return A.M__index__()
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for index: '%s'", a.Type().Name))
}

// Index the python Object returning an int
//
// Will raise TypeError if Index can't be run on this object
//
// or IndexError if the Int won't fit!
func IndexInt(a Object) int {
	i := Index(a)
	intI := int(i)

	// Int might not fit in an int
	if Int(intI) != i {
		// FIXME IndexError
		panic(fmt.Sprintf("IndexError: cannot fit %d into an index-sized integer", i))
	}

	return intI
}

// As IndexInt but if index is -ve addresses it from the end
//
// If index is out of range throws IndexError
func IndexIntCheck(a Object, max int) int {
	i := IndexInt(a)
	if i < 0 {
		i += max
	}
	if i < 0 || i >= max {
		// FIXME IndexError
		panic("IndexError: list index out of range")
	}
	return i
}

// Return the result of not a
func Not(a Object) Object {
	switch MakeBool(a) {
	case False:
		return True
	case True:
		return False
	}
	panic("bool() didn't return True or False")
}

// Calls function fnObj with args and kwargs in a new vm (or directly
// if Go code)
//
// kwargs should be nil if not required
//
// fnObj must be a callable type such as *py.Method or *py.Function
//
// The result is returned
func Call(fn Object, args Tuple, kwargs StringDict) Object {
	if I, ok := fn.(I__call__); ok {
		return I.M__call__(args, kwargs)
	}
	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: '%s' object is not callable", fn.Type().Name))
}

// GetItem
func GetItem(self Object, key Object) Object {
	if I, ok := self.(I__getitem__); ok {
		return I.M__getitem__(key)
	} else if res, ok := TypeCall1(self, "__getitem__", key); ok {
		return res
	}
	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: '%s' object is not subscriptable", self.Type().Name))
}

// SetItem
func SetItem(self Object, key Object, value Object) Object {
	if I, ok := self.(I__setitem__); ok {
		return I.M__setitem__(key, value)
	} else if res, ok := TypeCall2(self, "__setitem__", key, value); ok {
		return res
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: '%s' object does not support item assignment", self.Type().Name))
}

// GetAttrOrNil - returns the result nil if attribute not found
func GetAttrOrNil(self Object, key string) (res Object) {
	// Call __getattribute unconditionally if it exists
	if I, ok := self.(I__getattribute__); ok {
		res = I.M__getattribute__(key)
		goto found
	} else if res, ok = TypeCall1(self, "__getattribute__", Object(String(key))); ok {
		// FIXME catch AttributeError here
		goto found
	}

	if t, ok := self.(*Type); ok {
		// Now look in the instance dictionary etc
		res = t.GetAttrOrNil(key)
		if res != nil {
			goto found
		}
	} else {
		// FIXME introspection for M__methods__ on non *Type objects
	}

	// And now only if not found call __getattr__
	if I, ok := self.(I__getattr__); ok {
		res = I.M__getattr__(key)
		goto found
	} else if res, ok = TypeCall1(self, "__getattr__", Object(String(key))); ok {
		goto found
	}

	// Not found - return nil
	res = nil
	return

found:
	// FIXME if self is an instance then if it returning a function then it needs to return a bound method?
	// otherwise it should return a function
	//
	// >>> str.find
	// <method 'find' of 'str' objects>
	// >>> "".find
	// <built-in method find of str object at 0x7f929bd54c00>
	// >>>
	//
	// created by PyMethod_New defined in classobject.c
	// called by type.tp_descr_get

	// FIXME Not completely correct!
	// Should be using __get__
	switch res.(type) {
	case *Function, *Method:
		res = NewBoundMethod(self, res)
	}
	return
}

// GetAttrString
func GetAttrString(self Object, key string) Object {
	res := GetAttrOrNil(self, key)
	if res == nil {
		// FIXME should be AttributeError
		panic(fmt.Sprintf("AttributeError: '%s' has no attribute '%s'", self.Type().Name, key))
	}
	return res
}

// GetAttr
func GetAttr(self Object, keyObj Object) Object {
	if key, ok := keyObj.(String); ok {
		return GetAttrString(self, string(key))
	}
	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: attribute name must be string, not '%s'", self.Type().Name))
}

// SetAttrString
func SetAttrString(self Object, key string, value Object) Object {
	if I, ok := self.(I__setattr__); ok {
		return I.M__setattr__(key, value)
	} else if res, ok := TypeCall2(self, "__setattr__", String(key), value); ok {
		return res
	}

	// Set the attribute on *Type
	if t, ok := self.(*Type); ok {
		if t.Dict == nil {
			t.Dict = make(StringDict)
		}
		t.Dict[key] = value
		return None
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: '%s' object does not support setting attributes", self.Type().Name))
}

// SetAttr
func SetAttr(self Object, keyObj Object, value Object) Object {
	if key, ok := keyObj.(String); ok {
		return GetAttrString(self, string(key))
	}
	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: attribute name must be string, not '%s'", self.Type().Name))
}

// Call __next__ for the python object
func Next(self Object) Object {
	if I, ok := self.(I__next__); ok {
		return I.M__next__()
	} else if res, ok := TypeCall0(self, "__next__"); ok {
		return res
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: '%s' object is not iterable", self.Type().Name))
}
