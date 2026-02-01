package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

var (
	selectionStyle        = lipgloss.NewStyle().Background(lipgloss.Color("8"))
	checkedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	checkedSelectionStyle = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("2")).Bold(true)
	crossedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	crossedSelectionStyle = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("1")).Bold(true)
	goodPingStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	badPingStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	blockedPingStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
	timedoutPingStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	titleStyle            = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	statusStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	statusOkStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	statusWarningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	helpStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).PaddingBottom(1)
	modeAllowStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	modeBlockStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
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

	csPath string

	Err      error
	Quitting bool
}

func InitialModel(cfg *config.Config) *model {
	// change these to use prompt instead of what is used rn
	pi := textinput.New()
	pi.Placeholder = "/path/to/cs2/"
	pi.Blur()
	pi.Prompt = ""
	pi.CharLimit = 156
	pi.Width = 20
	pi.PromptStyle = lipgloss.NewStyle()
	pi.PlaceholderStyle = inputStyle

	si := textinput.New()
	si.Placeholder = "font name"
	si.Blur()
	si.Prompt = ""
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
