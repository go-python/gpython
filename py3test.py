#!/usr/bin/env python3.4

# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

"""
Check all the tests work with python3.4 and gpython

Note that this isn't quite the same as running the unit tests - the
unit tests should be preferred.  This is a quick check to make sure
the tests run with python3.
"""

import os
import sys
from subprocess import Popen, PIPE, STDOUT

testwith = ("python3.4", "gpython")

def runtests(dirpath, filenames):
    """Run the tests found"""
    print("Running tests in %s" % dirpath)
    for name in filenames:
        if not name.endswith(".py") or name.startswith("lib") or name.startswith("raise"):
            continue
        print("Testing %s" % name)
        fullpath = os.path.join(dirpath, name)
        for cmd in testwith:
            prog = [cmd, fullpath]
            p = Popen(prog, stdin=PIPE, stdout=PIPE, stderr=STDOUT, close_fds=True)
            stdout, stderr = p.communicate("")
            rc = p.returncode
            if rc != 0:
                print("*** %s %s Fail ***" % (cmd, fullpath))
                print("="*60)
                sys.stdout.write(stdout.decode("utf-8"))
                print("="*60)
        
def main():
    binary = os.path.abspath(__file__)
    home = os.path.dirname(binary)
    os.chdir(home)
    print("Scanning %s for tests" % home)

    for dirpath, dirnames, filenames in os.walk("."):
        if os.path.basename(dirpath) == "tests":
            runtests(dirpath, filenames)

if __name__ == "__main__":
    main()
