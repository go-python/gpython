# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="iter"
cnt = 0
def f():
    global cnt
    cnt += 1
    return cnt

l = list(iter(f,20))
assert len(l) == 19
for idx, v in enumerate(l):
    assert idx + 1 == v

words1 = ['g', 'p', 'y', 't', 'h', 'o', 'n']
words2 = list(iter(words1))
for w1, w2 in zip(words1, words2):
    assert w1 == w2

class SequenceClass:
    def __init__(self, n):
        self.n = n
    def __getitem__(self, i):
        if 0 <= i < self.n:
            return i
        else:
            raise IndexError

assert list(iter(SequenceClass(5))) == [0, 1, 2, 3, 4]

doc="finished"