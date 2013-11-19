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
