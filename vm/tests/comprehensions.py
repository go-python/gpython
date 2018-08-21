# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="List comprehensions"
A = [1,2,3,4]
B = [ 2*a for a in A ]
assert tuple(B) == tuple([2,4,6,8])
B = [ 2*a for a in A if a != 2]
assert tuple(B) == tuple([2,6,8])

# FIXME - exitYield not working?
doc="Generator expressions"
A = (1,2,3,4)
B = ( 2*a for a in A )
assert tuple(B) == (2,4,6,8)
B = [ 2*a for a in A if a != 2]
assert tuple(B) == (2,6,8)

doc="Set comprehensions"
A = {1,2,3,4}
B = { 2*a for a in A }
assert B == {2,4,6,8}
B = { 2*a for a in A if a != 2}
assert B == {2,6,8}

doc="Dict comprehensions"
A = {"a":1, "b":2, "c":3}
B = { k:k for k in ("a","b","c") }
assert B["b"] == "b"

doc="finished"
