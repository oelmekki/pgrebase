package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
)

/*
 * Load or reload all types found in FS.
 */
func LoadTypes() ( err error ) {
	successfulCount := len( Cfg.TypeFiles )
	errors := make( []string, 0 )

	files, err := ResolveDependencies( Cfg.TypeFiles, Cfg.SqlDirPath + "types" )
	if err != nil { return err }

	for _, file := range files {
		pgtype := Type{}
		pgtype.Path = file

		err = ProcessUnit( &pgtype, pgtype.Path )
		if err != nil {
			successfulCount--;
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
		}
	}

	Report( "types", successfulCount, len( Cfg.TypeFiles ), errors )

	return
}

type Type struct {
	CodeUnit
}

/*
 * Load type definition from file
 */
func ( pgtype *Type ) Load() ( err error ) {
	definition, err := ioutil.ReadFile( pgtype.Path )
	if err != nil { return err }
	pgtype.Definition = string( definition )

	return
}

/*
 * Parse type for name
 */
func ( pgtype *Type ) Parse() ( err error ) {
	nameFinder := regexp.MustCompile( `(?is)CREATE\s+TYPE\s+(\S+)` )
	subMatches := nameFinder.FindStringSubmatch( pgtype.Definition )

	if len( subMatches ) < 2 {
		return fmt.Errorf( "Can't find a type in %s", pgtype.Path )
	}

	pgtype.Name = subMatches[1]

	return
}

/*
 * Drop existing type from pg
 */
func ( pgtype *Type ) Drop() ( err error ) {
	return pgtype.CodeUnit.Drop( `DROP TYPE IF EXISTS ` + pgtype.Name )
}

/*
 * Create the type in pg
 */
func ( pgtype *Type ) Create() ( err error ) {
	return pgtype.CodeUnit.Create( pgtype.Definition )
}
