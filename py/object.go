// Functions to operate on any object

package py

import (
	"fmt"
)

// Gets the attribute attr from object or returns nil
func ObjectGetAttr(o Object, attr string) Object {
	// FIXME
	return nil
}

// Gets the repr for an object
func ObjectRepr(o Object) Object {
	// FIXME
	return String(fmt.Sprintf("<%s %v>", o.Type().Name, o))
}
