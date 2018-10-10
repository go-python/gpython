// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gpython binary

package main

import (
	"flag"
	"fmt"
	"runtime/pprof"

	_ "github.com/go-python/gpython/builtin"
	"github.com/go-python/gpython/repl/cli"

	//_ "github.com/go-python/gpython/importlib"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/marshal"
	_ "github.com/go-python/gpython/math"
	"github.com/go-python/gpython/py"
	pysys "github.com/go-python/gpython/sys"
	_ "github.com/go-python/gpython/time"
	"github.com/go-python/gpython/vm"
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

// Exit with the message
func fatal(message string, args ...interface{}) {
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	syntaxError()
	fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(1)
}

func main() {
	flag.Usage = syntaxError
	flag.Parse()
	args := flag.Args()
	py.MustGetModule("sys").Globals["argv"] = pysys.MakeArgv(args)
	if len(args) == 0 {
		cli.RunREPL()
		return
	}
	prog := args[0]
	// fmt.Printf("Running %q\n", prog)

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

	// FIXME should be using ImportModuleLevelObject() here
	f, err := os.Open(prog)
	if err != nil {
		log.Fatalf("Failed to open %q: %v", prog, err)
	}
	var obj py.Object
	if strings.HasSuffix(prog, ".pyc") {
		obj, err = marshal.ReadPyc(f)
		if err != nil {
			log.Fatalf("Failed to marshal %q: %v", prog, err)
		}
	} else if strings.HasSuffix(prog, ".py") {
		str, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("Failed to read %q: %v", prog, err)
		}
		obj, err = compile.Compile(string(str), prog, "exec", 0, true)
		if err != nil {
			log.Fatalf("Can't compile %q: %v", prog, err)
		}
	} else {
		log.Fatalf("Can't execute %q", prog)
	}
	if err = f.Close(); err != nil {
		log.Fatalf("Failed to close %q: %v", prog, err)
	}
	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	module.Globals["__file__"] = py.String(prog)
	res, err := vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		log.Fatal(err)
	}
	// fmt.Printf("Return = %v\n", res)
	_ = res

}
