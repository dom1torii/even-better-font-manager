package tui

import (
	"github.com/dom1torii/even-better-font-manager/internal/platform/cs2path"
	"github.com/dom1torii/even-better-font-manager/internal/platform/fonts"
)

type csPathMsg cs2path.CS2Path
type getCsPathMsg string

type pathConfirmedMsg struct{}
type fontPreviewMsg struct{}

type fontCollectionMsg []fonts.Font
type addFontToCollectionMsg struct{}
type removeFontFromCollectionMsg struct{}

type systemFontsMsg []fonts.Font

type customFontPathMsg string
type customFontMsg font

type writeFontConfigMsg string
type resetFontConfigMsg string
type clearApplyStatusMsg struct{}
type clearResetStatusMsg struct{}
