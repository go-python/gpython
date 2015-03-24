package compile

//go:generate ./make_compile_test.py

import (
	"testing"

	"github.com/ncw/gpython/py"
)

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
		if py.Eq(a[i], b[i]) != py.True {
			t.Errorf("%v[%d] has differs, want %#v, got %#v", name, i, a, b)
		}
	}
}

func EqCode(t *testing.T, a, b *py.Code) {
	// int32
	if a.Argcount != b.Argcount {
		t.Errorf("Argcount differs, want %d, got %d", a.Argcount, b.Argcount)
	}
	if a.Kwonlyargcount != b.Kwonlyargcount {
		t.Errorf("Kwonlyargcount differs, want %d, got %d", a.Kwonlyargcount, b.Kwonlyargcount)
	}
	if a.Nlocals != b.Nlocals {
		t.Errorf("Nlocals differs, want %d, got %d", a.Nlocals, b.Nlocals)
	}
	if a.Stacksize != b.Stacksize {
		t.Errorf("Stacksize differs, want %d, got %d", a.Stacksize, b.Stacksize)
	}
	if a.Flags != b.Flags {
		t.Errorf("Flags differs, want %d, got %d", a.Flags, b.Flags)
	}
	if a.Firstlineno != b.Firstlineno {
		t.Errorf("Firstlineno differs, want %d, got %d", a.Firstlineno, b.Firstlineno)
	}

	// string
	if a.Code != b.Code {
		t.Errorf("Code differs, want %q, got %q", a.Code, b.Code)
	}
	if a.Filename != b.Filename {
		t.Errorf("Filename differs, want %q, got %q", a.Filename, b.Filename)
	}
	if a.Name != b.Name {
		t.Errorf("Name differs, want %q, got %q", a.Name, b.Name)
	}
	if a.Lnotab != b.Lnotab {
		t.Errorf("Lnotab differs, want %q, got %q", a.Lnotab, b.Lnotab)
	}

	// Tuple
	EqObjs(t, "Names", a.Consts, b.Consts)

	// []string
	EqStrings(t, "Names", a.Names, b.Names)
	EqStrings(t, "Varnames", a.Varnames, b.Varnames)
	EqStrings(t, "Freevars", a.Freevars, b.Freevars)
	EqStrings(t, "Cellvars", a.Cellvars, b.Cellvars)

	// []byte
	// Cell2arg
}

func TestCompile(t *testing.T) {
	for _, test := range compileTestData {
		codeObj := Compile(test.in, "<string>", test.mode, 0, true)
		code, ok := codeObj.(*py.Code)
		if !ok {
			t.Fatalf("Expecting *py.Code but got %T", codeObj)
		}
		t.Logf("Testing %q", test.in)
		EqCode(t, &test.out, code)
	}
}
