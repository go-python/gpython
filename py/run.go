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

// type Sys interface {
// 	ResetArfs
// }

// Rename to Ctx?
type Ctx interface {

	// These each initiate blocking execution.
	RunCode(code *Code, globals, locals StringDict, closure Tuple) (res Object, err error)

	// Runs a given file in the given host module.
	RunFile(runPath string, opts RunOpts) (*Module, error)

	// // Execution of any of the above will stop when the next opcode runs
	// SignalHalt()

	//Sys() Sys

	GetModule(moduleName string) (*Module, error)
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

var DefaultCtxOpts = CtxOpts{}

type CtxOpts struct {
	Args       []string
	SetSysArgs bool
}

// High-level entry points
var (
	NewCtx  func(opts CtxOpts) Ctx
	Compile func(str, filename string, mode CompileMode, flags int, dont_inherit bool) (Object, error)
)
