package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case customFontPathMsg:
		return m, m.setCustomFont(string(msg))

	case customFontMsg:
		m.customFont.chosenFont = font(msg)
		return m, nil

	case systemFontsMsg:
		m.systemFont.fonts = msg
		m.systemFont.filteredFonts = msg
		return m, nil

	case fontCollectionMsg:
		m.fonts.collection = msg

	case csPathMsg:
		m.path.chosenPath = string(msg)
		return m, nil

	case getCsPathMsg:
		m.csPath = string(msg)

	case pathConfirmedMsg:
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case stateStart:
			return m.updateStartSelection(msg)
		case statePath:
			return m.updatePathSelection(msg)
		case stateFonts:
			return m.updateFontSelection(msg)
		case stateSystemFont:
			return m.updateSystemFontSelection(msg)
		case stateCustomFont:
			return m.updateCustomFontSelection(msg)
		}
	}

	return m, nil
}
