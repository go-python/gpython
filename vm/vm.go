// Python virtual machine
package vm

import (
	"github.com/ncw/gpython/py"
)

// Virtual machine state
type Vm struct {
	// Object stack
	stack []py.Object
	// Frame stack
	frames []py.Frame
	// Current frame
	frame *py.Frame
	// Whether ext should be added to the next arg
	extended bool
	// 16 bit extension for argument for next opcode
	ext int32
}

// Make a new VM
func NewVm() *Vm {
	return &Vm{
		stack:  make([]py.Object, 0, 16),
		frames: make([]py.Frame, 0, 16),
	}
}
