// Python virtual machine
package vm

import (
	"github.com/ncw/gpython/py"
)

//go:generate stringer -type=vmExit -output stringer.go

// VM exit type
type vmExit byte

// VM exit values
const (
	exitNot       vmExit = iota // No error
	exitException               // Exception occurred
	exitReturn                  // 'return' statement
	exitBreak                   // 'break' statement
	exitContinue                // 'continue' statement
	exitYield                   // 'yield' operator
	exitSilenced                // Exception silenced by 'with'
)

// Virtual machine state
type Vm struct {
	// Current frame
	frame *py.Frame
	// Whether ext should be added to the next arg
	extended bool
	// 16 bit extension for argument for next opcode
	ext int32
	// Return value
	result py.Object
	// Exit value
	exit vmExit
	// Current Pending exception type, value and traceback
	curexc py.ExceptionInfo
	// Previous exception type, value and traceback
	exc py.ExceptionInfo
}

// Make a new VM
func NewVm(frame *py.Frame) *Vm {
	return &Vm{
		frame: frame,
	}
}
