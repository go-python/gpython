// Read Eval Print Loop
package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ncw/gpython/compile"
	"github.com/ncw/gpython/py"
	"github.com/ncw/gpython/vm"
)

func Run() {
	fmt.Printf("Gpython 3.4.0\n")
	bio := bufio.NewReader(os.Stdin)
	module := py.NewModule("__main__", "", nil, nil)
	prog := "<stdin>"
	module.Globals["__file__"] = py.String(prog)
	for {
		fmt.Printf(">>> ")
		line, hasMoreInLine, err := bio.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error: %v", err)
			break
		}
		if hasMoreInLine {
			log.Printf("Line truncated")
		}
		// FIXME need +"\n" because "single" is broken
		obj, err := compile.Compile(string(line)+"\n", prog, "single", 0, true)
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
