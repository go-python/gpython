// Bytes objects

package py

var BytesType = NewType("bytes")

type Bytes []byte

// Type of this Bytes object
func (o Bytes) Type() *Type {
	return BytesType
}
