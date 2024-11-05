// utils.go
package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Define a secret key for JWT signing (in a real application, this should be stored securely)
var jwtSecret = []byte("your_secret_key")

// GenerateJWT generates a new JWT token with a 1-hour expiration for a specific user
func GenerateJWT(username string) (string, error) {
	// Define token expiration time
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create claims, including username and expiration time
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create token with claims and sign it using the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates the JWT token in the request header
func ValidateJWT(r *http.Request) bool {
	// Get the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return false
	}

	// Remove the "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Return whether the token is valid or not
	return err == nil && token.Valid
}