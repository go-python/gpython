// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code objects

package py

import (
	"strings"
)

// Code object
//
// Freevars are variables declared in the namespace where the
// code object was defined, they are used when closures are
// used - these become Cellvars in the function with a closure.
//
// Cellvars are names of local variables referenced by functions with
// a closure.
type Code struct {
	Argcount       int32    // #arguments, except *args
	Kwonlyargcount int32    // #keyword only arguments
	Nlocals        int32    // #local variables
	Stacksize      int32    // #entries needed for evaluation stack
	Flags          int32    // CO_..., see below
	Code           string   // instruction opcodes
	Consts         Tuple    // list (constants used)
	Names          []string // list of strings (names used)
	Varnames       []string // tuple of strings (local variable names)
	Freevars       []string // tuple of strings (free variable names)
	Cellvars       []string // tuple of strings (cell variable names)
	// The rest doesn't count for hash or comparisons
	Cell2arg    []byte // Maps cell vars which are arguments.
	Filename    string // unicode (where it was loaded from)
	Name        string // unicode (name, for reference)
	Firstlineno int32  // first source line number
	Lnotab      string // string (encoding addr<->lineno mapping) See Objects/lnotab_notes.txt for details.

	Weakreflist *List // to support weakrefs to code objects
}

var CodeType = NewType("code", "code(argcount, kwonlyargcount, nlocals, stacksize, flags, codestring,\n      constants, names, varnames, filename, name, firstlineno,\n      lnotab[, freevars[, cellvars]])\n\nCreate a code object.  Not for the faint of heart.")

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

	CO_COMPILER_FLAGS_MASK = CO_FUTURE_DIVISION | CO_FUTURE_ABSOLUTE_IMPORT | CO_FUTURE_WITH_STATEMENT | CO_FUTURE_PRINT_FUNCTION | CO_FUTURE_UNICODE_LITERALS | CO_FUTURE_BARRY_AS_BDFL

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
		if !strings.ContainsRune(NAME_CHARS, c) {
			return false
		}
	}
	return true
}

// Make a new code object
func NewCode(argcount int32, kwonlyargcount int32,
	nlocals int32, stacksize int32, flags int32,
	code_ Object, consts_ Object, names_ Object,
	varnames_ Object, freevars_ Object, cellvars_ Object,
	filename_ Object, name_ Object, firstlineno int32,
	lnotab_ Object) *Code {

	var cell2arg []byte

	// Type assert the objects
	consts := consts_.(Tuple)
	namesTuple := names_.(Tuple)
	varnamesTuple := varnames_.(Tuple)
	freevarsTuple := freevars_.(Tuple)
	cellvarsTuple := cellvars_.(Tuple)
	name := string(name_.(String))
	filename := string(filename_.(String))
	lnotab := string(lnotab_.(String))
	code := string(code_.(String))

	// Convert Tuples to native []string for speed
	names := make([]string, len(namesTuple))
	for i := range namesTuple {
		names[i] = string(namesTuple[i].(String))
	}
	varnames := make([]string, len(varnamesTuple))
	for i := range varnamesTuple {
		varnames[i] = string(varnamesTuple[i].(String))
	}
	freevars := make([]string, len(freevarsTuple))
	for i := range freevarsTuple {
		freevars[i] = string(freevarsTuple[i].(String))
	}
	cellvars := make([]string, len(cellvarsTuple))
	for i := range cellvarsTuple {
		cellvars[i] = string(cellvarsTuple[i].(String))
	}

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
	intern_strings(namesTuple)
	intern_strings(varnamesTuple)
	intern_strings(freevarsTuple)
	intern_strings(cellvarsTuple)
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
		if flags&CO_VARARGS != 0 {
			total_args++
		}
		if flags&CO_VARKEYWORDS != 0 {
			total_args++
		}
		used_cell2arg := false
		cell2arg = make([]byte, n_cellvars)
		for i := range cell2arg {
			cell2arg[i] = CO_CELL_NOT_AN_ARG
		}
		// Find cells which are also arguments.
		for i, cell := range cellvars {
			for j := int32(0); j < total_args; j++ {
				arg := varnames[j]
				if cell == arg {
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
		Argcount:       argcount,
		Kwonlyargcount: kwonlyargcount,
		Nlocals:        nlocals,
		Stacksize:      stacksize,
		Flags:          flags,
		Code:           code,
		Consts:         consts,
		Names:          names,
		Varnames:       varnames,
		Freevars:       freevars,
		Cellvars:       cellvars,
		Cell2arg:       cell2arg,
		Filename:       filename,
		Name:           name,
		Firstlineno:    firstlineno,
		Lnotab:         lnotab,
		Weakreflist:    nil,
	}
}

// Return number of free variables
func (co *Code) GetNumFree() int {
	return len(co.Freevars)
}

// Use co_lnotab to compute the line number from a bytecode index,
// addrq.  See lnotab_notes.txt for the details of the lnotab
// representation.
func (co *Code) Addr2Line(addrq int32) int32 {
	line := co.Firstlineno
	addr := int32(0)
	for i := 0; i < len(co.Lnotab); i += 2 {
		addr += int32(co.Lnotab[i])
		if addr > addrq {
			break
		}
		line += int32(co.Lnotab[i+1])
	}
	return line
}

// FIXME this should be the default?
func (co *Code) M__eq__(other Object) (Object, error) {
	if otherCo, ok := other.(*Code); ok && co == otherCo {
		return True, nil
	}
	return False, nil
}

// FIXME this should be the default?
func (co *Code) M__ne__(other Object) (Object, error) {
	if otherCo, ok := other.(*Code); ok && co == otherCo {
		return False, nil
	}
	return True, nil
}

// Check interface is satisfied
var _ I__eq__ = (*Code)(nil)
var _ I__ne__ = (*Code)(nil)
