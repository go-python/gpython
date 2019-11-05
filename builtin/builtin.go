// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Built-in functions
package builtin

import (
	"fmt"
	"math/big"
	"unicode/utf8"

	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/py"
)

const builtin_doc = `Built-in functions, exceptions, and other objects.

Noteworthy: None is the 'nil' object; Ellipsis represents '...' in slices.`

// Initialise the module
func init() {
	methods := []*py.Method{
		py.MustNewMethod("__build_class__", builtin___build_class__, 0, build_class_doc),
		py.MustNewMethod("__import__", py.InternalMethodImport, 0, import_doc),
		py.MustNewMethod("abs", builtin_abs, 0, abs_doc),
		py.MustNewMethod("all", builtin_all, 0, all_doc),
		py.MustNewMethod("any", builtin_any, 0, any_doc),
		py.MustNewMethod("ascii", builtin_ascii, 0, ascii_doc),
		py.MustNewMethod("bin", builtin_bin, 0, bin_doc),
		// py.MustNewMethod("callable", builtin_callable, 0, callable_doc),
		py.MustNewMethod("chr", builtin_chr, 0, chr_doc),
		py.MustNewMethod("compile", builtin_compile, 0, compile_doc),
		py.MustNewMethod("delattr", builtin_delattr, 0, delattr_doc),
		// py.MustNewMethod("dir", builtin_dir, 0, dir_doc),
		py.MustNewMethod("divmod", builtin_divmod, 0, divmod_doc),
		py.MustNewMethod("eval", py.InternalMethodEval, 0, eval_doc),
		py.MustNewMethod("exec", py.InternalMethodExec, 0, exec_doc),
		// py.MustNewMethod("format", builtin_format, 0, format_doc),
		py.MustNewMethod("getattr", builtin_getattr, 0, getattr_doc),
		py.MustNewMethod("globals", py.InternalMethodGlobals, 0, globals_doc),
		py.MustNewMethod("hasattr", builtin_hasattr, 0, hasattr_doc),
		// py.MustNewMethod("hash", builtin_hash, 0, hash_doc),
		py.MustNewMethod("hex", builtin_hex, 0, hex_doc),
		// py.MustNewMethod("id", builtin_id, 0, id_doc),
		// py.MustNewMethod("input", builtin_input, 0, input_doc),
		py.MustNewMethod("isinstance", builtin_isinstance, 0, isinstance_doc),
		// py.MustNewMethod("issubclass", builtin_issubclass, 0, issubclass_doc),
		py.MustNewMethod("iter", builtin_iter, 0, iter_doc),
		py.MustNewMethod("len", builtin_len, 0, len_doc),
		py.MustNewMethod("locals", py.InternalMethodLocals, 0, locals_doc),
		py.MustNewMethod("max", builtin_max, 0, max_doc),
		py.MustNewMethod("min", builtin_min, 0, min_doc),
		py.MustNewMethod("next", builtin_next, 0, next_doc),
		py.MustNewMethod("open", builtin_open, 0, open_doc),
		// py.MustNewMethod("oct", builtin_oct, 0, oct_doc),
		py.MustNewMethod("ord", builtin_ord, 0, ord_doc),
		py.MustNewMethod("pow", builtin_pow, 0, pow_doc),
		py.MustNewMethod("print", builtin_print, 0, print_doc),
		py.MustNewMethod("repr", builtin_repr, 0, repr_doc),
		py.MustNewMethod("round", builtin_round, 0, round_doc),
		py.MustNewMethod("setattr", builtin_setattr, 0, setattr_doc),
		py.MustNewMethod("sorted", builtin_sorted, 0, sorted_doc),
		py.MustNewMethod("sum", builtin_sum, 0, sum_doc),
		// py.MustNewMethod("vars", builtin_vars, 0, vars_doc),
	}
	globals := py.StringDict{
		"None":     py.None,
		"Ellipsis": py.Ellipsis,
		"False":    py.False,
		"True":     py.True,
		"bool":     py.BoolType,
		// "memoryview":     py.MemoryViewType,
		// "bytearray":      py.ByteArrayType,
		"bytes":       py.BytesType,
		"classmethod": py.ClassMethodType,
		"complex":     py.ComplexType,
		"dict":        py.StringDictType, // FIXME
		"enumerate":   py.EnumerateType,
		// "filter":         py.FilterType,
		"float":     py.FloatType,
		"frozenset": py.FrozenSetType,
		// "property":       py.PropertyType,
		"int":  py.IntType, // FIXME LongType?
		"list": py.ListType,
		// "map":            py.MapType,
		"object": py.ObjectType,
		"range":  py.RangeType,
		// "reversed":       py.ReversedType,
		"set":          py.SetType,
		"slice":        py.SliceType,
		"staticmethod": py.StaticMethodType,
		"str":          py.StringType,
		// "super":          py.SuperType,
		"tuple": py.TupleType,
		"type":  py.TypeType,
		"zip":   py.ZipType,

		// Exceptions
		"ArithmeticError":           py.ArithmeticError,
		"AssertionError":            py.AssertionError,
		"AttributeError":            py.AttributeError,
		"BaseException":             py.BaseException,
		"BlockingIOError":           py.BlockingIOError,
		"BrokenPipeError":           py.BrokenPipeError,
		"BufferError":               py.BufferError,
		"BytesWarning":              py.BytesWarning,
		"ChildProcessError":         py.ChildProcessError,
		"ConnectionAbortedError":    py.ConnectionAbortedError,
		"ConnectionError":           py.ConnectionError,
		"ConnectionRefusedError":    py.ConnectionRefusedError,
		"ConnectionResetError":      py.ConnectionResetError,
		"DeprecationWarning":        py.DeprecationWarning,
		"EOFError":                  py.EOFError,
		"EnvironmentError":          py.OSError,
		"Exception":                 py.ExceptionType,
		"FileExistsError":           py.FileExistsError,
		"FileNotFoundError":         py.FileNotFoundError,
		"FloatingPointError":        py.FloatingPointError,
		"FutureWarning":             py.FutureWarning,
		"GeneratorExit":             py.GeneratorExit,
		"IOError":                   py.OSError,
		"ImportError":               py.ImportError,
		"ImportWarning":             py.ImportWarning,
		"IndentationError":          py.IndentationError,
		"IndexError":                py.IndexError,
		"InterruptedError":          py.InterruptedError,
		"IsADirectoryError":         py.IsADirectoryError,
		"KeyError":                  py.KeyError,
		"KeyboardInterrupt":         py.KeyboardInterrupt,
		"LookupError":               py.LookupError,
		"MemoryError":               py.MemoryError,
		"NameError":                 py.NameError,
		"NotADirectoryError":        py.NotADirectoryError,
		"NotImplemented":            py.NotImplemented,
		"NotImplementedError":       py.NotImplementedError,
		"OSError":                   py.OSError,
		"OverflowError":             py.OverflowError,
		"PendingDeprecationWarning": py.PendingDeprecationWarning,
		"PermissionError":           py.PermissionError,
		"ProcessLookupError":        py.ProcessLookupError,
		"ReferenceError":            py.ReferenceError,
		"ResourceWarning":           py.ResourceWarning,
		"RuntimeError":              py.RuntimeError,
		"RuntimeWarning":            py.RuntimeWarning,
		"StopIteration":             py.StopIteration,
		"SyntaxError":               py.SyntaxError,
		"SyntaxWarning":             py.SyntaxWarning,
		"SystemError":               py.SystemError,
		"SystemExit":                py.SystemExit,
		"TabError":                  py.TabError,
		"TimeoutError":              py.TimeoutError,
		"TypeError":                 py.TypeError,
		"UnboundLocalError":         py.UnboundLocalError,
		"UnicodeDecodeError":        py.UnicodeDecodeError,
		"UnicodeEncodeError":        py.UnicodeEncodeError,
		"UnicodeError":              py.UnicodeError,
		"UnicodeTranslateError":     py.UnicodeTranslateError,
		"UnicodeWarning":            py.UnicodeWarning,
		"UserWarning":               py.UserWarning,
		"ValueError":                py.ValueError,
		"Warning":                   py.Warning,
		"ZeroDivisionError":         py.ZeroDivisionError,
	}
	py.NewModule("builtins", builtin_doc, methods, globals)
}

