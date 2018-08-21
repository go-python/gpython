// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ellipsis objects

package py

type EllipsisType struct{}

var (
	EllipsisTypeType = NewType("EllipsisType", "")
	Ellipsis         = EllipsisType(struct{}{})
)

// Type of this object
func (s EllipsisType) Type() *Type {
	return EllipsisTypeType
}

func (a EllipsisType) M__bool__() (Object, error) {
	return False, nil
}

func (a EllipsisType) M__repr__() (Object, error) {
	return String("Ellipsis"), nil
}

func (a EllipsisType) M__eq__(other Object) (Object, error) {
	if _, ok := other.(EllipsisType); ok {
		return True, nil
	}
	return False, nil
}

func (a EllipsisType) M__ne__(other Object) (Object, error) {
	if _, ok := other.(EllipsisType); ok {
		return False, nil
	}
	return True, nil
}

// Check interface is satisfied
var _ I__bool__ = Ellipsis
var _ I__repr__ = Ellipsis
var _ I__eq__ = Ellipsis
var _ I__eq__ = Ellipsis
