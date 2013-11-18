// name lookup

package vm

import (
	"fmt"
	"github.com/ncw/gpython/py"
)

// Python names are looked up in three scopes
//
// First the local scope
// Next the module global scope
// And finally the builtins
func (vm *Vm) lookup(name string) (obj py.Object) {
	var ok bool

	// Lookup in locals
	fmt.Printf("locals = %v\n", vm.locals)
	if obj, ok = vm.locals[name]; ok {
		return
	}

	// Lookup in globals
	fmt.Printf("globals = %v\n", vm.globals)
	if obj, ok = vm.globals[name]; ok {
		return
	}

	// Lookup in builtins
	fmt.Printf("builtins = %v\n", py.Builtins.Globals)
	if obj, ok = py.Builtins.Globals[name]; ok {
		return
	}

	// FIXME this should be a NameError
	panic(fmt.Sprintf("NameError: name '%s' is not defined", name))
}

// Lookup a method
func (vm *Vm) lookupMethod(name string) py.Callable {
	obj := vm.lookup(name)
	method, ok := obj.(py.Callable)
	if !ok {
		// FIXME should be TypeError
		panic(fmt.Sprintf("TypeError: '%s' object is not callable", obj.Type().Name))
	}
	return method
}

// Calls function fn with args and kwargs
//
// kwargs is a sequence of name, value pairs
func (vm *Vm) call(fn py.Object, args []py.Object, kwargs []py.Object) py.Object {
	fmt.Printf("Call %v with args = %v, kwargs = %v\n", fn, args, kwargs)
	fn_name := string(fn.(py.String))
	method := vm.lookupMethod(fn_name)
	self := py.None // FIXME should be the module
	if len(kwargs) > 0 {
		// Convert kwargs into dictionary
		if len(kwargs)%2 != 0 {
			panic("Odd length kwargs")
		}
		kwargsd := py.NewStringDict()
		for i := 0; i < len(kwargs); i += 2 {
			kwargsd[string(kwargs[i].(py.String))] = kwargs[i+1]
		}
		return method.CallWithKeywords(self, args, kwargsd)
	} else {
		return method.Call(self, args)
	}
}
