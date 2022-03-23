// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytest

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/py"

	_ "github.com/go-python/gpython/stdlib"
)

var gContext = py.NewContext(py.DefaultContextOpts())

// Compile the program in the file prog to code in the module that is returned
func compileProgram(t testing.TB, prog string) (*py.Module, *py.Code) {
	f, err := os.Open(prog)
	if err != nil {
		t.Fatalf("%s: Open failed: %v", prog, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatalf("%s: Close failed: %v", prog, err)
		}
	}()

	str, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("%s: ReadAll failed: %v", prog, err)
	}
	return CompileSrc(t, gContext, string(str), prog)
}

func CompileSrc(t testing.TB, ctx py.Context, pySrc string, prog string) (*py.Module, *py.Code) {
	code, err := compile.Compile(string(pySrc), prog, py.ExecMode, 0, true)
	if err != nil {
		t.Fatalf("%s: Compile failed: %v", prog, err)
	}

	module, err := ctx.Store().NewModule(ctx, &py.ModuleImpl{
		Info: py.ModuleInfo{
			FileDesc: prog,
		},
	})
	if err != nil {
		t.Fatalf("%s: NewModule failed: %v", prog, err)
	}

	return module, code
}

// Run the code in the module
func run(t testing.TB, module *py.Module, code *py.Code) {
	_, err := gContext.RunCode(code, module.Globals, module.Globals, nil)
	if err != nil {
		if wantErrObj, ok := module.Globals["err"]; ok {
			gotExc, ok := err.(py.ExceptionInfo)
			if !ok {
				t.Fatalf("got err is not ExceptionInfo: %#v", err)
			}
			if gotExc.Value.Type() != wantErrObj.Type() {
				t.Fatalf("Want exception %v got %v", wantErrObj, gotExc.Value)
			}
			// t.Logf("matched exception")
			return
		} else {
			py.TracebackDump(err)
			t.Fatalf("Run failed: %v at %q", err, module.Globals["doc"])
		}
	}

	// t.Logf("%s: Return = %v", prog, res)
	if doc, ok := module.Globals["doc"]; ok {
		if docStr, ok := doc.(py.String); ok {
			if string(docStr) != "finished" {
				t.Fatalf("Didn't finish at %q", docStr)
			}
		} else {
			t.Fatalf("Set doc variable to non string: %#v", doc)
		}
	} else {
		t.Fatalf("Didn't set doc variable at all")
	}
}

// find the python files in the directory passed in
func findFiles(t testing.TB, testDir string) (names []string) {
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	for _, f := range files {
		name := f.Name()
		if !strings.HasPrefix(name, "lib") && strings.HasSuffix(name, ".py") {
			names = append(names, name)
		}
	}
	return names
}

// RunTests runs the tests in the directory passed in
func RunTests(t *testing.T, testDir string) {
	for _, name := range findFiles(t, testDir) {
		t.Run(name, func(t *testing.T) {
			module, code := compileProgram(t, path.Join(testDir, name))
			run(t, module, code)
		})
	}
}

// RunBenchmarks runs the benchmarks in the directory passed in
func RunBenchmarks(b *testing.B, testDir string) {
	for _, name := range findFiles(b, testDir) {
		module, code := compileProgram(b, path.Join(testDir, name))
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				run(b, module, code)
			}
		})
	}
}
