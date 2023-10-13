# Copyright 2023 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Imitate the calling method of unittest

def assertRaises(expecting, fn, *args, **kwargs):
    """Check the exception was raised - don't check the text"""
    try:
        fn(*args, **kwargs)
    except expecting as e:
        pass
    else:
        assert False, "%s not raised" % (expecting,)

def assertEqual(first, second, msg=None):
    if msg:
        assert first == second, "%s not equal" % (msg,)
    else:
        assert first == second

def assertIs(expr1, expr2, msg=None):
    if msg:
        assert expr1 is expr2, "%s is not None" % (msg,)
    else:
        assert expr1 is expr2

def assertIsNone(obj, msg=None):
    if msg:
        assert obj is None, "%s is not None" % (msg,)
    else:
        assert obj is None

def assertTrue(obj, msg=None):
    if msg:
        assert obj, "%s is not True" % (msg,)
    else:
        assert obj

def assertRaisesText(expecting, text, fn, *args, **kwargs):
    """Check the exception with text in is raised"""
    try:
        fn(*args, **kwargs)
    except expecting as e:
        assert text in e.args[0], "'%s' not found in '%s'" % (text, e.args[0])
    else:
        assert False, "%s not raised" % (expecting,)

def assertTypedEqual(actual, expect, msg=None):
    assertEqual(actual, expect, msg)
    def recurse(actual, expect):
        if isinstance(expect, (tuple, list)):
            for x, y in zip(actual, expect):
                recurse(x, y)
        else:
            assertIs(type(actual), type(expect))
    recurse(actual, expect)
