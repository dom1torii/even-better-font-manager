package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"

	"github.com/dom1torii/even-better-font-manager/internal/fs"
)

type Config struct {
	General GeneralConfig `toml:"general"`
	Log     LogConfig     `toml:"logging"`

	// cli mode
}

type GeneralConfig struct {
	CS2Path string `toml:"cs2_path"`
}

type LogConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
}

type Collection struct {
	Fonts []Font `json:"fonts"`
}

type Font struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func Init() *Config {
	cfg := &Config{}

	homeDir := fs.GetHomeDir()

	configPath := getConfigPath()
	collectionPath := getCollectionPath()

	defaultLogPath := filepath.Join(homeDir, "ebfm.log")

	fs.EnsureDirectory(configPath)
	fs.EnsureDirectory(collectionPath)

	info, err := os.Stat(configPath)
	if err == nil && info.Size() == 0 {
		defaultConfig(configPath, defaultLogPath)
	}

	if _, err := toml.DecodeFile(configPath, cfg); err != nil {
		log.Fatalln("Failed to decode config: ", err)
	}

	logFlag := pflag.BoolP("log", "l", cfg.Log.Enabled, "Enable logging. Default path: {homedir}/ebfm.log")
	logPath := pflag.String("logpath", getFlag(cfg.Log.Path, defaultLogPath), "Specify custom log file path.")

	pflag.Parse()

	cfg.Log.Path = *logPath

	isLogFlagSet := pflag.Lookup("log").Changed
	isPathFlagSet := pflag.Lookup("logpath").Changed

	// we can just use --logpath instead of using both -l and --logpath to enable logging + change path
	if isLogFlagSet || isPathFlagSet || cfg.Log.Enabled {
		if isLogFlagSet && !*logFlag {
			cfg.Log.Enabled = false
		} else {
			cfg.Log.Enabled = true
		}
	}

	if cfg.Log.Enabled {
		fs.EnsureDirectory(cfg.Log.Path)
		initLogger(cfg.Log.Path)
		log.Println("Started logging at: ", cfg.Log.Path)
	} else {
		log.SetOutput(io.Discard)
	}

	return cfg
}

func initLogger(loc string) {
	f, err := os.OpenFile(loc, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open log file: ", err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func defaultConfig(path, defaultLog string) {
	// fix windows paths
	defaultLog = strings.ReplaceAll(defaultLog, "\\", "\\\\")

	content := []byte(strings.Join([]string{
		"[logging]",
		"enabled = false",
		"path = \"" + defaultLog + "\"",
	}, "\n"))

	if err := os.WriteFile(path, content, 0644); err != nil {
		log.Fatalln("Failed to write default config: ", err)
	}
}

func getFlag(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func (cfg *Config) Save() {
	path := getConfigPath()

	f, err := os.Create(path)
	if err != nil {
		log.Fatalln("Failed to create config file: ", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		log.Fatalln("Failed to encode toml: ", err)
	}
}

// need a separate file for collection later uwu
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

func getConfigPath() string {
 	homeDir := fs.GetHomeDir()
  return filepath.Join(homeDir, ".config", "ebfm", "config.toml")
}

func (c *Collection) List() []Font {
  return c.Fonts
}

func (c *Collection) Add(f Font) {
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
