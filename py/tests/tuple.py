# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="str"
assert str(()) == "()"
assert str((1,2,3)) == "(1, 2, 3)"
assert str((1,(2,3),4)) == "(1, (2, 3), 4)"
assert str(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"

doc="repr"
assert repr(()) == "()"
assert repr((1,2,3)) == "(1, 2, 3)"
assert repr((1,(2,3),4)) == "(1, (2, 3), 4)"
assert repr(("1",(2.5,17,()))) == "('1', (2.5, 17, ()))"

doc="mul"
a = (1, 2, 3)
assert a * 2  == (1, 2, 3, 1, 2, 3)
assert a * 0 == ()
assert a * -1 == ()

doc="finished"
