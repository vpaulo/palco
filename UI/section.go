package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func Section() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle)
}
