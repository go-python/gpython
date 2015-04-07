package compile

//go:generate ./make_compile_test.py

import (
	"testing"

	"github.com/ncw/gpython/py"
)

func EqString(t *testing.T, name string, a, b string) {
	if a != b {
		t.Errorf("%s want %q, got %q", name, a, b)
	}
}

func EqInt32(t *testing.T, name string, a, b int32) {
	if a != b {
		t.Errorf("%s want %d, got %d", name, a, b)
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

func EqObjs(t *testing.T, name string, a, b []py.Object) {
	if len(a) != len(b) {
		t.Errorf("%s has differing length, want %v, got %v", name, a, b)
		return
	}
	for i := range a {
		equal := a[i].Type() == b[i].Type()
		if equal {
			if a[i].Type() == py.CodeType {
				A := a[i].(*py.Code)
				B := b[i].(*py.Code)
				EqCode(t, name, A, B)
			} else {
				equal = py.Eq(a[i], b[i]) == py.True
			}
		}
		if !equal {
			t.Errorf("%v[%d] has differs, want %#v, got %#v", name, i, a, b)
		}
	}
}

func EqCode(t *testing.T, name string, a, b *py.Code) {
	// int32
	EqInt32(t, name+": Argcount", a.Argcount, b.Argcount)
	EqInt32(t, name+": Kwonlyargcount", a.Kwonlyargcount, b.Kwonlyargcount)
	EqInt32(t, name+": Nlocals", a.Nlocals, b.Nlocals)
	EqInt32(t, name+": Stacksize", a.Stacksize, b.Stacksize)
	EqInt32(t, name+": Flags", a.Flags, b.Flags)
	EqInt32(t, name+": Firstlineno", a.Firstlineno, b.Firstlineno)

	// string
	EqString(t, name+": Code", a.Code, b.Code)
	EqString(t, name+": Filename", a.Filename, b.Filename)
	EqString(t, name+": Name", a.Name, b.Name)
	EqString(t, name+": Lnotab", a.Lnotab, b.Lnotab)

	// Tuple
	EqObjs(t, name+": Consts", a.Consts, b.Consts)

	// []string
	EqStrings(t, name+": Names", a.Names, b.Names)
	EqStrings(t, name+": Varnames", a.Varnames, b.Varnames)
	EqStrings(t, name+": Freevars", a.Freevars, b.Freevars)
	EqStrings(t, name+": Cellvars", a.Cellvars, b.Cellvars)

	// []byte
	// Cell2arg
}

func TestCompile(t *testing.T) {
	for _, test := range compileTestData {
		var codeObj py.Object
		func() {
			defer func() {
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
			codeObj = Compile(test.in, "<string>", test.mode, 0, true)
		}()
		if test.out == nil {
			if codeObj != nil {
				t.Errorf("%s: Expecting nil *py.Code but got %T", test.in, codeObj)
			}
		} else {
			code, ok := codeObj.(*py.Code)
			if !ok {
				t.Errorf("%s: Expecting *py.Code but got %T", test.in, codeObj)
			} else {
				//t.Logf("Testing %q", test.in)
				EqCode(t, test.in, test.out, code)
			}
		}
	}
}
