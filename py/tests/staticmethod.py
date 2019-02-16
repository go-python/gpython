# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="staticmethod"

class A:
    @staticmethod
    def fn(p):
        return p+1

a = A()
assert a.fn(1) == 2

a.x = 3
assert a.x == 3
    
doc="finished"
