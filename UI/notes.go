package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Notes(m Model) string {
	col1Width := int(float64(m.width) * 0.40)
	row1Height := int(float64(m.height) * 0.30)

	// Build content
	var content string
	if len(m.tasks) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No task selected")
	} else if len(m.notes) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No notes found")
	} else {
		content = renderNotesList(m)
	}

	return Section().Width(col1Width).Height(row1Height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Notes"),
			content,
		),
	)
}

func renderNotesList(m Model) string {
	// Filter out description notes (they're shown in details panel)
	var displayNotes []int
	for i, note := range m.notes {
		if !note.IsDescription {
			displayNotes = append(displayNotes, i)
		}
	}

	if len(displayNotes) == 0 {
		return lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No notes (description shown in Details)")
	}

	items := make([]string, len(displayNotes))
	for idx, i := range displayNotes {
		note := m.notes[i]

		// Truncate note content for list view
		content := note.Content
		if len(content) > 35 {
			content = content[:32] + "..."
		}

		// Replace newlines with spaces for single line display
		content = lipgloss.NewStyle().
			MaxWidth(35).
			Render(content)

		cursor := " "
		if m.activeSection == 2 && idx == 0 {
			// For now, just highlight first note when notes section is active
			cursor = ">"
		}

		items[idx] = fmt.Sprintf("%s • %s", cursor, content)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
