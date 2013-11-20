// Implement unmarshal and marshal
package marshal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ncw/gpython/py"
	"io"
	"math/big"
	"strconv"
)

const (
	TYPE_NULL           = '0'
	TYPE_NONE           = 'N'
	TYPE_FALSE          = 'F'
	TYPE_TRUE           = 'T'
	TYPE_STOPITER       = 'S'
	TYPE_ELLIPSIS       = '.'
	TYPE_INT            = 'i'
	TYPE_FLOAT          = 'f'
	TYPE_BINARY_FLOAT   = 'g'
	TYPE_COMPLEX        = 'x'
	TYPE_BINARY_COMPLEX = 'y'
	TYPE_LONG           = 'l'
	TYPE_STRING         = 's'
	TYPE_INTERNED       = 't'
	TYPE_REF            = 'r'
	TYPE_TUPLE          = '('
	TYPE_LIST           = '['
	TYPE_DICT           = '{'
	TYPE_CODE           = 'c'
	TYPE_UNICODE        = 'u'
	TYPE_UNKNOWN        = '?'
	TYPE_SET            = '<'
	TYPE_FROZENSET      = '>'
	FLAG_REF            = 0x80 // with a type, add obj to index
	SIZE32_MAX          = 0x7FFFFFFF

	// We assume that Python ints are stored internally in base some power of
	// 2**15; for the sake of portability we'll always read and write them in base
	// exactly 2**15.

	PyLong_MARSHAL_SHIFT = 15
	PyLong_MARSHAL_BASE  = (1 << PyLong_MARSHAL_SHIFT)
	PyLong_MARSHAL_MASK  = (PyLong_MARSHAL_BASE - 1)
)

