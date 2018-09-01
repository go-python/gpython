// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file holds the go generate command to run yacc on the grammar in grammar.y.
// To build y.go:
//      % go generate
//      % go build

//go:generate goyacc -v y.output grammar.y
package parser
