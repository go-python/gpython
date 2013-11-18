// Built-in functions
package builtin

import (
	"fmt"
	"github.com/ncw/gpython/py"
)

const builtin_doc = `Built-in functions, exceptions, and other objects.

Noteworthy: None is the 'nil' object; Ellipsis represents '...' in slices.`

// Initialise the module
func init() {
	methods := []*py.Method{
		// py.NewMethodWithKeywords("__build_class__", builtin___build_class__, py.METH_VARARGS|py.METH_KEYWORDS, build_class_doc),
		// py.NewMethodWithKeywords("__import__", builtin___import__, py.METH_VARARGS|py.METH_KEYWORDS, import_doc),
		// py.NewMethod("abs", builtin_abs, py.METH_O, abs_doc),
		// py.NewMethod("all", builtin_all, py.METH_O, all_doc),
		// py.NewMethod("any", builtin_any, py.METH_O, any_doc),
		// py.NewMethod("ascii", builtin_ascii, py.METH_O, ascii_doc),
		// py.NewMethod("bin", builtin_bin, py.METH_O, bin_doc),
		// py.NewMethod("callable", builtin_callable, py.METH_O, callable_doc),
		// py.NewMethod("chr", builtin_chr, py.METH_VARARGS, chr_doc),
		// py.NewMethodWithKeywords("compile", builtin_compile, py.METH_VARARGS|py.METH_KEYWORDS, compile_doc),
		// py.NewMethod("delattr", builtin_delattr, py.METH_VARARGS, delattr_doc),
		// py.NewMethod("dir", builtin_dir, py.METH_VARARGS, dir_doc),
		// py.NewMethod("divmod", builtin_divmod, py.METH_VARARGS, divmod_doc),
		// py.NewMethod("eval", builtin_eval, py.METH_VARARGS, eval_doc),
		// py.NewMethod("exec", builtin_exec, py.METH_VARARGS, exec_doc),
		// py.NewMethod("format", builtin_format, py.METH_VARARGS, format_doc),
		// py.NewMethod("getattr", builtin_getattr, py.METH_VARARGS, getattr_doc),
		// py.NewMethod("globals", builtin_globals, py.METH_NOARGS, globals_doc),
		// py.NewMethod("hasattr", builtin_hasattr, py.METH_VARARGS, hasattr_doc),
		// py.NewMethod("hash", builtin_hash, py.METH_O, hash_doc),
		// py.NewMethod("hex", builtin_hex, py.METH_O, hex_doc),
		// py.NewMethod("id", builtin_id, py.METH_O, id_doc),
		// py.NewMethod("input", builtin_input, py.METH_VARARGS, input_doc),
		// py.NewMethod("isinstance", builtin_isinstance, py.METH_VARARGS, isinstance_doc),
		// py.NewMethod("issubclass", builtin_issubclass, py.METH_VARARGS, issubclass_doc),
		// py.NewMethod("iter", builtin_iter, py.METH_VARARGS, iter_doc),
		// py.NewMethod("len", builtin_len, py.METH_O, len_doc),
		// py.NewMethod("locals", builtin_locals, py.METH_NOARGS, locals_doc),
		// py.NewMethodWithKeywords("max", builtin_max, py.METH_VARARGS|py.METH_KEYWORDS, max_doc),
		// py.NewMethodWithKeywords("min", builtin_min, py.METH_VARARGS|py.METH_KEYWORDS, min_doc),
		// py.NewMethod("next", builtin_next, py.METH_VARARGS, next_doc),
		// py.NewMethod("oct", builtin_oct, py.METH_O, oct_doc),
		// py.NewMethod("ord", builtin_ord, py.METH_O, ord_doc),
		// py.NewMethod("pow", builtin_pow, py.METH_VARARGS, pow_doc),
		py.NewMethodWithKeywords("print", builtin_print, py.METH_VARARGS|py.METH_KEYWORDS, print_doc),
		// py.NewMethod("repr", builtin_repr, py.METH_O, repr_doc),
		// py.NewMethodWithKeywords("round", builtin_round, py.METH_VARARGS|py.METH_KEYWORDS, round_doc),
		// py.NewMethod("setattr", builtin_setattr, py.METH_VARARGS, setattr_doc),
		// py.NewMethodWithKeywords("sorted", builtin_sorted, py.METH_VARARGS|py.METH_KEYWORDS, sorted_doc),
		// py.NewMethod("sum", builtin_sum, py.METH_VARARGS, sum_doc),
		// py.NewMethod("vars", builtin_vars, py.METH_VARARGS, vars_doc),
	}
	py.NewModule("builtins", builtin_doc, methods)
}

const print_doc = `print(value, ..., sep=' ', end='\\n', file=sys.stdout, flush=False)

Prints the values to a stream, or to sys.stdout by default.
Optional keyword arguments:
file:  a file-like object (stream); defaults to the current sys.stdout.
sep:   string inserted between values, default a space.
end:   string appended after the last value, default a newline.
flush: whether to forcibly flush the stream.`

func builtin_print(self py.Object, args py.Tuple, kwargs py.StringDict) py.Object {
	fmt.Printf("print %v, %v, %v\n", self, args, kwargs)
	return py.None
}
