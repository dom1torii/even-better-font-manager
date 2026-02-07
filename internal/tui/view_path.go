package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type pathModel struct {
	selection int
	pathInput textinput.Model
}

func initialPathModel() pathModel {
	return pathModel{
		selection: 0,
		pathInput: createInput("/path/to/cs2/", "Path: ", 20),
	}
}

func (m *model) updatePathSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// keypresses inside focused input
		if m.path.pathInput.Focused() {
			switch msg.String() {
			case "enter":
				m.path.pathInput.Blur()
				m.csPath = m.path.pathInput.Value()
				return m, nil
			case "esc":
				m.path.pathInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.path.pathInput, cmd = m.path.pathInput.Update(msg)
			return m, cmd
		}

		// normal keypresses
		switch msg.String() {
		case "1":
			m.path.selection = 0
			return m, nil
		case "2":
			m.path.selection = 1
			return m, nil
		case "3":
			m.path.selection = 2
			return m, nil
		case "4":
			m.path.selection = 3
			return m, nil
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.path.selection < len(pathItems)-1 {
				m.path.selection++
			} else {
				m.path.selection = 0
			}
		case "k", "up":
			if m.path.selection > 0 {
				m.path.selection--
			} else {
				m.path.selection = len(pathItems) - 1
			}
		case "enter", " ":
			if m.path.selection == 0 {
				return m, m.path.pathInput.Focus()
			}
			if m.path.selection == 1 {
				return m, m.chooseCsPath("Choose your Counter-Strike Global Offensive/ folder")
			}
			if m.path.selection == 2 {
				return m, m.detectCsPath()
			}
			if m.path.selection == 3 {
				m.state = stateStart
				return m, m.confirmCsPath()
			}
			if m.path.selection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		case "i":
			if m.path.selection == 0 {
				return m, m.path.pathInput.Focus()
			}
		}
	}
	return m, nil
}

func (m *model) pathView() string {
	var pathChoices []string

	for i, label := range pathItems {
		var row string
		// add actual input instead of *input*
		if i == 0 {
			inputView := m.path.pathInput.View()
			row = pathItem("(1) "+inputView, m.path.selection == i)
		} else {
			row = pathItem(label, m.path.selection == i)
		}
		pathChoices = append(pathChoices, row)

		if i == 3 {
			status := fmt.Sprintf("    Chosen path: %s", m.csPath)
			pathChoices = append(pathChoices, statusStyle.Render(status))
		}
	}

	items := strings.Join(pathChoices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Choose cs2 installation path:"), m.width),
		lipgloss.NewStyle().Width(35).Render(items),
		wordwrap.String(helpStyle.Render("(↓↑: move | space/enter: select | q/esc: quit)"), m.width),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

var pathItems = []string{
	"(1) Path: *input*",
	"(2) Open file chooser",
	"(3) Try to detect",
	"(4) Confirm",
	"(q) Quit",
}

func pathItem(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
