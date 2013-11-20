// Bool objects

package py

type Bool bool

var (
	BoolType = NewType("bool")
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
