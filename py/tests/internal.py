# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

def fn(x):
    return x

doc="check internal bound methods"
assert (1).__str__() == "1"
assert (1).__add__(2) == 3
assert fn.__call__(4) == 4
assert fn.__get__(fn, None)()(1) == 1
assertRaises(TypeError, fn.__get__, fn, None, None)
# These tests don't work on python3.4
# assert Exception().__getattr__("a") is not None # check doesn't explode only
# assertRaises(TypeError, Exception().__getattr__, "a", "b")
# assertRaises(ValueError, Exception().__getattr__, 42)

doc="finished"

