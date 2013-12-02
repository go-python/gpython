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
	Name        string     // The __name__ attribute, a string object
	Dict        StringDict // The __dict__ attribute, a dict or NULL
	Weakreflist List       // List of weak references
	Module      Object     // The __module__ attribute, can be anything
	Annotations StringDict // Annotations, a dict or NULL
	Qualname    string     // The qualified name
}

var FunctionType = NewType("function", "A python function")

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
func NewFunction(code *Code, globals StringDict, qualname string) *Function {
	var doc Object
	var module Object = None
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
		module = moduleobj
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

// Setup locals for calling the function with the given arguments
func (f *Function) LocalsForCall(args Tuple) StringDict {
	// fmt.Printf("call f %#v with %v\n", f, args)
	max := int(f.Code.Argcount)
	min := max - len(f.Defaults)
	if len(args) > max || len(args) < min {
		if min == max {
			panic(ExceptionNewf(TypeError, "%s() takes %d positional arguments but %d were given", f.Name, max))
		} else {
			panic(ExceptionNewf(TypeError, "%s() takes from %d to %d positional arguments but %d were given", f.Name, min, max))
		}
	}

	// FIXME not sure this is right!
	// Copy the args into the local variables
	locals := NewStringDict()
	for i := range args {
		locals[f.Code.Varnames[i]] = args[i]
	}
	for i := len(args); i < max; i++ {
		locals[f.Code.Varnames[i]] = f.Defaults[i-min]
	}
	// fmt.Printf("locals = %v\n", locals)
	return locals
}

// Call the function with the given arguments
func (f *Function) LocalsForCallWithKeywords(args Tuple, kwargs StringDict) StringDict {
	locals := NewStringDict()
	fmt.Printf("FIXME LocalsForCallWithKeywords NOT IMPLEMENTED\n")
	return locals
}

// Call a function
func (f *Function) M__call__(args Tuple, kwargs StringDict) Object {
	var locals StringDict
	if kwargs != nil {
		locals = f.LocalsForCallWithKeywords(args, kwargs)
	} else {
		locals = f.LocalsForCall(args)
	}
	result, err := Run(f.Globals, locals, f.Code, f.Closure)
	if err != nil {
		// Propagate the error
		panic(err)
	}
	return result
}

// Make sure it satisfies the interface
var _ Object = (*Function)(nil)
var _ I__call__ = (*Function)(nil)
