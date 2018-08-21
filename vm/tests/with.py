# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

class Context():
    def __init__(self):
        self.state = "__init__"
    def __enter__(self):
        self.state = "__enter__"
        return 42
    def __exit__(self, type, value, traceback):
        self.state = "__exit__"

doc="with"
c = Context()
assert c.state == "__init__"
with c:
    assert c.state == "__enter__"

assert c.state == "__exit__"

doc="with as"
c = Context()
assert c.state == "__init__"
with c as a:
    assert a == 42
    assert c.state == "__enter__"

assert c.state == "__exit__"

doc="with exception"
c = Context()
ok = False
try:
    assert c.state == "__init__"
    with c:
        assert c.state == "__enter__"
        raise ValueError("potato")
except ValueError:
    ok = True
assert c.state == "__exit__"
assert ok, "ValueError not raised"

class SilencedContext():
    def __init__(self):
        self.state = "__init__"
    def __enter__(self):
        self.state = "__enter__"
    def __exit__(self, type, value, traceback):
        """Return True to silence the error"""
        self.type = type
        self.value = value
        self.traceback = traceback
        self.state = "__exit__"
        return True

doc="with silenced error"
c = SilencedContext()
assert c.state == "__init__"
with c:
    assert c.state == "__enter__"
    raise ValueError("potato")
assert c.state == "__exit__"
assert c.type == ValueError
assert c.value is not None
assert c.traceback is not None

doc="with silenced error no error"
c = SilencedContext()
assert c.state == "__init__"
with c:
    assert c.state == "__enter__"
assert c.state == "__exit__"
assert c.type is None 
assert c.value is None
assert c.traceback is None

doc="with in loop: break"
c = Context()
assert c.state == "__init__"
while True:
    with c:
        assert c.state == "__enter__"
        break
assert c.state == "__exit__"

doc="with in loop: continue"
c = Context()
assert c.state == "__init__"
first = True
while True:
    if not first:
        break
    with c:
        assert c.state == "__enter__"
        first = False
        continue
assert c.state == "__exit__"

doc="return in with"
c = Context()
def return_in_with():
    assert c.state == "__init__"
    first = True
    with c:
        assert c.state == "__enter__"
        first = False
        return "potato"
assert return_in_with() == "potato"
assert c.state == "__exit__"

doc="finished"
