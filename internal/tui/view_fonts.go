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
}

func initialFontsModel() fontsModel {
	return fontsModel{
		leftSelection:  0,
		rightSelection: 0,
	}
}

func (m *model) updateFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxViewHeight := 10

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			if m.fonts.colActive == 0 && m.fonts.leftSelection > 0 {
				m.fonts.leftSelection--
			} else if m.fonts.colActive == 1 && m.fonts.rightSelection > 0 {
				m.fonts.rightSelection--
			}

		case "j", "down":
			if m.fonts.colActive == 0 && m.fonts.leftSelection < len(m.fonts.collection.List())-1 {
				m.fonts.leftSelection++
			} else if m.fonts.colActive == 1 && m.fonts.rightSelection < len(fontItems)-1 {
				m.fonts.rightSelection++
			}

		case "h", "left":
			if m.fonts.colActive == 1 {
				m.fonts.colActive = 0
			}

		case "l", "right":
			if m.fonts.colActive == 0 {
				m.fonts.colActive = 1
			}

		case "enter", " ":
			if m.fonts.colActive == 0 {
				selectedFont := m.fonts.collection.List()[m.fonts.leftSelection]
				m.chosenFont = font{
					Font: selectedFont,
					Error: nil,
				}
				m.state = stateStart
			}
			if m.fonts.colActive == 1 && m.fonts.rightSelection == 0 {
				m.state = stateCustomFont
				return m, nil
			}
			if m.fonts.colActive == 1 && m.fonts.rightSelection == 1 {
				m.state = stateSystemFont
				return m, nil
			}

		case "d":
			m.fonts.collection.Remove(m.fonts.leftSelection)
			return m, nil

		case "p":
			return m, m.previewFont(m.fonts.collection.List()[m.fonts.leftSelection].Path)

		case "q", "esc":
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
			leftChoices = append(leftChoices, fontItem(m.fonts.collection.List()[i].Name, isSelected, leftWidth-4))
		}
	}

	var rightChoices []string
	for i, label := range fontItems {
		isSelected := (m.fonts.colActive == 1 && m.fonts.rightSelection == i)
		rightChoices = append(rightChoices, fontItem(label, isSelected, 0))
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

var fontItems = []string{
	"(1) Custom",
	"(2) From system",
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
