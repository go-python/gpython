// Exception objects

package py

import (
	"fmt"
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

var (
	// Exception heirachy
	BaseException             = ObjectType.NewType("BaseException", "Common base class for all exceptions", ExceptionNew, nil)
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
	SystemError               = ExceptionType.NewType("SystemError", "Internal error in the Python interpreter.\n\nPlease report this to the Python maintainer, along with the traceback,\nthe Python version, and the hardware/OS platform and version.", nil, nil)
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
	NotImplemented = NotImplementedError.M__call__(nil, nil)
)

// Type of this object
func (o *Exception) Type() *Type {
	return o.Base
}

// RangeNew
func ExceptionNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	if len(kwargs) != 0 {
		// TypeError
		panic(fmt.Sprintf("TypeError: %s does not take keyword arguments", metatype.Name))
	}
	return &Exception{
		Base: metatype,
		Args: args.Copy(),
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
