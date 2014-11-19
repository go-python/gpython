// None objects

package py

type NoneType struct{}

var (
	NoneTypeType = NewType("NoneType", "")
	// And the ubiquitous
	None = NoneType(struct{}{})
)

// Type of this object
func (s NoneType) Type() *Type {
	return NoneTypeType
}

func (a NoneType) M__bool__() Object {
	return False
}

func (a NoneType) M__str__() Object {
	return String("None")
}

// Check interface is satisfied
var _ I__bool__ = None
var _ I__str__ = None
