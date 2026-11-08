// Copyright 2026 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdlib_test

import (
	"sync"
	"testing"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
)

// A registered ModuleImpl is shared by every Context, so importing a CodeSrc
// module from several contexts at once must not race (run with -race).
func TestModuleInitConcurrentCodeSrc(t *testing.T) {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name:     "codesrc_test",
			FileDesc: "<codesrc_test>",
		},
		CodeSrc: "X = 42\ndef get():\n    return X\n",
	})

	const numContexts = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < numContexts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := py.NewContext(py.DefaultContextOpts())
			defer ctx.Close()
			<-start
			err := py.Import(ctx, "codesrc_test")
			if err != nil {
				t.Error(err)
				return
			}
			module, err := ctx.GetModule("codesrc_test")
			if err != nil {
				t.Error(err)
				return
			}
			res, err := module.Call("get", nil, nil)
			if err != nil {
				t.Error(err)
				return
			}
			if res != py.Int(42) {
				t.Errorf("get() = %v, want 42", res)
			}
		}()
	}
	close(start)
	wg.Wait()
}
