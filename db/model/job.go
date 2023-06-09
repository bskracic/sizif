package model

import (
	"gorm.io/gorm"
	"time"
)

type Job struct {
	gorm.Model
	Name          string     `gorm:"column:name"`
	Requirements  string     `gorm:"requirements"`
	Script        string     `gorm:"column:script"`
	Type          string     `gorm:"column:type"`
	ScheduleUnit  string     `gorm:"column:schedule_unit"`
	ScheduleValue int        `gorm:"column:schedule_value"`
	RunnerTypeId  uint       `gorm:"column:runner_type_id"`
	Runner        RunnerType `gorm:"foreignKey:RunnerTypeId"`
	Runs          []Run      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:job_id"`
}

func (Job) TableName() string {
	return "job"
}

type JobView struct {
	gorm.Model
	Name          string     `gorm:"column:name"`
	Requirements  string     `gorm:"requirements"`
	Script        string     `gorm:"column:script"`
	Type          string     `gorm:"column:type"`
	ScheduleUnit  string     `gorm:"column:schedule_unit"`
	ScheduleValue int        `gorm:"column:schedule_value"`
	LastRun       time.Time  `gorm:"column:last_run"`
	LastRunStatus RunStatus  `gorm:"column:last_run_status"`
	RunnerTypeId  uint       `gorm:"column:runner_type_id"`
	Runner        RunnerType `gorm:"foreignKey:RunnerTypeId"`
}

func RetrieveJobViews(db *gorm.DB) []JobView {
	var jobs []JobView
	query := `select
    j.*,
    (select finished_at from run r where r.job_id = j.id order by finished_at desc limit 1) as last_run,
    (select status from run r where r.job_id = j.id order by finished_at desc limit 1) as last_run_status
	from job j;`

	db.Raw(query).Scan(&jobs)

	return jobs
}
