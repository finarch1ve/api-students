package database

import (
	"siakad-mini/app/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Student{},
		&model.Course{},
		&model.Enrollment{},
	)
}