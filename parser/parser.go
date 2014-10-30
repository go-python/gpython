// This file holds the go generate command to run yacc on the grammar in grammar.y.
// To build y.go:
//      % go generate
//      % go build

//go:generate go tool yacc -v y.output grammar.y
package parser
