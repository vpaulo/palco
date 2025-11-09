package models

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	DueDate     sql.NullTime   `json:"due_date"`
	Archived    bool           `json:"archived"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
