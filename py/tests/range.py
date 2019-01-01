# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="range"
a = range(255)
b = [e for e in a]
assert len(a) == len(b)
a = range(5, 100, 5)
b = [e for e in a]
assert len(a) == len(b)
a = range(100 ,0, 1)
b = [e for e in a]
assert len(a) == len(b)

a = range(100, 0, -1)
b = [e for e in a]
assert len(a) == 100
assert len(b) == 100

doc="range_get_item"
a = range(3)
assert a[2] == 2
assert a[1] == 1
assert a[0] == 0
assert a[-1] == 2
assert a[-2] == 1
assert a[-3] == 0

b = range(0, 10, 2)
assert b[4] == 8
assert b[3] == 6
assert b[2] == 4
assert b[1] == 2
assert b[0] == 0
assert b[-4] == 2
assert b[-3] == 4
assert b[-2] == 6
assert b[-1] == 8

doc="finished"
