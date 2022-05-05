# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os

print("test os")
print("os.error: ", os.error)
print("os.getenv($GPYTHON_TEST_HOME)=", os.getenv("GPYTHON_TEST_HOME"))
os.putenv("GPYTHON_TEST_HOME", "/home/go")
print("os.environ($GPYTHON_TEST_HOME)=", os.environ.get("GPYTHON_TEST_HOME"))
print("os.getenv($GPYTHON_TEST_HOME)=", os.getenv("GPYTHON_TEST_HOME"))
os.unsetenv("GPYTHON_TEST_HOME")
print("os.unsetenv($GPYTHON_TEST_HOME)=", os.getenv("GPYTHON_TEST_HOME"))

if not os.error is OSError:
    print("os.error is not OSError!")
else:
    print("os.error is OSError [OK]")

## FIXME(sbinet): check returned value with a known one
## (ie: when os.mkdir is implemented)
if os.getcwd() == None:
    print("os.getcwd() == None !")
else:
    print("os.getcwd() != None [OK]")

## FIXME(sbinet): check returned value with a known one
## (ie: when os.mkdir is implemented)
if os.getcwdb() == None:
    print("os.getcwdb() == None !")
else:
    print("os.getcwdb() != None [OK]")

print("os.system('echo hello')...")
if os.name != "nt":
    os.system('echo hello')
else: ## FIXME(sbinet): find a way to test this nicely
    print("hello\n")

if os.getpid() > 1:
    print("os.getpid is greater than 1 [OK]")
else:
    print("invalid os.getpid: ", os.getpid())

orig = os.getcwd()
testdir = "/"
if os.name == "nt":
    testdir = "C:\\"
os.chdir(testdir)
if os.getcwd() != testdir:
    print("invalid getcwd() after os.chdir:",os.getcwd())
else:
    print("os.chdir(testdir) [OK]")
os.chdir(orig)

try:
    os.chdir(1)
    print("expected an error with os.chdir(1)")
except TypeError:
    print("os.chdir(1) failed [OK]")

try:
    os.environ.get(15)
    print("expected an error with os.environ.get(15)")
except KeyError:
    print("os.environ.get(15) failed [OK]")

try:
    os.putenv()
    print("expected an error with os.putenv()")
except TypeError:
    print("os.putenv() failed [OK]")

try:
    os.unsetenv()
    print("expected an error with os.unsetenv()")
except TypeError:
    print("os.unsetenv() failed [OK]")

try:
    os.getenv()
    print("expected an error with os.getenv()")
except TypeError:
    print("os.getenv() failed [OK]")

try:
    os.unsetenv("FOO", "BAR")
    print("expected an error with os.unsetenv(\"FOO\", \"BAR\")")
except TypeError:
    print("os.unsetenv(\"FOO\", \"BAR\") failed [OK]")

if bytes(os.getcwd(), "utf-8") == os.getcwdb():
    print('bytes(os.getcwd(), "utf-8") == os.getcwdb() [OK]')
else:
    print('expected: bytes(os.getcwd(), "utf-8") == os.getcwdb()')

golden = {
        "posix": {
            "sep": "/",
            "pathsep": ":",
            "linesep": "\n",
            "devnull": "/dev/null",
            "altsep": None
        },
        "nt": {
            "sep": "\\",
            "pathsep": ";",
            "linesep": "\r\n",
            "devnull": "nul",
            "altsep": "/"
        },
}[os.name]

for k in ("sep", "pathsep", "linesep", "devnull", "altsep"):
    if getattr(os, k) != golden[k]:
        print("invalid os."+k+": got=",getattr(os,k),", want=", golden[k])
    else:
        print("os."+k+": [OK]")

print("OK")
