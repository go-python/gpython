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

doc="range_slice"
a = range(10)
assert a[::-1][0] == 9
assert a[::-1][9] == 0
assert a[0:3][0] == 0
assert a[0:3][2] == 2
assert a[-3:10][0] == 7
assert a[-100:13][0] == 0
assert a[-100:13][9] == 9

try:
    a[0:3][3]
except IndexError:
    pass
else:
    assert False, "IndexError not raised"
try:
    a[100:13][0]
except IndexError:
    pass
else:
    assert False, "IndexError not raised"
try:
    a[0:3:0]
except ValueError:
    pass
else:
    assert False, "ValueError not raised"

doc="range_index"
class Index:
    def __index__(self):
        return 1

a = range(10)
b = Index()
assert a[b] == 1
assert a[b:10] == a[1:10]
assert a[10:b:-1] == a[10:1:-1]

class NonIntegerIndex:
    def __index__(self):
        return 1.1

a = range(10)
b = NonIntegerIndex()
try:
    a[b]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

try:
    a[b:10]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="finished"
