// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// System module
//
// Various bits of information used by the interpreter are collected in
// module 'sys'.
// Function member:
// - exit(sts): raise SystemExit
// Data members:
// - stdin, stdout, stderr: standard file objects
// - modules: the table of modules (dictionary)
// - path: module search path (list of strings)
// - argv: script arguments (list of strings)
// - ps1, ps2: optional primary and secondary prompts (strings)

package sys

import (
	"os"

	"github.com/go-python/gpython/py"
)

const module_doc = `This module provides access to some objects used or maintained by the
interpreter and to functions that interact strongly with the interpreter.

Dynamic objects:

argv -- command line arguments; argv[0] is the script pathname if known
path -- module search path; path[0] is the script directory, else ''
modules -- dictionary of loaded modules

displayhook -- called to show results in an interactive session
excepthook -- called to handle any uncaught exception other than SystemExit
  To customize printing in an interactive session or to install a custom
  top-level exception handler, assign other functions to replace these.

stdin -- standard input file object; used by input()
stdout -- standard output file object; used by print()
stderr -- standard error object; used for error messages
  By assigning other file objects (or objects that behave like files)
  to these, it is possible to redirect all of the interpreter's I/O.

last_type -- type of last uncaught exception
last_value -- value of last uncaught exception
last_traceback -- traceback of last uncaught exception
  These three are only available in an interactive session after a
  traceback has been printed.

 objects:

builtin_module_names -- tuple of module names built into this interpreter
copyright -- copyright notice pertaining to this interpreter
exec_prefix -- prefix used to find the machine-specific Python library
executable -- absolute path of the executable binary of the Python interpreter
float_info -- a struct sequence with information about the float implementation.
float_repr_style -- string indicating the style of repr() output for floats
hexversion -- version information encoded as a single integer
implementation -- Python implementation information.
int_info -- a struct sequence with information about the int implementation.
maxsize -- the largest supported length of containers.
maxunicode -- the value of the largest Unicode codepoint
platform -- platform identifier
prefix -- prefix used to find the Python library
thread_info -- a struct sequence with information about the thread implementation.
version -- the version of this interpreter as a string
version_info -- version information as a named tuple
__stdin__ -- the original stdin; don't touch!
__stdout__ -- the original stdout; don't touch!
__stderr__ -- the original stderr; don't touch!
__displayhook__ -- the original displayhook; don't touch!
__excepthook__ -- the original excepthook; don't touch!

Functions:

displayhook() -- print an object to the screen, and save it in builtins._
excepthook() -- print an exception and its traceback to sys.stderr
exc_info() -- return thread-safe information about the current exception
exit() -- exit the interpreter by raising SystemExit
getdlopenflags() -- returns flags to be used for dlopen() calls
getprofile() -- get the global profiling function
getrefcount() -- return the reference count for an object (plus one :-)
getrecursionlimit() -- return the max recursion depth for the interpreter
getsizeof() -- return the size of an object in bytes
gettrace() -- get the global debug tracing function
setcheckinterval() -- control how often the interpreter checks for events
setdlopenflags() -- set the flags to be used for dlopen() calls
setprofile() -- set the global profiling function
setrecursionlimit() -- set the max recursion depth for the interpreter
settrace() -- set the global debug tracing function
`

const displayhook_doc = `displayhook(object) -> None

Print an object to sys.stdout and also save it in builtins._`

func sys_displayhook(self, o py.Object) (py.Object, error) {
	return nil, py.NotImplementedError
}

const excepthook_doc = `excepthook(exctype, value, traceback) -> None

Handle an exception by displaying it with a traceback on sys.stderr.`

func sys_excepthook(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
}

const exc_info_doc = `exc_info() -> (type, value, traceback)

Return information about the most recent exception caught by an except
clause in the current stack frame or in an older stack frame.`

func sys_exc_info(self py.Object) (py.Object, error) {
	return nil, py.NotImplementedError
}

const exit_doc = `exit([status])

Exit the interpreter by raising SystemExit(status).
If the status is omitted or None, it defaults to zero (i.e., success).
If the status is an integer, it will be used as the system exit status.
If it is another kind of object, it will be printed and the system
exit status will be one (i.e., failure).`

