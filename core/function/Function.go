package function

import (
	"fmt"
	"github.com/oelmekki/pgrebase/core/codeunit"
	"github.com/oelmekki/pgrebase/core/connection"
	"io/ioutil"
	"regexp"
	"strings"
)

// Function is the code unit for functions.
type Function struct {
	codeunit.CodeUnit
	Signature string // function signature, unparsed
}

// Parse parses function for name and signature.
func (function *Function) Parse() (err error) {
	signatureFinder := regexp.MustCompile(`(?is)CREATE(?:\s+OR\s+REPLACE)?\s+FUNCTION\s+(\S+?)\((.*?)\)`)
	subMatches := signatureFinder.FindStringSubmatch(function.Definition)

	if len(subMatches) < 3 {
		return fmt.Errorf("Can't find a function in %s", function.Path)
	}

	function.Name = subMatches[1]

	if function.ParseSignature {
		function.Signature = subMatches[2]
	} else {
		function.Signature, function.PreviousExists, err = function.previousSignature()
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

// Load loads function definition from file.
func (function *Function) Load() (err error) {
	definition, err := ioutil.ReadFile(function.Path)
	if err != nil {
		return err
	}
	function.Definition = string(definition)

	return
}

// Drop removes existing function from pg.
func (function *Function) Drop() (err error) {
	if function.PreviousExists {
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
	rows, err := connection.Query(`SELECT pg_get_functiondef(oid) FROM pg_proc WHERE proname = $1`, function.Name)
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
		oldFunction := Function{CodeUnit: codeunit.CodeUnit{Definition: body, ParseSignature: true}}
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
