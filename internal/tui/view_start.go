package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type startModel struct {
	selection int
	menuItems []menuItem
}

func initialStartModel() startModel {
	return startModel{
		selection: 0,
		menuItems: []menuItem{
			{
				render: func(m *model) string {
					return "(1) Choose Font"
				},
				action: func(m *model) tea.Cmd {
					m.state = stateFonts
					return nil
				},
			},
			{
				render: func(m *model) string {
					return "(2) Preview font"
				},
				action: func(m *model) tea.Cmd {
					return m.previewFont(m.chosenFont.Path)
				},
			},
			{
				render: func(m *model) string {
					return "(3) Apply"
				},
				status: func(m *model) string {
					return fmt.Sprintf("    Will apply: %s %s", m.chosenFont.Name, m.chosenFont.Style)
				},
				action: func(m *model) tea.Cmd {
					return m.writeFontConfig(m.chosenFont.Name, m.chosenFont.Path)
				},
			},
			{
				render: func(m *model) string {
					return "(4) Change CS2 path"
				},
				action: func(m *model) tea.Cmd {
					m.state = statePath
					return nil
				},
			},
			{
				render: func(m *model) string {
					return "(q) Quit"
				},
				action: func(m *model) tea.Cmd {
					m.Quitting = true
					return tea.Quit
				},
			},
		},
	}
}

func (m *model) updateStartSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1", "2", "3", "4":
			idx := int(msg.String()[0] - '1')
			if idx < len(m.start.menuItems) {
				m.path.selection = idx
				return m, m.start.menuItems[idx].action(m)
			}
		case "j", "down":
			if m.path.selection < len(m.start.menuItems)-1 {
				m.path.selection++
			}
		case "k", "up":
			m.start.selection--
			if m.start.selection < 0 {
				m.start.selection = len(m.start.menuItems) - 1
			}
		case "enter", " ":
			return m, m.start.menuItems[m.start.selection].action(m)
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) startView() string {
	var choices []string
	for i, item := range m.start.menuItems {
		choices = append(choices, startItemStyle(item.render(m), m.start.selection == i))
		if item.status != nil && item.status(m) != "" {
			choices = append(choices, statusStyle.Render(item.status(m)))
		}
	}

	list := strings.Join(choices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Even Better Font Manager"), m.width),
		lipgloss.NewStyle().Width(35).Render(list),
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

func startItemStyle(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
