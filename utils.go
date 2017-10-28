package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
)

/*
 * Check if file exists and is a directory
 */
func IsDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

/*
 * Check if provided file is an sql file (only check for extension)
 */
func IsSqlFile(filePath string) bool {
	isSqlFile := regexp.MustCompile(`.*\.sql$`)
	return isSqlFile.MatchString(filePath)
}

/*
 * Check if file is hidden
 */
func IsHiddenFile(filePath string) bool {
	basename := path.Base(filePath)
	return string(basename[0]) == "."
}

/*
 * Pretty print the result of an import
 */
func Report(name string, successCount, totalCount int, errors []string) {
	fmt.Printf("Loaded %d %s", successCount, name)
	if successCount < totalCount {
		fmt.Printf(" - %d with error", totalCount-successCount)
	}
	fmt.Printf("\n")

	for _, err := range errors {
		fmt.Printf(err)
	}
}
