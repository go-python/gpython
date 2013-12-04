// Module objects

package py

import (
	"fmt"
)

var (
	// Registry of installed modules
	modules = make(map[string]*Module)
	// Builtin module
	Builtins *Module
	// this should be the frozen module importlib/_bootstrap.py generated
	// by Modules/_freeze_importlib.c into Python/importlib.h
	Importlib *Module
)

// A python Module object
type Module struct {
	Name    string
	Doc     string
	Globals StringDict
	//	dict Dict
}

var ModuleType = NewType("module", "module object")

// Type of this object
func (o *Module) Type() *Type {
	return ModuleType
}

// Define a new module
func NewModule(name, doc string, methods []*Method, globals StringDict) *Module {
	m := &Module{
		Name:    name,
		Doc:     doc,
		Globals: globals.Copy(),
	}
	// Insert the methods into the module dictionary
	for _, method := range methods {
		m.Globals[method.Name] = method
	}
	// Set some module globals
	m.Globals["__name__"] = String(name)
	m.Globals["__doc__"] = String(doc)
	m.Globals["__package__"] = None
	// Register the module
	modules[name] = m
	// Make a note of some modules
	switch name {
	case "builtins":
		Builtins = m
	case "importlib":
		Importlib = m
	}
	fmt.Printf("Registering module %q\n", name)
	return m
}

// Calls a named method of a module
func (m *Module) Call(name string, args Tuple, kwargs StringDict) Object {
	return Call(m.M__getattribute__(name), args, kwargs)
}

// Get an attribute from the module
func (m *Module) M__getattribute__(name string) Object {
	res, ok := m.Globals[name]
	if !ok {
		panic(ExceptionNewf(AttributeError, "module '%s' has no attribute '%s'", m.Name, name))
	}
	return res
}

func (m *Module) M__setattr__(name string, value Object) Object {
	m.Globals[name] = value
	return None
}

// Interfaces
var _ I__getattribute__ = (*Module)(nil)
var _ I__setattr__ = (*Module)(nil)
