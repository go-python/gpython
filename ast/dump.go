package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ncw/gpython/py"
)

// Dump an Ast node as a string
func Dump(ast Ast) string {
	if ast == nil {
		return "<nil>"
	}
	name := ast.Type().Name
	if name == "ExprStmt" {
		name = "Expr"
	}
	astValue := reflect.Indirect(reflect.ValueOf(ast))
	astType := astValue.Type()
	args := make([]string, 0)
	for i := 0; i < astType.NumField(); i++ {
		fieldType := astType.Field(i)
		fieldValue := astValue.Field(i)
		fname := strings.ToLower(fieldType.Name)
		if fieldValue.Kind() == reflect.Slice {
			strs := make([]string, fieldValue.Len())
			for i := 0; i < fieldValue.Len(); i++ {
				element := fieldValue.Index(i)
				if element.CanInterface() {
					if x, ok := element.Interface().(Ast); ok {
						strs[i] = Dump(x)
					} else {
						strs[i] = fmt.Sprintf("%v", element)
					}
				} else {
					strs[i] = fmt.Sprintf("%v", element)
				}
			}
			args = append(args, fmt.Sprintf("%s=[%s]", fname, strings.Join(strs, ", ")))
		} else if fieldValue.CanInterface() {
			v := fieldValue.Interface()
			switch x := v.(type) {
			case py.String:
				args = append(args, fmt.Sprintf("%s=%q", fname, string(x)))
			case ModBase:
			case StmtBase:
			case ExprBase:
			case SliceBase:
			case Pos:
			case Ast:
				args = append(args, fmt.Sprintf("%s=%s", fname, Dump(x)))
			default:
				args = append(args, fmt.Sprintf("%s=%v", fname, x))
			}
		}
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(args, ", "))
}
