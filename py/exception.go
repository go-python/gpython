// Exception objects

package py

import (
	"fmt"
	"io"
)

// A python Exception object
type Exception struct {
	Base            *Type
	Args            Object
	Traceback       Object
	Context         Object
	Cause           Object
	SuppressContext bool
	Other           StringDict // anything else that we want to stuff in
}

// A python exception info block
type ExceptionInfo struct {
	Type      *Type
	Value     Object
	Traceback *Traceback
}

// Make Exception info statisfy the error interface

var (
	// Exception heirachy
	BaseException             = ObjectType.NewTypeFlags("BaseException", "Common base class for all exceptions", ExceptionNew, nil, ObjectType.Flags|TPFLAGS_BASE_EXC_SUBCLASS)
	SystemExit                = BaseException.NewType("SystemExit", "Request to exit from the interpreter.", nil, nil)
	KeyboardInterrupt         = BaseException.NewType("KeyboardInterrupt", "Program interrupted by user.", nil, nil)
	GeneratorExit             = BaseException.NewType("GeneratorExit", "Request that a generator exit.", nil, nil)
	ExceptionType             = BaseException.NewType("Exception", "Common base class for all non-exit exceptions.", nil, nil)
	StopIteration             = ExceptionType.NewType("StopIteration", "Signal the end from iterator.__next__().", nil, nil)
	ArithmeticError           = ExceptionType.NewType("ArithmeticError", "Base class for arithmetic errors.", nil, nil)
	FloatingPointError        = ArithmeticError.NewType("FloatingPointError", "Floating point operation failed.", nil, nil)
	OverflowError             = ArithmeticError.NewType("OverflowError", "Result too large to be represented.", nil, nil)
	ZeroDivisionError         = ArithmeticError.NewType("ZeroDivisionError", "Second argument to a division or modulo operation was zero.", nil, nil)
	AssertionError            = ExceptionType.NewType("AssertionError", "Assertion failed.", nil, nil)
	AttributeError            = ExceptionType.NewType("AttributeError", "Attribute not found.", nil, nil)
	BufferError               = ExceptionType.NewType("BufferError", "Buffer error.", nil, nil)
	EOFError                  = ExceptionType.NewType("EOFError", "Read beyond end of file.", nil, nil)
	ImportError               = ExceptionType.NewType("ImportError", "Import can't find module, or can't find name in module.", nil, nil)
	LookupError               = ExceptionType.NewType("LookupError", "Base class for lookup errors.", nil, nil)
	IndexError                = LookupError.NewType("IndexError", "Sequence index out of range.", nil, nil)
	KeyError                  = LookupError.NewType("KeyError", "Mapping key not found.", nil, nil)
	MemoryError               = ExceptionType.NewType("MemoryError", "Out of memory.", nil, nil)
	NameError                 = ExceptionType.NewType("NameError", "Name not found globally.", nil, nil)
	UnboundLocalError         = NameError.NewType("UnboundLocalError", "Local name referenced but not bound to a value.", nil, nil)
	OSError                   = ExceptionType.NewType("OSError", "Base class for I/O related errors.", nil, nil)
	BlockingIOError           = OSError.NewType("BlockingIOError", "I/O operation would block.", nil, nil)
	ChildProcessError         = OSError.NewType("ChildProcessError", "Child process error.", nil, nil)
	ConnectionError           = OSError.NewType("ConnectionError", "Connection error.", nil, nil)
	BrokenPipeError           = ConnectionError.NewType("BrokenPipeError", "Broken pipe.", nil, nil)
	ConnectionAbortedError    = ConnectionError.NewType("ConnectionAbortedError", "Connection aborted.", nil, nil)
	ConnectionRefusedError    = ConnectionError.NewType("ConnectionRefusedError", "Connection refused.", nil, nil)
	ConnectionResetError      = ConnectionError.NewType("ConnectionResetError", "Connection reset.", nil, nil)
	FileExistsError           = OSError.NewType("FileExistsError", "File already exists.", nil, nil)
	FileNotFoundError         = OSError.NewType("FileNotFoundError", "File not found.", nil, nil)
	InterruptedError          = OSError.NewType("InterruptedError", "Interrupted by signal.", nil, nil)
	IsADirectoryError         = OSError.NewType("IsADirectoryError", "Operation doesn't work on directories.", nil, nil)
	NotADirectoryError        = OSError.NewType("NotADirectoryError", "Operation only works on directories.", nil, nil)
	PermissionError           = OSError.NewType("PermissionError", "Not enough permissions.", nil, nil)
	ProcessLookupError        = OSError.NewType("ProcessLookupError", "Process not found.", nil, nil)
	TimeoutError              = OSError.NewType("TimeoutError", "Timeout expired.", nil, nil)
	ReferenceError            = ExceptionType.NewType("ReferenceError", "Weak ref proxy used after referent went away.", nil, nil)
	RuntimeError              = ExceptionType.NewType("RuntimeError", "Unspecified run-time error.", nil, nil)
	NotImplementedError       = RuntimeError.NewType("NotImplementedError", "Method or function hasn't been implemented yet.", nil, nil)
	SyntaxError               = ExceptionType.NewType("SyntaxError", "Invalid syntax.", nil, nil)
	IndentationError          = SyntaxError.NewType("IndentationError", "Improper indentation.", nil, nil)
	TabError                  = IndentationError.NewType("TabError", "Improper mixture of spaces and tabs.", nil, nil)
	SystemError               = ExceptionType.NewType("SystemError", "Internal error in the Gpython interpreter.\n\nPlease report this to the Gpython maintainer, along with the traceback,\nthe Gpython version, and the hardware/OS platform and version.", nil, nil)
	TypeError                 = ExceptionType.NewType("TypeError", "Inappropriate argument type.", nil, nil)
	ValueError                = ExceptionType.NewType("ValueError", "Inappropriate argument value (of correct type).", nil, nil)
	UnicodeError              = ValueError.NewType("UnicodeError", "Unicode related error.", nil, nil)
	UnicodeDecodeError        = UnicodeError.NewType("UnicodeDecodeError", "Unicode decoding error.", nil, nil)
	UnicodeEncodeError        = UnicodeError.NewType("UnicodeEncodeError", "Unicode encoding error.", nil, nil)
	UnicodeTranslateError     = UnicodeError.NewType("UnicodeTranslateError", "Unicode translation error.", nil, nil)
	Warning                   = ExceptionType.NewType("Warning", "Base class for warning categories.", nil, nil)
	DeprecationWarning        = Warning.NewType("DeprecationWarning", "Base class for warnings about deprecated features.", nil, nil)
	PendingDeprecationWarning = Warning.NewType("PendingDeprecationWarning", "Base class for warnings about features which will be deprecated\nin the future.", nil, nil)
	RuntimeWarning            = Warning.NewType("RuntimeWarning", "Base class for warnings about dubious runtime behavior.", nil, nil)
	SyntaxWarning             = Warning.NewType("SyntaxWarning", "Base class for warnings about dubious syntax.", nil, nil)
	UserWarning               = Warning.NewType("UserWarning", "Base class for warnings generated by user code.", nil, nil)
	FutureWarning             = Warning.NewType("FutureWarning", "Base class for warnings about constructs that will change semantically\nin the future.", nil, nil)
	ImportWarning             = Warning.NewType("ImportWarning", "Base class for warnings about probable mistakes in module imports", nil, nil)
	UnicodeWarning            = Warning.NewType("UnicodeWarning", "Base class for warnings about Unicode related problems, mostly\nrelated to conversion problems.", nil, nil)
	BytesWarning              = Warning.NewType("BytesWarning", "Base class for warnings about bytes and buffer related problems, mostly\nrelated to conversion from str or comparing to str.", nil, nil)
	ResourceWarning           = Warning.NewType("ResourceWarning", "Base class for warnings about resource usage.", nil, nil)
	// Singleton exceptions
	NotImplemented = ExceptionNew(NotImplementedError, nil, nil)
)

