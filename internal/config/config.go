package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func Init() *Config {
	cfg := &Config{}

	homeDir := fs.GetHomeDir()

	configDir := filepath.Join(homeDir, ".config", "ebfm")
	configFile := filepath.Join(configDir, "config.toml")

	defaultLogPath := filepath.Join(homeDir, "ebfm.log")

	fs.EnsureDirectory(configFile)

	info, err := os.Stat(configFile)
	if err == nil && info.Size() == 0 {
		defaultConfig(configFile, defaultLogPath)
	}

	if _, err := toml.DecodeFile(configFile, cfg); err != nil {
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

func (cfg *Config) Apply() {
	homeDir := fs.GetHomeDir()
	configDir := filepath.Join(homeDir, ".config", "ebfm")
	configFile := filepath.Join(configDir, "config.toml")

	f, err := os.Create(configFile)
	if err != nil {
		log.Fatalln("Failed to create config file: ", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		log.Fatalln("Failed to encode toml: ", err)
	}
}
