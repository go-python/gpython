// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Python global definitions
package py

// Generate arithmetic boilerplate
//go:generate go run gen.go

// A python object
type Object interface {
	Type() *Type
}

// Optional interfaces
type IGetDict interface {
	GetDict() StringDict
}
type IGoInt interface {
	GoInt() (int, error)
}
type IGoInt64 interface {
	GoInt64() (int64, error)
}

var (
	// Set in vm/eval.go - to avoid circular import
	VmEvalCode func(ctx Context, code *Code, globals, locals StringDict, args []Object, kws StringDict, defs []Object, kwdefs StringDict, closure Tuple) (retval Object, err error)
	VmRunFrame func(frame *Frame) (res Object, err error)
)

// Called to create a new instance of class cls. __new__() is a static method (special-cased so you need not declare it as such) that takes the class of which an instance was requested as its first argument. The remaining arguments are those passed to the object constructor expression (the call to the class). The return value of __new__() should be the new object instance (usually an instance of cls).
//
// Typical implementations create a new instance of the class by invoking the superclass’s __new__() method using super(currentclass, cls).__new__(cls[, ...]) with appropriate arguments and then modifying the newly-created instance as necessary before returning it.
//
// If __new__() returns an instance of cls, then the new instance’s __init__() method will be invoked like __init__(self[, ...]), where self is the new instance and the remaining arguments are the same as were passed to __new__().
//
// If __new__() does not return an instance of cls, then the new instance’s __init__() method will not be invoked.
//
// __new__() is intended mainly to allow subclasses of immutable types (like int, str, or tuple) to customize instance creation. It is also commonly overridden in custom metaclasses in order to customize class creation.
// object.__new__(cls[, ...])
type I__new__ interface {
	M__new__(cls, args, kwargs Object) (Object, error)
}

// Called when the instance is created. The arguments are those passed
// to the class constructor expression. If a base class has an
// __init__() method, the derived class’s __init__() method, if any,
// must explicitly call it to ensure proper initialization of the base
// class part of the instance; for example: BaseClass.__init__(self,
// [args...]). As a special constraint on constructors, no value may
// be returned; doing so will cause a TypeError to be raised at
// runtime.
// object.__init__(self[, ...])
type I__init__ interface {
	M__init__(self, args, kwargs Object) (Object, error)
}

// Called when the instance is about to be destroyed. This is also
// called a destructor. If a base class has a __del__() method, the
// derived class’s __del__() method, if any, must explicitly call it
// to ensure proper deletion of the base class part of the
// instance. Note that it is possible (though not recommended!) for
// the __del__() method to postpone destruction of the instance by
// creating a new reference to it. It may then be called at a later
// time when this new reference is deleted. It is not guaranteed that
// __del__() methods are called for objects that still exist when the
// interpreter exits.
//
// Note del x doesn’t directly call x.__del__() — the former
// decrements the reference count for x by one, and the latter is only
// called when x‘s reference count reaches zero. Some common
// situations that may prevent the reference count of an object from
// going to zero include: circular references between objects (e.g., a
// doubly-linked list or a tree data structure with parent and child
// pointers); a reference to the object on the stack frame of a
// function that caught an exception (the traceback stored in
// sys.exc_info()[2] keeps the stack frame alive); or a reference to
// the object on the stack frame that raised an unhandled exception in
// interactive mode (the traceback stored in sys.last_traceback keeps
// the stack frame alive). The first situation can only be remedied by
// explicitly breaking the cycles; the latter two situations can be
// resolved by storing None in sys.last_traceback. Circular references
// which are garbage are detected and cleaned up when the cyclic
// garbage collector is enabled (it’s on by default). Refer to the
// documentation for the gc module for more information about this
// topic.
//
// Warning Due to the precarious circumstances under which __del__()
// methods are invoked, exceptions that occur during their execution
// are ignored, and a warning is printed to sys.stderr instead. Also,
// when __del__() is invoked in response to a module being deleted
// (e.g., when execution of the program is done), other globals
// referenced by the __del__() method may already have been deleted or
// in the process of being torn down (e.g. the import machinery
// shutting down). For this reason, __del__() methods should do the
// absolute minimum needed to maintain external invariants. Starting
// with version 1.5, Python guarantees that globals whose name begins
// with a single underscore are deleted from their module before other
// globals are deleted; if no other references to such globals exist,
// this may help in assuring that imported modules are still available
// at the time when the __del__() method is called.
// object.__del__(self)
type I__del__ interface {
	M__del__() (Object, error)
}

