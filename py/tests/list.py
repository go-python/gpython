# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="str"
assert str([]) == "[]"
assert str([1,2,3]) == "[1, 2, 3]"
assert str([1,[2,3],4]) == "[1, [2, 3], 4]"
assert str(["1",[2.5,17,[]]]) == "['1', [2.5, 17, []]]"
assert str([1, 1.0]) == "[1, 1.0]"

doc="repr"
assert repr([]) == "[]"
assert repr([1,2,3]) == "[1, 2, 3]"
assert repr([1,[2,3],4]) == "[1, [2, 3], 4]"
assert repr(["1",[2.5,17,[]]]) == "['1', [2.5, 17, []]]"
assert repr([1, 1.0]) == "[1, 1.0]"

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
# [].sort
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
s5 = [0, 1]
# Sorting a list of len >= 2 with uncallable key must fail on all Python implementations.
assertRaises(TypeError, lambda: s5.sort(key=1))

# list.sort([])
a = [3, 1.1, 1, 2]
s1 = list(a)
assert list.sort(s1) is None
assert s1 == [1, 1.1, 2, 3]
assert list.sort(s1) is None # sort a sorted list
assert s1 == [1, 1.1, 2, 3]
s2 = list(a)
list.sort(s2, reverse=True)
assert s2 == [3, 2, 1.1, 1]
list.sort(s2) # sort a reversed list
assert s2 == [1, 1.1, 2, 3]
s3 = list(a)
list.sort(s3, key=lambda l: l+1) # test lambda key
assert s3 == [1, 1.1, 2, 3]
s4 = [2.0, 2, 1, 1.0]
list.sort(s4, key=lambda l: 0) # test stability
assert s4 == [2.0, 2, 1, 1.0]
assert [type(t) for t in s4] == [float, int, int, float]
s4 = [2.0, 2, 1, 1.0]
list.sort(s4) # test stability
assert s4 == [1, 1.0, 2.0, 2]
assert [type(t) for t in s4] == [int, float, float, int]
s5 = [2.0, "abc"]
assertRaises(TypeError, lambda: list.sort(s5))
s5 = []
list.sort(s5)
assert s5 == []
s5 = [0]
list.sort(s5)
assert s5 == [0]
s5 = [0, 1]
# Sorting a list of len >= 2 with uncallable key must fail on all Python implementations.
assertRaises(TypeError, lambda: list.sort(s5, key=1))
assertRaises(TypeError, lambda: list.sort(1))

class Index:
    def __index__(self):
        return 1

a = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
b = Index()
assert a[b] == 1
assert a[b:10] == a[1:10]
assert a[10:b:-1] == a[10:1:-1]

class NonIntegerIndex:
    def __index__(self):
        return 1.1

a = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
b = NonIntegerIndex()
try:
    a[b]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

try:
    a[b:10]
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="finished"
