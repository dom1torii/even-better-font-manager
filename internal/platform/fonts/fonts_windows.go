//go:build windows

package fonts

import ()

type SystemFont struct {
	Path string
	Name string
}

func GetSystemFonts() {

}

func Preview(path string) {

}

func GetName(path string) (string, error) {
	return "", nil
}
