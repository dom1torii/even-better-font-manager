package fonts

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dom1torii/even-better-font-manager/internal/config"
	"github.com/dom1torii/even-better-font-manager/internal/fs"
	"golang.org/x/image/font/sfnt"
)

type Font struct {
	Name  string
	Style string
	Path  string
}

func GetCollection(cfg *config.Config) []Font {
	fontsDir := filepath.Join(cfg.General.CS2Path, "game/core/panorama/fonts/ebfm-custom")

	var fonts []Font

	if _, err := os.Stat(fontsDir); os.IsNotExist(err) {
    return fonts
  }

	fontFiles, err := os.ReadDir(fontsDir)
	if err != nil {
		log.Println("Failed to read custom font files: ", err)
		return fonts
	}

	for _, fontFile := range fontFiles {
		extension := strings.ToLower(filepath.Ext(fontFile.Name()))
		if extension == ".ttf" || extension == ".otf" {
			path := filepath.Join(fontsDir, fontFile.Name())
			f, err := ParseFont(path)
			if err != nil {
				log.Println("Something happened when parsing font: ", err)
				continue
			}
			fonts = append(fonts, f)
		}
	}

	return fonts
}

func AddToCollection(cfg *config.Config, path string) {
	fontsDir := filepath.Join(cfg.General.CS2Path, "game/core/panorama/fonts/ebfm-custom")
	fs.EnsureDirectory(fontsDir)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalln("Failed to open font: ", err)
	}
	defer f.Close()

	fileName := filepath.Base(path)
	destPath := filepath.Join(fontsDir, fileName)

	dest, err := os.Create(destPath)
	if err != nil {
		log.Fatalln("Failed to create destination file: ", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, f)
	if err != nil {
		log.Fatalln("Failed to copy font into custom fonts directory: ", err)
	}
}

func RemoveFromCollection(path string) {
	fontsDir := filepath.Dir(path)

	if err := os.RemoveAll(path); err != nil {
		log.Println("Failed to delete font: ", err)
	}

	fontFiles, err := os.ReadDir(fontsDir)
  if err != nil {
  	log.Println("Failed to read custom font files: ", err)
   	return
  }

  if len(fontFiles) == 0 {
  	if err := os.RemoveAll(fontsDir); err != nil {
   		log.Println("Failed to remove empty custom font directory: ", err)
   	}
  }
}

func ParseFont(path string) (Font, error) {
	fontBytes, err := os.ReadFile(path)
  if err != nil {
    return Font{}, err
  }

  f, err := sfnt.Parse(fontBytes)
  if err != nil {
  	return Font{}, err
  }

 	name, err := f.Name(nil, sfnt.NameIDFull)
	if err != nil {
		return Font{}, err
	}

	style, err := f.Name(nil, sfnt.NameIDSubfamily)
	if err != nil {
		return Font{
			Name: name,
			Path: path,
		}, nil
	}

	return Font{
		Name: name,
		Style: style,
		Path: path,
	}, nil
}
