package repository

import (
	"database/sql"
	"fmt"
	"palco/internal/database/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task and optionally a description note
func (r *TaskRepository) Create(projectID int64, parentTaskID *int64, title string, description *string, priority int) (*models.Task, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert task
	taskQuery := `
		INSERT INTO tasks (project_id, parent_task_id, title, priority)
		VALUES (?, ?, ?, ?)
		RETURNING id, project_id, parent_task_id, title, priority, completed, created_at, updated_at
	`

	var task models.Task
	err = tx.QueryRow(taskQuery, projectID, parentTaskID, title, priority).Scan(
		&task.ID,
		&task.ProjectID,
		&task.ParentTaskID,
		&task.Title,
		&task.Priority,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// If description is provided, create a description note
	if description != nil && *description != "" {
		noteQuery := `
			INSERT INTO notes (task_id, content, is_description)
			VALUES (?, ?, 1)
		`
		_, err = tx.Exec(noteQuery, task.ID, *description)
		if err != nil {
			return nil, fmt.Errorf("failed to create description note: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	query := `
		SELECT id, project_id, parent_task_id, title, priority, completed, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	var task models.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.ParentTaskID,
		&task.Title,
		&task.Priority,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// GetByProjectID retrieves all tasks for a project
func (r *TaskRepository) GetByProjectID(projectID int64) ([]models.Task, error) {
	query := `
		SELECT id, project_id, parent_task_id, title, priority, completed, created_at, updated_at
		FROM tasks
		WHERE project_id = ?
		ORDER BY priority DESC, created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.ParentTaskID,
			&task.Title,
			&task.Priority,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetSubtasks retrieves all subtasks for a parent task
func (r *TaskRepository) GetSubtasks(parentTaskID int64) ([]models.Task, error) {
	query := `
		SELECT id, project_id, parent_task_id, title, priority, completed, created_at, updated_at
		FROM tasks
		WHERE parent_task_id = ?
		ORDER BY priority DESC, created_at DESC
	`

	rows, err := r.db.Query(query, parentTaskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.ParentTaskID,
			&task.Title,
			&task.Priority,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(id int64, title string, priority int, completed bool) (*models.Task, error) {
	query := `
		UPDATE tasks
		SET title = ?, priority = ?, completed = ?
		WHERE id = ?
		RETURNING id, project_id, parent_task_id, title, priority, completed, created_at, updated_at
	`

	var task models.Task
	err := r.db.QueryRow(query, title, priority, completed, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.ParentTaskID,
		&task.Title,
		&task.Priority,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &task, nil
}

// Delete deletes a task
func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