// Called by the repr() built-in function to compute the “official”
// string representation of an object. If at all possible, this should
// look like a valid Python expression that could be used to recreate
// an object with the same value (given an appropriate
// environment). If this is not possible, a string of the form
// <...some useful description...> should be returned. The return
// value must be a string object. If a class defines __repr__() but
// not __str__(), then __repr__() is also used when an “informal”
// string representation of instances of that class is required.
//
// This is typically used for debugging, so it is important that the
// representation is information-rich and unambiguous.
// object.__repr__(self)
type I__repr__ interface {
	M__repr__() (Object, error)
}

// Called by str(object) and the built-in functions format() and
// print() to compute the “informal” or nicely printable string
// representation of an object. The return value must be a string
// object.
//
// This method differs from object.__repr__() in that there is no
// expectation that __str__() return a valid Python expression: a more
// convenient or concise representation can be used.
//
// The default implementation defined by the built-in type object
// calls object.__repr__().
// object.__str__(self)
type I__str__ interface {
	M__str__() (Object, error)
}

// Called by bytes() to compute a byte-string representation of an
// object. This should return a bytes object.
// object.__bytes__(self)
type I__bytes__ interface {
	M__bytes__() (Object, error)
}

// Called by the format() built-in function (and by extension, the
// str.format() method of class str) to produce a “formatted” string
// representation of an object. The format_spec argument is a string
// that contains a description of the formatting options desired. The
// interpretation of the format_spec argument is up to the type
// implementing __format__(), however most classes will either
// delegate formatting to one of the built-in types, or use a similar
// formatting option syntax.
//
// See Format Specification Mini-Language for a description of the
// standard formatting syntax.
//
// The return value must be a string object.
// object.__format__(self, format_spec)
type I__format__ interface {
	M__format__(format_spec Object) (Object, error)
}

// These are the so-called “rich comparison” methods. The
// correspondence between operator symbols and method names is as
// follows: x<y calls x.__lt__(y), x<=y calls x.__le__(y), x==y calls
// x.__eq__(y), x!=y calls x.__ne__(y), x>y calls x.__gt__(y), and
// x>=y calls x.__ge__(y).
//
// A rich comparison method may return the singleton NotImplemented if
// it does not implement the operation for a given pair of
// arguments. By convention, False and True are returned for a
// successful comparison. However, these methods can return any value,
// so if the comparison operator is used in a Boolean context (e.g.,
// in the condition of an if statement), Python will call bool() on
// the value to determine if the result is true or false.
//
// There are no implied relationships among the comparison
// operators. The truth of x==y does not imply that x!=y is
// false. Accordingly, when defining __eq__(), one should also define
// __ne__() so that the operators will behave as expected. See the
// paragraph on __hash__() for some important notes on creating
// hashable objects which support custom comparison operations and are
// usable as dictionary keys.
//
// There are no swapped-argument versions of these methods (to be used
// when the left argument does not support the operation but the right
// argument does); rather, __lt__() and __gt__() are each other’s
// reflection, __le__() and __ge__() are each other’s reflection, and
// __eq__() and __ne__() are their own reflection.
//
// Arguments to rich comparison methods are never coerced.
//
// To automatically generate ordering operations from a single root
// operation, see functools.total_ordering().
// object.__lt__(self, other)
type I__lt__ interface {
	M__lt__(other Object) (Object, error)
}

// object.__le__(self, other)
type I__le__ interface {
	M__le__(other Object) (Object, error)
}

// object.__eq__(self, other)
type I__eq__ interface {
	M__eq__(other Object) (Object, error)
}

// object.__ne__(self, other)
type I__ne__ interface {
	M__ne__(other Object) (Object, error)
}

// object.__gt__(self, other)
type I__gt__ interface {
	M__gt__(other Object) (Object, error)
}

// object.__ge__(self, other)
type I__ge__ interface {
	M__ge__(other Object) (Object, error)
}