const print_doc = `print(value, ..., sep=' ', end='\\n', file=sys.stdout, flush=False)

Prints the values to a stream, or to sys.stdout by default.
Optional keyword arguments:
file:  a file-like object (stream); defaults to the current sys.stdout.
sep:   string inserted between values, default a space.
end:   string appended after the last value, default a newline.
flush: whether to forcibly flush the stream.`

func builtin_print(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		sepObj py.Object = py.String(" ")
		endObj py.Object = py.String("\n")
		file   py.Object = py.MustGetModule("sys").Globals["stdout"]
		flush  py.Object
	)
	kwlist := []string{"sep", "end", "file", "flush"}
	err := py.ParseTupleAndKeywords(nil, kwargs, "|ssOO:print", kwlist, &sepObj, &endObj, &file, &flush)
	if err != nil {
		return nil, err
	}
	sep := sepObj.(py.String)
	end := endObj.(py.String)

	write, err := py.GetAttrString(file, "write")
	if err != nil {
		return nil, err
	}

	for i, v := range args {
		v, err := py.Str(v)
		if err != nil {
			return nil, err
		}

		_, err = py.Call(write, py.Tuple{v}, nil)
		if err != nil {
			return nil, err
		}

		if i != len(args)-1 {
			_, err = py.Call(write, py.Tuple{sep}, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err = py.Call(write, py.Tuple{end}, nil)
	if err != nil {
		return nil, err
	}

	if shouldFlush, _ := py.MakeBool(flush); shouldFlush == py.True {
		fflush, err := py.GetAttrString(file, "flush")
		if err == nil {
			return py.Call(fflush, nil, nil)
		}
	}

	return py.None, nil
}

const repr_doc = `repr(object) -> string

Return the canonical string representation of the object.
For most object types, eval(repr(object)) == object.`

func builtin_repr(self py.Object, obj py.Object) (py.Object, error) {
	return py.Repr(obj)
}

const pow_doc = `pow(x, y[, z]) -> number

With two arguments, equivalent to x**y.  With three arguments,
equivalent to (x**y) % z, but may be more efficient (e.g. for ints).`

func builtin_pow(self py.Object, args py.Tuple) (py.Object, error) {
	var v, w, z py.Object
	z = py.None
	err := py.UnpackTuple(args, nil, "pow", 2, 3, &v, &w, &z)
	if err != nil {
		return nil, err
	}
	return py.Pow(v, w, z)
}

const abs_doc = `"abs(number) -> number

Return the absolute value of the argument.`

func builtin_abs(self, v py.Object) (py.Object, error) {
	return py.Abs(v)
}

const all_doc = `all(iterable) -> bool

Return True if bool(x) is True for all values x in the iterable.
If the iterable is empty, return True.
`

func builtin_all(self, seq py.Object) (py.Object, error) {
	iter, err := py.Iter(seq)
	if err != nil {
		return nil, err
	}
	for {
		item, err := py.Next(iter)
		if err != nil {
			if py.IsException(py.StopIteration, err) {
				break
			}
			return nil, err
		}
		if !py.ObjectIsTrue(item) {
			return py.False, nil
		}
	}
	return py.True, nil
}

const any_doc = `any(iterable) -> bool

Return True if bool(x) is True for any x in the iterable.
If the iterable is empty, Py_RETURN_FALSE."`

func builtin_any(self, seq py.Object) (py.Object, error) {
	iter, err := py.Iter(seq)
	if err != nil {
		return nil, err
	}
	for {
		item, err := py.Next(iter)
		if err != nil {
			if py.IsException(py.StopIteration, err) {
				break
			}
			return nil, err
		}
		if py.ObjectIsTrue(item) {
			return py.True, nil
		}
	}
	return py.False, nil
}

const ascii_doc = `Return an ASCII-only representation of an object.

As repr(), return a string containing a printable representation of an
object, but escape the non-ASCII characters in the string returned by
repr() using \\x, \\u or \\U escapes. This generates a string similar
to that returned by repr() in Python 2.
`

func builtin_ascii(self, o py.Object) (py.Object, error) {
	reprObj, err := py.Repr(o)
	if err != nil {
		return nil, err
	}
	repr := reprObj.(py.String)
	out := py.StringEscape(repr, true)
	return py.String(out), err
}

const bin_doc = `Return the binary representation of an integer.

>>> bin(2796202)
'0b1010101010101010101010'
`

func builtin_bin(self, o py.Object) (py.Object, error) {
	bigint, ok := py.ConvertToBigInt(o)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "'%s' object cannot be interpreted as an integer", o.Type().Name)
	}

	value := (*big.Int)(bigint)
	var out string
	if value.Sign() < 0 {
		value = new(big.Int).Abs(value)
		out = fmt.Sprintf("-0b%b", value)
	} else {
		out = fmt.Sprintf("0b%b", value)
	}
	return py.String(out), nil
}

