package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const embeddingTestOutput = "testdata/embedding_out_golden.txt"

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
		t.Fatalf("failed to compile embedding example: %v", err)
	}

	out := new(bytes.Buffer)
	cmd = exec.Command(exe, "mylib-demo.py")
	cmd.Stdout = out

	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to run embedding binary: %v", err)
	}

	testOutput := out.Bytes()

	flag.Parse()
	if *regen {
		err = os.WriteFile(embeddingTestOutput, testOutput, 0644)
		if err != nil {
			t.Fatalf("failed to write test output: %v", err)
		}
	}

	mustMatch, err := os.ReadFile(embeddingTestOutput)
	if err != nil {
		t.Fatalf("failed read %q", embeddingTestOutput)
	}
	if !bytes.Equal(testOutput, mustMatch) {
		t.Fatalf("embedded test output did not match accepted output from %q", embeddingTestOutput)
	}
}
