// handlers/user.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prashg008/jobBoards/models"
	"github.com/prashg008/jobBoards/utils"
)

type Credentials struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterUser registers a new user
// @Summary Register a new user
// @Description Registers a new user by saving them to the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Success 201 {object} models.User "User created"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Failed to create user"
// @Router /users/register [post]
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password before saving to database
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Save user to database
	if err := models.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// LoginUser logs in a user
// @Summary Log in a user
// @Description Logs in a user by verifying credentials and returning a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Login credentials"
// @Success 200 {object} ErrorResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Failed to generate token"
// @Router /users/login [post]
func LoginUser(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user from database by email or username
	user, err := models.GetUserByEmailOrUsername(credentials.Identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserDetails retrieves user details
// @Summary Get user details
// @Description Retrieves details for the currently authenticated user
// @Tags users
// @Security jwt
// @Accept json
// @Produce json
// @Success 200 {object} models.User "User details"
// @Failure 500 {object} ErrorResponse "Failed to get user"
// @Router /users/me [get]
func GetUserDetails(c *gin.Context) {
	userID := c.GetInt("userID")

	// Retrieve user from database
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListApplicationsByUser lists applications by user ID
// @Summary List applications for a user
// @Description Retrieves all job applications for the specified user ID
// @Tags applications
// @Security jwt
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.Application
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 500 {object} ErrorResponse "Failed to retrieve applications"
// @Router /users/{id}/applications [get]
func ListApplicationsByUser(c *gin.Context) {
	// Parse user ID from request
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get applications by user ID
	applications, err := models.GetApplicationsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}
