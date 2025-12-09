package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func RenderForm(m Model) string {
	if m.mode == ModeNormal {
		return ""
	}

	var title string
	var fields []string

	if m.mode == ModeCreateProject {
		title = "Create New Project"
		fields = []string{"Name:", "Description:"}
	} else if m.mode == ModeCreateTask {
		title = "Create New Task"
		fields = []string{"Title:", "Description:", "Priority:"}
	} else if m.mode == ModeEditProject {
		title = "Edit Project"
		fields = []string{"Name:", "Description:"}
	} else if m.mode == ModeEditTask {
		title = "Edit Task"
		fields = []string{"Title:", "Description:", "Priority:"}
	} else if m.mode == ModeCreateNote {
		title = "Create New Note"
		fields = []string{"Content:"}
	}

	// Build form
	var formParts []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		MarginBottom(1)
	formParts = append(formParts, titleStyle.Render(title))

	// Input fields
	for i, field := range fields {
		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(special)

		formParts = append(formParts, labelStyle.Render(field))
		formParts = append(formParts, m.formInputs[i].View())
		formParts = append(formParts, "")
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(subtle).
		MarginTop(1)

	helpText := "Tab/Shift+Tab: Switch fields • Enter: Submit • Esc: Cancel"
	formParts = append(formParts, helpStyle.Render(helpText))

	// Container
	formContent := lipgloss.JoinVertical(lipgloss.Left, formParts...)

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(2, 4).
		Width(60)

	formBox := formStyle.Render(formContent)

	// Overlay background
	overlayStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return overlayStyle.Render(formBox)
}

// Helper function to render a form field
func renderFormField(label string, input string) string {
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(special)

	return fmt.Sprintf("%s\n%s", labelStyle.Render(label), input)
}
