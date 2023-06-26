package core

import (
	"fmt"

	"gitlab.com/oelmekki/pgrebase/core/utils"
)

// Sanity type encapsulates requirement checks.
type sanity struct{}

// Check makes sure the fs is ready to be used.
func (checker *sanity) Check() (err error) {
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
	if !utils.IsDir(conf.SqlDirPath) {
		return fmt.Errorf("%s is not a directory", conf.SqlDirPath)
	}

	return
}

// typedDirExists makes sure that at least one of functions/, views/, triggers/, types/ exists.
func (checker *sanity) typedDirExists() (err error) {
	directories := make([]string, 0)

	for _, typedDir := range []string{"functions", "triggers", "types", "views"} {
		path := conf.SqlDirPath + "/" + typedDir
		if utils.IsDir(path) {
			directories = append(directories, path)
		}
	}

	if len(directories) == 0 {
		return fmt.Errorf("No functions/, triggers/, types/ or views/ directory found in %s", conf.SqlDirPath)
	}

	return
}

// sqlFilesPresent checks there are source file.
//
// No need to process any further if there are no sql files to load.
func (checker *sanity) sqlFilesPresent() (err error) {
	if len(conf.FunctionFiles)+len(conf.TriggerFiles)+len(conf.ViewFiles) == 0 {
		return fmt.Errorf("Didn't find any sql file in %s", conf.SqlDirPath)
	}

	return
}