const round_doc = `round(number[, ndigits]) -> number

Round a number to a given precision in decimal digits (default 0 digits).
This returns an int when called with one argument, otherwise the
same type as the number. ndigits may be negative.`

func builtin_round(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var number, ndigits py.Object
	ndigits = py.Int(0)
	// var kwlist = []string{"number", "ndigits"}
	// FIXME py.ParseTupleAndKeywords(args, kwargs, "O|O:round", kwlist, &number, &ndigits)
	err := py.UnpackTuple(args, nil, "round", 1, 2, &number, &ndigits)
	if err != nil {
		return nil, err
	}

	numberRounder, ok := number.(py.I__round__)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "type %s doesn't define __round__ method", number.Type().Name)
	}

	return numberRounder.M__round__(ndigits)
}

const build_class_doc = `__build_class__(func, name, *bases, metaclass=None, **kwds) -> class

Internal helper function used by the class statement.`

func builtin___build_class__(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	// fmt.Printf("__build_class__(self=%#v, args=%#v, kwargs=%#v\n", self, args, kwargs)
	var prep, cell, cls py.Object
	var mkw, ns py.StringDict
	var meta, winner *py.Type
	var isclass bool
	var err error

	if len(args) < 2 {
		return nil, py.ExceptionNewf(py.TypeError, "__build_class__: not enough arguments")
	}

	// Better be callable
	fn, ok := args[0].(*py.Function)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "__build__class__: func must be a function")
	}

	name := args[1].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "__build_class__: name is not a string")
	}
	bases := args[2:]

	if kwargs != nil {
		mkw = kwargs.Copy()      // Don't modify kwds passed in!
		meta := mkw["metaclass"] // _PyDict_GetItemId(mkw, &PyId_metaclass)
		if meta != nil {
			delete(mkw, "metaclass")
			// metaclass is explicitly given, check if it's indeed a class
			_, isclass = meta.(*py.Type)
		}
	}
	if meta == nil {
		// if there are no bases, use type:
		if len(bases) == 0 {
			meta = py.TypeType
		} else {
			// else get the type of the first base
			meta = bases[0].Type()
		}
		isclass = true // meta is really a class
	}

	if isclass {
		// meta is really a class, so check for a more derived
		// metaclass, or possible metaclass conflicts:
		winner, err = meta.CalculateMetaclass(bases)
		if err != nil {
			return nil, err
		}
		if winner != meta {
			meta = winner
		}
	}
	// else: meta is not a class, so we cannot do the metaclass
	// calculation, so we will use the explicitly given object as it is
	prep = meta.Type().Dict["___prepare__"] // FIXME should be using _PyObject_GetAttr
	if prep == nil {
		ns = py.NewStringDict()
	} else {
		nsObj, err := py.Call(prep, py.Tuple{name, bases}, mkw)
		if err != nil {
			return nil, err
		}
		ns = nsObj.(py.StringDict)
	}
	// fmt.Printf("Calling %v with %v and %v\n", fn.Name, fn.Globals, ns)
	// fmt.Printf("Code = %#v\n", fn.Code)
	cell, err = py.VmRun(fn.Globals, ns, fn.Code, fn.Closure)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("result = %#v err = %s\n", cell, err)
	// fmt.Printf("locals = %#v\n", locals)
	// fmt.Printf("ns = %#v\n", ns)
	if cell != nil {
		// fmt.Printf("Calling %v\n", meta)
		cls, err = py.Call(meta, py.Tuple{name, bases, ns}, mkw)
		if err != nil {
			return nil, err
		}
		if c, ok := cell.(*py.Cell); ok {
			c.Set(cls)
		}
	}
	// fmt.Printf("Globals = %v, Locals = %v\n", fn.Globals, ns)
	return cls, nil
}

