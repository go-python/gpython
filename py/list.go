// List objects

package py

var ListType = NewType("list")

type List []Object

// Type of this List object
func (o List) Type() *Type {
	return ListType
}
