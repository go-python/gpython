// Copyright 2023 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package array provides the implementation of the python's 'array' module.
package array

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-python/gpython/py"
)

type array struct {
	descr byte // typecode of elements
	esize int  // element size in bytes
	data  any

	append func(v py.Object) (py.Object, error)
	extend func(seq py.Object) (py.Object, error)
}

// Type of this StringDict object
func (*array) Type() *py.Type {
	return ArrayType
}

var (
	_ py.Object       = (*array)(nil)
	_ py.I__getitem__ = (*array)(nil)
	_ py.I__setitem__ = (*array)(nil)
	_ py.I__len__     = (*array)(nil)
	_ py.I__repr__    = (*array)(nil)
	_ py.I__str__     = (*array)(nil)
)

var (
	typecodes = py.String("bBuhHiIlLqQfd")
	ArrayType = py.ObjectType.NewType("array.array", array_doc, array_new, nil)

	descr2esize = map[byte]int{
		'b': 1,
		'B': 1,
		'u': 2,
		'h': 2,
		'H': 2,
		'i': 2,
		'I': 2,
		'l': 8,
		'L': 8,
		'q': 8,
		'Q': 8,
		'f': 4,
		'd': 8,
	}
)

func init() {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "array",
			Doc: "This module defines an object type which can efficiently represent\n" +
				"an array of basic values: characters, integers, floating point\n" +
				"numbers.  Arrays are sequence types and behave very much like lists,\n" +
				"except that the type of objects stored in them is constrained.\n",
		},
		Methods: []*py.Method{},
		Globals: py.StringDict{
			"typecodes": typecodes,
			"array":     ArrayType,
			"ArrayType": ArrayType,
		},
	})

	ArrayType.Dict["itemsize"] = &py.Property{
		Fget: func(self py.Object) (py.Object, error) {
			arr := self.(*array)
			return py.Int(arr.esize), nil
		},
		Doc: "the size, in bytes, of one array item",
	}

	ArrayType.Dict["typecode"] = &py.Property{
		Fget: func(self py.Object) (py.Object, error) {
			arr := self.(*array)
			return py.String(arr.descr), nil
		},
		Doc: "the typecode character used to create the array",
	}

	ArrayType.Dict["append"] = py.MustNewMethod("append", array_append, 0, array_append_doc)
	ArrayType.Dict["extend"] = py.MustNewMethod("extend", array_extend, 0, array_extend_doc)
}

const array_doc = `array(typecode [, initializer]) -> array

Return a new array whose items are restricted by typecode, and
initialized from the optional initializer value, which must be a list,
string or iterable over elements of the appropriate type.

Arrays represent basic values and behave very much like lists, except
the type of objects stored in them is constrained. The type is specified
at object creation time by using a type code, which is a single character.
The following type codes are defined:

    Type code   C Type             Minimum size in bytes
    'b'         signed integer     1
    'B'         unsigned integer   1
    'u'         Unicode character  2 (see note)
    'h'         signed integer     2
    'H'         unsigned integer   2
    'i'         signed integer     2
    'I'         unsigned integer   2
    'l'         signed integer     4
    'L'         unsigned integer   4
    'q'         signed integer     8 (see note)
    'Q'         unsigned integer   8 (see note)
    'f'         floating point     4
    'd'         floating point     8

NOTE: The 'u' typecode corresponds to Python's unicode character. On
narrow builds this is 2-bytes on wide builds this is 4-bytes.

NOTE: The 'q' and 'Q' type codes are only available if the platform
C compiler used to build Python supports 'long long', or, on Windows,
'__int64'.

Methods:

append() -- append a new item to the end of the array
buffer_info() -- return information giving the current memory info
byteswap() -- byteswap all the items of the array
count() -- return number of occurrences of an object
extend() -- extend array by appending multiple elements from an iterable
fromfile() -- read items from a file object
fromlist() -- append items from the list
frombytes() -- append items from the string
index() -- return index of first occurrence of an object
insert() -- insert a new item into the array at a provided position
pop() -- remove and return item (default last)
remove() -- remove first occurrence of an object
reverse() -- reverse the order of the items in the array
tofile() -- write all items to a file object
tolist() -- return the array converted to an ordinary list
tobytes() -- return the array converted to a string

Attributes:

typecode -- the typecode character used to create the array
itemsize -- the length in bytes of one array item

`

