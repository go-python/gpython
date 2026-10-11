// Copyright 2026 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
)

func TestRunWithNilModulePointer(t *testing.T) {
	for _, entry := range []string{"RunCode", "RunSrc", "RunFile"} {
		t.Run(entry, func(t *testing.T) {
			ctx := py.NewContext(py.DefaultContextOpts())
			defer ctx.Close()
			src := "value = 42\n"
			var module *py.Module
			var err error
			switch entry {
			case "RunCode":
				var code *py.Code
				code, err = py.Compile(src, "<test>", py.ExecMode, 0, true)
				if err == nil {
					module, err = py.RunCode(ctx, code, "<test>", module)
				}
			case "RunSrc":
				module, err = py.RunSrc(ctx, src, "<test>", module)
			case "RunFile":
				path := filepath.Join(t.TempDir(), "test.py")
				if err = os.WriteFile(path, []byte(src), 0600); err == nil {
					module, err = py.RunFile(ctx, filepath.Base(path), py.CompileOpts{CurDir: filepath.Dir(path)}, module)
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			if module == nil || module.Globals["value"] != py.Int(42) {
				t.Fatalf("expected a new module containing value = 42, got %v", module)
			}
			reused, err := py.RunSrc(ctx, "value += 1", "<reuse>", module)
			if err != nil {
				t.Fatal(err)
			}
			if reused != module || module.Globals["value"] != py.Int(43) {
				t.Fatal("expected the existing module and globals to be reused")
			}
		})
	}
}