func sys_exit(self py.Object, args py.Tuple) (py.Object, error) {
	var exit_code py.Object
	err := py.UnpackTuple(args, nil, "exit", 0, 1, &exit_code)
	if err != nil {
		return nil, err
	}
	// Raise SystemExit so callers may catch it or clean up.
	return py.ExceptionNew(py.SystemExit, args, nil)
}

const getdefaultencoding_doc = `getdefaultencoding() -> string

Return the current default string encoding used by the Unicode 
implementation.`

func sys_getdefaultencoding(self py.Object) (py.Object, error) {
	return nil, py.NotImplementedError
	// return PyUnicode_FromString(PyUnicode_GetDefaultEncoding());
}

const getfilesystemencoding_doc = `getfilesystemencoding() -> string

Return the encoding used to convert Unicode filenames in
operating system filenames.`

func sys_getfilesystemencoding(self py.Object) (py.Object, error) {
	return nil, py.NotImplementedError
	// if (Py_FileSystemDefaultEncoding) {
	//     return PyUnicode_FromString(Py_FileSystemDefaultEncoding);
	// }
	// PyErr_SetString(PyExc_RuntimeError,
	//                 "filesystem encoding is not initialized");
	// return nil;
}

const intern_doc = `intern(string) -> string

"Intern" the given string.  This enters the string in the (global)
table of interned strings whose purpose is to speed up dictionary lookups.
Return the string itself or the previously interned string object with the
same value.`

func sys_intern(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// py.Object s;
	// if (!PyArg_ParseTuple(args, "U:intern", &s)) {
	//     return nil;
	// }
	// if (PyUnicode_CheckExact(s)) {
	//     Py_INCREF(s);
	//     PyUnicode_InternInPlace(&s);
	//     return s;
	// } else {
	//     PyErr_Format(PyExc_TypeError,
	//                  "can't intern %.400s", s->ob_type->tp_name);
	//     return nil;
	// }
}

const settrace_doc = `settrace(function)

Set the global debug tracing function.  It will be called on each
function call.  See the debugger chapter in the library manual.`

func sys_settrace(self py.Object, args py.Tuple) (py.Object, error) {
	// if (trace_init() == -1) {
	//     return nil;
	// }
	// if (args == Py_None) {
	//     PyEval_SetTrace(nil, nil);
	// } else {
	//     PyEval_SetTrace(trace_trampoline, args);
	// }
	// Py_INCREF(Py_None);
	// return Py_None;
	return nil, py.NotImplementedError
}

const gettrace_doc = `gettrace()

Return the global debug tracing function set with sys.settrace.
See the debugger chapter in the library manual.`

func sys_gettrace(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// PyThreadState *tstate = PyThreadState_GET();
	// py.Object temp = tstate->c_traceobj;

	// if (temp == nil) {
	//     temp = Py_None;
	// }
	// Py_INCREF(temp);
	// return temp;
}

const setprofile_doc = `setprofile(function)

Set the profiling function.  It will be called on each function call
and return.  See the profiler chapter in the library manual.`

func sys_setprofile(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// if (trace_init() == -1) {
	//     return nil;
	// }
	// if (args == Py_None) {
	//     PyEval_SetProfile(nil, nil);
	// } else {
	//     PyEval_SetProfile(profile_trampoline, args);
	// }
	// Py_INCREF(Py_None);
	// return Py_None;
}

const getprofile_doc = `getprofile()

Return the profiling function set with sys.setprofile.
See the profiler chapter in the library manual.`

func sys_getprofile(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// PyThreadState *tstate = PyThreadState_GET();
	// py.Object temp = tstate->c_profileobj;

	// if (temp == nil) {
	//     temp = Py_None;
	// }
	// Py_INCREF(temp);
	// return temp;
}

// int _check_interval = 100;

const setcheckinterval_doc = `setcheckinterval(n)

Tell the Python interpreter to check for asynchronous events every
n instructions.  This also affects how often thread switches occur.`

func sys_setcheckinterval(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// if (PyErr_WarnEx(PyExc_DeprecationWarning,
	//                  "sys.getcheckinterval() and sys.setcheckinterval() "
	//                  "are deprecated.  Use sys.setswitchinterval() "
	//                  "instead.", 1) < 0) {
	//     return nil;
	// }
	// if (!PyArg_ParseTuple(args, "i:setcheckinterval", &_check_interval)) {
	//     return nil;
	// }
	// Py_INCREF(Py_None);
	// return Py_None;
}

