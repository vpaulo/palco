package ui

import (
	"fmt"
	"palco/internal/database"
	"palco/internal/database/models"
	"palco/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// width = 96

	columnWidth = 30
)

// View modes
const (
	ModeNormal = iota
	ModeCreateProject
	ModeCreateTask
	ModeEditProject
	ModeEditTask
	ModeCreateNote
	ModeHelp
)

// Messages
type projectsLoadedMsg struct {
	projects []models.Project
}

type tasksLoadedMsg struct {
	tasks  []models.Task
	depths []int
}

type notesLoadedMsg struct {
	notes []models.Note
}

type projectCreatedMsg struct {
	project *models.Project
}

type taskCreatedMsg struct {
	task *models.Task
}

type taskUpdatedMsg struct {
	task *models.Task
}

type projectDeletedMsg struct{}

type taskDeletedMsg struct{}

type noteCreatedMsg struct {
	note *models.Note
}

var (
	// General.

	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	base      = lipgloss.NewStyle().Foreground(normal)

	// List.

	list = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle)
	// MarginRight(2).
	// Height(8).
	// Width(columnWidth + 1)

	listHeader = base.
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = base.PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Page.

	docStyle = lipgloss.NewStyle()
)

type Model struct {
	Db *database.DB

	// Repositories
	ProjectRepo *repository.ProjectRepository
	TaskRepo    *repository.TaskRepository
	NoteRepo    *repository.NoteRepository

	// Terminal dimensions
	width  int
	height int

	// State
	projects             []models.Project
	tasks                []models.Task
	taskDepths           []int // Depth level for each task (for indentation)
	notes                []models.Note
	selectedProjectIndex int
	selectedTaskIndex    int
	activeSection        int // 0: projects, 1: tasks, 2: notes, 3: details, 4: drafts
	noteContext          int // 0: project notes, 1: task notes

	// Form state
	mode         int
	formInputs   []textinput.Model
	focusedInput int
	parentTaskID *int64 // Used when creating a subtask
}

func (m Model) Init() tea.Cmd {
	return m.loadProjects
}

// initProjectForm initializes the form for creating a new project
func (m *Model) initProjectForm() {
	m.mode = ModeCreateProject
	m.formInputs = make([]textinput.Model, 2)
	m.focusedInput = 0

	// Name input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Project name"
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 100
	m.formInputs[0].Width = 50

	// Description input
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Description (optional)"
	m.formInputs[1].CharLimit = 500
	m.formInputs[1].Width = 50
}

// initTaskForm initializes the form for creating a new task
func (m *Model) initTaskForm() {
	if len(m.projects) == 0 {
		return
	}

	m.mode = ModeCreateTask
	m.parentTaskID = nil // Creating a top-level task
	m.formInputs = make([]textinput.Model, 3)
	m.focusedInput = 0

	// Title input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Task title"
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 200
	m.formInputs[0].Width = 50

	// Description input
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Description (optional)"
	m.formInputs[1].CharLimit = 1000
	m.formInputs[1].Width = 50

	// Priority input
	m.formInputs[2] = textinput.New()
	m.formInputs[2].Placeholder = "Priority (0=None, 1=Low, 2=Medium, 3=High, 4=Urgent)"
	m.formInputs[2].SetValue("0")
	m.formInputs[2].CharLimit = 1
	m.formInputs[2].Width = 50
}

// initSubtaskForm initializes the form for creating a subtask
func (m *Model) initSubtaskForm() {
	if len(m.tasks) == 0 || len(m.projects) == 0 {
		return
	}

	// Set parent task ID to the currently selected task
	taskID := m.tasks[m.selectedTaskIndex].ID
	m.parentTaskID = &taskID

	m.mode = ModeCreateTask
	m.formInputs = make([]textinput.Model, 3)
	m.focusedInput = 0

	// Title input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Subtask title"
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 200
	m.formInputs[0].Width = 50

	// Description input
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Description (optional)"
	m.formInputs[1].CharLimit = 1000
	m.formInputs[1].Width = 50

	// Priority input
	m.formInputs[2] = textinput.New()
	m.formInputs[2].Placeholder = "Priority (0=None, 1=Low, 2=Medium, 3=High, 4=Urgent)"
	m.formInputs[2].SetValue("0")
	m.formInputs[2].CharLimit = 1
	m.formInputs[2].Width = 50
}

