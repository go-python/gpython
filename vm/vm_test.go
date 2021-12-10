// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm_test

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/pytest"
)

func TestVm(t *testing.T) {
	pytest.RunTests(t, "tests")
}

func BenchmarkVM(b *testing.B) {
	pytest.RunBenchmarks(b, "benchmarks")
}

var jobSrcTemplate = `

doc="multi py.Ctx text"
WORKER_ID = "{{WORKER_ID}}"
def fib(n):
    if n == 0:
        return 0
    elif n == 1:
        return 1
    return fib(n - 2) + fib(n - 1)

x = {{FIB_TO}}
fx = fib(x)
print("%s says fib(%d) is %d" % (WORKER_ID, x, fx))
`

type worker struct {
	name string
	ctx  py.Ctx
}

func (w *worker) run(b testing.TB, pySrc string, countUpto int) {
	pySrc = strings.Replace(pySrc, "{{WORKER_ID}}", w.name, -1)
	pySrc = strings.Replace(pySrc, "{{FIB_TO}}", strconv.Itoa(countUpto), -1)

	module, code := pytest.CompileSrc(b, w.ctx, pySrc, w.name)
	_, err := w.ctx.RunCode(code, module.Globals, module.Globals, nil)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkCtx(b *testing.B) {
	numWorkers := 4
	workersRunning := sync.WaitGroup{}

	numJobs := 35
	fmt.Printf("Starting %d workers to process %d jobs...\n", numWorkers, numJobs)

	jobPipe := make(chan int)
	go func() {
		for i := 1; i <= numJobs; i++ {
			jobPipe <- i
		}
		close(jobPipe)
	}()

	workers := make([]worker, numWorkers)
	for i := 0; i < numWorkers; i++ {

		workers[i] = worker{
			name: fmt.Sprintf("Worker #%d", i+1),
			ctx:  py.NewCtx(py.DefaultCtxOpts()),
		}

		workersRunning.Add(1)
		w := workers[i]
		go func() {
			for jobID := range jobPipe {
				w.run(b, jobSrcTemplate, jobID)
				//fmt.Printf("### %s finished job %v ###\n", w.name, jobID)
			}
			workersRunning.Done()
		}()
	}

	workersRunning.Wait()
}
