package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dom1torii/even-better-font-manager/internal/platform/detectpath"
	"github.com/dom1torii/even-better-font-manager/internal/platform/filechooser"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Even Better Font Manager"),
	)
}

func (m *model) detectCsPath() tea.Cmd {
	return func() tea.Msg {
		path := detectpath.DetectCS2Path()
		return csPathMsg(path)
	}
}

func (m *model) chooseCsPath() tea.Cmd {
	return func() tea.Msg {
		path := filechooser.ChoosePath()
		return csPathMsg(path)
	}
}
