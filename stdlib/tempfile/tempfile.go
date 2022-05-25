// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tempfile provides the implementation of the python's 'tempfile' module.
package tempfile

import (
	"fmt"
	"os"

	"github.com/go-python/gpython/py"
)

var (
	gblTempDir py.Object = py.None
)

const tempfile_doc = `Temporary files.

This module provides generic, low- and high-level interfaces for
creating temporary files and directories.  All of the interfaces
provided by this module can be used without fear of race conditions
except for 'mktemp'.  'mktemp' is subject to race conditions and
should not be used; it is provided for backward compatibility only.

The default path names are returned as str.  If you supply bytes as
input, all return values will be in bytes.  Ex:

    >>> tempfile.mkstemp()
    (4, '/tmp/tmptpu9nin8')
    >>> tempfile.mkdtemp(suffix=b'')
    b'/tmp/tmppbi8f0hy'

This module also provides some data items to the user:

  TMP_MAX  - maximum number of names that will be tried before
             giving up.
  tempdir  - If this is set to a string before the first use of
             any routine from this module, it will be considered as
			 another candidate location to store temporary files.`

func init() {
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "tempfile",
			Doc:  tempfile_doc,
		},
		Methods: []*py.Method{
			py.MustNewMethod("gettempdir", gettempdir, 0, gettempdir_doc),
			py.MustNewMethod("gettempdirb", gettempdirb, 0, gettempdirb_doc),
			py.MustNewMethod("mkdtemp", mkdtemp, 0, mkdtemp_doc),
			py.MustNewMethod("mkstemp", mkstemp, 0, mkstemp_doc),
		},
		Globals: py.StringDict{
			"tempdir": gblTempDir,
		},
	})
}

const gettempdir_doc = `Returns tempfile.tempdir as str.`

func gettempdir(self py.Object) (py.Object, error) {
	// FIXME(sbinet): lock access to glbTempDir?
	if gblTempDir != py.None {
		switch dir := gblTempDir.(type) {
		case py.String:
			return dir, nil
		case py.Bytes:
			return py.String(dir), nil
		default:
			return nil, py.ExceptionNewf(py.TypeError, "expected str, bytes or os.PathLike object, not %s", dir.Type().Name)
		}
	}
	return py.String(os.TempDir()), nil
}

const gettempdirb_doc = `Returns tempfile.tempdir as bytes.`

func gettempdirb(self py.Object) (py.Object, error) {
	// FIXME(sbinet): lock access to glbTempDir?
	if gblTempDir != py.None {
		switch dir := gblTempDir.(type) {
		case py.String:
			return py.Bytes(dir), nil
		case py.Bytes:
			return dir, nil
		default:
			return nil, py.ExceptionNewf(py.TypeError, "expected str, bytes or os.PathLike object, not %s", dir.Type().Name)
		}
	}
	return py.Bytes(os.TempDir()), nil
}

const mkdtemp_doc = `mkdtemp(suffix=None, prefix=None, dir=None)
    User-callable function to create and return a unique temporary
    directory.  The return value is the pathname of the directory.
    
    Arguments are as for mkstemp, except that the 'text' argument is
    not accepted.
    
    The directory is readable, writable, and searchable only by the
    creating user.
    
    Caller is responsible for deleting the directory when done with it.`

