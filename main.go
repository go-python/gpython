// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gpython binary

package main

import (
	"flag"
	"fmt"
	"runtime"
	"runtime/pprof"

	"github.com/go-python/gpython/repl"
	"github.com/go-python/gpython/repl/cli"

	"log"
	"os"

	"github.com/go-python/gpython/py"
)

// Globals
var (
	// Flags
	debug      = flag.Bool("d", false, "Print lots of debugging")
	cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
)

// syntaxError prints the syntax
func syntaxError() {
	fmt.Fprintf(os.Stderr, `GPython

A python implementation in Go

Full options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = syntaxError
	flag.Parse()
	args := flag.Args()
	
	opts := py.DefaultCtxOpts()
	opts.SysArgs = flag.Args()
	ctx := py.NewCtx(opts)
	
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
	
	// IF no args, enter REPL mode
	if len(args) == 0 {

		fmt.Printf("Python 3.4.0 (%s, %s)\n", commit, date)
		fmt.Printf("[Gpython %s]\n", version)
		fmt.Printf("- os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("- go version: %s\n", runtime.Version())

		replCtx := repl.New(ctx)
		cli.RunREPL(replCtx)
		
	} else {
		err := py.RunFile(ctx, args[0], py.CompileOpts{})
		if err != nil {
			py.TracebackDump(err)
			log.Fatal(err)
		}
	}

}
