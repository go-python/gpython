// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetPrefix("ci: ")
	log.SetFlags(0)

	var (
		race  = flag.Bool("race", false, "enable race detector")
		cover = flag.Bool("cover", false, "enable code coverage")
		tags  = flag.String("tags", "", "build tags")
	)

	flag.Parse()

	out := new(bytes.Buffer)
	cmd := exec.Command("go", "list", "./...")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("coverage.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	args := []string{"test"}

	if *cover {
		args = append(args, "-coverprofile=profile.out", "-covermode=atomic")
	}
	if *tags != "" {
		args = append(args, "-tags="+*tags)
	}
	if *race {
		args = append(args, "-race")
	}
	args = append(args, "")

	scan := bufio.NewScanner(out)
	for scan.Scan() {
		pkg := scan.Text()
		if strings.Contains(pkg, "vendor") {
			continue
		}
		args[len(args)-1] = pkg
		cmd := exec.Command("go", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		if *cover {
			profile, err := ioutil.ReadFile("profile.out")
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.Write(profile)
			if err != nil {
				log.Fatal(err)
			}
			os.Remove("profile.out")
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
