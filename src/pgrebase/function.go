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

	for i := 0 ; i < len( Cfg.FunctionFiles ) ; i += Cfg.MaxConnection {
		next := Cfg.MaxConnection
		rest := len( Cfg.FunctionFiles ) - i
		if rest < next { next = rest }

		errChan := make( chan error )

		for _, file := range Cfg.FunctionFiles[i:i+next] {
			function := Function{ Path: file }
			go function.Process( errChan )
		}

		for j := 0 ; j < i + next ; j++ {
			err = <-errChan
			if err != nil {
				successfulCount--;
				errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			}
		}
	}

	FunctionsReport( successfulCount, errors )

	return
}

/*
 * Pretty print of function loading result
 */
func FunctionsReport( successfulCount int, errors []string ) {
	fmt.Printf( "Loaded %d functions", successfulCount )
	if successfulCount < len( Cfg.FunctionFiles ) {
		fmt.Printf( " - %d with error", len( Cfg.FunctionFiles ) - successfulCount )
	}
	fmt.Printf( "\n" )

	for _, err := range errors {
		fmt.Printf( err )
	}
}

type Function struct {
	Path            string
	Name            string
	Signature       string
	Definition      string
	previousExists  bool
	parseSignature  bool
}

/*
 * Create or update a function found in FS
 */
func ( function *Function ) Process( errChan chan error ) {
	var err error

	errFmt := "  error while loading %s\n  %v\n"

	if err = function.Load() ; err != nil { errChan <- fmt.Errorf( errFmt, function.Path, err ) ; return }
	if err = function.Parse() ; err != nil { errChan <- fmt.Errorf( errFmt, function.Path, err ) ; return }
	if err = function.Drop() ; err != nil { errChan <- fmt.Errorf( errFmt, function.Path, err ) ; return }
	if err = function.Create() ; err != nil { errChan <- fmt.Errorf( errFmt, function.Path, err ) ; return }

	errChan <- err
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
	signatureFinder := regexp.MustCompile( `(?is)CREATE(?: OR REPLACE)? FUNCTION (\S+?)\((.*?)\)` )
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
		rows, err := Query( `DROP FUNCTION IF EXISTS ` + function.Name + `(` + function.Signature + `)` )
		if err != nil { return err }
		rows.Close()
	}
	return
}

/*
 * Create the function in pg
 */
func ( function *Function ) Create() ( err error ) {
	rows, err := Query( function.Definition )
	if err != nil { return err }
	rows.Close()

	return
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
		oldFunction := Function{ Definition: body, parseSignature: true }
		if err = oldFunction.Parse() ; err != nil { return }
		signature = oldFunction.Signature
	}

	return
}
