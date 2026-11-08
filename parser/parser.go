// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file holds the go generate command to run yacc on the grammar in grammar.y.
// To build y.go:
//      % go generate
//      % go build

//go:generate goyacc -v y.output grammar.y
package parser

import (
	"strings"

	"github.com/go-python/gpython/ast"
)

// dottedName returns the expression for a dotted name such as "a.b.c": a
// Name for the first part and an Attribute for each further part, all at
// pos, as ast_for_dotted_name does in CPython.
func dottedName(pos ast.Pos, name string) ast.Expr {
	parts := strings.Split(name, ".")
	var expr ast.Expr = &ast.Name{ExprBase: ast.ExprBase{Pos: pos}, Id: ast.Identifier(parts[0]), Ctx: ast.Load}
	for _, part := range parts[1:] {
		expr = &ast.Attribute{ExprBase: ast.ExprBase{Pos: pos}, Value: expr, Attr: ast.Identifier(part), Ctx: ast.Load}
	}
	return expr
}