// createProject creates a new project from form inputs
func (m Model) createProject() tea.Msg {
	name := m.formInputs[0].Value()
	if name == "" {
		return nil
	}

	var description *string
	if desc := m.formInputs[1].Value(); desc != "" {
		description = &desc
	}

	project, err := m.ProjectRepo.Create(name, description, nil)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return projectCreatedMsg{project: project}
}

// createTask creates a new task from form inputs
func (m Model) createTask() tea.Msg {
	title := m.formInputs[0].Value()
	if title == "" {
		return nil
	}

	if len(m.projects) == 0 {
		return nil
	}

	projectID := m.projects[m.selectedProjectIndex].ID

	var description *string
	if desc := m.formInputs[1].Value(); desc != "" {
		description = &desc
	}

	// Parse priority from form input (default to 0 if invalid)
	priority := 0
	if priorityStr := m.formInputs[2].Value(); priorityStr != "" {
		if p, err := strconv.Atoi(priorityStr); err == nil && p >= 0 && p <= 4 {
			priority = p
		}
	}

	// Use parentTaskID if creating a subtask, otherwise nil for top-level task
	task, err := m.TaskRepo.Create(projectID, m.parentTaskID, title, description, priority)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return taskCreatedMsg{task: task}
}

// initEditProjectForm initializes the form for editing the selected project
func (m *Model) initEditProjectForm() {
	if len(m.projects) == 0 || m.selectedProjectIndex >= len(m.projects) {
		return
	}

	project := m.projects[m.selectedProjectIndex]

	m.mode = ModeEditProject
	m.formInputs = make([]textinput.Model, 2)
	m.focusedInput = 0

	// Name input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Project name"
	m.formInputs[0].SetValue(project.Name)
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 100
	m.formInputs[0].Width = 50

	// Description input
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Description (optional)"
	if project.Description.Valid {
		m.formInputs[1].SetValue(project.Description.String)
	}
	m.formInputs[1].CharLimit = 500
	m.formInputs[1].Width = 50
}

// initNoteForm initializes the form for creating a new note
func (m *Model) initNoteForm() {
	// Check if we have a project or task to attach the note to
	hasProject := len(m.projects) > 0
	hasTask := len(m.tasks) > 0

	// Determine context based on current noteContext (which is set when notes are loaded)
	// Override if we're in the Projects or Tasks section directly
	if m.activeSection == 0 {
		// In Projects section - create project note
		if !hasProject {
			return
		}
		m.noteContext = 0
	} else if m.activeSection == 1 {
		// In Tasks section - create task note
		if !hasTask {
			return
		}
		m.noteContext = 1
	} else {
		// In Notes section or other - use current noteContext
		// Validate we have the right entity
		if m.noteContext == 0 && !hasProject {
			return
		}
		if m.noteContext == 1 && !hasTask {
			return
		}
	}

	m.mode = ModeCreateNote
	m.formInputs = make([]textinput.Model, 1)
	m.focusedInput = 0

	// Content input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Note content"
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 1000
	m.formInputs[0].Width = 50
}

// initEditTaskForm initializes the form for editing the selected task
func (m *Model) initEditTaskForm() {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return
	}

	task := m.tasks[m.selectedTaskIndex]

	m.mode = ModeEditTask
	m.formInputs = make([]textinput.Model, 3)
	m.focusedInput = 0

	// Title input
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Task title"
	m.formInputs[0].SetValue(task.Title)
	m.formInputs[0].Focus()
	m.formInputs[0].CharLimit = 200
	m.formInputs[0].Width = 50

	// Description input - need to fetch from notes
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Description (optional)"

	// Find description note
	for _, note := range m.notes {
		if note.IsDescription {
			m.formInputs[1].SetValue(note.Content)
			break
		}
	}

	m.formInputs[1].CharLimit = 1000
	m.formInputs[1].Width = 50

	// Priority input
	m.formInputs[2] = textinput.New()
	m.formInputs[2].Placeholder = "Priority (0=None, 1=Low, 2=Medium, 3=High, 4=Urgent)"
	m.formInputs[2].SetValue(fmt.Sprintf("%d", task.Priority))
	m.formInputs[2].CharLimit = 1
	m.formInputs[2].Width = 50
}

