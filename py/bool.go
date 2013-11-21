// Bool objects

package py

type Bool bool

var (
	BoolType = NewType("bool", "bool(x) -> bool\n\nReturns True when the argument x is true, False otherwise.\nThe builtins True and False are the only two instances of the class bool.\nThe class bool is a subclass of the class int, and cannot be subclassed.")
	// Some well known bools
	False = Bool(false)
	True  = Bool(true)
)

// Type of this object
func (s Bool) Type() *Type {
	return BoolType
}

func (a Bool) M__bool__() Object {
	return a
}

// Check interface is satisfied
var _ I__bool__ = Bool(false)
