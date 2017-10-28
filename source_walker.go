package main

import (
	"os"
	"path/filepath"
)

type SourceWalker struct {
	Config *Config
}

/*
 * Load all source files paths
 */
func (walker *SourceWalker) Process() {
	walker.Config.FunctionFiles = walker.findFunctions()
	walker.Config.TriggerFiles = walker.findTriggers()
	walker.Config.TypeFiles = walker.findTypes()
	walker.Config.ViewFiles = walker.findViews()

	return
}

/*
 * Find path of function files
 */
func (walker *SourceWalker) findFunctions() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/functions")
}

/*
 * Find path of trigger files
 */
func (walker *SourceWalker) findTriggers() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/triggers")
}

/*
 * Find path of type files
 */
func (walker *SourceWalker) findTypes() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/types")
}

/*
 * Find path of view files
 */
func (walker *SourceWalker) findViews() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/views")
}

/*
 * Walk a directory to find sql files
 */
func (walker *SourceWalker) sqlFilesIn(path string) (paths []string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if IsSqlFile(path) {
			paths = append(paths, path)
		}

		return nil
	})

	return
}
