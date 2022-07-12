// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"runtime"

	"github.com/go-python/gpython/py"
)

// These gpython py.Object type delcarations are the bridge between gpython and embedded Go types.
var (
	PyVacationStopType = py.NewType("Stop", "")
	PyVacationType     = py.NewType("Vacation", "")
)

// init is where you register your embedded module and attach methods to your embedded class types.
func init() {

	// For each of your embedded python types, attach instance methods.
	// When an instance method is invoked, the "self" py.Object is the instance.
	PyVacationStopType.Dict["Set"] = py.MustNewMethod("Set", VacationStop_Set, 0, "")
	PyVacationStopType.Dict["Get"] = py.MustNewMethod("Get", VacationStop_Get, 0, "")
	PyVacationType.Dict["add_stops"] = py.MustNewMethod("Vacation.add_stops", Vacation_add_stops, 0, "")
	PyVacationType.Dict["num_stops"] = py.MustNewMethod("Vacation.num_stops", Vacation_num_stops, 0, "")
	PyVacationType.Dict["get_stop"] = py.MustNewMethod("Vacation.get_stop", Vacation_get_stop, 0, "")

	// Bind methods attached at the module (global) level.
	// When these are invoked, the first py.Object param (typically "self") is the bound *Module instance.
	methods := []*py.Method{
		py.MustNewMethod("VacationStop_new", VacationStop_new, 0, ""),
		py.MustNewMethod("Vacation_new", Vacation_new, 0, ""),
	}

	// Register a ModuleImpl instance used by the gpython runtime to instantiate new py.Module when first imported.
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "mylib_go",
			Doc:  "Example embedded python module",
		},
		Methods: methods,
		Globals: py.StringDict{
			"PY_VERSION": py.String("Python 3.4 (github.com/go-python/gpython)"),
			"GO_VERSION": py.String(fmt.Sprintf("%s on %s %s", runtime.Version(), runtime.GOOS, runtime.GOARCH)),
			"MYLIB_VERS": py.String("Vacation 1.0 by Fletch F. Fletcher"),
		},
		OnContextClosed: func(instance *py.Module) {
			py.Println(instance, "<<< host py.Context of py.Module instance closing >>>\n+++")
		},
	})
}

// VacationStop is an example Go struct to embed.
type VacationStop struct {
	Desc      py.String
	NumNights py.Int
}

// Type comprises the py.Object interface, allowing a Go struct to be cast as a py.Object.
// Instance methods of an type are then attached to this type object
func (stop *VacationStop) Type() *py.Type {
	return PyVacationStopType
}

func (stop *VacationStop) M__str__() (py.Object, error) {
	line := fmt.Sprintf(" %-16v  |  %2v nights", stop.Desc, stop.NumNights)
	return py.String(line), nil
}

func (stop *VacationStop) M__repr__() (py.Object, error) {
	return stop.M__str__()
}

func VacationStop_new(module py.Object, args py.Tuple) (py.Object, error) {
	stop := &VacationStop{}
	VacationStop_Set(stop, args)
	return stop, nil
}

// VacationStop_Set is an embedded instance method of VacationStop
func VacationStop_Set(self py.Object, args py.Tuple) (py.Object, error) {
	stop := self.(*VacationStop)

	// Check out other convenience functions in py/util.go
	// Also available is py.ParseTuple(args, "si", ...)
	err := py.LoadTuple(args, []interface{}{&stop.Desc, &stop.NumNights})
	if err != nil {
		return nil, err
	}

	/* Alternative util func is ParseTuple():
	var desc, nights py.Object
	err := py.ParseTuple(args, "si", &desc, &nights)
	if err != nil {
		return nil, err
	}
	stop.Desc = desc.(py.String)
	stop.NumNights = desc.(py.Int)
	*/

	return py.None, nil
}

// VacationStop_Get is an embedded instance method of VacationStop
func VacationStop_Get(self py.Object, args py.Tuple) (py.Object, error) {
	stop := self.(*VacationStop)

	return py.Tuple{
		stop.Desc,
		stop.NumNights,
	}, nil
}

type Vacation struct {
	Stops  []*VacationStop
	MadeBy string
}

func (v *Vacation) Type() *py.Type {
	return PyVacationType
}

func Vacation_new(module py.Object, args py.Tuple) (py.Object, error) {
	v := &Vacation{}

	// For Module-bound methods, we have easy access to the parent Module
	py.LoadAttr(module, "MYLIB_VERS", &v.MadeBy)

	ret := py.Tuple{
		v,
		py.String(v.MadeBy),
	}
	return ret, nil
}

func Vacation_num_stops(self py.Object, args py.Tuple) (py.Object, error) {
	v := self.(*Vacation)
	return py.Int(len(v.Stops)), nil
}

func Vacation_get_stop(self py.Object, args py.Tuple) (py.Object, error) {
	v := self.(*Vacation)

	// Check out other convenience functions in py/util.go
	// If you would like to be a contributor for gpython, improving these or adding more is a great place to start!
	stopNum, err := py.GetInt(args[0])
	if err != nil {
		return nil, err
	}

	if stopNum < 1 || int(stopNum) > len(v.Stops) {
		return nil, py.ExceptionNewf(py.IndexError, "invalid stop index")
	}

	return py.Object(v.Stops[stopNum-1]), nil
}

func Vacation_add_stops(self py.Object, args py.Tuple) (py.Object, error) {
	v := self.(*Vacation)
	srcStops, ok := args[0].(py.Tuple)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "expected Tuple, got %T", args[0])
	}

	for _, arg := range srcStops {
		stop, ok := arg.(*VacationStop)
		if !ok {
			return nil, py.ExceptionNewf(py.TypeError, "expected Stop, got %T", arg)
		}

		v.Stops = append(v.Stops, stop)
	}

	return py.None, nil
}
