// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/py"
)

func TestCountIndent(t *testing.T) {
	for _, test := range []struct {
		in       string
		expected int
	}{
		{"", 0},
		{" ", 1},
		{"  ", 2},
		{"   ", 3},
		{"    ", 4},
		{"     ", 5},
		{"      ", 6},
		{"       ", 7},
		{"        ", 8},
		{"\t", 8},
		{"\t\t", 16},
		{"\t\t\t", 24},
		{"\t ", 9},
		{"\t  ", 10},
		{" \t", 8},
		{"  \t", 8},
		{"   \t", 8},
		{"    \t", 8},
		{"     \t", 8},
		{"      \t", 8},
		{"       \t", 8},
		{"        \t", 16},
		{"         \t", 16},
		{"               \t", 16},
		{"                \t", 24},
		{"               \t ", 17},
		{"                \t ", 25},
	} {
		got := countIndent(test.in)
		if got != test.expected {
			t.Errorf("countIndent(%q) expecting %d got %d", test.in, test.expected, got)
		}
	}
}

func TestLexToken(t *testing.T) {
	yylval := &yySymType{}
	for _, test := range []struct {
		token       int
		valueString string
		valueObj    py.Object
		expected    string
	}{
		{NAME, "potato", nil, `"NAME" (57348) = py.String{potato} 0:0`},
		{STRING, "", py.String("potato"), `"STRING" (57351) = py.String{potato} 0:0`},
		{NUMBER, "", py.Int(1), `"NUMBER" (57352) = py.Int{1} 0:0`},
		{'+', "", nil, `"+" (43) 0:0`},
		{LTLTEQ, "", nil, `"<<=" (57367) 0:0`},
	} {
		yylval.str = test.valueString
		yylval.obj = test.valueObj
		lt := newLexToken(test.token, yylval)
		got := lt.String()
		if got != test.expected {
			t.Errorf("newLexToken(%d,%q,%v) expecting %q got %q", test.token, test.valueString, test.valueObj, test.expected, got)
		}
	}

}

func TestLexTokensEq(t *testing.T) {
	for _, test := range []struct {
		as       LexTokens
		bs       LexTokens
		expected bool
	}{
		{
			LexTokens{},
			LexTokens{},
			true,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			true,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			false,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			LexTokens{
				{NUMBER, py.Int(2), ast.Pos{1, 0}},
			},
			false,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			LexTokens{
				{NEWLINE, nil, ast.Pos{1, 0}},
			},
			false,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
				{NEWLINE, nil, ast.Pos{1, 0}},
				{ENDMARKER, nil, ast.Pos{1, 0}},
			},
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
				{NEWLINE, nil, ast.Pos{1, 0}},
				{ENDMARKER, nil, ast.Pos{1, 0}},
			},
			true,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
				{NEWLINE, nil, ast.Pos{1, 0}},
				{ENDMARKER, nil, ast.Pos{1, 0}},
			},
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
				{NEWLINE, nil, ast.Pos{1, 0}},
				{NEWLINE, nil, ast.Pos{1, 0}},
			},
			false,
		},
	} {
		got := test.as.Eq(test.bs)
		if got != test.expected {
			t.Errorf("LexTokensEq(%v, %v) expecting %v got %v", test.as, test.bs, test.expected, got)
		}
	}
}

func TestLexTokensString(t *testing.T) {
	for _, test := range []struct {
		as       LexTokens
		expected string
	}{
		{
			LexTokens{},
			"[]",
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 0}},
			},
			`[{"NUMBER" (57352) = py.Int{1} 1:0}, ]`,
		},
		{
			LexTokens{
				{NUMBER, py.Int(1), ast.Pos{1, 2}},
				{NUMBER, py.Int(1), ast.Pos{3, 4}},
			},
			`[{"NUMBER" (57352) = py.Int{1} 1:2}, {"NUMBER" (57352) = py.Int{1} 3:4}, ]`,
		},
	} {
		got := test.as.String()
		if got != test.expected {
			t.Errorf("LexTokensString(%v) expecting %q got %q", test.as, test.expected, got)
		}
	}
}