// Type of this object
func (e *Exception) Type() *Type {
	return e.Base
}

// Go error interface
func (e *Exception) Error() string {
	return fmt.Sprintf("%s: %v", e.Base.Name, e.Args)
}

// Go error interface
func (e ExceptionInfo) Error() string {
	if exception, ok := e.Value.(*Exception); ok {
		return exception.Error()
	}
	return e.Value.Type().Name
}

// Dump a traceback for exc to w
func (exc *ExceptionInfo) TracebackDump(w io.Writer) {
	fmt.Fprintf(w, "Traceback (most recent call last):\n")
	exc.Traceback.TracebackDump(w)
	fmt.Fprintf(w, "%v: %v\n", exc.Type.Name, exc.Value)
}

// Test for being set
func (exc *ExceptionInfo) IsSet() bool {
	return exc.Type != nil
}

// ExceptionNew
func ExceptionNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	if len(kwargs) != 0 {
		// FIXME this causes an initialization loop
		//panic(ExceptionNewf(TypeError, "%s does not take keyword arguments", metatype.Name))
		panic(fmt.Sprintf("TypeError: %s does not take keyword arguments", metatype.Name))
	}
	return &Exception{
		Base: metatype,
		Args: args.Copy(),
	}
}

// ExceptionNewf - make a new exception with fmt parameters
func ExceptionNewf(metatype *Type, format string, a ...interface{}) *Exception {
	message := fmt.Sprintf(format, a...)
	return &Exception{
		Base: metatype,
		Args: Tuple{String(message)},
	}
}

