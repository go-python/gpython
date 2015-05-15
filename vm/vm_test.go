package vm_test

import (
	"fmt"
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
		t.Fatalf("Open failed: %v", err)
	}
	defer f.Close()

	str, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	obj, err := compile.Compile(string(str), prog, "exec", 0, true)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	code := obj.(*py.Code)
	module := py.NewModule("__main__", "", nil, nil)
	module.Globals["__file__"] = py.String(prog)

	res, err := vm.Run(module.Globals, module.Globals, code, nil)
	if err != nil {
		py.TracebackDump(err)
		t.Fatalf("Run failed: %v", err)
	}

	fmt.Printf("Return = %v\n", res)
}

func TestVm(t *testing.T) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	for _, f := range files {
		name := path.Join(testDir, f.Name())
		if strings.HasSuffix(name, ".py") {
			t.Logf("Running %q", name)
			run(t, name)
		}
	}
}
