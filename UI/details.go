package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Details(m Model) string {
	col2Width := int(float64(m.width) * 0.40)

	// Build content based on active section
	var content string
	if m.activeSection == 0 {
		// Show project details
		content = renderProjectDetails(m)
	} else if m.activeSection == 1 {
		// Show task details
		content = renderTaskDetails(m)
	} else {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("Select a project or task to view details")
	}

	return Section().Width(col2Width).Height(m.height - 3).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Details"),
			content,
		),
	)
}

func renderProjectDetails(m Model) string {
	if len(m.projects) == 0 {
		return lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No project selected")
	}

	project := m.projects[m.selectedProjectIndex]

	// Build the details
	var parts []string

	// Project name (header)
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		MarginBottom(1)
	parts = append(parts, nameStyle.Render(project.Name))

	// Description
	if project.Description.Valid && project.Description.String != "" {
		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(special)

		descStyle := lipgloss.NewStyle().
			MarginBottom(1)

		parts = append(parts, labelStyle.Render("Description:"))

		// Wrap description text
		wrapped := wrapText(project.Description.String, 45)
		parts = append(parts, descStyle.Render(wrapped))
	}

	// Due date
	if project.DueDate.Valid {
		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(special)

		dateStr := project.DueDate.Time.Format("2006-01-02")
		parts = append(parts, labelStyle.Render("Due Date: ")+dateStr)
	}

	// Status
	statusLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		Render("Status: ")

	status := "Active"
	if project.Archived {
		status = "Archived"
	}
	parts = append(parts, statusLabel+status)

	// Created date
	createdLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		Render("Created: ")

	createdStr := project.CreatedAt.Format("2006-01-02")
	parts = append(parts, createdLabel+createdStr)

	contentStyle := lipgloss.NewStyle().Padding(1)
	return contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

func renderTaskDetails(m Model) string {
	if len(m.tasks) == 0 {
		return lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No task selected")
	}

	task := m.tasks[m.selectedTaskIndex]

	// Build the details
	var parts []string

	// Task title (header)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		MarginBottom(1)
	parts = append(parts, titleStyle.Render(task.Title))

	// Status
	statusLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		Render("Status: ")

	status := "[ ] Incomplete"
	if task.Completed {
		status = "[✓] Completed"
	}
	parts = append(parts, statusLabel+status)

	// Priority
	priorityLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		Render("Priority: ")

	priorityText := getPriorityText(task.Priority)
	parts = append(parts, priorityLabel+priorityText)

	// Description (from notes)
	if len(m.notes) > 0 {
		for _, note := range m.notes {
			if note.IsDescription {
				labelStyle := lipgloss.NewStyle().
					Bold(true).
					Foreground(special).
					MarginTop(1)

				descStyle := lipgloss.NewStyle().
					MarginBottom(1)

				parts = append(parts, labelStyle.Render("Description:"))

				// Wrap description text
				wrapped := wrapText(note.Content, 45)
				parts = append(parts, descStyle.Render(wrapped))
				break
			}
		}
	}

	// Other notes
	var otherNotes []string
	for _, note := range m.notes {
		if !note.IsDescription {
			otherNotes = append(otherNotes, note.Content)
		}
	}

	if len(otherNotes) > 0 {
		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(special).
			MarginTop(1)

		parts = append(parts, labelStyle.Render(fmt.Sprintf("Notes (%d):", len(otherNotes))))

		for i, noteContent := range otherNotes {
			noteStyle := lipgloss.NewStyle().
				MarginLeft(2)

			wrapped := wrapText(noteContent, 43)
			noteText := fmt.Sprintf("• %s", wrapped)
			parts = append(parts, noteStyle.Render(noteText))

			if i < len(otherNotes)-1 {
				parts = append(parts, "")
			}
		}
	}

	// Created date
	createdLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		MarginTop(1).
		Render("Created: ")

	createdStr := task.CreatedAt.Format("2006-01-02")
	parts = append(parts, createdLabel+createdStr)

	contentStyle := lipgloss.NewStyle().Padding(1)
	return contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

func getPriorityText(priority int) string {
	switch priority {
	case 0:
		return "None"
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "Urgent"
	default:
		return "Unknown"
	}
}

func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
