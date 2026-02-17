package tui

import (
	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/dom1torii/even-better-font-manager/internal/platform/cs2path"
)

type csPathMsg cs2path.CS2Path
type getCsPathMsg string

type pathConfirmedMsg struct{}
type fontPreviewMsg struct{}

type fontCollectionMsg *config.Collection

type systemFontsMsg []config.Font

type customFontPathMsg string
type customFontMsg font

type fontConfigMsg struct{}
