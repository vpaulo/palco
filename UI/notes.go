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

	// Determine what to show based on note context
	if m.noteContext == 0 {
		// Project notes
		if len(m.projects) == 0 {
			content = lipgloss.NewStyle().
				Foreground(subtle).
				Padding(1).
				Render("No project selected")
		} else if len(m.notes) == 0 {
			content = lipgloss.NewStyle().
				Foreground(subtle).
				Padding(1).
				Render("No notes for this project")
		} else {
			content = renderNotesList(m)
		}
	} else {
		// Task notes
		if len(m.tasks) == 0 {
			content = lipgloss.NewStyle().
				Foreground(subtle).
				Padding(1).
				Render("No task selected")
		} else if len(m.notes) == 0 {
			content = lipgloss.NewStyle().
				Foreground(subtle).
				Padding(1).
				Render("No notes for this task")
		} else {
			content = renderNotesList(m)
		}
	}

	return Section(m.activeSection == 2).Width(col1Width).Height(row1Height - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Notes [3]"),
			content,
		),
	)
}

func renderNotesList(m Model) string {
	// Filter out description notes for tasks (they're shown in details panel)
	// Project notes don't have descriptions, so show all
	var displayNotes []int
	for i, note := range m.notes {
		// Only filter description notes for task context
		if m.noteContext == 0 || !note.IsDescription {
			displayNotes = append(displayNotes, i)
		}
	}

	if len(displayNotes) == 0 {
		if m.noteContext == 1 {
			return lipgloss.NewStyle().
				Foreground(subtle).
				Padding(1).
				Render("No notes (description shown in Details)")
		}
		return lipgloss.NewStyle().
			Foreground(subtle).
			Padding(1).
			Render("No notes")
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
