package models

import (
	"database/sql"
	"time"
)

type Note struct {
	ID            int64        `json:"id"`
	ProjectID     sql.NullInt64 `json:"project_id"`
	TaskID        sql.NullInt64 `json:"task_id"`
	Content       string       `json:"content"`
	IsDescription bool         `json:"is_description"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
