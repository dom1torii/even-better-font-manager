package fonts

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/image/font/sfnt"
)

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
