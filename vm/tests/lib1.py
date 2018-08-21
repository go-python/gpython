# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Some targets to be imported

__all__ = [
    "lib1fn",
    "lib1var",
]

def lib1fn():
    return 42

lib1var = 43

class lib1class:
    def method(self):
        return 44

_lib1private = 45
