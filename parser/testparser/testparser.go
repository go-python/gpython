package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ncw/gpython/parser"
)

var (
	lex        = flag.Bool("l", false, "Lex the file only")
	debugLevel = flag.Int("d", 0, "Debug level 0-4")
)

func main() {
	flag.Parse()
	parser.SetDebug(*debugLevel)
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
			_, err = parser.Lex(in)
		} else {
			err = parser.Parse(in)
		}
		fmt.Printf("-----------------\n")
		in.Close()
		if err != nil {
			log.Fatalf("Failed on %q: %v", path, err)
		}
	}
}
