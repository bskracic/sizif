package db

import (
	"github.com/bskracic/sizif/db/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	//c := config.GetFromEnv()
	c := "host=localhost user=sizifUser password=supersecretpassword dbname=sizifdb port=5432 sslmode=disable"
	for {
		db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		dbInit(db)
		return db
	}
}

func dbInit(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.RunnerType{},
		&model.Job{},
	)
	if err != nil {
		panic(err)
	}
	python := &model.RunnerType{Name: "python"}
	db.Create(&python)
}
