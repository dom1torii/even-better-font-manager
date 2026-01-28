package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Even Better Font Manager"),
	)
}
