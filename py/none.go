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

// Check interface is satisfied
var _ I__bool__ = None
