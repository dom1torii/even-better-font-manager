package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (m *model) View() string {
	switch m.state {
	case statePath:
		return m.pathView()
	case stateStart:
		return m.startView()
	default:
		return ""
	}
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

func (m *model) pathView() string {
	var pathChoices []string

	for i, label := range pathItems {
		var row string
		// add actual input instead of *input*
		if i == 0 {
			inputView := m.pathInput.View()
			row = pathItem("(1) Path: "+inputView, m.PathSelection == i)
		} else {
			row = pathItem(label, m.PathSelection == i)
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

var startItems = []string{
	"(1) Choose font",
	"(2) Apply",
	"(3) Preview font",
	"(4) Change CS2 path",
	"(q) Quit",
}

func startItem(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}

func (m *model) startView() string {
	var startChoices []string
	for i, label := range startItems {
		startChoices = append(startChoices, startItem(label, m.StartSelection == i))
		// add status lines
		if i == 2 {
			status := fmt.Sprintf("    Will apply: %s", "font name")
			startChoices = append(startChoices, statusStyle.Render(status))
		}
	}

	items := strings.Join(startChoices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Even Better Font Manager"), m.width),
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
