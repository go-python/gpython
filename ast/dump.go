package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ncw/gpython/py"
)

func dumpItem(v interface{}) string {
	switch x := v.(type) {
	case py.String:
		return fmt.Sprintf("'%s'", string(x))
	case py.Bytes:
		return fmt.Sprintf("b'%s'", string(x))
	case Identifier:
		return fmt.Sprintf("'%s'", string(x))
	case ModBase:
	case StmtBase:
	case ExprBase:
	case SliceBase:
	case Pos:
	case Ast:
		return Dump(x)
	case py.I__str__:
		return string(x.M__str__().(py.String))
	case Comprehension:
		return dump(v, "comprehension")
	}
	return fmt.Sprintf("%v", v)
}

// Dump ast as a string with name
func dump(ast interface{}, name string) string {
	astValue := reflect.Indirect(reflect.ValueOf(ast))
	astType := astValue.Type()
	args := make([]string, 0)
	for i := 0; i < astType.NumField(); i++ {
		fieldType := astType.Field(i)
		fieldValue := astValue.Field(i)
		fname := strings.ToLower(fieldType.Name)
		if fname == "stmtbase" || fname == "exprbase" || fname == "modbase" {
			continue
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
