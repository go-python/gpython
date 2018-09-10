// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Method objects
//
// This is about the type 'builtin_function_or_method', not Python
// methods in user-defined classes.  See class.go for the latter.

package py

import (
	"fmt"
)

// Types for methods

// Called with self and a tuple of args
type PyCFunction func(self Object, args Tuple) (Object, error)

// Called with self, a tuple of args and a stringdic of kwargs
type PyCFunctionWithKeywords func(self Object, args Tuple, kwargs StringDict) (Object, error)

// Called with self only
type PyCFunctionNoArgs func(Object) (Object, error)

// Called with one (unnamed) parameter only
type PyCFunction1Arg func(Object, Object) (Object, error)

const (
	// These two constants are not used to indicate the calling convention
	// but the binding when use with methods of classes. These may not be
	// used for functions defined for modules. At most one of these flags
	// may be set for any given method.

	// The method will be passed the type object as the first parameter
	// rather than an instance of the type. This is used to create class
	// methods, similar to what is created when using the classmethod()
	// built-in function.
	METH_CLASS = 0x0010

	// The method will be passed NULL as the first parameter rather than
	// an instance of the type. This is used to create static methods,
	// similar to what is created when using the staticmethod() built-in
	// function.
	METH_STATIC = 0x0020

	// One other constant controls whether a method is loaded in
	// place of another definition with the same method name.

	// The method will be loaded in place of existing definitions. Without
	// METH_COEXIST, the default is to skip repeated definitions. Since
	// slot wrappers are loaded before the method table, the existence of
	// a sq_contains slot, for example, would generate a wrapped method
	// named __contains__() and preclude the loading of a corresponding
	// PyCFunction with the same name. With the flag defined, the
	// PyCFunction will be loaded in place of the wrapper object and will
	// co-exist with the slot. This is helpful because calls to
	// PyCFunctions are optimized more than wrapper object calls.
	METH_COEXIST = 0x0040
)

// A python Method object
type Method struct {
	// Name of this function
	Name string
	// Doc string
	Doc string
	// Flags - see METH_* flags
	Flags int
	// Go function implementation
	method interface{}
}

// Internal method types implemented within eval.go
type InternalMethod int

const (
	InternalMethodNone InternalMethod = iota
	InternalMethodGlobals
	InternalMethodLocals
	InternalMethodImport
	InternalMethodEval
	InternalMethodExec
)

var MethodType = NewType("method", "method object")

// Type of this object
func (o *Method) Type() *Type {
	return MethodType
}

// Define a new method
func NewMethod(name string, method interface{}, flags int, doc string) (*Method, error) {
	// have to write out the function arguments - can't use the
	// type aliases as they are different types :-(
	switch method.(type) {
	case func(self Object, args Tuple) (Object, error):
	case func(self Object, args Tuple, kwargs StringDict) (Object, error):
	case func(Object) (Object, error):
	case func(Object, Object) (Object, error):
	case InternalMethod:
	default:
		return nil, ExceptionNewf(SystemError, "Unknown function type for NewMethod %q, %T", name, method)
	}
	return &Method{
		Name:   name,
		Doc:    doc,
		Flags:  flags,
		method: method,
	}, nil
}

// As NewMethod but panics on error
func MustNewMethod(name string, method interface{}, flags int, doc string) *Method {
	m, err := NewMethod(name, method, flags, doc)
	if err != nil {
		panic(err)
	}
	return m
}

// Returns the InternalMethod type of this method
func (m *Method) Internal() InternalMethod {
	if internalMethod, ok := m.method.(InternalMethod); ok {
		return internalMethod
	}
	return InternalMethodNone
}

// Call the method with the given arguments
func (m *Method) Call(self Object, args Tuple) (Object, error) {
	switch f := m.method.(type) {
	case func(self Object, args Tuple) (Object, error):
		return f(self, args)
	case func(self Object, args Tuple, kwargs StringDict) (Object, error):
		return f(self, args, NewStringDict())
	case func(Object) (Object, error):
		if len(args) != 0 {
			return nil, ExceptionNewf(TypeError, "%s() takes no arguments (%d given)", m.Name, len(args))
		}
		return f(self)
	case func(Object, Object) (Object, error):
		if len(args) != 1 {
			return nil, ExceptionNewf(TypeError, "%s() takes exactly 1 argument (%d given)", m.Name, len(args))
		}
		return f(self, args[0])
	}
	panic(fmt.Sprintf("Unknown method type: %T", m.method))
}

