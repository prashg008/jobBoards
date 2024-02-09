// models/job.go
package models

import (
	"github.com/jinzhu/gorm"
)

// Job represents a job posting
type Job struct {
	gorm.Model
	PostedByID          int    `gorm:"not null;constraint:OnDelete:CASCADE;"`
	PostedBy            User   `gorm:"foreignkey:UserID"`
	Title               string `gorm:"not null"`
	Description         string `gorm:"not null"`
	JD                  string `gorm:"not null"` // JD stands for Job Description
	MaximumApplications int    `gorm:"not null"`
	// TODO: add more fields related to job posting
}

// CreateJob creates a new job posting in the database
func CreateJob(job *Job) error {
	if err := db.Create(job).Error; err != nil {
		return err
	}
	return nil
}

// GetJobsPaginated retrieves a list of job postings with pagination
func GetJobsPaginated(limit, offset int) ([]Job, error) {
	var jobs []Job
	if err := db.Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJobs retrieves a list of job postings from the database
func GetJobs() ([]Job, error) {
	var jobs []Job
	if err := db.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// DeleteJob deletes a job posting from the database
func DeleteJob(job *Job) error {
	if err := db.Delete(job).Error; err != nil {
		return err
	}
	return nil
}

// GetJobByID retrieves a job posting by its ID from the database
func GetJobByID(jobID int) (*Job, error) {
	var job Job
	if err := db.First(&job, jobID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &job, nil
}

// JobFilter represents filtering criteria for job listings
type JobFilter struct {
	Title       string
	Description string
	// Add other filtering criteria as needed
}

// GetFilteredJobs retrieves a list of job postings filtered by criteria
func GetFilteredJobs(filter JobFilter, limit, offset int) ([]Job, error) {
	var jobs []Job
	db := db

	// Apply filtering criteria
	if filter.Title != "" {
		db = db.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Description != "" {
		db = db.Where("description LIKE ?", "%"+filter.Description+"%")
	}
	// Add other filtering conditions as needed

	// Retrieve filtered jobs
	if err := db.Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// notifyUsersAboutJob notifies all users about the new job posting asynchronously
func NotifyUsersAboutJob(job *Job) {
	users, err := GetAllUsers()
	if err != nil {
		// Handle error
		return
	}

	// Iterate through users and send notifications asynchronously
	for _, user := range users {
		go func(u User) {
			notification := Notification{
				UserID:  u.ID,
				Message: "New job posted: " + job.Title,
				// Add other notification details as needed
				Channels: []string{"email", "sms"},
			}
			notification.SendNotification() // This function is defined in the Notification model
		}(user)
	}
}
