// Bytes objects

package py

import (
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
func BytesNew(metatype *Type, args Tuple, kwargs StringDict) (res Object) {
	var x Object
	var encoding Object
	var errors Object
	var New Object
	kwlist := []string{"source", "encoding", "errors"}

	ParseTupleAndKeywords(args, kwargs, "|Oss:bytes", kwlist, &x, &encoding, &errors)
	if x == nil {
		if encoding != nil || errors != nil {
			panic(ExceptionNewf(TypeError, "encoding or errors without sequence argument"))
		}
		return Bytes{}
	}

	if s, ok := x.(String); ok {
		// Encode via the codec registry
		if encoding == nil {
			panic(ExceptionNewf(TypeError, "string argument without an encoding"))
		}
		encodingStr := strings.ToLower(string(encoding.(String)))
		if encodingStr == "utf-8" || encodingStr == "utf8" {
			return Bytes([]byte(s))
		}
		// FIXME
		// New = PyUnicode_AsEncodedString(x, encoding, errors)
		// assert(PyBytes_Check(New))
		// return New
		panic(ExceptionNewf(NotImplementedError, "String decode for %q not implemented", encodingStr))
	}

	// We'd like to call PyObject_Bytes here, but we need to check for an
	// integer argument before deferring to PyBytes_FromObject, something
	// PyObject_Bytes doesn't do.
	var ok bool
	if I, ok := x.(I__bytes__); ok {
		New = I.M__bytes__()
	} else if New, ok = TypeCall0(x, "__bytes__"); ok {
	} else {
		goto no_bytes_method
	}
	if _, ok = New.(Bytes); !ok {
		panic(ExceptionNewf(TypeError, "__bytes__ returned non-bytes (type %s)", New.Type().Name))
	}
no_bytes_method:

	// Is it an integer?
	if _, ok := x.(Int); ok {
		size := IndexInt(x)
		if size < 0 {
			panic(ExceptionNewf(ValueError, "negative count"))
		}
		return make(Bytes, size)
	}

	// If it's not unicode, there can't be encoding or errors
	if encoding != nil || errors != nil {
		panic(ExceptionNewf(TypeError, "encoding or errors without a string argument"))
	}

	return BytesFromObject(x)
}

// Converts an object into bytes
func BytesFromObject(x Object) Bytes {
	// Look for special cases
	// FIXME implement converting from any object implementing the buffer API.
	switch z := x.(type) {
	case Bytes:
		// Immutable type so just return what was passed in
		return z
	case String:
		panic(ExceptionNewf(TypeError, "cannot convert unicode object to bytes"))
	}
	// Otherwise iterate through the whatever converting it into ints
	b := Bytes{}
	Iterate(x, func(item Object) {
		value := IndexInt(item)
		if value < 0 || value >= 256 {
			panic(ExceptionNewf(ValueError, "bytes must be in range(0, 256)"))
		}
		b = append(b, byte(value))
	})
	return b
}
