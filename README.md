# Palco

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A project management platform built with Wails, Go, and SQLite.

## About

Palco is a desktop application for managing projects, tasks, and notes. It provides a structured way to organize your work with priority-based task management, hierarchical subtasks, and integrated note-taking.

## Features

- **Project Management**: Create and manage projects with descriptions and due dates
- **Task Organization**:
  - Priority-based task system (None, Low, Medium, High, Urgent)
  - Hierarchical subtasks for breaking down complex tasks
  - Automatic task description management via linked notes
- **Note Taking**:
  - Project-level notes for general information
  - Task-specific notes and descriptions
  - Automatic note creation when creating tasks with descriptions
- **SQLite Database**:
  - Local-first data storage with `palco.db`
  - WAL (Write-Ahead Logging) mode for better concurrency
  - Automatic migrations on startup
  - Foreign key constraints with cascading deletes

## Project Structure

```
palco/
├── internal/
│   ├── database/          # Database connection and migrations
│   │   ├── db.go          # SQLite connection with WAL mode
│   │   ├── migrate.go     # Migration runner
│   │   └── models/        # Data models (Project, Task, Note)
│   └── repository/        # Data access layer
│       ├── project.go     # Project CRUD operations
│       ├── task.go        # Task CRUD with auto-note creation
│       └── note.go        # Note CRUD operations
├── migrations/            # SQL migration files
│   ├── 001_create_projects_table.up.sql
│   ├── 002_create_tasks_table.up.sql
│   └── 003_create_notes_table.up.sql
├── frontend/              # Frontend application
├── main.go               # Application entry point
└── app.go                # App struct with lifecycle hooks
```
## Getting Started

### Prerequisites
- Go 1.24.2 or higher
- Node.js and npm
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Development

To run in live development mode:
```bash
wails dev
```

This will run a Vite development server with hot reload for frontend changes. A dev server also runs on http://localhost:34115 where you can access your Go methods from the browser devtools.

### Building

To build a production package:
```bash
wails build
```

The built application will be in the `build/bin` directory.

