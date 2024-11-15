package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var DB *sql.DB

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	log.Println("db.go.InitDB()::Init --------------------------------")
	// Read database connection parameters from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Set the PostgreSQL connection string using environment variables
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("db.go.InitDB()::Could not connect to PostgreSQL database: " + err.Error())
	}

	// Set connection pool
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Run migrations
	createMigrationTable()
	runMigrations()
}

// createMigrationTable creates a table to track applied migrations
func createMigrationTable() {
	log.Println("db.go.createMigrationTable()::Creating migration table")

	// Migration table
	var sql string = `
	CREATE TABLE IF NOT EXISTS migrations (
      migration_id SERIAL PRIMARY KEY,
      migration_module VARCHAR(255),
      migration_name VARCHAR(255) NOT NULL,
      migration_run_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
	`

	_, err := DB.Exec(sql)
	if err != nil {
		panic("db.go.createMigrationTable()::Error could not create migration table: " + err.Error())
	}
}

// runMigrations applies migrations from the migrations folder
func runMigrations() {
	log.Println("db.go.runMigrations()::Init")

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Could not get the current working directory: " + err.Error())
	}

	// Construct the path to db/migrations based on the current working directory
	migrationsPath := filepath.Join(wd, "db", "migrations")

	// Check if the migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Println("db/migrations directory does not exist, skipping migrations.")
		return
	}

	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		panic("db.go.runMigrations()::Error reading migrations directory: " + err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationName := file.Name()
			if !isMigrationApplied(migrationName) {
				err := applyMigration(migrationName)
				if err != nil {
					log.Printf("db.go.runMigrations()::Error applying migration %s: %v", migrationName, err)
					panic(err)
				}
				logMigration(migrationName)
			}
		}
	}
}

// isMigrationApplied checks if a migration has already been applied
func isMigrationApplied(migrationName string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM migrations WHERE migration_name = $1)`
	err := DB.QueryRow(query, migrationName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("db.go.isMigrationApplied()::Error checking migration %s: %v", migrationName, err)
	}
	return exists
}

// applyMigration applies a migration to the database
func applyMigration(migrationName string) error {
	log.Printf("db.go.applyMigration()::Applying migration %s", migrationName)
	content, err := os.ReadFile(filepath.Join("db/migrations", migrationName))
	if err != nil {
		return err
	}
	_, err = DB.Exec(string(content))
	return err
}

// logMigration logs an applied migration in the migrations table
func logMigration(migrationName string) {
	sql := `INSERT INTO migrations (migration_module, migration_name) VALUES ($1, $2)`
	_, err := DB.Exec(sql, "default", migrationName)
	if err != nil {
		log.Fatalf("db.go.logMigration()::Error logging migration %s: %v", migrationName, err)
	}
	log.Printf("db.go.logMigration()::Migration %s logged successfully", migrationName)
}
