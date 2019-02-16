// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Get the Dict
func (f *Function) GetDict() StringDict {
	return f.Dict
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
		Dict:     make(StringDict),
	}
}

// Call a function
func (f *Function) M__call__(args Tuple, kwargs StringDict) (Object, error) {
	result, err := VmEvalCodeEx(f.Code, f.Globals, NewStringDict(), args, kwargs, f.Defaults, f.KwDefaults, f.Closure)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Read a function from a class which makes a bound method
func (f *Function) M__get__(instance, owner Object) (Object, error) {
	if instance != None {
		return NewBoundMethod(instance, f), nil
	}
	return f, nil
}

// Properties
func init() {
	FunctionType.Dict["__code__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Function).Code, nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			// Not legal to set f.func_code to anything other than a code object.
			code, ok := value.(*Code)
			if !ok {
				return ExceptionNewf(TypeError, "__code__ must be set to a code object")
			}
			nfree := len(code.Freevars)
			nclosure := len(f.Closure)
			if nfree != nclosure {
				return ExceptionNewf(ValueError, "%s() requires a code object with %d free vars, not %d", f.Name, nclosure, nfree)
			}
			f.Code = code
			return nil
		},
	}
	FunctionType.Dict["__defaults__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Function).Defaults, nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			defaults, ok := value.(Tuple)
			if !ok {
				return ExceptionNewf(TypeError, "__defaults__ must be set to a tuple object")
			}
			f.Defaults = defaults
			return nil
		},
		Fdel: func(self Object) error {
			self.(*Function).Defaults = nil
			return nil
		},
	}
	FunctionType.Dict["__kwdefaults__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Function).KwDefaults, nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			kwdefaults, ok := value.(StringDict)
			if !ok {
				return ExceptionNewf(TypeError, "__kwdefaults__ must be set to a dict object")
			}
			f.KwDefaults = kwdefaults
			return nil
		},
		Fdel: func(self Object) error {
			self.(*Function).KwDefaults = nil
			return nil
		},
	}
	FunctionType.Dict["__annotations__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Function).Annotations, nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			annotations, ok := value.(StringDict)
			if !ok {
				return ExceptionNewf(TypeError, "__annotations__ must be set to a dict object")
			}
			f.Annotations = annotations
			return nil
		},
		Fdel: func(self Object) error {
			self.(*Function).Annotations = nil
			return nil
		},
	}
	FunctionType.Dict["__dict__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Function).Dict, nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			dict, ok := value.(StringDict)
			if !ok {
				return ExceptionNewf(TypeError, "__dict__ must be set to a dict object")
			}
			f.Dict = dict
			return nil
		},
	}
	FunctionType.Dict["__name__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return String(self.(*Function).Name), nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			name, ok := value.(String)
			if !ok {
				return ExceptionNewf(TypeError, "__name__ must be set to a string object")
			}
			f.Name = string(name)
			return nil
		},
	}
	FunctionType.Dict["__qualname__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return String(self.(*Function).Qualname), nil
		},
		Fset: func(self, value Object) error {
			f := self.(*Function)
			qualname, ok := value.(String)
			if !ok {
				return ExceptionNewf(TypeError, "__qualname__ must be set to a string object")
			}
			f.Qualname = string(qualname)
			return nil
		},
	}
}

// Make sure it satisfies the interface
var _ Object = (*Function)(nil)
var _ I__call__ = (*Function)(nil)
var _ IGetDict = (*Function)(nil)
var _ I__get__ = (*Function)(nil)
