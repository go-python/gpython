// Method objects
//
// This is about the type 'builtin_function_or_method', not Python
// methods in user-defined classes.  See class.go for the latter.

package py

// Types for methods
type PyCFunction func(Object, Tuple) Object
type PyCFunctionWithKeywords func(Object, Tuple, StringDict) Object

const (

	// This is the typical calling convention, where the methods have the
	// type PyCFunction. The function expects two PyObject* values. The
	// first one is the self object for methods; for module functions, it
	// is the module object. The second parameter (often called args) is a
	// tuple object representing all arguments. This parameter is
	// typically processed using PyArg_ParseTuple() or
	// PyArg_UnpackTuple().
	METH_VARARGS = 0x0001

	// Methods with these flags must be of type
	// PyCFunctionWithKeywords. The function expects three parameters:
	// self, args, and a dictionary of all the keyword arguments. The flag
	// is typically combined with METH_VARARGS, and the parameters are
	// typically processed using PyArg_ParseTupleAndKeywords().
	METH_KEYWORDS = 0x0002

	// Methods without parameters don’t need to check whether arguments
	// are given if they are listed with the METH_NOARGS flag. They need
	// to be of type PyCFunction. The first parameter is typically named
	// self and will hold a reference to the module or object instance. In
	// all cases the second parameter will be NULL.
	METH_NOARGS = 0x0004

	// Methods with a single object argument can be listed with the METH_O
	// flag, instead of invoking PyArg_ParseTuple() with a "O"
	// argument. They have the type PyCFunction, with the self parameter,
	// and a PyObject* parameter representing the single argument.
	METH_O = 0x0008

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
	// C function implementation (two definitions, only one is used)
	method             PyCFunction
	methodWithKeywords PyCFunctionWithKeywords
}

var MethodType = NewType("method")

// Type of this object
func (o *Method) Type() *Type {
	return MethodType
}

// Define a new method
func NewMethod(name string, method PyCFunction, flags int, doc string) *Method {
	if flags&METH_KEYWORDS != 0 {
		panic("Can't set METH_KEYWORDS")
	}
	return &Method{
		Name:   name,
		Doc:    doc,
		Flags:  flags,
		method: method,
	}
}

// Define a new method with keyword arguments
func NewMethodWithKeywords(name string, method PyCFunctionWithKeywords, flags int, doc string) *Method {
	if flags&METH_KEYWORDS == 0 {
		panic("Must set METH_KEYWORDS")
	}
	return &Method{
		Name:               name,
		Doc:                doc,
		Flags:              flags,
		methodWithKeywords: method,
	}
}

// Call the method with the given arguments
func (m *Method) Call(self Object, args Tuple) Object {
	if m.method != nil {
		return m.method(self, args)
	}
	// FIXME or call with empty dict?
	return m.methodWithKeywords(self, args, NewStringDict())
}

// Call the method with the given arguments
func (m *Method) CallWithKeywords(self Object, args Tuple, kwargs StringDict) Object {
	if m.method != nil {
		panic("Can't call method with kwargs")
	}
	return m.methodWithKeywords(self, args, kwargs)
}

// Check it implements the interface
var _ Callable = (*Method)(nil)
