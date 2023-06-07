// Copyright 2023 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// A python Filter object
type Filter struct {
	it  Object
	fun Object
}

var FilterType = NewTypeX("filter", `filter(function or None, iterable) --> filter object

Return an iterator yielding those items of iterable for which function(item)
is true. If function is None, return the items that are true.`,
	FilterTypeNew, nil)

// Type of this object
func (f *Filter) Type() *Type {
	return FilterType
}

// FilterTypeNew
func FilterTypeNew(metatype *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var fun, seq Object
	var it Object
	err = UnpackTuple(args, kwargs, "filter", 2, 2, &fun, &seq)
	if err != nil {
		return nil, err
	}
	it, err = Iter(seq)
	if err != nil {
		return nil, err
	}
	return &Filter{it: it, fun: fun}, nil
}

func (f *Filter) M__iter__() (Object, error) {
	return f, nil
}

func (f *Filter) M__next__() (Object, error) {
	var ok bool
	for {
		item, err := Next(f.it)
		if err != nil {
			return nil, err
		}
		// if (lz->func == Py_None || lz->func == (PyObject *)&PyBool_Type)
		if _, _ok := f.fun.(Bool); _ok || f.fun == None {
			ok, err = ObjectIsTrue(item)
		} else {
			var good Object
			good, err = Call(f.fun, Tuple{item}, nil)
			if err != nil {
				return nil, err
			}
			ok, err = ObjectIsTrue(good)
		}
		if ok {
			return item, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

// Check interface is satisfied
var _ I__iter__ = (*Filter)(nil)
var _ I__next__ = (*Filter)(nil)
