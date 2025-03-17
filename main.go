package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "golang-multi-tenant/docs" // This will be generated
	"golang-multi-tenant/internal/api"
	"golang-multi-tenant/internal/database"
	"golang-multi-tenant/internal/middleware"
)

// @title           Multi-Tenant API
// @version         1.0
// @description     A multi-tenant API with authentication and user management.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token directly without Bearer prefix
// @Security BearerAuth[]
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	database.InitDB()
	defer func() {
		if err := database.MainDB.Close(); err != nil {
			log.Printf("Error closing main database connection: %v", err)
		}
		// Close all tenant database connections
		for dbName, db := range database.TenantDBs {
			if err := db.Close(); err != nil {
				log.Printf("Error closing tenant database connection %s: %v", dbName, err)
			}
		}
	}()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins not recommended for production
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.POST("/register", api.Register)
	r.POST("/login", api.Login)
	r.POST("/tenants", api.CreateTenant)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", api.Me)
		
		// Post routes
		protected.POST("/posts", api.CreatePost)
		protected.GET("/posts", api.GetPosts)
		protected.GET("/posts/:id", api.GetPost)
	}

	// Start server
	log.Println("Server starting on 0.0.0.0:8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
} 