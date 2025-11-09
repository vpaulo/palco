package repository

import (
	"database/sql"
	"fmt"
	"palco/internal/database/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(name string, description *string, dueDate *string) (*models.Project, error) {
	query := `
		INSERT INTO projects (name, description, due_date)
		VALUES (?, ?, ?)
		RETURNING id, name, description, due_date, archived, created_at, updated_at
	`

	var project models.Project
	err := r.db.QueryRow(query, name, description, dueDate).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.DueDate,
		&project.Archived,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &project, nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id int64) (*models.Project, error) {
	query := `
		SELECT id, name, description, due_date, archived, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	var project models.Project
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.DueDate,
		&project.Archived,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	query := `
		SELECT id, name, description, due_date, archived, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DueDate,
			&project.Archived,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// Update updates a project
func (r *ProjectRepository) Update(id int64, name string, description *string, dueDate *string) (*models.Project, error) {
	query := `
		UPDATE projects
		SET name = ?, description = ?, due_date = ?
		WHERE id = ?
		RETURNING id, name, description, due_date, archived, created_at, updated_at
	`

	var project models.Project
	err := r.db.QueryRow(query, name, description, dueDate, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.DueDate,
		&project.Archived,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &project, nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id int64) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

// Archive archives a project
func (r *ProjectRepository) Archive(id int64) error {
	query := `UPDATE projects SET archived = 1 WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

// Unarchive unarchives a project
func (r *ProjectRepository) Unarchive(id int64) error {
	query := `UPDATE projects SET archived = 0 WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to unarchive project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

// GetAllActive retrieves all active (non-archived) projects
func (r *ProjectRepository) GetAllActive() ([]models.Project, error) {
	query := `
		SELECT id, name, description, due_date, archived, created_at, updated_at
		FROM projects
		WHERE archived = 0
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DueDate,
			&project.Archived,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// GetAllArchived retrieves all archived projects
func (r *ProjectRepository) GetAllArchived() ([]models.Project, error) {
	query := `
		SELECT id, name, description, due_date, archived, created_at, updated_at
		FROM projects
		WHERE archived = 1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get archived projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.DueDate,
			&project.Archived,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}
