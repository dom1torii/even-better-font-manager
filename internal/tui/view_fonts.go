package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/truncate"
)

type fontsModel struct {
	colActive int
	leftRow   int
	rightRow  int
	startRow  int
}

func (m *model) updateFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxViewHeight := 10

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			if m.fonts.colActive == 0 && m.fonts.leftRow > 0 {
				m.fonts.leftRow--
			} else if m.fonts.colActive == 1 && m.fonts.rightRow > 0 {
				m.fonts.rightRow--
			}

		case "j", "down":
			if m.fonts.colActive == 0 && m.fonts.leftRow < len(fontItemsTemp)-1 {
				m.fonts.leftRow++
			} else if m.fonts.colActive == 1 && m.fonts.rightRow < len(fontItems)-1 {
				m.fonts.rightRow++
			}

		case "h", "left":
			if m.fonts.colActive == 1 {
				m.fonts.colActive = 0
			}

		case "l", "right":
			if m.fonts.colActive == 0 {
				m.fonts.colActive = 1
			}

		case "q", "esc":
			return m, tea.Quit
		}

		// left column scrolling
		if m.fonts.colActive == 0 {
			if m.fonts.leftRow >= m.fonts.startRow+maxViewHeight {
				m.fonts.startRow = m.fonts.leftRow - maxViewHeight + 1
			} else if m.fonts.leftRow < m.fonts.startRow {
				m.fonts.startRow = m.fonts.leftRow
			}
		}
	}
	return m, nil
}

func (m *model) fontsView() string {
	maxViewHeight := 10
	leftWidth := max(min(m.width-30, 40), 20)

	var leftChoices []string
	start := m.fonts.startRow
	end := min(start+maxViewHeight, len(fontItemsTemp))

	for i := start; i < end; i++ {
		isSelected := (m.fonts.colActive == 0 && m.fonts.leftRow == i)
		leftChoices = append(leftChoices, fontItem(fontItemsTemp[i], isSelected, leftWidth-4))
	}

	var rightChoices []string
	for i, label := range fontItems {
		isSelected := (m.fonts.colActive == 1 && m.fonts.rightRow == i)
		rightChoices = append(rightChoices, fontItem(label, isSelected, 0))
	}

	leftCol := lipgloss.NewStyle().PaddingRight(4).Border(lipgloss.NormalBorder()).Width(leftWidth).Height(maxViewHeight).Render(lipgloss.JoinVertical(lipgloss.Left, leftChoices...))
	rightCol := lipgloss.NewStyle().PaddingTop(1).Render(lipgloss.JoinVertical(lipgloss.Left, rightChoices...))

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Your fonts"),
		columns,
		wordwrap.String(helpStyle.Render("(←↓↑→: move | d: remove font | q: quit)"), m.width),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

var fontItems = []string{
	"(1) Custom",
	"(2) From system",
}

// temporary for testing
var fontItemsTemp = []string{
	"Font 1dsadasdasadsasdadsdasdas",
	"Font 2",
	"Font 3",
	"Font 4",
	"Font 5",
	"Font 6",
	"Font 7",
	"Font 8",
	"Font 9",
	"Font 10",
	"Font 11",
	"Font 12",
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
