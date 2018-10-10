# gpython

[![Build Status](https://travis-ci.org/go-python/gpython.svg?branch=master)](https://travis-ci.org/go-python/gpython)
[![codecov](https://codecov.io/gh/go-python/gpython/branch/master/graph/badge.svg)](https://codecov.io/gh/go-python/gpython)
[![GoDoc](https://godoc.org/github.com/go-python/gpython?status.svg)](https://godoc.org/github.com/go-python/gpython)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-python/gpython/blob/master/LICENSE)

gpython is a part re-implementation / part port of the Python 3.4
interpreter to the Go language, "batteries not included".

It includes:

  * runtime - using compatible byte code to python3.4
  * lexer
  * parser
  * compiler
  * interactive mode (REPL) ([try online!](https://gpython.org))

It does not include very many python modules as many of the core
modules are written in C not python.  The converted modules are:

  * builtins
  * marshal
  * math
  * time
  * sys

## Install

Gpython is a Go program and comes as a single binary file.

Download the relevant binary from here: https://github.com/go-python/gpython/releases

Or alternatively if you have Go installed use

    go get github.com/go-python/gpython

and this will build the binary in `$GOPATH/bin`.  You can then modify
the source and submit patches.

## Objectives

Gpython was written as a learning experiment to investigate how hard
porting Python to Go might be.  It turns out that all those C modules
are a significant barrier to making a fully functional port.

## Status

The project works well enough to parse all the code in the python 3.4
distribution and to compile and run python 3 programs which don't
depend on a module gpython doesn't support.

See the examples directory for some python programs which run with
gpython.

Speed hasn't been a goal of the conversions however it runs pystone at
about 20% of the speed of cpython.  The pi test runs quicker under
gpython as I think the Go long integer primitives are faster than the
Python ones.

There are many directions this project could go in.  I think the most
profitable would be to re-use the
[grumpy](https://github.com/grumpyhome/grumpy) runtime (which would mean
changing the object model).  This would give access to the C modules
that need to be ported and would give grumpy access to a compiler and
interpreter (gpython does support `eval` for instance).

I (@ncw) haven't had much time to work on gpython (I started it in
2013 and have worked on it very sporadically) so someone who wants to
take it in the next direction would be much appreciated.

## Limitations and Bugs

Lots!

## Similar projects

  * [grumpy](https://github.com/grumpyhome/grumpy) - a python to go transpiler

## License

This is licensed under the MIT licence, however it contains code which
was ported fairly directly directly from the cpython source code under
the (PSF LICENSE)[https://github.com/python/cpython/blob/master/LICENSE].
