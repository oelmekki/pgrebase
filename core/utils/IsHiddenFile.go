package utils

import (
	"path"
)

// IsHiddenFile checks if file is hidden.
func IsHiddenFile(filePath string) bool {
	basename := path.Base(filePath)
	return string(basename[0]) == "."
}
