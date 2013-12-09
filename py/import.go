// Import modules

package py

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// This will become sys.path one day ;-)
	modulePath = []string{"", "/usr/lib/python3.3", "/usr/local/lib/python3.3/dist-packages", "/usr/lib/python3/dist-packages"}
)

// The workings of __import__
//
// __import__(name, globals=None, locals=None, fromlist=(), level=0)
//
// This function is invoked by the import statement. It can be
// replaced (by importing the builtins module and assigning to
// builtins.__import__) in order to change semantics of the import
// statement, but doing so is strongly discouraged as it is usually
// simpler to use import hooks (see PEP 302) to attain the same goals
// and does not cause issues with code which assumes the default
// import implementation is in use. Direct use of __import__() is also
// discouraged in favor of importlib.import_module().
//
// The function imports the module name, potentially using the given
// globals and locals to determine how to interpret the name in a
// package context. The fromlist gives the names of objects or
// submodules that should be imported from the module given by
// name. The standard implementation does not use its locals argument
// at all, and uses its globals only to determine the package context
// of the import statement.
//
// level specifies whether to use absolute or relative imports. 0 (the
// default) means only perform absolute imports. Positive values for
// level indicate the number of parent directories to search relative
// to the directory of the module calling __import__() (see PEP 328
// for the details).
//
// When the name variable is of the form package.module, normally, the
// top-level package (the name up till the first dot) is returned, not
// the module named by name. However, when a non-empty fromlist
// argument is given, the module named by name is returned.
//
// For example, the statement import spam results in bytecode
// resembling the following code:
//
// spam = __import__('spam', globals(), locals(), [], 0)
// The statement import spam.ham results in this call:
//
// spam = __import__('spam.ham', globals(), locals(), [], 0)
//
// Note how __import__() returns the toplevel module here because this
// is the object that is bound to a name by the import statement.
//
// On the other hand, the statement from spam.ham import eggs, sausage
// as saus results in
//
// _temp = __import__('spam.ham', globals(), locals(), ['eggs', 'sausage'], 0)
// eggs = _temp.eggs
// saus = _temp.sausage
//
// Here, the spam.ham module is returned from __import__(). From this
// object, the names to import are retrieved and assigned to their
// respective names.
//
// If you simply want to import a module (potentially within a
// package) by name, use importlib.import_module().
//
// Changed in version 3.3: Negative values for level are no longer
// supported (which also changes the default value to 0).
func ImportModuleLevelObject(name string, globals, locals StringDict, fromlist Tuple, level int) Object {
	// Module already loaded - return that
	if module, ok := modules[name]; ok {
		return module
	}

	if level != 0 {
		panic("Relative import not supported yet")
	}

	parts := strings.Split(name, ".")
	pathParts := path.Join(parts...)

	for _, mpath := range modulePath {
		if mpath == "" {
			mpathObj, ok := globals["__file__"]
			if !ok {
				panic(ExceptionNewf(SystemError, "Couldn't find __file__ in globals"))
			}
			mpath = string(mpathObj.(String))
		}
		fullPath := path.Join(mpath, pathParts)
		// FIXME Read pyc/pyo too
		fullPath, err := filepath.Abs(fullPath + ".py")
		if err != nil {
			continue
		}
		// Check if file exists
		if _, err := os.Stat(fullPath); err == nil {
			str, err := ioutil.ReadFile(fullPath)
			if err != nil {
				panic(ExceptionNewf(OSError, "Couldn't read %q: %v", fullPath, err))
			}
			codeObj := Compile(string(str), fullPath, "exec", 0, true)
			code, ok := codeObj.(*Code)
			if !ok {
				panic(ExceptionNewf(ImportError, "Compile didn't return code object"))
			}
			module := NewModule(name, "", nil, nil)
			_, err = Run(module.Globals, module.Globals, code, nil)
			if err != nil {
				panic(err)
			}
			return module
		}
	}
	panic(ExceptionNewf(ImportError, "No module named '%s'", name))

	// Convert to absolute path if relative
	// Use __file__ from globals to work out what we are relative to

	// '' in path seems to mean use the current __file__

	// Find a valid path which we need to check for the correct __init__.py in subdirectories etc

	// Look for .py and .pyc files

	// Make absolute module path too if we can for sys.modules

	//How do we uniquely identify modules?

	// SystemError: Parent module '' not loaded, cannot perform relative import

}