func TestLex(t *testing.T) {
	for _, test := range []struct {
		in        string
		errString string
		mode      string
		lts       LexTokens
	}{
		{"", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{ENDMARKER, nil, ast.Pos{1, 0}},
		}},
		{"", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{NEWLINE, nil, ast.Pos{1, 0}},
		}},
		{"\n", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{NEWLINE, nil, ast.Pos{2, 0}},
		}},
		{"pass", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{PASS, nil, ast.Pos{1, 0}},
		}},
		{"pass\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{PASS, nil, ast.Pos{1, 0}},
			{NEWLINE, nil, ast.Pos{1, 4}},
			{ENDMARKER, nil, ast.Pos{2, 0}},
		}},
		{"\n#hello\n  #comment\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{ENDMARKER, nil, ast.Pos{4, 0}},
		}},
		{"\n#hello\n\f  #comment\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{ENDMARKER, nil, ast.Pos{4, 0}},
		}},
		{"1\n 2\n", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{NEWLINE, nil, ast.Pos{1, 1}},
			{INDENT, nil, ast.Pos{2, 0}},
			{NUMBER, py.Int(2), ast.Pos{2, 1}},
			{NEWLINE, nil, ast.Pos{2, 2}},
			{DEDENT, nil, ast.Pos{3, 0}},
			{ENDMARKER, nil, ast.Pos{3, 0}},
		}},
		{"1", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{ENDMARKER, nil, ast.Pos{1, 1}},
		}},
		{"01", "illegal decimal with leading zero 1:0", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
		}},
		{"1", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{NEWLINE, nil, ast.Pos{1, 1}},
			{ENDMARKER, nil, ast.Pos{1, 1}},
		}},
		{"1 2 3", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{NUMBER, py.Int(2), ast.Pos{1, 2}},
			{NUMBER, py.Int(3), ast.Pos{1, 4}},
			{NEWLINE, nil, ast.Pos{1, 5}},
			{ENDMARKER, nil, ast.Pos{1, 5}},
		}},
		{"01", "illegal decimal with leading zero 1:0", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
		}},
		{"1\n 2\n  3\n4\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{NEWLINE, nil, ast.Pos{1, 1}},
			{INDENT, nil, ast.Pos{2, 0}},
			{NUMBER, py.Int(2), ast.Pos{2, 1}},
			{NEWLINE, nil, ast.Pos{2, 2}},
			{INDENT, nil, ast.Pos{3, 0}},
			{NUMBER, py.Int(3), ast.Pos{3, 2}},
			{NEWLINE, nil, ast.Pos{3, 3}},
			{DEDENT, nil, ast.Pos{4, 0}},
			{DEDENT, nil, ast.Pos{4, 0}},
			{NUMBER, py.Int(4), ast.Pos{4, 0}},
			{NEWLINE, nil, ast.Pos{4, 1}},
			{ENDMARKER, nil, ast.Pos{5, 0}},
		}},
		{"if 1:\n  pass \n pass\n", "Inconsistent indent 3:1", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{IF, nil, ast.Pos{1, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 3}},
			{':', nil, ast.Pos{1, 4}},
			{NEWLINE, nil, ast.Pos{1, 5}},
			{INDENT, nil, ast.Pos{2, 0}},
			{PASS, nil, ast.Pos{2, 2}},
			{NEWLINE, nil, ast.Pos{2, 6}},
		}},
		{"(\n  1\n)", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{'(', nil, ast.Pos{1, 0}},
			{NUMBER, py.Int(1), ast.Pos{2, 2}},
			{')', nil, ast.Pos{3, 0}},
			{NEWLINE, nil, ast.Pos{3, 1}},
			{ENDMARKER, nil, ast.Pos{3, 1}},
		}},
		{"{\n  1\n}", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{'{', nil, ast.Pos{1, 0}},
			{NUMBER, py.Int(1), ast.Pos{2, 2}},
			{'}', nil, ast.Pos{3, 0}},
		}},
		{"[\n  1\n]", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{'[', nil, ast.Pos{1, 0}},
			{NUMBER, py.Int(1), ast.Pos{2, 2}},
			{']', nil, ast.Pos{3, 0}},
			{ENDMARKER, nil, ast.Pos{3, 1}},
		}},
		{"1\\\n2", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{NUMBER, py.Int(2), ast.Pos{2, 0}},
			{ENDMARKER, nil, ast.Pos{2, 1}},
		}},
		{"1\\\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{ENDMARKER, nil, ast.Pos{2, 0}},
		}},
		{"1\\", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{ENDMARKER, nil, ast.Pos{1, 1}},
		}},
		{"'1\\\n2'", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{STRING, py.String("12"), ast.Pos{1, 0}},
		}},
		{"0x1234 +\t0.1-6.1j", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(0x1234), ast.Pos{1, 0}},
			{'+', nil, ast.Pos{1, 7}},
			{NUMBER, py.Float(0.1), ast.Pos{1, 9}},
			{'-', nil, ast.Pos{1, 12}},
			{NUMBER, py.Complex(complex(0, 6.1)), ast.Pos{1, 13}},
			{ENDMARKER, nil, ast.Pos{1, 17}},
		}},
		{"001", "illegal decimal with leading zero 1:0", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
		}},
		{"u'''1\n2\n'''", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{STRING, py.String("1\n2\n"), ast.Pos{1, 0}},
			{ENDMARKER, nil, ast.Pos{3, 3}},
		}},
		{"\"hello\n", "EOL while scanning string literal 1:1", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
		}},
		{"1 >>-3\na <<=+12", "", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
			{NUMBER, py.Int(1), ast.Pos{1, 0}},
			{GTGT, nil, ast.Pos{1, 2}},
			{'-', nil, ast.Pos{1, 4}},
			{NUMBER, py.Int(3), ast.Pos{1, 5}},
			{NEWLINE, nil, ast.Pos{1, 6}},
			{NAME, py.String("a"), ast.Pos{2, 0}},
			{LTLTEQ, nil, ast.Pos{2, 2}},
			{'+', nil, ast.Pos{2, 5}},
			{NUMBER, py.Int(12), ast.Pos{2, 6}},
			{ENDMARKER, nil, ast.Pos{2, 8}},
		}},
		{"$asdasd", "invalid syntax 1:0", "eval", LexTokens{
			{EVAL_INPUT, nil, ast.Pos{0, 0}},
		}},
		{"if True:\n   pass\n\n", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{IF, nil, ast.Pos{1, 0}},
			{TRUE, nil, ast.Pos{1, 3}},
			{':', nil, ast.Pos{1, 7}},
			{NEWLINE, nil, ast.Pos{1, 8}},
			{INDENT, nil, ast.Pos{2, 0}},
			{PASS, nil, ast.Pos{2, 3}},
			{NEWLINE, nil, ast.Pos{2, 7}},
			{DEDENT, nil, ast.Pos{4, 0}},
			{NEWLINE, nil, ast.Pos{4, 0}},
		}},
		{"while True:\n pass\nelse:\n return\n", "", "single", LexTokens{
			{SINGLE_INPUT, nil, ast.Pos{0, 0}},
			{WHILE, nil, ast.Pos{1, 0}},
			{TRUE, nil, ast.Pos{1, 6}},
			{':', nil, ast.Pos{1, 10}},
			{NEWLINE, nil, ast.Pos{1, 11}},
			{INDENT, nil, ast.Pos{2, 0}},
			{PASS, nil, ast.Pos{2, 1}},
			{NEWLINE, nil, ast.Pos{2, 5}},
			{DEDENT, nil, ast.Pos{3, 0}},
			{ELSE, nil, ast.Pos{3, 0}},
			{':', nil, ast.Pos{3, 4}},
			{NEWLINE, nil, ast.Pos{3, 5}},
			{INDENT, nil, ast.Pos{4, 0}},
			{RETURN, nil, ast.Pos{4, 1}},
			{NEWLINE, nil, ast.Pos{4, 7}},
			{DEDENT, nil, ast.Pos{5, 0}},
			{NEWLINE, nil, ast.Pos{5, 0}},
		}},
		{"while True:\n pass\nelse:\n return\n", "", "exec", LexTokens{
			{FILE_INPUT, nil, ast.Pos{0, 0}},
			{WHILE, nil, ast.Pos{1, 0}},
			{TRUE, nil, ast.Pos{1, 6}},
			{':', nil, ast.Pos{1, 10}},
			{NEWLINE, nil, ast.Pos{1, 11}},
			{INDENT, nil, ast.Pos{2, 0}},
			{PASS, nil, ast.Pos{2, 1}},
			{NEWLINE, nil, ast.Pos{2, 5}},
			{DEDENT, nil, ast.Pos{3, 0}},
			{ELSE, nil, ast.Pos{3, 0}},
			{':', nil, ast.Pos{3, 4}},
			{NEWLINE, nil, ast.Pos{3, 5}},
			{INDENT, nil, ast.Pos{4, 0}},
			{RETURN, nil, ast.Pos{4, 1}},
			{NEWLINE, nil, ast.Pos{4, 7}},
			{DEDENT, nil, ast.Pos{5, 0}},
			{ENDMARKER, nil, ast.Pos{5, 0}},
		}},
	} {
		lts, err := LexString(test.in, test.mode)
		errString := ""
		if err != nil {
			lineno := -1
			offset := -1
			if exc, ok := err.(*py.Exception); ok {
				lineno = int(exc.Dict["lineno"].(py.Int))
				offset = int(exc.Dict["offset"].(py.Int))
				errString = fmt.Sprintf("%s %d:%d", exc.Args.(py.Tuple)[0], lineno, offset)
			} else {
				panic("bad exception")
			}
		}
		if errString != test.errString || !lts.Eq(test.lts) {
			t.Errorf("Lex(%q) expecting (%v, %q) got (%v, %q)", test.in, test.lts, test.errString, lts, errString)
			n := len(lts)
			if len(test.lts) > n {
				n = len(test.lts)
			}
			for i := 0; i < n; i++ {
				var want, got LexToken
				if i < len(lts) {
					got = lts[i]
				}
				if i < len(test.lts) {
					want = test.lts[i]
				}
				if want != got {
					t.Logf(">>> want[%d] = %v", i, &want)
					t.Logf(">>>  got[%d] = %v", i, &got)
				}
			}
		}
	}
}

