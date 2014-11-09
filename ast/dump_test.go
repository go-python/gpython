package ast

import (
	"testing"

	"github.com/ncw/gpython/py"
)

func TestDump(t *testing.T) {
	for _, test := range []struct {
		in  Ast
		out string
	}{
		{nil, `<nil>`},
		{&Pass{}, `Pass()`},
		{&Str{S: py.String("potato")}, `Str(S="potato")`},
		{&Str{S: py.String("potato")}, `Str(S="potato")`},
		{&BinOp{Left: &Str{S: py.String("one")}, Op: Add, Right: &Str{S: py.String("two")}},
			`BinOp(Left=Str(S="one"),Op=Add,Right=Str(S="two"))`},
		{&Module{}, `Module(Body=[])`},
		{&Module{Body: []Stmt{&Pass{}}}, `Module(Body=[Pass()])`},
	} {
		out := Dump(test.in)
		if out != test.out {
			t.Errorf("Dump(%#v) got %q expected %q", test.in, out, test.out)
		}
	}
}
