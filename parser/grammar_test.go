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
		{"1234", "eval", "Expression(body=Num(n=1234))"},
		{"0x1234", "eval", "Expression(body=Num(n=4660))"},
		{"12.34", "eval", "Expression(body=Num(n=12.34))"},
		{"{ }", "eval", "Expression(body=Dict(keys=[], values=[]))"},
		{"{1}", "eval", "Expression(body=Set(elts=[Num(n=1)]))"},
		{"{1,2}", "eval", "Expression(body=Set(elts=[Num(n=1), Num(n=2)]))"},
		{"{1,2,3,}", "eval", "Expression(body=Set(elts=[Num(n=1), Num(n=2), Num(n=3)]))"},
		{"{ 'a':1 }", "eval", "Expression(body=Dict(keys=[Str(s='a')], values=[Num(n=1)]))"},
		{"{ 'a':1, 'b':2 }", "eval", "Expression(body=Dict(keys=[Str(s='a'), Str(s='b')], values=[Num(n=1), Num(n=2)]))"},
		{"{ 'a':{'aa':11, 'bb':{'aa':11, 'bb':22}}, 'b':{'aa':11, 'bb':22} }", "eval", "Expression(body=Dict(keys=[Str(s='a'), Str(s='b')], values=[Dict(keys=[Str(s='aa'), Str(s='bb')], values=[Num(n=11), Dict(keys=[Str(s='aa'), Str(s='bb')], values=[Num(n=11), Num(n=22)])]), Dict(keys=[Str(s='aa'), Str(s='bb')], values=[Num(n=11), Num(n=22)])]))"},
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
