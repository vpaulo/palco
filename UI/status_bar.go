package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Status Bar.

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))
)

func StatusBar(m Model) string {
	// Build status bar
	w := lipgloss.Width

	// Section name
	sectionName := ""
	switch m.activeSection {
	case 0:
		sectionName = "Projects"
	case 1:
		sectionName = "Tasks"
	case 2:
		sectionName = "Notes"
	case 3:
		sectionName = "Details"
	case 4:
		sectionName = "Drafts"
	}

	statusKey := statusStyle.Render(sectionName)

	// Help hint
	helpHint := encodingStyle.Render("? Help")

	// Status message based on mode and context
	var statusMsg string
	if m.mode != ModeNormal {
		statusMsg = "Editing..."
	} else {
		// Context-aware hints
		switch m.activeSection {
		case 0:
			statusMsg = "n:New  e:Edit  d:Delete  ↑↓:Navigate  Tab:Switch"
		case 1:
			statusMsg = "n:New  e:Edit  d:Delete  Space:Toggle  ↑↓:Navigate"
		case 2:
			statusMsg = "n:New Note  ↑↓:Navigate  Tab:Switch"
		default:
			statusMsg = "Tab:Switch Sections  ?:Help  q:Quit"
		}
	}

	statusVal := statusText.
		Width(m.width - w(statusKey) - w(helpHint) - 4).
		Render(statusMsg)

	return statusBarStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, statusKey, statusVal, helpHint),
	)
}
