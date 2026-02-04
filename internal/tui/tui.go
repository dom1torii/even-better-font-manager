package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

var (
	selectionStyle        = lipgloss.NewStyle().Background(lipgloss.Color("8"))
	titleStyle            = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	statusStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	statusOkStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	statusWarningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	helpStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).PaddingBottom(1)
	inputStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Italic(true)
)

type sessionState int

const (
	statePath sessionState = iota
	stateStart
	stateFonts
	stateSystemFont
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

	csPath     string
	chosenFont chosenFont

	Err      error
	Quitting bool
}

type chosenFont struct {
	Name string
	Path string
}

func InitialModel(cfg *config.Config) *model {
	pi := textinput.New()
	pi.Placeholder = "/path/to/cs2/"
	pi.Blur()
	pi.Prompt = "Path: "
	pi.CharLimit = 156
	pi.Width = 20
	pi.PromptStyle = lipgloss.NewStyle()
	pi.PlaceholderStyle = inputStyle

	si := textinput.New()
	si.Placeholder = "font name"
	si.Blur()
	si.Prompt = "Search: "
	si.CharLimit = 156
	si.Width = 50
	si.PromptStyle = lipgloss.NewStyle()
	si.PlaceholderStyle = inputStyle

	initialState := stateStart
	if cfg.General.CS2Path == "" {
		initialState = statePath
	}

	return &model{
		cfg: cfg,
		state: initialState,

		path: pathModel{
			selection: 0,
			pathInput: pi,
		},
		start: startModel{
			selection: 0,
		},
		fonts: fontsModel{
			leftSelection: 0,
			rightSelection: 0,
		},
		systemFont: systemFontModel{
			selection: 0,
			searchInput: si,
		},

		Quitting: false,
	}
}
