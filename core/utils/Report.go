package utils

import (
	"fmt"
)

// Report pretty prints the result of an import.
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
