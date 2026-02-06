package tui

import (
	"github.com/dom1torii/even-better-font-manager/internal/platform/fonts"
)

type csPathMsg string

type pathConfirmedMsg struct{}
type fontPreviewMsg struct{}

type systemFontsMsg []fonts.SystemFont