// Comparison operations
type richComparison interface {
	I__lt__
	I__le__
	I__eq__
	I__ne__
	I__gt__
	I__ge__
}

// Called by built-in function hash() and for operations on members of
// hashed collections including set, frozenset, and dict. __hash__()
// should return an integer. The only required property is that
// objects which compare equal have the same hash value; it is advised
// to somehow mix together (e.g. using exclusive or) the hash values
// for the components of the object that also play a part in
// comparison of objects.
//
// Note hash() truncates the value returned from an object’s custom
// __hash__() method to the size of a Py_ssize_t. This is typically 8
// bytes on 64-bit builds and 4 bytes on 32-bit builds. If an object’s
// __hash__() must interoperate on builds of different bit sizes, be
// sure to check the width on all supported builds. An easy way to do
// this is with python -c "import sys; print(sys.hash_info.width)"
//
// If a class does not define an __eq__() method it should not define
// a __hash__() operation either; if it defines __eq__() but not
// __hash__(), its instances will not be usable as items in hashable
// collections. If a class defines mutable objects and implements an
// __eq__() method, it should not implement __hash__(), since the
// implementation of hashable collections requires that a key’s hash
// value is immutable (if the object’s hash value changes, it will be
// in the wrong hash bucket).
//
// User-defined classes have __eq__() and __hash__() methods by
// default; with them, all objects compare unequal (except with
// themselves) and x.__hash__() returns an appropriate value such that
// x == y implies both that x is y and hash(x) == hash(y).
//
// A class that overrides __eq__() and does not define __hash__() will
// have its __hash__() implicitly set to None. When the __hash__()
// method of a class is None, instances of the class will raise an
// appropriate TypeError when a program attempts to retrieve their
// hash value, and will also be correctly identified as unhashable
// when checking isinstance(obj, collections.Hashable).
//
// If a class that overrides __eq__() needs to retain the
// implementation of __hash__() from a parent class, the interpreter
// must be told this explicitly by setting __hash__ =
// <ParentClass>.__hash__.
//
// If a class that does not override __eq__() wishes to suppress hash
// support, it should include __hash__ = None in the class
// definition. A class which defines its own __hash__() that
// explicitly raises a TypeError would be incorrectly identified as
// hashable by an isinstance(obj, collections.Hashable) call.
//
// Note By default, the __hash__() values of str, bytes and datetime
// objects are “salted” with an unpredictable random value. Although
// they remain constant within an individual Python process, they are
// not predictable between repeated invocations of Python.
//
// This is intended to provide protection against a denial-of-service
// caused by carefully-chosen inputs that exploit the worst case
// performance of a dict insertion, O(n^2) complexity. See
// http://www.ocert.org/advisories/ocert-2011-003.html for details.
//
// Changing hash values affects the iteration order of dicts, sets and
// other mappings. Python has never made guarantees about this
// ordering (and it typically varies between 32-bit and 64-bit
// builds).
//
// See also PYTHONHASHSEED.
//
// Changed in version 3.3: Hash randomization is enabled by default.
// object.__hash__(self)
type I__hash__ interface {
	M__hash__() (Object, error)
}

// Called to implement truth value testing and the built-in operation
// bool(); should return False or True. When this method is not
// defined, __len__() is called, if it is defined, and the object is
// considered true if its result is nonzero. If a class defines
// neither __len__() nor __bool__(), all its instances are considered
// true.
// object.__bool__(self)
type I__bool__ interface {
	M__bool__() (Object, error)
}

//The following methods can be defined to customize the meaning of attribute access (use of, assignment to, or deletion of x.name) for class instances.

// Called when an attribute lookup has not found the attribute in the
// usual places (i.e. it is not an instance attribute nor is it found
// in the class tree for self). name is the attribute name. This
// method should return the (computed) attribute value or raise an
// AttributeError exception.
//
// Note that if the attribute is found through the normal mechanism,
// __getattr__() is not called. (This is an intentional asymmetry
// between __getattr__() and __setattr__().) This is done both for
// efficiency reasons and because otherwise __getattr__() would have
// no way to access other attributes of the instance. Note that at
// least for instance variables, you can fake total control by not
// inserting any values in the instance attribute dictionary (but
// instead inserting them in another object). See the
// __getattribute__() method below for a way to actually get total
// control over attribute access.
// object.__getattr__(self, name)
type I__getattr__ interface {
	M__getattr__(name string) (Object, error)
}

