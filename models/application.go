// models/application.go
package models

import (
	"github.com/jinzhu/gorm"
)

// ApplicationState represents the state of a job application
type ApplicationState string

const (
	Pending   ApplicationState = "pending"
	Reviewed  ApplicationState = "reviewed"
	Accepted  ApplicationState = "accepted"
	Rejected  ApplicationState = "rejected"
	Withdrawn ApplicationState = "withdrawn"
)

// Application represents a job application
type Application struct {
	gorm.Model
	JobID       int    `gorm:"not null;unique_index:idx_job_user;index;constraint:OnDelete:CASCADE;"`
	Job         Job    `gorm:"foreignkey:JobID"`
	ApplicantID int    `gorm:"not null;unique_index:idx_job_user;index;constraint:OnDelete:CASCADE;"`
	Applicant   User   `gorm:"foreignkey:UserID"`
	ResumeURL   string `gorm:"not null"`
	CoverNote   string // Optional field
	State       ApplicationState
}

// CreateApplication creates a new job application in the database
func CreateApplication(application *Application) error {
	if err := db.Create(application).Error; err != nil {
		return err
	}
	return nil
}

// GetApplicationByJobAndUser retrieves a job application by job ID and user ID
func GetApplicationByJobAndUser(jobID int, userID int) (*Application, error) {
	var application Application
	if err := db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

// GetApplicationByID retrieves an application by its ID
func GetApplicationByID(id uint) (*Application, error) {
	var application Application
	if err := db.First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

// UpdateApplicationStatus updates the status of an application
func UpdateApplicationStatus(id uint, status ApplicationState) error {
	var application Application
	if err := db.First(&application, id).Error; err != nil {
		return err
	}

	// Update the status
	application.State = status
	if err := db.Save(&application).Error; err != nil {
		return err
	}
	return nil
}
