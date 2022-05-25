// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package os implements the Python os module.
package os

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-python/gpython/py"
)

var (
	osSep     = py.String("/")
	osName    = py.String("posix")
	osPathsep = py.String(":")
	osLinesep = py.String("\n")
	osDefpath = py.String(":/bin:/usr/bin")
	osDevnull = py.String("/dev/null")

	osAltsep py.Object = py.None
)

func initGlobals() {
	switch runtime.GOOS {
	case "android":
		osName = py.String("java")
	case "windows":
		osSep = py.String(`\`)
		osName = py.String("nt")
		osPathsep = py.String(";")
		osLinesep = py.String("\r\n")
		osDefpath = py.String(`C:\bin`)
		osDevnull = py.String("nul")
		osAltsep = py.String("/")
	}
}

func init() {
	initGlobals()

	methods := []*py.Method{
		py.MustNewMethod("_exit", _exit, 0, "Immediate program termination."),
		py.MustNewMethod("getcwd", getCwd, 0, "Get the current working directory"),
		py.MustNewMethod("getcwdb", getCwdb, 0, "Get the current working directory in a byte slice"),
		py.MustNewMethod("chdir", chdir, 0, "Change the current working directory"),
		py.MustNewMethod("getenv", getenv, 0, "Return the value of the environment variable key if it exists, or default if it doesn’t. key, default and the result are str."),
		py.MustNewMethod("getpid", getpid, 0, "Return the current process id."),
		py.MustNewMethod("makedirs", makedirs, 0, makedirs_doc),
		py.MustNewMethod("mkdir", mkdir, 0, mkdir_doc),
		py.MustNewMethod("putenv", putenv, 0, "Set the environment variable named key to the string value."),
		py.MustNewMethod("remove", remove, 0, remove_doc),
		py.MustNewMethod("removedirs", removedirs, 0, removedirs_doc),
		py.MustNewMethod("rmdir", rmdir, 0, rmdir_doc),
		py.MustNewMethod("system", system, 0, "Run shell commands, prints stdout directly to default"),
		py.MustNewMethod("unsetenv", unsetenv, 0, "Unset (delete) the environment variable named key."),
	}
	globals := py.StringDict{
		"error":   py.OSError,
		"environ": getEnvVariables(),
		"sep":     osSep,
		"name":    osName,
		"curdir":  py.String("."),
		"pardir":  py.String(".."),
		"extsep":  py.String("."),
		"altsep":  osAltsep,
		"pathsep": osPathsep,
		"linesep": osLinesep,
		"defpath": osDefpath,
		"devnull": osDevnull,
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "os",
			Doc:  "Miscellaneous operating system interfaces",
		},
		Methods: methods,
		Globals: globals,
	})
}

// getEnvVariables returns the dictionary of environment variables.
func getEnvVariables() py.StringDict {
	vs := os.Environ()
	dict := py.NewStringDictSized(len(vs))
	for _, evar := range vs {
		key_value := strings.SplitN(evar, "=", 2) // returns a []string containing [key,value]
		dict.M__setitem__(py.String(key_value[0]), py.String(key_value[1]))
	}

	return dict
}

// getCwd returns the current working directory.
func getCwd(self py.Object, args py.Tuple) (py.Object, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, py.ExceptionNewf(py.OSError, "Unable to get current working directory.")
	}
	return py.String(dir), nil
}

// getCwdb returns the current working directory as a byte list.
func getCwdb(self py.Object, args py.Tuple) (py.Object, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, py.ExceptionNewf(py.OSError, "Unable to get current working directory.")
	}
	return py.Bytes(dir), nil
}

// chdir changes the current working directory to the provided path.
func chdir(self py.Object, args py.Tuple) (py.Object, error) {
	if len(args) == 0 {
		return nil, py.ExceptionNewf(py.TypeError, "Missing required argument 'path' (pos 1)")
	}
	dir, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected, not "+args[0].Type().Name)
	}
	err := os.Chdir(string(dir))
	if err != nil {
		return nil, py.ExceptionNewf(py.NotADirectoryError, "Couldn't change cwd; "+err.Error())
	}
	return py.None, nil
}

// getenv returns the value of the environment variable key.
// If no such environment variable exists and a default value was provided, that value is returned.
func getenv(self py.Object, args py.Tuple) (py.Object, error) {
	if len(args) < 1 {
		return nil, py.ExceptionNewf(py.TypeError, "missing one required argument: 'name:str'")
	}
	k, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected (pos 1), not "+args[0].Type().Name)
	}
	v, ok := os.LookupEnv(string(k))
	if ok {
		return py.String(v), nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return py.None, nil
}

// getpid returns the current process id.
func getpid(self py.Object, args py.Tuple) (py.Object, error) {
	return py.Int(os.Getpid()), nil
}

const makedirs_doc = `makedirs(name [, mode=0o777][, exist_ok=False])

Super-mkdir; create a leaf directory and all intermediate ones.  Works like
mkdir, except that any intermediate path segment (not just the rightmost)
will be created if it does not exist. If the target directory already
exists, raise an OSError if exist_ok is False. Otherwise no exception is
raised.  This is recursive.`

func makedirs(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pypath py.Object
		pymode py.Object = py.Int(0o777)
		pyok   py.Object = py.False
	)
	err := py.ParseTupleAndKeywords(
		args, kwargs,
		"s#|ip:makedirs", []string{"path", "mode", "exist_ok"},
		&pypath, &pymode, &pyok,
	)
	if err != nil {
		return nil, err
	}

	var (
		path = ""
		mode = os.FileMode(pymode.(py.Int))
	)
	switch v := pypath.(type) {
	case py.String:
		path = string(v)
	case py.Bytes:
		path = string(v)
	}

	if pyok.(py.Bool) == py.False {
		// check if leaf exists.
		_, err := os.Stat(path)
		// FIXME(sbinet): handle other errors.
		if err == nil {
			return nil, py.ExceptionNewf(py.FileExistsError, "File exists: '%s'", path)
		}
	}

	err = os.MkdirAll(path, mode)
	if err != nil {
		return nil, err
	}

	return py.None, nil
}

const mkdir_doc = `Create a directory.

If dir_fd is not None, it should be a file descriptor open to a directory,
  and path should be relative; path will then be relative to that directory.
dir_fd may not be implemented on your platform.
  If it is unavailable, using it will raise a NotImplementedError.

The mode argument is ignored on Windows.`

func mkdir(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pypath  py.Object
		pymode  py.Object = py.Int(511)
		pydirfd py.Object = py.None
	)
	err := py.ParseTupleAndKeywords(
		args, kwargs,
		"s#|ii:mkdir", []string{"path", "mode", "dir_fd"},
		&pypath, &pymode, &pydirfd,
	)
	if err != nil {
		return nil, err
	}

	var (
		path = ""
		mode = os.FileMode(pymode.(py.Int))
	)
	switch v := pypath.(type) {
	case py.String:
		path = string(v)
	case py.Bytes:
		path = string(v)
	}

	if pydirfd != py.None {
		// FIXME(sbinet)
		return nil, py.ExceptionNewf(py.NotImplementedError, "mkdir(dir_fd=XXX) not implemented")
	}

	err = os.Mkdir(path, mode)
	if err != nil {
		return nil, err
	}

	return py.None, nil
}

// putenv sets the value of an environment variable named by the key.
func putenv(self py.Object, args py.Tuple) (py.Object, error) {
	if len(args) != 2 {
		return nil, py.ExceptionNewf(py.TypeError, "missing required arguments: 'key:str' and 'value:str'")
	}
	k, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected (pos 1), not "+args[0].Type().Name)
	}
	v, ok := args[1].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected (pos 2), not "+args[1].Type().Name)
	}
	err := os.Setenv(string(k), string(v))
	if err != nil {
		return nil, py.ExceptionNewf(py.OSError, "Unable to set enviroment variable")
	}
	return py.None, nil
}

// Unset (delete) the environment variable named key.
func unsetenv(self py.Object, args py.Tuple) (py.Object, error) {
	if len(args) != 1 {
		return nil, py.ExceptionNewf(py.TypeError, "missing one required argument: 'key:str'")
	}
	k, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected (pos 1), not "+args[0].Type().Name)
	}
	err := os.Unsetenv(string(k))
	if err != nil {
		return nil, py.ExceptionNewf(py.OSError, "Unable to unset enviroment variable")
	}
	return py.None, nil
}

// os._exit() immediate program termination; unlike sys.exit(), which raises a SystemExit, this function will termninate the program immediately.
func _exit(self py.Object, args py.Tuple) (py.Object, error) { // can never return
	if len(args) == 0 {
		os.Exit(0)
	}
	arg, ok := args[0].(py.Int)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "expected int (pos 1), not "+args[0].Type().Name)
	}
	os.Exit(int(arg))
	return nil, nil
}

const remove_doc = `Remove a file (same as unlink()).

If dir_fd is not None, it should be a file descriptor open to a directory,
  and path should be relative; path will then be relative to that directory.
dir_fd may not be implemented on your platform.
  If it is unavailable, using it will raise a NotImplementedError.`

func remove(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pypath py.Object
		pydir  py.Object = py.None
	)
	err := py.ParseTupleAndKeywords(args, kwargs, "s#|i:remove", []string{"path", "dir_fd"}, &pypath, &pydir)
	if err != nil {
		return nil, err
	}

	if pydir != py.None {
		// FIXME(sbinet) ?
		return nil, py.ExceptionNewf(py.NotImplementedError, "remove(dir_fd=XXX) not implemented")
	}

	var name string
	switch v := pypath.(type) {
	case py.String:
		name = string(v)
	case py.Bytes:
		name = string(v)
	}

	err = os.Remove(name)
	if err != nil {
		return nil, err
	}

	return py.None, nil
}

const removedirs_doc = `removedirs(name)

Super-rmdir; remove a leaf directory and all empty intermediate
ones.  Works like rmdir except that, if the leaf directory is
successfully removed, directories corresponding to rightmost path
segments will be pruned away until either the whole path is
consumed or an error occurs.  Errors during this latter phase are
ignored -- they generally mean that a directory was not empty.`

func removedirs(self py.Object, args py.Tuple) (py.Object, error) {
	var pypath py.Object
	err := py.ParseTuple(args, "s#:rmdir", &pypath)
	if err != nil {
		return nil, err
	}

	var name string
	switch v := pypath.(type) {
	case py.String:
		name = string(v)
	case py.Bytes:
		name = string(v)
	}

	err = os.RemoveAll(name)
	if err != nil {
		return nil, err
	}

	return py.None, nil
}

const rmdir_doc = `Remove a directory.

If dir_fd is not None, it should be a file descriptor open to a directory,
  and path should be relative; path will then be relative to that directory.
dir_fd may not be implemented on your platform.
  If it is unavailable, using it will raise a NotImplementedError.`

func rmdir(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pypath py.Object
		pydir  py.Object = py.None
	)
	err := py.ParseTupleAndKeywords(args, kwargs, "s#|i:rmdir", []string{"path", "dir_fd"}, &pypath, &pydir)
	if err != nil {
		return nil, err
	}

	if pydir != py.None {
		// FIXME(sbinet) ?
		return nil, py.ExceptionNewf(py.NotImplementedError, "rmdir(dir_fd=XXX) not implemented")
	}

	var name string
	switch v := pypath.(type) {
	case py.String:
		name = string(v)
	case py.Bytes:
		name = string(v)
	}

	err = os.Remove(name)
	if err != nil {
		return nil, err
	}

	return py.None, nil
}

// os.system(command string) this function runs a shell command and directs the output to standard output.
func system(self py.Object, args py.Tuple) (py.Object, error) {
	if len(args) != 1 {
		return nil, py.ExceptionNewf(py.TypeError, "missing one required argument: 'command:str'")
	}
	arg, ok := args[0].(py.String)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "str expected (pos 1), not "+args[0].Type().Name)
	}

	var command *exec.Cmd
	if runtime.GOOS != "windows" {
		command = exec.Command("/bin/sh", "-c", string(arg))
	} else {
		command = exec.Command("cmd.exe", string(arg))
	}
	outb, err := command.CombinedOutput() // - commbinedoutput to get both stderr and stdout -
	if err != nil {
		return nil, py.ExceptionNewf(py.OSError, err.Error())
	}
	ok = py.Println(self, string(outb))
	if !ok {
		return py.Int(1), nil
	}

	return py.Int(0), nil
}