// Reads an object from the input
func ReadObject(r io.Reader) (obj py.Object, err error) {
	var code byte
	// defer func() { fmt.Printf("ReadObject(%c) returning %#v with error %v\n", code, obj, err) }()
	err = binary.Read(r, binary.LittleEndian, &code)
	if err != nil {
		return
	}

	//flag := code & FLAG_REF
	Type := code &^ FLAG_REF

	switch Type {
	case TYPE_NULL:
		// A null object
		return nil, nil
	case TYPE_NONE:
		// The Python None object
		return py.None, nil
	case TYPE_FALSE:
		// The python False object
		return py.False, nil
	case TYPE_TRUE:
		// The python True object
		return py.True, nil
	case TYPE_STOPITER:
		// The python StopIteration Exception
		return py.StopIteration, nil
	case TYPE_ELLIPSIS:
		// The python elipsis object
		return py.Elipsis, nil
	case TYPE_INT:
		// 4 bytes of signed integer
		var n int32
		err = binary.Read(r, binary.LittleEndian, &n)
		if err != nil {
			return
		}
		return py.Int(n), nil
	case TYPE_FLOAT:
		// Floating point number as a string
		var length uint8
		err = binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return
		}
		buf := make([]byte, int(length))
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return
		}
		var f float64
		f, err = strconv.ParseFloat(string(buf), 64)
		if err != nil {
			return
		}
		return py.Float(f), nil
	case TYPE_BINARY_FLOAT:
		var f float64
		err = binary.Read(r, binary.LittleEndian, &f)
		if err != nil {
			return
		}
		return py.Float(f), nil
	case TYPE_COMPLEX:
		// Complex number as a string
		// FIXME this is using Go conversion not Python conversion which may differ
		var length uint8
		err = binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return
		}
		buf := make([]byte, int(length))
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return
		}
		var c complex128
		// FIXME c, err = strconv.ParseComplex(string(buf), 64)
		if err != nil {
			return
		}
		return py.Complex(c), nil
	case TYPE_BINARY_COMPLEX:
		var c complex128
		err = binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			return
		}
		return py.Complex(c), nil
	case TYPE_LONG:
		var size int32
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return
		}
		// FIXME negative := false
		if size < 0 {
			// FIXME negative = true
			size = -size
		}
		if size < 0 || size > SIZE32_MAX {
			return nil, errors.New("bad marshal data (long size out of range)")
		}
		// FIXME not sure what -ve size means!
		// Now read shorts which have 15 bits of the number in
		digits := make([]int16, size)
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return
		}
		if digits[0] == 0 {
			// FIXME should be ValueError
			return nil, errors.New("bad marshal data (digit out of range in long)")
		}
		// Convert into a big.Int
		r := new(big.Int)
		t := new(big.Int)
		for _, digit := range digits {
			r.Lsh(r, 15)
			t.SetInt64(int64(digit))
			r.Add(r, t)
		}
		// FIXME try to fit into int64 if possible
		return (*py.BigInt)(r), nil
	case TYPE_STRING, TYPE_INTERNED, TYPE_UNICODE:
		var size int32
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return
		}
		if size < 0 || size > SIZE32_MAX {
			return nil, errors.New("bad marshal data (string size out of range)")
		}
		buf := make([]byte, int(size))
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return
		}
		// FIXME do something different for unicode & interned?
		return py.String(buf), nil
	case TYPE_TUPLE, TYPE_LIST, TYPE_SET, TYPE_FROZENSET:
		var size int32
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return
		}
		if size < 0 || size > SIZE32_MAX {
			return nil, errors.New("bad marshal data (tuple size out of range)")
		}
		tuple := make([]py.Object, int(size))
		for i := range tuple {
			tuple[i], err = ReadObject(r)
			if err != nil {
				return
			}
		}
		switch Type {
		case TYPE_TUPLE:
			return py.Tuple(tuple), nil
		case TYPE_LIST:
			return py.List(tuple), nil
		}

		set := make(py.Set, len(tuple))
		for _, obj := range tuple {
			set[obj] = py.SetValue{}
		}
		switch Type {
		case TYPE_SET:
			return py.Set(set), nil
		case TYPE_FROZENSET:
			return py.FrozenSet(set), nil
		}
	case TYPE_DICT:
		// FIXME should be py.Dict
		dict := py.NewStringDict()
		var key, value py.Object
		for {
			key, err = ReadObject(r)
			if err != nil {
				return
			}
			if key == nil {
				break
			}
			value, err = ReadObject(r)
			if err != nil {
				return
			}
			if value != nil {
				// FIXME should be objects as key
				dict[string(key.(py.String))] = value
			}
		}
		return dict, nil
	case TYPE_REF:
		// Reference to something???
		var n int32
		err = binary.Read(r, binary.LittleEndian, &n)
		if err != nil {
			return
		}
		fmt.Printf("FIXME unimplemented TYPE_REF in unmarshal\n")
		// FIXME
	case TYPE_CODE:
		var argcount int32
		var kwonlyargcount int32
		var nlocals int32
		var stacksize int32
		var flags int32
		var code py.Object
		var consts py.Object
		var names py.Object
		var varnames py.Object
		var freevars py.Object
		var cellvars py.Object
		var filename py.Object
		var name py.Object
		var firstlineno int32
		var lnotab py.Object

		if err = binary.Read(r, binary.LittleEndian, &argcount); err != nil {
			return
		}
		if err = binary.Read(r, binary.LittleEndian, &kwonlyargcount); err != nil {
			return
		}
		if err = binary.Read(r, binary.LittleEndian, &nlocals); err != nil {
			return
		}
		if err = binary.Read(r, binary.LittleEndian, &stacksize); err != nil {
			return
		}
		if err = binary.Read(r, binary.LittleEndian, &flags); err != nil {
			return
		}
		if code, err = ReadObject(r); err != nil {
			return
		}
		if consts, err = ReadObject(r); err != nil {
			return
		}
		if names, err = ReadObject(r); err != nil {
			return
		}
		if varnames, err = ReadObject(r); err != nil {
			return
		}
		if freevars, err = ReadObject(r); err != nil {
			return
		}
		if cellvars, err = ReadObject(r); err != nil {
			return
		}
		if filename, err = ReadObject(r); err != nil {
			return
		}
		if name, err = ReadObject(r); err != nil {
			return
		}
		if err = binary.Read(r, binary.LittleEndian, &firstlineno); err != nil {
			return
		}
		if lnotab, err = ReadObject(r); err != nil {
			return
		}

		// fmt.Printf("argcount = %v\n", argcount)
		// fmt.Printf("kwonlyargcount = %v\n", kwonlyargcount)
		// fmt.Printf("nlocals = %v\n", nlocals)
		// fmt.Printf("stacksize = %v\n", stacksize)
		// fmt.Printf("flags = %v\n", flags)
		// fmt.Printf("code = %x\n", code)
		// fmt.Printf("consts = %v\n", consts)
		// fmt.Printf("names = %v\n", names)
		// fmt.Printf("varnames = %v\n", varnames)
		// fmt.Printf("freevars = %v\n", freevars)
		// fmt.Printf("cellvars = %v\n", cellvars)
		// fmt.Printf("filename = %v\n", filename)
		// fmt.Printf("name = %v\n", name)
		// fmt.Printf("firstlineno = %v\n", firstlineno)
		// fmt.Printf("lnotab = %x\n", lnotab)

		v := py.NewCode(
			argcount, kwonlyargcount,
			nlocals, stacksize, flags,
			code, consts, names, varnames,
			freevars, cellvars, filename, name,
			firstlineno, lnotab)
		return v, nil
	default:
		return nil, errors.New("bad marshal data (unknown type code)")
	}

	return
}

// The header on a .pyc file
type PycHeader struct {
	Magic     uint32
	Timestamp int32
	Length    int32
}

// Reads a pyc file
func ReadPyc(r io.Reader) (obj py.Object, err error) {
	var header PycHeader
	if err = binary.Read(r, binary.LittleEndian, &header); err != nil {
		return
	}
	// FIXME do something with timestamp & lengt?
	if header.Magic>>16 != 0x0a0d {
		return nil, errors.New("Bad magic in .pyc file")
	}
	// fmt.Printf("header = %v\n", header)
	return ReadObject(r)
}
