package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

type systemFontModel struct {
	selection     int
	startRow      int
	fonts         []config.Font
	filteredFonts []config.Font
	searchInput   textinput.Model
}

func initialSystemFontModel() systemFontModel {
	return systemFontModel{
		selection:   0,
		searchInput: createInput("font name", "Search: ", 50),
	}
}

func (m *model) updateSystemFontSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxViewHeight := 12
	numItems := len(m.systemFont.filteredFonts)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// keypresses inside focused input
		if m.systemFont.searchInput.Focused() {
			switch msg.String() {
			case "up", "down":
				m.systemFont.searchInput.Blur()
				return m, nil
			case "enter":
				m.systemFont.searchInput.Blur()
				return m, nil
			case "esc":
				m.systemFont.searchInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.systemFont.searchInput, cmd = m.systemFont.searchInput.Update(msg)

			m.filterFonts()

			return m, cmd
		}

		// normal keypresses
		switch msg.String() {
		case "k", "up":
			if m.systemFont.selection > 0 {
				m.systemFont.selection--
			}
		case "j", "down":
			if m.systemFont.selection < numItems-1 {
				m.systemFont.selection++
			}
		case "p":
			return m, m.previewFont(m.systemFont.filteredFonts[m.systemFont.selection].Path)
		case "q", "esc":
			m.state = stateFonts
			return m, nil
		case "f":
			return m, m.systemFont.searchInput.Focus()
		case "enter", " ":
			if len(m.systemFont.filteredFonts) > 0 {
			  selected := m.systemFont.filteredFonts[m.systemFont.selection]
			  m.fonts.collection.Add(selected)
			  m.state = stateFonts
			}
			return m, nil
		}

		// font list scrolling
		if m.systemFont.selection >= m.systemFont.startRow+maxViewHeight {
			m.systemFont.startRow = m.systemFont.selection - maxViewHeight + 1
		} else if m.systemFont.selection < m.systemFont.startRow {
			m.systemFont.startRow = m.systemFont.selection
		}
	}
	return m, nil
}

func (m *model) systemFontView() string {
	maxViewHeight := 12
	fontsWidth := max(min(m.width-10, 50), 20)

	var choices []string
	start := m.systemFont.startRow
	end := min(start+maxViewHeight, len(m.systemFont.filteredFonts))

	for i := start; i < end; i++ {
		isSelected := m.systemFont.selection == i
		choices = append(choices, systemFontItem(m.systemFont.filteredFonts[i].Name + " " + m.systemFont.filteredFonts[i].Style, isSelected, fontsWidth-4))
	}

	fontList := lipgloss.NewStyle().PaddingRight(4).Border(lipgloss.NormalBorder()).Width(fontsWidth).Height(maxViewHeight).Render(lipgloss.JoinVertical(lipgloss.Left, choices...))
	search := lipgloss.NewStyle().Width(fontsWidth).Align(lipgloss.Center).Render(m.systemFont.searchInput.View())
	content := lipgloss.JoinVertical(lipgloss.Center, fontList, search)

	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Add system font"),
		content,
		wordwrap.String(helpStyle.Render("(↑↓: move | p: preview | f: search | enter/space: add | q: back)"), m.width),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
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

func (m *model) filterFonts() {
	value := strings.ToLower(m.systemFont.searchInput.Value())
	if value == "" {
		m.systemFont.filteredFonts = m.systemFont.fonts
	} else {
		var filtered []config.Font
		for _, f := range m.systemFont.fonts {
			if strings.Contains(strings.ToLower(f.Name), value) {
				filtered = append(filtered, f)
			}
		}
		m.systemFont.filteredFonts = filtered
	}

	// selection might go further that expected so we reset it to 0
	m.systemFont.selection = 0
	m.systemFont.startRow = 0
}
