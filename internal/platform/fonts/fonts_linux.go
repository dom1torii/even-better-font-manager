//go:build linux

package fonts

import (
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

func GetSystemFonts() []config.Font {
	cmd := exec.Command("fc-list")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalln("Failed to run fc-list: ", err)
	}

	lines := strings.Split(string(output), "\n")

	var fonts []config.Font
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		path := strings.TrimSpace(parts[0])

		// one font (for some reason) can have multiple names separated my comma,
		// so we use last one since it usually describes font the best
		fontNames := strings.TrimSpace(parts[1])
		nameList := strings.Split(fontNames, ",")
		finalName := strings.TrimSpace(nameList[len(nameList)-1])

		stylesPart := strings.TrimPrefix(strings.TrimSpace(parts[2]), "style=")
		styleList := strings.Split(stylesPart, ",")
		style := strings.TrimSpace(styleList[0])

		fonts = append(fonts, config.Font{
			Name:  finalName,
			Style: style,
			Path:  path,
		})
	}

	sort.SliceStable(fonts, func(i, j int) bool {
		return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name)
	})

	return fonts
}

func Preview(path string) {
	cmd := exec.Command("gnome-font-viewer", path)
	if err := cmd.Run(); err != nil {
		log.Println("Failed to preview font: ", err)
	}
}