func mkdtemp(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pysuffix py.Object = py.None
		pyprefix py.Object = py.None
		pydir    py.Object = py.None
	)
	err := py.ParseTupleAndKeywords(args, kwargs,
		"|z#z#z#:mkdtemp",
		[]string{"suffix", "prefix", "dir"},
		&pysuffix, &pyprefix, &pydir,
	)
	if err != nil {
		return nil, err
	}

	str := func(v py.Object, typ *uint8) string {
		switch v := v.(type) {
		case py.Bytes:
			*typ = 2
			return string(v)
		case py.String:
			*typ = 1
			return string(v)
		case py.NoneType:
			*typ = 0
			return ""
		default:
			panic(fmt.Errorf("tempfile: invalid type %T (v=%+v)", v, v))
		}
	}

	var (
		t1, t2, t3 uint8

		suffix  = str(pysuffix, &t1)
		prefix  = str(pyprefix, &t2)
		dir     = str(pydir, &t3)
		pattern = prefix + "*" + suffix
	)

	cmp := func(t1, t2 uint8) bool {
		if t1 > 0 && t2 > 0 {
			return t1 == t2
		}
		return true
	}

	if !cmp(t1, t2) || !cmp(t1, t3) || !cmp(t2, t3) {
		return nil, py.ExceptionNewf(py.TypeError, "Can't mix bytes and non-bytes in path components")
	}

	tmp, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return nil, err
	}

	typ := t1
	if typ == 0 {
		typ = t2
	}
	if typ == 0 {
		typ = t3
	}

	switch typ {
	case 2:
		return py.Bytes(tmp), nil
	default:
		return py.String(tmp), nil
	}
}

const mkstemp_doc = `mkstemp(suffix=None, prefix=None, dir=None, text=False)

User-callable function to create and return a unique temporary
file.  The return value is a pair (fd, name) where fd is the
file descriptor returned by os.open, and name is the filename.

If 'suffix' is not None, the file name will end with that suffix,
otherwise there will be no suffix.

If 'prefix' is not None, the file name will begin with that prefix,
otherwise a default prefix is used.

If 'dir' is not None, the file will be created in that directory,
otherwise a default directory is used.

If 'text' is specified and true, the file is opened in text
mode.  Else (the default) the file is opened in binary mode.

If any of 'suffix', 'prefix' and 'dir' are not None, they must be the
same type.  If they are bytes, the returned name will be bytes; str
otherwise.

The file is readable and writable only by the creating user ID.
If the operating system uses permission bits to indicate whether a
file is executable, the file is executable by no one. The file
descriptor is not inherited by children of this process.

Caller is responsible for deleting the file when done with it.`

func mkstemp(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var (
		pysuffix py.Object = py.None
		pyprefix py.Object = py.None
		pydir    py.Object = py.None
		pytext   py.Object = py.False // FIXME(sbinet): can we do something with that?
	)

	err := py.ParseTupleAndKeywords(args, kwargs,
		"|z#z#z#p:mkstemp",
		[]string{"suffix", "prefix", "dir", "text"},
		&pysuffix, &pyprefix, &pydir, &pytext,
	)
	if err != nil {
		return nil, err
	}

	str := func(v py.Object, typ *uint8) string {
		switch v := v.(type) {
		case py.Bytes:
			*typ = 2
			return string(v)
		case py.String:
			*typ = 1
			return string(v)
		case py.NoneType:
			*typ = 0
			return ""
		default:
			panic(fmt.Errorf("tempfile: invalid type %T (v=%+v)", v, v))
		}
	}

	var (
		t1, t2, t3 uint8

		suffix  = str(pysuffix, &t1)
		prefix  = str(pyprefix, &t2)
		dir     = str(pydir, &t3)
		pattern = prefix + "*" + suffix
	)

	cmp := func(t1, t2 uint8) bool {
		if t1 > 0 && t2 > 0 {
			return t1 == t2
		}
		return true
	}

	if !cmp(t1, t2) || !cmp(t1, t3) || !cmp(t2, t3) {
		return nil, py.ExceptionNewf(py.TypeError, "Can't mix bytes and non-bytes in path components")
	}

	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, err
	}

	typ := t1
	if typ == 0 {
		typ = t2
	}
	if typ == 0 {
		typ = t3
	}

	tuple := py.Tuple{py.Int(f.Fd())}
	switch typ {
	case 2:
		tuple = append(tuple, py.Bytes(f.Name()))
	default:
		tuple = append(tuple, py.String(f.Name()))
	}

	return tuple, nil
}
