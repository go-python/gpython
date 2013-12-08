// Gpython binary

package main

import (
	"flag"
	"fmt"
	_ "github.com/ncw/gpython/builtin"
	//_ "github.com/ncw/gpython/importlib"
	"github.com/ncw/gpython/marshal"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/vm"
	"log"
	"os"
	"strings"
)

// Globals
var (
	// Flags
	debug = flag.Bool("d", false, "Print lots of debugging")
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
	if len(args) != 1 {
		fatal("Need program to run")
	}
	prog := args[0]
	fmt.Printf("Running %q\n", prog)

	f, err := os.Open(prog)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	obj, err := marshal.ReadPyc(f)
	if err != nil {
		log.Fatal(err)
	}
	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	res, err := vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		log.Fatal(err)
	}
	fmt.Printf("Return = %v\n", res)

}
