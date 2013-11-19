// Tuple objects

package py

var TupleType = NewType("tuple")

type Tuple []Object

// Type of this Tuple object
func (o Tuple) Type() *Type {
	return TupleType
}
