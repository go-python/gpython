// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytest

import (
	"testing"
)

func TestCompileSrc(t *testing.T) {
	for _, tc := range []struct {
		name string
		code string
	}{
		{
			name: "hello",
			code: `print("hello")`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = CompileSrc(t, gContext, tc.code, tc.name)
		})
	}
}

func TestRunTests(t *testing.T) {
	RunTests(t, "./testdata/tests")
}

func TestRunScript(t *testing.T) {
	RunScript(t, "./testdata/hello.py")
}
