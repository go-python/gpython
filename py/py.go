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
type Bytes []byte

type Int64 int64
type Float64 float64
type Complex64 complex64

type Dict map[Object]Object
type Set map[Object]struct{}
