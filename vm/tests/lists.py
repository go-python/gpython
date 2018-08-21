# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

L = [1,2,3]
T = (1,2,3)
assert L == [1,2,3]
assert L != [1,4,3]

doc="UNPACK_SEQUENCE"
a, b, c = L
assert a == 1
assert b == 2
assert c == 3

a, b, c = T
assert a == 1
assert b == 2
assert c == 3

a, b, c = range(3)
assert a == 0
assert b == 1
assert c == 2

ok = False
try:
    a, b = L
except ValueError:
    ok = True
assert ok

doc="UNPACK_EX"
a, *b = L
assert a == 1
assert b == [2,3]

doc="SETITEM"
LL = [1,2,3]
LL[1] = 17
assert LL == [1,17,3]

L=[1,2,3]
L[:] = [4,5,6,7]
assert L == [4,5,6,7]

L=[1,2,3]
L[3:3] = [4,5,6]
assert L == [1,2,3,4,5,6]

L=[1,2,3,4]
L[1:3] = [5,6,7]
assert L == [1,5,6,7,4]

L=[1,2,3,4]
L[:2] = [5,6,7]
assert L == [5,6,7,3,4]

L=[1,2,3,4]
L[2:] = [5,6,7]
assert L == [1,2,5,6,7]

L=[1,2,3,4]
L[::2] = [7,8]
assert L == [7, 2, 8, 4]

doc="GETITEM"
assert LL[0] == 1
assert LL[1] == 17
assert LL[2] == 3

L=[1,2,3,4,5,6]
assert L[:] == [1,2,3,4,5,6]
assert L[:3] == [1,2,3]
assert L[3:] == [4,5,6]
assert L[1:3] == [2,3]
assert L[1:5:2] == [2,4]

doc="DELITEM"
del LL[1]
assert LL == [1,3]

L=[1,2,3,4,5,6]
del L[:3]
assert L == [4,5,6]

L=[1,2,3,4,5,6]
del L[3:]
assert L == [1,2,3]

L=[1,2,3,4,5,6]
del L[1:5]
assert L == [1,6]

L=[1,2,3,4,5,6]
del L[1:5:2]
assert L == [1,3,5,6]

L=[1,2,3,4,5,6]
del L[:]
assert L == []

doc="finished"