const next_doc = `next(iterator[, default])

Return the next item from the iterator. If default is given and the iterator
is exhausted, it is returned instead of raising StopIteration.`

func builtin_next(self py.Object, args py.Tuple) (res py.Object, err error) {
	var it, def py.Object

	err = py.UnpackTuple(args, nil, "next", 1, 2, &it, &def)
	if err != nil {
		return nil, err
	}

	res, err = py.Next(it)
	if err != nil && def != nil && py.IsException(py.StopIteration, err) {
		// Return defult on StopIteration
		res = def
		err = nil
	}
	return res, err
}

const import_doc = `__import__(name, globals=None, locals=None, fromlist=(), level=0) -> module

Import a module. Because this function is meant for use by the Python
interpreter and not for general use it is better to use
importlib.import_module() to programmatically import a module.

The globals argument is only used to determine the context;
they are not modified.  The locals argument is unused.  The fromlist
should be a list of names to emulate ''from name import ...'', or an
empty list to emulate ''import name''.
When importing a module from a package, note that __import__('A.B', ...)
returns package A when fromlist is empty, but its submodule B when
fromlist is not empty.  Level is used to determine whether to perform 
absolute or relative imports. 0 is absolute while a positive number
is the number of parent directories to search relative to the current module.`

const open_doc = `open(name[, mode[, buffering]]) -> file object

Open a file using the file() type, returns a file object.  This is the
preferred way to open a file.  See file.__doc__ for further information.`

