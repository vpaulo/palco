CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER,
    task_id INTEGER,
    content TEXT NOT NULL,
    is_description BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CHECK (
        (project_id IS NOT NULL AND task_id IS NULL) OR
        (project_id IS NULL AND task_id IS NOT NULL)
    )
);

-- Index for faster lookups by project
CREATE INDEX IF NOT EXISTS idx_notes_project_id ON notes(project_id);

-- Index for faster lookups by task
CREATE INDEX IF NOT EXISTS idx_notes_task_id ON notes(task_id);

-- Index for finding task description notes
CREATE INDEX IF NOT EXISTS idx_notes_task_description ON notes(task_id, is_description);

-- Trigger to automatically update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_notes_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
