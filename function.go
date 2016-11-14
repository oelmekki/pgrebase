package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
)

/*
 * Load or reload all functions found in FS.
 */
func LoadFunctions() ( err error ) {
	successfulCount := len( Cfg.FunctionFiles )
	errors := make( []string, 0 )


	for _, file := range Cfg.FunctionFiles {
		function := Function{}
		function.Path = file

		err = ProcessUnit( &function, function.Path )
		if err != nil {
			successfulCount--;
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
		}
	}

	Report( "functions", successfulCount, len( Cfg.FunctionFiles ), errors )

	return
}

type Function struct {
	CodeUnit
	Signature string
}

/*
 * Load function definition from file
 */
func ( function *Function ) Load() ( err error ) {
	definition, err := ioutil.ReadFile( function.Path )
	if err != nil { return err }
	function.Definition = string( definition )

	return
}

/*
 * Parse function for name and signature
 */
func ( function *Function ) Parse() ( err error ) {
	signatureFinder := regexp.MustCompile( `(?is)CREATE(?:\s+OR\s+REPLACE)?\s+FUNCTION\s+(\S+?)\((.*?)\)` )
	subMatches := signatureFinder.FindStringSubmatch( function.Definition )

	if len( subMatches ) < 3 {
		return fmt.Errorf( "Can't find a function in %s", function.Path )
	}

	function.Name = subMatches[1]

	if function.parseSignature {
		function.Signature = subMatches[2]
	} else {
		function.Signature, function.previousExists, err = function.previousSignature()
		if err != nil { return err }
	}

	return
}

/*
 * Drop existing function from pg
 */
func ( function *Function ) Drop() ( err error ) {
	if function.previousExists {
		return function.CodeUnit.Drop( `DROP FUNCTION IF EXISTS ` + function.Name + `(` + function.Signature + `)` )
	}

	return
}

/*
 * Create the function in pg
 */
func ( function *Function ) Create() ( err error ) {
	return function.CodeUnit.Create( function.Definition )
}

/*
 * Retrieve old signature from function in database, if any
 */
func ( function *Function ) previousSignature() ( signature string, exists bool, err error ) {
	rows, err := Query( `SELECT pg_get_functiondef(oid) FROM pg_proc WHERE proname = $1`, function.Name )
	if err != nil { return }
	defer rows.Close()

	if rows.Next() {
		exists = true

		var body string
		if err = rows.Scan( &body ) ; err != nil { return }
		oldFunction := Function{ CodeUnit: CodeUnit{ Definition: body, parseSignature: true } }
		if err = oldFunction.Parse() ; err != nil { return }
		signature = oldFunction.Signature
	}

	return
}
