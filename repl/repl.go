// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Read Eval Print Loop
package repl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-python/gpython/compile"
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/vm"
)

// Possible prompts for the REPL
const (
	NormalPrompt       = ">>> "
	ContinuationPrompt = "... "
)

// Repl state
type REPL struct {
	module       *py.Module
	prog         string
	continuation bool
	previous     string
	term         UI
}

// UI implements the user interface for the REPL
type UI interface {
	// Set the prompt for the start of line
	SetPrompt(string)

	// Print a line of output
	Print(string)
}

// New create a new REPL and initialises the state machine
func New() *REPL {
	r := &REPL{
		module:       py.NewModule("__main__", "", nil, nil),
		prog:         "<stdin>",
		continuation: false,
		previous:     "",
	}
	r.module.Globals["__file__"] = py.String(r.prog)
	return r
}

// SetUI initialises the output user interface
func (r *REPL) SetUI(term UI) {
	r.term = term
	r.term.SetPrompt(NormalPrompt)
}

// Run runs a single line of the REPL
func (r *REPL) Run(line string) {
	// Override the PrintExpr output temporarily
	oldPrintExpr := vm.PrintExpr
	vm.PrintExpr = r.term.Print
	defer func() {
		vm.PrintExpr = oldPrintExpr
	}()
	if r.continuation {
		if line != "" {
			r.previous += string(line) + "\n"
			return
		}
	}
	// need +"\n" because "single" expects \n terminated input
	toCompile := r.previous + string(line)
	if toCompile == "" {
		return
	}
	obj, err := compile.Compile(toCompile+"\n", r.prog, "single", 0, true)
	if err != nil {
		// Detect that we should start a continuation line
		// FIXME detect EOF properly!
		errText := err.Error()
		if strings.Contains(errText, "unexpected EOF while parsing") || strings.Contains(errText, "EOF while scanning triple-quoted string literal") {
			stripped := strings.TrimSpace(toCompile)
			isComment := len(stripped) > 0 && stripped[0] == '#'
			if !isComment {
				r.continuation = true
				r.previous += string(line) + "\n"
				r.term.SetPrompt(ContinuationPrompt)
			}
			return
		}
	}
	r.continuation = false
	r.term.SetPrompt(NormalPrompt)
	r.previous = ""
	if err != nil {
		r.term.Print(fmt.Sprintf("Compile error: %v", err))
		return
	}
	code := obj.(*py.Code)
	_, err = vm.Run(r.module.Globals, r.module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
	}
}

// WordCompleter takes the currently edited line with the cursor
// position and returns the completion candidates for the partial word
// to be completed. If the line is "Hello, wo!!!" and the cursor is
// before the first '!', ("Hello, wo!!!", 9) is passed to the
// completer which may returns ("Hello, ", {"world", "Word"}, "!!!")
// to have "Hello, world!!!".
func (r *REPL) Completer(line string, pos int) (head string, completions []string, tail string) {
	head = line[:pos]
	tail = line[pos:]
	lastSpace := strings.LastIndex(head, " ")
	head, partial := line[:lastSpace+1], line[lastSpace+1:]
	// log.Printf("head = %q, partial = %q, tail = %q", head, partial, tail)
	found := make(map[string]struct{})
	match := func(d py.StringDict) {
		for k := range d {
			if strings.HasPrefix(k, partial) {
				if _, ok := found[k]; !ok {
					completions = append(completions, k)
					found[k] = struct{}{}
				}
			}
		}
	}
	match(r.module.Globals)
	match(py.Builtins.Globals)
	sort.Strings(completions)
	return head, completions, tail
}
