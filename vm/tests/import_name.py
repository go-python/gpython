# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

doc="IMPORT_NAME"

import lib

assert lib.libfn() == 42
assert lib.libvar == 43
assert lib.libclass().method() == 44

doc="finished"