// updateProject updates the selected project from form inputs
func (m Model) updateProject() tea.Msg {
	if len(m.projects) == 0 || m.selectedProjectIndex >= len(m.projects) {
		return nil
	}

	name := m.formInputs[0].Value()
	if name == "" {
		return nil
	}

	project := m.projects[m.selectedProjectIndex]

	var description *string
	if desc := m.formInputs[1].Value(); desc != "" {
		description = &desc
	}

	updatedProject, err := m.ProjectRepo.Update(project.ID, name, description, nil)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return projectCreatedMsg{project: updatedProject}
}

// createNote creates a new note from form inputs
func (m Model) createNote() tea.Msg {
	content := m.formInputs[0].Value()
	if content == "" {
		return nil
	}

	var note *models.Note
	var err error

	if m.noteContext == 0 {
		// Create project note
		if len(m.projects) == 0 {
			return nil
		}
		projectID := m.projects[m.selectedProjectIndex].ID
		note, err = m.NoteRepo.CreateForProject(projectID, content)
	} else {
		// Create task note
		if len(m.tasks) == 0 {
			return nil
		}
		taskID := m.tasks[m.selectedTaskIndex].ID
		note, err = m.NoteRepo.CreateForTask(taskID, content)
	}

	if err != nil {
		// TODO: Handle error
		return nil
	}

	return noteCreatedMsg{note: note}
}

// updateTask updates the selected task from form inputs
func (m Model) updateTask() tea.Msg {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return nil
	}

	title := m.formInputs[0].Value()
	if title == "" {
		return nil
	}

	task := m.tasks[m.selectedTaskIndex]

	// Parse priority from form input (default to current priority if invalid)
	priority := task.Priority
	if priorityStr := m.formInputs[2].Value(); priorityStr != "" {
		if p, err := strconv.Atoi(priorityStr); err == nil && p >= 0 && p <= 4 {
			priority = p
		}
	}

	// Update task
	updatedTask, err := m.TaskRepo.Update(task.ID, title, priority, task.Completed)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	// Update description note if provided
	if desc := m.formInputs[1].Value(); desc != "" {
		// Check if description note exists
		hasDescription := false
		for _, note := range m.notes {
			if note.IsDescription {
				hasDescription = true
				break
			}
		}

		if hasDescription {
			// Update existing description
			err = m.NoteRepo.UpdateTaskDescription(task.ID, desc)
		} else {
			// Create new description note
			_, err = m.NoteRepo.CreateForTask(task.ID, desc)
		}

		if err != nil {
			// TODO: Handle error
		}
	}

	return taskUpdatedMsg{task: updatedTask}
}

// toggleTaskCompletion toggles the completion status of the selected task
func (m Model) toggleTaskCompletion() tea.Msg {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return nil
	}

	task := m.tasks[m.selectedTaskIndex]

	// Toggle completion
	newCompleted := !task.Completed

	updatedTask, err := m.TaskRepo.Update(task.ID, task.Title, task.Priority, newCompleted)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return taskUpdatedMsg{task: updatedTask}
}

// deleteProject deletes the currently selected project
func (m Model) deleteProject() tea.Msg {
	if len(m.projects) == 0 || m.selectedProjectIndex >= len(m.projects) {
		return nil
	}

	project := m.projects[m.selectedProjectIndex]
	err := m.ProjectRepo.Delete(project.ID)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return projectDeletedMsg{}
}

// deleteTask deletes the currently selected task
func (m Model) deleteTask() tea.Msg {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return nil
	}

	task := m.tasks[m.selectedTaskIndex]
	err := m.TaskRepo.Delete(task.ID)
	if err != nil {
		// TODO: Handle error
		return nil
	}

	return taskDeletedMsg{}
}

