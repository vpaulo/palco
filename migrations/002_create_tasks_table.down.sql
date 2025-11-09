DROP TRIGGER IF EXISTS update_tasks_timestamp;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_parent_task_id;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP TABLE IF EXISTS tasks;
