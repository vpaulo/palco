package main

import (
	"fmt"
	"os"

	"palco/UI"
	"palco/internal/database"
	"palco/internal/repository"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(init_model(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func init_model() ui.Model {
	db := database.Run()
	return ui.Model{
		Db: db,

		// Initialize repositories
		ProjectRepo: repository.NewProjectRepository(db.DB),
		TaskRepo:    repository.NewTaskRepository(db.DB),
		NoteRepo:    repository.NewNoteRepository(db.DB),
	}
}
