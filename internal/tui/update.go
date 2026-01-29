package tui

import (
	// "log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case csPathMsg:
		m.csPath = string(msg)
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
		}
	}

	return m, nil
}

func (m *model) updateStartSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.StartSelection = 0
			// m.state = stateRelays
			return m, nil
		case "2":
			m.StartSelection = 1
			// m.state = statePresets
		case "3":
			m.StartSelection = 2
			return m, nil
		case "4":
			m.StartSelection = 3
			return m, nil
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.StartSelection < len(startItems)-1 {
				m.StartSelection++
			} else {
				m.StartSelection = 0
			}
		case "k", "up":
			if m.StartSelection > 0 {
				m.StartSelection--
			} else {
				m.StartSelection = len(startItems) - 1
			}
		case "enter", " ":
			if m.StartSelection == 0 {
				// m.state = stateRelays
				return m, nil
			}
			if m.StartSelection == 1 {
				// m.state = statePresets
				return m, nil
			}
			if m.StartSelection == 2 {
				return m, nil
			}
			if m.StartSelection == 3 {
				return m, nil
			}
			if m.StartSelection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *model) updatePathSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.pathInput.Focused() {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.pathInput.Blur()
				m.csPath = m.pathInput.Value()
				return m, nil
			case "esc":
				m.pathInput.Blur()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.PathSelection = 0
			// m.state = stateRelays
			return m, nil
		case "2":
			m.PathSelection = 1
			// m.state = statePresets
		case "3":
			m.PathSelection = 2
			return m, nil
		case "4":
			m.PathSelection = 3
			return m, nil
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.PathSelection < len(pathItems)-1 {
				m.PathSelection++
			} else {
				m.PathSelection = 0
			}
		case "k", "up":
			if m.PathSelection > 0 {
				m.PathSelection--
			} else {
				m.PathSelection = len(pathItems) - 1
			}
		case "enter", " ":
			if m.PathSelection == 0 {
				return m, m.pathInput.Focus()
			}
			if m.PathSelection == 1 {
				return m, m.chooseCsPath()
			}
			if m.PathSelection == 2 {
				return m, m.detectCsPath()
			}
			if m.PathSelection == 3 {
				return m, nil
			}
			if m.PathSelection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
