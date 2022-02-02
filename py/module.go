// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Module objects

package py

import (
	"fmt"
	"sync"
)

type ModuleFlags int32

const (
	// ShareModule signals that an embedded module is threadsafe and read-only, meaninging it could be shared across multiple py.Context instances (for efficiency).
	// Otherwise, ModuleImpl will create a separate py.Module instance for each py.Context that imports it.
	// This should be used with extreme caution since any module mutation (write) means possible cross-context data corruption.
	ShareModule ModuleFlags = 0x01

	MainModuleName = "__main__"
)

// ModuleInfo contains info and about a module and can specify flags that affect how it is imported into a py.Context
type ModuleInfo struct {
	Name     string // __name__ (if nil, "__main__" is used)
	Doc      string // __doc__
	FileDesc string // __file__
	Flags    ModuleFlags
}

// ModuleImpl is used for modules that are ready to be imported into a py.Context.
// The model is that a ModuleImpl is read-only and instantiates a Module into a py.Context when imported.
//
// By convention, .Code is executed when a module instance is initialized. If nil,
// then .CodeBuf or .CodeSrc will be auto-compiled to set .Code.
type ModuleImpl struct {
	Info            ModuleInfo
	Methods         []*Method     // Module-bound global method functions
	Globals         StringDict    // Module-bound global variables
	CodeSrc         string        // Module code body (source code to be compiled)
	CodeBuf         []byte        // Module code body (serialized py.Code object)
	Code            *Code         // Module code body
	OnContextClosed func(*Module) // Callback for when a py.Context is closing to release resources
}

// ModuleStore is a container of Module imported into an owning py.Context.
type ModuleStore struct {
	// Registry of installed modules
	modules map[string]*Module
	// Builtin module
	Builtins *Module
	// this should be the frozen module importlib/_bootstrap.py generated
	// by Modules/_freeze_importlib.c into Python/importlib.h
	Importlib *Module
}

func RegisterModule(module *ModuleImpl) {
	gRuntime.RegisterModule(module)
}

func GetModuleImpl(moduleName string) *ModuleImpl {
	gRuntime.mu.RLock()
	defer gRuntime.mu.RUnlock()
	impl := gRuntime.ModuleImpls[moduleName]
	return impl
}

type Runtime struct {
	mu          sync.RWMutex
	ModuleImpls map[string]*ModuleImpl
}

var gRuntime = Runtime{
	ModuleImpls: make(map[string]*ModuleImpl),
}

func (rt *Runtime) RegisterModule(impl *ModuleImpl) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.ModuleImpls[impl.Info.Name] = impl
}

func NewModuleStore() *ModuleStore {
	return &ModuleStore{
		modules: make(map[string]*Module),
	}
}

// Module is a runtime instance of a ModuleImpl bound to the py.Context that imported it.
type Module struct {
	ModuleImpl *ModuleImpl // Parent implementation of this Module instance
	Globals    StringDict  // Initialized from ModuleImpl.Globals
	Context    Context     // Parent context that "owns" this Module instance
}

var ModuleType = NewType("module", "module object")

// Type of this object
func (o *Module) Type() *Type {
	return ModuleType
}

func (m *Module) M__repr__() (Object, error) {
	name, ok := m.Globals["__name__"].(String)
	if !ok {
		name = "???"
	}
	return String(fmt.Sprintf("<module %s>", string(name))), nil
}

// Get the Dict
func (m *Module) GetDict() StringDict {
	return m.Globals
}

// Calls a named method of a module
func (m *Module) Call(name string, args Tuple, kwargs StringDict) (Object, error) {
	attr, err := GetAttrString(m, name)
	if err != nil {
		return nil, err
	}
	return Call(attr, args, kwargs)
}

// Interfaces
var _ IGetDict = (*Module)(nil)

// NewModule adds a new Module instance to this ModuleStore.
// Each given Method prototype is used to create a new "live" Method bound this the newly created Module.
// This func also sets appropriate module global attribs based on the given ModuleInfo (e.g. __name__).
func (store *ModuleStore) NewModule(ctx Context, impl *ModuleImpl) (*Module, error) {
	name := impl.Info.Name
	if name == "" {
		name = MainModuleName
	}
	m := &Module{
		ModuleImpl: impl,
		Globals:    impl.Globals.Copy(),
		Context:    ctx,
	}
	// Insert the methods into the module dictionary
	// Copy each method an insert each "live" with a ptr back to the module (which can also lead us to the host Context)
	for _, method := range impl.Methods {
		methodInst := new(Method)
		*methodInst = *method
		methodInst.Module = m
		m.Globals[method.Name] = methodInst
	}
	// Set some module globals
	m.Globals["__name__"] = String(name)
	m.Globals["__doc__"] = String(impl.Info.Doc)
	m.Globals["__package__"] = None
	if len(impl.Info.FileDesc) > 0 {
		m.Globals["__file__"] = String(impl.Info.FileDesc)
	}
	// Register the module
	store.modules[name] = m
	// Make a note of some modules
	switch name {
	case "builtins":
		store.Builtins = m
	case "importlib":
		store.Importlib = m
	}
	// fmt.Printf("Registered module %q\n", moduleName)
	return m, nil
}

// Gets a module
func (store *ModuleStore) GetModule(name string) (*Module, error) {
	m, ok := store.modules[name]
	if !ok {
		return nil, ExceptionNewf(ImportError, "Module '%s' not found", name)
	}
	return m, nil
}

// Gets a module or panics
func (store *ModuleStore) MustGetModule(name string) *Module {
	m, err := store.GetModule(name)
	if err != nil {
		panic(err)
	}
	return m
}

// OnContextClosed signals all module instances that the parent py.Context has closed
func (store *ModuleStore) OnContextClosed() {
	for _, m := range store.modules {
		if m.ModuleImpl.OnContextClosed != nil {
			m.ModuleImpl.OnContextClosed(m)
		}
	}
}
