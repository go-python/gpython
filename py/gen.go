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
	Opposite            string
	Conversion          string
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
	ComparisonOps: Ops{
		{Name: "gt", Title: "Gt", Operator: ">", Opposite: "le"},
		{Name: "ge", Title: "Ge", Operator: ">=", Opposite: "lt"},
		{Name: "lt", Title: "Lt", Operator: "<", Opposite: "ge"},
		{Name: "le", Title: "Le", Operator: "<=", Opposite: "gt"},
		{Name: "eq", Title: "Eq", Operator: "==", Opposite: "ne"},
		{Name: "ne", Title: "Ne", Operator: "!=", Opposite: "eq"},
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

{{ range .UnaryOps }}
// {{.Title}} the python Object returning an Object
//
// Will raise TypeError if {{.Title}} can't be run on this object
func {{.Title}}(a Object) Object {
{{ if .Conversion }}
	if _, ok := a.({{.Conversion}}); ok {
		return a
	}
{{end}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res := A.M__{{.Name}}__()
		if res != NotImplemented {
			return res
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for {{.Operator}}: '%s'", a.Type().Name))
}
{{ end }}

{{ range .BinaryOps }}
// {{.Title}} {{ if .Binary }}two{{ end }}{{ if .Ternary }}three{{ end }} python objects together returning an Object
{{ if .Ternary}}//
// If c != None then it won't attempt to call __r{{.Name}}__
{{ end }}//
// Will raise TypeError if can't be {{.Name}} can't be run on these objects
func {{.Title}}(a, b {{ if .Ternary }}, c{{ end }} Object) (Object {{ if .TwoReturnParameters}}, Object{{ end }}) {
	// Try using a to {{.Name}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res {{ if .TwoReturnParameters}}, res2{{ end }} := A.M__{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if res != NotImplemented {
			return res {{ if .TwoReturnParameters }}, res2{{ end }}
		}
	}

	// Now using b to r{{.Name}} if different in type to a
	if {{ if .Ternary }} c == None && {{ end }} a.Type() != b.Type() {
		if B, ok := b.(I__r{{.Name}}__); ok {
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
	if A, ok := a.(I__i{{.Name}}__); ok {
		res := A.M__i{{.Name}}__(b {{ if .Ternary }}, c{{ end }})
		if res != NotImplemented {
			return res
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
func {{.Title}}(a Object, b Object) Object {
	// Try using a to {{.Name}}
	if A, ok := a.(I__{{.Name}}__); ok {
		res := A.M__{{.Name}}__(b)
		if res != NotImplemented {
			return res
		}
	}

	// Try using b to {{.Opposite}} with reversed parameters
	if B, ok := a.(I__{{.Opposite}}__); ok {
		res := B.M__{{.Opposite}}__(b)
		if res == True {
			return False
		} else if res == False {
			return True
		}
	}

	// FIXME should be TypeError
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for {{.Operator}}: '%s' and '%s'", a.Type().Name, b.Type().Name))
}
{{ end }}
`
