package util

import (
	"path"
	"strings"
)

// GetNameBasedOnUrl will return the base name of the given URL.
func GetBaseName(s string) string {
	normalizedName := strings.ReplaceAll(s, "\\", "/")
	return path.Base(normalizedName)
}
