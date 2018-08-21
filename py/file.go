// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// File object
//
// FIXME cpython 3.3 has a compicated heirachy of types to implement
// this which we do not emulate yet

package py

import (
	"os"
)

var FileType = NewTypeX("file", `represents an open file`,
	nil, nil)

type File os.File

// Type of this object
func (o *File) Type() *Type {
	return FileType
}

// Check interface is satisfied
