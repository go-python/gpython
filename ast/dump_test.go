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
		{&Str{S: py.String("potato")}, `Str(s="potato")`},
		{&Str{S: py.String("potato")}, `Str(s="potato")`},
		{&BinOp{Left: &Str{S: py.String("one")}, Op: Add, Right: &Str{S: py.String("two")}},
			`BinOp(left=Str(s="one"), op=Add(), right=Str(s="two"))`},
		{&Module{}, `Module(body=[])`},
		{&Module{Body: []Stmt{&Pass{}}}, `Module(body=[Pass()])`},
		{&Module{Body: []Stmt{&ExprStmt{Value: &Tuple{}}}}, `Module(body=[Expr(value=Tuple(elts=[], ctx=UnknownExprContext(0)))])`},
		{&NameConstant{Value: py.True}, `NameConstant(value=True)`},
	} {
		out := Dump(test.in)
		if out != test.out {
			t.Errorf("Dump(%#v) got %q expected %q", test.in, out, test.out)
		}
	}
}
