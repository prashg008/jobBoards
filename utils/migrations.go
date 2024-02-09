package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/prashg008/jobBoards/models"
)

func MigrateDB(db *gorm.DB) error {
	// Migrate user, job, and application tables
	if err := db.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.Application{},
		&models.Notification{}).Error; err != nil {
		return err
	}

	return nil
}
