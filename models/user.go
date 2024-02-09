package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/prashg008/jobBoards/database"
)

var db = database.GetDB()

// UserRole represents the role of a user
type UserRole string

const (
	Admin     UserRole = "admin"
	Company   UserRole = "company"
	Applicant UserRole = "applicant"
)

type User struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email     string    `gorm:"unique_index;not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Role      UserRole  `gorm:"not null"`
}

// CreateUser creates a new user in the database
func CreateUser(user *User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(userID int) (*User, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmailOrUsername retrieves a user by email or username from the database
func GetUserByEmailOrUsername(identifier string) (*User, error) {
	var user User
	// Query for the user by email
	if err := db.Where("email = ?", identifier).First(&user).Error; err == nil {
		return &user, nil
	}
	// Query for the user by username (UUID)
	if err := db.Where("uuid = ?", identifier).First(&user).Error; err == nil {
		return &user, nil
	}
	// If user not found by email or username, return error
	return nil, gorm.ErrRecordNotFound
}

// GetApplicationsByUserID retrieves all applications by user ID
func GetApplicationsByUserID(userID uint) ([]Application, error) {
	var applications []Application
	if err := db.Where("user_id = ?", userID).Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}
