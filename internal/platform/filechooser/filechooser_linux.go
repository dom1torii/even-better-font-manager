//go:build linux

package filechooser

import (
	"log"
	"strings"
	"os/exec"
)

// dependency: zenity
func ChoosePath() string {
	cmd := exec.Command("zenity",
		"--file-selection",
		"--directory",
		"--title=Choose your Counter-Strike Global Offensive/ folder")
	output, err := cmd.Output()
	if err != nil {
	  log.Println("File chooser was closed or zenity failed: ", err)
	}
	path := strings.TrimSpace(string(output))

	return path
}
