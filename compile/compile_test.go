// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

// See FIXME for tests that need to be re-instanted

//go:generate ./make_compile_test.py

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/go-python/gpython/py"
)

func EqString(t *testing.T, name string, a, b string) {
	if a != b {
		t.Errorf("%s want %q, got %q", name, a, b)
	}
}

// Compact the lnotab as python3 seems to generate inefficient ones
// with spurious zero line number increments.
func lnotabCompact(t *testing.T, pxs *[]byte) {
	xs := *pxs
	newxs := make([]byte, 0, len(xs))
	carry := 0
	for i := 0; i < len(xs); i += 2 {
		d_offset, d_lineno := xs[i], xs[i+1]
		if d_lineno == 0 {
			carry += int(d_offset)
			continue
		}
		// FIXME ignoring d_offset overflow
		d_offset += byte(carry)
		carry = 0
		newxs = append(newxs, byte(d_offset), d_lineno)
	}
	// if string(newxs) != string(xs) {
	// 	t.Logf("Compacted\n% x\n% x", xs, newxs)
	// }
	*pxs = newxs
}

func EqLnotab(t *testing.T, name string, aStr, bStr string) {
	a := []byte(aStr)
	b := []byte(bStr)
	lnotabCompact(t, &a)
	lnotabCompact(t, &b)
	if string(a) != string(b) {
		t.Errorf("%s want % x, got % x", name, a, b)
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
	// The Lnotabs are mostly the same but not entirely
	// So it is probably not profitable to test them exactly
	// EqLnotab(t, name+": Lnotab", a.Lnotab, b.Lnotab)

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