// Straight port of the python code
//
// This calls functins from _bootstrap.py which is a frozen module
//
// Too much functionality for the moment
func XImportModuleLevelObject(nameObj, given_globals, locals, given_fromlist Object, level int) Object {
	var abs_name string
	var builtins_import Object
	var final_mod Object
	var mod Object
	var PackageObj Object
	var Package string
	var globals StringDict
	var fromlist Tuple
	var ok bool
	var name string

	// Make sure to use default values so as to not have
	// PyObject_CallMethodObjArgs() truncate the parameter list because of a
	// nil argument.
	if given_globals == nil {
		globals = StringDict{}
	} else {
		// Only have to care what given_globals is if it will be used
		// for something.
		globals, ok = given_globals.(StringDict)
		if level > 0 && !ok {
			panic(ExceptionNewf(TypeError, "globals must be a dict"))
		}
	}

	if given_fromlist == nil || given_fromlist == None {
		fromlist = Tuple{}
	} else {
		fromlist = SequenceTuple(given_fromlist)
	}
	if nameObj == nil {
		panic(ExceptionNewf(ValueError, "Empty module name"))
	}

	// The below code is importlib.__import__() & _gcd_import(), ported to Go
	// for added performance.

	_, ok = nameObj.(String)
	if !ok {
		panic(ExceptionNewf(TypeError, "module name must be a string"))
	}
	name = string(nameObj.(String))

	if level < 0 {
		panic(ExceptionNewf(ValueError, "level must be >= 0"))
	} else if level > 0 {
		PackageObj, ok = globals["__package__"]
		if ok && PackageObj != None {
			if _, ok = PackageObj.(String); !ok {
				panic(ExceptionNewf(TypeError, "package must be a string"))
			}
			Package = string(PackageObj.(String))
		} else {
			PackageObj, ok = globals["__name__"]
			if !ok {
				panic(ExceptionNewf(KeyError, "'__name__' not in globals"))
			} else if _, ok = PackageObj.(String); !ok {
				panic(ExceptionNewf(TypeError, "__name__ must be a string"))
			}
			Package = string(PackageObj.(String))

			if _, ok = globals["__path__"]; !ok {
				i := strings.LastIndex(string(Package), ".")
				if i < 0 {
					Package = ""
				} else {
					Package = Package[:i]
				}
			}
		}

		if _, ok = modules[string(Package)]; !ok {
			panic(ExceptionNewf(SystemError, "Parent module %q not loaded, cannot perform relative import", Package))
		}
	} else { // level == 0 */
		if len(name) == 0 {
			panic(ExceptionNewf(ValueError, "Empty module name"))
		}
		Package = ""
	}

	if level > 0 {
		last_dot := len(Package)
		var base string
		level_up := 1

		for level_up = 1; level_up < level; level_up += 1 {
			last_dot = strings.LastIndex(string(Package[:last_dot]), ".")
			if last_dot < 0 {
				panic(ExceptionNewf(ValueError, "attempted relative import beyond top-level Package"))
			}
		}

		base = Package[:last_dot]

		if len(name) > 0 {
			abs_name = strings.Join([]string{base, name}, ".")
		} else {
			abs_name = base
		}
	} else {
		abs_name = name
	}

	// FIXME _PyImport_AcquireLock()

	// From this point forward, goto error_with_unlock!
	builtins_import, ok = globals["__import__"]
	if !ok {
		builtins_import, ok = Builtins.Globals["__import__"]
		if !ok {
			panic(ExceptionNewf(ImportError, "__import__ not found"))
		}
	}

	mod, ok = modules[abs_name]
	if mod == None {
		panic(ExceptionNewf(ImportError, "import of %q halted; None in sys.modules", abs_name))
	} else if ok {
		var value Object
		var err error
		initializing := false

		// Optimization: only call _bootstrap._lock_unlock_module() if
		// __initializing__ is true.
		// NOTE: because of this, __initializing__ must be set *before*
		// stuffing the new module in sys.modules.

		value, err = GetAttrErr(mod, "__initializing__")
		if err == nil {
			initializing = bool(MakeBool(value).(Bool))
		}
		if initializing {
			// _bootstrap._lock_unlock_module() releases the import lock */
			value = Importlib.Call("_lock_unlock_module", Tuple{String(abs_name)}, nil)
		} else {
			// FIXME locking
			// if _PyImport_ReleaseLock() < 0 {
			// 	panic(ExceptionNewf(RuntimeError, "not holding the import lock"))
			// }
		}
	} else {
		// _bootstrap._find_and_load() releases the import lock
		mod = Importlib.Call("_find_and_load", Tuple{String(abs_name), builtins_import}, nil)
	}
	// From now on we don't hold the import lock anymore.

	if len(fromlist) == 0 {
		if level == 0 || len(name) > 0 {
			i := strings.Index(name, ".")
			if i < 0 {
				// No dot in module name, simple exit
				final_mod = mod
				goto error
			}
			front := name[:1]

			if level == 0 {
				final_mod = Call(builtins_import, Tuple{String(front)}, nil)
			} else {
				cut_off := len(name) - len(front)
				abs_name_len := len(abs_name)
				to_return := abs_name[:abs_name_len-cut_off]
				final_mod, ok = modules[to_return]
				if !ok {
					panic(ExceptionNewf(KeyError, "%q not in sys.modules as expected", to_return))
				}
			}
		} else {
			final_mod = mod
		}
	} else {
		final_mod = Importlib.Call("_handle_fromlist", Tuple{mod, fromlist, builtins_import}, nil)
	}
	goto error

	//error_with_unlock:
	// FIXME defer?
	// if _PyImport_ReleaseLock() < 0 {
	// 	panic(ExceptionNewf(RuntimeError, "not holding the import lock")
	// }
error:
	// FIXME defer?
	// if final_mod == nil {
	// 	remove_importlib_frames()
	// }
	return final_mod
}
