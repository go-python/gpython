// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

//go:generate ./make_grammar_test.py

import (
	"flag"
	"testing"

	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/py"
)

var debugLevel = flag.Int("debugLevel", 0, "Debug level 0-4")

// FIXME test pos is correct

func TestGrammar(t *testing.T) {
	SetDebug(*debugLevel)
	for _, test := range grammarTestData {
		Ast, err := ParseString(test.in, py.CompileMode(test.mode))
		if err != nil {
			if test.exceptionType == nil {
				t.Errorf("%s: Got exception %v when not expecting one", test.in, err)
				return
			} else if exc, ok := err.(*py.Exception); !ok {
				t.Errorf("%s: Got non python exception %T %v", test.in, err, err)
				return
			} else if exc.Type() != test.exceptionType {
				t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
				return
			} else if exc.Type() != test.exceptionType {
				t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
				return
			} else {
				msg := string(exc.Args.(py.Tuple)[0].(py.String))
				if msg != test.errString {
					t.Errorf("%s: want exception text %q got %q", test.in, test.errString, msg)
				}
			}
		} else {
			if test.exceptionType != nil {
				t.Errorf("%s: expecting exception %q", test.in, test.errString)
			} else {
				out := ast.Dump(Ast)
				if out != test.out {
					t.Errorf("Parse(%q)\nwant> %q\n got> %q\n", test.in, test.out, out)
				}
			}
		}
	}
}
