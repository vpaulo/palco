package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Run() *DB {
	// Get database path
	dbPath, err := GetDatabasePath()
	if err != nil {
		log.Fatalf("Failed to get database path: %v", err)
	}

	// Initialize database
	db, err := New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	if err := RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")

	return db
}

// New creates a new database connection
func New(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrency and performance
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetDataDir returns the user data directory path
func GetDataDir() (string, error) {
	// For now, use current directory. Later can be updated to use OS-specific paths
	return ".", nil
}

// GetDatabasePath returns the full path to the database file
func GetDatabasePath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "palco.db"), nil
}
