// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	// This initializes gpython for runtime execution and is essential.
	// It defines forward-declared symbols and registers native built-in modules, such as sys and time.
	_ "github.com/go-python/gpython/modules"

	// Commonly consumed gpython
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/repl"
	"github.com/go-python/gpython/repl/cli"
)

func main() {
	flag.Parse()
	runWithFile(flag.Arg(0))
}

func runWithFile(pyFile string) error {

	// See type Context interface and related docs
	ctx := py.NewContext(py.DefaultContextOpts())
	
	// This drives modules being able to perform cleanup and release resources 
	defer ctx.Close()

	var err error
	if len(pyFile) == 0 {
		replCtx := repl.New(ctx)

		fmt.Print("\n=======  Entering REPL mode, press Ctrl+D to exit  =======\n")

		_, err = py.RunFile(ctx, "lib/REPL-startup.py", py.CompileOpts{}, replCtx.Module)
		if err == nil {
			cli.RunREPL(replCtx)
		}

	} else {
		_, err = py.RunFile(ctx, pyFile, py.CompileOpts{}, nil)
	}

	if err != nil {
		py.TracebackDump(err)
	}

	return err
}
