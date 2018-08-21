# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

class C:
    attr1 = 42
    def __init__(self):
        self.attr2 = 43
        self.attr3 = 44
c = C()

doc="Test LOAD_ATTR"
assert c.attr1 == 42
assert C.attr1 == 42
assert c.attr2 == 43
assert c.attr3 == 44

doc="Test DELETE_ATTR"
del c.attr3

ok = False
try:
    c.attr3
except AttributeError:
    ok = True
assert ok

doc="Test STORE_ATTR"
c.attr1 = 100
c.attr2 = 101
c.attr3 = 102
assert c.attr1 == 100
assert C.attr1 == 42
assert c.attr2 == 101
assert c.attr3 == 102

doc="finished"
