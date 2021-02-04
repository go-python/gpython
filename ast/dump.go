// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-python/gpython/py"
)

func dumpItem(v interface{}) string {
	if v == nil {
		return "None"
	}
	switch x := v.(type) {
	case py.String:
		return fmt.Sprintf("'%s'", string(x))
	case py.Bytes:
		return fmt.Sprintf("b'%s'", string(x))
	case Identifier:
		if x == "" {
			return "None"
		}
		return fmt.Sprintf("'%s'", string(x))
	case *Keyword:
		return dump(x, "keyword")
	case *WithItem:
		return dump(x, "withitem")
	case *Arguments:
		if x == nil {
			return "None"
		}
		return dump(x, "arguments")
	case *Arg:
		if x == nil {
			return "None"
		}
		return dump(x, "arg")
	case ModBase:
	case StmtBase:
	case ExprBase:
	case SliceBase:
	case Pos:
	case *Alias:
		return dump(v, "alias")
	case Ast:
		return Dump(x)
	case py.I__str__:
		str, err := x.M__str__()
		if err != nil {
			panic(err)
		}
		return string(str.(py.String))
	case Comprehension:
		return dump(v, "comprehension")
	}
	return fmt.Sprintf("%v", v)
}

// Dump ast as a string with name
func dump(ast interface{}, name string) string {
	if name == "Ellipsis" {
		return name
	}
	astValue := reflect.Indirect(reflect.ValueOf(ast))
	astType := astValue.Type()
	args := make([]string, 0, astType.NumField()+1)
	for i := 0; i < astType.NumField(); i++ {
		fieldType := astType.Field(i)
		fieldValue := astValue.Field(i)
		fname := strings.ToLower(fieldType.Name)
		switch fname {
		case "stmtbase", "exprbase", "modbase", "slicebase", "pos":
			continue
		case "exprtype":
			fname = "type"
		case "contextexpr":
			fname = "context_expr"
		case "optionalvars":
			fname = "optional_vars"
		case "kwdefaults":
			fname = "kw_defaults"
		case "decoratorlist":
			fname = "decorator_list"
		}
		if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() != reflect.Uint8 {
			strs := make([]string, fieldValue.Len())
			for i := 0; i < fieldValue.Len(); i++ {
				element := fieldValue.Index(i)
				if element.CanInterface() {
					v := element.Interface()
					strs[i] = dumpItem(v)
				}
			}
			args = append(args, fmt.Sprintf("%s=[%s]", fname, strings.Join(strs, ", ")))
		} else if fieldValue.CanInterface() {
			v := fieldValue.Interface()
			args = append(args, fmt.Sprintf("%s=%s", fname, dumpItem(v)))
		}
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", "))
}

// Dump an Ast node as a string
func Dump(ast Ast) string {
	if ast == nil {
		return "<nil>"
	}
	name := ast.Type().Name
	switch name {
	case "ExprStmt":
		name = "Expr"
	}
	return dump(ast, name)
}
