# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="zip"
a = [1, 2, 3, 4, 5]
b = [10, 11, 12]
c = [e for e in zip(a, b)]
assert len(c) == 3
for idx, e in enumerate(c):
    assert a[idx] == c[idx][0]
    assert b[idx] == c[idx][1]
doc="finished"