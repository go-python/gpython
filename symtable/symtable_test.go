package symtable

//go:generate ./make_symtable_test.py

import (
	"fmt"
	"testing"

	"github.com/go-python/gpython/parser"
	"github.com/go-python/gpython/py"
)

func EqString(t *testing.T, name string, a, b string) {
	if a != b {
		t.Errorf("%s want %q, got %q", name, a, b)
	}
}

func EqStrings(t *testing.T, name string, a, b []string) {
	if len(a) != len(b) {
		t.Errorf("%s has differing length, want %v, got %v", name, a, b)
		return
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("%s[%d] has differs, want %v, got %v", name, i, a, b)
		}
	}
}

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
	EqScope(t, name+".Scope", a.Scope, b.Scope)
	EqInt(t, name+".Flags", int(a.Flags), int(b.Flags))
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
	for kb, _ := range b {
		if _, ok := a[kb]; ok {
			// Checked already
		} else {
			t.Errorf("%s[%s] extra", name, kb)
		}
	}
}

func EqChildren(t *testing.T, name string, a, b Children) {
	if len(a) != len(b) {
		t.Errorf("%s sizes, want %d got %d", name, len(a), len(b))
		missing := make(map[string]*SymTable)
		extra := make(map[string]*SymTable)
		for _, x := range a {
			missing[x.Name] = x
		}
		for _, x := range b {
			extra[x.Name] = x
		}
		for _, x := range a {
			delete(extra, x.Name)
		}
		for _, x := range b {
			delete(missing, x.Name)
		}
		for _, x := range extra {
			t.Errorf("%s Extra %#v", name, x)
		}
		for _, x := range missing {
			t.Errorf("%s Missing %#v", name, x)
		}
		return
	}
	for i := range a {
		EqSymTable(t, fmt.Sprintf("%s[%d]", name, i), a[i], b[i])
	}
}

func EqSymTable(t *testing.T, name string, a, b *SymTable) {
	EqBlockType(t, name+": Type", a.Type, b.Type)
	EqString(t, name+": Name", a.Name, b.Name)
	// FIXME EqInt(t, name+": Lineno", a.Lineno, b.Lineno)
	EqInt(t, name+": Unoptimized", int(a.Unoptimized), int(b.Unoptimized))
	EqBool(t, name+": Nested", a.Nested, b.Nested)
	EqBool(t, name+": Free", a.Free, b.Free)
	EqBool(t, name+": ChildFree", a.ChildFree, b.ChildFree)
	EqBool(t, name+": Generator", a.Generator, b.Generator)
	EqBool(t, name+": Varargs", a.Varargs, b.Varargs)
	EqBool(t, name+": Varkeywords", a.Varkeywords, b.Varkeywords)
	EqBool(t, name+": ReturnsValue", a.ReturnsValue, b.ReturnsValue)
	EqBool(t, name+": NeedsClassClosure", a.NeedsClassClosure, b.NeedsClassClosure)

	EqSymbols(t, name+": Symbols", a.Symbols, b.Symbols)
	//Global     *SymTable
	//Parent     *SymTable
	EqStrings(t, name+": Varnames", a.Varnames, b.Varnames)
	EqChildren(t, name+": Children", a.Children, b.Children)
}

func TestSymTable(t *testing.T) {
	for _, test := range symtableTestData {
		var symtab *SymTable
		Ast, err := parser.ParseString(test.in, test.mode)
		if err != nil {
			t.Fatalf("Unexpected parse error: %v", err)
		}
		symtab, err = NewSymTable(Ast, "<string>")
		if err != nil {
			if test.exceptionType == nil {
				t.Errorf("%s: Got exception %v when not expecting one", test.in, err)
			} else if exc, ok := err.(*py.Exception); !ok {
				t.Errorf("%s: Got non python exception %T %v", test.in, err, err)
			} else if exc.Type() != test.exceptionType {
				t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
			} else if exc.Type() != test.exceptionType {
				t.Errorf("%s: want exception type %v got %v", test.in, test.exceptionType, exc.Type())
			} else {
				msg := string(exc.Args.(py.Tuple)[0].(py.String))
				if msg != test.errString {
					t.Errorf("%s: want exception text %q got %q", test.in, test.errString, msg)
				}
				if lineno, ok := exc.Dict["lineno"]; ok {
					if lineno.(py.Int) == 0 {
						t.Errorf("%s: lineno not set in exception: %v", test.in, exc.Dict)
					}
				} else {
					t.Errorf("%s: lineno not found in exception: %v", test.in, exc.Dict)
				}
				if filename, ok := exc.Dict["filename"]; ok {
					if filename.(py.String) == py.String("") {
						t.Errorf("%s: filename not set in exception: %v", test.in, exc.Dict)
					}
				} else {
					t.Errorf("%s: filename not found in exception: %v", test.in, exc.Dict)
				}
			}
		} else {
			if test.exceptionType != nil {
				t.Errorf("%s: Didn't get exception %v", test.in, err)
			} else if test.out == nil && symtab != nil {
				t.Errorf("%s: Expecting nil *SymbolTab but got %T", test.in, symtab)
			} else {
				EqSymTable(t, test.in, test.out, symtab)
			}
		}
	}
}

func TestStringer(t *testing.T) {
	EqString(t, "Scope", "ScopeLocal", ScopeLocal.String())
	EqString(t, "Scope", "Scope(100)", Scope(100).String())
	EqString(t, "BlockType", "ClassBlock", ClassBlock.String())
	EqString(t, "BlockType", "BlockType(100)", BlockType(100).String())
}

func TestSymTableFind(t *testing.T) {
	st := &SymTable{
		Symbols: Symbols{
			"x": Symbol{
				Flags: DefLocal | DefNonlocal,
				Scope: ScopeFree,
			},
			"a": Symbol{
				Flags: DefLocal | DefNonlocal,
				Scope: ScopeFree,
			},
			"b": Symbol{
				Flags: DefLocal | DefNonlocal,
				Scope: ScopeFree,
			},
			"c": Symbol{
				Flags: DefNonlocal,
				Scope: ScopeCell,
			},
			"d": Symbol{
				Flags: DefLocal,
				Scope: ScopeCell,
			},
		},
	}

	for _, test := range []struct {
		scope Scope
		flag  DefUseFlags
		want  []string
	}{
		{scope: ScopeGlobalExplicit, flag: 0, want: []string{}},
		{scope: ScopeFree, flag: 0, want: []string{"a", "b", "x"}},
		{scope: ScopeFree, flag: DefLocal, want: []string{"a", "b", "d", "x"}},
		{scope: 0, flag: DefNonlocal, want: []string{"a", "b", "c", "x"}},
	} {
		got := st.Find(test.scope, test.flag)
		EqStrings(t, fmt.Sprintf("Scope %v, Flag %v", test.scope, test.flag), test.want, got)
	}
}
