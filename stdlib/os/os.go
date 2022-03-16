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
		py.MustNewMethod("getcwd", getCwd, 0, "Get the current working directory"),
		py.MustNewMethod("getcwdb", getCwdb, 0, "Get the current working directory in a byte slice"),
		py.MustNewMethod("chdir", chdir, 0, "Change the current working directory"),
		py.MustNewMethod("getenv", getenv, 0, "Return the value of the environment variable key if it exists, or default if it doesn’t. key, default and the result are str."),
		py.MustNewMethod("getpid", getpid, 0, "Return the current process id."),
		py.MustNewMethod("putenv", putenv, 0, "Set the environment variable named key to the string value."),
		py.MustNewMethod("unsetenv", unsetenv, 0, "Unset (delete) the environment variable named key."),
		py.MustNewMethod("_exit", _exit, 0, "Immediate program termination."),
		py.MustNewMethod("system", system, 0, "Run shell commands, prints stdout directly to deault"),
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
