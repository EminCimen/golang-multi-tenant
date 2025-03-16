package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	MainDB *sql.DB
	TenantDBs map[string]*sql.DB = make(map[string]*sql.DB)
)

// InitDB initializes the main database connection
func InitDB() {
	mainConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	var err error
	MainDB, err = sql.Open("postgres", mainConnStr)
	if err != nil {
		log.Fatal("Error connecting to main database:", err)
	}

	err = MainDB.Ping()
	if err != nil {
		log.Fatal("Error pinging main database:", err)
	}

	// Create main tenant management database
	_, err = MainDB.Exec("CREATE DATABASE tenant_management")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatal("Error creating tenant management database:", err)
	}

	// Connect to tenant management database
	managementConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=tenant_management sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	MainDB, err = sql.Open("postgres", managementConnStr)
	if err != nil {
		log.Fatal("Error connecting to tenant management database:", err)
	}

	// Create necessary tables in tenant management database
	createMainTables()
}

// createMainTables creates tables in the main tenant management database
func createMainTables() {
	// Create tenants table
	_, err := MainDB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			db_name VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating tenants table:", err)
	}
}

// GetTenantDB gets or creates a connection to a tenant's database
func GetTenantDB(tenantID int) (*sql.DB, error) {
	// Get tenant info from main database
	var dbName string
	err := MainDB.QueryRow("SELECT db_name FROM tenants WHERE id = $1", tenantID).Scan(&dbName)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %v", err)
	}

	// Check if we already have a connection
	if db, exists := TenantDBs[dbName]; exists {
		return db, nil
	}

	// Create new connection
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to tenant database: %v", err)
	}

	TenantDBs[dbName] = db
	return db, nil
}

// CreateTenantDB creates a new database for a tenant
func CreateTenantDB(tenantName string) (string, error) {
	// Generate database name
	dbName := fmt.Sprintf("tenant_%s", strings.ToLower(strings.ReplaceAll(tenantName, " ", "_")))

	// Create new database
	_, err := MainDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return "", fmt.Errorf("error creating tenant database: %v", err)
	}

	// Connect to new database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to new tenant database: %v", err)
	}

	// Create tenant-specific tables
	err = createTenantTables(db)
	if err != nil {
		return "", fmt.Errorf("error creating tenant tables: %v", err)
	}

	TenantDBs[dbName] = db
	return dbName, nil
}

// createTenantTables creates necessary tables in a tenant's database
func createTenantTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id),
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
} 