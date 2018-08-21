# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="CALL_FUNCTION"
def fn1(a, b=None):
    assert a == 1
    assert b == 2
fn1(1, 2)
fn1(1, b=2)
fn1(a=1, b=2)

ok = False
try:
    fn1(1, 2, c=4)
except TypeError as e:
    assert e.args[0] == "fn1() got an unexpected keyword argument 'c'"
    ok = True
assert ok, "TypeError not raised"

ok = False
try:
    fn1(1,*(2,3,4),b=42)
except TypeError as e:
    assert e.args[0] == "fn1() got multiple values for argument 'b'"
    ok = True
assert ok, "TypeError not raised"

doc="CALL_FUNCTION_VAR"
def fn2(a, b=None, *args):
    assert a == 1
    assert b == 2
    assert args == (3, 4)
fn2(1,2,*(3,4))
fn2(1,*(2,3,4))
fn2(*(1,2,3,4))

doc="CALL_FUNCTION_KW"
def fn3(a, b=None, **kwargs):
    assert a == 1
    assert b == 2
    assert kwargs['c'] == 3
    assert kwargs['d'] == 4
fn3(1, 2, **{'c':3, 'd':4})
fn3(1, **{'b':2, 'c':3, 'd':4})
fn3(**{'a':1, 'b':2, 'c':3, 'd':4})
    
doc="CALL_FUNCTION_VAR_KW"
def fn4(a, b=None, *args, **kwargs):
    assert a == 1
    assert b == 2
    assert args == (3, 4)
    assert kwargs['c'] == 5
    assert kwargs['d'] == 6
fn4(1, 2, *(3, 4), **{'c':5, 'd':6})

doc="finished"