const getcheckinterval_doc = `getcheckinterval() -> current check interval; see setcheckinterval().`

func sys_getcheckinterval(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// if (PyErr_WarnEx(PyExc_DeprecationWarning,
	//                  "sys.getcheckinterval() and sys.setcheckinterval() "
	//                  "are deprecated.  Use sys.getswitchinterval() "
	//                  "instead.", 1) < 0) {
	//     return nil;
	// }
	// return PyLong_FromLong(_check_interval);
}

const setswitchinterval_doc = `setswitchinterval(n)

Set the ideal thread switching delay inside the Python interpreter
The actual frequency of switching threads can be lower if the
interpreter executes long sequences of uninterruptible code
(this is implementation-specific and workload-dependent).

The parameter must represent the desired switching delay in seconds
A typical value is 0.005 (5 milliseconds).`

func sys_setswitchinterval(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
	// double d;
	// if (!PyArg_ParseTuple(args, "d:setswitchinterval", &d)) {
	//     return nil;
	// }
	// if (d <= 0.0) {
	//     PyErr_SetString(PyExc_ValueError,
	//                     "switch interval must be strictly positive");
	//     return nil;
	// }
	// _PyEval_SetSwitchInterval((unsigned long) (1e6 * d));
	// Py_INCREF(Py_None);
	// return Py_None;
}

const getswitchinterval_doc = `getswitchinterval() -> current thread switch interval; see setswitchinterval().`

func sys_getswitchinterval(self py.Object, args py.Tuple) (py.Object, error) {
	// return PyFloat_FromDouble(1e-6 * _PyEval_GetSwitchInterval());
	return nil, py.NotImplementedError
}

const setrecursionlimit_doc = `setrecursionlimit(n)

Set the maximum depth of the Python interpreter stack to n.  This
limit prevents infinite recursion from causing an overflow of the C
stack and crashing Python.  The highest possible limit is platform-
dependent.`

func sys_setrecursionlimit(self py.Object, args py.Tuple) (py.Object, error) {
	// int new_limit;
	// if (!PyArg_ParseTuple(args, "i:setrecursionlimit", &new_limit)) {
	//     return nil;
	// }
	// if (new_limit <= 0) {
	//     PyErr_SetString(PyExc_ValueError,
	//                     "recursion limit must be positive");
	//     return nil;
	// }
	// Py_SetRecursionLimit(new_limit);
	// Py_INCREF(Py_None);
	// return Py_None;
	return nil, py.NotImplementedError
}

const hash_info_doc = `hash_info

A struct sequence providing parameters used for computing
numeric hashes.  The attributes are read only.`

//  PyStructSequence_Field hash_info_fields[] = {
//     {"width", "width of the type used for hashing, in bits"},
//     {
//         "modulus", "prime number giving the modulus on which the hash "
//         "function is based"
//     },
//     {"inf", "value to be used for hash of a positive infinity"},
//     {"nan", "value to be used for hash of a nan"},
//     {"imag", "multiplier used for the imaginary part of a complex number"},
//     {nil, nil}
// };

//  PyStructSequence_Desc hash_info_desc = {
//     "sys.hash_info",
//     hash_info_doc,
//     hash_info_fields,
//     5,
// };

// get_hash_info(void) {
//     py.Object hash_info;
//     int field = 0;
//     hash_info = PyStructSequence_New(&Hash_InfoType);
//     if (hash_info == nil) {
//         return nil;
//     }
//     PyStructSequence_SET_ITEM(hash_info, field++,
//                               PyLong_FromLong(8*sizeof(Py_hash_t)));
//     PyStructSequence_SET_ITEM(hash_info, field++,
//                               PyLong_FromSsize_t(_PyHASH_MODULUS));
//     PyStructSequence_SET_ITEM(hash_info, field++,
//                               PyLong_FromLong(_PyHASH_INF));
//     PyStructSequence_SET_ITEM(hash_info, field++,
//                               PyLong_FromLong(_PyHASH_NAN));
//     PyStructSequence_SET_ITEM(hash_info, field++,
//                               PyLong_FromLong(_PyHASH_IMAG));
//     if (PyErr_Occurred()) {
//         Py_CLEAR(hash_info);
//         return nil;
//     }
//     return hash_info;
// }

