// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Return whether the object is True or not
func ObjectIsTrue(o Object) bool {
	if o == True {
		return true
	}
	if o == False {
		return false
	}

	if o == None {
		return false
	}

	if I, ok := o.(I__bool__); ok {
		cmp, err := I.M__bool__()
		if err == nil && cmp == True {
			return true
		} else if err == nil && cmp == False {
			return false
		}
	}
	return false
}
