package model

import "gorm.io/gorm"

type RunnerType struct {
	gorm.Model
	Name string
}

func (RunnerType) TableName() string {
	return "runner_type"
}
