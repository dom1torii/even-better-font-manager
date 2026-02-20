package filechooser

import (
	"log"

	"github.com/ncruces/zenity"
)

func Open(chooserType string, title string) string {
	var path string
	var err error

	if chooserType == "directory" {
		path, err = zenity.SelectFile(
			zenity.Title(title),
			zenity.Directory(),
		)
	} else {
		path, err = zenity.SelectFile(
			zenity.Title(title),
		)
	}

	if err != nil {
		if err == zenity.ErrCanceled {
			log.Println("File chooser was closed")
		} else {
			log.Println("Zenity failed: ", err)
		}
		return ""
	}

	return path
}
