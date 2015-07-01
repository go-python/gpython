#!/usr/bin/env python3
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
