// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/py"
	"github.com/google/go-cmp/cmp"

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

// RunScript runs the provided path to a script.
// RunScript captures the stdout and stderr while executing the script
// and compares it to a golden file, blocking until completion.
//
//	RunScript("./testdata/foo.py")
//
// will compare the output with "./testdata/foo_golden.txt".
func RunScript(t *testing.T, fname string) {

	RunTestTasks(t, []*TestTask{
		{
			PyFile: fname,
		},
	})
}

// RunTestTasks runs each given task in a newly created py.Context concurrently.
// If a fatal error is encountered, the given testing.T is signaled.
func RunTestTasks(t *testing.T, tasks []*TestTask) {
	onCompleted := make(chan *TestTask)

	numTasks := len(tasks)
	for ti := 0; ti < numTasks; ti++ {
		task := tasks[ti]
		go func() {
			err := task.run()
			task.Error = err
			onCompleted <- task
		}()
	}

	tasks = tasks[:0]
	for ti := 0; ti < numTasks; ti++ {
		task := <-onCompleted
		if task.Error != nil {
			t.Error(task.Error)
		}
		tasks = append(tasks, task)
	}
}

var (
	taskCounter int32

	// RegenGoldFiles will cause RunTestTasks() and RunScript() to overwrite their output as the "golden" file (rather than compare against it)
	RegenGoldFiles bool

	// GoldFileSuffix is the default expected suffix for a "golden" output file
	GoldFileSuffix = "_golden.txt"
)

type TestTask struct {
	TaskNum  int32                      // Assigned when this task is run
	TestID   string                     // unique key identifying this task.  If nil, this is autogeneratoed from a PyFile derivative
	PyFile   string                     // If set, this file pathname is executed in a newly created ctx
	PyTask   func(ctx py.Context) error // If set, a new created ctx is created and this blocks until completion
	GoldFile string                     // Filename containing the "gold standard" stdout+stderr.  If nil, autogenerated from PyFile and TestID
	Error    error                      // Non-nil if a fatal error is encountered with this task
}

func (task *TestTask) run() error {
	fileBase := ""

	opts := py.DefaultContextOpts()
	if task.PyFile != "" {
		opts.SysArgs = []string{task.PyFile}
		if task.TestID == "" {
			ext := filepath.Ext(task.PyFile)
			fileBase = task.PyFile[0 : len(task.PyFile)-len(ext)]
		}
	}

	task.TaskNum = atomic.AddInt32(&taskCounter, 1)
	if task.TestID == "" {
		if fileBase == "" {
			task.TestID = fmt.Sprintf("task-%04d", atomic.AddInt32(&taskCounter, 1))
		} else {
			task.TestID = strings.TrimPrefix(fileBase, "./")
			//nameID = strings.ReplaceAll(fileBase, "/", "_")
		}
	}

	if task.GoldFile == "" {
		task.GoldFile = fileBase + GoldFileSuffix
	}

	ctx := py.NewContext(opts)
	defer ctx.Close()

	sys := ctx.Store().MustGetModule("sys")
	tmp, err := os.MkdirTemp("", "gpython-pytest-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	out, err := os.Create(filepath.Join(tmp, "combined"))
	if err != nil {
		return fmt.Errorf("could not create stdout+stderr output file: %w", err)
	}
	defer out.Close()

	sys.Globals["stdout"] = &py.File{File: out, FileMode: py.FileWrite}
	sys.Globals["stderr"] = &py.File{File: out, FileMode: py.FileWrite}

	if task.PyTask != nil {
		err := task.PyTask(ctx)
		if err != nil {
			return err
		}
	}

	if task.PyFile != "" {
		_, err := py.RunFile(ctx, task.PyFile, py.CompileOpts{}, nil)
		if err != nil {
			return fmt.Errorf("could not run target script %q: %+v", task.PyFile, err)
		}
	}

	// Close the ctx explicitly as it may legitimately generate output
	ctx.Close()
	<-ctx.Done()

	err = out.Close()
	if err != nil {
		return fmt.Errorf("could not close output file: %+v", err)
	}

	got, err := os.ReadFile(out.Name())
	if err != nil {
		return fmt.Errorf("could not read script output file: %+v", err)
	}

	if RegenGoldFiles {
		err := os.WriteFile(task.GoldFile, got, 0644)
		if err != nil {
			return fmt.Errorf("could not write golden output %q: %+v", task.GoldFile, err)
		}
	}

	want, err := os.ReadFile(task.GoldFile)
	if err != nil {
		return fmt.Errorf("could not read golden output %q: %+v", task.GoldFile, err)
	}

	diff := cmp.Diff(string(want), string(got))
	if !bytes.Equal(got, want) {
		out := fileBase + ".txt"
		_ = os.WriteFile(out, got, 0644)
		return fmt.Errorf("output differ: -- (-ref +got)\n%s", diff)
	}

	return nil
}
