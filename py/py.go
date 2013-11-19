// Python global definitions
package py

// A python object
type Object interface {
	Type() *Type
}

// Some well known objects
var (
	None, Elipsis Object
)

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
	M__add__(other Object) Object
}

// object.__sub__(self, other)
type I__sub__ interface {
	M__sub__(other Object) Object
}

// object.__mul__(self, other)
type I__mul__ interface {
	M__mul__(other Object) Object
}

// object.__truediv__(self, other)
type I__truediv__ interface {
	M__truediv__(other Object) Object
}

// object.__floordiv__(self, other)
type I__floordiv__ interface {
	M__floordiv__(other Object) Object
}

// object.__mod__(self, other)
type I__mod__ interface {
	M__mod__(other Object) Object
}

// object.__divmod__(self, other)
type I__divmod__ interface {
	M__divmod__(other Object) Object
}

// object.__pow__(self, other[, modulo])
type I__pow__ interface {
	M__pow__(other, modulo Object) Object
}

// object.__lshift__(self, other)
type I__lshift__ interface {
	M__lshift__(other Object) Object
}

// object.__rshift__(self, other)
type I__rshift__ interface {
	M__rshift__(other Object) Object
}

// object.__and__(self, other)
type I__and__ interface {
	M__and__(other Object) Object
}

// object.__xor__(self, other)
type I__xor__ interface {
	M__xor__(other Object) Object
}

// object.__or__(self, other)
type I__or__ interface {
	M__or__(other Object) Object
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
	M__radd__(other Object) Object
}

// object.__rsub__(self, other)
type I__rsub__ interface {
	M__rsub__(other Object) Object
}

// object.__rmul__(self, other)
type I__rmul__ interface {
	M__rmul__(other Object) Object
}

// object.__rtruediv__(self, other)
type I__rtruediv__ interface {
	M__rtruediv__(other Object) Object
}

// object.__rfloordiv__(self, other)
type I__rfloordiv__ interface {
	M__rfloordiv__(other Object) Object
}

// object.__rmod__(self, other)
type I__rmod__ interface {
	M__rmod__(other Object) Object
}

// object.__rdivmod__(self, other)
type I__rdivmod__ interface {
	M__rdivmod__(other Object) Object
}

// object.__rpow__(self, other)
type I__rpow__ interface {
	M__rpow__(other Object) Object
}

// object.__rlshift__(self, other)
type I__rlshift__ interface {
	M__rlshift__(other Object) Object
}

// object.__rrshift__(self, other)
type I__rrshift__ interface {
	M__rrshift__(other Object) Object
}

// object.__rand__(self, other)
type I__rand__ interface {
	M__rand__(other Object) Object
}

// object.__rxor__(self, other)
type I__rxor__ interface {
	M__rxor__(other Object) Object
}

// object.__ror__(self, other)
type I__ror__ interface {
	M__ror__(other Object) Object
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
	M__iadd__(other Object) Object
}

// object.__isub__(self, other)
type I__isub__ interface {
	M__isub__(other Object) Object
}

// object.__imul__(self, other)
type I__imul__ interface {
	M__imul__(other Object) Object
}

// object.__itruediv__(self, other)
type I__itruediv__ interface {
	M__itruediv__(other Object) Object
}

// object.__ifloordiv__(self, other)
type I__ifloordiv__ interface {
	M__ifloordiv__(other Object) Object
}

// object.__imod__(self, other)
type I__imod__ interface {
	M__imod__(other Object) Object
}

// object.__ipow__(self, other[, modulo])

type I__ipow__ interface {
	M__ipow__(other, modulo Object) Object
}

// object.__ilshift__(self, other)
type I__ilshift__ interface {
	M__ilshift__(other Object) Object
}

// object.__irshift__(self, other)
type I__irshift__ interface {
	M__irshift__(other Object) Object
}

// object.__iand__(self, other)
type I__iand__ interface {
	M__iand__(other Object) Object
}

// object.__ixor__(self, other)
type I__ixor__ interface {
	M__ixor__(other Object) Object
}

// object.__ior__(self, other)
type I__ior__ interface {
	M__ior__(other Object) Object
}

// Called to implement the unary arithmetic operations (-, +, abs() and ~).

// object.__neg__(self)
type I__neg__ interface {
	M__neg__() Object
}

// object.__pos__(self)
type I__pos__ interface {
	M__pos__() Object
}

// object.__abs__(self)
type I__abs__ interface {
	M__abs__() Object
}

// object.__invert__(self)
type I__invert__ interface {
	M__invert__() Object
}

// Called to implement the built-in functions complex(), int(),
// float() and round(). Should return a value of the appropriate type.

// object.__complex__(self)
type I__complex__ interface {
	M__complex__() Object
}

// object.__int__(self)
type I__int__ interface {
	M__int__() Object
}

// object.__float__(self)
type I__float__ interface {
	M__float__() Object
}

// object.__round__(self, n)
type I__round__ interface {
	M__round__(n Object) Object
}

// Called to implement operator.index(). Also called whenever Python
// needs an integer object (such as in slicing, or in the built-in
// bin(), hex() and oct() functions). Must return an integer.

// object.__index__(self)
type I__index__ interface {
	M__index__() Object
}
