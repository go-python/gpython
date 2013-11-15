// Gpython binary

package main

import (
	"flag"
	"fmt"
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
	v := vm.NewVm()
	v.Run(code)

}
