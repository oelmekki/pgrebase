package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
)

/*
 * Load or reload all views found in FS.
 */
func LoadViews() ( err error ) {
	successfulCount := len( Cfg.ViewFiles )
	errors := make( []string, 0 )


	for _, file := range Cfg.ViewFiles {
		view := View{}
		view.Path = file

		err = ProcessUnit( &view, view.Path )
		if err != nil {
			successfulCount--;
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
		}
	}

	Report( "triggers", successfulCount, len( Cfg.ViewFiles ), errors )

	return
}

type View struct {
	CodeUnit
}

/*
 * Load view definition from file
 */
func ( view *View ) Load() ( err error ) {
	definition, err := ioutil.ReadFile( view.Path )
	if err != nil { return err }
	view.Definition = string( definition )

	return
}

/*
 * Parse view for name
 */
func ( view *View ) Parse() ( err error ) {
	nameFinder := regexp.MustCompile( `(?is)CREATE(?:\s+OR\s+REPLACE)?\s+VIEW\s+(\S+)` )
	subMatches := nameFinder.FindStringSubmatch( view.Definition )

	if len( subMatches ) < 2 {
		return fmt.Errorf( "Can't find a view in %s", view.Path )
	}

	view.Name = subMatches[1]

	return
}

/*
 * Drop existing view from pg
 */
func ( view *View ) Drop() ( err error ) {
	return view.CodeUnit.Drop( `DROP VIEW IF EXISTS ` + view.Name + ` CASCADE` )
}

/*
 * Create the view in pg
 */
func ( view *View ) Create() ( err error ) {
	return view.CodeUnit.Create( view.Definition )
}
