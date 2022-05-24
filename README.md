# gpython

[![Build Status](https://github.com/go-python/gpython/workflows/CI/badge.svg)](https://github.com/go-python/gpython/actions)
[![codecov](https://codecov.io/gh/go-python/gpython/branch/main/graph/badge.svg)](https://codecov.io/gh/go-python/gpython)
[![GoDoc](https://godoc.org/github.com/go-python/gpython?status.svg)](https://godoc.org/github.com/go-python/gpython)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-python/gpython/blob/main/LICENSE)

gpython is a part re-implementation, part port of the Python 3.4
interpreter in Go.  Although there are many areas of improvement,
it stands as an noteworthy achievement in capability and potential.
 
gpython includes:

  * lexer, parser, and compiler
  * runtime and high-level convenience functions
  * multi-context interpreter instancing
  * easy embedding into your Go application
  * interactive mode (REPL) ([try online!](https://gpython.org))


gpython does not include many python modules as many of the core
modules are written in C not python.  The converted modules are:

  * builtins
  * marshal
  * math
  * time
  * sys

## Install

Download directly from the [releases page](https://github.com/go-python/gpython/releases) 

Or if you have Go installed:

    go install github.com/go-python/gpython

## Objectives

gpython started as an experiment to investigate how hard
porting Python to Go might be.  It turns out that all those C modules
are a significant barrier to making gpython a complete replacement
to CPython.  

However, to those who want to embed a highly popular and known language
into their Go application, gpython could be a great choice over less
capable (or lesser known) alternatives.

## Status

gpython currently:
 - Parses all the code in the Python 3.4 distribution
 - Runs Python 3 for the modules that are currently supported
 - Supports concurrent multi-interpreter ("multi-context") execution

Speed hasn't been a goal of the conversions however it runs pystone at
about 20% of the speed of CPython.  A [π computation test](https://github.com/go-python/gpython/tree/main/examples/pi_chudnovsky_bs.py) runs quicker under
gpython as the Go long integer primitives are likely faster than the
Python ones.

@ncw started gpython in 2013 and work on is sporadic. If you or someone
you know would be interested to take it futher, it would be much appreciated.

## Getting Started

The [embedding example](https://github.com/go-python/gpython/tree/main/examples/embedding) demonstrates how to
easily embed and invoke gpython from any Go application.

Of interest, gpython is able to run multiple interpreter instances simultaneously,
allowing you to embed gpython naturally into your Go application.  This makes it
possible to use gpython in a server situation where complete interpreter 
independence is paramount.  See this in action in the [multi-context example](https://github.com/go-python/gpython/tree/main/examples/multi-context).
 
If you are looking to get involved, a light and easy place to start is adding more convenience functions to [py/util.go](https://github.com/go-python/gpython/tree/main/py/util.go).  See [notes.txt](https://github.com/go-python/gpython/blob/main/notes.txt) for bigger ideas.


## Other Projects of Interest

  * [grumpy](https://github.com/grumpyhome/grumpy) - a python to go transpiler

## Community

You can chat with the go-python community (or which gpython is part)
at [go-python@googlegroups.com](https://groups.google.com/forum/#!forum/go-python)
or on the [Gophers Slack](https://gophers.slack.com/) in the `#go-python` channel.

## License

This is licensed under the MIT licence, however it contains code which
was ported fairly directly directly from the CPython source code under
the [PSF LICENSE](https://github.com/python/cpython/blob/main/LICENSE).
