// Float objects

package py

var FloatType = NewType("float")

type Float float64

// Type of this Float64 object
func (o Float) Type() *Type {
	return FloatType
}
