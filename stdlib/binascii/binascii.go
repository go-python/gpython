// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binascii provides the implementation of the python's 'binascii' module.
package binascii

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash/crc32"

	"github.com/go-python/gpython/py"
)

var (
	Incomplete = py.ExceptionType.NewType("binascii.Incomplete", "", nil, nil)
	Error      = py.ValueError.NewType("binascii.Error", "", nil, nil)
)

func init() {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "binascii",
			Doc:  "Conversion between binary data and ASCII",
		},
		Methods: []*py.Method{
			py.MustNewMethod("a2b_base64", a2b_base64, 0, "Decode a line of base64 data."),
			py.MustNewMethod("b2a_base64", b2a_base64, 0, "Base64-code line of data."),
			py.MustNewMethod("a2b_hex", a2b_hex, 0, a2b_hex_doc),
			py.MustNewMethod("b2a_hex", b2a_hex, 0, b2a_hex_doc),
			py.MustNewMethod("crc32", crc32_, 0, "Compute CRC-32 incrementally."),
			py.MustNewMethod("unhexlify", a2b_hex, 0, unhexlify_doc),
			py.MustNewMethod("hexlify", b2a_hex, 0, hexlify_doc),
		},
		Globals: py.StringDict{
			"Incomplete": Incomplete,
			"Error":      Error,
		},
	})
}

func b2a_base64(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pydata py.Object
		pynewl py.Object = py.True
	)
	err := py.ParseTupleAndKeywords(args, kwargs, "y*|p:binascii.b2a_base64", []string{"data", "newline"}, &pydata, &pynewl)
	if err != nil {
		return nil, err
	}

	var (
		buf     = []byte(pydata.(py.Bytes))
		newline = bool(pynewl.(py.Bool))
	)

	out := base64.StdEncoding.EncodeToString(buf)
	if newline {
		out += "\n"
	}
	return py.Bytes(out), nil
}

func a2b_base64(self py.Object, args py.Tuple) (py.Object, error) {
	var pydata py.Object
	err := py.ParseTuple(args, "s:binascii.a2b_base64", &pydata)
	if err != nil {
		return nil, err
	}

	out, err := base64.StdEncoding.DecodeString(string(pydata.(py.String)))
	if err != nil {
		return nil, py.ExceptionNewf(Error, "could not decode base64 data: %+v", err)
	}

	return py.Bytes(out), nil
}

func crc32_(self py.Object, args py.Tuple) (py.Object, error) {
	var (
		pydata py.Object
		pycrc  py.Object = py.Int(0)
	)

	err := py.ParseTuple(args, "y*|i:binascii.crc32", &pydata, &pycrc)
	if err != nil {
		return nil, err
	}

	crc := crc32.Update(uint32(pycrc.(py.Int)), crc32.IEEETable, []byte(pydata.(py.Bytes)))
	return py.Int(crc), nil

}

const a2b_hex_doc = `Binary data of hexadecimal representation.

hexstr must contain an even number of hex digits (upper or lower case).
This function is also available as "unhexlify()".`

func a2b_hex(self py.Object, args py.Tuple) (py.Object, error) {
	var (
		hexErr hex.InvalidByteError
		pydata py.Object
		src    string
	)
	err := py.ParseTuple(args, "s*:binascii.a2b_hex", &pydata)
	if err != nil {
		return nil, err
	}

	switch v := pydata.(type) {
	case py.String:
		src = string(v)
	case py.Bytes:
		src = string(v)
	}

	o, err := hex.DecodeString(src)
	if err != nil {
		switch {
		case errors.Is(err, hex.ErrLength):
			return nil, py.ExceptionNewf(Error, "Odd-length string")
		case errors.As(err, &hexErr):
			return nil, py.ExceptionNewf(Error, "Non-hexadecimal digit found")
		default:
			return nil, py.ExceptionNewf(Error, "could not decode hex data: %+v", err)
		}
	}

	return py.Bytes(o), nil
}

const b2a_hex_doc = `Hexadecimal representation of binary data.

The return value is a bytes object.  This function is also
available as "hexlify()".`

func b2a_hex(self py.Object, args py.Tuple) (py.Object, error) {
	var pydata py.Object
	err := py.ParseTuple(args, "y*:binascii.b2a_hex", &pydata)
	if err != nil {
		return nil, err
	}

	o := hex.EncodeToString([]byte(pydata.(py.Bytes)))
	return py.Bytes(o), nil
}

const unhexlify_doc = `Binary data of hexadecimal representation.

hexstr must contain an even number of hex digits (upper or lower case).`

const hexlify_doc = `Hexadecimal representation of binary data.

The return value is a bytes object.  This function is also
available as "b2a_hex()".`