func TestLexerIsIdentifier(t *testing.T) {
	for _, test := range []struct {
		in    rune
		start bool
		char  bool
	}{
		{'a', true, true},
		{'r', true, true},
		{'z', true, true},
		{'A', true, true},
		{'R', true, true},
		{'Z', true, true},
		{'0', false, true},
		{'4', false, true},
		{'9', false, true},
		{'_', true, true},
		{'@', false, false},
		{'[', false, false},
		{' ', false, false},
		{'\t', false, false},
		{'µ', true, true},
		{'©', false, false},
		{'—', false, false},
	} {
		got := isIdentifierStart(test.in)
		if got != test.start {
			t.Errorf("isIdentifierStart(%q) got %v expected %v", test.in, got, test.start)
		}
		got = isIdentifierChar(test.in)
		if got != test.char {
			t.Errorf("isIdentifierChar(%q) got %v expected %v", test.in, got, test.char)
		}
	}

}

func TestLexerReadIdentifier(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in       string
		expected string
		out      string
	}{
		{"", "", ""},
		{"1", "", "1"},
		{"potato", "potato", ""},
		{"break+", "break", "+"},
		{"_aAzZ09ß²¹", "_aAzZ09ß", "²¹"},
		{"123abc", "", "123abc"},
		{" abc", "", " abc"},
		{"+abc", "", "+abc"},
	} {
		x.line = test.in
		got := x.readIdentifier()
		if got != test.expected || x.line != test.out {
			t.Errorf("readIdentifier(%q) got %q remainder %q, expected %q remainder %q", test.in, got, x.line, test.expected, test.out)
		}
	}
}

