// Python global definitions
package py

// A python object
type Object interface{}

// Some well known objects
var (
	None, False, True, StopIteration, Elipsis Object
)

// Some python types
// FIXME factor into own files probably

type Tuple []Object
type List []Object
type Set []Object
type Bytes []byte
type String string
