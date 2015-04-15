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
		{&Module{}, []string{"*ast.Module"}},
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
