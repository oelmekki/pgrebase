package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// LoadFunctions loads or reload all functions found in FS.
func LoadFunctions() (err error) {
	successfulCount := len(Cfg.FunctionFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := ResolveDependencies(Cfg.FunctionFiles, Cfg.SqlDirPath+"functions")
	if err != nil {
		return err
	}

	functions := make([]*Function, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		function := Function{}
		function.Path = file
		functions = append(functions, &function)

		err = DownPass(&function, function.Path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[function.Path] = true
		}
	}

	for i := len(functions) - 1; i >= 0; i-- {
		function := functions[i]
		if _, ignore := bypass[function.Path]; !ignore {
			err = UpPass(function, function.Path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	Report("functions", successfulCount, len(Cfg.FunctionFiles), errors)

	return
}

// Function is the code unit for functions.
type Function struct {
	CodeUnit
	Signature string // function signature, unparsed
}

// Load loads function definition from file.
func (function *Function) Load() (err error) {
	definition, err := ioutil.ReadFile(function.Path)
	if err != nil {
		return err
	}
	function.Definition = string(definition)

	return
}

// Parse parses function for name and signature.
func (function *Function) Parse() (err error) {
	signatureFinder := regexp.MustCompile(`(?is)CREATE(?:\s+OR\s+REPLACE)?\s+FUNCTION\s+(\S+?)\((.*?)\)`)
	subMatches := signatureFinder.FindStringSubmatch(function.Definition)

	if len(subMatches) < 3 {
		return fmt.Errorf("Can't find a function in %s", function.Path)
	}

	function.Name = subMatches[1]

	if function.parseSignature {
		function.Signature = subMatches[2]
	} else {
		function.Signature, function.previousExists, err = function.previousSignature()
		if err != nil {
			return err
		}
	}

	err = function.removeDefaultFromSignature()
	if err != nil {
		return
	}

	return
}

// Drop removes existing function from pg.
func (function *Function) Drop() (err error) {
	if function.previousExists {
		return function.CodeUnit.Drop(`DROP FUNCTION IF EXISTS ` + function.Name + `(` + function.Signature + `)`)
	}

	return
}

// Create adds the function in pg.
func (function *Function) Create() (err error) {
	return function.CodeUnit.Create(function.Definition)
}

// previousSignature retrieves old signature from function in database, if any.
func (function *Function) previousSignature() (signature string, exists bool, err error) {
	rows, err := Query(`SELECT pg_get_functiondef(oid) FROM pg_proc WHERE proname = $1`, function.Name)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Next() {
		exists = true

		var body string
		if err = rows.Scan(&body); err != nil {
			return
		}
		oldFunction := Function{CodeUnit: CodeUnit{Definition: body, parseSignature: true}}
		if err = oldFunction.Parse(); err != nil {
			return
		}
		signature = oldFunction.Signature
	}

	return
}

// removeDefaultFromSignature sanitizes function signature.
//
// `DROP FUNCTION` raises error if the signature contains `DEFAULT` values for
// parameters, so we must remove them.
func (function *Function) removeDefaultFromSignature() (err error) {
	anyDefault, err := regexp.MatchString("(?i)DEFAULT", function.Signature)
	if err != nil {
		return
	}

	if anyDefault {
		arguments := strings.Split(function.Signature, ",")
		newArgs := make([]string, 0)

		for _, arg := range arguments {
			arg = strings.Replace(arg, " DEFAULT ", " default ", -1)
			newArg := strings.Split(arg, " default ")[0]
			newArgs = append(newArgs, newArg)
		}

		function.Signature = strings.Join(newArgs, ",")
	}

	return
}
