// Internal interface for use from Go
//
// See arithmetic.go for the auto generated stuff

package py

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

	panic(ExceptionNewf(TypeError, "unsupported operand type(s) for index: '%s'", a.Type().Name))
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
		panic(ExceptionNewf(IndexError, "cannot fit %d into an index-sized integer", i))
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
		panic(ExceptionNewf(IndexError, "list index out of range"))
	}
	return i
}

// Returns the number of items of a sequence or mapping
func Len(self Object) Object {
	if I, ok := self.(I__len__); ok {
		return I.M__len__()
	} else if res, ok := TypeCall0(self, "__len__"); ok {
		return res
	}
	panic(ExceptionNewf(TypeError, "object of type '%s' has no len()", self.Type().Name))
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
	panic(ExceptionNewf(TypeError, "'%s' object is not callable", fn.Type().Name))
}

// GetItem
func GetItem(self Object, key Object) Object {
	if I, ok := self.(I__getitem__); ok {
		return I.M__getitem__(key)
	} else if res, ok := TypeCall1(self, "__getitem__", key); ok {
		return res
	}
	panic(ExceptionNewf(TypeError, "'%s' object is not subscriptable", self.Type().Name))
}

// SetItem
func SetItem(self Object, key Object, value Object) Object {
	if I, ok := self.(I__setitem__); ok {
		return I.M__setitem__(key, value)
	} else if res, ok := TypeCall2(self, "__setitem__", key, value); ok {
		return res
	}

	panic(ExceptionNewf(TypeError, "'%s' object does not support item assignment", self.Type().Name))
}

// Delitem
func DelItem(self Object, key Object) Object {
	if I, ok := self.(I__delitem__); ok {
		return I.M__delitem__(key)
	} else if res, ok := TypeCall1(self, "__delitem__", key); ok {
		return res
	}
	panic(ExceptionNewf(TypeError, "'%s' object does not support item deletion", self.Type().Name))
}

// GetAttrErr - returns the result or an err to be raised if not found
//
// Only AttributeErrors will be returned in err, everything else will be raised
func GetAttrErr(self Object, key string) (res Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			if IsException(AttributeError, r) {
				// AttributeError caught - return nil and error
				res = nil
				err = r.(error)
			} else {
				// Propagate the exception
				panic(r)
			}
		}
	}()

	// Call __getattribute__ unconditionally if it exists
	if I, ok := self.(I__getattribute__); ok {
		res = I.M__getattribute__(key)
		return
	} else if res, ok = TypeCall1(self, "__getattribute__", Object(String(key))); ok {
		return
	}

	// Look in the instance dictionary if it exists
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		res, ok = dict[key]
		if ok {
			return
		}
	}

	// Now look in type's dictionary etc
	t := self.Type()
	res = t.NativeGetAttrOrNil(key)
	if res != nil {
		// Call __get__ which creates bound methods, reads properties etc
		if I, ok := res.(I__get__); ok {
			res = I.M__get__(self, t)
		}
		return
	}

	// And now only if not found call __getattr__
	if I, ok := self.(I__getattr__); ok {
		res = I.M__getattr__(key)
		return
	} else if res, ok = TypeCall1(self, "__getattr__", Object(String(key))); ok {
		return
	}

	// Not found - return nil
	res = nil
	err = ExceptionNewf(AttributeError, "'%s' has no attribute '%s'", self.Type().Name, key)
	return
}

// GetAttrString gets the attribute, raising an error if not found
func GetAttrString(self Object, key string) Object {
	res, err := GetAttrErr(self, key)
	if err != nil {
		panic(err)
	}
	return res
}

// GetAttr gets the attribute rasing an error if key isn't a string or
// attribute not found
func GetAttr(self Object, keyObj Object) Object {
	if key, ok := keyObj.(String); ok {
		return GetAttrString(self, string(key))
	}
	panic(ExceptionNewf(TypeError, "attribute name must be string, not '%s'", self.Type().Name))
}

// SetAttrString
func SetAttrString(self Object, key string, value Object) Object {
	// First look in type's dictionary etc for a property that could
	// be set - do this before looking in the instance dictionary
	setter := self.Type().NativeGetAttrOrNil(key)
	if setter != nil {
		// Call __set__ which writes properties etc
		if I, ok := setter.(I__set__); ok {
			return I.M__set__(self, value)
		}
	}

	// If we have __setattr__ then use that
	if I, ok := self.(I__setattr__); ok {
		return I.M__setattr__(key, value)
	} else if res, ok := TypeCall2(self, "__setattr__", String(key), value); ok {
		return res
	}

	// Otherwise set the attribute in the instance dictionary if
	// possible
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		if dict == nil {
			panic(ExceptionNewf(SystemError, "nil Dict in %s", self.Type().Name))
		}
		dict[key] = value
		return None
	}

	// If not blow up
	panic(ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", self.Type().Name, key))
}

// SetAttr
func SetAttr(self Object, keyObj Object, value Object) Object {
	if key, ok := keyObj.(String); ok {
		return GetAttrString(self, string(key))
	}
	panic(ExceptionNewf(TypeError, "attribute name must be string, not '%s'", self.Type().Name))
}

// DeleteAttrString
func DeleteAttrString(self Object, key string) {
	// First look in type's dictionary etc for a property that could
	// be set - do this before looking in the instance dictionary
	deleter := self.Type().NativeGetAttrOrNil(key)
	if deleter != nil {
		// Call __set__ which writes properties etc
		if I, ok := deleter.(I__delete__); ok {
			I.M__delete__(self)
			return
		}
	}

	// If we have __delattr__ then use that
	if I, ok := self.(I__delattr__); ok {
		I.M__delattr__(key)
		return
	} else if _, ok := TypeCall1(self, "__delattr__", String(key)); ok {
		return
	}

	// Otherwise delete the attribute from the instance dictionary
	// if possible
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		if dict == nil {
			panic(ExceptionNewf(SystemError, "nil Dict in %s", self.Type().Name))
		}
		if _, ok := dict[key]; ok {
			delete(dict, key)
			return
		}
	}

	// If not blow up
	panic(ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", self.Type().Name, key))
}

// DeleteAttr
func DeleteAttr(self Object, keyObj Object) {
	if key, ok := keyObj.(String); ok {
		DeleteAttrString(self, string(key))
		return
	}
	panic(ExceptionNewf(TypeError, "attribute name must be string, not '%s'", self.Type().Name))
}
