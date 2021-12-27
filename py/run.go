package py

type CompileMode string

const (
	ExecMode   CompileMode = "exec"   // Compile a module
	EvalMode   CompileMode = "eval"   // Compile an expression
	SingleMode CompileMode = "single" // Compile a single (interactive) statement
)

type RunFlags int32

const (
	// RunOpts.FilePath is intelligently interepreted to load the appropriate pyc object (otherwise new code is generated from the implied .py file)
	SmartCodeAcquire RunFlags = 0x01
)

// Ctx is gpython virtual environment instance ("context"), providing a mechanism
// for multiple gpython interpreters to run concurrently without restriction.
//
// In general, one creates a py.Ctx (via py.NewCtx) for each concurrent goroutine to be running an interpreter.
// In other words, ensure that a py.Ctx is never concurrently accessed across goroutines.
//
// RunFile() and RunCode() block until code execution is complete.
// In the future, they will abort early if the parent associated py.Ctx is signaled.
//
// See examples/multi-ctx
type Ctx interface {

	// Resolves then compiles (if applicable) the given file system pathname into a py.Code ready to be executed.
	ResolveAndCompile(pathname string, opts CompileOpts) (CompileOut, error)

	// Creates a new py.Module instance and initializes ModuleImpl's code in the new module (if applicable).
	ModuleInit(impl *ModuleImpl) (*Module, error)

	// RunCode is a lower-level invocation to execute the given py.Code.
	RunCode(code *Code, globals, locals StringDict, closure Tuple) (result Object, err error)

	// Execution of any of the above will stop when the next opcode runs
	// @@TODO
	// SignalHalt()

	// Returns the named module for this context (or an error if not found)
	GetModule(moduleName string) (*Module, error)

	// Gereric access to this context's modules / state.
	Store() *Store
}

const (
	SrcFileExt  = ".py"
	CodeFileExt = ".pyc"
)

type StdLib int32

const (
	Lib_sys StdLib = 1 << iota
	Lib_time

	CoreLibs = Lib_sys | Lib_time
)

type CompileOpts struct {
	UseSysPaths bool   // If set, sys.path will be used to resolve relative pathnames
	CurDir      string // If non-nil, this is the path of the current working directory.  If nil, os.Getwd() is used.
}

type CompileOut struct {
	SrcPathname string // Resolved pathname the .py file that was compiled (if applicable)
	PycPathname string // Pathname of the .pyc file read and/or written (if applicable)
	FileDesc    string // Pathname to be used for a a module's "__file__" attrib
	Code        *Code  // Read/Output code object ready for execution
}

// Can be changed during runtime and will \play nice with others using DefaultCtxOpts()
var CoreSysPaths = []string{
	".",
	"lib",
}

// Can be changed during runtime and will \play nice with others using DefaultCtxOpts()
var AuxSysPaths = []string{
	"/usr/lib/python3.4",
	"/usr/local/lib/python3.4/dist-packages",
	"/usr/lib/python3/dist-packages",
}

type CtxOpts struct {
	SysArgs  []string // sys.argv initializer
	SysPaths []string // sys.path initializer
}

var (
	// DefaultCtxOpts should be the default opts created for py.NewCtx.
	// Calling this ensure that you future proof you code for suggested/default settings.
	DefaultCtxOpts = func() CtxOpts {
		opts := CtxOpts{
			SysPaths: CoreSysPaths,
		}
		opts.SysPaths = append(opts.SysPaths, AuxSysPaths...)
		return opts
	}

	// NewCtx is a high-level call to create a new gpython interpreter context.
	// See type Ctx interface.
	NewCtx func(opts CtxOpts) Ctx

	// Compiles a python buffer into a py.Code object.
	// Returns a py.Code object or otherwise an error.
	Compile func(src, srcDesc string, mode CompileMode, flags int, dont_inherit bool) (*Code, error)
)

// Resolves the given pathname, compiles as needed, and runs that code in the given module, returning the Module to indicate success.
// If inModule is a *Module, then the code is run in that module.
// If inModule is nil, the code is run in a new __main__ module (and the new Module is returned).
// If inModule is a string, the code is run in a new module with the given name (and the new Module is returned).
func RunFile(ctx Ctx, pathname string, opts CompileOpts, inModule interface{}) (*Module, error) {
	out, err := ctx.ResolveAndCompile(pathname, opts)
	if err != nil {
		return nil, err
	}

	var moduleName string
	createNew := false
	var module *Module

	switch mod := inModule.(type) {

	case string:
		moduleName = mod
		createNew = true
	case nil:
		createNew = true
	case *Module:
		_, err = ctx.RunCode(out.Code, mod.Globals, mod.Globals, nil)
		module = mod
	default:
		err = ExceptionNewf(TypeError, "unsupported module type: %v", inModule)
	}

	if err == nil && createNew {
		moduleImpl := ModuleImpl{
			Info: ModuleInfo{
				Name:     moduleName,
				FileDesc: out.FileDesc,
			},
			Code: out.Code,
		}
		module, err = ctx.ModuleInit(&moduleImpl)
	}

	if err != nil {
		return nil, err
	}

	return module, nil
}
