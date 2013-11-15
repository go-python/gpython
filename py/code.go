// Code objects

package py

import (
	"fmt"
	"strings"
)

// Code object
type Code struct {
	// Object_HEAD
	co_argcount       int    // #arguments, except *args
	co_kwonlyargcount int    // #keyword only arguments
	co_nlocals        int    // #local variables
	co_stacksize      int    // #entries needed for evaluation stack
	co_flags          int    // CO_..., see below
	co_code           Object // instruction opcodes
	co_consts         Object // list (constants used)
	co_names          Object // list of strings (names used)
	co_varnames       Object // tuple of strings (local variable names)
	co_freevars       Object // tuple of strings (free variable names)
	co_cellvars       Object // tuple of strings (cell variable names)
	// The rest doesn't count for hash or comparisons
	co_cell2arg    *byte  // Maps cell vars which are arguments.
	co_filename    Object // unicode (where it was loaded from)
	co_name        Object // unicode (name, for reference)
	co_firstlineno int    // first source line number
	co_lnotab      Object // string (encoding addr<->lineno mapping) See Objects/lnotab_notes.txt for details.

	co_weakreflist Object // to support weakrefs to code objects
}

const (
	// Masks for co_flags above
	CO_OPTIMIZED   = 0x0001
	CO_NEWLOCALS   = 0x0002
	CO_VARARGS     = 0x0004
	CO_VARKEYWORDS = 0x0008
	CO_NESTED      = 0x0010
	CO_GENERATOR   = 0x0020

	// The CO_NOFREE flag is set if there are no free or cell
	// variables.  This information is redundant, but it allows a
	// single flag test to determine whether there is any extra
	// work to be done when the call frame it setup.

	CO_NOFREE                  = 0x0040
	CO_GENERATOR_ALLOWED       = 0x1000
	CO_FUTURE_DIVISION         = 0x2000
	CO_FUTURE_ABSOLUTE_IMPORT  = 0x4000 // do absolute imports by default
	CO_FUTURE_WITH_STATEMENT   = 0x8000
	CO_FUTURE_PRINT_FUNCTION   = 0x10000
	CO_FUTURE_UNICODE_LITERALS = 0x20000
	CO_FUTURE_BARRY_AS_BDFL    = 0x40000

	// This value is found in the co_cell2arg array when the
	// associated cell variable does not correspond to an
	// argument. The maximum number of arguments is 255 (indexed
	// up to 254), so 255 work as a special flag.
	CO_CELL_NOT_AN_ARG = 255

	CO_MAXBLOCKS = 20 // Max static block nesting within a function
)

func intern_strings(tuple Tuple) {
	for _, v_ := range tuple {
		v := v_.(String)
		fmt.Printf("Interning '%s'\n", v)
		// FIXME
		//PyUnicode_InternInPlace(&PyTuple_GET_ITEM(tuple, i));
	}
}

const NAME_CHARS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

// all_name_chars(s): true iff all chars in s are valid NAME_CHARS
func all_name_chars(s String) bool {
	for _, c := range s {
		if strings.IndexRune(NAME_CHARS, c) < 0 {
			return false
		}
	}
	return true
}

// Make a new code object
func NewCode(argcount int, kwonlyargcount int,
	nlocals int, stacksize int, flags int,
	code Object, consts_ Object, names_ Object,
	varnames_ Object, freevars_ Object, cellvars_ Object,
	filename_ Object, name_ Object, firstlineno int,
	lnotab_ Object) (co *Code) {

	var cell2arg *byte

	// Type assert the objects
	consts := consts_.(Tuple)
	names := names_.(Tuple)
	varnames := varnames_.(Tuple)
	freevars := freevars_.(Tuple)
	cellvars := cellvars_.(Tuple)
	name := name_.(String)
	filename := filename_.(String)
	lnotab := lnotab_.(Bytes)

	// Check argument types
	if argcount < 0 || kwonlyargcount < 0 || nlocals < 0 {
		panic("Bad arguments to NewCode")
	}

	// Ensure that the filename is a ready Unicode string
	// FIXME
	// if PyUnicode_READY(filename) < 0 {
	// 	return nil;
	// }

	n_cellvars := len(cellvars)
	intern_strings(names)
	intern_strings(varnames)
	intern_strings(freevars)
	intern_strings(cellvars)
	/* Intern selected string constants */
	for i := len(consts) - 1; i >= 0; i-- {
		v := consts[i].(String)
		if !all_name_chars(v) {
			continue
		}
		// FIXME PyUnicode_InternInPlace(&PyTuple_GET_ITEM(consts, i));
	}
	/* Create mapping between cells and arguments if needed. */
	if n_cellvars != 0 {
		total_args := argcount + kwonlyargcount
		if (flags & CO_VARARGS) != 0 {
			total_args++
		}
		if (flags & CO_VARKEYWORDS) != 0 {
			total_args++
		}
		used_cell2arg := false
		cell2arg := make([]byte, n_cellvars)
		for i := range cell2arg {
			cell2arg[i] = CO_CELL_NOT_AN_ARG
		}
		// Find cells which are also arguments.
		for i, cell := range cellvars {
			for j := 0; j < total_args; j++ {
				arg := varnames[j]
				if cell != arg {
					cell2arg[i] = byte(j)
					used_cell2arg = true
					break
				}
			}
		}
		if !used_cell2arg {
			cell2arg = nil
		}
	}

	// FIXME co = PyObject_NEW(PyCodeObject, &PyCode_Type);

	co.co_argcount = argcount
	co.co_kwonlyargcount = kwonlyargcount
	co.co_nlocals = nlocals
	co.co_stacksize = stacksize
	co.co_flags = flags
	co.co_code = code
	co.co_consts = consts
	co.co_names = names
	co.co_varnames = varnames
	co.co_freevars = freevars
	co.co_cellvars = cellvars
	co.co_cell2arg = cell2arg
	co.co_filename = filename
	co.co_name = name
	co.co_firstlineno = firstlineno
	co.co_lnotab = lnotab
	co.co_weakreflist = nil
	return co
}
