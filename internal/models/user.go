package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the user model
type User struct {
    ID        int       `json:"id"`
    TenantID  int       `json:"tenant_id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // "-" means this field won't be included in JSON
    CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
    TenantID int    `json:"tenant_id" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
    TenantID int    `json:"tenant_id" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPassword checks if the provided password matches the hash
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
} 