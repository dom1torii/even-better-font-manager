package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	tea "github.com/charmbracelet/bubbletea"
)

type startModel struct {
	selection int
}

func initialStartModel() startModel {
	return startModel{
		selection: 0,
	}
}

func (m *model) updateStartSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.start.selection = 0
			m.state = stateFonts
			return m, nil
		case "2":
			m.start.selection = 1
		case "3":
			m.start.selection = 2
			return m, nil
		case "4":
			m.start.selection = 3
			return m, nil
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.start.selection < len(startItems)-1 {
				m.start.selection++
			} else {
				m.start.selection = 0
			}
		case "k", "up":
			if m.start.selection > 0 {
				m.start.selection--
			} else {
				m.start.selection = len(startItems) - 1
			}
		case "enter", " ":
			if m.start.selection == 0 {
				m.state = stateFonts
				return m, nil
			}
			if m.start.selection == 1 {
				return m, nil
			}
			if m.start.selection == 2 {
				return m, nil
			}
			if m.start.selection == 3 {
				m.state = statePath
				return m, nil
			}
			if m.start.selection == 4 {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *model) startView() string {
	var startChoices []string
	for i, label := range startItems {
		startChoices = append(startChoices, startItem(label, m.start.selection == i))
		// add status lines
		if i == 2 {
			status := fmt.Sprintf("    Will apply: %s", m.chosenFont.Name)
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

var startItems = []string{
	"(1) Choose font",
	"(2) Preview font",
	"(3) Apply",
	"(4) Change CS2 path",
	"(q) Quit",
}

func startItem(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
