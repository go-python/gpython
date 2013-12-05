// Import modules

package py

import (
	"strings"
)

// The workings of __import__
func ImportModuleLevelObject(nameObj, given_globals, locals, given_fromlist Object, level int) Object {
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
