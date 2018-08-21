# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="IMPORT_FROM"

from lib import libfn, libvar, libclass

assert libfn() == 42
assert libvar == 43
assert libclass().method() == 44

doc="finished"
