## Embedding gpython

This is an example demonstrating how to embed gpython into a Go application.  


### Why embed gpython?

Embedding a highly capable and familiar "interpreted" language allows your users
to easily augment app behavior, configuration, and customization -- all post-deployment.

Have you ever discovered an exciting software project but lost interest when you had to also
learn its esoteric language schema?  In an era of limited attention span, 
most people are generally turned off if they have to learn a new language in addition to learning
to use your app.

If you consider [why use Python](https://www.stxnext.com/what-is-python-used-for/), then perhaps also 
consider that your users will be interested to hear that your software offers
even more value that it can be driven from a scripting language they already know.

Python is widespread in finance, sciences, hobbyist programming and is often 
endearingly regarded as most popular programming language for non-developers. 
If your application can be driven by embedded Python, then chances are others will 
feel excited and empowered that your project can be used out of the box 
and feel like familiar territory.

### But what about the lack of python modules?

There are only be a small number of native modules available, but don't forget you have the entire
Go standard library and *any* Go package you can name at your fingertips to expose!  
This plus multi-context capability gives gpython enormous potential on how it can
serve you.

So basically, gpython is only off the table if you need to run python that makes heavy use of 
modules that are only available in CPython.

### Packing List

|                          |                                                                   |
|------------------------- | ------------------------------------------------------------------|
| `main.go`                | if no args, runs in REPL mode, otherwise runs the given file      |
| `lib/mylib.py`           | models a library that your application would expose for users     |
| `lib/REPL-startup.py`    | invoked by `main.go` when starting REPL mode                      |
| `testdata/mylib-demo.py` | models a user-authored script that consumes `mylib`               |
| `mylib.module.go`        | Go implementation of `mylib_go` consumed by `mylib`               |


### Invoking a Python Script

```bash
$ cd examples/embedding/
$ go build .
$ ./embedding ./testdata/mylib-demo.py
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
  - Embedding a Go `struct` as a Python object only requires that it implement `py.Object`, which is a single function: 
    `Type() *py.Type`
  - See [py/run.go](https://github.com/go-python/gpython/tree/main/py/run.go) for more about interpreter instances and `py.Context`
  - Helper functions are available in [py/util.go](https://github.com/go-python/gpython/tree/main/py/util.go) and your contributions are welcome!
