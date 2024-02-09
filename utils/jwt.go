package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/spf13/viper"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token
func GenerateJWT(userID uint, username, role string) (string, error) {
	// Load JWT secret key from configuration
	jwtSecret := viper.GetString("jwt.secret_key")
	if jwtSecret == "" {
		return "", errors.New("JWT secret key not found in configuration")
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
