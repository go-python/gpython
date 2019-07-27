# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="str"
assert str([]) == "[]"
assert str([1,2,3]) == "[1, 2, 3]"
assert str([1,[2,3],4]) == "[1, [2, 3], 4]"
assert str(["1",[2.5,17,[]]]) == "['1', [2.5, 17, []]]"

doc="repr"
assert repr([]) == "[]"
assert repr([1,2,3]) == "[1, 2, 3]"
assert repr([1,[2,3],4]) == "[1, [2, 3], 4]"
assert repr(["1",[2.5,17,[]]]) == "['1', [2.5, 17, []]]"

doc="enumerate"
a = [e for e in enumerate([3,4,5,6,7], 4)]
idxs = [4, 5, 6, 7, 8]
values = [3, 4, 5, 6, 7]
for idx, value in enumerate(values):
    assert idxs[idx] == a[idx][0]
    assert values[idx] == a[idx][1]

doc="append"
a = [1,2,3]
a.append(4)
assert repr(a) == "[1, 2, 3, 4]"
a = ['a', 'b', 'c']
a.extend(['d', 'e', 'f'])
assert repr(a) == "['a', 'b', 'c', 'd', 'e', 'f']"
assertRaises(TypeError, lambda: [].append())

doc="mul"
a = [1, 2, 3]
assert a * 2  == [1, 2, 3, 1, 2, 3]
assert a * 0 == []
assert a * -1 == []

doc="finished"