func array_new(metatype *py.Type, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	switch n := len(args); n {
	case 0:
		return nil, py.ExceptionNewf(py.TypeError, "array() takes at least 1 argument (0 given)")
	case 1, 2:
		// ok
	default:
		return nil, py.ExceptionNewf(py.TypeError, "array() takes at most 2 arguments (%d given)", n)
	}

	if len(kwargs) != 0 {
		return nil, py.ExceptionNewf(py.TypeError, "array.array() takes no keyword arguments")
	}

	descr, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "array() argument 1 must be a unicode character, not %s", args[0].Type().Name)
	}

	if len(descr) != 1 {
		return nil, py.ExceptionNewf(py.TypeError, "array() argument 1 must be a unicode character, not str")
	}

	if !strings.ContainsAny(string(descr), string(typecodes)) {
		ts := new(strings.Builder)
		for i, v := range typecodes {
			if i > 0 {
				switch {
				case i == len(typecodes)-1:
					ts.WriteString(" or ")
				default:
					ts.WriteString(", ")
				}
			}
			ts.WriteString(string(v))
		}
		return nil, py.ExceptionNewf(py.ValueError, "bad typecode (must be %s)", ts)
	}

	arr := &array{
		descr: descr[0],
		esize: descr2esize[descr[0]],
	}

	if descr[0] == 'u' {
		// FIXME(sbinet)
		return nil, py.NotImplementedError
	}

	switch descr[0] {
	case 'b':
		var data []int8
		arr.data = data
		arr.append = arr.appendI8
		arr.extend = arr.extendI8
	case 'h':
		var data []int16
		arr.data = data
		arr.append = arr.appendI16
		arr.extend = arr.extendI16
	case 'i':
		var data []int32
		arr.data = data
		arr.append = arr.appendI32
		arr.extend = arr.extendI32
	case 'l', 'q':
		var data []int64
		arr.data = data
		arr.append = arr.appendI64
		arr.extend = arr.extendI64
	case 'B':
		var data []uint8
		arr.data = data
		arr.append = arr.appendU8
		arr.extend = arr.extendU8
	case 'H':
		var data []uint16
		arr.data = data
		arr.append = arr.appendU16
		arr.extend = arr.extendU16
	case 'I':
		var data []uint32
		arr.data = data
		arr.append = arr.appendU32
		arr.extend = arr.extendU32
	case 'L', 'Q':
		var data []uint64
		arr.data = data
		arr.append = arr.appendU64
		arr.extend = arr.extendU64
	case 'f':
		var data []float32
		arr.data = data
		arr.append = arr.appendF32
		arr.extend = arr.extendF32
	case 'd':
		var data []float64
		arr.data = data
		arr.append = arr.appendF64
		arr.extend = arr.extendF64
	}

	if len(args) == 2 {
		_, err := arr.extend(args[1])
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

const array_append_doc = `Append new value v to the end of the array.`

func array_append(self py.Object, args py.Tuple) (py.Object, error) {
	arr, ok := self.(*array)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "expected an array, got '%s'", self.Type().Name)
	}
	if len(args) != 1 {
		return nil, py.ExceptionNewf(py.TypeError, "array.append() takes exactly one argument (%d given)", len(args))
	}

	return arr.append(args[0])
}

const array_extend_doc = `Append items to the end of the array.`

func array_extend(self py.Object, args py.Tuple) (py.Object, error) {
	arr, ok := self.(*array)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "expected an array, got '%s'", self.Type().Name)
	}
	if len(args) == 0 {
		return nil, py.ExceptionNewf(py.TypeError, "extend() takes exactly 1 positional argument (%d given)", len(args))
	}
	if len(args) != 1 {
		return nil, py.ExceptionNewf(py.TypeError, "extend() takes at most 1 argument (%d given)", len(args))
	}

	return arr.extend(args[0])
}

func (arr *array) M__repr__() (py.Object, error) {
	o := new(strings.Builder)
	o.WriteString("array('" + string(arr.descr) + "'")
	if data := reflect.ValueOf(arr.data); arr.data != nil && data.Len() > 0 {
		o.WriteString(", [")
		for i := 0; i < data.Len(); i++ {
			if i > 0 {
				o.WriteString(", ")
			}
			fmt.Fprintf(o, "%v", data.Index(i))
		}
		o.WriteString("]")
	}
	o.WriteString(")")
	return py.String(o.String()), nil
}

func (arr *array) M__str__() (py.Object, error) {
	return arr.M__repr__()
}

func (arr *array) M__len__() (py.Object, error) {
	if arr.data == nil {
		return py.Int(0), nil
	}
	sli := reflect.ValueOf(arr.data)
	return py.Int(sli.Len()), nil
}

