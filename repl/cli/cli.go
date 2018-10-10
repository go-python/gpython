// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Read Eval Print Loop for CLI
package cli

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/repl"
	"github.com/peterh/liner"
)

const HistoryFileName = ".gpyhistory"

// homeDirectory finds the home directory or returns ""
func homeDirectory() string {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}
	// Fall back to reading $HOME - work around user.Current() not
	// working for cross compiled binaries on OSX.
	// https://github.com/golang/go/issues/6376
	return os.Getenv("HOME")
}

// Holds state for readline services
type readline struct {
	*liner.State
	repl        *repl.REPL
	historyFile string
	module      *py.Module
	prompt      string
}

// newReadline creates a new instance of readline
func newReadline(repl *repl.REPL) *readline {
	rl := &readline{
		State: liner.NewLiner(),
		repl:  repl,
	}
	home := homeDirectory()
	if home != "" {
		rl.historyFile = filepath.Join(home, HistoryFileName)
	}
	rl.SetTabCompletionStyle(liner.TabPrints)
	rl.SetWordCompleter(rl.Completer)
	return rl
}

// readHistory reads the history into the term
func (rl *readline) ReadHistory() error {
	f, err := os.Open(rl.historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = rl.State.ReadHistory(f)
	if err != nil {
		return err
	}
	return nil
}

// writeHistory writes the history from the term
func (rl *readline) WriteHistory() error {
	f, err := os.OpenFile(rl.historyFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = rl.State.WriteHistory(f)
	if err != nil {
		return err
	}
	return nil
}

// Close the readline and write history
func (rl *readline) Close() error {
	err := rl.State.Close()
	if err != nil {
		return err
	}
	if rl.historyFile != "" {
		err := rl.WriteHistory()
		if err != nil {
			return err
		}
	}
	return nil
}

// Completer takes the currently edited line with the cursor
// position and returns the completion candidates for the partial word
// to be completed. If the line is "Hello, wo!!!" and the cursor is
// before the first '!', ("Hello, wo!!!", 9) is passed to the
// completer which may returns ("Hello, ", {"world", "Word"}, "!!!")
// to have "Hello, world!!!".
func (rl *readline) Completer(line string, pos int) (head string, completions []string, tail string) {
	return rl.repl.Completer(line, pos)
}

// SetPrompt sets the current terminal prompt
func (rl *readline) SetPrompt(prompt string) {
	rl.prompt = prompt
}

// Print prints the output
func (rl *readline) Print(out string) {
	_, _ = os.Stdout.WriteString(out + "\n")
}

// RunREPL starts the REPL loop
func RunREPL() {
	repl := repl.New()
	rl := newReadline(repl)
	repl.SetUI(rl)
	defer rl.Close()
	err := rl.ReadHistory()
	if err != nil {
		fmt.Printf("Failed to open history: %v\n", err)
	}

	fmt.Printf("Gpython 3.4.0\n")

	for {
		line, err := rl.Prompt(rl.prompt)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\n")
				break
			}
			fmt.Printf("Problem reading line: %v\n", err)
			continue
		}
		if line != "" {
			rl.AppendHistory(line)
		}
		rl.repl.Run(line)
	}
}
