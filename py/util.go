// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedObjType = errors.New("unsupported obj type")
)

// GetLen is a high-level convenience function that returns the length of the given Object.
func GetLen(obj Object) (Int, error) {
	getlen, ok := obj.(I__len__)
	if !ok {
		return 0, nil
	}

	lenObj, err := getlen.M__len__()
	if err != nil {
		return 0, err
	}

	return GetInt(lenObj)
}

// GetInt is a high-level convenience function that converts the given value to an int.
func GetInt(obj Object) (Int, error) {
	toIdx, ok := obj.(I__index__)
	if !ok {
		_, err := cantConvert(obj, "int")
		return 0, err
	}

	return toIdx.M__index__()
}

// LoadTuple attempts to convert each element of the given list and store into each destination value (based on its type).
func LoadTuple(args Tuple, vars []interface{}) error {

	if len(args) > len(vars) {
		return ExceptionNewf(RuntimeError, "%d args given, expected %d", len(args), len(vars))
	}

	if len(vars) > len(args) {
		vars = vars[:len(args)]
	}

	for i, rval := range vars {
		err := loadValue(args[i], rval)
		if err == ErrUnsupportedObjType {
			return ExceptionNewf(TypeError, "arg %d has unsupported object type: %s", i, args[i].Type().Name)
		}
	}

	return nil
}

// LoadAttr gets the named attribute and attempts to store it into the given destination value (based on its type).
func LoadAttr(obj Object, attrName string, dst interface{}) error {
	attr, err := GetAttrString(obj, attrName)
	if err != nil {
		return err
	}
	err = loadValue(attr, dst)
	if err == ErrUnsupportedObjType {
		return ExceptionNewf(TypeError, "attribute \"%s\" has unsupported object type: %s", attrName, attr.Type().Name)
	}
	return nil
}

// LoadIntsFromList extracts a list of ints contained given a py.List or py.Tuple
func LoadIntsFromList(list Object) ([]int64, error) {
	N, err := GetLen(list)
	if err != nil {
		return nil, err
	}

	getter, ok := list.(I__getitem__)
	if !ok {
		return nil, nil
	}

	if N <= 0 {
		return nil, nil
	}

	intList := make([]int64, N)
	for i := Int(0); i < N; i++ {
		item, err := getter.M__getitem__(i)
		if err != nil {
			return nil, err
		}

		var intVal Int
		intVal, err = GetInt(item)
		if err != nil {
			return nil, err
		}

		intList[i] = int64(intVal)
	}

	return intList, nil
}

func loadValue(src Object, data interface{}) error {
	var (
		v_str   string
		v_float float64
		v_int   int64
	)

	haveStr := false

	switch v := src.(type) {
	case Bool:
		if v {
			v_int = 1
			v_float = 1
			v_str = "True"
		} else {
			v_str = "False"
		}
		haveStr = true
	case Int:
		v_int = int64(v)
		v_float = float64(v)
	case Float:
		v_int = int64(v)
		v_float = float64(v)
	case String:
		v_str = string(v)
		haveStr = true
	case NoneType:
		// No-op
	default:
		return ErrUnsupportedObjType
	}

	switch dst := data.(type) {
	case *Int:
		*dst = Int(v_int)
	case *bool:
		*dst = v_int != 0
	case *int8:
		*dst = int8(v_int)
	case *uint8:
		*dst = uint8(v_int)
	case *int16:
		*dst = int16(v_int)
	case *uint16:
		*dst = uint16(v_int)
	case *int32:
		*dst = int32(v_int)
	case *uint32:
		*dst = uint32(v_int)
	case *int:
		*dst = int(v_int)
	case *uint:
		*dst = uint(v_int)
	case *int64:
		*dst = v_int
	case *uint64:
		*dst = uint64(v_int)
	case *float32:
		if haveStr {
			v_float, _ = strconv.ParseFloat(v_str, 32)
		}
		*dst = float32(v_float)
	case *float64:
		if haveStr {
			v_float, _ = strconv.ParseFloat(v_str, 64)
		}
		*dst = v_float
	case *Float:
		if haveStr {
			v_float, _ = strconv.ParseFloat(v_str, 64)
		}
		*dst = Float(v_float)
	case *string:
		*dst = v_str
	case *String:
		*dst = String(v_str)
	// case []uint64:
	// 	for i := range data {
	// 		dst[i] = order.Uint64(bs[8*i:])
	// 	}
	// case []float32:
	// 	for i := range data {
	// 		dst[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
	// 	}
	// case []float64:
	// 	for i := range data {
	// 		dst[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
	// 	}

	default:
		return ExceptionNewf(NotImplementedError, "%s", "unsupported Go data type")
	}
	return nil
}

// Println prints the provided strings to gpython's stdout.
func Println(self Object, args ...string) bool {
	sysModule, err := self.(*Module).Context.GetModule("sys")
	if err != nil {
		return false
	}
	stdout := sysModule.Globals["stdout"]
	write, err := GetAttrString(stdout, "write")
	if err != nil {
		return false
	}
	call, ok := write.(I__call__)
	if !ok {
		return false
	}
	for _, v := range args {
		if !strings.Contains(v, "\n") {
			v += " "
		}
		_, err := call.M__call__(Tuple{String(v)}, nil)
		if err != nil {
			return false
		}

	}
	_, err = call.M__call__(Tuple{String("\n")}, nil) // newline
	return err == nil
}
