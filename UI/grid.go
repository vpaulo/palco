package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func Grid(m Model) string {
	// Calculate column widths as percentages
	col1Width := int(float64(m.width) * 0.40)
	col2Width := int(float64(m.width) * 0.40)
	// Make the third column take remaining width to avoid rounding gaps
	// 6 is for the 6 borders
	col3Width := m.width - col1Width - col2Width - 6

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				Projects(m),
				Tasks(m),
				Notes(m),
			),
		),
		lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				Details(m),
			),
		),
		lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				Section(m.activeSection == 4).Width(col3Width).Height(m.height-3).Render("Drafts [5]"),
			),
		),
	)
}
