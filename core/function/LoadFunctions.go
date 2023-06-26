package function

import (
	"fmt"

	"gitlab.com/oelmekki/pgrebase/core/codeunit"
	"gitlab.com/oelmekki/pgrebase/core/resolver"
	"gitlab.com/oelmekki/pgrebase/core/utils"
)

// LoadFunctions loads or reload all functions found in FS.
func LoadFunctions() (err error) {
	successfulCount := len(conf.FunctionFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := resolver.ResolveDependencies(conf.FunctionFiles, conf.SqlDirPath+"functions")
	if err != nil {
		return err
	}

	functions := make([]*Function, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		function := Function{}
		function.Path = file
		functions = append(functions, &function)

		err = codeunit.DownPass(&function, function.Path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[function.Path] = true
		}
	}

	for i := len(functions) - 1; i >= 0; i-- {
		function := functions[i]
		if _, ignore := bypass[function.Path]; !ignore {
			err = codeunit.UpPass(function, function.Path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	utils.Report("functions", successfulCount, len(conf.FunctionFiles), errors)

	return
}
