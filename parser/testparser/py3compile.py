#!/usr/bin/env python3

# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""
Compile any files passed in on the command line
"""

import sys
for path in sys.argv[1:]:
    print("Compiling %s" % path)
    with open(path) as f:
        try:
            data = f.read()
            compile(data, path, "exec")
        except Exception as e:
            print("Failed: %s" % e)
