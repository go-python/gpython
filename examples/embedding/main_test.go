package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

const embeddingTestOutput = "embedding-test-output.txt"

func TestEmbeddedExample(t *testing.T) {

	cmd := exec.Command("go", "build", ".")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to compile embedding example: %v", err)
	}

	out := new(bytes.Buffer)
	cmd = exec.Command("./embedding", "mylib-demo.py")
	cmd.Stdout = out

	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to run embedding binary: %v", err)
	}

	resetTest := false // true
	testOutput := out.Bytes()
	if resetTest {
		err = os.WriteFile(embeddingTestOutput, testOutput, 0644)
		if err != nil {
			t.Fatalf("failed to write test output: %v", err)
		}
	} else {
		mustMatch, err := os.ReadFile(embeddingTestOutput)
		if err != nil {
			t.Fatalf("failed read %q", embeddingTestOutput)
		}
		if !bytes.Equal(testOutput, mustMatch) {
			t.Fatalf("embedded test output did not match accepted output from %q", embeddingTestOutput)
		}
	}
}
