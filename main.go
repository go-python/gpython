// Gpython binary

package main

import (
	"flag"
	"fmt"

	_ "github.com/ncw/gpython/builtin"
	"github.com/ncw/gpython/repl"
	//_ "github.com/ncw/gpython/importlib"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ncw/gpython/compile"
	"github.com/ncw/gpython/marshal"
	"github.com/ncw/gpython/py"
	_ "github.com/ncw/gpython/sys"
	_ "github.com/ncw/gpython/time"
	"github.com/ncw/gpython/vm"
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
	if len(args) == 0 {
		repl.Run()
		return
	}
	prog := args[0]
	fmt.Printf("Running %q\n", prog)

	// FIXME should be using ImportModuleLevelObject() here
	f, err := os.Open(prog)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // FIXME don't leave f open for the whole program!
	var obj py.Object
	if strings.HasSuffix(prog, ".pyc") {
		obj, err = marshal.ReadPyc(f)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasSuffix(prog, ".py") {
		str, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		obj, err = compile.Compile(string(str), prog, "exec", 0, true)
		if err != nil {
			log.Fatalf("Can't compile %q: %v", prog, err)
		}
	} else {
		log.Fatalf("Can't execute %q", prog)
	}
	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	module.Globals["__file__"] = py.String(prog)
	res, err := vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		log.Fatal(err)
	}
	fmt.Printf("Return = %v\n", res)

}
