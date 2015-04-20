#!/bin/sh
set -e
rm -rf build
PYTHON=/opt/python3.4/bin/python3.4
INCLUDE=/opt/python3.4/include/python3.4m
$PYTHON setup.py build_ext --inplace
cp -av *.so *.dll ..
