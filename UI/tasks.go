package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Tasks(m Model) string {
	col1Width := int(float64(m.width) * 0.40)
	row2Height := int(float64(m.height) * 0.40)

	// Build content
	var content string
	if len(m.projects) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No project selected")
	} else if len(m.tasks) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No tasks found")
	} else {
		content = renderTaskList(m)
	}

	return Section().Width(col1Width).Height(row2Height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Tasks"),
			content,
		),
	)
}

func renderTaskList(m Model) string {
	items := make([]string, len(m.tasks))
	for i, task := range m.tasks {
		cursor := " "
		if i == m.selectedTaskIndex && m.activeSection == 1 {
			cursor = ">"
		}

		title := task.Title
		if len(title) > 23 {
			title = title[:20] + "..."
		}

		// Add completion indicator
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}

		items[i] = fmt.Sprintf("%s %s %s", cursor, status, title)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
