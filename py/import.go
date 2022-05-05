// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Import modules

package py

import (
	"path"
	"strings"
)

func Import(ctx Context, names ...string) error {
	for _, name := range names {
		_, err := ImportModuleLevelObject(ctx, name, nil, nil, nil, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

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
func ImportModuleLevelObject(ctx Context, name string, globals, locals StringDict, fromlist Tuple, level int) (Object, error) {
	// Module already loaded - return that
	if module, err := ctx.GetModule(name); err == nil {
		return module, nil
	}

	// See if the module is a registered embeddded module that has not been loaded into this ctx yet.
	if impl := GetModuleImpl(name); impl != nil {
		module, err := ctx.ModuleInit(impl)
		if err != nil {
			return nil, err
		}
		return module, nil
	}

	if level != 0 {
		return nil, ExceptionNewf(SystemError, "Relative import not supported yet")
	}

	// Convert import's dot separators into path seps
	parts := strings.Split(name, ".")
	srcPathname := path.Join(parts...)

	opts := CompileOpts{
		UseSysPaths: true,
	}

	if fromFile, ok := globals["__file__"]; ok {
		opts.CurDir = path.Dir(string(fromFile.(String)))
	}

	module, err := RunFile(ctx, srcPathname, opts, name)
	if err != nil {
		return nil, err
	}

	return module, nil
}

// Straight port of the python code
//
// This calls functins from _bootstrap.py which is a frozen module
//
// Too much functionality for the moment
func XImportModuleLevelObject(ctx Context, nameObj, given_globals, locals, given_fromlist Object, level int) (Object, error) {
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
	var err error
	store := ctx.Store()

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
			return nil, ExceptionNewf(TypeError, "globals must be a dict")
		}
	}

	if given_fromlist == nil || given_fromlist == None {
		fromlist = Tuple{}
	} else {
		fromlist, err = SequenceTuple(given_fromlist)
		if err != nil {
			return nil, err
		}
	}
	if nameObj == nil {
		return nil, ExceptionNewf(ValueError, "Empty module name")
	}

	// The below code is importlib.__import__() & _gcd_import(), ported to Go
	// for added performance.

	_, ok = nameObj.(String)
	if !ok {
		return nil, ExceptionNewf(TypeError, "module name must be a string")
	}
	name = string(nameObj.(String))

	if level < 0 {
		return nil, ExceptionNewf(ValueError, "level must be >= 0")
	} else if level > 0 {
		PackageObj, ok = globals["__package__"]
		if ok && PackageObj != None {
			if _, ok = PackageObj.(String); !ok {
				return nil, ExceptionNewf(TypeError, "package must be a string")
			}
			Package = string(PackageObj.(String))
		} else {
			PackageObj, ok = globals["__name__"]
			if !ok {
				return nil, ExceptionNewf(KeyError, "'__name__' not in globals")
			} else if _, ok = PackageObj.(String); !ok {
				return nil, ExceptionNewf(TypeError, "__name__ must be a string")
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

		if _, err = ctx.GetModule(string(Package)); err != nil {
			return nil, ExceptionNewf(SystemError, "Parent module %q not loaded, cannot perform relative import", Package)
		}
	} else { // level == 0 */
		if len(name) == 0 {
			return nil, ExceptionNewf(ValueError, "Empty module name")
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
				return nil, ExceptionNewf(ValueError, "attempted relative import beyond top-level Package")
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
		builtins_import, ok = store.Builtins.Globals["__import__"]
		if !ok {
			return nil, ExceptionNewf(ImportError, "__import__ not found")
		}
	}

	mod, err = ctx.GetModule(abs_name)
	if err != nil || mod == None {
		return nil, ExceptionNewf(ImportError, "import of %q halted; None in sys.modules", abs_name)
	} else if err == nil {
		var value Object
		var err error
		initializing := false

		// Optimization: only call _bootstrap._lock_unlock_module() if
		// __initializing__ is true.
		// NOTE: because of this, __initializing__ must be set *before*
		// stuffing the new module in sys.modules.

		value, err = GetAttrString(mod, "__initializing__")
		if err == nil {
			x, err := MakeBool(value)
			if err != nil {
				return nil, err
			}
			initializing = bool(x.(Bool))
		}
		if initializing {
			// _bootstrap._lock_unlock_module() releases the import lock */
			_, err = store.Importlib.Call("_lock_unlock_module", Tuple{String(abs_name)}, nil)
			if err != nil {
				return nil, err
			}
			//	} else { // not initializing
			// FIXME locking
			// if _PyImport_ReleaseLock() < 0 {
			// 	return nil, ExceptionNewf(RuntimeError, "not holding the import lock")
			// }
		}
	} else {
		// _bootstrap._find_and_load() releases the import lock
		mod, err = store.Importlib.Call("_find_and_load", Tuple{String(abs_name), builtins_import}, nil)
		if err != nil {
			return nil, err
		}
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
				var err error
				final_mod, err = Call(builtins_import, Tuple{String(front)}, nil)
				if err != nil {
					return nil, err
				}
			} else {
				cut_off := len(name) - len(front)
				abs_name_len := len(abs_name)
				to_return := abs_name[:abs_name_len-cut_off]
				final_mod, err = ctx.GetModule(to_return)
				if err != nil {
					return nil, ExceptionNewf(KeyError, "%q not in sys.modules as expected", to_return)
				}
			}
		} else {
			final_mod = mod
		}
	} else {
		final_mod, err = store.Importlib.Call("_handle_fromlist", Tuple{mod, fromlist, builtins_import}, nil)
		if err != nil {
			return nil, err
		}

	}
	goto error

	//error_with_unlock:
	// FIXME defer?
	// if _PyImport_ReleaseLock() < 0 {
	// 	return nil, ExceptionNewf(RuntimeError, "not holding the import lock"
	// }
error:
	// FIXME defer?
	// if final_mod == nil {
	// 	remove_importlib_frames()
	// }
	return final_mod, nil
}

// The actual import code
func BuiltinImport(ctx Context, self Object, args Tuple, kwargs StringDict, currentGlobal StringDict) (Object, error) {
	kwlist := []string{"name", "globals", "locals", "fromlist", "level"}
	var name Object
	var globals Object = currentGlobal
	var locals Object = NewStringDict()
	var fromlist Object = Tuple{}
	var level Object = Int(0)

	err := ParseTupleAndKeywords(args, kwargs, "U|OOOi:__import__", kwlist, &name, &globals, &locals, &fromlist, &level)
	if err != nil {
		return nil, err
	}
	if fromlist == None {
		fromlist = Tuple{}
	}
	return ImportModuleLevelObject(ctx, string(name.(String)), globals.(StringDict), locals.(StringDict), fromlist.(Tuple), int(level.(Int)))
}
