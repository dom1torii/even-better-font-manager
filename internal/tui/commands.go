package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/dom1torii/even-better-font-manager/internal/fontconfig"
	"github.com/dom1torii/even-better-font-manager/internal/platform/detectpath"
	"github.com/dom1torii/even-better-font-manager/internal/platform/fonts"
	"github.com/dom1torii/even-better-font-manager/internal/zenity"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.getCsPath(),
		m.getSystemFonts(),
		m.loadFontCollection(),
		tea.SetWindowTitle("Even Better Font Manager"),
	)
}

func (m *model) detectCsPath() tea.Cmd {
	return func() tea.Msg {
		path := detectpath.DetectCS2Path()
		return csPathMsg(path)
	}
}

func (m *model) chooseCsPath(title string) tea.Cmd {
	return func() tea.Msg {
		path := zenity.Open("directory", title)
		return csPathMsg(path)
	}
}

func (m *model) confirmCsPath() tea.Cmd {
	return func() tea.Msg {
		m.cfg.General.CS2Path = m.csPath
		m.cfg.Save()
		return pathConfirmedMsg{}
	}
}

func (m *model) getCsPath() tea.Cmd {
	return func() tea.Msg {
		path := m.cfg.General.CS2Path
		return getCsPathMsg(path)
	}
}

func (m *model) getSystemFonts() tea.Cmd {
	return func() tea.Msg {
		systemFonts := fonts.GetSystemFonts()
		return systemFontsMsg(systemFonts)
	}
}

func (m *model) previewFont(path string) tea.Cmd {
	return func() tea.Msg {
		fonts.Preview(path)
		return fontPreviewMsg{}
	}
}

func (m *model) chooseCustomFontPath(title string) tea.Cmd {
	return func() tea.Msg {
		path := zenity.Open("", title)
		return customFontPathMsg(path)
	}
}

func (m *model) setCustomFont(path string) tea.Cmd {
	return func() tea.Msg {
		fontName, fontStyle, err := fonts.GetName(path)
		if err != nil {
			return customFontMsg{
				Font: config.Font{
					Name:  fontName,
					Style: fontStyle,
					Path:  path,
				},
				Error: err,
			}
		}
		return customFontMsg{
			Font: config.Font{
				Name:  fontName,
				Style: fontStyle,
				Path:  path,
			},
			Error: nil,
		}
	}
}

func (m *model) loadFontCollection() tea.Cmd {
	return func() tea.Msg {
		return fontCollectionMsg(config.LoadCollection())
	}
}

func (m *model) writeFontConfig(fontName string, fontPath string) tea.Cmd {
	return func() tea.Msg {
		fontconfig.Apply(m.cfg, fontName, fontPath)
		return fontConfigMsg{}
	}
}
