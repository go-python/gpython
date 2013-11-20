// String objects

package py

import (
//	"fmt"
)

type String string

var StringType = NewType("string")

// Type of this object
func (s String) Type() *Type {
	return StringType
}

// Intern s possibly returning a reference to an already interned string
func (s String) Intern() String {
	// fmt.Printf("FIXME interning %q\n", s)
	return s
}
