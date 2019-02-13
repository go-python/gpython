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
doc="finished"