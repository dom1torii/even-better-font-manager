//go:build linux

package filechooser

import (
	"log"
	"os/exec"
	"strings"
)

// dependency: zenity
func Open(chooserType string, title string) string {
	cmd := exec.Command("zenity",
		"--file-selection",
		"--"+chooserType,
		"--title="+title)
	output, err := cmd.Output()
	if err != nil {
		log.Println("File chooser was closed or zenity failed: ", err)
	}
	path := strings.TrimSpace(string(output))

	return path
}
