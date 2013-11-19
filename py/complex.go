// Complex objects

package py

var ComplexType = NewType("complex64")

type Complex complex64

// Type of this Complex object
func (o Complex) Type() *Type {
	return ComplexType
}
