// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var regen = flag.Bool("regen", false, "regenerate golden files")

func TestEmbeddedExample(t *testing.T) {

	tmp, err := os.MkdirTemp("", "go-python-embedding-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	exe := filepath.Join(tmp, "out.exe")
	cmd := exec.Command("go", "build", "-o", exe, ".")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to compile embedding example: %+v", err)
	}

	out := new(bytes.Buffer)
	cmd = exec.Command(exe, "mylib-demo.py")
	cmd.Stdout = out
	cmd.Stderr = out

	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to run embedding binary: %+v", err)
	}

	const fname = "testdata/embedding_out_golden.txt"

	got := out.Bytes()

	flag.Parse()
	if *regen {
		err = os.WriteFile(fname, got, 0644)
		if err != nil {
			t.Fatalf("could not write golden file: %+v", err)
		}
	}

	want, err := os.ReadFile(fname)
	if err != nil {
		t.Fatalf("could not read golden file: %+v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("stdout differ:\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}
