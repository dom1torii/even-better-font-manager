//go:build windows

package detectpath

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

// not tested yet, too lazy to open windows vm, surely it works
func DetectCS2Path() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\cs2`, registry.QUERY_VALUE)
	if err != nil {
		log.Println("No CS2 path detected")
	}
	defer k.Close()

	s, _, err := k.GetStringValue("installpath")
	if err != nil {
		log.Println("No CS2 path detected")
	}

	return s
}
