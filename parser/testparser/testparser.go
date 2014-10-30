package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ncw/gpython/parser"
)

var (
	lex = flag.Bool("l", false, "Lex the file only")
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		if *lex {
			fmt.Printf("Lexing %q\n", path)
		} else {
			fmt.Printf("Parsing %q\n", path)
		}
		in, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("-----------------\n")
		if *lex {
			parser.Lex(in)
		} else {
			parser.Parse(in)
		}
		fmt.Printf("-----------------\n")
		in.Close()
	}
}
