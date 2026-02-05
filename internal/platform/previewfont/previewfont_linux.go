//go:build linux

package previewfont

import (
	"log"
	"os/exec"
)

func PreviewFont(path string) {
	cmd := exec.Command("fontforge", path)
	if err := cmd.Run(); err != nil {
		log.Println("Failed to open fontforge", err)
	}
}
