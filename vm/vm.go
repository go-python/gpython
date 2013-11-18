// Python virtual machine
package vm

import (
	"github.com/ncw/gpython/py"
)

// Virtual machine state
type Vm struct {
	// Object stack
	stack []py.Object
	// Current code object we are interpreting
	co *py.Code
	// Current globals
	globals py.StringDict
	// Current locals
	locals py.StringDict
	// Whether ext should be added to the next arg
	extended bool
	// 16 bit extension for argument for next opcode
	ext int32
	// Whether we should exit
	exit bool
}

// Make a new VM
func NewVm() *Vm {
	vm := new(Vm)
	vm.stack = make([]py.Object, 0, 1024)
	vm.locals = py.NewStringDict()
	return vm
}
