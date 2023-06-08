package runner

type RunStatus string

const (
	Interrupted RunStatus = "Interrupted"
	Failed                = "Failed"
	Finished              = "Finished"
)

type RunOptions struct {
	Stdin  string
	Script string
}

type RunResult struct {
	Stdout   string
	Message  string
	ExitCode int
	Status   RunStatus
}

type Runner interface {
	Run(options *RunOptions) *RunResult
	Kill()
}
