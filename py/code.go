// Code objects

package py

import (
	"strings"
)

// Code object
type Code struct {
	// Object_HEAD
	argcount       int32  // #arguments, except *args
	kwonlyargcount int32  // #keyword only arguments
	nlocals        int32  // #local variables
	stacksize      int32  // #entries needed for evaluation stack
	flags          int32  // CO_..., see below
	code           Object // instruction opcodes
	consts         Object // list (constants used)
	names          Object // list of strings (names used)
	varnames       Object // tuple of strings (local variable names)
	freevars       Object // tuple of strings (free variable names)
	cellvars       Object // tuple of strings (cell variable names)
	// The rest doesn't count for hash or comparisons
	cell2arg    *byte  // Maps cell vars which are arguments.
	filename    Object // unicode (where it was loaded from)
	name        Object // unicode (name, for reference)
	firstlineno int32  // first source line number
	lnotab      Object // string (encoding addr<->lineno mapping) See Objects/lnotab_notes.txt for details.

	weakreflist Object // to support weakrefs to code objects
}

var CodeType = NewType("code")

// Type of this object
func (o *Code) Type() *Type {
	return CodeType
}

// Make sure it satisfies the interface
var _ Object = (*Code)(nil)

const (
	// Masks for flags above
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

	// This value is found in the cell2arg array when the
	// associated cell variable does not correspond to an
	// argument. The maximum number of arguments is 255 (indexed
	// up to 254), so 255 work as a special flag.
	CO_CELL_NOT_AN_ARG = 255

	CO_MAXBLOCKS = 20 // Max static block nesting within a function
)

// Intern all the strings in the tuple
func intern_strings(tuple Tuple) {
	for i, v_ := range tuple {
		v := v_.(String)
		tuple[i] = v.Intern()
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
func NewCode(argcount int32, kwonlyargcount int32,
	nlocals int32, stacksize int32, flags int32,
	code Object, consts_ Object, names_ Object,
	varnames_ Object, freevars_ Object, cellvars_ Object,
	filename_ Object, name_ Object, firstlineno int32,
	lnotab_ Object) *Code {

	var cell2arg *byte

	// Type assert the objects
	consts := consts_.(Tuple)
	names := names_.(Tuple)
	varnames := varnames_.(Tuple)
	freevars := freevars_.(Tuple)
	cellvars := cellvars_.(Tuple)
	name := name_.(String)
	filename := filename_.(String)
	lnotab := lnotab_.(String)

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
		if v, ok := consts[i].(String); ok {
			if all_name_chars(v) {
				consts[i] = v.Intern()
			}
		}
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
			for j := int32(0); j < total_args; j++ {
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

	return &Code{
		argcount:       argcount,
		kwonlyargcount: kwonlyargcount,
		nlocals:        nlocals,
		stacksize:      stacksize,
		flags:          flags,
		code:           code,
		consts:         consts,
		names:          names,
		varnames:       varnames,
		freevars:       freevars,
		cellvars:       cellvars,
		cell2arg:       cell2arg,
		filename:       filename,
		name:           name,
		firstlineno:    firstlineno,
		lnotab:         lnotab,
		weakreflist:    nil,
	}
}
