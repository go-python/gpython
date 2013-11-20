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
	Name                string
	Title               string
	Operator            string
	TwoReturnParameters bool
	Unary               bool
	Binary              bool
	Ternary             bool
	NoInplace           bool
}

type Data struct {
	UnaryOps   Ops
	BinaryOps  Ops
	TernaryOps Ops
}

var data = Data{
	UnaryOps: Ops{
		{Name: "neg", Title: "Neg", Operator: "-", Unary: true},
		{Name: "pos", Title: "Pos", Operator: "+", Unary: true},
		{Name: "abs", Title: "Abs", Operator: "abs", Unary: true},
		{Name: "invert", Title: "Invert", Operator: "invert", Unary: true},
		{Name: "complex", Title: "MakeComplex", Operator: "complex", Unary: true},
		{Name: "int", Title: "MakeInt", Operator: "int", Unary: true},
		{Name: "float", Title: "MakeFloat", Operator: "float", Unary: true},
		{Name: "index", Title: "Index", Operator: "index", Unary: true},
	},
	BinaryOps: Ops{
		{Name: "add", Title: "Add", Operator: "+", Binary: true},
		{Name: "sub", Title: "Sub", Operator: "-", Binary: true},
		{Name: "mul", Title: "Mul", Operator: "*", Binary: true},
		{Name: "truediv", Title: "TrueDiv", Operator: "/", Binary: true},
		{Name: "floordiv", Title: "FloorDiv", Operator: "//", Binary: true},
		{Name: "mod", Title: "Mod", Operator: "%", Binary: true},
		{Name: "divmod", Title: "DivMod", Operator: "divmod", Binary: true, TwoReturnParameters: true, NoInplace: true},
		{Name: "lshift", Title: "Lshift", Operator: "<<", Binary: true},
		{Name: "rshift", Title: "Rshift", Operator: ">>", Binary: true},
		{Name: "and", Title: "And", Operator: "&", Binary: true},
		{Name: "xor", Title: "Xor", Operator: "^", Binary: true},
		{Name: "or", Title: "Or", Operator: "|", Binary: true},
		{Name: "pow", Title: "Pow", Operator: "** or pow()", Ternary: true},
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
// {{.Title}} {{ if .Binary }}two{{ end }}{{ if .Ternary }}three{{ end }} python objects together returning an Object
{{ if .Ternary}}//
// If c != None then it won't attempt to call __r{{.Name}}__
{{ end }}//
// Will raise TypeError if can't be {{.Name}} can't be run on these objects
func {{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object {{ if .TwoReturnParameters}}, Object{{ end }}) {
	// Try using a to {{.Name}}
	A, ok := a.(I__{{.Name}}__)
	if ok {
		res {{ if .TwoReturnParameters}}, res2{{ end }} := A.M__{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if res != NotImplemented {
			return res {{ if .TwoReturnParameters }}, res2{{ end }}
		}
	}

	// Now using b to r{{.Name}} if different in type to a
	if {{ if .Ternary }} c == None && {{ end }} a.Type() != b.Type() {
		B, ok := b.(I__r{{.Name}}__)
		if ok {
			res {{ if .TwoReturnParameters}}, res2 {{ end }} := B.M__r{{.Name}}__(a)
			if res != NotImplemented {
				return res{{ if .TwoReturnParameters}}, res2{{ end }}
			}
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for {{.Operator}}: '%s' and '%s'", a.Type().Name, b.Type().Name))
}

{{ if not .NoInplace }}
// Inplace {{.Name}}
func I{{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) Object {
	A, ok := a.(I__i{{.Name}}__)
	if ok {
		res := A.M__i{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if res != NotImplemented {
			return res
		}
	}
	return {{.Title}}(a, b {{ if .Ternary }}, c{{ end }})
}
{{end}}
{{end}}
`
