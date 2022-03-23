// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package string provides the implementation of the python's 'string' module.
package string

import (
	"strings"

	"github.com/go-python/gpython/py"
)

func init() {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "string",
			Doc:  module_doc,
		},
		Methods: []*py.Method{
			py.MustNewMethod("capwords", capwords, 0, capwords_doc),
		},
		Globals: py.StringDict{
			"whitespace":      whitespace,
			"ascii_lowercase": ascii_lowercase,
			"ascii_uppercase": ascii_uppercase,
			"ascii_letters":   ascii_letters,
			"digits":          digits,
			"hexdigits":       hexdigits,
			"octdigits":       octdigits,
			"punctuation":     punctuation,
			"printable":       printable,
		},
	})
}

const module_doc = `A collection of string constants.

Public module variables:

whitespace -- a string containing all ASCII whitespace
ascii_lowercase -- a string containing all ASCII lowercase letters
ascii_uppercase -- a string containing all ASCII uppercase letters
ascii_letters -- a string containing all ASCII letters
digits -- a string containing all ASCII decimal digits
hexdigits -- a string containing all ASCII hexadecimal digits
octdigits -- a string containing all ASCII octal digits
punctuation -- a string containing all ASCII punctuation characters
printable -- a string containing all ASCII characters considered printable
`

var (
	whitespace      = py.String(" \t\n\r\x0b\x0c")
	ascii_lowercase = py.String("abcdefghijklmnopqrstuvwxyz")
	ascii_uppercase = py.String("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ascii_letters   = ascii_lowercase + ascii_uppercase
	digits          = py.String("0123456789")
	hexdigits       = py.String("0123456789abcdefABCDEF")
	octdigits       = py.String("01234567")
	punctuation     = py.String("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	printable       = py.String("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r\x0b\x0c")
)

const capwords_doc = `capwords(s [,sep]) -> string

Split the argument into words using split, capitalize each
word using capitalize, and join the capitalized words using
join.  If the optional second argument sep is absent or None,
runs of whitespace characters are replaced by a single space
and leading and trailing whitespace are removed, otherwise
sep is used to split and join the words.`

func capwords(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pystr py.Object
		pysep py.Object = py.None
	)
	err := py.ParseTupleAndKeywords(args, kwargs, "s|z", []string{"s", "sep"}, &pystr, &pysep)
	if err != nil {
		return nil, err
	}

	pystr = py.String(strings.ToLower(string(pystr.(py.String))))
	pyvs, err := pystr.(py.String).Split(py.Tuple{pysep}, nil)
	if err != nil {
		return nil, err
	}

	var (
		lst   = pyvs.(*py.List).Items
		vs    = make([]string, len(lst))
		sep   = ""
		title = func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		}
	)

	switch pysep {
	case py.None:
		for i := range vs {
			v := string(lst[i].(py.String))
			vs[i] = title(strings.Trim(v, string(whitespace)))
		}
		sep = " "
	default:
		sep = string(pysep.(py.String))
		for i := range vs {
			v := string(lst[i].(py.String))
			vs[i] = title(v)
		}
	}

	return py.String(strings.Join(vs, sep)), nil
}
