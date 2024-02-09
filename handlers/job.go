package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prashg008/jobBoards/models"
)

// CreateJob creates a new job posting
// @Summary Create a new job posting
// @Description Creates a new job posting by saving to database
// @Tags jobs
// @Security jwt
// @Accept json
// @Produce json
// @Param job body models.Job true "Job details"
// @Success 201 {object} models.Job "Job created"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Failed to create job"
// @Router /jobs [post]
func CreateJob(c *gin.Context) {
	// Check if the user is a company
	role := c.GetString("role")
	if role != "company" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only companies can post jobs"})
		return
	}

	var job models.Job
	job.PostedByID = c.GetInt("userID")
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save job to database
	if err := models.CreateJob(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}
	go models.NotifyUsersAboutJob(&job)

	c.JSON(http.StatusCreated, job)
}

// ListJobs lists all job postings
// @Summary List all job postings
// @Description Retrieves paginated list of job postings
// @Tags jobs
// @Accept json
// @Produce json
// @Param limit query int false "Results per page"
// @Param page query int false "Page number"
// @Param title query string false "Filter by title"
// @Param description query string false "Filter by description"
// @Success 200 {array} models.Job
// @Failure 400 {object} ErrorResponse "Invalid params"
// @Failure 500 {object} ErrorResponse "Failed to get jobs"
// @Router /jobs [get]
func ListJobs(c *gin.Context) {
	// Parse pagination parameters (if needed)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	// Parse filtering parameters
	filter := models.JobFilter{
		Title:       c.Query("title"),
		Description: c.Query("description"),
		// Add other filtering parameters as needed
	}

	// Retrieve list of filtered jobs from the database
	jobs, err := models.GetFilteredJobs(filter, limit, (page-1)*limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// DeleteJob deletes a job posting
// @Summary Delete a job posting
// @Description Deletes a job posting by ID
// @Tags jobs
// @Security jwt
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 204
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Failed to delete job"
// @Router /jobs/{id} [delete]
func DeleteJob(c *gin.Context) {
	// Get job ID from path parameter
	jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	// Check if the job exists
	job, err := models.GetJobByID(int(jobID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Check if the authenticated user is the owner of the job
	userID := uint(c.GetInt("userID"))
	if uint(job.PostedByID) != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this job"})
		return
	}

	// Delete the job
	if err := models.DeleteJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}
