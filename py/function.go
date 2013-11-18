// Function objects
//
// Function objects and code objects should not be confused with each other:
//
// Function objects are created by the execution of the 'def' statement.
// They reference a code object in their __code__ attribute, which is a
// purely syntactic object, i.e. nothing more than a compiled version of some
// source code lines.  There is one code object per source code "fragment",
// but each code object can be referenced by zero or many function objects
// depending only on how many times the 'def' statement in the source was
// executed so far.
package py

import (
	"fmt"
)

// A python Function object
type Function struct {
	Code        *Code      // A code object, the __code__ attribute
	Globals     StringDict // A dictionary (other mappings won't do)
	Defaults    Tuple      // NULL or a tuple
	KwDefaults  StringDict // NULL or a dict
	Closure     Tuple      // NULL or a tuple of cell objects
	Doc         Object     // The __doc__ attribute, can be anything
	Name        String     // The __name__ attribute, a string object
	Dict        StringDict // The __dict__ attribute, a dict or NULL
	Weakreflist List       // List of weak references
	Module      *Module    // The __module__ attribute, can be anything
	Annotations StringDict // Annotations, a dict or NULL
	Qualname    String     // The qualified name
}

var FunctionType = NewType("function")

// Type of this object
func (o *Function) Type() *Type {
	return FunctionType
}

// Define a new function
//
// Return a new function object associated with the code object
// code. globals must be a dictionary with the global variables
// accessible to the function.
//
// The function’s docstring, name and __module__ are retrieved from
// the code object, the argument defaults and closure are set to NULL.
//
// Allows to set the function object’s __qualname__
// attribute. qualname should be a unicode object or ""; if "", the
// __qualname__ attribute is set to the same value as its __name__
// attribute.
func NewFunction(code *Code, globals StringDict, qualname String) *Function {
	var doc Object
	var module *Module
	if len(code.Consts) >= 1 {
		doc = code.Consts[0]
		if _, ok := doc.(String); !ok {
			doc = None
		}
	} else {
		doc = None
	}

	// __module__: If module name is in globals, use it. Otherwise, use None.

	if moduleobj, ok := globals["__name__"]; ok {
		module = (moduleobj).(*Module)
	}

	if qualname == "" {
		qualname = code.Name
	}

	return &Function{
		Code:     code,
		Qualname: qualname,
		Globals:  globals,
		Name:     code.Name,
		Doc:      doc,
		Module:   module,
	}
}

// Call the function with the given arguments
func (f *Function) Call(self Object, args Tuple) Object {
	fmt.Printf("call f %#v with %v and %v\n", f, self, args)
	if len(f.Code.Varnames) < len(args) {
		panic("Too many args!")
		// FIXME don't know how to deal with default args
	}
	// FIXME not sure this is right!
	// Copy the args into the local variables
	locals := NewStringDict()
	for i := range args {
		locals[string(f.Code.Varnames[i].(String))] = args[i]
	}
	fmt.Printf("locals = %v\n", locals)
	// FIXME return vm.Run(f.Globals, locals, f.Code)
	return None
}

// Call the function with the given arguments
func (f *Function) CallWithKeywords(self Object, args Tuple, kwargs StringDict) Object {
	return None
}

// Check it implements the interface
var _ Callable = (*Function)(nil)
