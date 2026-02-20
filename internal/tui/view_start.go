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
				disabled: func(m *model) bool {
				  return m.chosenFont.Name == ""
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
				disabled: func(m *model) bool {
				  return m.chosenFont.Name == ""
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
			return m, m.handleNumberPress(msg.String(), &m.start.selection, m.start.menuItems)
		case "j", "down":
			m.start.selection = m.nextMenuItem(m.start.selection, m.start.menuItems)
		case "k", "up":
			m.start.selection = m.prevMenuItem(m.start.selection, m.start.menuItems)
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
		if item.disabled != nil && item.disabled(m) {
			choices = append(choices, disabledStyle.Render(item.render(m)))
		} else {
			choices = append(choices, menuItemStyle(item.render(m), m.start.selection == i, 0))
			if item.status != nil && item.status(m) != "" {
				choices = append(choices, statusStyle.Render(item.status(m)))
			}
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
