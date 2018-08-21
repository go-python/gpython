// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	out := []string{}
	accumulate := func(ast Ast) bool {
		out = append(out, fmt.Sprintf("%T", ast))
		return true
	}

	for _, test := range []struct {
		in  Ast
		out []string
	}{
		{nil, []string{}},

		// An empty one of everything
		{&Module{}, []string{"*ast.Module"}},
		{&Interactive{}, []string{"*ast.Interactive"}},
		{&Expression{}, []string{"*ast.Expression"}},
		{&Suite{}, []string{"*ast.Suite"}},
		{&FunctionDef{}, []string{"*ast.FunctionDef"}},
		{&ClassDef{}, []string{"*ast.ClassDef"}},
		{&Return{}, []string{"*ast.Return"}},
		{&Delete{}, []string{"*ast.Delete"}},
		{&Assign{}, []string{"*ast.Assign"}},
		{&AugAssign{}, []string{"*ast.AugAssign"}},
		{&For{}, []string{"*ast.For"}},
		{&While{}, []string{"*ast.While"}},
		{&If{}, []string{"*ast.If"}},
		{&With{}, []string{"*ast.With"}},
		{&Raise{}, []string{"*ast.Raise"}},
		{&Try{}, []string{"*ast.Try"}},
		{&Assert{}, []string{"*ast.Assert"}},
		{&Import{}, []string{"*ast.Import"}},
		{&ImportFrom{}, []string{"*ast.ImportFrom"}},
		{&Global{}, []string{"*ast.Global"}},
		{&Nonlocal{}, []string{"*ast.Nonlocal"}},
		{&ExprStmt{}, []string{"*ast.ExprStmt"}},
		{&Pass{}, []string{"*ast.Pass"}},
		{&Break{}, []string{"*ast.Break"}},
		{&Continue{}, []string{"*ast.Continue"}},
		{&BoolOp{}, []string{"*ast.BoolOp"}},
		{&BinOp{}, []string{"*ast.BinOp"}},
		{&UnaryOp{}, []string{"*ast.UnaryOp"}},
		{&Lambda{}, []string{"*ast.Lambda"}},
		{&IfExp{}, []string{"*ast.IfExp"}},
		{&Dict{}, []string{"*ast.Dict"}},
		{&Set{}, []string{"*ast.Set"}},
		{&ListComp{}, []string{"*ast.ListComp"}},
		{&SetComp{}, []string{"*ast.SetComp"}},
		{&DictComp{}, []string{"*ast.DictComp"}},
		{&GeneratorExp{}, []string{"*ast.GeneratorExp"}},
		{&Yield{}, []string{"*ast.Yield"}},
		{&YieldFrom{}, []string{"*ast.YieldFrom"}},
		{&Compare{}, []string{"*ast.Compare"}},
		{&Call{}, []string{"*ast.Call"}},
		{&Num{}, []string{"*ast.Num"}},
		{&Str{}, []string{"*ast.Str"}},
		{&Bytes{}, []string{"*ast.Bytes"}},
		{&NameConstant{}, []string{"*ast.NameConstant"}},
		{&Ellipsis{}, []string{"*ast.Ellipsis"}},
		{&Attribute{}, []string{"*ast.Attribute"}},
		{&Subscript{}, []string{"*ast.Subscript"}},
		{&Starred{}, []string{"*ast.Starred"}},
		{&Name{}, []string{"*ast.Name"}},
		{&List{}, []string{"*ast.List"}},
		{&Tuple{}, []string{"*ast.Tuple"}},
		{&Slice{}, []string{"*ast.Slice"}},
		{&ExtSlice{}, []string{"*ast.ExtSlice"}},
		{&Index{}, []string{"*ast.Index"}},
		{&ExceptHandler{}, []string{"*ast.ExceptHandler"}},
		{&Arguments{}, []string{"*ast.Arguments"}},
		{&Arg{}, []string{"*ast.Arg"}},
		{&Keyword{}, []string{"*ast.Keyword"}},
		{&Alias{}, []string{"*ast.Alias"}},
		{&WithItem{}, []string{"*ast.WithItem"}},

		// Excercise the walk* closures
		{&Module{Body: []Stmt{&Pass{}}}, []string{"*ast.Module", "*ast.Pass"}},
		{&Module{Body: []Stmt{&Pass{}, &Pass{}}}, []string{"*ast.Module", "*ast.Pass", "*ast.Pass"}},
		{&Expression{Body: &Num{}}, []string{"*ast.Expression", "*ast.Num"}},
		{&Attribute{Value: &Num{}}, []string{"*ast.Attribute", "*ast.Num"}},
		{&List{Elts: []Expr{&Num{}, &Str{}}}, []string{"*ast.List", "*ast.Num", "*ast.Str"}},
		{&ListComp{Elt: &Num{}, Generators: []Comprehension{{Target: &Num{}, Iter: &Str{}, Ifs: []Expr{&Num{}, &Str{}}}}}, []string{"*ast.ListComp", "*ast.Num", "*ast.Num", "*ast.Str", "*ast.Num", "*ast.Str"}},
	} {
		out = nil
		Walk(test.in, accumulate)
		if len(out) != len(test.out) {
			t.Errorf("%q: differing length: want %#v got %#v", Dump(test.in), test.out, out)
		}
		for i := range out {
			if out[i] != test.out[i] {
				t.Errorf("%q: out[%d]: want %q got %q", Dump(test.in), i, test.out[i], out[i])
			}
		}
	}
}
