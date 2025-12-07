package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Projects(m Model) string {
	col1Width := int(float64(m.width) * 0.40)
	row1Height := int(float64(m.height) * 0.30)

	// Build content
	var content string
	if len(m.projects) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No projects found")
	} else {
		content = renderProjectList(m)
	}

	return Section().Width(col1Width).Height(row1Height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Projects"),
			content,
		),
	)
}

func renderProjectList(m Model) string {
	items := make([]string, len(m.projects))
	for i, project := range m.projects {
		cursor := " "
		if i == m.selectedProjectIndex && m.activeSection == 0 {
			cursor = ">"
		}

		name := project.Name
		if len(name) > 25 {
			name = name[:22] + "..."
		}

		items[i] = fmt.Sprintf("%s %s", cursor, name)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