func builtin_open(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	kwlist := []string{
		"file",
		"mode",
		"buffering",
		"encoding",
		"errors",
		"newline",
		"closefd",
		"opener",
	}

	var (
		filename  py.Object
		mode      py.Object = py.String("r")
		buffering py.Object = py.Int(-1)
		encoding  py.Object = py.None
		errors    py.Object = py.None
		newline   py.Object = py.None
		closefd   py.Object = py.Bool(true)
		opener    py.Object = py.None
	)

	err := py.ParseTupleAndKeywords(args, kwargs, "s|sizzzpO:open", kwlist,
		&filename,
		&mode,
		&buffering,
		&encoding,
		&errors,
		&newline,
		&closefd,
		&opener)
	if err != nil {
		return nil, err
	}

	if encoding != py.None && encoding.(py.String) != py.String("utf-8") {
		return nil, py.ExceptionNewf(py.NotImplementedError, "encoding not implemented yet")
	}

	if errors != py.None {
		return nil, py.ExceptionNewf(py.NotImplementedError, "errors not implemented yet")
	}

	if newline != py.None {
		return nil, py.ExceptionNewf(py.NotImplementedError, "newline not implemented yet")
	}

	if opener != py.None {
		return nil, py.ExceptionNewf(py.NotImplementedError, "opener not implemented yet")
	}

	return py.OpenFile(string(filename.(py.String)),
		string(mode.(py.String)),
		int(buffering.(py.Int)))
}

const ord_doc = `ord(c) -> integer

Return the integer ordinal of a one-character string.`

func builtin_ord(self, obj py.Object) (py.Object, error) {
	var size int
	switch x := obj.(type) {
	case py.Bytes:
		size = len(x)
		if size == 1 {
			return py.Int(x[0]), nil
		}
	case py.String:
		size = len(x)
		rune, runeSize := utf8.DecodeRuneInString(string(x))
		if size == runeSize && rune != utf8.RuneError {
			return py.Int(rune), nil
		}
	//case py.ByteArray:
	// XXX Hopefully this is temporary
	// FIXME implement
	// size = PyByteArray_GET_SIZE(obj)
	// if size == 1 {
	// 	ord = (long)((char) * PyByteArray_AS_STRING(obj))
	// 	return PyLong_FromLong(ord)
	// }
	default:
		return nil, py.ExceptionNewf(py.TypeError, "ord() expected string of length 1, but %s found", obj.Type().Name)
	}

	return nil, py.ExceptionNewf(py.TypeError, "ord() expected a character, but string of length %d found", size)
}

const getattr_doc = `getattr(object, name[, default]) -> value

Get a named attribute from an object; getattr(x, 'y') is equivalent to x.y.
When a default argument is given, it is returned when the attribute doesn't
exist; without it, an exception is raised in that case.`

func builtin_getattr(self py.Object, args py.Tuple) (py.Object, error) {
	var v, result, dflt py.Object
	var name py.Object

	err := py.UnpackTuple(args, nil, "getattr", 2, 3, &v, &name, &dflt)
	if err != nil {
		return nil, err
	}

	result, err = py.GetAttr(v, name)
	if err != nil {
		if dflt == nil {
			return nil, err
		}
		result = dflt
	}
	return result, nil
}

const hasattr_doc = `hasattr(object, name) -> bool

Return whether the object has an attribute with the given name.
(This is done by calling getattr(object, name) and catching AttributeError.)`

func builtin_hasattr(self py.Object, args py.Tuple) (py.Object, error) {
	var v py.Object
	var name py.Object
	err := py.UnpackTuple(args, nil, "hasattr", 2, 2, &v, &name)
	if err != nil {
		return nil, err
	}
	_, err = py.GetAttr(v, name)
	return py.NewBool(err == nil), nil
}

const setattr_doc = `setattr(object, name, value)

Set a named attribute on an object; setattr(x, 'y', v) is equivalent to
"x.y = v".`

func builtin_setattr(self py.Object, args py.Tuple) (py.Object, error) {
	var v py.Object
	var name py.Object
	var value py.Object

	err := py.UnpackTuple(args, nil, "setattr", 3, 3, &v, &name, &value)
	if err != nil {
		return nil, err
	}

	return py.SetAttr(v, name, value)
}

// Reads the source as a string
func source_as_string(cmd py.Object, funcname, what string /*, PyCompilerFlags *cf */) (string, error) {
	// FIXME only understands strings, not bytes etc at the moment
	if str, ok := cmd.(py.String); ok {
		// FIXME cf->cf_flags |= PyCF_IGNORE_COOKIE;
		return string(str), nil
	}
	// } else if (!PyObject_CheckReadBuffer(cmd)) {
	return "", py.ExceptionNewf(py.TypeError, "%s() arg 1 must be a %s object", funcname, what)
	// } else if (PyObject_AsReadBuffer(cmd, (const void **)&str, &size) < 0) {
	// 	return nil;
}

const delattr_doc = `Deletes the named attribute from the given object.

delattr(x, 'y') is equivalent to  "del x.y"
`

func builtin_delattr(self py.Object, args py.Tuple) (py.Object, error) {
	var v py.Object
	var name py.Object

	err := py.UnpackTuple(args, nil, "delattr", 2, 2, &v, &name)
	if err != nil {
		return nil, err
	}

	err = py.DeleteAttr(v, name)
	if err != nil {
		return nil, err
	}
	return py.None, nil
}

