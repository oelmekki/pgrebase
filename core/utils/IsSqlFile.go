package utils

import (
	"regexp"
)

// IsSqlFile checks if provided file is an sql file (only check for extension).
func IsSqlFile(filePath string) bool {
	isSqlFile := regexp.MustCompile(`.*\.sql$`)
	return isSqlFile.MatchString(filePath)
}
