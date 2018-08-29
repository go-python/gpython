# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc = "open"
assertRaises(FileNotFoundError, open, "not-existent.file")

assertRaises(IsADirectoryError, open, ".")

f = open(__file__)
assert f is not None

doc = "read"
b = f.read(12)
assert b == '# Copyright '

b = f.read(4)
assert b == '2018'

b = f.read()
assert b != ''

b = f.read()
assert b == ''

doc = "write"
assertRaises(TypeError, f.write, 42)

# assertRaises(io.UnsupportedOperation, f.write, 'hello')

import sys
n = sys.stdout.write('hello')
assert n == 5

doc = "close"
f.close()

# closing a closed file should not throw an error
f.close()

doc = "finished"
