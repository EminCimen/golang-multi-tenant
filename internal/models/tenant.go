package models

import "time"

// Tenant represents the tenant model
type Tenant struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

// CreateTenantRequest represents the create tenant request body
type CreateTenantRequest struct {
    Name string `json:"name" binding:"required" example:"Example Company"`
} 