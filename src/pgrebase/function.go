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
			if err != nil { return err }
		}
	}

	return
}

type Function struct {
	Path            string
	Name            string
	Signature       string
	Definition      string
	parseSignature  bool
}

/*
 * Create or update a function found in FS
 */
func ( function *Function ) Process( errChan chan error ) {
	var err error

	if err = function.Load() ; err != nil { errChan <- err ; return }
	if err = function.Parse() ; err != nil { errChan <- err ; return }
	if err = function.Drop() ; err != nil { errChan <- err ; return }
	if err = function.Create() ; err != nil { errChan <- err ; return }

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
	signatureFinder := regexp.MustCompile( `(?is)CREATE(?: OR REPLACE) FUNCTION (\S+?)\((.*?)\)` )
	subMatches := signatureFinder.FindStringSubmatch( function.Definition )

	if len( subMatches ) < 3 {
		return fmt.Errorf( "Can't find a function in %s", function.Path )
	}

	function.Name = subMatches[1]

	if function.parseSignature {
		function.Signature = subMatches[2]
	} else {
		function.Signature, err = function.previousSignature()
		if err != nil { return err }
	}

	return
}

/*
 * Drop existing function from pg
 */
func ( function *Function ) Drop() ( err error ) {
	if len( function.Signature ) > 0 {
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
func ( function *Function ) previousSignature() ( signature string, err error ) {
	rows, err := Query( `SELECT pg_get_functiondef(oid) FROM pg_proc WHERE proname = $1`, function.Name )
	if err != nil { return signature, err }
	defer rows.Close()

	if rows.Next() {
		var body string
		if err = rows.Scan( &body ) ; err != nil { return signature, err }
		oldFunction := Function{ Definition: body, parseSignature: true }
		if err = oldFunction.Parse() ; err != nil { return signature, err }
		return oldFunction.Signature, nil
	}

	return
}
