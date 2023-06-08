package model

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	Name      string `gorm:"column:name"`
	RunnerId  uint   `gorm:"column:runner_id;foreignKey:ID;references:runner_type(ID)"`
	Script    string `gorm:"column:script"`
	Type      string `gorm:"column:type"`
	NextJobID *uint  `gorm:"column:next_job_id,omitempty;foreignKey:ID;references:job(ID)"`
	Schedule  string `json:"schedule"`
}

func (Job) TableName() string {
	return "job"
}
