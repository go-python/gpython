# Gpython Web

This implements a web viewable version of the gpython REPL.

This is done by compiling gpython into wasm and running that in the
browser.

[Try it online.](https://www.craig-wood.com/nick/gpython/)

## Build and run

`make build` will build with go wasm (you'll need go1.11 minimum)

`make serve` will run a local webserver you can see the results on

## Thanks

Thanks to [jQuery Terminal](https://terminal.jcubic.pl/) for the
terminal emulator and the go team for great [wasm
support](https://github.com/golang/go/wiki/WebAssembly).