const getrecursionlimit_doc = `getrecursionlimit()

Return the current value of the recursion limit, the maximum depth
of the Python interpreter stack.  This limit prevents infinite
recursion from causing an overflow of the C stack and crashing Python.`

func sys_getrecursionlimit(self py.Object) (py.Object, error) {
	// return PyLong_FromLong(Py_GetRecursionLimit());
	return nil, py.NotImplementedError
}

const getsizeof_doc = `getsizeof(object, default) -> int

Return the size of object in bytes.`

func sys_getsizeof(self py.Object, args py.Tuple, kwds py.StringDict) (py.Object, error) {
	// py.Object res = nil;
	//  py.Object gc_head_size = nil;
	//  char *kwlist[] = {"object", "default", 0};
	// py.Object o, *dflt = nil;
	// py.Object method;
	// _Py_IDENTIFIER(__sizeof__);

	// if (!PyArg_ParseTupleAndKeywords(args, kwds, "O|O:getsizeof",
	//                                  kwlist, &o, &dflt)) {
	//     return nil;
	// }
	return nil, py.NotImplementedError
}

const getrefcount_doc = `getrefcount(object) -> integer

Return the reference count of object.  The count returned is generally
one higher than you might expect, because it includes the (temporary)
reference as an argument to getrefcount().`

func sys_getrefcount(self, arg py.Object) (py.Object, error) {
	return py.Int(2), nil
}

const getframe_doc = `_getframe([depth]) -> frameobject

Return a frame object from the call stack.  If optional integer depth is
given, return the frame object that many calls below the top of the stack.
If that is deeper than the call stack, ValueError is raised.  The default
for depth is zero, returning the frame at the top of the call stack.

This function should be used for internal and specialized
purposes only.`

func sys_getframe(self py.Object, args py.Tuple) (py.Object, error) {
	// PyFrameObject *f = PyThreadState_GET()->frame;
	// int depth = -1;

	// if (!PyArg_ParseTuple(args, "|i:_getframe", &depth)) {
	//     return nil;
	// }

	// while (depth > 0 && f != nil) {
	//     f = f->f_back;
	//     --depth;
	// }
	// if (f == nil) {
	//     PyErr_SetString(PyExc_ValueError,
	//                     "call stack is not deep enough");
	//     return nil;
	// }
	// Py_INCREF(f);
	// return (PyObject*)f;
	return nil, py.NotImplementedError
}

const current_frames_doc = `_current_frames() -> dictionary

Return a dictionary mapping each current thread T's thread id to T's
current stack frame.

This function should be used for specialized purposes only.`

func sys_current_frames(self py.Object) (py.Object, error) {
	// return _PyThread_CurrentFrames();
	return nil, py.NotImplementedError
}

const call_tracing_doc = `call_tracing(func, args) -> object

Call func(*args), while tracing is enabled.  The tracing state is
saved, and restored afterwards.  This is intended to be called from
a debugger from a checkpoint, to recursively debug some other code.`

func sys_call_tracing(self py.Object, args py.Tuple) (py.Object, error) {
	// py.Object func, *funcargs;
	// if (!PyArg_ParseTuple(args, "OO!:call_tracing", &func, &PyTuple_Type, &funcargs)) {
	//     return nil;
	// }
	// return _PyEval_CallTracing(func, funcargs);
	return nil, py.NotImplementedError
}

const callstats_doc = `callstats() -> tuple of integers

Return a tuple of function call statistics, if CALL_PROFILE was defined
when Python was built.  Otherwise, return None.

When enabled, this function returns detailed, implementation-specific
details about the number of function calls executed. The return value is
a 11-tuple where the entries in the tuple are counts of:
0. all function calls
1. calls to PyFunction_Type objects
2. PyFunction calls that do not create an argument tuple
3. PyFunction calls that do not create an argument tuple
   and bypass PyEval_EvalCodeEx()
4. PyMethod calls
5. PyMethod calls on bound methods
6. PyType calls
7. PyCFunction calls
8. generator calls
9. All other calls
10. Number of stack pops performed by call_function()`

