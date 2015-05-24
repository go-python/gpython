#!/usr/bin/env python3.4

doc="IMPORT_NAME"

import lib

assert lib.libfn() == 42
assert lib.libvar == 43
assert lib.libclass().method() == 44

# End with this
doc="finished"
