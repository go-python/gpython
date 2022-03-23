# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""
Simple test harness
"""

def assertRaises(expecting, fn, *args, **kwargs):
    """Check the exception was raised - don't check the text"""
    try:
        fn(*args, **kwargs)
    except expecting as e:
        pass
    else:
        assert False, "%s not raised" % (expecting,)

def assertRaisesText(expecting, text, fn, *args, **kwargs):
    """Check the exception with text in is raised"""
    try:
        fn(*args, **kwargs)
    except expecting as e:
        assert text in e.args[0], "'%s' not found in '%s'" % (text, e.args[0])
    else:
        assert False, "%s not raised" % (expecting,)

def assertTrue(x):
    """assert x is True"""
    assert x

def assertFalse(x):
    """assert x is False"""
    assert not x

def assertEqual(x, y):
    """assert x == y"""
    assert x == y

def assertAlmostEqual(x, y, places=7):
    """assert x == y to places"""
    assert round(abs(y-x), places) == 0

def fail(x):
    """Fails with error message"""
    assert False, x
