# Palco

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A project management terminal user interface (TUI) built with Bubbletea, Go, and SQLite.

## About

Palco is a terminal-based application for managing projects, tasks, and notes. It provides a structured way to organize your work with priority-based task management, hierarchical subtasks, and integrated note-taking through an intuitive keyboard-driven interface.

## Features

- **Terminal User Interface**:
  - Clean, keyboard-driven interface built with Bubbletea
  - Multi-panel layout for efficient navigation
  - Vim-style keybindings (j/k for navigation)
  - Context-aware help system (press `?`)
- **Project Management**: Create and manage projects with descriptions and due dates
- **Task Organization**:
  - Priority-based task system (None, Low, Medium, High, Urgent)
  - Hierarchical subtasks for breaking down complex tasks
  - Task completion tracking
  - Automatic task description management via linked notes
- **Note Taking**:
  - Project-level notes for general information
  - Task-specific notes and descriptions
  - Automatic note creation when creating tasks with descriptions
  - Context-aware note creation (project or task notes)
- **SQLite Database**:
  - Local-first data storage with `palco.db`
  - WAL (Write-Ahead Logging) mode for better concurrency
  - Automatic migrations on startup
  - Foreign key constraints with cascading deletes

## Project Structure

```
palco/
├── cmd/
│   └── palco/
│       └── main.go        # Application entry point
├── internal/
│   ├── database/          # Database connection and migrations
│   │   ├── db.go          # SQLite connection with WAL mode
│   │   ├── migrate.go     # Migration runner
│   │   └── models/        # Data models (Project, Task, Note)
│   └── repository/        # Data access layer
│       ├── project.go     # Project CRUD operations
│       ├── task.go        # Task CRUD with auto-note creation
│       └── note.go        # Note CRUD operations
├── UI/                    # Bubbletea TUI components
│   ├── model.go           # Main app model and state
│   ├── projects.go        # Projects panel
│   ├── tasks.go           # Tasks panel
│   ├── notes.go           # Notes panel
│   ├── details.go         # Details panel
│   ├── form.go            # Form components
│   ├── help.go            # Help screen
│   └── status_bar.go      # Status bar
├── migrations/            # SQL migration files
│   ├── 001_create_projects_table.up.sql
│   ├── 002_create_tasks_table.up.sql
│   └── 003_create_notes_table.up.sql
└── palco.db              # SQLite database (auto-created)
```
## Getting Started

### Prerequisites
- Go 1.24 or higher

### Installation

Clone the repository and build the application:
```bash
git clone https://github.com/vpaulo/palco.git
cd palco
go build -o palco cmd/palco/main.go
```

### Running

Run the application directly:
```bash
./palco
```

Or run without building:
```bash
go run cmd/palco/main.go
```

### Installing System-wide

To install Palco system-wide:
```bash
go install ./cmd/palco
```

Then run from anywhere:
```bash
palco
```

## Usage

Palco features a multi-panel TUI interface with five main sections:
1. **Projects** - View and manage all projects
2. **Tasks** - View and manage tasks for the selected project
3. **Notes** - View and manage notes for selected projects or tasks
4. **Details** - View detailed information about selected items
5. **Drafts** - Create and manage draft content

### Keybindings

#### Navigation
- `↑/↓` or `k/j` - Move cursor up/down in active section
- `Tab` - Switch to next section
- `Shift+Tab` - Switch to previous section
- `1, 2, 3, 4, 5` - Jump directly to a section (Projects, Tasks, Notes, Details, Drafts)

#### Projects Section
- `n` - Create new project
- `e` - Edit selected project
- `d` - Delete selected project

#### Tasks Section
- `n` - Create new task
- `s` - Create subtask (child of selected task)
- `e` - Edit selected task
- `d` - Delete selected task
- `Space/Enter` - Toggle task completion

#### Notes Section
- `n` - Create new note (project or task note based on context)

#### Forms
- `Tab/Shift+Tab` - Switch between form fields
- `Enter` - Submit form
- `Esc` - Cancel form

#### General
- `?` - Show help screen with all keybindings
- `q` or `Ctrl+C` - Quit application

### Quick Start Guide

1. **Create a Project**: Press `1` to go to Projects, then press `n` to create a new project
2. **Add Tasks**: Press `2` to go to Tasks, then press `n` to create a task for the selected project
3. **Add Notes**: Press `3` to go to Notes, then press `n` to add a note to the selected project or task
4. **View Details**: Press `4` to see full details of the selected item
5. **Get Help**: Press `?` anytime to view the help screen

## License

MIT License - see LICENSE file for details

