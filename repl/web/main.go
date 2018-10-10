// An online REPL for gpython using wasm

// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package main

import (
	"log"
	"runtime"

	"github.com/gopherjs/gopherwasm/js" // gopherjs to wasm converter shim

	// import required modules
	_ "github.com/go-python/gpython/builtin"
	_ "github.com/go-python/gpython/math"
	"github.com/go-python/gpython/repl"
	_ "github.com/go-python/gpython/sys"
	_ "github.com/go-python/gpython/time"
)

// Implement the replUI interface
type termIO struct {
	js.Value
}

// SetPrompt sets the UI prompt
func (t *termIO) SetPrompt(prompt string) {
	t.Call("set_prompt", prompt)
}

// Print outputs the string to the output
func (t *termIO) Print(out string) {
	t.Call("echo", out)
}

var document js.Value

func isUndefined(node js.Value) bool {
	return node == js.Undefined()
}

func getElementById(name string) js.Value {
	node := document.Call("getElementById", name)
	if isUndefined(node) {
		log.Fatalf("Couldn't find element %q", name)
	}
	return node
}

func running() string {
	switch {
	case runtime.GOOS == "js" && runtime.GOARCH == "wasm":
		return "go/wasm"
	case runtime.GOARCH == "js":
		return "gopherjs"
	}
	return "unknown"
}

func main() {
	document = js.Global().Get("document")
	if isUndefined(document) {
		log.Fatalf("Didn't find document - not running in browser")
	}

	// Clear the loading text
	termNode := getElementById("term")
	termNode.Set("innerHTML", "")

	// work out what we are running on and mark active
	tech := running()
	node := getElementById(tech)
	node.Get("classList").Call("add", "active")

	// Make a repl referring to an empty term for the moment
	REPL := repl.New()
	cb := js.NewCallback(func(args []js.Value) {
		REPL.Run(args[0].String())
	})

	// Create a jquery terminal instance
	opts := js.ValueOf(map[string]interface{}{
		"greetings": "Gpython 3.4.0 running in your browser with " + tech,
		"name":      "gpython",
		"prompt":    repl.NormalPrompt,
	})
	terminal := js.Global().Call("$", "#term").Call("terminal", cb, opts)

	// Send the console log direct to the terminal
	js.Global().Get("console").Set("log", terminal.Get("echo"))

	// Set the implementation of term
	REPL.SetUI(&termIO{terminal})

	// wait for callbacks
	select {}
}