func TestLexerReadIdentifierOrKeyword(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in    string
		token int
		value string
		out   string
	}{
		{"", eof, "", ""},
		{"1", eof, "", "1"},
		{"potato", NAME, "potato", ""},
		{"break+", BREAK, "break", "+"},
		{"breaker+", NAME, "breaker", "+"},
		{"_aAzZ09ß²¹", NAME, "_aAzZ09ß", "²¹"},
		{"123abc", eof, "", "123abc"},
		{" abc", eof, "", " abc"},
		{"+abc", eof, "", "+abc"},
	} {
		x.line = test.in
		token, value := x.readIdentifierOrKeyword()
		if token != test.token || value != test.value || x.line != test.out {
			t.Errorf("readIdentifierOrKeyword(%q) got (%q,%q) remainder %q, expected (%q,%q) remainder %q", test.in, tokenToString[token], value, x.line, tokenToString[test.token], test.value, test.out)
		}
	}
}

func TestLexerReadOperator(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in       string
		expected int
		out      string
	}{
		{"", eof, ""},
		{" <", eof, " <"},
		{"<", '<', ""},
		{"<+", '<', "+"},
		{"<< ", LTLT, " "},
		{"<<=", LTLTEQ, ""},
		{"<<< ", LTLT, "< "},
		{"<==", LTEQ, "="},
		{"/", '/', ""},
		{"//", DIVDIV, ""},
		{"=//", '=', "//"},
		{"//=", DIVDIVEQ, ""},
		{"....", ELIPSIS, "."},
	} {
		x.line = test.in
		got := x.readOperator()
		if got != test.expected || x.line != test.out {
			t.Errorf("readOperator(%q) got %q remainder %q, expected %q remainder %q", test.in, tokenToString[got], x.line, tokenToString[test.expected], test.out)
		}
	}
}

