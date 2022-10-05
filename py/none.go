// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// None objects

package py

type NoneType struct{}

var (
	NoneTypeType = NewType("NoneType", "")
	// And the ubiquitous
	None = NoneType(struct{}{})
)

// Type of this object
func (s NoneType) Type() *Type {
	return NoneTypeType
}

func (a NoneType) M__bool__() (Object, error) {
	return False, nil
}

func (a NoneType) M__str__() (Object, error) {
	return a.M__repr__()
}

func (a NoneType) M__repr__() (Object, error) {
	return String("None"), nil
}

// Convert an Object to an NoneType
//
// Returns ok as to whether the conversion worked or not
func convertToNoneType(other Object) (NoneType, bool) {
	switch b := other.(type) {
	case NoneType:
		return b, true
	}
	return None, false
}

func (a NoneType) M__eq__(other Object) (Object, error) {
	if _, ok := convertToNoneType(other); ok {
		return True, nil
	}
	return False, nil
}

func (a NoneType) M__ne__(other Object) (Object, error) {
	if _, ok := convertToNoneType(other); ok {
		return False, nil
	}
	return True, nil
}

// Check interface is satisfied
var _ I__bool__ = None
var _ I__str__ = None
var _ I__repr__ = None
var _ I__eq__ = None
var _ I__ne__ = None
