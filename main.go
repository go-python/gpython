// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gpython binary

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/repl"
	"github.com/go-python/gpython/repl/cli"

	_ "github.com/go-python/gpython/stdlib"
)

var (
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
	xmain(flag.Args())
}

func xmain(args []string) {
	opts := py.DefaultContextOpts()
	opts.SysArgs = args
	ctx := py.NewContext(opts)
	defer ctx.Close()

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
		_, err := py.RunFile(ctx, args[0], py.CompileOpts{}, nil)
		if err != nil {
			if py.IsException(py.SystemExit, err) {
				args := err.(py.ExceptionInfo).Value.(*py.Exception).Args.(py.Tuple)
				if len(args) == 0 {
					os.Exit(0)
				} else if len(args) == 1 {
					if code, ok := args[0].(py.Int); ok {
						c, err := code.GoInt()
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
						os.Exit(c)
					}
					msg, err := py.ReprAsString(args[0])
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					} else {
						fmt.Fprintln(os.Stderr, msg)
					}
					os.Exit(1)
				} else {
					msg, err := py.ReprAsString(args)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					} else {
						fmt.Fprintln(os.Stderr, msg)
					}
					os.Exit(1)
				}
			}
			py.TracebackDump(err)
			os.Exit(1)
		}
	}
}