// Whether two floats are more or less the same
func approxEq(a, b float64) bool {
	log.Printf("ApproxEq(a = %#v, b = %#v)", a, b)
	diff := a - b
	log.Printf("ApproxEq(diff = %e)", diff)
	if math.Abs(diff) > 1e-10 {
		log.Printf("ApproxEq(false)")
		return false
	}
	log.Printf("ApproxEq(true)")
	return true
}

func TestLexerReadNumber(t *testing.T) {
	x := yyLex{}
	for _, test := range []struct {
		in    string
		token int
		value py.Object
		out   string
	}{
		{"", eof, nil, ""},
		{"break", eof, py.Object(nil), "break"},

		{"0o0", NUMBER, py.Int(0), ""},
		{"0O765a", NUMBER, py.Int(0765), "a"},
		{"0o0007779", NUMBER, py.Int(0777), "9"},

		{"0x0", NUMBER, py.Int(0), ""},
		{"0XaBcDeFg", NUMBER, py.Int(0xABCDEF), "g"},
		{"0x000123z", NUMBER, py.Int(0x123), "z"},
		{"0x0b", NUMBER, py.Int(11), ""},

		{"0b0", NUMBER, py.Int(0), ""},
		{"0B100", NUMBER, py.Int(4), ""},
		{"0B0001112", NUMBER, py.Int(7), "2"},

		{"1.", NUMBER, py.Float(1), ""},
		{".1", NUMBER, py.Float(.1), ""},
		{"0.1", NUMBER, py.Float(0.1), ""},
		{"00000.10000", NUMBER, py.Float(0.1), ""},
		{"1.E1", NUMBER, py.Float(10), ""},
		{".1e1", NUMBER, py.Float(1), ""},
		{"0.1e-01", NUMBER, py.Float(0.01), ""},
		{"00000.10000E-000001", NUMBER, py.Float(0.01), ""},
		{"1.j", NUMBER, py.Complex(complex(0, 1)), ""},
		{".1j", NUMBER, py.Complex(complex(0, .1)), ""},
		{"0.1j", NUMBER, py.Complex(complex(0, 0.1)), ""},
		{"00000.10000j", NUMBER, py.Complex(complex(0, 0.1)), ""},
		{"1.E1j", NUMBER, py.Complex(complex(0, 10)), ""},
		{".1e1j", NUMBER, py.Complex(complex(0, 1)), ""},
		{"0.1e-01j", NUMBER, py.Complex(complex(0, 0.01)), ""},
		{"00000.10000E-000001j", NUMBER, py.Complex(complex(0, 0.01)), ""},

		{"1", NUMBER, py.Int(1), ""},
		{"1+2", NUMBER, py.Int(1), "+2"},
		{"01", eofError, nil, "01"},
		{"00", NUMBER, py.Int(0), ""},
		{"123", NUMBER, py.Int(123), ""},
		{"0123", eofError, nil, "0123"},
		{"0123j", NUMBER, py.Complex(complex(0, 123)), ""},
		{"00j", NUMBER, py.Complex(complex(0, 0)), ""},
	} {
		x.line = test.in
		token, value := x.readNumber()
		if token != test.token || value != test.value || x.line != test.out {
			t.Errorf("readNumber(%q) got (%q,%T,%#v) remainder %q, expected (%q,%T,%#v) remainder %q", test.in, tokenToString[token], value, value, x.line, tokenToString[test.token], test.value, test.value, test.out)
		}
	}
}

