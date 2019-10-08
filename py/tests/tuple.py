# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="str"
assert str(()) == "()"
assert str((1,2,3)) == "(1, 2, 3)"
assert str((1,(2,3),4)) == "(1, (2, 3), 4)"
assert str(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"
assert str((1, 1.0)) == "(1, 1.0)"

doc="repr"
assert repr(()) == "()"
assert repr((1,2,3)) == "(1, 2, 3)"
assert repr((1,(2,3),4)) == "(1, (2, 3), 4)"
assert repr(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"
assert repr((1, 1.0)) == "(1, 1.0)"

doc="mul"
a = (1, 2, 3)
assert a * 2  == (1, 2, 3, 1, 2, 3)
assert a * 0 == ()
assert a * -1 == ()

class Index:
    def __index__(self):
        return 1

a = (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
b = Index()
assert a[b] == 1
assert a[b:10] == a[1:10]
assert a[10:b:-1] == a[10:1:-1]

class NonIntegerIndex:
    def __index__(self):
        return 1.1

a = (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
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
