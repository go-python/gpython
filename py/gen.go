// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"text/template"
)

const filename = "arithmetic.go"

type Ops []struct {
	Name                string
	Title               string
	Operator            string
	TwoReturnParameters bool
	Unary               bool
	Binary              bool
	Ternary             bool
	NoInplace           bool
	Reversed            string
	Conversion          string
	FailReturn          string
}

type Data struct {
	UnaryOps      Ops
	BinaryOps     Ops
	ComparisonOps Ops
}

var data = Data{
	UnaryOps: Ops{
		{Name: "neg", Title: "Neg", Operator: "-", Unary: true},
		{Name: "pos", Title: "Pos", Operator: "+", Unary: true},
		{Name: "abs", Title: "Abs", Operator: "abs", Unary: true},
		{Name: "invert", Title: "Invert", Operator: "~", Unary: true},
		{Name: "complex", Title: "MakeComplex", Operator: "complex", Unary: true, Conversion: "Complex"},
		{Name: "int", Title: "MakeInt", Operator: "int", Unary: true, Conversion: "Int"},
		{Name: "float", Title: "MakeFloat", Operator: "float", Unary: true, Conversion: "Float"},
		{Name: "iter", Title: "Iter", Operator: "iter", Unary: true},
	},
	BinaryOps: Ops{
		{Name: "add", Title: "Add", Operator: "+", Binary: true},
		{Name: "sub", Title: "Sub", Operator: "-", Binary: true},
		{Name: "mul", Title: "Mul", Operator: "*", Binary: true},
		{Name: "truediv", Title: "TrueDiv", Operator: "/", Binary: true},
		{Name: "floordiv", Title: "FloorDiv", Operator: "//", Binary: true},
		{Name: "mod", Title: "Mod", Operator: "%%", Binary: true},
		{Name: "divmod", Title: "DivMod", Operator: "divmod", Binary: true, TwoReturnParameters: true, NoInplace: true},
		{Name: "lshift", Title: "Lshift", Operator: "<<", Binary: true},
		{Name: "rshift", Title: "Rshift", Operator: ">>", Binary: true},
		{Name: "and", Title: "And", Operator: "&", Binary: true},
		{Name: "xor", Title: "Xor", Operator: "^", Binary: true},
		{Name: "or", Title: "Or", Operator: "|", Binary: true},
		{Name: "pow", Title: "Pow", Operator: "** or pow()", Ternary: true},
	},
	ComparisonOps: Ops{
		{Name: "gt", Title: "Gt", Operator: ">", Reversed: "lt"},
		{Name: "ge", Title: "Ge", Operator: ">=", Reversed: "le"},
		{Name: "lt", Title: "Lt", Operator: "<", Reversed: "gt"},
		{Name: "le", Title: "Le", Operator: "<=", Reversed: "ge"},
		{Name: "eq", Title: "Eq", Operator: "==", Reversed: "eq", FailReturn: "False"},
		{Name: "ne", Title: "Ne", Operator: "!=", Reversed: "ne", FailReturn: "True"},
	},
}

func main() {
	t := template.Must(template.New("main").Parse(program))
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to open %q: %v", filename, err)
	}
	if err := t.Execute(out, data); err != nil {
		log.Fatal(err)
	}
	err = out.Close()
	if err != nil {
		log.Fatalf("Failed to close %q: %v", filename, err)
	}
	err = exec.Command("go", "fmt", filename).Run()
	if err != nil {
		log.Fatalf("Failed to gofmt %q: %v", filename, err)
	}
}

var program = `// Automatically generated - DO NOT EDIT
// Regenerate with: go generate

// Arithmetic operations

package py

{{ range .UnaryOps }}
// {{.Title}} the python Object returning an Object
//
// Will raise TypeError if {{.Title}} can't be run on this object
func {{.Title}}(a Object) (Object, error) {
{{ if .Conversion }}
	if _, ok := a.({{.Conversion}}); ok {
		return a, nil
	}
{{end}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res, err := A.M__{{.Name}}__()
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for {{.Operator}}: '%s'", a.Type().Name)
}
{{ end }}

{{ range .BinaryOps }}
// {{.Title}} {{ if .Binary }}two{{ end }}{{ if .Ternary }}three{{ end }} python objects together returning an Object
{{ if .Ternary}}//
// If c != None then it won't attempt to call __r{{.Name}}__
{{ end }}//
// Will raise TypeError if can't be {{.Name}} can't be run on these objects
func {{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object {{ if .TwoReturnParameters}}, Object{{ end }}, error) {
	// Try using a to {{.Name}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res {{ if .TwoReturnParameters}}, res2{{ end }}, err := A.M__{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if err != nil {
			return nil {{ if .TwoReturnParameters }}, nil{{ end }}, err
		}
		if res != NotImplemented {
			return res {{ if .TwoReturnParameters }}, res2{{ end }}, nil
		}
	}

	// Now using b to r{{.Name}} if different in type to a
	if {{ if .Ternary }} c == None && {{ end }} a.Type() != b.Type() {
		if B, ok := b.(I__r{{.Name}}__); ok {
			res {{ if .TwoReturnParameters}}, res2 {{ end }}, err := B.M__r{{.Name}}__(a)
			if err != nil {
				return nil {{ if .TwoReturnParameters }}, nil{{ end }}, err
			}
			if res != NotImplemented {
				return res{{ if .TwoReturnParameters}}, res2{{ end }}, nil
			}
		}
	}
	return nil{{ if .TwoReturnParameters}}, nil{{ end }}, ExceptionNewf(TypeError, "unsupported operand type(s) for {{.Operator}}: '%s' and '%s'", a.Type().Name, b.Type().Name)
}

{{ if not .NoInplace }}
// Inplace {{.Name}}
func I{{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object, error) {
	if A, ok := a.(I__i{{.Name}}__); ok {
		res, err := A.M__i{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}
	return {{.Title}}(a, b {{ if .Ternary }}, c{{ end }})
}
{{end}}
{{end}}

{{ range .ComparisonOps }}
// {{.Title}} two python objects returning a boolean result
//
// Will raise TypeError if {{.Title}} can't be run on this object
func {{.Title}}(a Object, b Object) (Object, error) {
	// Try using a to {{.Name}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res, err := A.M__{{.Name}}__(b)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to {{.Reversed}} with reversed parameters
	if B, ok := b.(I__{{.Reversed}}__); ok {
		res, err := B.M__{{.Reversed}}__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

{{ if .FailReturn}}
if a.Type() != b.Type() {
	return {{ .FailReturn }}, nil
}
{{ end }}
	return nil, ExceptionNewf(TypeError, "unsupported operand type(s) for {{.Operator}}: '%s' and '%s'", a.Type().Name, b.Type().Name)
}
{{ end }}
`
