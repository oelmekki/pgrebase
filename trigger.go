package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
)

/*
 * Load or reload all triggers found in FS.
 */
func LoadTriggers() ( err error ) {
	successfulCount := len( Cfg.TriggerFiles )
	errors := make( []string, 0 )
	var bypass map[string]bool

	files, err := ResolveDependencies( Cfg.TriggerFiles, Cfg.SqlDirPath + "triggers" )
	if err != nil { return err }

	triggers := make( []*Trigger, 0 )
	for i := len( files ) - 1 ; i >= 0 ; i-- {
		file := files[ i ]
		trigger := Trigger{}
		trigger.Path = file
		triggers = append( triggers, &trigger )

		err = DownPass( &trigger, trigger.Path )
		if err != nil {
			successfulCount--
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			bypass[ trigger.Path ] = true
		}
	}

	for i := len( triggers ) - 1 ; i >= 0 ; i-- {
		trigger := triggers[ i ]
		if _, ignore := bypass[ trigger.Path ] ; ! ignore {
			err = UpPass( trigger, trigger.Path )
			if err != nil {
				successfulCount--
				errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			}
		}
	}

	Report( "triggers", successfulCount, len( Cfg.TriggerFiles ), errors )

	return
}

type Trigger struct {
	CodeUnit
	Table           string
	Function        Function
}

/*
 * Load trigger definition from file
 */
func ( trigger *Trigger ) Load() ( err error ) {
	definition, err := ioutil.ReadFile( trigger.Path )
	if err != nil { return err }
	trigger.Definition = string( definition )

	return
}

/*
 * Parse trigger for name and signature
 */
func ( trigger *Trigger ) Parse() ( err error ) {
	triggerFinder := regexp.MustCompile( `(?is)CREATE(?:\s+CONSTRAINT)?\s+TRIGGER\s+(\S+).*?ON\s+(\S+)` )
	subMatches := triggerFinder.FindStringSubmatch( trigger.Definition )

	if len( subMatches ) < 3 {
		return fmt.Errorf( "Can't find a trigger in %s", trigger.Path )
	}

	trigger.Function = Function{ CodeUnit: CodeUnit{ Definition: trigger.Definition, Path: trigger.Path } }
	trigger.Function.Parse()

	trigger.Name = subMatches[1]
	trigger.Table = subMatches[2]

	return
}

/*
 * Drop existing trigger from pg
 */
func ( trigger *Trigger ) Drop() ( err error ) {
	err = trigger.CodeUnit.Drop( `DROP TRIGGER IF EXISTS ` + trigger.Name + ` ON ` + trigger.Table )
	if err != nil { return err }

	return trigger.Function.Drop()
}

/*
 * Create the trigger in pg
 */
func ( trigger *Trigger ) Create() ( err error ) {
	return trigger.CodeUnit.Create( trigger.Definition )
}

