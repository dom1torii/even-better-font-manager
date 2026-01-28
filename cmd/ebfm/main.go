package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/dom1torii/even-better-font-manager/internal/tui"
)

func main() {
	cfg := config.Init()

	p := tea.NewProgram(tui.InitialModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalln("Error starting bubbletea: ", err)
		os.Exit(1)
	}
}
