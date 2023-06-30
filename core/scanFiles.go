package core

import (
	"os"
	"path/filepath"
)

// sourceWalker type encapsulates fs walking functions.
type sourceWalker struct {
	config *config
}

// process loads all source files paths.
func (walker *sourceWalker) process() {
	walker.config.functionFiles = walker.findFunctions()
	walker.config.triggerFiles = walker.findTriggers()
	walker.config.viewFiles = walker.findViews()

	return
}

// findFunctions finds path of function files.
func (walker *sourceWalker) findFunctions() (paths []string) {
	return walker.sqlFilesIn(walker.config.sqlDirPath + "/functions")
}

// findTriggers dinds path of trigger files.
func (walker *sourceWalker) findTriggers() (paths []string) {
	return walker.sqlFilesIn(walker.config.sqlDirPath + "/triggers")
}

// findViews finds path of view files.
func (walker *sourceWalker) findViews() (paths []string) {
	return walker.sqlFilesIn(walker.config.sqlDirPath + "/views")
}

// sqlFilesIn walks a directory to find sql files.
func (walker *sourceWalker) sqlFilesIn(path string) (paths []string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if isSqlFile(path) {
			paths = append(paths, path)
		}

		return nil
	})

	return
}

// scanFiles scans sql directory for sql files.
func scanFiles(cfg *config) {
	cfg.functionFiles = make([]string, 0)
	cfg.triggerFiles = make([]string, 0)
	cfg.viewFiles = make([]string, 0)

	walker := sourceWalker{config: cfg}
	walker.process()
}
