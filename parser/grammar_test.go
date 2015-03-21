package parser

//go:generate ./make_grammar_test.py

import (
	"flag"
	"testing"

	"github.com/ncw/gpython/ast"
)

var debugLevel = flag.Int("debugLevel", 0, "Debug level 0-4")

// FIXME test pos is correct

// FIXME add tests to test the error cases

func TestGrammar(t *testing.T) {
	SetDebug(*debugLevel)
	for _, test := range grammarTestData {
		Ast, err := ParseString(test.in, test.mode)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", test.in, err)
		} else {
			out := ast.Dump(Ast)
			if out != test.out {
				t.Errorf("Parse(%q)\nwant> %q\n got> %q\n", test.in, test.out, out)
			}
		}
	}
}
