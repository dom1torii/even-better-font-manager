package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/truncate"
)

type systemFontModel struct {
	selection int
	startRow  int
}

func (m *model) updateSystemFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
  maxViewHeight := 10
  numItems := len(systemFontTemp)

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "k", "up":
      if m.systemFont.selection > 0 {
        m.systemFont.selection--
      }
    case "j", "down":
      if m.systemFont.selection < numItems-1 {
        m.systemFont.selection++
      }
    case "q", "esc":
    	m.state = stateFonts
      return m, nil
    }

    // font list scrolling
    if m.systemFont.selection >= m.systemFont.startRow + maxViewHeight {
      m.systemFont.startRow = m.systemFont.selection - maxViewHeight + 1
    } else if m.systemFont.selection < m.systemFont.startRow {
      m.systemFont.startRow = m.systemFont.selection
    }
  }
  return m, nil
}

func (m *model) systemFontView() string {
  maxViewHeight := 10
  fontsWidth := max(min(m.width-30, 40), 20)

  var choices []string
  start := m.systemFont.startRow
  end := min(start+maxViewHeight, len(systemFontTemp))

  for i := start; i < end; i++ {
    isSelected := m.systemFont.selection == i
    choices = append(choices, systemFontItem(systemFontTemp[i], isSelected, fontsWidth-4))
  }

  fontList := lipgloss.NewStyle().PaddingRight(4).Border(lipgloss.NormalBorder()).Width(fontsWidth).Height(maxViewHeight).Render(lipgloss.JoinVertical(lipgloss.Left, choices...))
  search := "Search: *input*"
  content := lipgloss.JoinVertical(lipgloss.Center, fontList, search)

  view := fmt.Sprintf(
    "%s\n\n%s\n\n%s",
    titleStyle.Render("Add system font"),
    content,
    helpStyle.Render("(↑↓: move | p: preview | enter/space: add | q: back)"),
  )

  return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

// temporary for testing
var systemFontTemp = []string{
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

func systemFontItem(label string, isSelected bool, width int) string {
	truncated := label
	if width > 0 {
		truncated = truncate.StringWithTail(label, uint(width), "...")
  }
	if isSelected {
		return selectionStyle.Render(truncated)
	}
	return truncated
}
