# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="str"
assert str(3+4j) == "(3+4j)"

doc="repr"
assert repr(3+4j) == "(3+4j)"

doc="real"
assert (3+4j).real == 3.0

doc="imag"
assert (3+4j).imag == 4.0

doc="conjugate"
assert (3+4j).conjugate() == 3-4j

doc="add"
assert (3+4j) + 2 == 5+4j
assert (3+4j) + 2j == 3+6j

doc="sub"
assert (3+4j) - 1 == 2+4j
assert (3+4j) - 1j == 3+3j

doc="mul"
assert (3+4j) * 2 == 6+8j
assert (3+4j) * 2j == -8+6j

doc="finished"