// Called unconditionally to implement attribute accesses for
// instances of the class. If the class also defines __getattr__(),
// the latter will not be called unless __getattribute__() either
// calls it explicitly or raises an AttributeError. This method should
// return the (computed) attribute value or raise an AttributeError
// exception. In order to avoid infinite recursion in this method, its
// implementation should always call the base class method with the
// same name to access any attributes it needs, for example,
// object.__getattribute__(self, name).
//
// Note This method may still be bypassed when looking up special
// methods as the result of implicit invocation via language syntax or
// built-in functions. See Special method lookup.
// object.__getattribute__(self, name)
type I__getattribute__ interface {
	M__getattribute__(name string) (Object, error)
}

// Called when an attribute assignment is attempted. This is called
// instead of the normal mechanism (i.e. store the value in the
// instance dictionary). name is the attribute name, value is the
// value to be assigned to it.
//
// If __setattr__() wants to assign to an instance attribute, it
// should call the base class method with the same name, for example,
// object.__setattr__(self, name, value).
// object.__setattr__(self, name, value)
type I__setattr__ interface {
	M__setattr__(name string, value Object) (Object, error)
}

// Like __setattr__() but for attribute deletion instead of
// assignment. This should only be implemented if del obj.name is
// meaningful for the object.
// object.__delattr__(self, name)
type I__delattr__ interface {
	M__delattr__(name string) (Object, error)
}

// Called when dir() is called on the object. A sequence must be
// returned. dir() converts the returned sequence to a list and sorts
// it.
// object.__dir__(self)
type I__dir__ interface {
	M__dir__() (Object, error)
}

// The following methods only apply when an instance of the class
// containing the method (a so-called descriptor class) appears in an
// owner class (the descriptor must be in either the owner’s class
// dictionary or in the class dictionary for one of its parents). In
// the examples below, “the attribute” refers to the attribute whose
// name is the key of the property in the owner class’ __dict__.

// Called to get the attribute of the owner class (class attribute
// access) or of an instance of that class (instance attribute
// access). owner is always the owner class, while instance is the
// instance that the attribute was accessed through, or None when the
// attribute is accessed through the owner. This method should return
// the (computed) attribute value or raise an AttributeError
// exception.
// object.__get__(self, instance, owner)
type I__get__ interface {
	M__get__(instance, owner Object) (Object, error)
}

// Called to set the attribute on an instance of the owner
// class to a new value.
// object.__set__(self, instance, value)
type I__set__ interface {
	M__set__(instance, value Object) (Object, error)
}

// Called to delete the attribute on an instance instance of the owner
// class.
// object.__delete__(self, instance)
type I__delete__ interface {
	M__delete__(instance Object) (Object, error)
}

// The following methods are used to override the default behavior of
// the isinstance() and issubclass() built-in functions.
//
// In particular, the metaclass abc.ABCMeta implements these methods
// in order to allow the addition of Abstract Base Classes (ABCs) as
// “virtual base classes” to any class or type (including built-in
// types), including other ABCs.
//
// Note that these methods are looked up on the type (metaclass) of a
// class. They cannot be defined as class methods in the actual
// class. This is consistent with the lookup of special methods that
// are called on instances, only in this case the instance is itself a
// class.

// Return true if instance should be considered a (direct or indirect)
// instance of class. If defined, called to implement
// isinstance(instance, class).
// object.__instancecheck__(self, instance)
type I__instancecheck__ interface {
	M__instancecheck__(instance Object) (Object, error)
}

// Return true if subclass should be considered a (direct or indirect)
// subclass of class. If defined, called to implement
// issubclass(subclass, class).
// object.__subclasscheck__(self, subclass)
type I__subclasscheck__ interface {
	M__subclasscheck__(subclass Object) (Object, error)
}

// Called when the instance is “called” as a function; if this method
// is defined, x(arg1, arg2, ...) is a shorthand for x.__call__(arg1,
// arg2, ...).
// object.__call__(self[, args...])
type I__call__ interface {
	M__call__(args Tuple, kwargs StringDict) (Object, error)
}

