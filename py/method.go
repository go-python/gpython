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
type PyCFunction func(self Object, args Tuple) Object

// Called with self, a tuple of args and a stringdic of kwargs
type PyCFunctionWithKeywords func(self Object, args Tuple, kwargs StringDict) Object

// Called with self only
type PyCFunctionNoArgs func(Object) Object

// Called with one (unnamed) parameter only
type PyCFunction1Arg func(Object, Object) Object

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

var MethodType = NewType("method", "method object")

// Type of this object
func (o *Method) Type() *Type {
	return MethodType
}

// Define a new method
func NewMethod(name string, method interface{}, flags int, doc string) *Method {
	// have to write out the function arguments - can't use the
	// type aliases as they are different types :-(
	switch method.(type) {
	case func(self Object, args Tuple) Object:
	case func(self Object, args Tuple, kwargs StringDict) Object:
	case func(Object) Object:
	case func(Object, Object) Object:
	default:
		panic(fmt.Sprintf("Unknown function type for NewMethod %q: %T\n", name, method))
	}
	return &Method{
		Name:   name,
		Doc:    doc,
		Flags:  flags,
		method: method,
	}
}

// Call the method with the given arguments
func (m *Method) Call(self Object, args Tuple) Object {
	switch f := m.method.(type) {
	case func(self Object, args Tuple) Object:
		return f(self, args)
	case func(self Object, args Tuple, kwargs StringDict) Object:
		return f(self, args, NewStringDict())
	case func(Object) Object:
		if len(args) != 0 {
			// FIXME type error
			panic(fmt.Sprintf("TypeError: %s() takes no arguments (%d given)", m.Name, len(args)))
		}
		return f(self)
	case func(Object, Object) Object:
		fmt.Printf("*** CALL %v %v\n", self, args)
		if len(args) != 1 {
			// FIXME type error
			panic(fmt.Sprintf("FOO TypeError: %s() takes exactly 1 argument (%d given)", m.Name, len(args)))
		}
		return f(self, args[0])
	}
	panic("Unknown method type")
}

// Call the method with the given arguments
func (m *Method) CallWithKeywords(self Object, args Tuple, kwargs StringDict) Object {
	switch f := m.method.(type) {
	case func(self Object, args Tuple, kwargs StringDict) Object:
		return f(self, args, kwargs)
	case func(self Object, args Tuple) Object:
	case func(Object) Object:
	case func(Object, Object) Object:
		// FIXME type error
		panic(fmt.Sprintf("TypeError: %s() takes no keyword arguments", m.Name))
	}
	panic("Unknown method type")
}

// Call a method
func (m *Method) M__call__(args Tuple, kwargs StringDict) Object {
	self := None // FIXME should be the module
	var result Object
	if kwargs != nil {
		result = m.CallWithKeywords(self, args, kwargs)
	} else {
		result = m.Call(self, args)
	}
	return result
}

// Make sure it satisfies the interface
var _ Object = (*Method)(nil)
var _ I__call__ = (*Method)(nil)
