package core

import (
	"fmt"
)

// Sanity type encapsulates requirement checks.
type sanity struct{}

// check makes sure the fs is ready to be used.
func (checker *sanity) check() (err error) {
	if err = checker.directoryExists(); err != nil {
		return err
	}
	if err = checker.typedDirExists(); err != nil {
		return err
	}
	scanFiles(&conf)
	if err = checker.sqlFilesPresent(); err != nil {
		return err
	}

	return
}

// directoryExists checks the provided sql directory is indeed a directory.
func (checker *sanity) directoryExists() (err error) {
	if !isDir(conf.sqlDirPath) {
		return fmt.Errorf("%s is not a directory", conf.sqlDirPath)
	}

	return
}

// typedDirExists makes sure that at least one of functions/, views/, triggers/, types/ exists.
func (checker *sanity) typedDirExists() (err error) {
	directories := make([]string, 0)

	for _, typedDir := range []string{"functions", "triggers", "types", "views"} {
		path := conf.sqlDirPath + "/" + typedDir
		if isDir(path) {
			directories = append(directories, path)
		}
	}

	if len(directories) == 0 {
		return fmt.Errorf("No functions/, triggers/, types/ or views/ directory found in %s", conf.sqlDirPath)
	}

	return
}

// sqlFilesPresent checks there are source file.
//
// No need to process any further if there are no sql files to load.
func (checker *sanity) sqlFilesPresent() (err error) {
	if len(conf.functionFiles)+len(conf.triggerFiles)+len(conf.viewFiles) == 0 {
		return fmt.Errorf("Didn't find any sql file in %s", conf.sqlDirPath)
	}

	return
}
