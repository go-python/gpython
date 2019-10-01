// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Internal interface for use from Go
//
// See arithmetic.go for the auto generated stuff

package py

import (
	"fmt"
	"reflect"
	"strings"
)

// AttributeName converts an Object to a string, raising a TypeError
// if it wasn't a String
func AttributeName(keyObj Object) (string, error) {
	if key, ok := keyObj.(String); ok {
		return string(key), nil
	}
	return "", ExceptionNewf(TypeError, "attribute name must be string, not '%s'", keyObj.Type().Name)
}

// Bool is called to implement truth value testing and the built-in
// operation bool(); should return False or True. When this method is
// not defined, __len__() is called, if it is defined, and the object
// is considered true if its result is nonzero. If a class defines
// neither __len__() nor __bool__(), all its instances are considered
// true.
func MakeBool(a Object) (Object, error) {
	if _, ok := a.(Bool); ok {
		return a, nil
	}

	if A, ok := a.(I__bool__); ok {
		res, err := A.M__bool__()
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	if B, ok := a.(I__len__); ok {
		res, err := B.M__len__()
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return MakeBool(res)
		}
	}

	return True, nil
}

// Turns a into a go int if possible
func MakeGoInt(a Object) (int, error) {
	a, err := MakeInt(a)
	if err != nil {
		return 0, err
	}
	A, ok := a.(IGoInt)
	if ok {
		return A.GoInt()
	}
	return 0, ExceptionNewf(TypeError, "'%v' object cannot be interpreted as a go integer", a.Type().Name)
}

// Turns a into a go int64 if possible
func MakeGoInt64(a Object) (int64, error) {
	a, err := MakeInt(a)
	if err != nil {
		return 0, err
	}
	A, ok := a.(IGoInt64)
	if ok {
		return A.GoInt64()
	}
	return 0, ExceptionNewf(TypeError, "'%v' object cannot be interpreted as a go int64", a.Type().Name)
}

// Index the python Object returning an Int
//
// Will raise TypeError if Index can't be run on this object
func Index(a Object) (Int, error) {
	if A, ok := a.(I__index__); ok {
		return A.M__index__()
	}

	if A, ok, err := TypeCall0(a, "__index__"); ok {
		if err != nil {
			return 0, err
		}

		if res, ok := A.(Int); ok {
			return res, nil
		}

		return 0, ExceptionNewf(TypeError, "__index__ returned non-int: (type %s)", A.Type().Name)
	}

	return 0, ExceptionNewf(TypeError, "unsupported operand type(s) for index: '%s'", a.Type().Name)
}

// Index the python Object returning an int
//
// Will raise TypeError if Index can't be run on this object
//
// or IndexError if the Int won't fit!
func IndexInt(a Object) (int, error) {
	i, err := Index(a)
	if err != nil {
		return 0, err
	}
	intI := int(i)

	// Int might not fit in an int
	if Int(intI) != i {
		return 0, ExceptionNewf(IndexError, "cannot fit %d into an index-sized integer", i)
	}

	return intI, nil
}

// As IndexInt but if index is -ve addresses it from the end
//
// If index is out of range throws IndexError
func IndexIntCheck(a Object, max int) (int, error) {
	i, err := IndexInt(a)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		i += max
	}
	if i < 0 || i >= max {
		return 0, ExceptionNewf(IndexError, "index out of range")
	}
	return i, nil
}

// Returns the number of items of a sequence or mapping
func Len(self Object) (Object, error) {
	if I, ok := self.(I__len__); ok {
		return I.M__len__()
	} else if res, ok, err := TypeCall0(self, "__len__"); ok {
		return res, err
	}
	return nil, ExceptionNewf(TypeError, "object of type '%s' has no len()", self.Type().Name)
}

// Return the result of not a
func Not(a Object) (Object, error) {
	b, err := MakeBool(a)
	if err != nil {
		return nil, err
	}
	switch b {
	case False:
		return True, nil
	case True:
		return False, nil
	}
	return nil, ExceptionNewf(TypeError, "bool() didn't return True or False")
}

// Calls function fnObj with args and kwargs in a new vm (or directly
// if Go code)
//
// kwargs should be nil if not required
//
// fnObj must be a callable type such as *py.Method or *py.Function
//
// The result is returned
func Call(fn Object, args Tuple, kwargs StringDict) (Object, error) {
	if I, ok := fn.(I__call__); ok {
		return I.M__call__(args, kwargs)
	}
	return nil, ExceptionNewf(TypeError, "'%s' object is not callable", fn.Type().Name)
}

// GetItem
func GetItem(self Object, key Object) (Object, error) {
	if I, ok := self.(I__getitem__); ok {
		return I.M__getitem__(key)
	} else if res, ok, err := TypeCall1(self, "__getitem__", key); ok {
		return res, err
	}
	return nil, ExceptionNewf(TypeError, "'%s' object is not subscriptable", self.Type().Name)
}

// SetItem
func SetItem(self Object, key Object, value Object) (Object, error) {
	if I, ok := self.(I__setitem__); ok {
		return I.M__setitem__(key, value)
	} else if res, ok, err := TypeCall2(self, "__setitem__", key, value); ok {
		return res, err
	}

	return nil, ExceptionNewf(TypeError, "'%s' object does not support item assignment", self.Type().Name)
}

// Delitem
func DelItem(self Object, key Object) (Object, error) {
	if I, ok := self.(I__delitem__); ok {
		return I.M__delitem__(key)
	} else if res, ok, err := TypeCall1(self, "__delitem__", key); ok {
		return res, err
	}
	return nil, ExceptionNewf(TypeError, "'%s' object does not support item deletion", self.Type().Name)
}

