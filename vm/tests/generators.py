# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="generator 1"
def g1():
    yield 1
    yield 2
    yield 3
i = g1()
assert next(i) == 1
assert next(i) == 2
assert next(i) == 3
for _ in (1, 2):
    ok = False
    try:
        next(i)
    except StopIteration:
        ok = True
    assert ok, "StopIteration not raised"

doc="generator 2"
def g2():
    for i in range(5):
        yield i
assert tuple(g2()) == (0,1,2,3,4)

doc="generator 3"
ok = False
try:
    list(ax for x in range(10))
except NameError:
    ok = True
assert ok, "NameError not raised"

doc="yield from"
def g3():
    yield "potato"
    yield from g1()
    yield "sausage"
assert tuple(g3()) == ("potato",1,2,3,"sausage")

doc="yield into"
state = "not started"
def echo(value=None):
    """Example from python docs"""
    global state
    state = "started"
    try:
        while True:
            try:
                value = (yield value)
            except Exception as e:
                value = e
    finally:
        # clean up when close is called
        state = "finally"

assert state == "not started"
generator = echo(1)

assert state == "not started"

assert next(generator) == 1
assert state == "started"

assert next(generator) == None
assert state == "started"

assert generator.send(2) == 2
assert state == "started"

assert generator.send(3) == 3
assert state == "started"

assert next(generator) == None
assert state == "started"

# FIXME not implemented
# err = ValueError("potato")
# e = generator.throw(ValueError, "potato")
# assert isinstance(e, ValueError)
# assert state == "started"

# FIXME not implemented
# generator.close()
# assert state == "finally"


doc="finished"
