package utils

import (
	"os"
)

// IsDir checks if file exists and is a directory.
func IsDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
