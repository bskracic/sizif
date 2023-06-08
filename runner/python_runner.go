package runner

import (
	"time"

	"github.com/bskracic/sizif/runtime"
)

const (
	imageName          = "python:3.10-alpine"
	language           = "python3"
	fileName           = "main.py"
	contBasePath       = "/"
	executionTimeLimit = 10000
)

var executeCmd = []string{"python3", "main.py"}

var spec = runtime.Specs{
	Lang:  language,
	Image: imageName,
}

type PythonRunner struct {
	runtime     runtime.Runtime
	containerID string
}

func NewPythonRunner(r runtime.Runtime) *PythonRunner {
	id := r.Prepare(spec)
	return &PythonRunner{runtime: r, containerID: id}
}

func (pr *PythonRunner) copyFiles(options *RunOptions) error {
	return pr.runtime.CopyFile(pr.containerID, options.Script, fileName, contBasePath)
}

func (pr *PythonRunner) Run(options *RunOptions) *RunResult {
	// First copy the main script
	err := pr.copyFiles(options)
	if err != nil {
		return createErrorResult(err)
	}

	ch := make(chan *runtime.ExecResult, 1)
	go func(ch chan *runtime.ExecResult) {
		err := pr.runtime.Exec(pr.containerID, executeCmd, ch)
		if err != nil {
			ch <- &runtime.ExecResult{ExitCode: -1, Stderr: err.Error()}
		}
	}(ch)

	var rr RunResult
	select {
	case res := <-ch:
		rr.ExitCode = res.ExitCode
		if res.ExitCode != 0 {
			rr.Status = Failed
			rr.Stdout = res.Stderr
		} else {
			rr.Status = Finished
			rr.Stdout = res.Stdout
		}
	case <-time.After(time.Duration(executionTimeLimit) * time.Millisecond):
		rr.Status = Interrupted
	}
	pr.runtime.Kill(pr.containerID)

	return &rr
}

func createErrorResult(err error) *RunResult {
	return &RunResult{Message: err.Error(), Status: Failed, ExitCode: -1}
}

func (pr *PythonRunner) Kill() {
	go pr.runtime.Kill(pr.containerID)
}
