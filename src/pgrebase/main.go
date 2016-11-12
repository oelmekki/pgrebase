package main

import (
	"fmt"
	"os"
)

var Cfg Config

/*
 * Usage :
 *   DATABASE_URL=url pgrebase [-w] sql_dir/
 */
func ParseConfig() {
	if err := Cfg.Parse() ; err != nil {
		fmt.Printf( "Error: %v\n", err )
		Usage()
	}
}

/*
 * Expected target structure:
 *
 * <sql_dir>/
 *   functions/
 *   triggers/
 *   views/
 *
 * At least one of functions/triggers/views/ should exist.
 *
 */
func CheckSanity() {
	sanity := Sanity{}
	if err := sanity.Check() ; err != nil {
		fmt.Printf( "Error: %v\n", err )
		Usage()
	}
}

/*
 * Give user a chance to know what went wrong
 */
func Usage() {
	usage := `
USAGE:
	DATABASE_URL=url pgrebase [-w] [-n <num>] <sql_directory>

PgRebase is a tool that allows you to easily handle your postgres codebase for
functions, triggers and views.

Your expected to provide a postgresql connection url as DATABASE_URL and
a sql directory as <sql_directory>.

<sql_directory> should be structured that way:
	<sql_directory>/
	├── functions/
	├── triggers/
	└── views/

At least one of functions/triggers/views/ should exist.

OPTIONS:
	-w: enter watch mode.
		In watch mode, pgrebase will keep watching for file changes and will
		automatically reload your sql code when it happens.

	-n <num>: max concurrency
	  Specifies a number of maximum files to be loaded concurrently.
		Default: pg max connections / 2
	`;

	fmt.Println( usage )
	os.Exit(1)
}

/*
 * Start the actual work
 */
func Process() ( err error ) {
	if err = LoadFunctions() ; err != nil { return err }
	// if err = LoadTriggers() ; err != nil { return err }
	// if err = LoadViews() ; err != nil { return err }

	return
}

func main() {
	ParseConfig()
	CheckSanity()

	if err := Process() ; err != nil {
		fmt.Printf( "Error: %v\n", err )
		os.Exit(1)
	}
}
