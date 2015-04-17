package compile

//go:generate ./make_symtable_test.py

import (
	"testing"

	"github.com/ncw/gpython/parser"
	"github.com/ncw/gpython/py"
)

func EqInt(t *testing.T, name string, a, b int) {
	if a != b {
		t.Errorf("%s want %v, got %v", name, a, b)
	}
}

func EqScope(t *testing.T, name string, a, b Scope) {
	if a != b {
		t.Errorf("%s want %v, got %v", name, a, b)
	}
}

func EqBlockType(t *testing.T, name string, a, b BlockType) {
	if a != b {
		t.Errorf("%s want %v, got %v", name, a, b)
	}
}

func EqBool(t *testing.T, name string, a, b bool) {
	if a != b {
		t.Errorf("%s want %v, got %v", name, a, b)
	}
}

func EqSymbol(t *testing.T, name string, a, b Symbol) {
	EqString(t, name+".Name", a.Name, b.Name)
	EqScope(t, name+".Scope", a.Scope, b.Scope)
	EqInt(t, name+".Flags", int(a.Flags), int(b.Flags))
	if a.Namespace == nil {
		if b.Namespace == nil {
			// OK
		} else {
			t.Errorf("%s.Namespace want == nil but got != nil", name)
		}
	} else {
		if b.Namespace == nil {
			t.Errorf("%s.Namespace got == nil but want != nil", name)
		} else {
			EqSymTable(t, name+".Namespace", a.Namespace, b.Namespace)
		}
	}
}

func EqSymbols(t *testing.T, name string, a, b Symbols) {
	if len(a) != len(b) {
		t.Errorf("%s sizes, want %d got %d", name, len(a), len(b))
	}
	for ka, va := range a {
		if vb, ok := b[ka]; ok {
			EqSymbol(t, name+"["+ka+"]", va, vb)
		} else {
			t.Errorf("%s[%s] not found", name, ka)
		}
	}
	for kb, vb := range b {
		if va, ok := a[kb]; !ok || vb != va {
			EqSymbol(t, name+"["+kb+"]", va, vb)
		} else {
			t.Errorf("%s[%s] extra", name, kb)
		}
	}
}

func EqSymTable(t *testing.T, name string, a, b *SymTable) {
	EqBlockType(t, name+": Type", a.Type, b.Type)
	EqString(t, name+": Name", a.Name, b.Name)
	EqInt(t, name+": Lineno", a.Lineno, b.Lineno)
	EqBool(t, name+": Optimized", a.Optimized, b.Optimized)
	EqBool(t, name+": Nested", a.Nested, b.Nested)
	EqBool(t, name+": Exec", a.Exec, b.Exec)
	EqBool(t, name+": ImportStar", a.ImportStar, b.ImportStar)
	EqSymbols(t, name+": Symbols", a.Symbols, b.Symbols)
	//Global     *SymTable
	//Parent     *SymTable
	EqStrings(t, name+": Varnames", a.Varnames, b.Varnames)
}

func TestSymTable(t *testing.T) {
	for _, test := range symtableTestData {
		var symtab *SymTable
		func() {
			defer func() {
				return
				if r := recover(); r != nil {
					if test.exceptionType == nil {
						t.Errorf("%s: Got exception %v when not expecting one", test.in, r)
						return
					}
					exc, ok := r.(*py.Exception)
					if !ok {
						t.Errorf("%s: Got non python exception %T %v", test.in, r, r)
						return
					}
					if exc.Type() != test.exceptionType {
						t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
						return
					}
					if exc.Type() != test.exceptionType {
						t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
						return
					}
					msg := string(exc.Args.(py.Tuple)[0].(py.String))
					if msg != test.errString {
						t.Errorf("%s: want exception text %q got %q", test.in, test.errString, msg)
					}

				}
			}()
			Ast, err := parser.ParseString(test.in, test.mode)
			if err != nil {
				panic(err) // FIXME error handling!
			}
			symtab = NewSymTable(Ast)
		}()
		if test.out == nil {
			if symtab != nil {
				t.Errorf("%s: Expecting nil *py.Code but got %T", test.in, symtab)
			}
		} else {
			EqSymTable(t, test.in, test.out, symtab)
		}
	}
}
