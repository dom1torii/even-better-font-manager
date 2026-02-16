//go:build linux

package cs2path

import (
	"log"
	"os"
	"path/filepath"

	"github.com/andygrunwald/vdf"
)

func Detect() string {
	// on linux, steam always (probably) adds a file to your home directory
	// that is called .steam that has folders that symlink to your actual
	// steam installation, so we use it to get libraryfolders.vdf and
	// parse it to get our cs installation folder
	libraryFoldersPath := os.ExpandEnv("$HOME/.steam/steam/steamapps/libraryfolders.vdf")
	f, err := os.Open(libraryFoldersPath)
	if err != nil {
		log.Println("Failed to open libraryfolders.vdf: ", err)
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		log.Fatalln("Failed to parse vdf file: ", err)
	}

	csAppId := "730"
	// we try to find "path" of the folder that has "730" inside its "apps"
	if libraryFolders, ok := m["libraryfolders"].(map[string]any); ok {
		for _, folderData := range libraryFolders {
			if folder, ok := folderData.(map[string]any); ok {
				if apps, ok := folder["apps"].(map[string]any); ok {
					if _, found := apps[csAppId]; found {
						log.Printf("CS2 Path: %s", filepath.Join(folder["path"].(string), "steamapps/common/Counter-Strike Global Offensive"))
						return filepath.Join(folder["path"].(string), "steamapps/common/Counter-Strike Global Offensive")
					}
				}
			}
		}
	}

	log.Println("No CS2 path detected")
	return ""
}
