package py

type CompileMode string

const (
	ExecMode   CompileMode = "exec"
	EvalMode   CompileMode = "eval"
	SingleMode CompileMode = "single"
)

type RunFlags int32

const (
	// RunOpts.FilePath is intelligently interepreted to load the appropriate pyc object (otherwise new code is generated from the implied .py file)
	SmartCodeAcquire RunFlags = 0x01
)

// Ctx is gpython virtual environment instance ("context"), providing a mechanism
// for multple gpython interpreters to run concurrently without restriction.
//
// See BenchmarkCtx() in vm/vm_test.go
type Ctx interface {

	// These each initiate blocking execution.
	RunCode(code *Code, globals, locals StringDict, closure Tuple) (res Object, err error)

	// Runs a given file in the given host module.
	RunFile(runPath string, opts RunOpts) (*Module, error)

	// Execution of any of the above will stop when the next opcode runs
	// @@TODO
	// SignalHalt()

	// Returns the named module for this context (or an error if not found)
	GetModule(moduleName string) (*Module, error)

	// WIP
	Store() *Store
}

const (
	SrcFileExt  = ".py"
	CodeFileExt = ".pyc"

	DefaultModuleName = "__main__"
)

type StdLib int32

const (
	Lib_sys StdLib = 1 << iota
	Lib_time

	CoreLibs = Lib_sys | Lib_time
)

type RunOpts struct {
	HostModule  *Module // Host module to execute within (if nil, a new module is created)
	ModuleName  string  // If HostModule == nil, this is the name of the newly created module.  If nil, "__main__" is used.
	CurDir      string  // If non-nil, this is the path of the current working directory.  If nil, os.Getwd() is used
	UseSysPaths bool
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

// High-level entry points
var (
	DefaultCtxOpts = func() CtxOpts {
		opts := CtxOpts{
			SysPaths: CoreSysPaths,
		}
		opts.SysPaths = append(opts.SysPaths, AuxSysPaths...)
		return opts
	}
	NewCtx  func(opts CtxOpts) Ctx
	Compile func(str, filename string, mode CompileMode, flags int, dont_inherit bool) (Object, error)
)
