package dto

import (
	"github.com/bskracic/sizif/db/model"
	"github.com/gin-gonic/gin"
)

type JobDto struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	RunnerId  uint   `json:"runner"`
	Script    string `json:"script"`
	Type      string `json:"type"`
	NextJobID *uint  `json:"next_job_id,omitempty"`
	Schedule  string `json:"schedule"`
}

func ToJob(c *gin.Context) (*model.Job, error) {

	var dto JobDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		return nil, err
	}

	var j model.Job
	j.ID = dto.Id
	j.Name = dto.Name
	j.Script = dto.Script
	j.RunnerId = dto.RunnerId
	j.Type = dto.Type
	j.NextJobID = dto.NextJobID
	j.Schedule = dto.Schedule

	return &j, nil
}
