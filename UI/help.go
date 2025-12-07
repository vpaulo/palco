package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func RenderHelp(m Model) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		MarginBottom(1).
		Align(lipgloss.Center)

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(special).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		Width(20)

	descStyle := lipgloss.NewStyle().
		Foreground(normal)

	helpContent := []string{
		titleStyle.Render("Palco - Keybindings"),
		"",
		sectionTitleStyle.Render("Navigation"),
		keyStyle.Render("↑/↓ or k/j") + descStyle.Render("Move cursor up/down in active section"),
		keyStyle.Render("Tab") + descStyle.Render("Switch to next section"),
		keyStyle.Render("Shift+Tab") + descStyle.Render("Switch to previous section"),
		"",
		sectionTitleStyle.Render("Projects Section"),
		keyStyle.Render("n") + descStyle.Render("Create new project"),
		keyStyle.Render("e") + descStyle.Render("Edit selected project"),
		keyStyle.Render("d") + descStyle.Render("Delete selected project"),
		"",
		sectionTitleStyle.Render("Tasks Section"),
		keyStyle.Render("n") + descStyle.Render("Create new task"),
		keyStyle.Render("e") + descStyle.Render("Edit selected task"),
		keyStyle.Render("d") + descStyle.Render("Delete selected task"),
		keyStyle.Render("Space/Enter") + descStyle.Render("Toggle task completion"),
		"",
		sectionTitleStyle.Render("Notes Section"),
		keyStyle.Render("n") + descStyle.Render("Create new note for selected task"),
		"",
		sectionTitleStyle.Render("Forms"),
		keyStyle.Render("Tab/Shift+Tab") + descStyle.Render("Switch between form fields"),
		keyStyle.Render("Enter") + descStyle.Render("Submit form"),
		keyStyle.Render("Esc") + descStyle.Render("Cancel form"),
		"",
		sectionTitleStyle.Render("General"),
		keyStyle.Render("?") + descStyle.Render("Show this help screen"),
		keyStyle.Render("q or Ctrl+C") + descStyle.Render("Quit application"),
		"",
		"",
		lipgloss.NewStyle().
			Foreground(subtle).
			Align(lipgloss.Center).
			Render("Press any key to close this help screen"),
	}

	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(2, 4).
		Width(70)

	helpText := lipgloss.JoinVertical(lipgloss.Left, helpContent...)
	styledHelp := helpBox.Render(helpText)

	// Center the help box
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		styledHelp,
	)
}
