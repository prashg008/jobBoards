// handlers/application.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/prashg008/jobBoards/models"
)

// ApplyForJob applies for a job
// @Summary Apply for a job
// @Description Creates a new job application
// @Tags applications
// @Security jwt
// @Accept json
// @Produce json
// @Param application body models.Application true "Application details"
// @Success 201 {object} models.Application "Application created"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Failed to create application"
// @Router /applications [post]
func ApplyForJob(c *gin.Context) {
	var application models.Application

	application.ApplicantID = c.GetInt("userID")
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if application for this job already exists
	existingApplication, err := models.GetApplicationByJobAndUser(application.JobID, application.ApplicantID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing application"})
		return
	}
	if existingApplication != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application for this job already exists"})
		return
	}

	// Save application to database
	if err := models.CreateApplication(&application); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, application)
}

// UpdateApplicationStatus updates application status
// @Summary Update application status
// @Description Updates status of an application
// @Tags applications
// @Security jwt
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param status body models.ApplicationState true "New status"
// @Success 200 {object} ErrorResponse "Status updated"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Failed to update status"
// @Router /applications/{id} [put]
func UpdateApplicationStatus(c *gin.Context) {
	// Parse application ID from request
	applicationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Parse request body to extract updated status
	var status models.ApplicationState
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status update data"})
		return
	}

	// Get the current user from the context (assuming authentication middleware)
	currentUserID := c.GetInt("userID")

	// Check if the current user is the job owner
	application, err := models.GetApplicationByID(uint(applicationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve application"})
		return
	}

	// Check if the current user is the job owner
	job, err := models.GetJobByID(application.JobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}

	if job.PostedByID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this application"})
		return
	}

	// Update the status of the application
	err = models.UpdateApplicationStatus(uint(applicationID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application status"})
		return
	}

	// Send notification to the user about the updated status
	go sendApplicationStatusNotification(application)

	c.JSON(http.StatusOK, gin.H{"message": "Application status updated successfully"})
}

// sendApplicationStatusNotification sends a notification to the user about the updated status of their application
func sendApplicationStatusNotification(application *models.Application) {
	// Implement logic to send notification to the user
	// Example: notification := models.Notification{UserID: userID, Message: "Your application status is now " + status}
	// notification.SendNotification()
}