// The following methods can be defined to implement container
// objects. Containers usually are sequences (such as lists or tuples)
// or mappings (like dictionaries), but can represent other containers
// as well. The first set of methods is used either to emulate a
// sequence or to emulate a mapping; the difference is that for a
// sequence, the allowable keys should be the integers k for which 0
// <= k < N where N is the length of the sequence, or slice objects,
// which define a range of items. It is also recommended that mappings
// provide the methods keys(), values(), items(), get(), clear(),
// setdefault(), pop(), popitem(), copy(), and update() behaving
// similar to those for Python’s standard dictionary objects. The
// collections module provides a MutableMapping abstract base class to
// help create those methods from a base set of __getitem__(),
// __setitem__(), __delitem__(), and keys(). Mutable sequences should
// provide methods append(), count(), index(), extend(), insert(),
// pop(), remove(), reverse() and sort(), like Python standard list
// objects. Finally, sequence types should implement addition (meaning
// concatenation) and multiplication (meaning repetition) by defining
// the methods __add__(), __radd__(), __iadd__(), __mul__(),
// __rmul__() and __imul__() described below; they should not define
// other numerical operators. It is recommended that both mappings and
// sequences implement the __contains__() method to allow efficient
// use of the in operator; for mappings, in should search the
// mapping’s keys; for sequences, it should search through the
// values. It is further recommended that both mappings and sequences
// implement the __iter__() method to allow efficient iteration
// through the container; for mappings, __iter__() should be the same
// as keys(); for sequences, it should iterate through the values.

// Called to implement the built-in function len(). Should return the
// length of the object, an integer >= 0. Also, an object that doesn’t
// define a __bool__() method and whose __len__() method returns zero
// is considered to be false in a Boolean context.
// object.__len__(self)
type I__len__ interface {
	M__len__() (Object, error)
}

// Called to implement operator.length_hint(). Should return an
// estimated length for the object (which may be greater or less than
// the actual length). The length must be an integer >= 0. This method
// is purely an optimization and is never required for correctness.
//
// New in version 3.4.
// object.__length_hint__(self)
type I__length_hint__ interface {
	M__length_hint__() (Object, error)
}

// Note Slicing is done exclusively with the following three methods. A call like
// a[1:2] = b
// is translated to
// a[slice(1, 2, None)] = b
// and so forth. Missing slice items are always filled in with None.

// Called to implement evaluation of self[key]. For sequence types,
// the accepted keys should be integers and slice objects. Note that
// the special interpretation of negative indexes (if the class wishes
// to emulate a sequence type) is up to the __getitem__() method. If
// key is of an inappropriate type, TypeError may be raised; if of a
// value outside the set of indexes for the sequence (after any
// special interpretation of negative values), IndexError should be
// raised. For mapping types, if key is missing (not in the
// container), KeyError should be raised.
//
// Note for loops expect that an IndexError will be raised for illegal
// indexes to allow proper detection of the end of the sequence.
// object.__getitem__(self, key)
type I__getitem__ interface {
	M__getitem__(key Object) (Object, error)
}

// Called to implement assignment to self[key]. Same note as for
// __getitem__(). This should only be implemented for mappings if the
// objects support changes to the values for keys, or if new keys can
// be added, or for sequences if elements can be replaced. The same
// exceptions should be raised for improper key values as for the
// __getitem__() method.
// object.__setitem__(self, key, value)
type I__setitem__ interface {
	M__setitem__(key, value Object) (Object, error)
}

// Called to implement deletion of self[key]. Same note as for
// __getitem__(). This should only be implemented for mappings if the
// objects support removal of keys, or for sequences if elements can
// be removed from the sequence. The same exceptions should be raised
// for improper key values as for the __getitem__() method.
// object.__delitem__(self, key)
type I__delitem__ interface {
	M__delitem__(key Object) (Object, error)
}

// This method is called when an iterator is required for a
// container. This method should return a new iterator object that can
// iterate over all the objects in the container. For mappings, it
// should iterate over the keys of the container, and should also be
// made available as the method keys().
//
// Iterator objects also need to implement this method; they are
// required to return themselves. For more information on iterator
// objects, see Iterator Types.
// object.__iter__(self)
type I__iter__ interface {
	M__iter__() (Object, error)
}

// The next method for iterators
type I__next__ interface {
	M__next__() (Object, error)
}

