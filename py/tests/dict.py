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

doc="check __iter__"
a = {"a":"b","c":5.5}
l =  list(iter(a))
assert "a" in l
assert "c" in l
assert len(l) == 2

doc="finished"
