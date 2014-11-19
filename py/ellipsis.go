// Ellipsis objects

package py

type EllipsisType struct{}

var (
	EllipsisTypeType = NewType("EllipsisType", "")
	Ellipsis         = EllipsisType(struct{}{})
)

// Type of this object
func (s EllipsisType) Type() *Type {
	return EllipsisTypeType
}

func (a EllipsisType) M__bool__() Object {
	return False
}

func (a EllipsisType) M__str__() Object {
	return String("Ellipsis")
}

// Check interface is satisfied
var _ I__bool__ = Ellipsis
var _ I__str__ = Ellipsis