// Interface all iterators must satisfy
type I_iterator interface {
	I__iter__
	I__next__
}

// Generator interfaces
type I_send interface {
	Send(value Object) (Object, error)
}
type I_throw interface {
	Throw(args Tuple, kwargs StringDict) (Object, error)
}
type I_close interface {
	Close() (Object, error)
}

// Interface all generators must satisfy
type I_generator interface {
	I_iterator
	I_send
	I_throw
	I_close
}

// Called (if present) by the reversed() built-in to implement reverse
// iteration. It should return a new iterator object that iterates
// over all the objects in the container in reverse order.
//
// If the __reversed__() method is not provided, the reversed()
// built-in will fall back to using the sequence protocol (__len__()
// and __getitem__()). Objects that support the sequence protocol
// should only provide __reversed__() if they can provide an
// implementation that is more efficient than the one provided by
// reversed().
// object.__reversed__(self)
type I__reversed__ interface {
	M__reversed__() (Object, error)
}

// The membership test operators (in and not in) are normally
// implemented as an iteration through a sequence. However, container
// objects can supply the following special method with a more
// efficient implementation, which also does not require the object be
// a sequence.

// Called to implement membership test operators. Should return true
// if item is in self, false otherwise. For mapping objects, this
// should consider the keys of the mapping rather than the values or
// the key-item pairs.
//
// For objects that don’t define __contains__(), the membership test
// first tries iteration via __iter__(), then the old sequence
// iteration protocol via __getitem__(), see this section in the
// language reference.
// object.__contains__(self, item)
type I__contains__ interface {
	M__contains__(item Object) (Object, error)
}

// These methods are called to implement the binary arithmetic
// operations (+, -, *, /, //, %, divmod(), pow(), **, <<, >>, &, ^,
// |). For instance, to evaluate the expression x + y, where x is an
// instance of a class that has an __add__() method, x.__add__(y) is
// called. The __divmod__() method should be the equivalent to using
// __floordiv__() and __mod__(); it should not be related to
// __truediv__(). Note that __pow__() should be defined to accept an
// optional third argument if the ternary version of the built-in
// pow() function is to be supported.
//
// If one of those methods does not support the operation with the
// supplied arguments, it should return NotImplemented.

// object.__add__(self, other)
type I__add__ interface {
	M__add__(other Object) (Object, error)
}

// object.__sub__(self, other)
type I__sub__ interface {
	M__sub__(other Object) (Object, error)
}

// object.__mul__(self, other)
type I__mul__ interface {
	M__mul__(other Object) (Object, error)
}

// object.__truediv__(self, other)
type I__truediv__ interface {
	M__truediv__(other Object) (Object, error)
}

// object.__floordiv__(self, other)
type I__floordiv__ interface {
	M__floordiv__(other Object) (Object, error)
}

// object.__mod__(self, other)
type I__mod__ interface {
	M__mod__(other Object) (Object, error)
}

// object.__divmod__(self, other)
type I__divmod__ interface {
	M__divmod__(other Object) (Object, Object, error)
}

// object.__pow__(self, other[, modulo])
type I__pow__ interface {
	M__pow__(other, modulo Object) (Object, error)
}

// object.__lshift__(self, other)
type I__lshift__ interface {
	M__lshift__(other Object) (Object, error)
}

// object.__rshift__(self, other)
type I__rshift__ interface {
	M__rshift__(other Object) (Object, error)
}

// object.__and__(self, other)
type I__and__ interface {
	M__and__(other Object) (Object, error)
}

// object.__xor__(self, other)
type I__xor__ interface {
	M__xor__(other Object) (Object, error)
}

// object.__or__(self, other)
type I__or__ interface {
	M__or__(other Object) (Object, error)
}

// These methods are called to implement the binary arithmetic
// operations (+, -, *, /, //, %, divmod(), pow(), **, <<, >>, &, ^,
// |) with reflected (swapped) operands. These functions are only
// called if the left operand does not support the corresponding
// operation and the operands are of different types. [2] For
// instance, to evaluate the expression x - y, where y is an instance
// of a class that has an __rsub__() method, y.__rsub__(x) is called
// if x.__sub__(y) returns NotImplemented.
//
// Note that ternary pow() will not try calling __rpow__() (the
// coercion rules would become too complicated).
//
// Note If the right operand’s type is a subclass of the left
// operand’s type and that subclass provides the reflected method for
// the operation, this method will be called before the left operand’s
// non-reflected method. This behavior allows subclasses to override
// their ancestors’ operations.