func sys_callstats(self py.Object, args py.Tuple) (py.Object, error) {
	return py.None, nil
}

const debugmallocstats_doc = `_debugmallocstats()

Print summary info to stderr about the state of
pymalloc's structures.

In Py_DEBUG mode, also perform some expensive internal consistency
checks.`

func sys_debugmallocstats(self py.Object, args py.Tuple) (py.Object, error) {
	return nil, py.NotImplementedError
}

const sys_clear_type_cache__doc__ = `_clear_type_cache() -> None
Clear the internal type lookup cache.`

func sys_clear_type_cache(self py.Object, args py.Tuple) (py.Object, error) {
	// PyType_ClearCache()
	return nil, py.NotImplementedError
}

const flags__doc__ = `sys.flags

Flags provided through command line arguments or environment vars.`

//  PyTypeObject FlagsType;

//  PyStructSequence_Field flags_fields[] = {
//     {"debug",                   "-d"},
//     {"inspect",                 "-i"},
//     {"interactive",             "-i"},
//     {"optimize",                "-O or -OO"},
//     {"dont_write_bytecode",     "-B"},
//     {"no_user_site",            "-s"},
//     {"no_site",                 "-S"},
//     {"ignore_environment",      "-E"},
//     {"verbose",                 "-v"},
//     /* {"unbuffered",                   "-u"}, */
//     /* {"skip_first",                   "-x"}, */
//     {"bytes_warning",           "-b"},
//     {"quiet",                   "-q"},
//     {"hash_randomization",      "-R"},
//     {0}
// };

//  PyStructSequence_Desc flags_desc = {
//     "sys.flags",        /* name */
//     flags__doc__,       /* doc */
//     flags_fields,       /* fields */
// };

//  PyObject*
// make_flags(void) {
//     int pos = 0;
//     py.Object seq;

//     seq = PyStructSequence_New(&FlagsType);
//     if (seq == nil) {
//         return nil;
//     }

// #define SetFlag(flag) \
//     PyStructSequence_SET_ITEM(seq, pos++, PyLong_FromLong(flag))

//     SetFlag(Py_DebugFlag);
//     SetFlag(Py_InspectFlag);
//     SetFlag(Py_InteractiveFlag);
//     SetFlag(Py_OptimizeFlag);
//     SetFlag(Py_DontWriteBytecodeFlag);
//     SetFlag(Py_NoUserSiteDirectory);
//     SetFlag(Py_NoSiteFlag);
//     SetFlag(Py_IgnoreEnvironmentFlag);
//     SetFlag(Py_VerboseFlag);
//     /* SetFlag(saw_unbuffered_flag); */
//     /* SetFlag(skipfirstline); */
//     SetFlag(Py_BytesWarningFlag);
//     SetFlag(Py_QuietFlag);
//     SetFlag(Py_HashRandomizationFlag);

//     if (PyErr_Occurred()) {
//         return nil;
//     }
//     return seq;
// }

const version_info__doc__ = `sys.version_info

Version information as a named tuple.`

//  PyStructSequence_Field version_info_fields[] = {
//     {"major", "Major release number"},
//     {"minor", "Minor release number"},
//     {"micro", "Patch release number"},
//     {"releaselevel", "'alpha', 'beta', 'candidate', or 'release'"},
//     {"serial", "Serial release number"},
//     {0}
// };

//  PyStructSequence_Desc version_info_desc = {
//     "sys.version_info",     /* name */
//     version_info__doc__,    /* doc */
//     version_info_fields,    /* fields */
//     5
// };

