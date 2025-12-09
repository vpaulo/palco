package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func Section(isActive bool) lipgloss.Style {
	borderColor := subtle
	if isActive {
		borderColor = highlight
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)
}
