// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package glob provides the implementation of the python's 'glob' module.
package glob

import (
	"path/filepath"

	"github.com/go-python/gpython/py"
)

func init() {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "glob",
			Doc:  "Filename globbing utility.",
		},
		Methods: []*py.Method{
			py.MustNewMethod("glob", glob, 0, glob_doc),
		},
	})
}

const glob_doc = `Return a list of paths matching a pathname pattern.
The pattern may contain simple shell-style wildcards a la
fnmatch. However, unlike fnmatch, filenames starting with a
dot are special cases that are not matched by '*' and '?'
patterns.`

func glob(self py.Object, args py.Tuple) (py.Object, error) {
	var (
		pypathname py.Object
	)
	err := py.ParseTuple(args, "s*:glob", &pypathname)
	if err != nil {
		return nil, err
	}

	var (
		pathname string
		cnv      func(v string) py.Object
	)
	switch n := pypathname.(type) {
	case py.String:
		pathname = string(n)
		cnv = func(v string) py.Object { return py.String(v) }
	case py.Bytes:
		pathname = string(n)
		cnv = func(v string) py.Object { return py.Bytes(v) }
	}
	matches, err := filepath.Glob(pathname)
	if err != nil {
		return nil, err
	}

	lst := py.List{Items: make([]py.Object, len(matches))}
	for i, v := range matches {
		lst.Items[i] = cnv(v)
	}

	return &lst, nil
}
