// Generator objects

package py

// A python Generator object
type Generator struct {
	// Note: gi_frame can be NULL if the generator is "finished"
	Frame *Frame

	// True if generator is being executed.
	Running bool

	// The code object backing the generator
	Code *Code

	// List of weak reference.
	Weakreflist Object
}

var GeneratorType = NewType("generator", "generator object")

// Type of this object
func (o *Generator) Type() *Type {
	return GeneratorType
}

// Define a new generator
func NewGenerator(frame *Frame) *Generator {
	g := &Generator{
		Frame:   frame,
		Running: false,
		Code:    frame.Code,
	}
	return g
}

func (it *Generator) M__iter__() Object {
	return it
}

// generator.__next__()
//
// Starts the execution of a generator function or resumes it at the
// last executed yield expression. When a generator function is
// resumed with a __next__() method, the current yield expression
// always evaluates to None. The execution then continues to the next
// yield expression, where the generator is suspended again, and the
// value of the expression_list is returned to next()‘s caller. If the
// generator exits without yielding another value, a StopIteration
// exception is raised.
//
// This method is normally called implicitly, e.g. by a for loop, or by the built-in next() function.
func (it *Generator) M__next__() Object {
	it.Running = true
	res, err := RunFrame(it.Frame)
	it.Running = false
	// Push a None on the stack for next time
	// FIXME this value is the one sent by Send
	it.Frame.Stack = append(it.Frame.Stack, None)
	// FIXME not correct
	if err != nil {
		panic(err)
	}
	if it.Frame.Yielded {
		return res
	}
	panic(StopIteration)
}

// generator.send(value)
//
// Resumes the execution and “sends” a value into the generator
// function. The value argument becomes the result of the current
// yield expression. The send() method returns the next value yielded
// by the generator, or raises StopIteration if the generator exits
// without yielding another value. When send() is called to start the
// generator, it must be called with None as the argument, because
// there is no yield expression that could receive the value.
func (it *Generator) Send(value Object) {
	panic("generator send not implemented")
}

// generator.throw(type[, value[, traceback]])
//
// Raises an exception of type type at the point where generator was
// paused, and returns the next value yielded by the generator
// function. If the generator exits without yielding another value, a
// StopIteration exception is raised. If the generator function does
// not catch the passed-in exception, or raises a different exception,
// then that exception propagates to the caller.
func (it *Generator) Throw(args Tuple, kwargs StringDict) {
	panic("generator throw not implemented")
}

// generator.close()
//
// Raises a GeneratorExit at the point where the generator function
// was paused. If the generator function then raises StopIteration (by
// exiting normally, or due to already being closed) or GeneratorExit
// (by not catching the exception), close returns to its caller. If
// the generator yields a value, a RuntimeError is raised. If the
// generator raises any other exception, it is propagated to the
// caller. close() does nothing if the generator has already exited
// due to an exception or normal exit.
func (it *Generator) Close() {
	panic("generator close not implemented")
}

// Check interface is satisfied
var _ I_iterator = (*Generator)(nil)