// Initialise the module
func init() {
	methods := []*py.Method{
		py.MustNewMethod("callstats", sys_callstats, 0, callstats_doc),
		py.MustNewMethod("_clear_type_cache", sys_clear_type_cache, 0, sys_clear_type_cache__doc__),
		py.MustNewMethod("_current_frames", sys_current_frames, 0, current_frames_doc),
		py.MustNewMethod("displayhook", sys_displayhook, 0, displayhook_doc),
		py.MustNewMethod("exc_info", sys_exc_info, 0, exc_info_doc),
		py.MustNewMethod("excepthook", sys_excepthook, 0, excepthook_doc),
		py.MustNewMethod("exit", sys_exit, 0, exit_doc),
		py.MustNewMethod("getdefaultencoding", sys_getdefaultencoding, 0, getdefaultencoding_doc),
		py.MustNewMethod("getfilesystemencoding", sys_getfilesystemencoding, 0, getfilesystemencoding_doc),
		py.MustNewMethod("getrefcount", sys_getrefcount, 0, getrefcount_doc),
		py.MustNewMethod("getrecursionlimit", sys_getrecursionlimit, 0, getrecursionlimit_doc),
		py.MustNewMethod("getsizeof", sys_getsizeof, 0, getsizeof_doc),
		py.MustNewMethod("_getframe", sys_getframe, 0, getframe_doc),
		py.MustNewMethod("intern", sys_intern, 0, intern_doc),
		py.MustNewMethod("setcheckinterval", sys_setcheckinterval, 0, setcheckinterval_doc),
		py.MustNewMethod("getcheckinterval", sys_getcheckinterval, 0, getcheckinterval_doc),
		py.MustNewMethod("setswitchinterval", sys_setswitchinterval, 0, setswitchinterval_doc),
		py.MustNewMethod("getswitchinterval", sys_getswitchinterval, 0, getswitchinterval_doc),
		py.MustNewMethod("setprofile", sys_setprofile, 0, setprofile_doc),
		py.MustNewMethod("getprofile", sys_getprofile, 0, getprofile_doc),
		py.MustNewMethod("setrecursionlimit", sys_setrecursionlimit, 0, setrecursionlimit_doc),
		py.MustNewMethod("settrace", sys_settrace, 0, settrace_doc),
		py.MustNewMethod("gettrace", sys_gettrace, 0, gettrace_doc),
		py.MustNewMethod("call_tracing", sys_call_tracing, 0, call_tracing_doc),
		py.MustNewMethod("_debugmallocstats", sys_debugmallocstats, 0, debugmallocstats_doc),
	}

	stdin := &py.File{File: os.Stdin, FileMode: py.FileRead}
	stdout := &py.File{File: os.Stdout, FileMode: py.FileWrite}
	stderr := &py.File{File: os.Stderr, FileMode: py.FileWrite}

	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	globals := py.StringDict{
		"path":       py.NewList(),
		"argv":       py.NewListFromStrings(os.Args[1:]),
		"stdin":      stdin,
		"stdout":     stdout,
		"stderr":     stderr,
		"__stdin__":  stdin,
		"__stdout__": stdout,
		"__stderr__": stderr,
		"executable": py.String(executable),

		//"version": py.Int(MARSHAL_VERSION),
		//     /* stdin/stdout/stderr are now set by pythonrun.c */

		//     PyDict_SetItemString(sysdict, "__displayhook__",
		//                          PyDict_GetItemString(sysdict, "displayhook"));
		//     PyDict_SetItemString(sysdict, "__excepthook__",
		//                          PyDict_GetItemString(sysdict, "excepthook"));
		//     SET_SYS_FROM_STRING("version",
		//                         PyUnicode_FromString(Py_GetVersion()));
		//     SET_SYS_FROM_STRING("hexversion",
		//                         PyLong_FromLong(PY_VERSION_HEX));
		//     SET_SYS_FROM_STRING("_mercurial",
		//                         Py_BuildValue("(szz)", "CPython", _Py_hgidentifier(),
		//                                       _Py_hgversion()));
		//     SET_SYS_FROM_STRING("dont_write_bytecode",
		//                         PyBool_FromLong(Py_DontWriteBytecodeFlag));
		//     SET_SYS_FROM_STRING("api_version",
		//                         PyLong_FromLong(PYTHON_API_VERSION));
		//     SET_SYS_FROM_STRING("copyright",
		//                         PyUnicode_FromString(Py_GetCopyright()));
		//     SET_SYS_FROM_STRING("platform",
		//                         PyUnicode_FromString(Py_GetPlatform()));
		//     SET_SYS_FROM_STRING("executable",
		//                         PyUnicode_FromWideChar(
		//                             Py_GetProgramFullPath(), -1));
		//     SET_SYS_FROM_STRING("prefix",
		//                         PyUnicode_FromWideChar(Py_GetPrefix(), -1));
		//     SET_SYS_FROM_STRING("exec_prefix",
		//                         PyUnicode_FromWideChar(Py_GetExecPrefix(), -1));
		//     SET_SYS_FROM_STRING("base_prefix",
		//                         PyUnicode_FromWideChar(Py_GetPrefix(), -1));
		//     SET_SYS_FROM_STRING("base_exec_prefix",
		//                         PyUnicode_FromWideChar(Py_GetExecPrefix(), -1));
		//     SET_SYS_FROM_STRING("maxsize",
		//                         PyLong_FromSsize_t(PY_SSIZE_T_MAX));
		//     SET_SYS_FROM_STRING("float_info",
		//                         PyFloat_GetInfo());
		//     SET_SYS_FROM_STRING("int_info",
		//                         PyLong_GetInfo());
		//     /* initialize hash_info */
		//     if (Hash_InfoType.tp_name == 0) {
		//         PyStructSequence_InitType(&Hash_InfoType, &hash_info_desc);
		//     }
		//     SET_SYS_FROM_STRING("hash_info",
		//                         get_hash_info());
		//     SET_SYS_FROM_STRING("maxunicode",
		//                         PyLong_FromLong(0x10FFFF));
		//     SET_SYS_FROM_STRING("builtin_module_names",
		//                         list_builtin_module_names());
		//     {
		//         /* Assumes that longs are at least 2 bytes long.
		//            Should be safe! */
		//         unsigned long number = 1;
		//         char *value;

		//         s = (char *) &number;
		//         if (s[0] == 0) {
		//             value = "big";
		//         } else {
		//             value = "little";
		//         }
		//         SET_SYS_FROM_STRING("byteorder",
		//                             PyUnicode_FromString(value));
		//     }
		// #ifdef MS_COREDLL
		//     SET_SYS_FROM_STRING("dllhandle",
		//                         PyLong_FromVoidPtr(PyWin_DLLhModule));
		//     SET_SYS_FROM_STRING("winver",
		//                         PyUnicode_FromString(PyWin_DLLVersionString));
		// #endif
		// #ifdef ABIFLAGS
		//     SET_SYS_FROM_STRING("abiflags",
		//                         PyUnicode_FromString(ABIFLAGS));
		// #endif
		//     if (warnoptions == nil) {
		//         warnoptions = PyList_New(0);
		//     } else {
		//         Py_INCREF(warnoptions);
		//     }
		//     if (warnoptions != nil) {
		//         PyDict_SetItemString(sysdict, "warnoptions", warnoptions);
		//     }

		//     v = get_xoptions();
		//     if (v != nil) {
		//         PyDict_SetItemString(sysdict, "_xoptions", v);
		//     }

		//     /* version_info */
		//     if (VersionInfoType.tp_name == 0) {
		//         PyStructSequence_InitType(&VersionInfoType, &version_info_desc);
		//     }
		//     version_info = make_version_info();
		//     SET_SYS_FROM_STRING("version_info", version_info);
		//     /* prevent user from creating new instances */
		//     VersionInfoType.tp_init = nil;
		//     VersionInfoType.tp_new = nil;

		//     /* implementation */
		//     SET_SYS_FROM_STRING("implementation", make_impl_info(version_info));

		//     /* flags */
		//     if (FlagsType.tp_name == 0) {
		//         PyStructSequence_InitType(&FlagsType, &flags_desc);
		//     }
		//     SET_SYS_FROM_STRING("flags", make_flags());
		//     /* prevent user from creating new instances */
		//     FlagsType.tp_init = nil;
		//     FlagsType.tp_new = nil;

		//     /* float repr style: 0.03 (short) vs 0.029999999999999999 (legacy) */
		// #ifndef PY_NO_SHORT_FLOAT_REPR
		//     SET_SYS_FROM_STRING("float_repr_style",
		//                         PyUnicode_FromString("short"));
		// #else
		//     SET_SYS_FROM_STRING("float_repr_style",
		//                         PyUnicode_FromString("legacy"));
		// #endif

		// #ifdef WITH_THREAD
		//     SET_SYS_FROM_STRING("thread_info", PyThread_GetInfo());
		// #endif
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "sys",
			Doc:  module_doc,
		},
		Methods: methods,
		Globals: globals,
	})

}
