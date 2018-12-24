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

doc="finished"
