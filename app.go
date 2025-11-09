package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"palco/internal/database"
	"palco/internal/repository"
)

// App struct
type App struct {
	ctx context.Context
	db  *database.DB

	// Repositories
	projectRepo *repository.ProjectRepository
	taskRepo    *repository.TaskRepository
	noteRepo    *repository.NoteRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) domready(ctx context.Context) {
	// runtime.EventsEmit(a.ctx, "loaded")
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// TODO:
	return false
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Get database path
	dbPath, err := database.GetDatabasePath()
	if err != nil {
		log.Fatalf("Failed to get database path: %v", err)
	}

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	a.db = db

	// Run migrations
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	a.projectRepo = repository.NewProjectRepository(db.DB)
	a.taskRepo = repository.NewTaskRepository(db.DB)
	a.noteRepo = repository.NewNoteRepository(db.DB)

	log.Println("Database initialized successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
