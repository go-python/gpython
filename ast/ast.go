// This file holds the go generate command to build Python-ast.go
// To build it:
//      % go generate
//      % go build

//go:generate python3 asdl_go.py Python.asdl
package ast
