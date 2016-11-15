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
	var bypass map[string]bool

	files, err := ResolveDependencies( Cfg.TypeFiles, Cfg.SqlDirPath + "types" )
	if err != nil { return err }

	pgtypes := make( []*Type, 0 )
	for i := len( files ) - 1 ; i >= 0 ; i-- {
		file := files[ i ]
		pgtype := Type{}
		pgtype.Path = file
		pgtypes = append( pgtypes, &pgtype )

		err = DownPass( &pgtype, pgtype.Path )
		if err != nil {
			successfulCount--
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			bypass[ pgtype.Path ] = true
		}
	}

	for i := len( pgtypes ) - 1 ; i >= 0 ; i-- {
		pgtype := pgtypes[ i ]
		if _, ignore := bypass[ pgtype.Path ] ; ! ignore {
			err = UpPass( pgtype, pgtype.Path )
			if err != nil {
				successfulCount--
				errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			}
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
