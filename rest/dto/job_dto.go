package dto

import (
	"github.com/bskracic/sizif/db/model"
	"github.com/gin-gonic/gin"
	"time"
)

type CreateJobDto struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	RunnerTypeId  uint   `json:"runner_type_id"`
	Script        string `json:"script"`
	Type          string `json:"type"`
	NextJobID     *uint  `json:"next_job_id"`
	ScheduleUnit  string `json:"schedule_unit"`
	ScheduleValue int    `json:"schedule_value"`
}

func ToJob(c *gin.Context) (*model.Job, error) {

	var dto CreateJobDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		return nil, err
	}

	var j model.Job
	j.ID = dto.Id
	j.Name = dto.Name
	j.Script = dto.Script
	j.RunnerTypeId = dto.RunnerTypeId
	j.Type = dto.Type
	//j.NextJobID = dto.NextJobID
	j.ScheduleValue = dto.ScheduleValue
	j.ScheduleUnit = dto.ScheduleUnit

	return &j, nil
}

type JobDto struct {
	Id            uint      `json:"id"`
	Name          string    `json:"name"`
	Requirements  string    `json:"requirements"`
	Script        string    `json:"script"`
	Type          string    `json:"type"`
	ScheduleUnit  string    `json:"schedule_unit"`
	ScheduleValue int       `json:"schedule_value"`
	LastRun       time.Time `json:"last_execution"`
	LastRunStatus string    `json:"last_run_status"`
	Runner        string    `json:"runner"`
}

func ToJobDto(job model.JobView) JobDto {
	var dto JobDto
	dto.Id = job.ID
	dto.Name = job.Name
	dto.Requirements = job.Requirements
	dto.Script = job.Script
	dto.Type = job.Type
	dto.ScheduleUnit = job.ScheduleUnit
	dto.ScheduleValue = job.ScheduleValue
	dto.LastRun = job.LastRun
	dto.LastRunStatus = job.LastRunStatus.String()
	dto.Runner = job.Runner.Name

	return dto
}

func ToJobListDto(jobs []model.JobView) []JobDto {
	var list []JobDto
	for _, job := range jobs {
		list = append(list, ToJobDto(job))
	}

	return list
}
