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
	chosenFont chosenFont

	Err      error
	Quitting bool
}

type chosenFont struct {
	Name string
	Path string
	Error error
}

func InitialModel(cfg *config.Config) *model {
	cs2PathInput := createInput("/path/to/cs2/", "Path: ", 20)
	systemFontSearchInput := createInput("font name", "Search: ", 50)
	customFontPathInput := createInput("/path/to/font", "Path: ", 20)

	// if don't have cs2 path set in the config, start with path chooser state
	initialState := stateStart
	if cfg.General.CS2Path == "" {
		initialState = statePath
	}

	return &model{
		cfg: cfg,
		state: initialState,
		path: pathModel{
			selection: 0,
			pathInput: cs2PathInput,
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
			searchInput: systemFontSearchInput,
		},
		customFont: customFontModel{
			selection: 0,
			pathInput: customFontPathInput,
		},
		Quitting: false,
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
