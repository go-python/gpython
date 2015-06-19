// Read Eval Print Loop
package repl

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ncw/gpython/compile"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/vm"
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
	historyFile string
	module      *py.Module
}

// newReadline creates a new instance of readline
func newReadline(module *py.Module) *readline {
	rl := &readline{
		State:  liner.NewLiner(),
		module: module,
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

// WordCompleter takes the currently edited line with the cursor
// position and returns the completion candidates for the partial word
// to be completed. If the line is "Hello, wo!!!" and the cursor is
// before the first '!', ("Hello, wo!!!", 9) is passed to the
// completer which may returns ("Hello, ", {"world", "Word"}, "!!!")
// to have "Hello, world!!!".
func (rl *readline) Completer(line string, pos int) (head string, completions []string, tail string) {
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
	match(rl.module.Globals)
	match(py.Builtins.Globals)
	sort.Strings(completions)
	return head, completions, tail
}

func Run() {
	module := py.NewModule("__main__", "", nil, nil)
	rl := newReadline(module)
	defer rl.Close()
	err := rl.ReadHistory()
	if err != nil {
		fmt.Printf("Failed to open history: %v\n", err)
	}

	fmt.Printf("Gpython 3.4.0\n")
	prog := "<stdin>"
	module.Globals["__file__"] = py.String(prog)
	continuation := false
	previous := ""
	for {
		prompt := ">>> "
		if continuation {
			prompt = "... "
		}
		line, err := rl.Prompt(prompt)
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
		if continuation {
			if line != "" {
				previous += string(line) + "\n"
				continue
			}

		}
		// need +"\n" because "single" expects \n terminated input
		toCompile := previous + string(line)
		if toCompile == "" {
			continue
		}
		obj, err := compile.Compile(toCompile+"\n", prog, "single", 0, true)
		if err != nil {
			// Detect that we should start a continuation line
			// FIXME detect EOF properly!
			errText := err.Error()
			if strings.Contains(errText, "unexpected EOF while parsing") || strings.Contains(errText, "EOF while scanning triple-quoted string literal") {
				continuation = true
				previous += string(line) + "\n"
				continue
			}
		}
		continuation = false
		previous = ""
		if err != nil {
			fmt.Printf("Compile error: %v\n", err)
			continue
		}
		code := obj.(*py.Code)
		_, err = vm.Run(module.Globals, module.Globals, code, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	}
}
