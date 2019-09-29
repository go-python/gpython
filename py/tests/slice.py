# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="slice"
a = slice(10)
assert a.start == None
assert a.stop == 10
assert a.step == None

a = slice(0, 10, 1)
assert a.start == 0
assert a.stop == 10
assert a.step == 1

assert slice(1).__eq__(slice(1))
assert slice(1) != slice(2)
assert slice(1) == slice(None, 1, None)
assert slice(0, 0, 0) == slice(0, 0, 0)

assert slice(0, 0, 1) != slice(0, 0, 0)
assert slice(0, 1, 0) != slice(0, 0, 0)
assert slice(1, 0, 0) != slice(0, 0, 0)
assert slice(0).__ne__(slice(1))
assert slice(0, None, 3).__ne__(slice(0, 0, 3))

doc="finished"