package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/even-better-font-manager/internal/platform/detectpath"
	"github.com/dom1torii/even-better-font-manager/internal/platform/filechooser"
	"github.com/dom1torii/even-better-font-manager/internal/platform/previewfont"
	"github.com/dom1torii/even-better-font-manager/internal/platform/systemfonts"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.getCsPath(),
		m.getSystemFonts(),
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

func (m *model) confirmCsPath() tea.Cmd {
  return func() tea.Msg {
    m.cfg.General.CS2Path = m.csPath
    m.cfg.Apply()
    return pathConfirmedMsg{}
  }
}

func (m *model) getCsPath() tea.Cmd {
	return func() tea.Msg {
		path := m.cfg.General.CS2Path
		return csPathMsg(path)
	}
}

func (m *model) getSystemFonts() tea.Cmd {
	return func() tea.Msg {
		fonts := systemfonts.GetFonts()
		return systemFontsMsg(fonts)
	}
}

func (m *model) previewFont() tea.Cmd {
	return func() tea.Msg {
		previewfont.PreviewFont(m.systemFont.filteredFonts[m.systemFont.selection].Path)
		return fontPreviewMsg{}
	}
}
