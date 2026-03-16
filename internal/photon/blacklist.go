package photon

import "strings"

func IsPathBlacklisted(path string) bool {
	splitPath := strings.Split(path, "/")

	if len(splitPath) < 2 {
		return true
	}

	switch splitPath[1] {
	case "dev":
		return true
	case "core":
		return true
	}

	return false
}
