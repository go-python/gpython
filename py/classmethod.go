// ClassMethod objects

package py

var ClassMethodType = ObjectType.NewType("classmethod",
	`classmethod(function) -> method

Convert a function to be a class method.

A class method receives the class as implicit first argument,
just like an instance method receives the instance.
To declare a class method, use this idiom:

  class C:
      def f(cls, arg1, arg2, ...): ...
      f = classmethod(f)

It can be called either on the class (e.g. C.f()) or on an instance
(e.g. C().f()).  The instance is ignored except for its class.
If a class method is called for a derived class, the derived class
object is passed as the implied first argument.

Class methods are different than C++ or Java static methods.
If you want those, see the staticmethod builtin.`, ClassMethodNew, nil)

type ClassMethod struct {
	Callable Object
	Dict     StringDict
}

// Type of this ClassMethod object
func (o ClassMethod) Type() *Type {
	return ClassMethodType
}

// Get the Dict
func (c *ClassMethod) GetDict() StringDict {
	return c.Dict
}

// ClassMethodNew
func ClassMethodNew(metatype *Type, args Tuple, kwargs StringDict) (res Object) {
	c := &ClassMethod{}
	UnpackTuple(args, kwargs, "classmethod", 1, 1, &c.Callable)
	return c
}

// Read a classmethod from a class which makes a bound method
func (c *ClassMethod) M__get__(instance, owner Object) Object {
	if owner == nil {
		owner = instance.Type()
	}
	return NewBoundMethod(owner, c.Callable)
}

// Properties
func init() {
	ClassMethodType.Dict["__func__"] = &Property{
		Fget: func(self Object) Object {
			return self.(*ClassMethod).Callable
		},
	}
}

// Check interface is satisfied
var _ IGetDict = (*ClassMethod)(nil)
var _ I__get__ = (*ClassMethod)(nil)
