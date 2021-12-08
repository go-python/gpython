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
	//RunAsModule(opts RunOpts) (*Module, error)
	//RunCodeAsModule(code *Code, moduleName string, fileDesc string) (*Module, error)
	Run(globals, locals StringDict, code *Code, closure Tuple) (res Object, err error)
	RunFrame(frame *Frame) (res Object, err error)
	EvalCodeEx(co *Code, globals, locals StringDict, args []Object, kws StringDict, defs []Object, kwdefs StringDict, closure Tuple) (retval Object, err error)

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

type RunParams struct {
	ModuleInfo ModuleInfo // Newly created module to execute within (if ModuleInfo.Name is nil, then "__main__" is used)
	NoopOnFNF  bool       // If set and Pathname did not resolve, (nil, nil) is returned (vs an error)
	Silent     bool       // If set and an error occurs, no error info is printed
	Pathname   string
}

var DefaultCtxOpts = CtxOpts{}

type CtxOpts struct {
	Args       []string
	SetSysArgs bool
}

// Some well known objects
var (
	NewCtx  func(opts CtxOpts) Ctx
	
	// // Called at least once before using gpython; multiple calls to it have no effect.
	// // Called each time NewCtx is called
	// Init    func()
	Compile func(str, filename string, mode CompileMode, flags int, dont_inherit bool) (Object, error)
	Run     func(ctx Ctx, params RunParams) (*Module, error)
)
