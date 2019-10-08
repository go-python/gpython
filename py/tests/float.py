# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaises

doc="float"
assert float("1.5") == 1.5
assert float(" 1.5") == 1.5
assert float(" 1.5 ") == 1.5
assert float("1.5 ") == 1.5
assert float("-1E9") == -1E9
assert float("1E400") == float("inf")
assert float(" -1E400") == float("-inf")
assertRaises(ValueError, float, "1 E200")

doc="repr"
assert repr(float("1.0")) == "1.0"
assert repr(float("1.")) == "1.0"
assert repr(float("1.1")) == "1.1"
assert repr(float("1.11")) == "1.11"
assert repr(float("-1.0")) == "-1.0"
assert repr(float("1.00101")) == "1.00101"
assert repr(float("1.00")) == "1.0"
assert repr(float("2.010")) == "2.01"

doc="str"
assert str(float("1.0")) == "1.0"
assert str(float("1.")) == "1.0"
assert str(float("1.1")) == "1.1"
assert str(float("1.11")) == "1.11"
assert str(float("-1.0")) == "-1.0"
assert str(float("1.00101")) == "1.00101"
assert str(float("1.00")) == "1.0"
assert str(float("2.010")) == "2.01"

doc="is_integer"
assert (1.0).is_integer() == True
assert (2.3).is_integer() == False

doc="finished"
