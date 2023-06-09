package worker

import (
	"github.com/bskracic/sizif/db/model"
	"github.com/bskracic/sizif/runner"
	"github.com/bskracic/sizif/runtime"
	"gorm.io/gorm"
	"log"
	"time"
)

type Worker struct {
	db   *gorm.DB
	rnt  runtime.Runtime
	jobs chan uint
}

func NewWorker(db *gorm.DB, rnt runtime.Runtime, jobs chan uint) *Worker {
	return &Worker{db: db, rnt: rnt, jobs: jobs}
}

func (w *Worker) Start() {
	for jobId := range w.jobs {
		w.processJob(jobId)
	}
}

func (w *Worker) processJob(jobId uint) {
	// Find job
	var job model.Job
	if err := w.db.First(&job, jobId).Error; err != nil {
		log.Printf("job [%v] run failed: %v", jobId, err)
		return
	}
	r := &model.Run{
		JobId:  jobId,
		Status: model.Running,
	}
	w.db.Create(&r)

	// Execute script
	pr := runner.NewPythonRunner(w.rnt)
	rs := pr.Run(&runner.RunOptions{
		Script: job.Script,
	})

	// Write results
	r.Output = rs.Stdout
	r.Message = rs.Message
	r.Status = model.ToRunStatus(rs.Status)
	if err := w.db.Save(&r).Error; err != nil {
		log.Printf("run [%v] for job [%v] failed to save output", r.ID, r.JobId)
		return
	}
}

func (w *Worker) CheckJobsToSchedule() {
	jobList := model.RetrieveJobViews(w.db)
	for _, job := range jobList {
		// Do not queue job that is currently running
		if job.LastRunStatus == model.Running {
			continue
		}

		now := time.Now()
		if now.After(job.LastRun.Add(getTimeInterval(job.ScheduleUnit, job.ScheduleValue))) {
			if !w.tryEnqueue(job.ID) {
				log.Printf("cannot enqueue job: %v \n", job.ID)
			}
		}
	}
}

func getTimeInterval(unit string, value int) time.Duration {

	var dur time.Duration
	switch unit {
	case "second":
		dur = time.Second
	case "minute":
		dur = time.Minute
	case "hour":
		dur = time.Hour
	}
	return time.Duration(value) * dur
}

func (w *Worker) tryEnqueue(jobId uint) bool {
	select {
	case w.jobs <- jobId:
		return true
	default:
		return false
	}
}