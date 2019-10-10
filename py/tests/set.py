# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="__and__"
a = {1, 2, 3}
b = {2, 3, 4, 5}
c = a.__and__(b)
assert 2 in c
assert 3 in c

d = a & b
assert 2 in d
assert 3 in d

doc="__or__"
a = {1, 2, 3}
b = {2, 3, 4, 5}
c = a.__or__(b)
assert 1 in c
assert 2 in c
assert 3 in c
assert 4 in c
assert 5 in c

d = a | b
assert 1 in c
assert 2 in c
assert 3 in c
assert 4 in c
assert 5 in c

doc="__sub__"
a = {1, 2, 3}
b = {2, 3, 4, 5}
c = a.__sub__(b)
d = b.__sub__(a)
assert 1 in c
assert 4 in d
assert 5 in d

e = a - b
f = b - a
assert 1 in c
assert 4 in d
assert 5 in d

doc="__xor__"
a = {1, 2, 3}
b = {2, 3, 4, 5}
c = a.__xor__(b)
assert 1 in c
assert 4 in c
assert 5 in c

d = a ^ b
assert 1 in c

doc="__repr__"
a = {1, 2, 3}
b = a.__repr__()
assert "{" in b
assert "1" in b
assert "2" in b
assert "3" in b
assert "}" in b

doc="set"
a = set([1,2,3])
b = set("set")
c = set((4,5))
assert len(a) == 3
assert len(b) == 3
assert len(c) == 2
assert 1 in a
assert 2 in a
assert 3 in a
assert "s" in b
assert "e" in b
assert "t" in b
assert 4 in c
assert 5 in c

doc="__eq__, __ne__"
a = set([1,2,3])
assert a.__eq__(3) != True
assert a.__ne__(3) != False
assert a.__ne__(3) != True
assert a.__ne__(3) != False # This part should be changed in comparison with NotImplemented

assert a.__ne__(set()) == True
assert a.__eq__({1,2,3}) == True
assert a.__ne__({1,2,3}) == False

doc="finished"
