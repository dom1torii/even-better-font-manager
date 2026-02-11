package config

import (
	"log"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/dom1torii/even-better-font-manager/internal/fs"
)

type Collection struct {
	Fonts []Font `json:"fonts"`
}

type Font struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (c *Collection) Save() {
  path := getCollectionPath()

  data, err := json.MarshalIndent(c, "", "  ")
  if err != nil {
    log.Println("Failed to encode collection json: ", err)
    return
  }

  if err := os.WriteFile(path, data, 0644); err != nil {
    log.Println("Failed to write collection file: ", err)
  }
}

func LoadCollection() *Collection {
  path := getCollectionPath()
  col := &Collection{Fonts: []Font{}}

  f, err := os.ReadFile(path)
  if err != nil {
  	// if file doesn't exist, it means no fonts were added yet
    if !os.IsNotExist(err) {
      log.Println("Failed to read collection file: ", err)
    }
    return col
  }

  if len(f) > 0 {
    if err := json.Unmarshal(f, col); err != nil {
      log.Println("Failed to decode collection json: ", err)
    }
  }
  return col
}

func getCollectionPath() string {
  homeDir := fs.GetHomeDir()
  return filepath.Join(homeDir, ".config", "ebfm", "collection.json")
}

func (c *Collection) List() []Font {
  return c.Fonts
}

func (c *Collection) Add(f Font) {
	// dont add if already exists
  for _, existing := range c.Fonts {
    if existing.Path == f.Path {
      return
    }
  }

  c.Fonts = append(c.Fonts, f)
  c.Save()
}

func (c *Collection) Remove(index int) {
  if index < 0 || index >= len(c.Fonts) {
    return
  }
  c.Fonts = append(c.Fonts[:index], c.Fonts[index+1:]...)
  c.Save()
}
