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
	// Set for modules that are threadsafe, stateless, and/or can be shared across multiple py.Ctx instances (for efficiency).
	// Otherwise, a separate module instance is created for each py.Ctx that imports it.
	ShareModule ModuleFlags = 0x01 // @@TODO

	MainModuleName = "__main__"
)

// ModuleInfo contains info and about a module and can specify flags that affect how it is imported into a py.Ctx
type ModuleInfo struct {
	Name     string // __name__ (if nil, "__main__" is used)
	Doc      string // __doc__
	FileDesc string // __file__
	Flags    ModuleFlags
}

// ModuleImpl is used for modules that are ready to be imported into a py.Ctx.
// If a module is threadsafe and stateless it can be shared across multiple py.Ctx instances (for efficiency).
// By convention, .Code is executed when a module instance is initialized.
// If .Code == nil, then .CodeBuf or .CodeSrc will be auto-compiled to set .Code.
type ModuleImpl struct {
	Info    ModuleInfo
	Methods []*Method
	Globals StringDict
	CodeSrc string // Module code body (py source code to be compiled)
	CodeBuf []byte // Module code body (serialized py.Code object)
	Code    *Code  // Module code body
}

type Store struct {
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

func NewStore() *Store {
	return &Store{
		modules: make(map[string]*Module),
	}
}

// A python Module object that has been initted for a given py.Ctx
type Module struct {
	ModuleInfo

	Globals StringDict
	Ctx     Ctx
}

var ModuleType = NewType("module", "module object")

// Type of this object
func (o *Module) Type() *Type {
	return ModuleType
}

func (m *Module) M__repr__() (Object, error) {
	return String(fmt.Sprintf("<module %s>", m.Name)), nil
}

// Get the Dict
func (m *Module) GetDict() StringDict {
	return m.Globals
}

// Define a new module
func (store *Store) NewModule(ctx Ctx, info ModuleInfo, methods []*Method, globals StringDict) (*Module, error) {
	if info.Name == "" {
		info.Name = MainModuleName
	}
	m := &Module{
		ModuleInfo: info,
		Globals:    globals.Copy(),
		Ctx:        ctx,
	}
	// Insert the methods into the module dictionary
	// Copy each method an insert each "live" with a ptr back to the module (which can also lead us to the host Ctx)
	for _, method := range methods {
		methodInst := new(Method)
		*methodInst = *method
		methodInst.Module = m
		m.Globals[method.Name] = methodInst
	}
	// Set some module globals
	m.Globals["__name__"] = String(info.Name)
	m.Globals["__doc__"] = String(info.Doc)
	m.Globals["__package__"] = None
	if len(info.FileDesc) > 0 {
		m.Globals["__file__"] = String(info.FileDesc)
	}
	// Register the module
	store.modules[info.Name] = m
	// Make a note of some modules
	switch info.Name {
	case "builtins":
		store.Builtins = m
	case "importlib":
		store.Importlib = m
	}
	// fmt.Printf("Registering module %q\n", name)
	return m, nil
}

// Gets a module
func (store *Store) GetModule(name string) (*Module, error) {
	m, ok := store.modules[name]
	if !ok {
		return nil, ExceptionNewf(ImportError, "Module '%s' not found", name)
	}
	return m, nil
}

// Gets a module or panics
func (store *Store) MustGetModule(name string) *Module {
	m, err := store.GetModule(name)
	if err != nil {
		panic(err)
	}
	return m
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