// object.__radd__(self, other)
type I__radd__ interface {
	M__radd__(other Object) (Object, error)
}

// object.__rsub__(self, other)
type I__rsub__ interface {
	M__rsub__(other Object) (Object, error)
}

// object.__rmul__(self, other)
type I__rmul__ interface {
	M__rmul__(other Object) (Object, error)
}

// object.__rtruediv__(self, other)
type I__rtruediv__ interface {
	M__rtruediv__(other Object) (Object, error)
}

// object.__rfloordiv__(self, other)
type I__rfloordiv__ interface {
	M__rfloordiv__(other Object) (Object, error)
}

// object.__rmod__(self, other)
type I__rmod__ interface {
	M__rmod__(other Object) (Object, error)
}

// object.__rdivmod__(self, other)
type I__rdivmod__ interface {
	M__rdivmod__(other Object) (Object, Object, error)
}

// object.__rpow__(self, other)
type I__rpow__ interface {
	M__rpow__(other Object) (Object, error)
}

// object.__rlshift__(self, other)
type I__rlshift__ interface {
	M__rlshift__(other Object) (Object, error)
}

// object.__rrshift__(self, other)
type I__rrshift__ interface {
	M__rrshift__(other Object) (Object, error)
}

// object.__rand__(self, other)
type I__rand__ interface {
	M__rand__(other Object) (Object, error)
}

// object.__rxor__(self, other)
type I__rxor__ interface {
	M__rxor__(other Object) (Object, error)
}

// object.__ror__(self, other)
type I__ror__ interface {
	M__ror__(other Object) (Object, error)
}

// These methods are called to implement the augmented arithmetic
// assignments (+=, -=, *=, /=, //=, %=, **=, <<=, >>=, &=, ^=,
// |=). These methods should attempt to do the operation in-place
// (modifying self) and return the result (which could be, but does
// not have to be, self). If a specific method is not defined, the
// augmented assignment falls back to the normal methods. For
// instance, to execute the statement x += y, where x is an instance
// of a class that has an __iadd__() method, x.__iadd__(y) is
// called. If x is an instance of a class that does not define a
// __iadd__() method, x.__add__(y) and y.__radd__(x) are considered,
// as with the evaluation of x + y.

// object.__iadd__(self, other)
type I__iadd__ interface {
	M__iadd__(other Object) (Object, error)
}

// object.__isub__(self, other)
type I__isub__ interface {
	M__isub__(other Object) (Object, error)
}

// object.__imul__(self, other)
type I__imul__ interface {
	M__imul__(other Object) (Object, error)
}

// object.__itruediv__(self, other)
type I__itruediv__ interface {
	M__itruediv__(other Object) (Object, error)
}

// object.__ifloordiv__(self, other)
type I__ifloordiv__ interface {
	M__ifloordiv__(other Object) (Object, error)
}

// object.__imod__(self, other)
type I__imod__ interface {
	M__imod__(other Object) (Object, error)
}

// object.__ipow__(self, other[, modulo])

type I__ipow__ interface {
	M__ipow__(other, modulo Object) (Object, error)
}

// object.__ilshift__(self, other)
type I__ilshift__ interface {
	M__ilshift__(other Object) (Object, error)
}

// object.__irshift__(self, other)
type I__irshift__ interface {
	M__irshift__(other Object) (Object, error)
}

// object.__iand__(self, other)
type I__iand__ interface {
	M__iand__(other Object) (Object, error)
}

// object.__ixor__(self, other)
type I__ixor__ interface {
	M__ixor__(other Object) (Object, error)
}

// object.__ior__(self, other)
type I__ior__ interface {
	M__ior__(other Object) (Object, error)
}

// Called to implement the unary arithmetic operations (-, +, abs() and ~).

// object.__neg__(self)
type I__neg__ interface {
	M__neg__() (Object, error)
}

// object.__pos__(self)
type I__pos__ interface {
	M__pos__() (Object, error)
}

