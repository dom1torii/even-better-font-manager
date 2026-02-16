package tui

import (
	"github.com/dom1torii/even-better-font-manager/internal/config"
)

type csPathMsg string
type getCsPathMsg string
type verifyCsPathMsg error

type pathConfirmedMsg struct{}
type fontPreviewMsg struct{}

type fontCollectionMsg *config.Collection

type systemFontsMsg []config.Font

type customFontPathMsg string
type customFontMsg font

type fontConfigMsg struct{}
