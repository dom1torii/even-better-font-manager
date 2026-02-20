package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) nextMenuItem(current int, menuItems []menuItem) int {
	for i := current + 1; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		if menuItem.disabled == nil || !menuItem.disabled(m) {
			return i
		}
	}
	return current
}

func (m *model) prevMenuItem(current int, menuItems []menuItem) int {
	for i := current - 1; i >= 0; i-- {
		menuItem := menuItems[i]
		if menuItem.disabled == nil || !menuItem.disabled(m) {
			return i
		}
	}
	return current
}

func (m *model) handleNumberPress(key string, selection *int, menuItems []menuItem) tea.Cmd {
	// we get an int out of a key, then if its not disabled change the selection and do the action
	idx := int(key[0] - '1')
	if idx < len(menuItems) {
		if menuItems[idx].disabled != nil && menuItems[idx].disabled(m) {
			return nil
		}
		*selection = idx
		return menuItems[idx].action(m)
	}
	return nil
}
