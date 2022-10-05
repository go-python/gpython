// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Bytes objects

package py

import (
	"bytes"
	"fmt"
	"strings"
)

var BytesType = ObjectType.NewType("bytes",
	`bytes(iterable_of_ints) -> bytes
bytes(string, encoding[, errors]) -> bytes
bytes(bytes_or_buffer) -> immutable copy of bytes_or_buffer
bytes(int) -> bytes object of size given by the parameter initialized with null bytes
bytes() -> empty bytes object

Construct an immutable array of bytes from:
  - an iterable yielding integers in range(256)
  - a text string encoded using the specified encoding
  - any object implementing the buffer API.
  - an integer`, BytesNew, nil)

type Bytes []byte

// Type of this Bytes object
func (o Bytes) Type() *Type {
	return BytesType
}

// BytesNew
func BytesNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var x Object
	var encoding Object
	var errors Object
	var New Object
	kwlist := []string{"source", "encoding", "errors"}

	err = ParseTupleAndKeywords(args, kwargs, "|Oss:bytes", kwlist, &x, &encoding, &errors)
	if err != nil {
		return nil, err
	}
	if x == nil {
		if encoding != nil || errors != nil {
			return nil, ExceptionNewf(TypeError, "encoding or errors without sequence argument")
		}
		return Bytes{}, nil
	}

	if s, ok := x.(String); ok {
		// Encode via the codec registry
		if encoding == nil {
			return nil, ExceptionNewf(TypeError, "string argument without an encoding")
		}
		encodingStr := strings.ToLower(string(encoding.(String)))
		if encodingStr == "utf-8" || encodingStr == "utf8" {
			return Bytes([]byte(s)), nil
		}
		// FIXME
		// New = PyUnicode_AsEncodedString(x, encoding, errors)
		// assert(PyBytes_Check(New))
		// return New
		return nil, ExceptionNewf(NotImplementedError, "String decode for %q not implemented", encodingStr)
	}

	// We'd like to call PyObject_Bytes here, but we need to check for an
	// integer argument before deferring to PyBytes_FromObject, something
	// PyObject_Bytes doesn't do.
	var ok bool
	if I, ok := x.(I__bytes__); ok {
		New, err = I.M__bytes__()
		if err != nil {
			return nil, err
		}
	} else if New, ok, err = TypeCall0(x, "__bytes__"); ok {
		if err != nil {
			return nil, err
		}
	} else {
		goto no_bytes_method
	}
	if _, ok = New.(Bytes); !ok {
		return nil, ExceptionNewf(TypeError, "__bytes__ returned non-bytes (type %s)", New.Type().Name)
	}
no_bytes_method:

	// Is it an integer?
	_, isInt := x.(Int)
	_, isBigInt := x.(*BigInt)
	if isInt || isBigInt {
		size, err := MakeGoInt(x)
		if err != nil {
			return nil, err
		}
		if size < 0 {
			return nil, ExceptionNewf(ValueError, "negative count")
		}
		return make(Bytes, size), nil
	}

	// If it's not unicode, there can't be encoding or errors
	if encoding != nil || errors != nil {
		return nil, ExceptionNewf(TypeError, "encoding or errors without a string argument")
	}

	return BytesFromObject(x)
}

// Converts an object into bytes
func BytesFromObject(x Object) (Bytes, error) {
	// Look for special cases
	// FIXME implement converting from any object implementing the buffer API.
	switch z := x.(type) {
	case Bytes:
		// Immutable type so just return what was passed in
		return z, nil
	case String:
		return nil, ExceptionNewf(TypeError, "cannot convert unicode object to bytes")
	}
	// Otherwise iterate through the whatever converting it into ints
	b := Bytes{}
	var loopErr error
	iterErr := Iterate(x, func(item Object) bool {
		var value int
		value, loopErr = IndexInt(item)
		if loopErr != nil {
			return true
		}
		if value < 0 || value >= 256 {
			loopErr = ExceptionNewf(ValueError, "bytes must be in range(0, 256)")
			return true
		}
		b = append(b, byte(value))
		return false
	})
	if iterErr != nil {
		return nil, iterErr
	}
	if loopErr != nil {
		return nil, loopErr
	}
	return b, nil
}

func (a Bytes) M__str__() (Object, error) {
	return a.M__repr__()
}

func (a Bytes) M__repr__() (Object, error) {
	// FIXME combine this with parser/stringescape.go into file in py?
	var out bytes.Buffer
	quote := '\''
	if bytes.IndexByte(a, byte('\'')) >= 0 && !(bytes.IndexByte(a, byte('"')) >= 0) {
		quote = '"'
	}
	out.WriteRune('b')
	out.WriteRune(quote)
	for _, c := range a {
		switch {
		case c < 0x20:
			switch c {
			case '\t':
				out.WriteString(`\t`)
			case '\n':
				out.WriteString(`\n`)
			case '\r':
				out.WriteString(`\r`)
			default:
				fmt.Fprintf(&out, `\x%02x`, c)
			}
		case c < 0x7F:
			if c == '\\' || (quote == '\'' && c == '\'') || (quote == '"' && c == '"') {
				out.WriteRune('\\')
			}
			out.WriteByte(c)
		default:
			fmt.Fprintf(&out, "\\x%02x", c)
		}
	}
	out.WriteRune(quote)
	return String(out.String()), nil
}

// Convert an Object to an Bytes
//
// Returns ok as to whether the conversion worked or not
func convertToBytes(other Object) (Bytes, bool) {
	switch b := other.(type) {
	case Bytes:
		return b, true
	}
	return []byte(nil), false
}

// Rich comparison

func (a Bytes) M__lt__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(bytes.Compare(a, b) < 0), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__le__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(bytes.Compare(a, b) <= 0), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__eq__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(bytes.Equal(a, b)), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__ne__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(!bytes.Equal(a, b)), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__gt__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(bytes.Compare(a, b) > 0), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__ge__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		return NewBool(bytes.Compare(a, b) >= 0), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__add__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		o := make([]byte, len(a)+len(b))
		copy(o[:len(a)], a)
		copy(o[len(a):], b)
		return Bytes(o), nil
	}
	return NotImplemented, nil
}

func (a Bytes) M__iadd__(other Object) (Object, error) {
	if b, ok := convertToBytes(other); ok {
		a = append(a, b...)
		return a, nil
	}
	return NotImplemented, nil
}

func (a Bytes) Replace(args Tuple) (Object, error) {
	var (
		pyold Object = None
		pynew Object = None
		pycnt Object = Int(-1)
	)
	err := ParseTuple(args, "yy|i:replace", &pyold, &pynew, &pycnt)
	if err != nil {
		return nil, err
	}

	var (
		old = []byte(pyold.(Bytes))
		new = []byte(pynew.(Bytes))
		cnt = int(pycnt.(Int))
	)

	return Bytes(bytes.Replace([]byte(a), old, new, cnt)), nil
}

// Check interface is satisfied
var (
	_ richComparison = (Bytes)(nil)
	_ I__add__       = (Bytes)(nil)
	_ I__iadd__      = (Bytes)(nil)
)

func init() {
	BytesType.Dict["replace"] = MustNewMethod("replace", func(self Object, args Tuple) (Object, error) {
		return self.(Bytes).Replace(args)
	}, 0, `replace(self, old, new, count=-1) -> return a copy with all occurrences of substring old replaced by new.

  count
    Maximum number of occurrences to replace.
    -1 (the default value) means replace all occurrences.

If the optional argument count is given, only the first count occurrences are
replaced.`)

}
