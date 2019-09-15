# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

l = []
l2 = sorted(l)
assert l == l2
assert not l is l2
assert sorted([5, 2, 3, 1, 4]) == [1, 2, 3, 4, 5]
a = [5, 2, 3, 1, 4]
assert a.sort() == None
assert a == [1, 2, 3, 4, 5]
assert sorted({"1": "D", "2": "B", "3": "B", "5": "E", "4": "A"}) == ["1", "2", "3", "4", "5"]

kwargs = {"key": lambda l: l&1+l, "reverse": True}
l = list(range(10))
l.sort(**kwargs)
assert l == sorted(range(10), **kwargs) == [8, 9, 6, 4, 5, 2, 0, 1, 3, 7]

assert sorted([1, 2, 1.1], reverse=1) == [2, 1.1, 1]

try:
    sorted()
except TypeError:
    pass
else:
    assert False

try:
    sorted([], 1)
except TypeError:
    pass
else:
    assert False

try:
    sorted(1)
except TypeError:
    pass
else:
    assert False

try:
    sorted(None)
except TypeError:
    pass
else:
    assert False

try:
    sorted([1, 2], key=1)
except TypeError:
    pass
else:
    assert False