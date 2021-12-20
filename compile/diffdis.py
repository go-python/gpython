#!/usr/bin/env python3.4

# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Diff two bytecode strings

import sys
import io
import os
from dis import dis
from difflib import unified_diff
from tempfile import NamedTemporaryFile

def disassemble(code):
    """Disassemble code into string"""
    out = io.StringIO()
    dis(code, file=out)
    return out.getvalue()

def main():
    assert len(sys.argv) == 3, "Need two arguments"
    a, b = (bytes(x, "latin1").decode("unicode_escape").encode("latin1") for x in sys.argv[1:])
    a_dissasembly = disassemble(a)
    b_dissasembly = disassemble(b)
    for line in unified_diff(a_dissasembly.split("\n"), b_dissasembly.split("\n"), fromfile="want", tofile="got"):
        print(line)

if __name__ == "__main__":
    main()
