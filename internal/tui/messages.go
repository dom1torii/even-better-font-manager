package tui

import (
	"github.com/dom1torii/even-better-font-manager/internal/platform/systemfonts"
)

type csPathMsg string

type pathConfirmedMsg struct{}

type systemFontsMsg []systemfonts.SystemFont
