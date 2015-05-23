#!/usr/bin/env python3.4

# test IMPORT_FROM

from lib import libfn, libvar, libclass

assert libfn() == 42
assert libvar == 43
assert libclass().method() == 44

# End with this
finished = True
