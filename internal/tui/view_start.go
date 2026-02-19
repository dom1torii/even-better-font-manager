package tui

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type startModel struct {
	selection int
	fontSize  float64
	menuItems []menuItem

	applyStatus string
	resetStatus string
}

func initialStartModel() startModel {
	return startModel{
		selection: 0,
		fontSize:  1.0,
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
					return fmt.Sprintf("(2) Font size: %.1f", m.start.fontSize)
				},
				action: func(m *model) tea.Cmd {
					return nil
				},
			},
			{
				render: func(m *model) string {
					return "(3) Preview font"
				},
				action: func(m *model) tea.Cmd {
					return m.previewFont(m.chosenFont.Path)
				},
			},
			{
				render: func(m *model) string {
					return "(4) Apply"
				},
				status: func(m *model) string {
					if m.start.applyStatus != "" {
						return fmt.Sprintf("    %s", statusOkStyle.Render(m.start.applyStatus))
					}
					return fmt.Sprintf("    Will apply: %s %s", m.chosenFont.Name, m.chosenFont.Style)
				},
				action: func(m *model) tea.Cmd {
					return m.writeFontConfig(m.chosenFont.Name, m.chosenFont.Path, m.start.fontSize)
				},
			},
			{
				render: func(m *model) string {
					return "(5) Reset"
				},
				status: func(m *model) string {
					if m.start.resetStatus != "" {
						return fmt.Sprintf("    %s", statusOkStyle.Render(m.start.resetStatus))
					}
					return ""
				},
				action: func(m *model) tea.Cmd {
					return m.resetFontConfig()
				},
			},
			{
				render: func(m *model) string {
					return "(6) Change CS2 path"
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
		case "1", "2", "3", "4", "5", "6":
			idx := int(msg.String()[0] - '1')
			if idx < len(m.start.menuItems) {
				m.start.selection = idx
				return m, m.start.menuItems[idx].action(m)
			}
		case "j", "down":
			if m.start.selection < len(m.start.menuItems)-1 {
				m.start.selection++
			}
		case "k", "up":
			if m.start.selection > 0 {
				m.start.selection--
			}
		case "h", "left":
			if m.start.selection == 1 {
				if m.start.fontSize > 0.1 {
					m.start.fontSize = math.Round((m.start.fontSize-0.1)*10) / 10
				} else {
					m.start.fontSize = 0
				}
			}
		case "l", "right":
			if m.start.selection == 1 {
				if m.start.fontSize < 9.9 {
					m.start.fontSize = math.Round((m.start.fontSize+0.1)*10) / 10
				} else {
					m.start.fontSize = 10.0
				}
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
		wordwrap.String(helpStyle.Render("(↓↑: move | ←→: adjust font size | space/enter: select | q/esc: quit)"), m.width),
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