// loadProjects loads all active projects from the database
func (m Model) loadProjects() tea.Msg {
	projects, err := m.ProjectRepo.GetAllActive()
	if err != nil {
		// For now, return empty slice on error
		// TODO: Add error handling
		return projectsLoadedMsg{projects: []models.Project{}}
	}
	return projectsLoadedMsg{projects: projects}
}

// loadTasks loads tasks for the currently selected project
func (m Model) loadTasks() tea.Msg {
	if len(m.projects) == 0 || m.selectedProjectIndex >= len(m.projects) {
		return tasksLoadedMsg{tasks: []models.Task{}}
	}

	projectID := m.projects[m.selectedProjectIndex].ID
	tasks, err := m.TaskRepo.GetByProjectID(projectID)
	if err != nil {
		// For now, return empty slice on error
		// TODO: Add error handling
		return tasksLoadedMsg{tasks: []models.Task{}}
	}

	// Organize tasks hierarchically (parents followed by their children, recursively)
	hierarchicalTasks, depths := organizeTasksHierarchically(tasks)

	return tasksLoadedMsg{tasks: hierarchicalTasks, depths: depths}
}

// organizeTasksHierarchically reorganizes tasks so subtasks appear under their parents (recursively)
func organizeTasksHierarchically(tasks []models.Task) ([]models.Task, []int) {
	if len(tasks) == 0 {
		return tasks, []int{}
	}

	// Build a map of parent ID to children
	subtasksByParent := make(map[int64][]models.Task)
	var rootTasks []models.Task

	for _, task := range tasks {
		if task.ParentTaskID.Valid {
			parentID := task.ParentTaskID.Int64
			subtasksByParent[parentID] = append(subtasksByParent[parentID], task)
		} else {
			rootTasks = append(rootTasks, task)
		}
	}

	// Recursively build the hierarchical list
	var result []models.Task
	var depths []int

	var addTaskAndChildren func(task models.Task, depth int)
	addTaskAndChildren = func(task models.Task, depth int) {
		result = append(result, task)
		depths = append(depths, depth)

		// Add children recursively
		if children, exists := subtasksByParent[task.ID]; exists {
			for _, child := range children {
				addTaskAndChildren(child, depth+1)
			}
		}
	}

	// Add all root tasks and their descendants
	for _, rootTask := range rootTasks {
		addTaskAndChildren(rootTask, 0)
	}

	return result, depths
}

// loadNotes loads notes for the currently selected task
func (m Model) loadNotes() tea.Msg {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return notesLoadedMsg{notes: []models.Note{}}
	}

	taskID := m.tasks[m.selectedTaskIndex].ID
	notes, err := m.NoteRepo.GetByTaskID(taskID)
	if err != nil {
		// For now, return empty slice on error
		// TODO: Add error handling
		return notesLoadedMsg{notes: []models.Note{}}
	}
	m.noteContext = 1 // Task notes context
	return notesLoadedMsg{notes: notes}
}

