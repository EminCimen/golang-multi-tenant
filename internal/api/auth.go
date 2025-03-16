package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-multi-tenant/internal/database"
	"golang-multi-tenant/internal/middleware"
	"golang-multi-tenant/internal/models"
)

// @Summary     Register a new user
// @Description Register a new user for a specific tenant
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.RegisterRequest true "Registration details"
// @Success     201 {object} map[string]interface{} "User registered successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     409 {object} map[string]string "User already exists"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant database
	tenantDB, err := database.GetTenantDB(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Check if user already exists
	var exists bool
	err = tenantDB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user in tenant database
	var userID int
	err = tenantDB.QueryRow(`
        INSERT INTO users (email, password)
        VALUES ($1, $2)
        RETURNING id`,
		req.Email, hashedPassword).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(userID, req.TenantID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

// @Summary     Login user
// @Description Authenticate a user and return a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.LoginRequest true "Login credentials"
// @Success     200 {object} map[string]interface{} "Login successful"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     401 {object} map[string]string "Invalid credentials"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant database
	tenantDB, err := database.GetTenantDB(req.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var user models.User
	var hashedPassword string
	err = tenantDB.QueryRow(`
        SELECT id, email, password
        FROM users
        WHERE email = $1`,
		req.Email).Scan(&user.ID, &user.Email, &hashedPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !models.CheckPassword(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, req.TenantID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// @Summary     Get user information
// @Description Get current user information based on JWT token
// @Tags        user
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} map[string]interface{} "User information"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Router      /me [get]
func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"tenant_id": tenantID,
		"email":     email,
	})
} 