//go:build windows

package fonts

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/image/font/sfnt"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

// nothing is tested, i just blindly wrote ts hoping it will work on windows. need to boot up a vm and im too lazy rn

func GetSystemFonts() []config.Font {
	var fonts []config.Font

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

    fontName, style, err := GetName(path)
    if err != nil {
      continue
    }

    fonts = append(fonts, config.Font{
      Name:  fontName,
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
  cmd := exec.Command("cmd", "/c", "start", "", path)
  if err := cmd.Run(); err != nil {
    log.Println("Failed to preview font: ", err)
  }
}

func GetName(path string) (string, string, error) {
	if path == "" {
		return "", "", fmt.Errorf("No path provided")
	}

	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	f, err := sfnt.Parse(fontBytes)
	if err != nil {
		return "", "", err
	}

	name, err := f.Name(nil, sfnt.NameIDFull)
	if err != nil {
		return "", "", err
	}

	style, err := f.Name(nil, sfnt.NameIDSubfamily)
	if err != nil {
		// just return name if we cant find style
		return name, "", nil
	}

	log.Println(name + " " + style)
	return name, style, nil
}