/*
	if py.ExceptionClassCheck(exc) {
		t = exc.(*py.Type)
		value = py.Call(exc, nil, nil)
		if value == nil {
			return exitException
		}
		if !py.ExceptionInstanceCheck(value) {
			PyErr_Format(PyExc_TypeError, "calling %s should have returned an instance of BaseException, not %s", t.Name, value.Type().Name)
			return exitException
		}
	} else if t = py.ExceptionInstanceCheck(exc); t != nil {
		value = exc
	} else {
		// Not something you can raise.  You get an exception
		// anyway, just not what you specified :-)
		PyErr_SetString(PyExc_TypeError, "exceptions must derive from BaseException")
		return exitException
	}
*/

// Coerce an object into an exception instance one way or another
func MakeException(r interface{}) *Exception {
	switch x := r.(type) {
	case *Exception:
		return x
	case *Type:
		if x.Flags&TPFLAGS_BASE_EXC_SUBCLASS != 0 {
			return ExceptionNew(x, nil, nil).(*Exception)
		} else {
			return ExceptionNewf(TypeError, "exceptions must derive from BaseException")
		}
	case error:
		return ExceptionNew(SystemError, Tuple{String(x.Error())}, nil).(*Exception)
	case string:
		return ExceptionNew(SystemError, Tuple{String(x)}, nil).(*Exception)
	default:
		return ExceptionNew(SystemError, Tuple{String(fmt.Sprintf("Unknown error %#v", r))}, nil).(*Exception)
	}
}

/*
#define PyType_HasFeature(t,f)  (((t)->tp_flags & (f)) != 0)

#define PyType_FastSubclass(t,f)  PyType_HasFeature(t,f)

#define PyType_Check(op) \
    PyType_FastSubclass(Py_TYPE(op), Py_TPFLAGS_TYPE_SUBCLASS)

#define PyType_CheckExact(op) (Py_TYPE(op) == &PyType_Type)

#define PyExceptionClass_Check(x)                                       \
    (PyType_Check((x)) &&                                               \
     PyType_FastSubclass((PyTypeObject*)(x), Py_TPFLAGS_BASE_EXC_SUBCLASS))

#define PyExceptionInstance_Check(x)                    \
    PyType_FastSubclass((x)->ob_type, Py_TPFLAGS_BASE_EXC_SUBCLASS)

#define PyExceptionClass_Name(x) \
     ((char *)(((PyTypeObject*)(x))->tp_name))

#define PyExceptionInstance_Class(x) ((PyObject*)((x)->ob_type))
*/

// Checks that the object passed in is a class and is an exception
func ExceptionClassCheck(err Object) bool {
	if t, ok := err.(*Type); ok {
		// FIXME not telling instances and classes apart
		// properly! This could be an instance of something
		// here
		return t.Flags&TPFLAGS_BASE_EXC_SUBCLASS != 0
	}
	return false
}

// Check to see if err matches exc
//
// exc can be a tuple
//
// Used in except statements
func ExceptionGivenMatches(err, exc Object) bool {
	if err == nil || exc == nil {
		// maybe caused by "import exceptions" that failed early on
		return false
	}

	// Test the tuple case recursively
	if excTuple, ok := exc.(Tuple); ok {
		for i := range excTuple {
			if ExceptionGivenMatches(err, excTuple[i]) {
				return true
			}
		}
		return false
	}

	// err might be an instance, so check its class.
	if exception, ok := err.(*Exception); ok {
		err = exception.Type()
	}

	if ExceptionClassCheck(err) && ExceptionClassCheck(exc) {
		res := false
		// PyObject *exception, *value, *tb;
		// PyErr_Fetch(&exception, &value, &tb);

		// PyObject_IsSubclass() can recurse and therefore is
		// not safe (see test_bad_getattr in test.pickletester).
		res = err.(*Type).IsSubtype(exc.(*Type))
		// This function must not fail, so print the error here
		// if (res == -1) {
		// 	PyErr_WriteUnraisable(err);
		// 	res = false
		// }
		// PyErr_Restore(exception, value, tb);
		return res
	}

	return err == exc
}

// IsException matches the result of recover to an exception
//
// For use to catch a single python exception from go code
//
// It can be an instance or the class itself
func IsException(exception *Type, r interface{}) bool {
	var t *Type
	switch ex := r.(type) {
	case *Exception:
		t = ex.Type()
	case *Type:
		t = ex
	default:
		return false
	}
	// Exact instance or subclass match
	if t == exception {
		return true
	}
	// Can't be a subclass of exception
	if t.Flags&TPFLAGS_BASE_EXC_SUBCLASS == 0 {
		return false
	}
	// Now the full match
	return t.IsSubtype(exception)
}

// Check Interfaces
var _ error = (*Exception)(nil)
var _ error = (*ExceptionInfo)(nil)
