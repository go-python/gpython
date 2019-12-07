// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dict and StringDict type
//
// The idea is that most dicts just have strings for keys so we use
// the simpler StringDict and promote it into a Dict when necessary

package py

import "bytes"

const dictDoc = `dict() -> new empty dictionary
dict(mapping) -> new dictionary initialized from a mapping object's
    (key, value) pairs
dict(iterable) -> new dictionary initialized as if via:
    d = {}
    for k, v in iterable:
        d[k] = v
dict(**kwargs) -> new dictionary initialized with the name=value pairs
    in the keyword argument list.  For example:  dict(one=1, two=2)`

var (
	StringDictType = NewType("dict", dictDoc)
	DictType       = NewType("dict", dictDoc)
	expectingDict  = ExceptionNewf(TypeError, "a dict is required")
)

func init() {
	StringDictType.Dict["items"] = MustNewMethod("items", func(self Object, args Tuple) (Object, error) {
		err := UnpackTuple(args, nil, "items", 0, 0)
		if err != nil {
			return nil, err
		}
		sMap := self.(StringDict)
		o := make([]Object, 0, len(sMap))
		for k, v := range sMap {
			o = append(o, Tuple{String(k), v})
		}
		return NewIterator(o), nil
	}, 0, "items() -> list of D's (key, value) pairs, as 2-tuples")

	StringDictType.Dict["get"] = MustNewMethod("get", func(self Object, args Tuple) (Object, error) {
		var length = len(args)
		switch {
		case length == 0:
			return nil, ExceptionNewf(TypeError, "%s expected at least 1 arguments, got %d", "items()", length)
		case length > 2:
			return nil, ExceptionNewf(TypeError, "%s expected at most 2 arguments, got %d", "items()", length)
		}
		sMap := self.(StringDict)
		if str, ok := args[0].(String); ok {
			if res, ok := sMap[string(str)]; ok {
				return res, nil
			}

			switch length {
			case 2:
				return args[1], nil
			default:
				return None, nil
			}
		}
		return nil, ExceptionNewf(KeyError, "%v", args[0])
	}, 0, "gets(key, default) -> If there is a val corresponding to key, return val, otherwise default")
	
	DictType.Dict["items"] = MustNewMethod("items", func(self Object, args Tuple) (Object, error) {
		dict := self.(Dict)
		o := make([]Object, 0, len(dict.keys))
		for _, item := range dict.items {
			if item[0] != nil {
				o = append(o, Tuple{item[0], item[1]})
			}
		}
		return NewIterator(o), nil
	}, 0, "items() -> list of D's (key, value) pairs, as 2-tuples")

	DictType.Dict["get"] = MustNewMethod("get", func(self Object, args Tuple) (Object, error) {
		var length = len(args)
		switch {
		case length == 0:
			return nil, ExceptionNewf(TypeError, "%s expected at least 1 arguments, got %d", "items()", length)
		case length > 2:
			return nil, ExceptionNewf(TypeError, "%s expected at most 2 arguments, got %d", "items()", length)
		}
		dict := self.(Dict)
		if res, ok := dict.keys[args[0]]; ok {
			return dict.items[res][1], nil
		}

		if length == 2 {
			return args[1], nil
		}
		return None, nil
	}, 0, "gets(key, default) -> If there is a val corresponding to key, return val, otherwise default")
}

// String to object dictionary
//
// Used for variables etc where the keys can only be strings
type StringDict map[string]Object

type Dict struct {
	keys  map[Object]int
	items [][2]Object // key-value pair
}

// Type of this StringDict object
func (o StringDict) Type() *Type {
	return StringDictType
}

// Type of this Dict object
func (o Dict) Type() *Type {
	return DictType
}

// Make a new dictionary
func NewStringDict() StringDict {
	return make(StringDict)
}

// Make a new dictionary
func NewDict() *Dict {
	return &Dict{}
}

// Make a new dictionary with reservation for n entries
func NewStringDictSized(n int) StringDict {
	return make(StringDict, n)
}

// Make a new dictionary with reservation for n entries
func NewDictSized(n int) *Dict {
	return &Dict{
		keys: make(map[Object]int, n),
		items: make([][2]Object, n),
	}
}

// Checks that obj is exactly a dictionary and returns an error if not
func DictCheckExact(obj Object) (StringDict, error) {
	dict, ok := obj.(StringDict)
	if !ok {
		return nil, expectingDict
	}
	return dict, nil
}

// Checks that obj is exactly a dictionary and returns an error if not
func _DictCheckExact(obj Object) (*Dict, error) {
	dict, ok := obj.(Dict)
	if !ok {
		return nil, expectingDict
	}
	return &dict, nil
}

// Checks that obj is exactly a dictionary and returns an error if not
func DictCheck(obj Object) (StringDict, error) {
	// FIXME should be checking subclasses
	return DictCheckExact(obj)
}

// Checks that obj is exactly a dictionary and returns an error if not
func _DictCheck(obj Object) (*Dict, error) {
	// FIXME should be checking subclasses
	return _DictCheckExact(obj)
}

// Copy a dictionary
func (d StringDict) Copy() StringDict {
	e := make(StringDict, len(d))
	for k, v := range d {
		e[k] = v
	}
	return e
}

