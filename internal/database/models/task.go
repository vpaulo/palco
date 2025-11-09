package models

import (
	"database/sql"
	"time"
)

// Priority levels for tasks
const (
	PriorityNone   = 0
	PriorityLow    = 1
	PriorityMedium = 2
	PriorityHigh   = 3
	PriorityUrgent = 4
)

type Task struct {
	ID           int64        `json:"id"`
	ProjectID    int64        `json:"project_id"`
	ParentTaskID sql.NullInt64 `json:"parent_task_id"`
	Title        string       `json:"title"`
	Priority     int          `json:"priority"`
	Completed    bool         `json:"completed"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
