# Copyright 2019 The go-python Authors.  All rights reserved.
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

doc="sort"
a = [3, 1.1, 1, 2]
s1 = list(a)
s1.sort()
assert s1 == [1, 1.1, 2, 3]
s1.sort() # sort a sorted list
assert s1 == [1, 1.1, 2, 3]
s2 = list(a)
s2.sort(reverse=True)
assert s2 == [3, 2, 1.1, 1]
s2.sort() # sort a reversed list
assert s2 == [1, 1.1, 2, 3]
s3 = list(a)
s3.sort(key=lambda l: l+1) # test lambda key
assert s3 == [1, 1.1, 2, 3]
s4 = [2.0, 2, 1, 1.0]
s4.sort(key=lambda l: 0) # test stability
assert s4 == [2.0, 2, 1, 1.0]
assert [type(t) for t in s4] == [float, int, int, float]
s4 = [2.0, 2, 1, 1.0]
s4.sort() # test stability
assert s4 == [1, 1.0, 2.0, 2]
assert [type(t) for t in s4] == [int, float, float, int]
s5 = [2.0, "abc"]
assertRaises(TypeError, lambda: s5.sort())
s5 = []
s5.sort()
assert s5 == []
s5 = [0]
s5.sort()
assert s5 == [0]
assertRaises(TypeError, lambda: s5.sort(key=1))

doc="finished"
