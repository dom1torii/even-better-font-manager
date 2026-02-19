//go:build windows

package fonts

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// nothing is tested, i just blindly wrote ts hoping it will work on windows. need to boot up a vm and im too lazy rn

func GetSystemFonts() []Font {
	var fonts []Font

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.QUERY_VALUE)
	if err != nil {
		log.Println("No System fonts found: ", err)
		return nil
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		log.Println("Failed to read font names: ", err)
	}

	fontsDir := filepath.Join(os.Getenv("WINDIR"), "Fonts")

	for _, name := range names {
		filename, _, err := k.GetStringValue(name)
		if err != nil {
			continue
		}

		path := filename
		if !filepath.IsAbs(filename) {
			path = filepath.Join(fontsDir, filename)
		}

		f, err := ParseFont(path)
		if err != nil {
			continue
		}

		fonts = append(fonts, Font{
			Name:  f.Name,
			Style: f.Style,
			Path:  path,
		})
	}

	sort.SliceStable(fonts, func(i, j int) bool {
		return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name)
	})

	return fonts
}

func Preview(path string) {
	cmd := exec.Command("cmd", "/c", "start", "", path)
	if err := cmd.Run(); err != nil {
		log.Println("Failed to preview font: ", err)
	}
}