// object.__abs__(self)
type I__abs__ interface {
	M__abs__() (Object, error)
}

// object.__invert__(self)
type I__invert__ interface {
	M__invert__() (Object, error)
}

// Called to implement the built-in functions complex(), int(),
// float() and round(). Should return a value of the appropriate type.

// object.__complex__(self)
type I__complex__ interface {
	M__complex__() (Object, error)
}

// object.__int__(self)
type I__int__ interface {
	M__int__() (Object, error)
}

// object.__float__(self)
type I__float__ interface {
	M__float__() (Object, error)
}

// object.__round__(self, n)
type I__round__ interface {
	M__round__(n Object) (Object, error)
}

// Called to implement operator.index(). Also called whenever Python
// needs an integer object (such as in slicing, or in the built-in
// bin(), hex() and oct() functions). Must return an integer.

// object.__index__(self)
type I__index__ interface {
	M__index__() (Int, error)
}

// Int, Float and Complex should satisfy this
type floatArithmetic interface {
	I__neg__
	I__pos__
	I__abs__
	I__add__
	I__sub__
	I__mul__
	I__truediv__
	I__floordiv__
	I__mod__
	I__divmod__
	I__pow__
	I__radd__
	I__rsub__
	I__rmul__
	I__rtruediv__
	I__rfloordiv__
	I__rmod__
	I__rdivmod__
	I__rpow__
	I__iadd__
	I__isub__
	I__imul__
	I__itruediv__
	I__ifloordiv__
	I__imod__
	I__ipow__
}

// Int should satisfy this
type booleanArithmetic interface {
	I__invert__
	I__lshift__
	I__rshift__
	I__and__
	I__xor__
	I__or__
	I__rlshift__
	I__rrshift__
	I__rand__
	I__rxor__
	I__ror__
	I__ilshift__
	I__irshift__
	I__iand__
	I__ixor__
	I__ior__
}

// Float and Int should statisfy this
type conversionBetweenTypes interface {
	I__complex__
	I__int__
	I__float__
	I__round__
}

// String, Tuple, List should statisfy this
type sequenceArithmetic interface {
	I__add__
	I__mul__
	I__radd__
	I__rmul__
	I__iadd__
	I__imul__
}

// FIXME everything should statisfy this ?
// Make a basics interface
// I__bool__

// Int should statisfy this
// I__index__

// A context manager is an object that defines the runtime context to
// be established when executing a with statement. The context manager
// handles the entry into, and the exit from, the desired runtime
// context for the execution of the block of code. Context managers
// are normally invoked using the with statement (described in section
// The with statement), but can also be used by directly invoking
// their methods.
//
// Typical uses of context managers include saving and restoring
// various kinds of global state, locking and unlocking resources,
// closing opened files, etc.
//
// For more information on context managers, see Context Manager
// Types.

// Enter the runtime context related to this object. The with
// statement will bind this method’s return value to the target(s)
// specified in the as clause of the statement, if any.
// object.__enter__(self)
type I__enter__ interface {
	M__enter__() (Object, error)
}

// Exit the runtime context related to this object. The parameters
// describe the exception that caused the context to be exited. If the
// context was exited without an exception, all three arguments will
// be None.
//
// If an exception is supplied, and the method wishes to suppress the
// exception (i.e., prevent it from being propagated), it should
// return a true value. Otherwise, the exception will be processed
// normally upon exit from this method.
//
// Note that __exit__() methods should not reraise the passed-in
// exception; this is the caller’s responsibility.
// object.__exit__(self, exc_type, exc_value, traceback)
type I__exit__ interface {
	M__exit__(exc_type, exc_value, traceback Object) (Object, error)
}

// Return the ceiling of x, the smallest integer greater than or equal
// to x. If x is not a float, delegates to x.__ceil__(), which should
// return an Integral value.
// object.__float__(self)
type I__ceil__ interface {
	M__ceil__() (Object, error)
}

// Return the floor of x, the largest integer less than or equal to
// x. If x is not a float, delegates to x.__floor__(), which should
// return an Integral value.
type I__floor__ interface {
	M__floor__() (Object, error)
}

// Return the Real value x truncated to an Integral (usually an
// integer). Delegates to x.__trunc__().
type I__trunc__ interface {
	M__trunc__() (Object, error)
}