const compile_doc = `compile(source, filename, mode[, flags[, dont_inherit]]) -> code object

Compile the source string (a Python module, statement or expression)
into a code object that can be executed by exec() or eval().
The filename will be used for run-time error messages.
The mode must be 'exec' to compile a module, 'single' to compile a
single (interactive) statement, or 'eval' to compile an expression.
The flags argument, if present, controls which future statements influence
the compilation of the code.
The dont_inherit argument, if non-zero, stops the compilation inheriting
the effects of any future statements in effect in the code calling
compile; if absent or zero these statements do influence the compilation,
in addition to any features explicitly specified.`

func builtin_compile(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	// FIXME lots of unsupported stuff here!
	var filename py.Object
	var startstr py.Object
	// var mode = -1
	var dont_inherit py.Object = py.Int(0)
	var supplied_flags py.Object = py.Int(0)
	var optimizeInt py.Object = py.Int(-1)
	//is_ast := false
	// var cf PyCompilerFlags
	var cmd py.Object
	kwlist := []string{"source", "filename", "mode", "flags", "dont_inherit", "optimize"}
	// start := []int{Py_file_input, Py_eval_input, Py_single_input}
	var result py.Object

	err := py.ParseTupleAndKeywords(args, kwargs, "Oss|iii:compile", kwlist,
		&cmd,
		&filename,
		&startstr,
		&supplied_flags,
		&dont_inherit,
		&optimizeInt)
	if err != nil {
		return nil, err
	}

	// cf.cf_flags = supplied_flags | PyCF_SOURCE_IS_UTF8

	// if supplied_flags&^(PyCF_MASK|PyCF_MASK_OBSOLETE|PyCF_DONT_IMPLY_DEDENT|PyCF_ONLY_AST) != 0 {
	// 	return nil, py.ExceptionNewf(py.ValueError, "compile(): unrecognised flags")
	// }
	// XXX Warn if (supplied_flags & PyCF_MASK_OBSOLETE) != 0?

	optimize := int(optimizeInt.(py.Int))
	if optimize < -1 || optimize > 2 {
		return nil, py.ExceptionNewf(py.ValueError, "compile(): invalid optimize value")
	}

	if dont_inherit.(py.Int) != 0 {
		// PyEval_MergeCompilerFlags(&cf)
	}

	// switch string(startstr.(py.String)) {
	// case "exec":
	// 	mode = 0
	// case "eval":
	// 	mode = 1
	// case "single":
	// 	mode = 2
	// default:
	// 	return nil, py.ExceptionNewf(py.ValueError, "compile() arg 3 must be 'exec', 'eval' or 'single'")
	// }

	// is_ast = PyAST_Check(cmd)
	// if is_ast {
	// 	if supplied_flags & PyCF_ONLY_AST {
	// 		result = cmd
	// 	} else {

	// 		arena := PyArena_New()
	// 		mod := PyAST_obj2mod(cmd, arena, mode)
	// 		PyAST_Validate(mod)
	// 		result = PyAST_CompileObject(mod, filename, &cf, optimize, arena)
	// 		PyArena_Free(arena)
	// 	}
	// } else {
	str, err := source_as_string(cmd, "compile", "string, bytes or AST" /*, &cf*/)
	if err != nil {
		return nil, err
	}
	// result = py.CompileStringExFlags(str, filename, start[mode], &cf, optimize)
	result, err = compile.Compile(str, string(filename.(py.String)), string(startstr.(py.String)), int(supplied_flags.(py.Int)), dont_inherit.(py.Int) != 0)
	if err != nil {
		return nil, err
	}
	// }

	return result, nil
}

const divmod_doc = `divmod(x, y) -> (quotient, remainder)

Return the tuple ((x-x%y)/y, x%y).  Invariant: div*y + mod == x.`

func builtin_divmod(self py.Object, args py.Tuple) (py.Object, error) {
	var x, y py.Object
	err := py.UnpackTuple(args, nil, "divmod", 2, 2, &x, &y)
	if err != nil {
		return nil, err
	}
	q, r, err := py.DivMod(x, y)
	if err != nil {
		return nil, err
	}
	return py.Tuple{q, r}, nil
}

const eval_doc = `"eval(source[, globals[, locals]]) -> value

Evaluate the source in the context of globals and locals.
The source may be a string representing a Python expression
or a code object as returned by compile().
The globals must be a dictionary and locals can be any mapping,
defaulting to the current globals and locals.
If only globals is given, locals defaults to it.`

// For code see vm/builtin.go

const exec_doc = `exec(object[, globals[, locals]])

Read and execute code from an object, which can be a string or a code
object.
The globals and locals are dictionaries, defaulting to the current
globals and locals.  If only globals is given, locals defaults to it.`

const hex_doc = `hex(number) -> string

Return the hexadecimal representation of an integer.

   >>> hex(12648430)
   '0xc0ffee'
`

