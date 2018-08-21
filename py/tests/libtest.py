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
