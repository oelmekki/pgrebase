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

	for i := 0 ; i < len( Cfg.TriggerFiles ) ; i += Cfg.MaxConnection {
		next := Cfg.MaxConnection
		rest := len( Cfg.TriggerFiles ) - i
		if rest < next { next = rest }

		errChan := make( chan error )

		for _, file := range Cfg.TriggerFiles[i:i+next] {
			trigger := Trigger{ Path: file }
			go trigger.Process( errChan )
		}

		for j := 0 ; j < i + next ; j++ {
			err = <-errChan
			if err != nil {
				successfulCount--;
				errors = append( errors, fmt.Sprintf( "%v\n", err ) )
			}
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
func ( trigger *Trigger ) Process( errChan chan error ) {
	var err error

	errFmt := "  error while loading %s\n  %v\n"

	if err = trigger.Load() ; err != nil { errChan <- fmt.Errorf( errFmt, trigger.Path, err ) ; return }
	if err = trigger.Parse() ; err != nil { errChan <- fmt.Errorf( errFmt, trigger.Path, err ) ; return }
	if err = trigger.Drop() ; err != nil { errChan <- fmt.Errorf( errFmt, trigger.Path, err ) ; return }
	if err = trigger.Create() ; err != nil { errChan <- fmt.Errorf( errFmt, trigger.Path, err ) ; return }

	errChan <- err
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

