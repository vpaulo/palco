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

	return Section(m.activeSection == 1).Width(col1Width).Height(row2Height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Tasks [2]"),
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

		// Determine indentation based on depth level
		depth := 0
		if i < len(m.taskDepths) {
			depth = m.taskDepths[i]
		}

		indent := ""
		prefix := ""
		if depth > 0 {
			// Two spaces per level of depth
			indent = lipgloss.NewStyle().Width(depth * 2).Render("")
			prefix = "└─"  // Tree branch character
		}

		title := task.Title
		maxTitleLen := 23 - (depth * 2) - len(prefix)
		if maxTitleLen < 5 {
			maxTitleLen = 5 // Minimum title length
		}
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}

		// Add completion indicator
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}

		items[i] = fmt.Sprintf("%s %s%s%s %s", cursor, indent, prefix, status, title)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
