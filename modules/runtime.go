package modules

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-python/gpython/marshal"
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/vm"

	_ "github.com/go-python/gpython/builtin"
	_ "github.com/go-python/gpython/math"
	_ "github.com/go-python/gpython/sys"
	_ "github.com/go-python/gpython/time"
)

func init() {
	// Assign the base-level py.Ctx creation function while also preventing an import cycle.
	py.NewCtx = NewCtx
}

// ctx implements py.Ctx
type ctx struct {
	store *py.Store
	opts  py.CtxOpts
}

// See py.Ctx interface
func NewCtx(opts py.CtxOpts) py.Ctx {
	ctx := &ctx{
		opts: opts,
	}

	ctx.store = py.NewStore()

	py.Import(ctx, "builtins", "sys")

	sys_mod := ctx.Store().MustGetModule("sys")
	sys_mod.Globals["argv"] = py.NewListFromStrings(opts.SysArgs)
	sys_mod.Globals["path"] = py.NewListFromStrings(opts.SysPaths)

	return ctx
}

func (ctx *ctx) RunFile(runPath string, opts py.RunOpts) (*py.Module, error) {
	tryPaths := defaultPaths
	if opts.UseSysPaths {
		tryPaths = ctx.Store().MustGetModule("sys").Globals["path"].(*py.List).Items
	}

	var codeObj py.Object
	var fileDesc string

	err := resolveRunPath(runPath, opts, tryPaths, func(fpath string) (bool, error) {

		stat, err := os.Stat(fpath)
		if err == nil && stat.IsDir() {
			// FIXME this is a massive simplification!
			fpath = path.Join(fpath, "__init__.py")
			_, err = os.Stat(fpath)
		}

		ext := strings.ToLower(filepath.Ext(fpath))
		if ext == "" && os.IsNotExist(err) {
			fpath += ".py"
			ext = ".py"
			_, err = os.Stat(fpath)
		}

		// Keep searching while we get FNFs, stop on an error
		if err != nil {
			if os.IsNotExist(err) {
				return true, nil
			}
			err = py.ExceptionNewf(py.OSError, "Error accessing %q: %v", fpath, err)
			return false, err
		}

		if ext == ".py" {
			var pySrc []byte
			pySrc, err = ioutil.ReadFile(fpath)
			if err != nil {
				return false, py.ExceptionNewf(py.OSError, "Error reading %q: %v", fpath, err)
			}

			codeObj, err = py.Compile(string(pySrc), fpath, py.ExecMode, 0, true)
			if err != nil {
				return false, err
			}
		} else if ext == ".pyc" {
			var file *os.File
			file, err = os.Open(fpath)
			if err != nil {
				return false, py.ExceptionNewf(py.OSError, "Error opening %q: %v", fpath, err)
			}
			defer file.Close()
			codeObj, err = marshal.ReadPyc(file)
			if err != nil {
				return false, py.ExceptionNewf(py.ImportError, "Failed to marshal %q: %v", fpath, err)
			}
		}

		fileDesc = fpath
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	var code *py.Code
	if codeObj != nil {
		code, _ = codeObj.(*py.Code)
	}
	if code == nil {
		return nil, py.ExceptionNewf(py.AssertionError, "Missing code object")
	}

	if opts.HostModule == nil {
		opts.HostModule = ctx.Store().NewModule(ctx, py.ModuleInfo{
			Name:     opts.ModuleName,
			FileDesc: fileDesc,
		}, nil, nil)
	}

	_, err = vm.EvalCode(ctx, code, opts.HostModule.Globals, opts.HostModule.Globals, nil, nil, nil, nil, nil)
	return opts.HostModule, err
}

var defaultPaths = []py.Object{
	py.String("."),
}

func resolveRunPath(runPath string, opts py.RunOpts, pathObjs []py.Object, tryPath func(pyPath string) (bool, error)) error {
	var cwd string

	// Remove trailing slash if present
	if runPath[len(runPath)-1] == '/' {
		runPath = runPath[:len(runPath)-1]
	}

	var err error

	cont := true

	for _, pathObj := range pathObjs {
		pathStr, ok := pathObj.(py.String)
		if !ok {
			continue
		}

		// If an absolute path, just try that.
		// Otherwise, check from the passed current dir then check from the current working dir.
		fpath := path.Join(string(pathStr), runPath)
		if filepath.IsAbs(fpath) {
			cont, err = tryPath(fpath)
		} else {
			if len(opts.CurDir) > 0 {
				subPath := path.Join(opts.CurDir, fpath)
				cont, err = tryPath(subPath)
			}
			if cont && err == nil {
				if len(cwd) == 0 {
					cwd, _ = os.Getwd()
				}
				subPath := path.Join(cwd, fpath)
				cont, err = tryPath(subPath)
			}
		}
		if !cont {
			break
		}
	}

	if err != nil {
		return err
	}

	if cont {
		return py.ExceptionNewf(py.FileNotFoundError, "Failed to resolve %q", runPath)
	}

	return err
}

func (ctx *ctx) RunCode(code *py.Code, globals, locals py.StringDict, closure py.Tuple) (py.Object, error) {
	return vm.EvalCode(ctx, code, globals, locals, nil, nil, nil, nil, closure)
}

func (ctx *ctx) GetModule(moduleName string) (*py.Module, error) {
	return ctx.store.GetModule(moduleName)
}

func (ctx *ctx) Store() *py.Store {
	return ctx.store
}
