package repl

import (
	"fmt"
	"reflect"
	"testing"

	// import required modules
	_ "github.com/go-python/gpython/stdlib"
)

type replTest struct {
	prompt string
	out    string
}

// SetPrompt sets the current terminal prompt
func (rt *replTest) SetPrompt(prompt string) {
	rt.prompt = prompt
}

// Print prints the output
func (rt *replTest) Print(out string) {
	rt.out = out
}

func (rt *replTest) assert(t *testing.T, what, wantPrompt, wantOut string) {
	t.Helper()
	if rt.prompt != wantPrompt {
		t.Errorf("%s: Prompt wrong:\ngot= %q\nwant=%q", what, rt.prompt, wantPrompt)
	}
	if rt.out != wantOut {
		t.Errorf("%s: Output wrong:\ngot= %q\nwant=%q", what, rt.out, wantOut)
	}
	rt.out = ""
}

func TestREPL(t *testing.T) {
	r := New(nil)
	rt := &replTest{}
	r.SetUI(rt)

	rt.assert(t, "init", NormalPrompt, "")

	r.Run("")
	rt.assert(t, "empty", NormalPrompt, "")

	r.Run("1+2")
	rt.assert(t, "1+2", NormalPrompt, "3")

	// FIXME this output goes to Stderr and Stdout
	r.Run("aksfjakf")
	rt.assert(t, "unbound", NormalPrompt, "")

	r.Run("sum = 0")
	rt.assert(t, "multi#1", NormalPrompt, "")
	r.Run("for i in range(10):")
	rt.assert(t, "multi#2", ContinuationPrompt, "")
	r.Run("    sum += i")
	rt.assert(t, "multi#3", ContinuationPrompt, "")
	r.Run("")
	rt.assert(t, "multi#4", NormalPrompt, "")
	r.Run("sum")
	rt.assert(t, "multi#5", NormalPrompt, "45")

	r.Run("if")
	rt.assert(t, "compileError", NormalPrompt, "Compile error: \n  File \"<stdin>\", line 1, offset 2\n    if\n\n\nSyntaxError: 'invalid syntax'")

	// test comments in the REPL work properly
	r.Run("# this is a comment")
	rt.assert(t, "comment", NormalPrompt, "")
	r.Run("a = 42")
	rt.assert(t, "comment continuation", NormalPrompt, "")
	r.Run("a")
	rt.assert(t, "comment check", NormalPrompt, "42")
}

func TestCompleter(t *testing.T) {
	r := New(nil)
	rt := &replTest{}
	r.SetUI(rt)

	for _, test := range []struct {
		line            string
		pos             int
		wantHead        string
		wantCompletions []string
		wantTail        string
	}{
		{
			line:            "di",
			pos:             2,
			wantHead:        "",
			wantCompletions: []string{"dict", "divmod"},
			wantTail:        "",
		},
		{
			line:            "div",
			pos:             3,
			wantHead:        "",
			wantCompletions: []string{"divmod"},
			wantTail:        "",
		},
		{
			line:            "doodle",
			pos:             6,
			wantHead:        "",
			wantCompletions: nil,
			wantTail:        "",
		},
		{
			line:            "divmod divm",
			pos:             9,
			wantHead:        "divmod ",
			wantCompletions: []string{"divmod"},
			wantTail:        "vm",
		},
	} {
		t.Run(fmt.Sprintf("line=%q,pos=%d)", test.line, test.pos), func(t *testing.T) {
			gotHead, gotCompletions, gotTail := r.Completer(test.line, test.pos)
			if test.wantHead != gotHead {
				t.Errorf("invalid head:\ngot= %q\nwant=%q", gotHead, test.wantHead)
			}
			if !reflect.DeepEqual(test.wantCompletions, gotCompletions) {
				t.Errorf("invalid completions:\ngot= %#v\nwant=%#v", gotCompletions, test.wantCompletions)
			}
			if test.wantTail != gotTail {
				t.Errorf("invalid tail:\ngot= %q\nwant=%q", gotTail, test.wantTail)
			}
		})
	}

}
