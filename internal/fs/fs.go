package fs

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dom1torii/even-better-font-manager/internal/platform/perms"
)

func GetHomeDir() string {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if usr, err := user.Lookup(sudoUser); err == nil {
			return usr.HomeDir
		}
	}

	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// should work for windows idk
	usr, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to get user:", err)
	}
	return usr.HomeDir
}

func EnsureDirectory(path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalln("Failed to create directory: ", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalln("Failed to create file: ", err)
		}
		if err := f.Close(); err != nil {
			log.Fatalln("Failed to close file: ", err)
		}
	}

	// on linux, we won't be able to edit files we create if we run with sudo,
	// so we need to change ownership of the files
	perms.FixPermissions(dir)
	perms.FixPermissions(path)
}
