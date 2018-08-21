# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="str"
assert str({}) == "{}"
a = str({"a":"b","c":5.5})
assert a == "{'a': 'b', 'c': 5.5}" or a == "{'c': 5.5, 'a': 'b'}"

doc="repr"
assert repr({}) == "{}"
a = repr({"a":"b","c":5.5})
assert a == "{'a': 'b', 'c': 5.5}" or a == "{'c': 5.5, 'a': 'b'}"

doc="finished"
