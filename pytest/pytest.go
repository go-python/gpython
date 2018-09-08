// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytest

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	_ "github.com/go-python/gpython/builtin"
	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/sys"
	"github.com/go-python/gpython/vm"
)

// Run the code in str
func Run(t *testing.T, prog string) {
	f, err := os.Open(prog)
	if err != nil {
		t.Fatalf("%s: Open failed: %v", prog, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatalf("%s: Close failed: %v", prog, err)
		}
	}()

	str, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("%s: ReadAll failed: %v", prog, err)
	}

	obj, err := compile.Compile(string(str), prog, "exec", 0, true)
	if err != nil {
		t.Fatalf("%s: Compile failed: %v", prog, err)
	}

	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	module.Globals["__file__"] = py.String(prog)

	_, err = vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		if wantErr, ok := module.Globals["err"]; ok {
			wantErrObj, ok := wantErr.(py.Object)
			if !ok {
				t.Fatalf("%s: want err is not py.Object: %#v", prog, wantErr)
			}
			gotExc, ok := err.(py.ExceptionInfo)
			if !ok {
				t.Fatalf("%s: got err is not ExceptionInfo: %#v", prog, err)
			}
			if gotExc.Value.Type() != wantErrObj.Type() {
				t.Fatalf("%s: Want exception %v got %v", prog, wantErrObj, gotExc.Value)
			}
			t.Logf("%s: matched exception", prog)
			return
		} else {
			py.TracebackDump(err)
			t.Fatalf("%s: Run failed: %v at %q", prog, err, module.Globals["doc"])
		}
	}

	// t.Logf("%s: Return = %v", prog, res)
	if doc, ok := module.Globals["doc"]; ok {
		if docStr, ok := doc.(py.String); ok {
			if string(docStr) != "finished" {
				t.Fatalf("%s: Didn't finish at %q", prog, docStr)
			}
		} else {
			t.Fatalf("%s: Set doc variable to non string: %#v", prog, doc)
		}
	} else {
		t.Fatalf("%s: Didn't set doc variable at all", prog)
	}
}

// Runs the tests in the directory passed in
func RunTests(t *testing.T, testDir string) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	for _, f := range files {
		name := f.Name()
		if !strings.HasPrefix(name, "lib") && strings.HasSuffix(name, ".py") {
			name := path.Join(testDir, name)
			t.Logf("%s: Running", name)
			Run(t, name)
		}
	}
}
