package cs2path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// we verify cs2 path by checking if it contains Counter-Strike Global Offensive or pak01_dir.vpk inside it
func Verify(path string) error {
	normalizedPath := filepath.ToSlash(path)
	pakFile := filepath.Join(path, "game/csgo/pak01_dir.vpk")

	_, err := os.Stat(pakFile)
	if strings.Contains(normalizedPath, "steamapps/common/Counter-Strike Global Offensive") && err == nil {
		return nil
	}

	// maybe change this error to something else, not sure what yet
	return fmt.Errorf("CS2 path isn't correct")
}

// probably wanna try to make something that will try to fix cs2 path for you if its wrong,
// since we expect Counter-Strike Global Offensive folder specifically.
