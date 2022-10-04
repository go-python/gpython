// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

type CompileMode string

const (
	ExecMode   CompileMode = "exec"   // Compile a module
	EvalMode   CompileMode = "eval"   // Compile an expression
	SingleMode CompileMode = "single" // Compile a single (interactive) statement
)

// Context is a gpython environment instance container, providing a high-level mechanism
// for multiple python interpreters to run concurrently without restriction.
//
// Context instances maintain completely independent environments, namely the modules that
// have been imported and their state.  Modules imported into a Context are instanced
// from a parent ModuleImpl.  For example, since Contexts each have their
// own sys module instance, each can set sys.path differently and independently.
//
// If you access a Context from multiple groutines, you are responsible that access is not concurrent,
// with the exception of Close() and Done().
//
// See examples/multi-context and examples/embedding.
type Context interface {

	// Resolves then compiles (if applicable) the given file system pathname into a py.Code ready to be executed.
	ResolveAndCompile(pathname string, opts CompileOpts) (CompileOut, error)

	// Creates a new py.Module instance and initializes ModuleImpl's code in the new module (if applicable).
	ModuleInit(impl *ModuleImpl) (*Module, error)

	// RunCode is a lower-level invocation to execute the given py.Code.
	// Blocks until execution is complete.
	RunCode(code *Code, globals, locals StringDict, closure Tuple) (result Object, err error)

	// Returns the named module for this context (or an error if not found)
	GetModule(moduleName string) (*Module, error)

	// Gereric access to this context's modules / state.
	Store() *ModuleStore

	// Close signals this context is about to go out of scope and any internal resources should be released.
	// Code execution on a py.Context that has been closed will result in an error.
	Close() error

	// Done returns a signal that can be used to detect when this Context has fully closed / completed.
	// If Close() is called while execution in progress, Done() will not signal until execution is complete.
	Done() <-chan struct{}
}

// CompileOpts specifies options for high-level compilation.
type CompileOpts struct {
	UseSysPaths bool   // If set, sys.path will be used to resolve relative pathnames
	CurDir      string // If non-empty, this is the path of the current working directory.  If empty, os.Getwd() is used.
}

// CompileOut the output of high-level compilation -- e.g. ResolveAndCompile()
type CompileOut struct {
	SrcPathname string // Resolved pathname the .py file that was compiled (if applicable)
	PycPathname string // Pathname of the .pyc file read and/or written (if applicable)
	FileDesc    string // Pathname to be used for a a module's "__file__" attrib
	Code        *Code  // Read/Output code object ready for execution
}

// DefaultCoreSysPaths specify default search paths for module sys
// This can be changed during runtime and plays nice with others using DefaultContextOpts()
var DefaultCoreSysPaths = []string{
	".",
	"lib",
}

// DefaultAuxSysPaths are secondary default search paths for module sys.
// This can be changed during runtime and plays nice with others using DefaultContextOpts()
// They are separated from the default core paths since they the more likley thing you will want to completely replace when using gpython.
var DefaultAuxSysPaths = []string{
	"/usr/lib/python3.4",
	"/usr/local/lib/python3.4/dist-packages",
	"/usr/lib/python3/dist-packages",
}

// ContextOpts specifies fundamental environment and input settings for creating a new py.Context
type ContextOpts struct {
	SysArgs  []string // sys.argv initializer
	SysPaths []string // sys.path initializer
}

var (
	// DefaultContextOpts should be the default opts created for py.NewContext.
	// Calling this ensure that you future proof you code for suggested/default settings.
	DefaultContextOpts = func() ContextOpts {
		opts := ContextOpts{
			SysPaths: DefaultCoreSysPaths,
		}
		opts.SysPaths = append(opts.SysPaths, DefaultAuxSysPaths...)
		return opts
	}

	// NewContext is a high-level call to create a new gpython interpreter context.
	// See type Context interface.
	NewContext func(opts ContextOpts) Context

	// Compiles a python buffer into a py.Code object.
	// Returns a py.Code object or otherwise an error.
	Compile func(src, srcDesc string, mode CompileMode, flags int, dont_inherit bool) (*Code, error)
)

// RunFile resolves the given pathname, compiles as needed, executes the code in the given module, and returns the Module to indicate success.
//
// See RunCode() for description of inModule.
func RunFile(ctx Context, pathname string, opts CompileOpts, inModule interface{}) (*Module, error) {
	out, err := ctx.ResolveAndCompile(pathname, opts)
	if err != nil {
		return nil, err
	}

	return RunCode(ctx, out.Code, out.FileDesc, inModule)
}

// RunSrc compiles the given python buffer and executes it within the given module and returns the Module to indicate success.
//
// See RunCode() for description of inModule.
func RunSrc(ctx Context, pySrc string, pySrcDesc string, inModule interface{}) (*Module, error) {
	if pySrcDesc == "" {
		pySrcDesc = "<run>"
	}
	code, err := Compile(pySrc+"\n", pySrcDesc, SingleMode, 0, true)
	if err != nil {
		return nil, err
	}

	return RunCode(ctx, code, pySrcDesc, inModule)
}

// RunCode executes the given code object within the given module and returns the Module to indicate success.
//
// If inModule is a *Module, then the code is run in that module.
//
// If inModule is nil, the code is run in a new __main__ module (and the new Module is returned).
//
// If inModule is a string, the code is run in a new module with the given name (and the new Module is returned).
func RunCode(ctx Context, code *Code, codeDesc string, inModule interface{}) (*Module, error) {
	var (
		module     *Module
		moduleName string
		err        error
	)

	createNew := false
	switch mod := inModule.(type) {

	case string:
		moduleName = mod
		createNew = true
	case nil:
		createNew = true
	case *Module:
		_, err = ctx.RunCode(code, mod.Globals, mod.Globals, nil)
		module = mod
	default:
		err = ExceptionNewf(TypeError, "unsupported module type: %v", inModule)
	}

	if err == nil && createNew {
		moduleImpl := ModuleImpl{
			Info: ModuleInfo{
				Name:     moduleName,
				FileDesc: codeDesc,
			},
			Code: code,
		}
		module, err = ctx.ModuleInit(&moduleImpl)
	}

	if err != nil {
		return nil, err
	}

	return module, nil
}