func (arr *array) M__getitem__(k py.Object) (py.Object, error) {
	switch k := k.(type) {
	case py.Int:
		var (
			sli = reflect.ValueOf(arr.data)
			i   = int(k)
		)
		if i < 0 {
			i = sli.Len() + i
		}
		if i < 0 || sli.Len() <= i {
			return nil, py.ExceptionNewf(py.IndexError, "array index out of range")
		}
		switch arr.descr {
		case 'b', 'h', 'i', 'l', 'q':
			return py.Int(sli.Index(i).Int()), nil
		case 'B', 'H', 'I', 'L', 'Q':
			return py.Int(sli.Index(i).Uint()), nil
		case 'u':
			// FIXME(sbinet)
			return nil, py.NotImplementedError
		case 'f', 'd':
			return py.Float(sli.Index(i).Float()), nil
		}
	case *py.Slice:
		return nil, py.NotImplementedError
	default:
		return nil, py.ExceptionNewf(py.TypeError, "array indices must be integers")
	}
	panic("impossible")
}

func (arr *array) M__setitem__(k, v py.Object) (py.Object, error) {
	switch k := k.(type) {
	case py.Int:
		var (
			sli = reflect.ValueOf(arr.data)
			i   = int(k)
		)
		if i < 0 {
			i = sli.Len() + i
		}
		if i < 0 || sli.Len() <= i {
			return nil, py.ExceptionNewf(py.IndexError, "array index out of range")
		}
		switch arr.descr {
		case 'b', 'h', 'i', 'l', 'q':
			vv := v.(py.Int)
			sli.Index(i).SetInt(int64(vv))
		case 'B', 'H', 'I', 'L', 'Q':
			vv := v.(py.Int)
			sli.Index(i).SetUint(uint64(vv))
		case 'u':
			// FIXME(sbinet)
			return nil, py.NotImplementedError
		case 'f', 'd':
			var vv float64
			switch v := v.(type) {
			case py.Int:
				vv = float64(v)
			case py.Float:
				vv = float64(v)
			default:
				return nil, py.ExceptionNewf(py.TypeError, "must be real number, not %s", v.Type().Name)
			}
			sli.Index(i).SetFloat(vv)
		}
		return py.None, nil
	case *py.Slice:
		return nil, py.NotImplementedError
	default:
		return nil, py.ExceptionNewf(py.TypeError, "array indices must be integers")
	}
	panic("impossible")
}

func (arr *array) appendI8(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]int8), int8(vv))
	return py.None, nil
}

func (arr *array) appendI16(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]int16), int16(vv))
	return py.None, nil
}

func (arr *array) appendI32(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]int32), int32(vv))
	return py.None, nil
}

func (arr *array) appendI64(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]int64), int64(vv))
	return py.None, nil
}

func (arr *array) appendU8(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]uint8), uint8(vv))
	return py.None, nil
}

func (arr *array) appendU16(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]uint16), uint16(vv))
	return py.None, nil
}

func (arr *array) appendU32(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]uint32), uint32(vv))
	return py.None, nil
}

func (arr *array) appendU64(v py.Object) (py.Object, error) {
	vv, err := asInt(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]uint64), uint64(vv))
	return py.None, nil
}

func (arr *array) appendF32(v py.Object) (py.Object, error) {
	vv, err := py.FloatAsFloat64(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]float32), float32(vv))
	return py.None, nil
}

func (arr *array) appendF64(v py.Object) (py.Object, error) {
	vv, err := py.FloatAsFloat64(v)
	if err != nil {
		return nil, err
	}
	arr.data = append(arr.data.([]float64), float64(vv))
	return py.None, nil
}

func (arr *array) extendI8(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendI8(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendI16(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendI16(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendI32(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendI32(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendI64(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendI64(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendU8(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendU8(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendU16(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendU16(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendU32(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendU32(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendU64(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendU64(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendF32(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendF32(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func (arr *array) extendF64(arg py.Object) (py.Object, error) {
	itr, err := py.Iter(arg)
	if err != nil {
		return nil, err
	}

	nxt := itr.(py.I__next__)

	for {
		o, err := nxt.M__next__()
		if err == py.StopIteration {
			break
		}
		_, err = arr.appendF64(o)
		if err != nil {
			return nil, err
		}
	}
	return py.None, nil
}

func asInt(o py.Object) (int64, error) {
	v, ok := o.(py.Int)
	if !ok {
		return 0, py.ExceptionNewf(py.TypeError, "unsupported operand type(s) for int: '%s'", o.Type().Name)
	}
	return int64(v), nil
}
