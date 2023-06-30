package core

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"

	_ "github.com/lib/pq"
)

// isDir checks if file exists and is a directory.
func isDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isHiddenFile checks if file is hidden.
func isHiddenFile(filePath string) bool {
	basename := path.Base(filePath)
	return string(basename[0]) == "."
}

// isSqlFile checks if provided file is an sql file (only check for extension).
func isSqlFile(filePath string) bool {
	sqlFile := regexp.MustCompile(`.*\.sql$`)
	return sqlFile.MatchString(filePath)
}

// report pretty prints the result of an import.
func report(name string, successCount, totalCount int, errors []string) {
	if os.Getenv("QUIET") != "true" || successCount < totalCount {
		fmt.Printf("Loaded %d %s", successCount, name)

		if successCount < totalCount {
			fmt.Printf(" - %d with error", totalCount-successCount)
		}

		fmt.Printf("\n")
	}

	for _, err := range errors {
		fmt.Printf(err)
	}
}

// query is a wrapper for sql.Query, meant to be the main query interface.
func query(q string, parameters ...interface{}) (rows *sql.Rows, err error) {
	var co *sql.DB

	co, err = sql.Open("postgres", conf.databaseUrl)
	if err != nil {
		fmt.Printf("can't connect to database : %v\n", err)
		return rows, err
	}
	defer co.Close()

	rows, err = co.Query(q, parameters...)
	if err != nil {
		return rows, err
	}

	return
}
