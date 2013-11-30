// Frame objects

package py

// Store information about try blocks
type TryBlock struct {
	Type    byte  // what kind of block this is (the opcode)
	Handler int32 // where to jump to find handler
	Level   int   // value stack level to pop to
}

// A python Frame object
type Frame struct {
	// Back       *Frame        // previous frame, or nil
	Code     *Code      // code segment
	Builtins StringDict // builtin symbol table
	Globals  StringDict // global symbol table
	Locals   StringDict // local symbol table
	Stack    []Object   // Valuestack
	// Next free slot in f_valuestack.  Frame creation sets to f_valuestack.
	// Frame evaluation usually NULLs it, but a frame that yields sets it
	// to the current stack top.
	// Stacktop *Object
	Yielded bool   // set if the function yielded, cleared otherwise
	Trace   Object // Trace function

	// In a generator, we need to be able to swap between the exception
	// state inside the generator and the exception state of the calling
	// frame (which shouldn't be impacted when the generator "yields"
	// from an except handler).
	// These three fields exist exactly for that, and are unused for
	// non-generator frames. See the save_exc_state and swap_exc_state
	// functions in ceval.c for details of their use.
	Exc_type      Object
	Exc_value     *Object
	Exc_traceback *Object
	// Borrowed reference to a generator, or NULL
	Gen Object

	// FIXME Tstate *PyThreadState
	Lasti int32 // Last instruction if called
	// Call PyFrame_GetLineNumber() instead of reading this field
	// directly.  As of 2.3 f_lineno is only valid when tracing is
	// active (i.e. when f_trace is set).  At other times we use
	// PyCode_Addr2Line to calculate the line from the current
	// bytecode index.
	Lineno     int        // Current line number
	Iblock     int        // index in f_blockstack
	Executing  byte       // whether the frame is still executing
	Blockstack []TryBlock // for try and loop blocks
	Block      *TryBlock  // pointer to current block or nil
	Localsplus []Object   // locals+stack, dynamically sized
}

var FrameType = NewType("frame", "Represents a stack frame")

// Type of this object
func (o *Frame) Type() *Type {
	return FrameType
}

// Make a new frame for a code object
func NewFrame(globals, locals StringDict, code *Code) *Frame {
	return &Frame{
		Globals:  globals,
		Locals:   locals,
		Code:     code,
		Builtins: Builtins.Globals,
		Stack:    make([]Object, 0, code.Stacksize),
	}
}

// Python names are looked up in three scopes
//
// First the local scope
// Next the module global scope
// And finally the builtins
func (f *Frame) Lookup(name string) (obj Object) {
	var ok bool

	// Lookup in locals
	// fmt.Printf("locals = %v\n", f.Locals)
	if obj, ok = f.Locals[name]; ok {
		return
	}

	// Lookup in globals
	// fmt.Printf("globals = %v\n", f.Globals)
	if obj, ok = f.Globals[name]; ok {
		return
	}

	// Lookup in builtins
	// fmt.Printf("builtins = %v\n", Builtins.Globals)
	if obj, ok = f.Builtins[name]; ok {
		return
	}

	panic(ExceptionNewf(NameError, "name '%s' is not defined", name))
}

// Make a new Block (try/for/while)
func (f *Frame) PushBlock(Type byte, Handler int32, Level int) {
	f.Blockstack = append(f.Blockstack, TryBlock{
		Type:    Type,
		Handler: Handler,
		Level:   Level,
	})
	f.Block = &f.Blockstack[len(f.Blockstack)-1]
}

// Pop the current block off
func (f *Frame) PopBlock() {
	f.Blockstack = f.Blockstack[:len(f.Blockstack)-1]
	if len(f.Blockstack) > 0 {
		f.Block = &f.Blockstack[len(f.Blockstack)-1]
	} else {
		f.Block = nil
	}
}
