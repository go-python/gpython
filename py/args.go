// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Argument parsing for Go functions called by python
//
// These functions are useful when creating your own extensions
// functions and methods. Additional information and examples are
// available in Extending and Embedding the Python Interpreter.
//
// The first three of these functions described, PyArg_ParseTuple(),
// PyArg_ParseTupleAndKeywords(), and PyArg_Parse(), all use format
// strings which are used to tell the function about the expected
// arguments. The format strings use the same syntax for each of these
// functions.
//
// Parsing arguments
//
// A format string consists of zero or more “format units.” A format
// unit describes one Python object; it is usually a single character
// or a parenthesized sequence of format units. With a few exceptions,
// a format unit that is not a parenthesized sequence normally
// corresponds to a single address argument to these functions. In the
// following description, the quoted form is the format unit; the
// entry in (round) parentheses is the Python object type that matches
// the format unit; and the entry in [square] brackets is the type of
// the C variable(s) whose address should be passed.
//
// s (str) [const char *]
//
// Convert a Unicode object to a C pointer to a character string. A
// pointer to an existing string is stored in the character pointer
// variable whose address you pass. The C string is
// NUL-terminated. The Python string must not contain embedded NUL
// bytes; if it does, a TypeError exception is raised. Unicode objects
// are converted to C strings using 'utf-8' encoding. If this
// conversion fails, a UnicodeError is raised.
//
// Note This format does not accept bytes-like objects. If you want to
// accept filesystem paths and convert them to C character strings, it
// is preferable to use the O& format with PyUnicode_FSConverter() as
// converter.
//
// s* (str, bytes, bytearray or buffer compatible object) [Py_buffer]
//
// This format accepts Unicode objects as well as bytes-like
// objects. It fills a Py_buffer structure provided by the caller. In
// this case the resulting C string may contain embedded NUL
// bytes. Unicode objects are converted to C strings using 'utf-8'
// encoding.
//
// s# (str, bytes or read-only buffer compatible object) [const char *, int or Py_ssize_t]
//
// Like s*, except that it doesn’t accept mutable buffer-like objects
// such as bytearray. The result is stored into two C variables, the
// first one a pointer to a C string, the second one its length. The
// string may contain embedded null bytes. Unicode objects are
// converted to C strings using 'utf-8' encoding.
//
// z (str or None) [const char *]
//
// Like s, but the Python object may also be None, in which case the C
// pointer is set to NULL.
//
// z* (str, bytes, bytearray, buffer compatible object or None)
// [Py_buffer]
//
// Like s*, but the Python object may also be None, in which case the
// buf member of the Py_buffer structure is set to NULL.
//
// z# (str, bytes, read-only buffer compatible object or None) [const
// char *, int]
//
// Like s#, but the Python object may also be None, in which case the
// C pointer is set to NULL.
//
// y (bytes) [const char *]
//
// This format converts a bytes-like object to a C pointer to a
// character string; it does not accept Unicode objects. The bytes
// buffer must not contain embedded NUL bytes; if it does, a TypeError
// exception is raised.
//
// y* (bytes, bytearray or bytes-like object) [Py_buffer]
//
// This variant on s* doesn’t accept Unicode objects, only bytes-like
// objects. This is the recommended way to accept binary data.
//
// y# (bytes) [const char *, int]
//
// This variant on s# doesn’t accept Unicode objects, only bytes-like
// objects.
//
// S (bytes) [PyBytesObject *]
//
// Requires that the Python object is a bytes object, without
// attempting any conversion. Raises TypeError if the object is not a
// bytes object. The C variable may also be declared as PyObject*.
//
// Y (bytearray) [PyByteArrayObject *]
//
// Requires that the Python object is a bytearray object, without
// attempting any conversion. Raises TypeError if the object is not a
// bytearray object. The C variable may also be declared as PyObject*.
//
// u (str) [Py_UNICODE *]
//
// Convert a Python Unicode object to a C pointer to a NUL-terminated
// buffer of Unicode characters. You must pass the address of a
// Py_UNICODE pointer variable, which will be filled with the pointer
// to an existing Unicode buffer. Please note that the width of a
// Py_UNICODE character depends on compilation options (it is either
// 16 or 32 bits). The Python string must not contain embedded NUL
// characters; if it does, a TypeError exception is raised.
//
// Note Since u doesn’t give you back the length of the string, and it
// may contain embedded NUL characters, it is recommended to use u# or
// U instead.
//
// u# (str) [Py_UNICODE *, int]
//
// This variant on u stores into two C variables, the first one a
// pointer to a Unicode data buffer, the second one its length.
//
// Z (str or None) [Py_UNICODE *]
//
// Like u, but the Python object may also be None, in which case the
// Py_UNICODE pointer is set to NULL.
//
// Z# (str or None) [Py_UNICODE *, int]
//
// Like u#, but the Python object may also be None, in which case the
// Py_UNICODE pointer is set to NULL.
//
// U (str) [PyObject *]
//
// Requires that the Python object is a Unicode object, without
// attempting any conversion. Raises TypeError if the object is not a
// Unicode object. The C variable may also be declared as PyObject*.
//
// w* (bytearray or read-write byte-oriented buffer) [Py_buffer]
//
// This format accepts any object which implements the read-write
// buffer interface. It fills a Py_buffer structure provided by the
// caller. The buffer may contain embedded null bytes. The caller have
// to call PyBuffer_Release() when it is done with the buffer.
//
// es (str) [const char *encoding, char **buffer]
//
// This variant on s is used for encoding Unicode into a character
// buffer. It only works for encoded data without embedded NUL bytes.
//
// This format requires two arguments. The first is only used as
// input, and must be a const char* which points to the name of an
// encoding as a NUL-terminated string, or NULL, in which case 'utf-8'
// encoding is used. An exception is raised if the named encoding is
// not known to Python. The second argument must be a char**; the
// value of the pointer it references will be set to a buffer with the
// contents of the argument text. The text will be encoded in the
// encoding specified by the first argument.
//
// PyArg_ParseTuple() will allocate a buffer of the needed size, copy
// the encoded data into this buffer and adjust *buffer to reference
// the newly allocated storage. The caller is responsible for calling
// PyMem_Free() to free the allocated buffer after use.
//
// et (str, bytes or bytearray) [const char *encoding, char **buffer]
//
// Same as es except that byte string objects are passed through
// without recoding them. Instead, the implementation assumes that the
// byte string object uses the encoding passed in as parameter.
//
// es# (str) [const char *encoding, char **buffer, int *buffer_length]
//
// This variant on s# is used for encoding Unicode into a character
// buffer. Unlike the es format, this variant allows input data which
// contains NUL characters.
//
// It requires three arguments. The first is only used as input, and
// must be a const char* which points to the name of an encoding as a
// NUL-terminated string, or NULL, in which case 'utf-8' encoding is
// used. An exception is raised if the named encoding is not known to
// Python. The second argument must be a char**; the value of the
// pointer it references will be set to a buffer with the contents of
// the argument text. The text will be encoded in the encoding
// specified by the first argument. The third argument must be a
// pointer to an integer; the referenced integer will be set to the
// number of bytes in the output buffer.
//
// There are two modes of operation:
//
// If *buffer points a NULL pointer, the function will allocate a
// buffer of the needed size, copy the encoded data into this buffer
// and set *buffer to reference the newly allocated storage. The
// caller is responsible for calling PyMem_Free() to free the
// allocated buffer after usage.
//
// If *buffer points to a non-NULL pointer (an already allocated
// buffer), PyArg_ParseTuple() will use this location as the buffer
// and interpret the initial value of *buffer_length as the buffer
// size. It will then copy the encoded data into the buffer and
// NUL-terminate it. If the buffer is not large enough, a ValueError
// will be set.
//
// In both cases, *buffer_length is set to the length of the encoded
// data without the trailing NUL byte.
//
// et# (str, bytes or bytearray) [const char *encoding, char **buffer,
// int *buffer_length]
//
// Same as es# except that byte string objects are passed through
// without recoding them. Instead, the implementation assumes that the
// byte string object uses the encoding passed in as parameter.
//
// Numbers
//
// b (int) [unsigned char]
//
// Convert a nonnegative Python integer to an unsigned tiny int,
// stored in a C unsigned char.
//
// B (int) [unsigned char]
//
// Convert a Python integer to a tiny int without overflow checking,
// stored in a C unsigned char.  h (int) [short int]
//
// Convert a Python integer to a C short int.
//
// H (int) [unsigned short int]
//
// Convert a Python integer to a C unsigned short int, without
// overflow checking.
//
// i (int) [int]
//
// Convert a Python integer to a plain C int.
//
// I (int) [unsigned int]
//
// Convert a Python integer to a C unsigned int, without overflow
// checking.
//
// l (int) [long int]
//
// Convert a Python integer to a C long int.
//
// k (int) [unsigned long]
//
// Convert a Python integer to a C unsigned long without overflow
// checking.
//
// L (int) [PY_LONG_LONG]
//
// Convert a Python integer to a C long long. This format is only
// available on platforms that support long long (or _int64 on
// Windows).
//
// K (int) [unsigned PY_LONG_LONG]
//
// Convert a Python integer to a C unsigned long long without overflow
// checking. This format is only available on platforms that support
// unsigned long long (or unsigned _int64 on Windows).
//
// n (int) [Py_ssize_t]
//
// Convert a Python integer to a C Py_ssize_t.
//
// c (bytes or bytearray of length 1) [char]
//
// Convert a Python byte, represented as a bytes or bytearray object
// of length 1, to a C char.
//
// Changed in version 3.3: Allow bytearray objects.
//
// C (str of length 1) [int]
//
// Convert a Python character, represented as a str object of length 1, to a C int.
//
// f (float) [float]
//
// Convert a Python floating point number to a C float.
//
// d (float) [double]
//
// Convert a Python floating point number to a C double.
//
// D (complex) [Py_complex]
//
// Convert a Python complex number to a C Py_complex structure.
//
// Other objects
//
// O (object) [PyObject *]
//
// Store a Python object (without any conversion) in a C object
// pointer. The C program thus receives the actual object that was
// passed. The object’s reference count is not increased. The pointer
// stored is not NULL.
//
// O! (object) [typeobject, PyObject *]
//
// Store a Python object in a C object pointer. This is similar to O,
// but takes two C arguments: the first is the address of a Python
// type object, the second is the address of the C variable (of type
// PyObject*) into which the object pointer is stored. If the Python
// object does not have the required type, TypeError is raised.
//
// O& (object) [converter, anything]
//
// Convert a Python object to a C variable through a converter
// function. This takes two arguments: the first is a function, the
// second is the address of a C variable (of arbitrary type),
// converted to void *. The converter function in turn is called as
// follows:
//
// status = converter(object, address);
//
// where object is the Python object to be converted and address is
// the void* argument that was passed to the PyArg_Parse*()
// function. The returned status should be 1 for a successful
// conversion and 0 if the conversion has failed. When the conversion
// fails, the converter function should raise an exception and leave
// the content of address unmodified.
//
// If the converter returns Py_CLEANUP_SUPPORTED, it may get called a
// second time if the argument parsing eventually fails, giving the
// converter a chance to release any memory that it had already
// allocated. In this second call, the object parameter will be NULL;
// address will have the same value as in the original call.
//
// Changed in version 3.1: Py_CLEANUP_SUPPORTED was added.
//
// p (bool) [int]
//
// Tests the value passed in for truth (a boolean predicate) and
// converts the result to its equivalent C true/false integer
// value. Sets the int to 1 if the expression was true and 0 if it was
// false. This accepts any valid Python value. See Truth Value Testing
// for more information about how Python tests values for truth.
//
// New in version 3.3.
//
// (items) (tuple) [matching-items]
//
// The object must be a Python sequence whose length is the number of
// format units in items. The C arguments must correspond to the
// individual format units in items. Format units for sequences may be
// nested.
//
// It is possible to pass “long” integers (integers whose value
// exceeds the platform’s LONG_MAX) however no proper range checking
// is done — the most significant bits are silently truncated when the
// receiving field is too small to receive the value (actually, the
// semantics are inherited from downcasts in C — your mileage may
// vary).
//
// A few other characters have a meaning in a format string. These may
// not occur inside nested parentheses. They are:
//
// |
//
// Indicates that the remaining arguments in the Python argument list
// are optional. The C variables corresponding to optional arguments
// should be initialized to their default value — when an optional
// argument is not specified, PyArg_ParseTuple() does not touch the
// contents of the corresponding C variable(s).
//
// $
//
// PyArg_ParseTupleAndKeywords() only: Indicates that the remaining
// arguments in the Python argument list are keyword-only. Currently,
// all keyword-only arguments must also be optional arguments, so |
// must always be specified before $ in the format string.
//
// New in version 3.3.
//
// :
//
// The list of format units ends here; the string after the colon is
// used as the function name in error messages (the “associated value”
// of the exception that PyArg_ParseTuple() raises).
//
// ;
//
// The list of format units ends here; the string after the semicolon
// is used as the error message instead of the default error
// message. : and ; mutually exclude each other.
//
// Note that any Python object references which are provided to the
// caller are borrowed references; do not decrement their reference
// count!
//
// Additional arguments passed to these functions must be addresses of
// variables whose type is determined by the format string; these are
// used to store values from the input tuple. There are a few cases,
// as described in the list of format units above, where these
// parameters are used as input values; they should match what is
// specified for the corresponding format unit in that case.
//
// For the conversion to succeed, the arg object must match the format
// and the format must be exhausted. On success, the PyArg_Parse*()
// functions return true, otherwise they return false and raise an
// appropriate exception. When the PyArg_Parse*() functions fail due
// to conversion failure in one of the format units, the variables at
// the addresses corresponding to that and the following format units
// are left untouched.

