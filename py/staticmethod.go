// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// StaticMethod objects

package py

var StaticMethodType = ObjectType.NewType("staticmethod",
	`staticmethod(function) -> method

Convert a function to be a static method.

A static method does not receive an implicit first argument.
To declare a static method, use this idiom:

     class C:
     def f(arg1, arg2, ...): ...
     f = staticmethod(f)

It can be called either on the class (e.g. C.f()) or on an instance
(e.g. C().f()).  The instance is ignored except for its class.

Static methods in Python are similar to those found in Java or C++.
For a more advanced concept, see the classmethod builtin.`, StaticMethodNew, nil)

type StaticMethod struct {
	Callable Object
	Dict     StringDict
}

// Type of this StaticMethod object
func (o StaticMethod) Type() *Type {
	return StaticMethodType
}

// Get the Dict
func (c *StaticMethod) GetDict() StringDict {
	return c.Dict
}

// StaticMethodNew
func StaticMethodNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	c := &StaticMethod{
		Dict: make(StringDict),
	}
	err = UnpackTuple(args, kwargs, "staticmethod", 1, 1, &c.Callable)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Read a staticmethod from a class - no bound method here
func (c *StaticMethod) M__get__(instance, owner Object) (Object, error) {
	return c.Callable, nil
}

// Properties
func init() {
	StaticMethodType.Dict["__func__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*StaticMethod).Callable, nil
		},
	}
}

// Check interface is satisfied
var _ IGetDict = (*StaticMethod)(nil)
var _ I__get__ = (*StaticMethod)(nil)
