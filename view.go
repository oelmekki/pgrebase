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
		view := View{ Path: file }
		err = view.Process()
		if err != nil {
			successfulCount--;
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
		}
	}

	ViewsReport( successfulCount, errors )

	return
}

/*
 * Pretty print of view loading result
 */
func ViewsReport( successfulCount int, errors []string ) {
	fmt.Printf( "Loaded %d views", successfulCount )
	if successfulCount < len( Cfg.ViewFiles ) {
		fmt.Printf( " - %d with error", len( Cfg.ViewFiles ) - successfulCount )
	}
	fmt.Printf( "\n" )

	for _, err := range errors {
		fmt.Printf( err )
	}
}

type View struct {
	Path            string
	Name            string
	Definition      string
	previousExists  bool
	parseSignature  bool
}

/*
 * Create or update a view found in FS
 */
func ( view *View ) Process() ( err error ) {
	errFmt := "  error while loading %s\n  %v\n"

	if err = view.Load() ; err != nil { return fmt.Errorf( errFmt, view.Path, err ) }
	if err = view.Parse() ; err != nil { return fmt.Errorf( errFmt, view.Path, err ) }
	if err = view.Drop() ; err != nil { return fmt.Errorf( errFmt, view.Path, err ) }
	if err = view.Create() ; err != nil { return fmt.Errorf( errFmt, view.Path, err ) }

	return
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
	rows, err := Query( `DROP VIEW IF EXISTS ` + view.Name + ` CASCADE` )
	if err != nil { return err }
	rows.Close()
	return
}

/*
 * Create the view in pg
 */
func ( view *View ) Create() ( err error ) {
	rows, err := Query( view.Definition )
	if err != nil { return err }
	rows.Close()

	return
}
