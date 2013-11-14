// Python global definitions
package py

// A python object
type Object interface{}

// Some well known objects
var (
	None, False, True, StopIteration, Elipsis Object
)