package py

// FIXME this would be a lot more useful if we could supply the
// address of a String rather than an Object - would then need
// introspection to set it properly

// ParseTupleAndKeywords
func ParseTupleAndKeywords(args Tuple, kwargs StringDict, format string, kwlist []string, results ...*Object) error {
	if kwlist != nil && len(results) != len(kwlist) {
		return ExceptionNewf(TypeError, "Internal error: supply the same number of results and kwlist")
	}
	min, max, name, ops := parseFormat(format)
	keywordOnly := false
	err := checkNumberOfArgs(name, len(args)+len(kwargs), len(results), min, max)
	if err != nil {
		return err
	}

	if len(ops) > 0 && ops[0] == "$" {
		keywordOnly = true
		ops = ops[1:]
	}
	// Check all the kwargs are in kwlist
	// O(N^2) Slow but kwlist is usually short
	for kwargName := range kwargs {
		for _, kw := range kwlist {
			if kw == kwargName {
				goto found
			}
		}
		return ExceptionNewf(TypeError, "%s() got an unexpected keyword argument '%s'", name, kwargName)
	found:
	}

	// Create args tuple with all the arguments we have in
	args = args.Copy()
	for i, kw := range kwlist {
		if value, ok := kwargs[kw]; ok {
			if len(args) > i {
				return ExceptionNewf(TypeError, "%s() got multiple values for argument '%s'", name, kw)
			}
			args = append(args, value)
		} else if keywordOnly {
			args = append(args, nil)
		}
	}
	for i, arg := range args {
		op := ops[i]
		result := results[i]
		switch op {
		case "O":
			*result = arg
		case "Z", "z":
			if _, ok := arg.(NoneType); ok {
				*result = arg
				break
			}
			fallthrough
		case "U", "s":
			if _, ok := arg.(String); !ok {
				return ExceptionNewf(TypeError, "%s() argument %d must be str, not %s", name, i+1, arg.Type().Name)
			}
			*result = arg
		case "i":
			if _, ok := arg.(Int); !ok {
				return ExceptionNewf(TypeError, "%s() argument %d must be int, not %s", name, i+1, arg.Type().Name)
			}
			*result = arg
		case "p":
			if _, ok := arg.(Bool); !ok {
				return ExceptionNewf(TypeError, "%s() argument %d must be bool, not %s", name, i+1, arg.Type().Name)
			}
			*result = arg
		case "d":
			switch x := arg.(type) {
			case Int:
				*result = Float(x)
			case Float:
				*result = x
			default:
				return ExceptionNewf(TypeError, "%s() argument %d must be float, not %s", name, i+1, arg.Type().Name)
			}

		default:
			return ExceptionNewf(TypeError, "Unknown/Unimplemented format character %q in ParseTupleAndKeywords called from %s", op, name)
		}
	}
	return nil
}

