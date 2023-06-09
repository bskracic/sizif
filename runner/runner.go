package runner

type Status string

const (
	Interrupted Status = "Interrupted"
	Failed             = "Failed"
	Finished           = "Finished"
)

type RunOptions struct {
	Stdin  string
	Script string
}

type RunResult struct {
	Stdout   string
	Message  string
	ExitCode int
	Status   Status
}

type Runner interface {
	Run(options *RunOptions) *RunResult
	Kill()
}
