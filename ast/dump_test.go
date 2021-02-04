// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"testing"

	"github.com/go-python/gpython/py"
)

func TestDump(t *testing.T) {
	for _, test := range []struct {
		in  Ast
		out string
	}{
		{nil, `<nil>`},
		{&Pass{}, `Pass()`},
		{&Constant{Value: py.String("potato")}, `Str(s='potato')`},
		{&Constant{Value: py.String("potato")}, `Str(s='potato')`},
		{&Constant{Value: py.Bytes("potato")}, `Bytes(s=b'potato')`},
		{&BinOp{Left: &Constant{Value: py.String("one")}, Op: Add, Right: &Constant{Value: py.String("two")}},
			`BinOp(left=Str(s='one'), op=Add(), right=Str(s='two'))`},
		{&Module{}, `Module(body=[], type_ignores=[])`},
		{&Module{Body: []Stmt{&Pass{}}}, `Module(body=[Pass()])`},
		{&Module{Body: []Stmt{&ExprStmt{Value: &Tuple{}}}}, `Module(body=[Expr(value=Tuple(elts=[], ctx=UnknownExprContext(0)))], type_ignores=[])`},
		{&Constant{Value: py.True}, `Constant(value=True)`},
		{&Name{Id: Identifier("hello"), Ctx: Load}, `Name(id='hello', ctx=Load())`},
		{&ListComp{Elt: &Constant{Value: py.String("potato")}, Generators: []Comprehension{{
			Target: &Name{Id: Identifier("hello"), Ctx: Load},
		}}}, `ListComp(elt=Str(s='potato'), generators=[comprehension(target=Name(id='hello', ctx=Load()), iter=None, ifs=[])])`},
	} {
		out := Dump(test.in)
		if out != test.out {
			t.Errorf("Dump(%#v) got %q expected %q", test.in, out, test.out)
		}
	}
}
