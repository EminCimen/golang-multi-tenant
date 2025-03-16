package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-multi-tenant/internal/database"
	"golang-multi-tenant/internal/models"
)

// @Summary     Create a new tenant
// @Description Create a new tenant in the system and set up its database
// @Tags        tenant
// @Accept      json
// @Produce     json
// @Param       request body models.CreateTenantRequest true "Tenant details"
// @Success     201 {object} models.Tenant "Tenant created successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     409 {object} map[string]string "Tenant already exists"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /tenants [post]
func CreateTenant(c *gin.Context) {
	var req models.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if tenant with same name exists
	var exists bool
	err := database.MainDB.QueryRow("SELECT EXISTS(SELECT 1 FROM tenants WHERE name = $1)", req.Name).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Tenant with this name already exists"})
		return
	}

	// Create tenant database
	dbName, err := database.CreateTenantDB(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating tenant database: " + err.Error()})
		return
	}

	// Create tenant record in management database
	var tenant models.Tenant
	err = database.MainDB.QueryRow(`
        INSERT INTO tenants (name, db_name)
        VALUES ($1, $2)
        RETURNING id, name, created_at`,
		req.Name, dbName).Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating tenant record"})
		return
	}

	c.JSON(http.StatusCreated, tenant)
} 