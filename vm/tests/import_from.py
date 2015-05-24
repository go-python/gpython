#!/usr/bin/env python3.4

doc="IMPORT_FROM"

from lib import libfn, libvar, libclass

assert libfn() == 42
assert libvar == 43
assert libclass().method() == 44

# End with this
doc="finished"
