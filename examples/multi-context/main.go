// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	// This initializes gpython for runtime execution and is critical.
	// It defines forward-declared symbols and registers native built-in modules, such as sys and time.
	_ "github.com/go-python/gpython/modules"

	// This is the primary import for gpython.
	// It contains all symbols needed to fully compile and run python.
	"github.com/go-python/gpython/py"
)

func main() {

	// The total job count implies a fixed amount of work.
	// The number of workers is how many py.Context (in concurrent goroutines) to pull jobs off the queue.
	// One worker does all the work serially while N number of workers will (ideally) divides up.
	totalJobs := 20

	for i := 0; i < 10; i++ {
		numWorkers := i + 1
		elapsed := RunMultiPi(numWorkers, totalJobs)
		fmt.Printf("=====> %2d worker(s): %v\n\n", numWorkers, elapsed)

		// Give each trial a fresh start
		runtime.GC()
	}
}

var jobScript = `
pi = chud.pi_chudnovsky_bs(numDigits)
last_5 = pi % 100000
print("%s: last 5 digits of %d is %d (job #%0d)" % (WORKER_ID, numDigits, last_5, jobID))
`

var jobSrcTemplate = `
import pi_chudnovsky_bs as chud

WORKER_ID = "{{WORKER_ID}}"

print("%s ready!" % (WORKER_ID))
`

type worker struct {
	name string
	ctx  py.Context
	main *py.Module
	job  *py.Code
}

func (w *worker) compileTemplate(pySrc string) {
	pySrc = strings.Replace(pySrc, "{{WORKER_ID}}", w.name, -1)

	mainImpl := py.ModuleImpl{
		CodeSrc: pySrc,
	}

	var err error
	w.main, err = w.ctx.ModuleInit(&mainImpl)
	if err != nil {
		log.Fatal(err)
	}
}

func RunMultiPi(numWorkers, numTimes int) time.Duration {
	var workersRunning sync.WaitGroup

	fmt.Printf("Starting %d worker(s) to calculate %d jobs...\n", numWorkers, numTimes)

	jobPipe := make(chan int)
	go func() {
		for i := 0; i < numTimes; i++ {
			jobPipe <- i + 1
		}
		close(jobPipe)
	}()

	// Note that py.Code can be shared (accessed concurrently) since it is an inherently read-only object
	jobCode, err := py.Compile(jobScript, "<jobScript>", py.ExecMode, 0, true)
	if err != nil {
		log.Fatal("jobScript failed to comple")
	}

	workers := make([]worker, numWorkers)
	for i := 0; i < numWorkers; i++ {

		opts := py.DefaultContextOpts()

		// Make sure our import statement will find pi_chudnovsky_bs
		opts.SysPaths = append(opts.SysPaths, "..")

		workers[i] = worker{
			name: fmt.Sprintf("Worker #%d", i+1),
			ctx:  py.NewContext(opts),
			job:  jobCode,
		}

		workersRunning.Add(1)
	}

	startTime := time.Now()

	for i := range workers {
		w := workers[i]
		go func() {

			// Compiling can be concurrent since there is no associated py.Context
			w.compileTemplate(jobSrcTemplate)

			for jobID := range jobPipe {
				numDigits := 100000
				if jobID%2 == 0 {
					numDigits *= 10
				}
				py.SetAttrString(w.main.Globals, "numDigits", py.Int(numDigits))
				py.SetAttrString(w.main.Globals, "jobID", py.Int(jobID))
				w.ctx.RunCode(jobCode, w.main.Globals, w.main.Globals, nil)
			}
			workersRunning.Done()

			// This drives modules being able to perform cleanup and release resources 
			w.ctx.Close()
		}()
	}

	workersRunning.Wait()

	return time.Since(startTime)
}
