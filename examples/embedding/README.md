## Embedding gpython

This example demonstrates how embed gpython into a Go application.  


### Why embed gpython?

Embedding a highly capable and familiar "interpreted" language allows your users
to easily augment app behavior, configuration, and customization -- all post-deployment.

Have you ever found an exciting software project but lose interest when you discover that you 
also need to learn an esoteric language schema?  In an era of limited attention span, 
most people are generally turned off if they have to learn a new language in addition to learning
to use your app.

If you consider [why use Python](https://www.stxnext.com/what-is-python-used-for/), then perhaps also 
consider your users could immediately feel interested to hear that your software offers
an additional familiar dimension of value. 

Python is widespread in finance, sciences of all kinds, hobbyist programming and is often 
endearingly regarded as most popular programming language for non-developers. 
If your application can be driven by embedded Python, then chances are others will 
feel excited and empowered that your project can be used out of the box 
with positive feelings of being in familiar territory.

### But what about the lack of python modules?

There are only be a small number of native modules available, but don't forget you have the entire
Go standard library and *any* Go package you can name at your fingertips to expose!  
This plus multi-context capability gives gpython enormous potential on how it can
serve your project.

So basically, gpython is only off the table if you need to run python that makes heavy use of 
modules that are only available in CPython.

### Packing List

|                       |                                                                   |
|---------------------- | ------------------------------------------------------------------|
| `main.go`             | if no args, runs in REPL mode, otherwise runs the given file      |
| `lib/mylib.py`        | models a library that your application would expose               |
| `lib/REPL-startup.py` | invoked by `main.go` when starting REPL mode                      |
| `mylib-demo.py`       | models a user-authored script that consumes `mylib`               |
| `mylib.module.go`     | Go implementation of `mylib_go` consumed by `mylib`               |


### Invoking a Python Script

```bash
$ cd examples/embedding/
$ go build .
$ ./embedding mylib-demo.py
```
```
Welcome to a gpython embedded example, 
    where your wildest Go-based python dreams come true!

==========================================================
        Python 3.4 (github.com/go-python/gpython)
        go1.17.6 on darwin amd64
==========================================================

Spring Break itinerary:
    Stop 1:   Miami, Florida    |   7 nights
    Stop 2:   Mallorca, Spain   |   3 nights
    Stop 3:   Ibiza, Spain      |  14 nights
    Stop 4:   Monaco            |  12 nights
###  Made with Vacaton 1.0 by Fletch F. Fletcher 

I bet Monaco will be the best!
```

### REPL Mode

```bash
$ ./embedding
```
```
=======  Entering REPL mode, press Ctrl+D to exit  =======

==========================================================
        Python 3.4 (github.com/go-python/gpython)
        go1.17.6 on darwin amd64
==========================================================

>>> v = Vacation("Spring Break", Stop("Florida", 3), Stop("Nice", 7))
>>> print(str(v))
Spring Break, 2 stop(s)
>>> v.PrintItinerary()
Spring Break itinerary:
    Stop 1:   Florida           |   3 nights
    Stop 2:   Nice              |   7 nights
###  Made with Vacaton 1.0 by Fletch F. Fletcher 
```

## Takeways

  - `main.go` demonstrates high-level convenience functions such as `py.RunFile()`.
  - Embedding any Go `struct` only requires that it implements `py.Object`, which is a single function: 
    `Type() *py.Type`
  - See [py/run.go](https://github.com/go-python/gpython/tree/master/py/run.go) for more about interpreter instances and `py.Context`
  - There are many helper functions available for you in [py/util.go](https://github.com/go-python/gpython/tree/master/py/util.go) and your contributions are welcome!