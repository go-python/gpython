# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
from libtest import assertRaises

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

doc="check get"
a = {"a":1}
assert a.get('a') == 1
assert a.get('a',100) == 1
assert a.get('b') == None
assert a.get('b',1) == 1
assert a.get('b',True) == True

doc="check items"
a = {"a":"b","c":5.5}
for k, v in a.items():
    assert k in ["a", "c"]
    if k == "a":
        assert v == "b"
    if k == "c":
        assert v == 5.5
assertRaises(TypeError, a.items, 'a')

doc="__contain__"
a = {'hello': 'world'}
assert a.__contains__('hello')
assert not a.__contains__('world')

doc="__eq__, __ne__"
a = {'a': 'b'}
assert a.__eq__(3) != True
assert a.__ne__(3) != False
assert a.__ne__(3) != True
assert a.__ne__(3) != False

assert a.__ne__({}) == True
assert a.__eq__({'a': 'b'}) == True
assert a.__ne__({'a': 'b'}) == False

doc="finished"