// Copy a dictionary
func (d Dict) Copy() *Dict {
	e := NewDictSized(len(d.keys))
	for k, v := range d.keys {
		e.keys[k] = v
		e.items[v][0] = d.items[v][0]
		e.items[v][1] = d.items[v][1]
	}
	return e
}

func (a StringDict) M__str__() (Object, error) {
	return a.M__repr__()
}


func (a Dict) M__str__() (Object, error) {
	return a.M__repr__()
}

func (a StringDict) M__repr__() (Object, error) {
	var out bytes.Buffer
	out.WriteRune('{')
	spacer := false
	for key, value := range a {
		if spacer {
			out.WriteString(", ")
		}
		keyStr, err := ReprAsString(String(key))
		if err != nil {
			return nil, err
		}
		valueStr, err := ReprAsString(value)
		if err != nil {
			return nil, err
		}
		out.WriteString(keyStr)
		out.WriteString(": ")
		out.WriteString(valueStr)
		spacer = true
	}
	out.WriteRune('}')
	return String(out.String()), nil
}

func (a Dict) M__repr__() (Object, error) {
	var out bytes.Buffer
	out.WriteRune('{')
	spacer := false
	for _, item := range a.items {
		if item[0] != nil {
			if spacer {
				out.WriteString(", ")
			}
			keyStr, err := ReprAsString(item[0])
			if err != nil {
				return nil, err
			}
			valueStr, err := ReprAsString(item[1])
			if err != nil {
				return nil, err
			}
			out.WriteString(keyStr)
			out.WriteString(": ")
			out.WriteString(valueStr)
			spacer = true
		}
	}
	out.WriteRune('}')
	return String(out.String()), nil
}

// Returns a list of keys from the dict
func (d StringDict) M__iter__() (Object, error) {
	o := make([]Object, 0, len(d))
	for k := range d {
		o = append(o, String(k))
	}
	return NewIterator(o), nil
}

// Returns a list of keys from the dict
func (d Dict) M__iter__() (Object, error) {
	o := make([]Object, 0, len(d.keys))
	for _, item := range d.items {
		if item[0] != nil {
			o = append(o, item[0])
		}
	}
	return NewIterator(o), nil
}

func (d StringDict) M__getitem__(key Object) (Object, error) {
	str, ok := key.(String)
	if ok {
		res, ok := d[string(str)]
		if ok {
			return res, nil
		}
	}
	return nil, ExceptionNewf(KeyError, "%v", key)
}

func (d Dict) M__getitem__(key Object) (Object, error) {
	// FIXME should be checking hash of Object
	res, ok := d.keys[key]
	if ok {
		return d.items[res][1], nil
	}
	return nil, ExceptionNewf(KeyError, "%v", key)
}

func (d StringDict) M__setitem__(key, value Object) (Object, error) {
	str, ok := key.(String)
	if !ok {
		return nil, ExceptionNewf(KeyError, "FIXME can only have string keys!: %v", key)
	}
	d[string(str)] = value
	return None, nil
}

func (d Dict) M__setitem__(key, value Object) (Object, error) {
	// FIXME should be checking hash of Object
	d.keys[key] = len(d.items)
	d.items = append(d.items, [2]Object{key, value})
	return None, nil
}

func (a StringDict) M__eq__(other Object) (Object, error) {
	b, ok := other.(StringDict)
	if !ok {
		return NotImplemented, nil
	}
	if len(a) != len(b) {
		return False, nil
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			return False, nil
		}
		res, err := Eq(av, bv)
		if err != nil {
			return nil, err
		}
		if res == False {
			return False, nil
		}
	}
	return True, nil
}

func (a Dict) M__eq__(other Object) (Object, error) {
	b, ok := other.(Dict)
	if !ok {
		return NotImplemented, nil
	}
	if len(a.keys) != len(b.keys) {
		return False, nil
	}
	for k, ai := range a.keys {
		// FIXME should be checking hash of Object
		bi, ok := b.keys[k]
		if !ok || len(a.keys) < ai || len(b.keys) < bi {
			return False, nil
		}
		aitem := a.items[ai]
		bitem := b.items[bi]

		res, err := Eq(aitem[0], bitem[0])
		if err != nil {
			return nil, err
		}
		if res == False {
			return False, nil
		}
		res, err = Eq(aitem[1], bitem[1])
		if err != nil {
			return nil, err
		}
		if res == False {
			return False, nil
		}
	}
	return True, nil
}

func (a StringDict) M__ne__(other Object) (Object, error) {
	res, err := a.M__eq__(other)
	if err != nil {
		return nil, err
	}
	if res == NotImplemented {
		return res, nil
	}
	if res == True {
		return False, nil
	}
	return True, nil
}

func (a Dict) M__ne__(other Object) (Object, error) {
	return notEq(a.M__eq__(other))
}

func (a StringDict) M__contains__(other Object) (Object, error) {
	key, ok := other.(String)
	if !ok {
		return nil, ExceptionNewf(KeyError, "FIXME can only have string keys!: %v", key)
	}

	if _, ok := a[string(key)]; ok {
		return True, nil
	}
	return False, nil
}

func (a Dict) M__contains__(other Object) (Object, error) {
	// FIXME should be checking hash of Object
	if i, ok := a.keys[other]; ok {
		return Eq(other, a.items[i][0])
	}
	return False, nil
}
