// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Traceback objects

package py

import (
	"fmt"
	"io"
	"os"
)

// A python Traceback object
type Traceback struct {
	Next   *Traceback
	Frame  *Frame
	Lasti  int32
	Lineno int32
}

var TracebackType = NewType("traceback", "A python traceback")

// Type of this object
func (o *Traceback) Type() *Type {
	return TracebackType
}

// Make a new traceback
func NewTraceback(next *Traceback, frame *Frame, lasti, lineno int32) *Traceback {
	return &Traceback{
		Next:   next,
		Frame:  frame,
		Lasti:  lasti,
		Lineno: lineno,
	}
}

/*
Traceback (most recent call last):
  File "throws.py", line 8, in <module>
    main()
  File "throws.py", line 5, in main
    throws()
  File "throws.py", line 2, in throws
    raise RuntimeError('this is the error message')
RuntimeError: this is the error message
*/

// Dump a traceback for tb to w
func (tb *Traceback) TracebackDump(w io.Writer) {
	for ; tb != nil; tb = tb.Next {
		fmt.Fprintf(w, "  File %q, line %d, in %s\n", tb.Frame.Code.Filename, tb.Lineno, tb.Frame.Code.Name)
		//fmt.Fprintf(w, "    %s\n", "FIXME line of source goes here")
	}
}

// Dumps a traceback to stderr
func TracebackDump(err interface{}) {
	switch e := err.(type) {
	case ExceptionInfo:
		e.TracebackDump(os.Stderr)
	case *ExceptionInfo:
		e.TracebackDump(os.Stderr)
	case *Exception:
		fmt.Fprintf(os.Stderr, "Exception %v\n", e)
		fmt.Fprintf(os.Stderr, "-- No traceback available --\n")
	default:
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		fmt.Fprintf(os.Stderr, "-- No traceback available --\n")
	}
}

// Properties
func init() {
	TracebackType.Dict["__tb_next__"] = &Property{
		Fget: func(self Object) (Object, error) {
			next := self.(*Traceback).Next
			if next == nil {
				return None, nil
			}
			return next, nil
		},
	}
	TracebackType.Dict["__tb_frame__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return self.(*Traceback).Frame, nil
		},
	}
	TracebackType.Dict["__tb_lasti__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return Int(self.(*Traceback).Lasti), nil
		},
	}
	TracebackType.Dict["__tb_lineno__"] = &Property{
		Fget: func(self Object) (Object, error) {
			return Int(self.(*Traceback).Lineno), nil
		},
	}
}

// Make sure it satisfies the interface
var _ Object = (*Traceback)(nil)
