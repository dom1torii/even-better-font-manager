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
)

type sessionState int

const (
	statePath sessionState = iota
	stateStart
	stateFonts
)

type model struct {
	cfg    *config.Config
	state  sessionState
	height int
	width  int

	path pathModel
	start startModel
	fonts fontsModel

	csPath    string

	Err      error
	Quitting bool
}

func InitialModel(cfg *config.Config) *model {
	ti := textinput.New()
	ti.Placeholder = "/path/to/cs2/"
	ti.Blur()
	ti.Prompt = ""
	ti.CharLimit = 156
	ti.Width = 20
	ti.PromptStyle = lipgloss.NewStyle()

	initialState := stateStart
	if cfg.General.CS2Path == "" {
		initialState = statePath
	}

	return &model{
		cfg: cfg,
		state: initialState,

		path: pathModel{
			pathInput: ti,
		},

		Quitting: false,
	}
}