// loadProjectNotes loads notes for the currently selected project
func (m Model) loadProjectNotes() tea.Msg {
	if len(m.projects) == 0 || m.selectedProjectIndex >= len(m.projects) {
		return notesLoadedMsg{notes: []models.Note{}}
	}

	projectID := m.projects[m.selectedProjectIndex].ID
	notes, err := m.NoteRepo.GetByProjectID(projectID)
	if err != nil {
		// For now, return empty slice on error
		// TODO: Add error handling
		return notesLoadedMsg{notes: []models.Note{}}
	}
	m.noteContext = 0 // Project notes context
	return notesLoadedMsg{notes: notes}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle projects loaded
	case projectsLoadedMsg:
		m.projects = msg.projects
		if len(m.projects) > 0 {
			m.selectedProjectIndex = 0
			return m, tea.Batch(m.loadTasks, m.loadProjectNotes)
		}
		return m, nil

	// Handle tasks loaded
	case tasksLoadedMsg:
		m.tasks = msg.tasks
		m.taskDepths = msg.depths
		m.selectedTaskIndex = 0
		if len(m.tasks) > 0 {
			return m, m.loadNotes
		}
		m.notes = []models.Note{}
		return m, nil

	// Handle notes loaded
	case notesLoadedMsg:
		m.notes = msg.notes
		return m, nil

	// Handle project created
	case projectCreatedMsg:
		m.mode = ModeNormal
		m.formInputs = nil
		return m, m.loadProjects

	// Handle task created
	case taskCreatedMsg:
		m.mode = ModeNormal
		m.formInputs = nil
		return m, m.loadTasks

	// Handle task updated
	case taskUpdatedMsg:
		m.mode = ModeNormal
		m.formInputs = nil
		return m, m.loadTasks

	// Handle project deleted
	case projectDeletedMsg:
		m.selectedProjectIndex = 0
		return m, m.loadProjects

	// Handle task deleted
	case taskDeletedMsg:
		m.selectedTaskIndex = 0
		return m, m.loadTasks

	// Handle note created
	case noteCreatedMsg:
		m.mode = ModeNormal
		m.formInputs = nil
		if m.noteContext == 0 {
			return m, m.loadProjectNotes
		}
		return m, m.loadNotes

	// Handle window resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		// Handle help screen
		if m.mode == ModeHelp {
			// Any key exits help
			m.mode = ModeNormal
			return m, nil
		}

		// Handle form inputs
		if m.mode == ModeCreateProject || m.mode == ModeCreateTask || m.mode == ModeEditProject || m.mode == ModeEditTask || m.mode == ModeCreateNote {
			switch msg.String() {
			case "esc":
				m.mode = ModeNormal
				m.formInputs = nil
				return m, nil

			case "enter":
				// Submit form
				if m.mode == ModeCreateProject {
					return m, m.createProject
				} else if m.mode == ModeCreateTask {
					return m, m.createTask
				} else if m.mode == ModeEditProject {
					return m, m.updateProject
				} else if m.mode == ModeEditTask {
					return m, m.updateTask
				} else if m.mode == ModeCreateNote {
					return m, m.createNote
				}

			case "tab", "down":
				// Move to next input
				m.formInputs[m.focusedInput].Blur()
				m.focusedInput = (m.focusedInput + 1) % len(m.formInputs)
				m.formInputs[m.focusedInput].Focus()
				return m, nil

			case "shift+tab", "up":
				// Move to previous input
				m.formInputs[m.focusedInput].Blur()
				m.focusedInput = (m.focusedInput - 1 + len(m.formInputs)) % len(m.formInputs)
				m.formInputs[m.focusedInput].Focus()
				return m, nil

			default:
				// Update the focused input
				var cmd tea.Cmd
				m.formInputs[m.focusedInput], cmd = m.formInputs[m.focusedInput].Update(msg)
				return m, cmd
			}
		}

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			m.Db.Close()
			return m, tea.Quit

		// Navigation
		case "up", "k":
			if m.activeSection == 0 && len(m.projects) > 0 {
				// Navigate projects
				if m.selectedProjectIndex > 0 {
					m.selectedProjectIndex--
					return m, tea.Batch(m.loadTasks, m.loadProjectNotes)
				}
			} else if m.activeSection == 1 && len(m.tasks) > 0 {
				// Navigate tasks
				if m.selectedTaskIndex > 0 {
					m.selectedTaskIndex--
					return m, m.loadNotes
				}
			}

		case "down", "j":
			if m.activeSection == 0 && len(m.projects) > 0 {
				// Navigate projects
				if m.selectedProjectIndex < len(m.projects)-1 {
					m.selectedProjectIndex++
					return m, tea.Batch(m.loadTasks, m.loadProjectNotes)
				}
			} else if m.activeSection == 1 && len(m.tasks) > 0 {
				// Navigate tasks
				if m.selectedTaskIndex < len(m.tasks)-1 {
					m.selectedTaskIndex++
					return m, m.loadNotes
				}
			}

		// Switch active section
		case "tab":
			m.activeSection = (m.activeSection + 1) % 5

		case "shift+tab":
			m.activeSection = (m.activeSection - 1 + 5) % 5

		// Create new item
		case "n":
			if m.activeSection == 0 {
				// Create new project
				m.initProjectForm()
			} else if m.activeSection == 1 {
				// Create new task
				m.initTaskForm()
			} else if m.activeSection == 2 {
				// Create new note
				m.initNoteForm()
			}
			return m, nil

		// Toggle task completion
		case " ", "enter":
			if m.activeSection == 1 && len(m.tasks) > 0 {
				return m, m.toggleTaskCompletion
			}

		// Delete item
		case "d":
			if m.activeSection == 0 && len(m.projects) > 0 {
				// Delete project
				return m, m.deleteProject
			} else if m.activeSection == 1 && len(m.tasks) > 0 {
				// Delete task
				return m, m.deleteTask
			}

		// Edit item
		case "e":
			if m.activeSection == 0 && len(m.projects) > 0 {
				// Edit project
				m.initEditProjectForm()
			} else if m.activeSection == 1 && len(m.tasks) > 0 {
				// Edit task
				m.initEditTaskForm()
			}
			return m, nil

		// Create subtask
		case "s":
			if m.activeSection == 1 && len(m.tasks) > 0 {
				// Create subtask for selected task
				m.initSubtaskForm()
			}
			return m, nil

		// Show help
		case "?":
			m.mode = ModeHelp
			return m, nil

		// Direct section navigation
		case "1":
			m.activeSection = 0 // Projects
			return m, nil
		case "2":
			m.activeSection = 1 // Tasks
			return m, nil
		case "3":
			m.activeSection = 2 // Notes
			return m, nil
		case "4":
			m.activeSection = 3 // Details
			return m, nil
		case "5":
			m.activeSection = 4 // Drafts
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	// Calculate column widths as percentages
	// col1Width := int(float64(m.width) * 0.0)
	// col2Width := int(float64(m.width) * 0.40)
	// // Make the third column take remaining width to avoid rounding gaps
	// // 6 is for the 6 borders
	// col3Width := m.width - col1Width - col2Width - 6

	// Build lists
	// lists := lipgloss.JoinHorizontal(lipgloss.Top,
	// 	list.Width(col1Width).Render(
	// 		lipgloss.JoinVertical(lipgloss.Left,
	// 			listHeader("List 1"),
	// 			listDone("Grapefruit"),
	// 			listDone("Yuzu"),
	// 			listItem("Citron"),
	// 			listItem("Kumquat"),
	// 			listItem("Pomelo"),
	// 		),
	// 	),
	// 	list.Width(col2Width).Render(
	// 		lipgloss.JoinVertical(lipgloss.Left,
	// 			listHeader("List 2"),
	// 			listItem("Glossier"),
	// 			listItem("Claire's Boutique"),
	// 			listDone("Nyx"),
	// 			listItem("Mac"),
	// 			listDone("Milk"),
	// 		),
	// 	),
	// 	list.Width(col3Width).Render(
	// 		lipgloss.JoinVertical(lipgloss.Left,
	// 			listHeader("List 3"),
	// 			listItem("Glossier"),
	// 			listItem("Claire's Boutique"),
	// 			listDone("Nyx"),
	// 			listItem("Mac"),
	// 			listDone("Milk"),
	// 		),
	// 	),
	// )

	// Build status bar
	statusBar := StatusBar(m)

	// Section/Container to span height
	content := base.Width(m.width).Height(m.height - lipgloss.Height(statusBar)).Render(Grid(m))

	// Combine content and status bar
	board := lipgloss.JoinVertical(lipgloss.Left, content, statusBar)

	// Place content with status bar at bottom
	mainView := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Bottom,
		board,
	)

	// If in help mode, show help screen
	if m.mode == ModeHelp {
		return RenderHelp(m)
	}

	// If in form mode, overlay the form
	if m.mode != ModeNormal {
		return RenderForm(m)
	}

	return mainView
}
