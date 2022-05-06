# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import glob

def norm(vs):
    if len(vs) == 0:
        return vs
    if type(vs[0]) == type(""):
        return normStr(vs)
    return normBytes(vs)

def normStr(vs):
    from os import sep
    x = []
    for v in vs:
        x.append(v.replace('/', sep))
    return x

def normBytes(vs):
    from os import sep
    x = []
    for v in vs:
        x.append(v.replace(b'/', bytes(sep, encoding="utf-8")))
    return x

def assertEqual(x, y):
    xx = norm(x)
    yy = norm(y)
    assert xx == yy, "got: %s, want: %s" % (repr(x), repr(y))


## test strings
assertEqual(glob.glob('*'), ["glob.go", "glob_test.go", "testdata"])
assertEqual(glob.glob('*test*'), ["glob_test.go", "testdata"])
assertEqual(glob.glob('*/test*'), ["testdata/test.py", "testdata/test_golden.txt"])
assertEqual(glob.glob('*/test*_*'), ["testdata/test_golden.txt"])
assertEqual(glob.glob('*/t??t*_*'), ["testdata/test_golden.txt"])
assertEqual(glob.glob('*/t[e]?t*_*'), ["testdata/test_golden.txt"])
assertEqual(glob.glob('*/t[oe]?t*_*'), ["testdata/test_golden.txt"])
assertEqual(glob.glob('*/t[o]?t*_*'), [])

## FIXME(sbinet)
## assertEqual(glob.glob('*/t[!o]?t*_*'), ["testdata/test_golden.txt"])

## test bytes
assertEqual(glob.glob(b'*'), [b"glob.go", b"glob_test.go", b"testdata"])
assertEqual(glob.glob(b'*test*'), [b"glob_test.go", b"testdata"])
assertEqual(glob.glob(b'*/test*'), [b"testdata/test.py", b"testdata/test_golden.txt"])
assertEqual(glob.glob(b'*/test*_*'), [b"testdata/test_golden.txt"])
assertEqual(glob.glob(b'*/t??t*_*'), [b"testdata/test_golden.txt"])
assertEqual(glob.glob(b'*/t[e]?t*_*'), [b"testdata/test_golden.txt"])
assertEqual(glob.glob(b'*/t[oe]?t*_*'), [b"testdata/test_golden.txt"])
assertEqual(glob.glob(b'*/t[o]?t*_*'), [])

## FIXME(sbinet)
## assertEqual(glob.glob(b'*/t[!o]?t*_*'), [b"testdata/test_golden.txt"])
