// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Python virtual machine
package vm

import (
	"github.com/go-python/gpython/py"
)

//go:generate stringer -type=vmStatus,OpCode -output stringer.go

// VM status code
type vmStatus byte

// VM Status code for main loop (reason for stack unwind)
const (
	whyNot       vmStatus = iota // No error
	whyException                 // Exception occurred
	whyReturn                    // 'return' statement
	whyBreak                     // 'break' statement
	whyContinue                  // 'continue' statement
	whyYield                     // 'yield' operator
	whySilenced                  // Exception silenced by 'with'
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
	retval py.Object
	// VM Status code for main loop
	why vmStatus
	// Current Pending exception type, value and traceback
	curexc py.ExceptionInfo
	// Previous exception type, value and traceback
	exc py.ExceptionInfo
	// VM access to state / modules
	context py.Context
}
