package core

import (
	"github.com/oelmekki/pgrebase/core/config"
	"github.com/oelmekki/pgrebase/core/utils"
	"os"
	"path/filepath"
)

// sourceWalker type encapsulates fs walking functions.
type sourceWalker struct {
	Config *config.Config
}

// Process loads all source files paths.
func (walker *sourceWalker) Process() {
	walker.Config.FunctionFiles = walker.findFunctions()
	walker.Config.TriggerFiles = walker.findTriggers()
	walker.Config.TypeFiles = walker.findTypes()
	walker.Config.ViewFiles = walker.findViews()

	return
}

// findFunctions finds path of function files.
func (walker *sourceWalker) findFunctions() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/functions")
}

// findTriggers dinds path of trigger files.
func (walker *sourceWalker) findTriggers() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/triggers")
}

// findTypes finds path of type files.
func (walker *sourceWalker) findTypes() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/types")
}

// findViews finds path of view files.
func (walker *sourceWalker) findViews() (paths []string) {
	return walker.sqlFilesIn(walker.Config.SqlDirPath + "/views")
}

// sqlFilesIn walks a directory to find sql files.
func (walker *sourceWalker) sqlFilesIn(path string) (paths []string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if utils.IsSqlFile(path) {
			paths = append(paths, path)
		}

		return nil
	})

	return
}