// Parse tuple only
func ParseTuple(args Tuple, format string, results ...*Object) error {
	return ParseTupleAndKeywords(args, nil, format, nil, results...)
}

// Parse the format
func parseFormat(format string) (min, max int, name string, ops []string) {
	name = "function"
	min = -1
	for format != "" {
		op := string(format[0])
		format = format[1:]
		if len(format) > 1 && (format[1] == '*' || format[1] == '#') {
			op += string(format[0])
			format = format[1:]
		}
		switch op {
		case ":", ";":
			name = format
			format = ""
		case "|":
			min = len(ops)
		default:
			ops = append(ops, op)
		}
	}
	max = len(ops)
	if min < 0 {
		min = max
	}
	return
}

// Checks the number of args passed in
func checkNumberOfArgs(name string, nargs, nresults, min, max int) error {
	if min == max {
		if nargs != max {
			return ExceptionNewf(TypeError, "%s() takes exactly %d arguments (%d given)", name, max, nargs)
		}
	} else {
		if nargs > max {
			return ExceptionNewf(TypeError, "%s() takes at most %d arguments (%d given)", name, max, nargs)
		}
		if nargs < min {
			return ExceptionNewf(TypeError, "%s() takes at least %d arguments (%d given)", name, min, nargs)
		}
	}

	if nargs > nresults {
		return ExceptionNewf(TypeError, "Internal error: not enough arguments supplied to Unpack*/Parse*")
	}
	return nil
}

// Unpack the args tuple into the results
//
// Up to the caller to set default values
func UnpackTuple(args Tuple, kwargs StringDict, name string, min int, max int, results ...*Object) error {
	if len(kwargs) != 0 {
		return ExceptionNewf(TypeError, "%s() does not take keyword arguments", name)
	}

	// Check number of arguments
	err := checkNumberOfArgs(name, len(args), len(results), min, max)
	if err != nil {
		return err
	}

	// Copy the results in
	for i := range args {
		*results[i] = args[i]
	}
	return nil
}
