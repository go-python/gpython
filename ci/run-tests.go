// Copyright ©2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	log.SetPrefix("ci: ")
	log.SetFlags(0)

	start := time.Now()
	defer func() {
		log.Printf("elapsed time: %v\n", time.Since(start))
	}()

	var (
		race    = flag.Bool("race", false, "enable race detector")
		cover   = flag.String("coverpkg", "", "apply coverage analysis in each test to packages matching the patterns.")
		tags    = flag.String("tags", "", "build tags")
		verbose = flag.Bool("v", false, "enable verbose output")
	)

	flag.Parse()

	pkgs, err := pkgList()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("coverage.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	args := []string{"test"}

	if *verbose {
		args = append(args, "-v")
	}
	if *cover != "" {
		args = append(args, "-coverprofile=profile.out", "-covermode=atomic", "-coverpkg="+*cover)
	}
	if *tags != "" {
		args = append(args, "-tags="+*tags)
	}
	switch {
	case *race:
		args = append(args, "-race", "-timeout=20m")
	default:
		args = append(args, "-timeout=10m")
	}
	args = append(args, "")

	for _, pkg := range pkgs {
		args[len(args)-1] = pkg
		cmd := exec.Command("go", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		if *cover != "" {
			profile, err := os.ReadFile("profile.out")
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

func pkgList() ([]string, error) {
	out := new(bytes.Buffer)
	cmd := exec.Command("go", "list", "./...")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("could not get package list: %w", err)
	}

	var pkgs []string
	scan := bufio.NewScanner(out)
	for scan.Scan() {
		pkg := scan.Text()
		if strings.Contains(pkg, "vendor") {
			continue
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}
