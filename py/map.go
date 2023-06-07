// Copyright 2023 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// A python Map object
type Map struct {
	iters Tuple
	fun   Object
}

var MapType = NewTypeX("filter", `map(func, *iterables) --> map object

Make an iterator that computes the function using arguments from
each of the iterables.  Stops when the shortest iterable is exhausted.`,
	MapTypeNew, nil)

// Type of this object
func (m *Map) Type() *Type {
	return FilterType
}

// MapType
func MapTypeNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	numargs := len(args)
	if numargs < 2 {
		return nil, ExceptionNewf(TypeError, "map() must have at least two arguments.")
	}
	iters := make(Tuple, numargs-1)
	for i := 1; i < numargs; i++ {
		iters[i-1], err = Iter(args[i])
		if err != nil {
			return nil, err
		}
	}
	return &Map{iters: iters, fun: args[0]}, nil
}

func (m *Map) M__iter__() (Object, error) {
	return m, nil
}

func (m *Map) M__next__() (Object, error) {
	numargs := len(m.iters)
	argtuple := make(Tuple, numargs)

	for i := 0; i < numargs; i++ {
		val, err := Next(m.iters[i])
		if err != nil {
			return nil, err
		}
		argtuple[i] = val
	}
	return Call(m.fun, argtuple, nil)
}

// Check interface is satisfied
var _ I__iter__ = (*Map)(nil)
var _ I__next__ = (*Map)(nil)
