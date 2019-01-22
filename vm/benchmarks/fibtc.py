# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Benchmark adapted from https://github.com/d5/tengobench/
doc="fib tail call recursion test"
def fib(n, a, b):
    if n == 0:
        return a
    elif n == 1:
        return b
    return fib(n-1, b, a+b)

fib(35, 0, 1)
doc="finished"