func TestLexerReadString(t *testing.T) {
	for _, test := range []struct {
		in    string
		token int
		value py.Object
		out   string
	}{
		{``, eof, nil, ``},
		{`1`, eof, nil, `1`},

		{`""a`, STRING, py.String(""), `a`},
		{`u"abc"`, STRING, py.String("abc"), ``},
		{`"a\nc"`, STRING, py.String("a\nc"), ``},
		{`r"a\nc"`, STRING, py.String(`a\nc`), ``},
		{`"a\"c"`, STRING, py.String("a\"c"), ``},
		{`"a\\"+`, STRING, py.String("a\\"), `+`},
		{`"a`, eofError, nil, "a"},
		{"\"a\n", eofError, nil, "a\n"},
		{"\"a\\\nb\"c", STRING, py.String(`ab`), `c`},

		{`''a`, STRING, py.String(``), `a`},
		{`U'abc'`, STRING, py.String(`abc`), ``},
		{`'a\nc'`, STRING, py.String("a\nc"), ``},
		{`R'a\nc'`, STRING, py.String(`a\nc`), ``},
		{`'a\'c'`, STRING, py.String("a'c"), ``},
		{`'\n`, eofError, nil, `\n`},
		{`'a`, eofError, nil, `a`},
		{"'\\\n\\\npotato\\\nX\\\n'c", STRING, py.String(`potatoX`), `c`},

		{`""""""a`, STRING, py.String(``), `a`},
		{`u"""abc"""`, STRING, py.String(`abc`), ``},
		{`"""a\nc"""`, STRING, py.String("a\nc"), ``},
		{`r"""a\"""c"""`, STRING, py.String(`a\"""c`), ``},
		{`"""a\"""c"""`, STRING, py.String(`a"""c`), ``},
		{`"""a`, eofError, nil, `a`},
		{"\"\"\"a\nb\nc\n\"\"\"\n", STRING, py.String("a\nb\nc\n"), "\n"},
		{"\"\"\"a\nb\nc\na", eofError, nil, "a"},
		{"\"\"\"a\\\nb\"\"\"c", STRING, py.String(`ab`), `c`},

		{`''''''a`, STRING, py.String(``), `a`},
		{`U'''abc'''`, STRING, py.String(`abc`), ``},
		{`'''a\nc'''`, STRING, py.String("a\nc"), ``},
		{`R'''a\nc'''`, STRING, py.String(`a\nc`), ``},
		{`'''a\'''c'''`, STRING, py.String(`a'''c`), ``},
		{`'''a`, eofError, nil, `a`},
		{"'''a\nb\nc\n'''\n", STRING, py.String("a\nb\nc\n"), "\n"},
		{"'''a\nb\nc\na", eofError, nil, "a"},
		{"'''\\\na\\\nb\\\n'''c", STRING, py.String(`ab`), `c`},

		{`b""a`, STRING, py.Bytes{}, "a"},
		{`b'abc'`, STRING, py.Bytes(string(`abc`)), ``},
		{`B"""a\nc"""`, STRING, py.Bytes(string("a\nc")), ``},
		{`B'''a\"c'''`, STRING, py.Bytes(string(`a"c`)), ``},

		{`rb""a`, STRING, py.Bytes{}, "a"},
		{`bR'abc'`, STRING, py.Bytes(string(`abc`)), ``},
		{`BR"""a\nc"""`, STRING, py.Bytes(string(`a\nc`)), ``},
		{`rB'''a\"c'''`, STRING, py.Bytes(string(`a\"c`)), ``},
	} {
		x, err := NewLex(bytes.NewBufferString(test.in), "<string>", "eval")
		if err != nil {
			t.Fatal(err)
		}
		x.refill()
		token, value := x.readString()
		equal := false
		if valueBytes, ok := value.(py.Bytes); ok {
			if testValueBytes, ok := test.value.(py.Bytes); !ok {
				t.Error("Expecting py.Bytes")
			} else {
				equal = (bytes.Compare(valueBytes, testValueBytes) == 0)
			}
		} else {
			equal = (value == test.value)
		}

		if token != test.token || !equal || x.line != test.out {
			t.Errorf("readString(%q) got (%q,%T,%#v) remainder %q, expected (%q,%T,%#v) remainder %q", test.in, tokenToString[token], value, value, x.line, tokenToString[test.token], test.value, test.value, test.out)
		}
	}
}
