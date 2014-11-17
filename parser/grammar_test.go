package parser

import (
	"testing"

	"github.com/ncw/gpython/ast"
)

// FIXME test pos is correct

func TestGrammar(t *testing.T) {
	for _, test := range []struct {
		in   string
		mode string
		out  string
	}{
		{"", "exec", "Module(Body=[])"},
		{"pass", "exec", "Module(Body=[Pass()])"},
		{"()", "eval", "Expression(Body=Tuple(Elts=[],Ctx=UnknownExprContext(0)))"},
		{"()", "exec", "Module(Body=[ExprStmt(Value=Tuple(Elts=[],Ctx=UnknownExprContext(0)))])"},
		{"[ ]", "exec", "Module(Body=[ExprStmt(Value=List(Elts=[],Ctx=UnknownExprContext(0)))])"},
	} {
		Ast, err := ParseString(test.in, test.mode)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", test.in, err)
		} else {
			out := ast.Dump(Ast)
			if out != test.out {
				t.Errorf("Parse(%q) expecting %q actual %q", test.in, test.out, out)
			}
		}
	}
}
