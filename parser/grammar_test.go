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
		// START TESTS
		{"", "exec", "Module(body=[])"},
		{"pass", "exec", "Module(body=[Pass()])"},
		{"()", "eval", "Expression(body=Tuple(elts=[], ctx=Load()))"},
		{"()", "exec", "Module(body=[Expr(value=Tuple(elts=[], ctx=Load()))])"},
		{"[ ]", "exec", "Module(body=[Expr(value=List(elts=[], ctx=Load()))])"},
		{"True\n", "eval", "Expression(body=NameConstant(value=True))"},
		{"False\n", "eval", "Expression(body=NameConstant(value=False))"},
		{"None\n", "eval", "Expression(body=NameConstant(value=None))"},
		{"...", "eval", "Expression(body=Ellipsis())"},
		{"abc123", "eval", "Expression(body=Name(id='abc123', ctx=Load()))"},
		{"\"abc\"", "eval", "Expression(body=Str(s='abc'))"},
		{"\"abc\" \"\"\"123\"\"\"", "eval", "Expression(body=Str(s='abc123'))"},
		{"b'abc'", "eval", "Expression(body=Bytes(s=b'abc'))"},
		{"b'abc' b'''123'''", "eval", "Expression(body=Bytes(s=b'abc123'))"},
		// END TESTS
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
