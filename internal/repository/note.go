package repository

import (
	"database/sql"
	"fmt"
	"palco/internal/database/models"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// CreateForProject creates a new note for a project
func (r *NoteRepository) CreateForProject(projectID int64, content string) (*models.Note, error) {
	query := `
		INSERT INTO notes (project_id, content)
		VALUES (?, ?)
		RETURNING id, project_id, task_id, content, is_description, created_at, updated_at
	`

	var note models.Note
	err := r.db.QueryRow(query, projectID, content).Scan(
		&note.ID,
		&note.ProjectID,
		&note.TaskID,
		&note.Content,
		&note.IsDescription,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project note: %w", err)
	}

	return &note, nil
}

// CreateForTask creates a new note for a task
func (r *NoteRepository) CreateForTask(taskID int64, content string) (*models.Note, error) {
	query := `
		INSERT INTO notes (task_id, content)
		VALUES (?, ?)
		RETURNING id, project_id, task_id, content, is_description, created_at, updated_at
	`

	var note models.Note
	err := r.db.QueryRow(query, taskID, content).Scan(
		&note.ID,
		&note.ProjectID,
		&note.TaskID,
		&note.Content,
		&note.IsDescription,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task note: %w", err)
	}

	return &note, nil
}

// GetByProjectID retrieves all notes for a project
func (r *NoteRepository) GetByProjectID(projectID int64) ([]models.Note, error) {
	query := `
		SELECT id, project_id, task_id, content, is_description, created_at, updated_at
		FROM notes
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.ProjectID,
			&note.TaskID,
			&note.Content,
			&note.IsDescription,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetByTaskID retrieves all notes for a task
func (r *NoteRepository) GetByTaskID(taskID int64) ([]models.Note, error) {
	query := `
		SELECT id, project_id, task_id, content, is_description, created_at, updated_at
		FROM notes
		WHERE task_id = ?
		ORDER BY is_description DESC, created_at DESC
	`

	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.ProjectID,
			&note.TaskID,
			&note.Content,
			&note.IsDescription,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetTaskDescription retrieves the description note for a task
func (r *NoteRepository) GetTaskDescription(taskID int64) (*models.Note, error) {
	query := `
		SELECT id, project_id, task_id, content, is_description, created_at, updated_at
		FROM notes
		WHERE task_id = ? AND is_description = 1
		LIMIT 1
	`

	var note models.Note
	err := r.db.QueryRow(query, taskID).Scan(
		&note.ID,
		&note.ProjectID,
		&note.TaskID,
		&note.Content,
		&note.IsDescription,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No description found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task description: %w", err)
	}

	return &note, nil
}

// UpdateTaskDescription updates the description note for a task
func (r *NoteRepository) UpdateTaskDescription(taskID int64, content string) error {
	query := `
		UPDATE notes
		SET content = ?
		WHERE task_id = ? AND is_description = 1
	`

	result, err := r.db.Exec(query, content, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task description: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task description not found")
	}

	return nil
}

// Update updates a note
func (r *NoteRepository) Update(id int64, content string) (*models.Note, error) {
	query := `
		UPDATE notes
		SET content = ?
		WHERE id = ?
		RETURNING id, project_id, task_id, content, is_description, created_at, updated_at
	`

	var note models.Note
	err := r.db.QueryRow(query, content, id).Scan(
		&note.ID,
		&note.ProjectID,
		&note.TaskID,
		&note.Content,
		&note.IsDescription,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return &note, nil
}

// Delete deletes a note
func (r *NoteRepository) Delete(id int64) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
