package model

import "gorm.io/gorm"

type RunnerType struct {
	gorm.Model
	Name string
	Jobs []Job `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:runner_type_id"`
}

func (RunnerType) TableName() string {
	return "runner_type"
}
