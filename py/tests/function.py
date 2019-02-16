# Copyright 2019 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="function"

def fn(p):
    "docstring"
    return p+1

assert fn(1) == 2

# FIXME this doesn't work yet
#assert fn.__doc__ == "docstring"
#fn.__doc__ = "hello"
#assert fn.__doc__ == "hello"

assert str(type(fn)) == "<class 'function'>"

fn.x = 3
assert fn.x == 3

def f2(p):
    return p+2

doc="check __code__"
fn.__code__ = f2.__code__
assert fn(1) == 3
try:
    fn.__code__ = "bad"
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="check __defaults__"
def f3(p=2):
    return p
assert f3.__defaults__ == (2,)
assert f3() == 2
f3.__defaults__ = (10,)
assert f3() == 10
assert f3.__defaults__ == (10,)
try:
    f3.__defaults__ = "bad"
except TypeError:
    pass
else:
    assert False, "TypeError not raised"
del f3.__defaults__
assert f3.__defaults__ == None or f3.__defaults__ == ()

doc="check __kwdefaults__"
def f4(*, b=2):
    return b
assert f4.__kwdefaults__ == {"b":2}
assert f4() == 2
f4.__kwdefaults__ = {"b":10}
assert f4() == 10
assert f4.__kwdefaults__ == {"b":10}
try:
    f4.__kwdefaults__ = "bad"
except TypeError:
    pass
else:
    assert False, "TypeError not raised"
del f4.__kwdefaults__
assert f4.__kwdefaults__ == None or f4.__kwdefaults__ == {}

doc="check __annotations__"
def f5(a: "potato") -> "sausage":
    pass
assert f5.__annotations__ == {'a': 'potato', 'return': 'sausage'}
f5.__annotations__ = {'a': 'potato', 'return': 'SAUSAGE'}
assert f5.__annotations__ == {'a': 'potato', 'return': 'SAUSAGE'}
try:
    f5.__annotations__ = "bad"
except TypeError:
    pass
else:
    assert False, "TypeError not raised"
del f5.__annotations__
assert f5.__annotations__ == None or f5.__annotations__ == {}

doc="check __dict__"
def f6():
    pass
assert f6.__dict__ == {}
f6.__dict__ = {'a': 'potato'}
assert f6.__dict__ == {'a': 'potato'}
try:
    f6.__dict__ = "bad"
except TypeError:
    pass
else:
    assert False, "TypeError not raised"
try:
    del f6.__dict__
except (TypeError, AttributeError):
    pass
else:
    assert False, "Error not raised"

doc="check __name__"
def f7():
    pass
assert f7.__name__ == "f7"
f7.__name__ = "new_name"
assert f7.__name__ == "new_name"
try:
    f7.__name__ = 1
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="check __qualname__"
def f8():
    pass
assert f8.__qualname__ == "f8"
f8.__qualname__ = "new_qualname"
assert f8.__qualname__ == "new_qualname"
try:
    f8.__qualname__ = 1
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

doc="finished"
