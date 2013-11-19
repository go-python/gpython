// Exception objects

package py

// A python Exception object
type Exception struct {
	Name            string // name of the exception FIXME should be part of class machinery
	Args            Object
	Traceback       Object
	Context         Object
	Cause           Object
	SuppressContext bool
	Other           StringDict // anything else that we want to stuff in
}

var (
	ExceptionType = NewType("exception")

	// Some well known exceptions - these should be types?
	// FIXME exceptions should be created in builtins probably
	// their names should certainly go in there!
	NotImplemented = NewException("NotImplemented")
	StopIteration  = NewException("StopIteration")
)

// Type of this object
func (o *Exception) Type() *Type {
	return ExceptionType
}

// Define a new exception
//
// FIXME need inheritance machinery to make this work properly
func NewException(name string) *Exception {
	m := &Exception{
		Name: name,
	}
	return m
}
