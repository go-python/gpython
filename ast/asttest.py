#!/usr/bin/env python3

# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""
Do some ast stuff
"""

import sys
import ast

def dump(path):
    print(path)
    a = ast.parse(open(path).read(), path)
    print(ast.dump(a, annotate_fields=True, include_attributes=False))

def main():
    for path in sys.argv[1:]:
        dump(path)

if __name__ == "__main__":
    main()
