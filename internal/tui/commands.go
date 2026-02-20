package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/even-better-font-manager/internal/filechooser"
	"github.com/dom1torii/even-better-font-manager/internal/fontconfig"
	"github.com/dom1torii/even-better-font-manager/internal/platform/cs2path"
	"github.com/dom1torii/even-better-font-manager/internal/platform/fonts"
)

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.getCsPath(),
		m.getSystemFonts(),
		m.loadFontCollection(),
		tea.SetWindowTitle("Even Better Font Manager"),
	)
}

// 3 functions below probably work perfectly, but maybe i need to simplify it all into 1 function?
// or 2 functions at least? verifyCsPath function doesn't really seem necessary idk
func (m *model) detectCsPath() tea.Cmd {
	return func() tea.Msg {
		path := cs2path.Detect()
		return csPathMsg(cs2path.Verify(path))
	}
}

func (m *model) chooseCsPath(title string) tea.Cmd {
	return func() tea.Msg {
		path := filechooser.Open("directory", title)
		return csPathMsg(cs2path.Verify(path))
	}
}

func (m *model) verifyCsPath(path string) tea.Cmd {
	return func() tea.Msg {
		return csPathMsg(cs2path.Verify(path))
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
		path := filechooser.Open("", title)
		return customFontPathMsg(path)
	}
}

func (m *model) setCustomFont(path string) tea.Cmd {
	return func() tea.Msg {
		f, err := fonts.ParseFont(path)
		if err != nil {
			return customFontMsg{
				Font: fonts.Font{
					Name:  f.Name,
					Style: f.Style,
					Path:  path,
				},
				Error: err,
			}
		}
		return customFontMsg{
			Font: fonts.Font{
				Name:  f.Name,
				Style: f.Style,
				Path:  path,
			},
			Error: nil,
		}
	}
}

func (m *model) loadFontCollection() tea.Cmd {
	return func() tea.Msg {
		return fontCollectionMsg(fonts.GetCollection(m.cfg))
	}
}

func (m *model) removeFontFromCollection(path string) tea.Cmd {
	return func() tea.Msg {
		fonts.RemoveFromCollection(path)
		return removeFontFromCollectionMsg{}
	}
}

func (m *model) addFontToCollection(path string) tea.Cmd {
	return func() tea.Msg {
		fonts.AddToCollection(m.cfg, path)
		return addFontToCollectionMsg{}
	}
}

func (m *model) writeFontConfig(fontName string, fontPath string, fontSize float64) tea.Cmd {
	return func() tea.Msg {
		fontconfig.Apply(m.cfg, fontName, fontPath, fontSize)
		return writeFontConfigMsg("Font applied successfully")
	}
}

func (m *model) resetFontConfig() tea.Cmd {
	return func() tea.Msg {
		fontconfig.Reset(m.cfg)
		return resetFontConfigMsg("Config reset successfully")
	}
}

// idk if these are the correct way to clear the statuses but its the only way i found
func (m *model) clearApplyStatus() tea.Cmd {
  return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
    return clearApplyStatusMsg{}
  })
}

func (m *model) clearResetStatus() tea.Cmd {
  return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
    return clearResetStatusMsg{}
  })
}
