package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
)

type fontsModel struct {
	colActive      int
	leftSelection  int
	rightSelection int
	startRow       int
	collection     *config.Collection
	menuItems      []menuItem
}

func initialFontsModel() fontsModel {
	return fontsModel{
		leftSelection:  0,
		rightSelection: 0,
		menuItems: []menuItem{
	    {
	      render: func(m *model) string {
					return "(1) Custom"
				},
	      action: func(m *model) tea.Cmd {
	        m.state = stateCustomFont
	        return nil
	      },
	    },
	    {
	      render: func(m *model) string {
					return "(2) From system"
				},
	      action: func(m *model) tea.Cmd {
	        m.state = stateSystemFont
	        return nil
	      },
	    },
  	},
	}
}

func (m *model) updateFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxViewHeight := 10

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1", "2":
			idx := int(msg.String()[0] - '1')
		  if idx < len(m.fonts.menuItems) {
		    m.fonts.rightSelection = idx
		    return m, m.fonts.menuItems[idx].action(m)
		  }
		case "j", "down":
			if m.fonts.colActive == 0 && m.fonts.leftSelection < len(m.fonts.collection.List())-1 {
				m.fonts.leftSelection++
			} else if m.fonts.colActive == 1 && m.fonts.rightSelection < len(m.fonts.menuItems)-1 {
				m.fonts.rightSelection++
			}
		case "k", "up":
			if m.fonts.colActive == 0 && m.fonts.leftSelection > 0 {
				m.fonts.leftSelection--
			} else if m.fonts.colActive == 1 && m.fonts.rightSelection > 0 {
				m.fonts.rightSelection--
			}
		case "h", "left":
			if m.fonts.colActive == 1 {
				m.fonts.colActive = 0
			}
		case "l", "right":
			if m.fonts.colActive == 0 {
				m.fonts.colActive = 1
			}
		case "d":
			if m.fonts.colActive == 0 {
	      m.fonts.collection.Remove(m.fonts.leftSelection)
	    }
			return m, nil
		case "p":
			if m.fonts.colActive == 0 {
	      return m, m.previewFont(m.fonts.collection.List()[m.fonts.leftSelection].Path)
	    }
		case "enter", " ":
			if m.fonts.colActive == 0 {
				selectedFont := m.fonts.collection.List()[m.fonts.leftSelection]
				m.chosenFont = font{
					Font: selectedFont,
					Error: nil,
				}
				m.state = stateStart
				return m, nil
			} else {
				return m, m.fonts.menuItems[m.fonts.rightSelection].action(m)
			}
    case "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		}

		// left column scrolling
		if m.fonts.colActive == 0 {
			if m.fonts.leftSelection >= m.fonts.startRow+maxViewHeight {
				m.fonts.startRow = m.fonts.leftSelection - maxViewHeight + 1
			} else if m.fonts.leftSelection < m.fonts.startRow {
				m.fonts.startRow = m.fonts.leftSelection
			}
		}
	}
	return m, nil
}

func (m *model) fontsView() string {
	maxViewHeight := 10
	leftWidth := max(min(m.width-30, 40), 20)

	var leftChoices []string
	if len(m.fonts.collection.List()) == 0 {
		emptySign := emptyListSignStyle.Render("No fonts added yet. Add them using options on the right")
		leftChoices = append(leftChoices, emptySign)
		m.fonts.colActive = 1
	} else {
		start := m.fonts.startRow
		end := min(start+maxViewHeight, len(m.fonts.collection.List()))
		for i := start; i < end; i++ {
			isSelected := (m.fonts.colActive == 0 && m.fonts.leftSelection == i)
			leftChoices = append(leftChoices, fontItem(m.fonts.collection.List()[i].Name + " " + m.fonts.collection.List()[i].Style, isSelected, leftWidth-4))
		}
	}

	var rightChoices []string
  for i, item := range m.fonts.menuItems {
    isSelected := (m.fonts.colActive == 1 && m.fonts.rightSelection == i)
    rightChoices = append(rightChoices, fontItem(item.render(m), isSelected, 0))
  }

	rightTitle := lipgloss.NewStyle().PaddingTop(1).Render("Add font:")

	leftCol := lipgloss.NewStyle().PaddingRight(4).Border(lipgloss.NormalBorder()).Width(leftWidth).Height(maxViewHeight).Render(lipgloss.JoinVertical(lipgloss.Left, leftChoices...))
	rightCol := lipgloss.JoinVertical(
		lipgloss.Left,
		rightTitle,
		lipgloss.NewStyle().Render(lipgloss.JoinVertical(lipgloss.Left, rightChoices...)),
	)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Your fonts"),
		columns,
		wordwrap.String(helpStyle.Render("(←↓↑→: move | p: preview font | d: remove font | q: quit)"), m.width),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}

func fontItem(label string, isSelected bool, width int) string {
	truncated := label
	if width > 0 {
		truncated = truncate.StringWithTail(label, uint(width), "...")
	}
	if isSelected {
		return selectionStyle.Render(truncated)
	}
	return truncated
}
