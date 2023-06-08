package runner

import (
	"errors"
	"github.com/bskracic/sizif/runtime"
)

func CreateRunner(r runtime.Runtime, runnerType string) (Runner, error) {
	var runner Runner
	var err error = nil
	switch runnerType {
	case "python":
		runner = NewPythonRunner(r)
	default:
		runner = nil
		err = errors.New("no such runner")
	}

	return runner, err
}
