// Copyright 2026 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-python/gpython/py"
)

func TestRunFileAbsolutePath(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "script.py")
	if err := os.WriteFile(filename, []byte("value = 42\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, useSysPaths := range []bool{false, true} {
		ctx := NewContext(py.DefaultContextOpts())
		defer ctx.Close()
		// Absolute paths must work even without any search directories.
		ctx.Store().MustGetModule("sys").Globals["path"] = py.NewList()
		module, err := py.RunFile(ctx, filename, py.CompileOpts{UseSysPaths: useSysPaths}, nil)
		if err != nil {
			t.Errorf("UseSysPaths=%v: %v", useSysPaths, err)
			continue
		}
		if module.Globals["value"] != py.Int(42) {
			t.Errorf("UseSysPaths=%v: script was not executed", useSysPaths)
		}
	}
}
