# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import tempfile
import os

print("test tempfile")

if not tempfile.tempdir is None:
    print("tempfile.tempdir is not None: %s" % (tempfile.tempdir,))
else:
    print("tempfile.tempdir is None [OK]")

v = tempfile.gettempdir()
if type(v) != type(""):
    print("tempfile.gettempdir() returned %s (type=%s)" % (v, type(v)))

v = tempfile.gettempdirb()
if type(v) != type(b""):
    print("tempfile.gettempdirb() returned %s (type=%s)" % (v, type(v)))

## mkdtemp
try:
    tmp = tempfile.mkdtemp()
    os.rmdir(tmp)
    print("mkdtemp() [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkdtemp(): %s" % e)

try:
    tmp = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix")
    os.rmdir(tmp)
    print("mkdtemp(prefix='prefix-', suffix='-suffix') [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkdtemp(prefix='prefix-', suffix='-suffix'): %s" % e)

try:
    top = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix")
    tmp = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix", dir=top)
    os.rmdir(tmp)
    os.rmdir(top)
    print("mkdtemp(prefix='prefix-', suffix='-suffix', dir=top) [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkdtemp(prefix='prefix-', suffix='-suffix', dir=top): %s" % e)

try:
    top = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix")
    tmp = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix", dir=top)
    os.removedirs(top)
    print("mkdtemp(prefix='prefix-', suffix='-suffix', dir=top) [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkdtemp(prefix='prefix-', suffix='-suffix', dir=top): %s" % e)

try:
    tmp = tempfile.mkdtemp(prefix=b"prefix-", suffix="-suffix")
    print("missing exception!")
    os.rmdir(tmp)
except TypeError as e:
    print("caught: %s [OK]" % e)
except Exception as e:
    print("INVALID error caught: %s" % e)

def remove(fd, name):
    os.close(fd)
    os.remove(name)

## mkstemp
try:
    fd, tmp = tempfile.mkstemp()
    remove(fd, tmp)
    print("mkstemp() [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkstemp(): %s" % e)

try:
    fd, tmp = tempfile.mkstemp(prefix="prefix-", suffix="-suffix")
    remove(fd, tmp)
    print("mkstemp(prefix='prefix-', suffix='-suffix') [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkstemp(prefix='prefix-', suffix='-suffix'): %s" % e)

try:
    top = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix")
    fd, tmp = tempfile.mkstemp(prefix="prefix-", suffix="-suffix", dir=top)
    remove(fd, tmp)
    os.remove(top)
    print("mkstemp(prefix='prefix-', suffix='-suffix', dir=top) [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkstemp(prefix='prefix-', suffix='-suffix', dir=top): %s" % e)

try:
    top = tempfile.mkdtemp(prefix="prefix-", suffix="-suffix")
    fd, tmp = tempfile.mkstemp(prefix="prefix-", suffix="-suffix", dir=top)
    os.fdopen(fd).close() ## needed on Windows.
    os.removedirs(top)
    print("mkstemp(prefix='prefix-', suffix='-suffix', dir=top) [OK]")
except Exception as e:
    print("could not create tmp dir w/ mkstemp(prefix='prefix-', suffix='-suffix', dir=top): %s" % e)

try:
    fd, tmp = tempfile.mkstemp(prefix=b"prefix-", suffix="-suffix")
    print("missing exception!")
    remove(fd, tmp)
except TypeError as e:
    print("caught: %s [OK]" % e)
except Exception as e:
    print("INVALID error caught: %s" % e)

print("OK")
