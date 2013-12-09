// Implement unmarshal and marshal
package marshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/vm"
	"io"
	"math/big"
	"strconv"
)

const (
	MARSHAL_VERSION     = 2
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
		return py.Ellipsis, nil
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
		err = binary.Read(r, binary.LittleEndian, &digits)
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
			return py.NewListFromItems(tuple), nil
		}

		switch Type {
		case TYPE_SET:
			return py.NewSetFromItems(tuple), nil
		case TYPE_FROZENSET:
			return py.NewFrozenSetFromItems(tuple), nil
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

// Unmarshals a frozen module
func LoadFrozenModule(name string, data []byte) *py.Module {
	r := bytes.NewBuffer(data)
	obj, err := ReadObject(r)
	if err != nil {
		panic(err)
	}
	code := obj.(*py.Code)
	module := py.NewModule(name, "", nil, nil)
	_, err = vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		panic(err)
	}
	return module
}

const dump_doc = `dump(value, file[, version])

Write the value on the open file. The value must be a supported type.
The file must be an open file object such as sys.stdout or returned by
open() or os.popen(). It must be opened in binary mode ('wb' or 'w+b').

If the value has (or contains an object that has) an unsupported type, a
ValueError exception is raised — but garbage data will also be written
to the file. The object will not be properly read back by load()

The version argument indicates the data format that dump should use.`

func marshal_dump(self py.Object, args py.Tuple) py.Object {
	/*
	   // XXX Quick hack -- need to do this differently
	   PyObject *x;
	   PyObject *f;
	   int version = Py_MARSHAL_VERSION;
	   PyObject *s;
	   PyObject *res;
	   _Py_IDENTIFIER(write);

	   if (!PyArg_ParseTuple(args, "OO|i:dump", &x, &f, &version))
	       return NULL;
	   s = PyMarshal_WriteObjectToString(x, version);
	   if (s == NULL)
	       return NULL;
	   res = _PyObject_CallMethodId(f, &PyId_write, "O", s);
	   Py_DECREF(s);
	   return res;
	*/
	panic("dump not implemented")
}

const load_doc = `load(file)

Read one value from the open file and return it. If no valid value is
read (e.g. because the data has a different Python version’s
incompatible marshal format), raise EOFError, ValueError or TypeError.
The file must be an open file object opened in binary mode ('rb' or
'r+b').

Note: If an object containing an unsupported type was marshalled with
dump(), load() will substitute None for the unmarshallable type.`

func marshal_load(self, f py.Object) py.Object {
	/*
	   PyObject *data, *result;
	   _Py_IDENTIFIER(read);
	   RFILE rf;

	    // Make a call to the read method, but read zero bytes.
	    // This is to ensure that the object passed in at least
	    // has a read method which returns bytes.
	   data = _PyObject_CallMethodId(f, &PyId_read, "i", 0);
	   if (data == NULL)
	       return NULL;
	   if (!PyBytes_Check(data)) {
	       PyErr_Format(PyExc_TypeError,
	                    "f.read() returned not bytes but %.100s",
	                    data->ob_type->tp_name);
	       result = NULL;
	   }
	   else {
	       rf.depth = 0;
	       rf.fp = NULL;
	       rf.readable = f;
	       rf.current_filename = NULL;
	       result = read_object(&rf);
	   }
	   Py_DECREF(data);
	   return result;
	*/
	panic("load not implemented")
}

const dumps_doc = `dumps(value[, version])

Return the string that would be written to a file by dump(value, file).
The value must be a supported type. Raise a ValueError exception if
value has (or contains an object that has) an unsupported type.

The version argument indicates the data format that dumps should use.`

func marshal_dumps(self py.Object, args py.Tuple) py.Object {
	/*
	   PyObject *x;
	   int version = Py_MARSHAL_VERSION;
	   if (!PyArg_ParseTuple(args, "O|i:dumps", &x, &version))
	       return NULL;
	   return PyMarshal_WriteObjectToString(x, version);
	*/
	panic("dumps not implemented")
}

const loads_doc = `loads(bytes)

Convert the bytes object to a value. If no valid value is found, raise
EOFError, ValueError or TypeError. Extra characters in the input are
ignored.`

func marshal_loads(self py.Object, args py.Tuple) py.Object {
	/*
	   RFILE rf;
	   Py_buffer p;
	   char *s;
	   Py_ssize_t n;
	   PyObject* result;
	   if (!PyArg_ParseTuple(args, "y*:loads", &p))
	       return NULL;
	   s = p.buf;
	   n = p.len;
	   rf.fp = NULL;
	   rf.readable = NULL;
	   rf.current_filename = NULL;
	   rf.ptr = s;
	   rf.end = s + n;
	   rf.depth = 0;
	   result = read_object(&rf);
	   PyBuffer_Release(&p);
	   return result;
	*/
	panic("loads not implemented")
}

const module_doc = `This module contains functions that can read and write Python values in
a binary format. The format is specific to Python, but independent of
machine architecture issues.

Not all Python object types are supported; in general, only objects
whose value is independent from a particular invocation of Python can be
written and read by this module. The following types are supported:
None, integers, floating point numbers, strings, bytes, bytearrays,
tuples, lists, sets, dictionaries, and code objects, where it
should be understood that tuples, lists and dictionaries are only
supported as long as the values contained therein are themselves
supported; and recursive lists and dictionaries should not be written
(they will cause infinite loops).

Variables:

version -- indicates the format that the module uses. Version 0 is the
    historical format, version 1 shares interned strings and version 2
    uses a binary format for floating point numbers.

Functions:

dump() -- write value to a file
load() -- read value from a file
dumps() -- write value to a string
loads() -- read value from a string`

// Initialise the module
func init() {
	methods := []*py.Method{
		py.NewMethod("dump", marshal_dump, 0, dump_doc),
		py.NewMethod("load", marshal_load, 0, load_doc),
		py.NewMethod("dumps", marshal_dumps, 0, dumps_doc),
		py.NewMethod("loads", marshal_loads, 0, loads_doc),
	}
	globals := py.StringDict{
		"version": py.Int(MARSHAL_VERSION),
	}
	py.NewModule("marshal", module_doc, methods, globals)
}
