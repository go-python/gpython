// +build ignore

// Generate arithmetic.go like this
//
// go run gen.go | gofmt >arithmetic.go

package main

import (
	"log"
	"os"
	"text/template"
)

type Ops []struct {
	Name     string
	Title    string
	Operator string
}

type Data struct {
	UnaryOps   Ops
	BinaryOps  Ops
	TrinaryOps Ops
}

var data = Data{
	UnaryOps: Ops{
		{"neg", "Neg", "-"},
		{"pos", "Pos", "+"},
		{"abs", "Abs", "abs"},
		{"invert", "Invert", "invert"},
		{"complex", "MakeComplex", "complex"},
		{"int", "MakeInt", "int"},
		{"float", "MakeFloat", "float"},
		{"index", "Index", "index"},
	},
	BinaryOps: Ops{
		{"add", "Add", "+"},
		{"sub", "Sub", "-"},
		{"mul", "Mul", "*"},
		{"truediv", "TrueDiv", "/"},
		{"floordiv", "FloorDiv", "//"},
		{"mod", "Mod", "%"},
		{"lshift", "Lshift", "<<"},
		{"rshift", "Rshift", ">>"},
		{"and", "And", "&"},
		{"xor", "Xor", "^"},
		{"or", "Or", "|"},
	},
	TrinaryOps: Ops{
		{"pow", "Pow", "**"},
	},
}

func main() {
	t := template.Must(template.New("main").Parse(program))
	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}

var program = `
// Automatically generated - DO NOT EDIT
// Regenerate with: go run gen.go | gofmt >arithmetic.go

// Arithmetic operations

package py

import (
	"fmt"
)

{{ range .BinaryOps }}
// {{.Title}} two python objects together returning an Object
//
// Will raise TypeError if can't be {{.Name}}ed
func {{.Title}}(a, b Object) Object {
	// Try using a to {{.Name}}
	A, ok := a.(I__{{.Name}}__)
	if ok {
		res := A.M__{{.Name}}__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Now using b to r{{.Name}} if different in type to a
	if a.Type() != b.Type() {
		B, ok := b.(I__r{{.Name}}__)
		if ok {
			res := B.M__r{{.Name}}__(a)
			if res != NotImplemented {
				return res
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for {{.Operator}}: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

// Inplace {{.Name}}
func I{{.Title}}(a, b Object) Object {
	A, ok := a.(I__i{{.Name}}__)
	if ok {
		res := A.M__i{{.Name}}__(b)
		if res != NotImplemented {
			return res
		}
	}
	return {{.Title}}(a, b)
}
{{end}}
`
