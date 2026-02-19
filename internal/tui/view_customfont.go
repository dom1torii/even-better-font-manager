package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type customFontModel struct {
	selection  int
	pathInput  textinput.Model
	chosenFont font
	menuItems  []menuItem
}

func initialCustomFontModel() customFontModel {
	return customFontModel{
		selection: 0,
		pathInput: createInput("/path/to/font", "Path: ", 20),
		menuItems: []menuItem{
			{
				render: func(m *model) string {
					return "(1) " + m.customFont.pathInput.View()
				},
				action: func(m *model) tea.Cmd {
					return m.customFont.pathInput.Focus()
				},
			},
			{
				render: func(m *model) string {
					return "(2) Open file chooser"
				},
				status: func(m *model) string {
					if m.customFont.chosenFont.Error != nil {
						return fmt.Sprintf("    Error: %s", m.customFont.chosenFont.Error)
					}
					return ""
				},
				action: func(m *model) tea.Cmd {
					return m.chooseCustomFontPath("Choose custom font to add")
				},
			},
			{
				render: func(m *model) string {
					return "(3) Preview"
				},
				action: func(m *model) tea.Cmd {
					return m.previewFont(m.customFont.chosenFont.Path)
				},
			},
			{
				render: func(m *model) string {
					return "(4) Confirm"
				},
				status: func(m *model) string {
					return fmt.Sprintf("    Will add: %s", m.customFont.chosenFont.Name+" "+m.customFont.chosenFont.Style)
				},
				action: func(m *model) tea.Cmd {
					m.state = stateFonts
					return m.addFontToCollection(m.customFont.chosenFont.Path)
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

func (m *model) updateCustomFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// keypresses inside focused input
		if m.customFont.pathInput.Focused() {
			switch msg.String() {
			case "enter":
				m.customFont.pathInput.Blur()
				return m, m.setCustomFont(m.customFont.pathInput.Value())
			case "esc":
				m.customFont.pathInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.customFont.pathInput, cmd = m.customFont.pathInput.Update(msg)
			return m, cmd
		}

		// normal keypresses
		switch msg.String() {
		case "1", "2", "3", "4":
			idx := int(msg.String()[0] - '1')
			if idx < len(m.customFont.menuItems) {
				m.customFont.selection = idx
				return m, m.customFont.menuItems[idx].action(m)
			}
		case "j", "down":
			if m.customFont.selection < len(m.customFont.menuItems)-1 {
				m.customFont.selection++
			}
		case "k", "up":
			if m.customFont.selection > 0 {
				m.customFont.selection--
			}
		case "i":
			if m.customFont.selection == 0 {
				return m, m.customFont.pathInput.Focus()
			}
		case "enter", " ":
			return m, m.customFont.menuItems[m.customFont.selection].action(m)
		case "q", "esc":
			m.state = stateFonts
			return m, nil
		}
	}
	return m, nil
}

func (m *model) customFontView() string {
	var choices []string
	for i, item := range m.customFont.menuItems {
		choices = append(choices, customFontItemStyle(item.render(m), m.customFont.selection == i))
		if item.status != nil && item.status(m) != "" {
			choices = append(choices, statusStyle.Render(item.status(m)))
		}
	}

	list := strings.Join(choices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Add custom font:"), m.width),
		lipgloss.NewStyle().Width(35).Render(list),
		wordwrap.String(helpStyle.Render("(↓↑: move | space/enter: select | q/esc: back)"), m.width),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

func customFontItemStyle(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
