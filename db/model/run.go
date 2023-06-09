package model

import (
	"fmt"
	"github.com/bskracic/sizif/runner"
	"gorm.io/gorm"
	"time"
)

type Run struct {
	gorm.Model
	Container  string    `gorm:"container"`
	Output     string    `gorm:"output"`
	Message    string    `gorm:"message"`
	FinishedAt time.Time `gorm:"finished_at"`
	Status     RunStatus `gorm:"status"`
	JobId      uint      `gorm:"column:job_id"`
	Job        Job
}

func (Run) TableName() string {
	return "run"
}

type RunStatus int

const (
	Running RunStatus = iota + 1
	Finished
	Interrupted
	Failed
)

func (e RunStatus) String() string {
	s := [...]string{"running", "finished", "interrupted", "failed"}
	if e < Running || e > Failed {
		return fmt.Sprintf("RunStatus(%d)", int(e))
	}
	return s[(e)-1]
}

func (e RunStatus) IsValid() bool {
	switch e {
	case Running, Finished, Interrupted, Failed:
		return true
	}
	return false
}

func ToRunStatus(s runner.Status) RunStatus {
	switch s {
	case runner.Interrupted:
		return Interrupted
	case runner.Failed:
		return Failed
	default:
		return Finished
	}
}
