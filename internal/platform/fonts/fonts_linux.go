//go:build linux

package fonts

import (
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/image/font/sfnt"
)

type SystemFont struct {
	Path string
	Name string
}

func GetSystemFonts() []SystemFont {
	cmd := exec.Command("fc-list")
	output, err := cmd.Output()
  if err != nil {
    log.Fatalln("Failed to run fc-list: ", err)
  }

  lines := strings.Split(string(output), "\n")

  var fonts []SystemFont
  for _, line := range lines {
 		if strings.TrimSpace(line) == "" {
      continue
    }

 		parts := strings.SplitN(line, ":", 2)
   	path := strings.TrimSpace(parts[0])
    secondPart := strings.Split(parts[1], ":")

    // one font (for some reason) can have multiple names separated my comma,
    // so we use last one since it usually describes font the best
    // EDIT: i need to do something else since sometimes it has the same name
    // for multiple fonts of different weight, example Noto Sans
    // EDIT2: maybe just append whatever is after style=?
    fontNames := strings.TrimSpace(secondPart[0])
    nameList := strings.Split(fontNames, ",")

    finalName := strings.TrimSpace(nameList[len(nameList)-1])

    fonts = append(fonts, SystemFont{
      Path: path,
      Name: finalName,
    })
  }

  sort.SliceStable(fonts, func(i, j int) bool {
    return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name)
  })

  return fonts
}

func Preview(path string) {
	cmd := exec.Command("fontforge", path)
	if err := cmd.Run(); err != nil {
		log.Println("Failed to open fontforge", err)
	}
}

func GetName(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	f, err := sfnt.Parse(fontBytes)
	if err != nil {
		log.Fatalln("Failed to parse font: ", err)
	}

	name, err := f.Name(nil, sfnt.NameIDFull)
	if err != nil {
		log.Fatalln("Failed to find font name: ", err)
	}

	log.Println(name)
	return name, nil
}
