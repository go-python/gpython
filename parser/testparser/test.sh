#!/bin/bash
# Run testparser over python3 source
#
# Pass in args to be passed to testparser, eg -c, -l
# Parses by default

PY3SOURCE=~/Code/cpython

go install

# Grep out python2 source code which we can't parse and files with deliberate syntax errors
find $PY3SOURCE -type f -name \*.py | egrep -v "Lib/(lib2to3/tests|test/bad.*py)|Tools/(hg|msi|test2to3)/" | xargs testparser "$@"

#find $PY3SOURCE -type f -name \*.py | egrep -v "Lib/(lib2to3/tests|test/bad.*py)|Tools/(hg|msi|test2to3)/" | xargs ./py3compile.py "$@"

