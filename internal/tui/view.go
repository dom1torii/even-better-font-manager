package tui

func (m *model) View() string {
	switch m.state {
	case statePath:
		return m.pathView()
	case stateStart:
		return m.startView()
	case stateFonts:
		return m.fontsView()
	default:
		return ""
	}
}
