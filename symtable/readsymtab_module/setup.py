# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from distutils.core import setup, Extension

readsymtab = Extension('readsymtab', sources = ['readsymtab.c'])

setup (name = 'readsymtab',
        version = '1.0',
        description = 'Read the symbol table',
        ext_modules = [readsymtab]
)
