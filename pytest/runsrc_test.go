// Copyright 2026 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytest

import (
	"testing"

	"github.com/go-python/gpython/py"

	_ "github.com/go-python/gpython/stdlib"
)

// Regression test for https://github.com/go-python/gpython/issues/257
//
// RunSrc used to compile in SingleMode which only accepts a single
// interactive statement, so source containing more than one top-level
// statement (e.g. two class definitions) failed with a SyntaxError.
func TestRunSrcMultipleStatements(t *testing.T) {
	for _, test := range []struct {
		name string
		src  string
	}{
		{
			name: "two classes",
			src: `
class A:
    pass


class B:
    pass
`,
		},
		{
			name: "two functions",
			src: `
def a():
    return 1

def b():
    return 2
`,
		},
		{
			name: "assignments and expression",
			src: `
x = 1
y = 2
z = x + y
`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := py.NewContext(py.DefaultContextOpts())
			defer ctx.Close()
			_, err := py.RunSrc(ctx, test.src, "run", nil)
			if err != nil {
				t.Fatalf("RunSrc failed: %v", err)
			}
		})
	}
}
