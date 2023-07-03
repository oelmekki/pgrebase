package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// function is the code unit for functions.
type function struct {
	codeUnit
	signature string // function signature, unparsed
}

// loadFunctions loads or reload all functions found in FS.
func loadFunctions() (err error) {
	successfulCount := len(conf.functionFiles)
	errors := make([]string, 0)
	bypass := make(map[string]bool)

	files, err := resolveDependencies(conf.functionFiles, conf.sqlDirPath+"functions")
	if err != nil {
		return err
	}

	functions := make([]*function, 0)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		f := function{}
		f.path = file
		functions = append(functions, &f)

		err = downPass(&f, f.path)
		if err != nil {
			successfulCount--
			errors = append(errors, fmt.Sprintf("%v\n", err))
			bypass[f.path] = true
		}
	}

	for i := len(functions) - 1; i >= 0; i-- {
		f := functions[i]
		if _, ignore := bypass[f.path]; !ignore {
			err = upPass(f, f.path)
			if err != nil {
				successfulCount--
				errors = append(errors, fmt.Sprintf("%v\n", err))
			}
		}
	}

	report("functions", successfulCount, len(conf.functionFiles), errors)

	return
}

// parse parses function for name and signature.
func (function *function) parse() (err error) {
	signatureFinder := regexp.MustCompile(`(?is)CREATE(?:\s+OR\s+REPLACE)?\s+FUNCTION\s+(\S+?)\((.*?)\)`)
	subMatches := signatureFinder.FindStringSubmatch(function.definition)

	if len(subMatches) < 3 {
		return fmt.Errorf("Can't find a function in %s", function.path)
	}

	function.name = subMatches[1]

	if function.parseSignature {
		function.signature = subMatches[2]
	} else {
		function.signature, function.previousExists, err = function.previousSignature()
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

// load loads function definition from file.
func (f *function) load() (err error) {
	definition, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}
	f.definition = string(definition)

	return
}

// drop removes existing function from pg.
func (f *function) drop() (err error) {
	if f.previousExists {
		return f.codeUnit.drop(`DROP FUNCTION IF EXISTS ` + f.name + `(` + f.signature + `)`)
	}

	return
}

// create adds the function in pg.
func (f *function) create() (err error) {
	return f.codeUnit.create(f.definition)
}

// previousSignature retrieves old signature from function in database, if any.
func (f *function) previousSignature() (signature string, exists bool, err error) {
	rows, err := query(`SELECT pg_get_functiondef(oid) FROM pg_proc WHERE proname = $1`, f.name)
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
		oldFunction := function{codeUnit: codeUnit{definition: body, parseSignature: true}}
		if err = oldFunction.parse(); err != nil {
			return
		}
		signature = oldFunction.signature
	}

	return
}

// removeDefaultFromSignature sanitizes function signature.
//
// `DROP FUNCTION` raises error if the signature contains `DEFAULT` values for
// parameters, so we must remove them.
func (f *function) removeDefaultFromSignature() (err error) {
	anyDefault, err := regexp.MatchString("(?i)DEFAULT", f.signature)
	if err != nil {
		return
	}

	if anyDefault {
		arguments := strings.Split(f.signature, ",")
		newArgs := make([]string, 0)

		for _, arg := range arguments {
			arg = strings.Replace(arg, " DEFAULT ", " default ", -1)
			newArg := strings.Split(arg, " default ")[0]
			newArgs = append(newArgs, newArg)
		}

		f.signature = strings.Join(newArgs, ",")
	}

	return
}
