// BoundMethod object
//
// Combines an object and a callable

package py

// A python BoundMethod object
type BoundMethod struct {
	Self   Object
	Method Object
}

var BoundMethodType = NewType("boundmethod", "boundmethod object")

// Type of this object
func (o *BoundMethod) Type() *Type {
	return BoundMethodType
}

// Define a new boundmethod
func NewBoundMethod(self, method Object) *BoundMethod {
	return &BoundMethod{Self: self, Method: method}
}

// Call the bound method
func (bm *BoundMethod) M__call__(args Tuple, kwargs StringDict) Object {
	newArgs := make(Tuple, len(args)+1)
	newArgs[0] = bm.Self
	copy(newArgs[1:], args)
	return Call(bm.Method, newArgs, kwargs)
}
