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

doc="range_eq"
assert range(10) == range(0, 10)
assert not range(10) == 3
assert range(20) != range(10)
assert range(100, 200, 1) == range(100, 200)
assert range(0, 10, 3) == range(0, 12, 3)
assert range(2000, 100) == range(3, 1)
assert range(0, 10, -3) == range(0, 12, -3)
assert not range(0, 20, 2) == range(0, 20, 4)
try:
    range('3', 10) == range(2)
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="range_ne"
assert range(10, 0, -3) != range(12, 0, -3)
assert range(10) != 3
assert not range(100, 200, 1) != range(100, 200)
assert range(0, 10) != range(0, 12)
assert range(0, 10) != range(0, 10, 2)
assert range(0, 20, 2) != range(0, 21, 2)
assert range(0, 20, 2) != range(0, 20, 4)
assert not range(0, 20, 3) != range(0, 20, 3)
try:
    range('3', 10) != range(2)
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="range_str"
assert str(range(10)) == 'range(0, 10)'
assert str(range(10, 0, 3)) == 'range(10, 0, 3)'
assert str(range(0, 3)) == 'range(0, 3)'
assert str(range(10, 3, -2)) == 'range(10, 3, -2)'

doc="finished"
