#!/usr/bin/env python3

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
from collections import defaultdict

py_version = "python3.4"

opt_install = "/opt/"+py_version

bin_dirs = os.environ["PATH"].split(os.pathsep) + [
    opt_install+"/bin",
    os.path.join(os.environ["HOME"], "bin/"+py_version+"/bin"),
]

def find_python():
    """Find a version of python to run"""
    for bin_dir in bin_dirs:
        path = os.path.join(bin_dir, py_version)
        if os.path.exists(path):
            return path
    print("Couldn't find "+py_version+" on $PATH or "+" or ".join(bin_dirs[-2:]))
    print("Install "+py_version+" by doing:")
    print("  sudo mkdir -p "+opt_install)
    print("  sudo chown $USER "+opt_install)
    print("  ./bin/install-python.sh "+opt_install+'"')
    sys.exit(1)

testwith = [find_python(), "gpython"]

def runtests(dirpath, filenames, failures):
    """Run the tests found accumulating failures"""
    print("Running tests in %s" % dirpath)
    for name in filenames:
        if not name.endswith(".py") or name.startswith("lib") or name.startswith("raise"):
            continue
        #print(" - %s" % name)
        fullpath = os.path.join(dirpath, name)
        for cmd in testwith:
            prog = [cmd, fullpath]
            p = Popen(prog, stdin=PIPE, stdout=PIPE, stderr=STDOUT, close_fds=True)
            stdout, stderr = p.communicate("")
            rc = p.returncode
            if rc != 0:
                failures[cmd][fullpath].append(stdout.decode("utf-8"))
    return failures

def main():
    binary = os.path.abspath(__file__)
    home = os.path.dirname(binary)
    os.chdir(home)
    print("Scanning %s for tests" % home)

    failures = defaultdict(lambda: defaultdict(list))
    for dirpath, dirnames, filenames in os.walk("."):
        if os.path.basename(dirpath) == "tests":
            runtests(dirpath, filenames, failures)

    if not failures:
        print("All OK")
        return

    print()

    sep = "="*60+"\n"
    sep2 = "-"*60+"\n"

    for cmd in sorted(failures.keys()):
        for path in sorted(failures[cmd].keys()):
            print(sep+"Failures for "+cmd+" in "+path)
            sys.stdout.write(sep+sep2.join(failures[cmd][path])+sep)
        print()
    sys.exit(1)


if __name__ == "__main__":
    main()
