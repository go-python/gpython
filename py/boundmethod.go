// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
func (bm *BoundMethod) M__call__(args Tuple, kwargs StringDict) (Object, error) {
	// Call built in methods slightly differently
	// FIXME not sure this is sensible! something is wrong with the call interface
	// as we aren't sure whether to call it with a self or not
	if m, ok := bm.Method.(*Method); ok {
		if kwargs != nil {
			return m.CallWithKeywords(bm.Self, args, kwargs)
		} else {
			return m.Call(bm.Self, args)
		}
	}
	newArgs := make(Tuple, len(args)+1)
	newArgs[0] = bm.Self
	copy(newArgs[1:], args)
	return Call(bm.Method, newArgs, kwargs)
}
