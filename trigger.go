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


	for _, file := range Cfg.TriggerFiles {
		trigger := Trigger{ Path: file }
		err = trigger.Process()
		if err != nil {
			successfulCount--;
			errors = append( errors, fmt.Sprintf( "%v\n", err ) )
		}
	}

	TriggersReport( successfulCount, errors )

	return
}

/*
 * Pretty print of trigger loading result
 */
func TriggersReport( successfulCount int, errors []string ) {
	fmt.Printf( "Loaded %d triggers", successfulCount )
	if successfulCount < len( Cfg.TriggerFiles ) {
		fmt.Printf( " - %d with error", len( Cfg.TriggerFiles ) - successfulCount )
	}
	fmt.Printf( "\n" )

	for _, err := range errors {
		fmt.Printf( err )
	}
}

type Trigger struct {
	Path            string
	Name            string
	Table           string
	Definition      string
	Function        Function
	previousExists  bool
	parseSignature  bool
}

/*
 * Create or update a trigger found in FS
 */
func ( trigger *Trigger ) Process() ( err error ) {
	errFmt := "  error while loading %s\n  %v\n"

	if err = trigger.Load() ; err != nil { return fmt.Errorf( errFmt, trigger.Path, err ) }
	if err = trigger.Parse() ; err != nil { return fmt.Errorf( errFmt, trigger.Path, err ) }
	if err = trigger.Drop() ; err != nil { return fmt.Errorf( errFmt, trigger.Path, err ) }
	if err = trigger.Create() ; err != nil { return fmt.Errorf( errFmt, trigger.Path, err ) }

	return
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
	triggerFinder := regexp.MustCompile( `(?is)CREATE(?: CONSTRAINT)? TRIGGER (\S+).*?ON (\S+)` )
	subMatches := triggerFinder.FindStringSubmatch( trigger.Definition )

	if len( subMatches ) < 3 {
		return fmt.Errorf( "Can't find a trigger in %s", trigger.Path )
	}

	trigger.Function = Function{ Definition: trigger.Definition, Path: trigger.Path }
	trigger.Function.Parse()

	trigger.Name = subMatches[1]
	trigger.Table = subMatches[2]

	return
}

/*
 * Drop existing trigger from pg
 */
func ( trigger *Trigger ) Drop() ( err error ) {
	rows, err := Query( `DROP TRIGGER IF EXISTS ` + trigger.Name + ` ON ` + trigger.Table + ` CASCADE` )
	if err != nil { fmt.Printf( "error on drop : DROP TRIGGER IF EXISTS " + trigger.Name + " ON " + trigger.Table + " CASCADE\n" ); return err }
	rows.Close()
	if err = trigger.Function.Drop() ; err != nil { return err }

	return
}

/*
 * Create the trigger in pg
 */
func ( trigger *Trigger ) Create() ( err error ) {
	rows, err := Query( trigger.Definition )
	if err != nil { return err }
	rows.Close()

	return
}

