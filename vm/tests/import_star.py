#!/usr/bin/env python3.4

# test IMPORT_STAR

from lib import *

assert libfn() == 42
assert libvar == 43
assert libclass().method() == 44

ok = False
try:
    _libprivate
except NameError:
    ok = True
assert ok

from lib1 import *

assert lib1fn() == 42
assert lib1var == 43

doc="IMPORT_START 1"
ok = False
try:
    lib1class
except NameError:
    ok = True
assert ok

doc="IMPORT_START 2"
ok = False
try:
    _libprivate
except NameError:
    ok = True
assert ok

# End with this
doc="finished"
