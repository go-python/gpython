// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sequence operations

package py

// Converts a sequence object v into a Tuple
func SequenceTuple(v Object) (Tuple, error) {
	switch x := v.(type) {
	case Tuple:
		return x, nil
	case *List:
		return Tuple(x.Items).Copy(), nil
	default:
		t := Tuple{}
		err := Iterate(v, func(item Object) bool {
			t = append(t, item)
			return false
		})
		if err != nil {
			return nil, err
		}
		return t, nil
	}
}

// Converts a sequence object v into a List
func SequenceList(v Object) (*List, error) {
	switch x := v.(type) {
	case Tuple:
		return NewListFromItems(x), nil
	case *List:
		return x.Copy(), nil
	default:
		l := NewList()
		err := l.ExtendSequence(v)
		if err != nil {
			return nil, err
		}
		return l, nil
	}
}

// Call __next__ for the python object
//
// Returns the next object
//
// err == StopIteration or subclass when finished
func Next(self Object) (obj Object, err error) {
	if I, ok := self.(I__next__); ok {
		return I.M__next__()
	} else if obj, ok, err = TypeCall0(self, "__next__"); ok {
		return obj, err
	}
	return nil, ExceptionNewf(TypeError, "'%s' object is not iterable", self.Type().Name)
}

// Create an iterator from obj and iterate the iterator until finished
// calling the function passed in on each object.  The iteration is
// finished if the function returns true
func Iterate(obj Object, fn func(Object) bool) error {
	// Some easy cases
	switch x := obj.(type) {
	case Tuple:
		for _, item := range x {
			if fn(item) {
				break
			}
		}
	case *List:
		for _, item := range x.Items {
			if fn(item) {
				break
			}
		}
	case String:
		for _, item := range x {
			if fn(String(item)) {
				break
			}
		}
	case Bytes:
		for _, item := range x {
			if fn(Int(item)) {
				break
			}
		}
	default:
		iterator, err := Iter(obj)
		if err != nil {
			return err
		}
		for {
			item, err := Next(iterator)
			if err == StopIteration {
				break
			}
			if err != nil {
				return err
			}
			if fn(item) {
				break
			}
		}
	}
	return nil
}

// Call send for the python object
func Send(self, value Object) (Object, error) {
	if I, ok := self.(I_send); ok {
		return I.Send(value)
	} else if res, ok, err := TypeCall1(self, "send", value); ok {
		return res, err
	}
	return nil, ExceptionNewf(TypeError, "'%s' object doesn't have send method", self.Type().Name)
}

// SequenceContains returns True if obj is in seq
func SequenceContains(seq, obj Object) (found bool, err error) {
	if I, ok := seq.(I__contains__); ok {
		result, err := I.M__contains__(obj)
		if err != nil {
			return false, err
		}
		return result == True, nil
	} else if result, ok, err := TypeCall1(seq, "__contains__", obj); ok {
		if err != nil {
			return false, err
		}
		return result == True, nil
	}
	var loopErr error
	err = Iterate(seq, func(item Object) bool {
		var eq Object
		eq, loopErr = Eq(item, obj)
		if loopErr != nil {
			return true
		}
		if eq == True {
			found = true
			return true
		}
		return false
	})
	if err == nil {
		err = loopErr
	}
	return found, err
}
