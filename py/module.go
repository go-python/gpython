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
	// Register the module
	modules[name] = m
	// Make a note of the builtin module
	if name == "builtins" {
		Builtins = m
	}
	fmt.Printf("Registering module %q\n", name)
	return m
}