// Call the method with the given arguments
func (m *Method) CallWithKeywords(self Object, args Tuple, kwargs StringDict) (Object, error) {
	if len(kwargs) == 0 {
		return m.Call(self, args)
	}
	switch f := m.method.(type) {
	case func(self Object, args Tuple, kwargs StringDict) (Object, error):
		return f(self, args, kwargs)
	case func(self Object, args Tuple) (Object, error),
		func(Object) (Object, error),
		func(Object, Object) (Object, error):
		return nil, ExceptionNewf(TypeError, "%s() takes no keyword arguments", m.Name)
	}
	panic(fmt.Sprintf("Unknown method type: %T", m.method))
}

// Return a new Method with the bound method passed in, or an error
//
// This needs to convert the methods into internally callable python
// methods
func newBoundMethod(name string, fn interface{}) (Object, error) {
	m := &Method{
		Name: name,
	}
	switch f := fn.(type) {
	case func(args Tuple) (Object, error):
		m.method = func(_ Object, args Tuple) (Object, error) {
			return f(args)
		}
	// M__call__(args Tuple, kwargs StringDict) (Object, error)
	case func(args Tuple, kwargs StringDict) (Object, error):
		m.method = func(_ Object, args Tuple, kwargs StringDict) (Object, error) {
			return f(args, kwargs)
		}
	// M__str__() (Object, error)
	case func() (Object, error):
		m.method = func(_ Object) (Object, error) {
			return f()
		}
	// M__add__(other Object) (Object, error)
	case func(Object) (Object, error):
		m.method = func(_ Object, other Object) (Object, error) {
			return f(other)
		}
	// M__getattr__(name string) (Object, error)
	case func(string) (Object, error):
		m.method = func(_ Object, stringObject Object) (Object, error) {
			name, err := StrAsString(stringObject)
			if err != nil {
				return nil, err
			}
			return f(name)
		}
	// M__get__(instance, owner Object) (Object, error)
	case func(Object, Object) (Object, error):
		m.method = func(_ Object, args Tuple) (Object, error) {
			var a, b Object
			err := UnpackTuple(args, nil, name, 2, 2, &a, &b)
			if err != nil {
				return nil, err
			}
			return f(a, b)
		}
	// M__new__(cls, args, kwargs Object) (Object, error)
	case func(Object, Object, Object) (Object, error):
		m.method = func(_ Object, args Tuple) (Object, error) {
			var a, b, c Object
			err := UnpackTuple(args, nil, name, 3, 3, &a, &b, &c)
			if err != nil {
				return nil, err
			}
			return f(a, b, c)
		}
	default:
		return nil, fmt.Errorf("Unknown bound method type for %q: %T", name, fn)
	}
	return m, nil
}

// Call a method
func (m *Method) M__call__(args Tuple, kwargs StringDict) (Object, error) {
	self := None // FIXME should be the module
	if kwargs != nil {
		return m.CallWithKeywords(self, args, kwargs)
	}
	return m.Call(self, args)
}

// Read a method from a class which makes a bound method
func (m *Method) M__get__(instance, owner Object) (Object, error) {
	if instance != None {
		return NewBoundMethod(instance, m), nil
	}
	return m, nil
}

// FIXME this should be the default?
func (m *Method) M__eq__(other Object) (Object, error) {
	if otherMethod, ok := other.(*Method); ok && m == otherMethod {
		return True, nil
	}
	return False, nil
}

// FIXME this should be the default?
func (m *Method) M__ne__(other Object) (Object, error) {
	if otherMethod, ok := other.(*Method); ok && m == otherMethod {
		return False, nil
	}
	return True, nil
}

// Make sure it satisfies the interface
var _ Object = (*Method)(nil)
var _ I__call__ = (*Method)(nil)
var _ I__get__ = (*Method)(nil)
var _ I__eq__ = (*Method)(nil)
var _ I__ne__ = (*Method)(nil)
