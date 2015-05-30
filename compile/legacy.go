package compile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ncw/gpython/marshal"
	"github.com/ncw/gpython/py"
)

// Compile with python3.4 - not used any more but keep for the moment!

// Compile(source, filename, mode, flags, dont_inherit) -> code object
//
// Compile the source string (a Python module, statement or expression)
// into a code object that can be executed by exec() or eval().
// The filename will be used for run-time error messages.
// The mode must be 'exec' to compile a module, 'single' to compile a
// single (interactive) statement, or 'eval' to compile an expression.
// The flags argument, if present, controls which future statements influence
// the compilation of the code.
// The dont_inherit argument, if non-zero, stops the compilation inheriting
// the effects of any future statements in effect in the code calling
// compile; if absent or zero these statements do influence the compilation,
// in addition to any features explicitly specified.
func LegacyCompile(str, filename, mode string, flags int, dont_inherit bool) py.Object {
	dont_inherit_str := "False"
	if dont_inherit {
		dont_inherit_str = "True"
	}
	// FIXME escaping in filename
	code := fmt.Sprintf(`import sys, marshal
str = sys.stdin.buffer.read().decode("utf-8")
code = compile(str, "%s", "%s", %d, %s)
marshalled_code = marshal.dumps(code)
sys.stdout.buffer.write(marshalled_code)
sys.stdout.close()`,
		filename,
		mode,
		flags,
		dont_inherit_str,
	)
	cmd := exec.Command("python3.4", "-c", code)
	cmd.Stdin = strings.NewReader(str)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "--- Failed to run python3.4 compile ---\n")
		fmt.Fprintf(os.Stderr, "--------------------\n")
		_, _ = os.Stderr.Write(stderr.Bytes())
		fmt.Fprintf(os.Stderr, "--------------------\n")
		panic(err)
	}
	obj, err := marshal.ReadObject(bytes.NewBuffer(out.Bytes()))
	if err != nil {
		panic(err)
	}
	return obj
}
