#!/bin/bash
# Run testparser over python3 source

PY3SOURCE=~/Code/cpython

go install

find $PY3SOURCE -type f -name \*.py | grep -v "lib2to3/tests" | xargs testparser
