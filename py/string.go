// String objects

package py

import (
	"fmt"
)

type String string

// Intern s possibly returning a reference to an already interned string
func (s String) Intern() String {
	fmt.Printf("FIXME interning %q\n", s)
	return s
}
