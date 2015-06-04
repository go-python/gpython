package compile

// See FIXME for tests that need to be re-instanted

//go:generate ./make_compile_test.py

import (
	"fmt"
	"io/ioutil"
	"os/exec"
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
			t.Errorf("%s[%d] has differs, want %v, got %v", name, i, a[i], b[i])
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
				eq, err := py.Eq(a[i], b[i])
				if err != nil {
					t.Fatalf("Eq error %v", err)
				}
				equal = eq == py.True
			}
		}
		if !equal {
			t.Errorf("%v[%d] has differs, want %#v, got %#v", name, i, a[i], b[i])
		}
	}
}

func EqCodeCode(t *testing.T, name string, a, b string) {
	if a == b {
		return
	}
	t.Errorf("%s code differs", name)
	want := fmt.Sprintf("%q", a)
	got := fmt.Sprintf("%q", b)
	cmd := exec.Command("./diffdis.py", want, got)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to open pipe: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		t.Errorf("Failed to run ./diffdis.py: %v", err)
		t.Errorf("%s code want %q, got %q", name, a, b)
		return
	}
	stdoutData, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Errorf("./diffdis.py returned error: %v", err)
	}
	t.Error(string(stdoutData))
}

func EqCode(t *testing.T, name string, a, b *py.Code) {
	// int32
	EqInt32(t, name+": Argcount", a.Argcount, b.Argcount)
	EqInt32(t, name+": Kwonlyargcount", a.Kwonlyargcount, b.Kwonlyargcount)
	EqInt32(t, name+": Nlocals", a.Nlocals, b.Nlocals)
	// FIXME EqInt32(t, name+": Stacksize", a.Stacksize, b.Stacksize)
	EqInt32(t, name+": Flags", a.Flags, b.Flags)
	// FIXME EqInt32(t, name+": Firstlineno", a.Firstlineno, b.Firstlineno)

	// string
	EqCodeCode(t, name+": Code", a.Code, b.Code)
	EqString(t, name+": Filename", a.Filename, b.Filename)
	EqString(t, name+": Name", a.Name, b.Name)
	// FIXME EqString(t, name+": Lnotab", a.Lnotab, b.Lnotab)

	// []string
	EqStrings(t, name+": Names", a.Names, b.Names)
	EqStrings(t, name+": Varnames", a.Varnames, b.Varnames)
	EqStrings(t, name+": Freevars", a.Freevars, b.Freevars)
	EqStrings(t, name+": Cellvars", a.Cellvars, b.Cellvars)

	// []byte
	// Cell2arg

	// Tuple
	EqObjs(t, name+": Consts", a.Consts, b.Consts)

}

func TestCompile(t *testing.T) {
	for _, test := range compileTestData {
		// log.Printf(">>> %s", test.in)
		codeObj, err := Compile(test.in, "<string>", test.mode, 0, true)
		if err != nil {
			if test.exceptionType == nil {
				t.Errorf("%s: Got exception %v when not expecting one", test.in, err)
				return
			} else if exc, ok := err.(*py.Exception); !ok {
				t.Errorf("%s: Got non python exception %T %v", test.in, err, err)
				return
			} else if exc.Type() != test.exceptionType {
				t.Errorf("%s: want exception type %v(%s) got %v(%v)", test.in, test.exceptionType, test.errString, exc.Type(), err)
				return
			} else {
				msg := string(exc.Args.(py.Tuple)[0].(py.String))
				if msg != test.errString {
					t.Errorf("%s: want exception text %q got %q", test.in, test.errString, msg)
				}
			}
		} else {
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
}