// GetAttrString - returns the result or an err to be raised if not found
//
// If not found err will be an AttributeError
func GetAttrString(self Object, key string) (res Object, err error) {
	// Call __getattribute__ unconditionally if it exists
	if I, ok := self.(I__getattribute__); ok {
		return I.M__getattribute__(key)
	} else if res, ok, err = TypeCall1(self, "__getattribute__", Object(String(key))); ok {
		return res, err
	}

	// Look up any __special__ methods as M__special__ and return a bound method
	if len(key) >= 5 && strings.HasPrefix(key, "__") && strings.HasSuffix(key, "__") {
		objectValue := reflect.ValueOf(self)
		methodValue := objectValue.MethodByName("M" + key)
		if methodValue.IsValid() {
			return newBoundMethod(key, methodValue.Interface())
		}
	}

	// Look in the instance dictionary if it exists
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		res, ok = dict[key]
		if ok {
			return res, err
		}
	}

	// Now look in type's dictionary etc
	t := self.Type()
	res = t.NativeGetAttrOrNil(key)
	if res != nil {
		// Call __get__ which creates bound methods, reads properties etc
		if I, ok := res.(I__get__); ok {
			res, err = I.M__get__(self, t)
		}
		return res, err
	}

	// And now only if not found call __getattr__
	if I, ok := self.(I__getattr__); ok {
		return I.M__getattr__(key)
	} else if res, ok, err = TypeCall1(self, "__getattr__", Object(String(key))); ok {
		return res, err
	}

	// Not found - return nil
	return nil, ExceptionNewf(AttributeError, "'%s' has no attribute '%s'", self.Type().Name, key)
}

// GetAttrErr - returns the result or an err to be raised if not found
//
// If not found an AttributeError will be returned
func GetAttr(self Object, keyObj Object) (res Object, err error) {
	key, err := AttributeName(keyObj)
	if err != nil {
		return nil, err
	}
	return GetAttrString(self, key)
}

// SetAttrString
func SetAttrString(self Object, key string, value Object) (Object, error) {
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
	} else if res, ok, err := TypeCall2(self, "__setattr__", String(key), value); ok {
		return res, err
	}

	// Otherwise set the attribute in the instance dictionary if
	// possible
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		if dict == nil {
			return nil, ExceptionNewf(SystemError, "nil Dict in %s", self.Type().Name)
		}
		dict[key] = value
		return None, nil
	}

	// If not blow up
	return nil, ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", self.Type().Name, key)
}

// SetAttr
func SetAttr(self Object, keyObj Object, value Object) (Object, error) {
	key, err := AttributeName(keyObj)
	if err != nil {
		return nil, err
	}
	return SetAttrString(self, key, value)
}

// DeleteAttrString
func DeleteAttrString(self Object, key string) error {
	// First look in type's dictionary etc for a property that could
	// be set - do this before looking in the instance dictionary
	deleter := self.Type().NativeGetAttrOrNil(key)
	if deleter != nil {
		// Call __set__ which writes properties etc
		if I, ok := deleter.(I__delete__); ok {
			_, err := I.M__delete__(self)
			return err
		}
	}

	// If we have __delattr__ then use that
	if I, ok := self.(I__delattr__); ok {
		_, err := I.M__delattr__(key)
		return err
	} else if _, ok, err := TypeCall1(self, "__delattr__", String(key)); ok {
		return err
	}

	// Otherwise delete the attribute from the instance dictionary
	// if possible
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		if dict == nil {
			return ExceptionNewf(SystemError, "nil Dict in %s", self.Type().Name)
		}
		if _, ok := dict[key]; ok {
			delete(dict, key)
			return nil
		}
	}

	// If not blow up
	return ExceptionNewf(AttributeError, "'%s' object has no attribute '%s'", self.Type().Name, key)
}

// DeleteAttr
func DeleteAttr(self Object, keyObj Object) error {
	key, err := AttributeName(keyObj)
	if err != nil {
		return err
	}
	return DeleteAttrString(self, key)
}

// Calls __str__ on the object
//
// Calls __repr__ on the object or returns a sensible default
func Repr(self Object) (Object, error) {
	if I, ok := self.(I__repr__); ok {
		return I.M__repr__()
	} else if res, ok, err := TypeCall0(self, "__repr__"); ok {
		return res, err
	}
	return String(fmt.Sprintf("<%s instance at %p>", self.Type().Name, self)), nil
}

// DebugRepr - see Repr but returns the repr or error as a string
func DebugRepr(self Object) string {
	res, err := Repr(self)
	if err != nil {
		return fmt.Sprintf("Repr(%s) returned %v", self.Type().Name, err)
	}
	str, ok := res.(String)
	if !ok {
		return fmt.Sprintf("Repr(%s) didn't return a string", self.Type().Name)
	}
	return string(str)
}

// Calls __str__ on the object and if not found calls __repr__
func Str(self Object) (Object, error) {
	if I, ok := self.(I__str__); ok {
		return I.M__str__()
	} else if res, ok, err := TypeCall0(self, "__str__"); ok {
		return res, err
	}
	return Repr(self)
}

// Returns object as a string
//
// Calls Str then makes sure the output is a string
func StrAsString(self Object) (string, error) {
	res, err := Str(self)
	if err != nil {
		return "", err
	}
	str, ok := res.(String)
	if !ok {
		return "", ExceptionNewf(TypeError, "result of __str__ must be string, not '%s'", res.Type().Name)
	}
	return string(str), nil
}

// Returns object as a string
//
// Calls Repr then makes sure the output is a string
func ReprAsString(self Object) (string, error) {
	res, err := Repr(self)
	if err != nil {
		return "", err
	}
	str, ok := res.(String)
	if !ok {
		return "", ExceptionNewf(TypeError, "result of __repr__ must be string, not '%s'", res.Type().Name)
	}
	return string(str), nil
}
