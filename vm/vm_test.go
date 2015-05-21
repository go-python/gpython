package vm_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	_ "github.com/ncw/gpython/builtin"
	"github.com/ncw/gpython/compile"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/vm"
)

const testDir = "tests"

// Run the code in str
func run(t *testing.T, prog string) {
	f, err := os.Open(prog)
	if err != nil {
		t.Fatalf("%s: Open failed: %v", prog, err)
	}
	defer f.Close()

	str, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("%s: ReadAll failed: %v", prog, err)
	}

	obj, err := compile.Compile(string(str), prog, "exec", 0, true)
	if err != nil {
		t.Fatalf("%s: Compile failed: %v", prog, err)
	}

	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	module.Globals["__file__"] = py.String(prog)

	_, err = vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		t.Fatalf("%s: Run failed: %v", prog, err)
	}

	// t.Logf("%s: Return = %v", prog, res)
	if module.Globals["finished"] != py.True {
		t.Fatalf("%s: Didn't finish", prog)
	}
}

func TestVm(t *testing.T) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	for _, f := range files {
		name := path.Join(testDir, f.Name())
		if strings.HasSuffix(name, ".py") {
			t.Logf("%s: Starting", name)
			run(t, name)
		}
	}
}