func builtin_hex(self, v py.Object) (py.Object, error) {
	var (
		i   int64
		err error
	)
	switch v := v.(type) {
	case *py.BigInt:
		// test bigint first to make sure we correctly handle the case
		// where int64 isn't large enough.
		vv := (*big.Int)(v)
		format := "%#x"
		if vv.Cmp(big.NewInt(0)) == -1 {
			format = "%+#x"
		}
		str := fmt.Sprintf(format, vv)
		return py.String(str), nil
	case py.IGoInt64:
		i, err = v.GoInt64()
	case py.IGoInt:
		var vv int
		vv, err = v.GoInt()
		i = int64(vv)
	default:
		return nil, py.ExceptionNewf(py.TypeError, "'%s' object cannot be interpreted as an integer", v.Type().Name)
	}

	if err != nil {
		return nil, err
	}

	format := "%#x"
	if i < 0 {
		format = "%+#x"
	}
	str := fmt.Sprintf(format, i)
	return py.String(str), nil
}

const isinstance_doc = `isinstance(obj, class_or_tuple) -> bool

Return whether an object is an instance of a class or of a subclass thereof.

A tuple, as in isinstance(x, (A, B, ...)), may be given as the target to
check against. This is equivalent to isinstance(x, A) or isinstance(x, B)
or ... etc.
`

func isinstance(obj py.Object, classOrTuple py.Object) (py.Bool, error) {
	switch classOrTuple.(type) {
	case py.Tuple:
		var class_tuple = classOrTuple.(py.Tuple)
		for idx := range class_tuple {
			res, _ := isinstance(obj, class_tuple[idx])
			if res {
				return res, nil
			}
		}
		return false, nil
	default:
		if classOrTuple.Type().ObjectType != py.TypeType {
			return false, py.ExceptionNewf(py.TypeError, "isinstance() arg 2 must be a type or tuple of types")
		}
		return obj.Type() == classOrTuple, nil
	}
}

func builtin_isinstance(self py.Object, args py.Tuple) (py.Object, error) {
	var obj py.Object
	var classOrTuple py.Object
	err := py.UnpackTuple(args, nil, "isinstance", 2, 2, &obj, &classOrTuple)
	if err != nil {
		return nil, err
	}

	return isinstance(obj, classOrTuple)
}

const iter_doc = `iter(iterable) -> iterator
iter(callable, sentinel) -> iterator

Get an iterator from an object.  In the first form, the argument must
supply its own iterator, or be a sequence.
In the second form, the callable is called until it returns the sentinel.
`

func builtin_iter(self py.Object, args py.Tuple) (py.Object, error) {
	nArgs := len(args)
	if nArgs < 1 {
		return nil, py.ExceptionNewf(py.TypeError,
			"iter expected at least 1 arguments, got %d",
			nArgs)
	} else if nArgs > 2 {
		return nil, py.ExceptionNewf(py.TypeError,
			"iter expected at most 2 arguments, got %d",
			nArgs)
	}

	v := args[0]
	if nArgs == 1 {
		return py.Iter(v)
	}
	_, ok := v.(*py.Function)
	sentinel := args[1]
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError,
			"iter(v, w): v must be callable")
	}
	return py.NewCallIterator(v, sentinel), nil
}

// For code see vm/builtin.go

const len_doc = `len(object) -> integer

Return the number of items of a sequence or mapping.`

func builtin_len(self, v py.Object) (py.Object, error) {
	return py.Len(v)
}

const max_doc = `
max(iterable, *[, default=obj, key=func]) -> value
max(arg1, arg2, *args, *[, key=func]) -> value

With a single iterable argument, return its biggest item. The
default keyword-only argument specifies an object to return if
the provided iterable is empty.
With two or more arguments, return the largest argument.`

func builtin_max(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	return min_max(args, kwargs, "max")
}

const min_doc = `
min(iterable, *[, default=obj, key=func]) -> value
min(arg1, arg2, *args, *[, key=func]) -> value

With a single iterable argument, return its smallest item. The
default keyword-only argument specifies an object to return if
the provided iterable is empty.
With two or more arguments, return the smallest argument.`

func builtin_min(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	return min_max(args, kwargs, "min")
}

