// Copyright 2021 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modules

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-python/gpython/marshal"
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/vm"

	_ "github.com/go-python/gpython/builtin"
	_ "github.com/go-python/gpython/math"
	_ "github.com/go-python/gpython/sys"
	_ "github.com/go-python/gpython/time"
)

func init() {
	// Assign the base-level py.Context creation function while also preventing an import cycle.
	py.NewContext = NewContext
}

// context implements py.Context
type context struct {
	store     *py.ModuleStore
	opts      py.ContextOpts
	closeOnce sync.Once
	closing   bool
	closed    bool
	running   sync.WaitGroup
	done      chan struct{}
}

// See py.Context interface
func NewContext(opts py.ContextOpts) py.Context {
	ctx := &context{
		opts:    opts,
		done:    make(chan struct{}),
		closing: false,
		closed:  false,
	}

	ctx.store = py.NewModuleStore()

	py.Import(ctx, "builtins", "sys")

	sys_mod := ctx.Store().MustGetModule("sys")
	sys_mod.Globals["argv"] = py.NewListFromStrings(opts.SysArgs)
	sys_mod.Globals["path"] = py.NewListFromStrings(opts.SysPaths)

	return ctx
}

func (ctx *context) ModuleInit(impl *py.ModuleImpl) (*py.Module, error) {
	err := ctx.pushBusy()
	defer ctx.popBusy()
	if err != nil {
		return nil, err
	}

	if impl.Code == nil && len(impl.CodeSrc) > 0 {
		impl.Code, err = py.Compile(string(impl.CodeSrc), impl.Info.FileDesc, py.ExecMode, 0, true)
		if err != nil {
			return nil, err
		}
	}

	if impl.Code == nil && len(impl.CodeBuf) > 0 {
		codeBuf := bytes.NewBuffer(impl.CodeBuf)
		obj, err := marshal.ReadObject(codeBuf)
		if err != nil {
			return nil, err
		}
		impl.Code, _ = obj.(*py.Code)
		if impl.Code == nil {
			return nil, py.ExceptionNewf(py.AssertionError, "Embedded code did not produce a py.Code object")
		}
	}

	module, err := ctx.Store().NewModule(ctx, impl)
	if err != nil {
		return nil, err
	}

	if impl.Code != nil {
		_, err = ctx.RunCode(impl.Code, module.Globals, module.Globals, nil)
		if err != nil {
			return nil, err
		}
	}

	return module, nil
}

func (ctx *context) ResolveAndCompile(pathname string, opts py.CompileOpts) (py.CompileOut, error) {
	err := ctx.pushBusy()
	defer ctx.popBusy()
	if err != nil {
		return py.CompileOut{}, err
	}

	tryPaths := defaultPaths
	if opts.UseSysPaths {
		tryPaths = ctx.Store().MustGetModule("sys").Globals["path"].(*py.List).Items
	}

	out := py.CompileOut{}

	err = resolveRunPath(pathname, opts, tryPaths, func(fpath string) (bool, error) {

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

		switch ext {
		case ".py":
			var pySrc []byte
			pySrc, err = ioutil.ReadFile(fpath)
			if err != nil {
				return false, py.ExceptionNewf(py.OSError, "Error reading %q: %v", fpath, err)
			}

			out.Code, err = py.Compile(string(pySrc), fpath, py.ExecMode, 0, true)
			if err != nil {
				return false, err
			}
			out.SrcPathname = fpath
		case ".pyc":
			file, err := os.Open(fpath)
			if err != nil {
				return false, py.ExceptionNewf(py.OSError, "Error opening %q: %v", fpath, err)
			}
			defer file.Close()
			codeObj, err := marshal.ReadPyc(file)
			if err != nil {
				return false, py.ExceptionNewf(py.ImportError, "Failed to marshal %q: %v", fpath, err)
			}
			out.Code, _ = codeObj.(*py.Code)
			out.PycPathname = fpath
		}

		out.FileDesc = fpath
		return false, nil
	})

	if out.Code == nil && err == nil {
		err = py.ExceptionNewf(py.AssertionError, "Missing code object")
	}

	if err != nil {
		return py.CompileOut{}, err
	}

	return out, nil
}

func (ctx *context) pushBusy() error {
	if ctx.closed {
		return py.ExceptionNewf(py.RuntimeError, "Context closed")
	}
	ctx.running.Add(1)
	return nil
}

func (ctx *context) popBusy() {
	ctx.running.Done()
}

// Close -- see type py.Context
func (ctx *context) Close() error {
	ctx.closeOnce.Do(func() {
		ctx.closing = true
		ctx.running.Wait()
		ctx.closed = true

		// Give each module a chance to release resources
		ctx.store.OnContextClosed()
		close(ctx.done)
	})
	return nil
}

// Done -- see type py.Context
func (ctx *context) Done() <-chan struct{} {
	return ctx.done
}

var defaultPaths = []py.Object{
	py.String("."),
}

func resolveRunPath(runPath string, opts py.CompileOpts, pathObjs []py.Object, tryPath func(pyPath string) (bool, error)) error {
	runPath = strings.TrimSuffix(runPath, "/")

	var (
		err  error
		cwd  string
		cont = true
	)

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
				if cwd == "" {
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

func (ctx *context) RunCode(code *py.Code, globals, locals py.StringDict, closure py.Tuple) (py.Object, error) {
	err := ctx.pushBusy()
	defer ctx.popBusy()
	if err != nil {
		return nil, err
	}

	return vm.EvalCode(ctx, code, globals, locals, nil, nil, nil, nil, closure)
}

func (ctx *context) GetModule(moduleName string) (*py.Module, error) {
	return ctx.store.GetModule(moduleName)
}

func (ctx *context) Store() *py.ModuleStore {
	return ctx.store
}
