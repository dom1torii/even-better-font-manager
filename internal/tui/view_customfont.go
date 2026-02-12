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
}

func initialCustomFontModel() customFontModel {
	return customFontModel{
		selection: 0,
		pathInput: createInput("/path/to/font", "Path: ", 20),
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
		case "1":
			m.customFont.selection = 0
			return m, nil
		case "2":
			m.customFont.selection = 1
			return m, nil
		case "3":
			m.customFont.selection = 2
			return m, nil
		case "4":
			m.customFont.selection = 3
			return m, nil
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.customFont.selection < len(customFontItems)-1 {
				m.customFont.selection++
			} else {
				m.customFont.selection = 0
			}
		case "k", "up":
			if m.customFont.selection > 0 {
				m.customFont.selection--
			} else {
				m.customFont.selection = len(customFontItems) - 1
			}
		case "enter", " ":
			if m.customFont.selection == 0 {
				return m, m.customFont.pathInput.Focus()
			}
			if m.customFont.selection == 1 {
				return m, m.chooseCustomFontPath("Choose custom font to add")
			}
			if m.customFont.selection == 2 {
				return m, m.previewFont(m.customFont.chosenFont.Path)
			}
			if m.customFont.selection == 3 {
			  m.fonts.collection.Add(m.customFont.chosenFont.Font)
			  m.state = stateFonts
			  return m, nil
			}
		case "i":
			if m.customFont.selection == 0 {
				return m, m.customFont.pathInput.Focus()
			}
		}
	}
	return m, nil
}

func (m *model) customFontView() string {
	var customFontChoices []string

	for i, label := range customFontItems {
		var row string
		// add actual input instead of *input*
		if i == 0 {
			inputView := m.customFont.pathInput.View()
			row = customFontItem("(1) "+inputView, m.customFont.selection == i)
		} else {
			row = customFontItem(label, m.customFont.selection == i)
		}
		customFontChoices = append(customFontChoices, row)

		if i == 3 {
			status := fmt.Sprintf("    Will add: %s", m.customFont.chosenFont.Name + " " + m.customFont.chosenFont.Style)
			customFontChoices = append(customFontChoices, statusStyle.Render(status))
		}

		if i == 1 && m.customFont.chosenFont.Error != nil {
			status := fmt.Sprintf("    Error: %s", m.customFont.chosenFont.Error)
			customFontChoices = append(customFontChoices, statusWarningStyle.Render(status))
		}
	}

	items := strings.Join(customFontChoices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Add custom font:"), m.width),
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

var customFontItems = []string{
	"(1) Path: *input*",
	"(2) Open file chooser",
	"(3) Preview",
	"(4) Confirm",
}

func customFontItem(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
