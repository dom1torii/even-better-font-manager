package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type pathModel struct {
	selection int
	pathInput textinput.Model
	menuItems []menuItem
	chosenPath string
}

func initialPathModel() pathModel {
	return pathModel{
		selection: 0,
		pathInput: createInput("/path/to/cs2/", "Path: ", 20),
		menuItems: []menuItem{
			{
				render: func(m *model) string {
					return "(1) " + m.path.pathInput.View()
				},
				action: func(m *model) tea.Cmd {
					return m.path.pathInput.Focus()
				},
			},
	  	{
				render: func(m *model) string {
					return "(2) Open file chooser"
				},
	      action: func(m *model) tea.Cmd {
	        return m.chooseCsPath("Choose your Counter-Strike Global Offensive/ folder")
	      },
	    },
	    {
				render: func(m *model) string {
					return "(3) Try to detect"
				},
	      action: func(m *model) tea.Cmd {
	        return m.detectCsPath()
	      },
	    },
	    {
				render: func(m *model) string {
					return "(4) Confirm"
				},
		    status: func(m *model) string {
		    	return fmt.Sprintf("    Chosen path: %s", m.path.chosenPath)
		    },
	      action: func(m *model) tea.Cmd {
					m.csPath = m.path.chosenPath
					m.state = stateStart
	        return m.confirmCsPath()
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

func (m *model) updatePathSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// keypresses inside focused input
		if m.path.pathInput.Focused() {
			switch msg.String() {
			case "enter":
				m.path.pathInput.Blur()
				m.path.chosenPath = m.path.pathInput.Value()
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
		case "1", "2", "3", "4":
		  idx := int(msg.String()[0] - '1')
		  if idx < len(m.path.menuItems) {
		    m.path.selection = idx
		    return m, m.path.menuItems[idx].action(m)
		  }
		case "j", "down":
			if m.path.selection < len(m.path.menuItems)-1 {
      	m.path.selection++
      }
		case "k", "up":
			if m.path.selection > 0 {
        m.path.selection--
      }
    case "i":
			if m.path.selection == 0 {
				return m, m.path.pathInput.Focus()
			}
		case "enter", " ":
			return m, m.path.menuItems[m.path.selection].action(m)
		case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) pathView() string {
	var choices []string
	for i, item := range m.path.menuItems {
		choices = append(choices, pathItemStyle(item.render(m), m.path.selection == i))
		if item.status != nil && item.status(m) != "" {
			choices = append(choices, statusStyle.Render(item.status(m)))
		}
	}

	list := strings.Join(choices, "\n")
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wordwrap.String(titleStyle.Render("Choose cs2 installation path:"), m.width),
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

func pathItemStyle(label string, isSelected bool) string {
	if isSelected {
		return selectionStyle.Render(label)
	}
	return label
}
