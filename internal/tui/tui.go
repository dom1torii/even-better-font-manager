package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"

	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/dom1torii/even-better-font-manager/internal/platform/fonts"
)

var (
	selectionStyle     = lipgloss.NewStyle().Background(lipgloss.Color("8"))
	titleStyle         = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	statusStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	statusOkStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	statusWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	helpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).PaddingBottom(1)
	inputStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Italic(true)
	emptyListSignStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
	disabledStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
)

type sessionState int

const (
	statePath sessionState = iota
	stateStart
	stateFonts
	stateSystemFont
	stateCustomFont
)

type model struct {
	cfg    *config.Config
	state  sessionState
	height int
	width  int

	path       pathModel
	start      startModel
	fonts      fontsModel
	systemFont systemFontModel
	customFont customFontModel

	csPath     string
	chosenFont font

	Err      error
	Quitting bool
}

type font struct {
	fonts.Font
	Error error
}

func InitialModel(cfg *config.Config) *model {
	// if don't have cs2 path set in the config, start with path chooser state
	initialState := stateStart
	if cfg.General.CS2Path == "" {
		initialState = statePath
	}

	return &model{
		cfg:        cfg,
		state:      initialState,
		path:       initialPathModel(),
		start:      initialStartModel(),
		fonts:      initialFontsModel(),
		systemFont: initialSystemFontModel(),
		customFont: initialCustomFontModel(),
		Quitting:   false,
	}
}

func createInput(placeholder, prompt string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = prompt
	ti.Width = width
	ti.CharLimit = 156
	ti.Blur()
	ti.PromptStyle = lipgloss.NewStyle()
	ti.PlaceholderStyle = inputStyle

	return ti
}

type menuItem struct {
	render   func(m *model) string
	status   func(m *model) string
	action   func(m *model) tea.Cmd
	disabled func(m *model) bool
}

// if width is 0 we just dont trucate the item
func menuItemStyle(label string, isSelected bool, width int) string {
	truncated := label
	if width > 0 {
		truncated = truncate.StringWithTail(label, uint(width), "...")
	}
	if isSelected {
		return selectionStyle.Render(truncated)
	}
	return truncated
}