func min_max(args py.Tuple, kwargs py.StringDict, name string) (py.Object, error) {
	kwlist := []string{"key", "default"}
	positional := len(args)
	var format string
	var values py.Object
	var cmp func(a py.Object, b py.Object) (py.Object, error)
	if name == "min" {
		format = "|$OO:min"
		cmp = py.Le
	} else if name == "max" {
		format = "|$OO:max"
		cmp = py.Ge
	}
	var defaultValue py.Object
	var keyFunc py.Object
	var maxVal, maxItem py.Object
	var kf *py.Function

	if positional > 1 {
		values = args
	} else {
		err := py.UnpackTuple(args, nil, name, 1, 1, &values)
		if err != nil {
			return nil, err
		}
	}
	err := py.ParseTupleAndKeywords(nil, kwargs, format, kwlist, &keyFunc, &defaultValue)
	if err != nil {
		return nil, err
	}
	if keyFunc == py.None {
		keyFunc = nil
	}
	if keyFunc != nil {
		var ok bool
		kf, ok = keyFunc.(*py.Function)
		if !ok {
			return nil, py.ExceptionNewf(py.TypeError, "'%s' object is not callable", keyFunc.Type())
		}
	}
	if defaultValue != nil {
		maxItem = defaultValue
		if keyFunc != nil {
			maxVal, err = py.Call(kf, py.Tuple{defaultValue}, nil)
			if err != nil {
				return nil, err
			}
		} else {
			maxVal = defaultValue
		}
	}
	iter, err := py.Iter(values)
	if err != nil {
		return nil, err
	}

	for {
		item, err := py.Next(iter)
		if err != nil {
			if py.IsException(py.StopIteration, err) {
				break
			}
			return nil, err
		}
		if maxVal == nil {
			if keyFunc != nil {
				maxVal, err = py.Call(kf, py.Tuple{item}, nil)
				if err != nil {
					return nil, err
				}
			} else {
				maxVal = item
			}
			maxItem = item
		} else {
			var compareVal py.Object
			if keyFunc != nil {
				compareVal, err = py.Call(kf, py.Tuple{item}, nil)
				if err != nil {
					return nil, err
				}
			} else {
				compareVal = item
			}
			changed, err := cmp(compareVal, maxVal)
			if err != nil {
				return nil, err
			}
			if changed == py.True {
				maxVal = compareVal
				maxItem = item
			}
		}

	}

	if maxItem == nil {
		return nil, py.ExceptionNewf(py.ValueError, "%s() arg is an empty sequence", name)
	}

	return maxItem, nil
}

const chr_doc = `chr(i) -> Unicode character

Return a Unicode string of one character with ordinal i; 0 <= i <= 0x10ffff.`

func builtin_chr(self py.Object, args py.Tuple) (py.Object, error) {
	var xObj py.Object

	err := py.ParseTuple(args, "i:chr", &xObj)
	if err != nil {
		return nil, err
	}

	x := xObj.(py.Int)
	if x < 0 || x >= 0x110000 {
		return nil, py.ExceptionNewf(py.ValueError, "chr() arg not in range(0x110000)")
	}
	buf := make([]byte, 8)
	n := utf8.EncodeRune(buf, rune(x))
	return py.String(buf[:n]), nil
}

const locals_doc = `locals() -> dictionary

Update and return a dictionary containing the current scope's local variables.`

const globals_doc = `globals() -> dictionary

Return the dictionary containing the current scope's global variables.`

const sum_doc = `sum($module, iterable, start=0, /)
--
Return the sum of a \'start\' value (default: 0) plus an iterable of numbers

When the iterable is empty, return the start value.
This function is intended specifically for use with numeric values and may
reject non-numeric types.
`

func builtin_sum(self py.Object, args py.Tuple) (py.Object, error) {
	var seq py.Object
	var start py.Object
	err := py.UnpackTuple(args, nil, "sum", 1, 2, &seq, &start)
	if err != nil {
		return nil, err
	}
	if start == nil {
		start = py.Int(0)
	} else {
		switch start.(type) {
		case py.Bytes:
			return nil, py.ExceptionNewf(py.TypeError, "sum() can't sum bytes [use b''.join(seq) instead]")
		case py.String:
			return nil, py.ExceptionNewf(py.TypeError, "sum() can't sum strings [use ''.join(seq) instead]")
		}
	}

	iter, err := py.Iter(seq)
	if err != nil {
		return nil, err
	}

	for {
		item, err := py.Next(iter)
		if err != nil {
			if py.IsException(py.StopIteration, err) {
				break
			}
			return nil, err
		}
		start, err = py.Add(start, item)
		if err != nil {
			return nil, err
		}
	}
	return start, nil
}

const sorted_doc = `sorted(iterable, key=None, reverse=False)

Return a new list containing all items from the iterable in ascending order.

A custom key function can be supplied to customize the sort order, and the
reverse flag can be set to request the result in descending order.`

func builtin_sorted(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	const funcName = "sorted"
	var iterable py.Object
	err := py.UnpackTuple(args, nil, funcName, 1, 1, &iterable)
	if err != nil {
		return nil, err
	}
	l, err := py.SequenceList(iterable)
	if err != nil {
		return nil, err
	}
	err = py.SortInPlace(l, kwargs, funcName)
	if err != nil {
		return nil, err
	}
	return l, nil
}
