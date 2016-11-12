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
		os.Exit(1)
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
}

/*
 * Give user a chance to know what went wrong
 */
func Usage() {
	usage := `
USAGE:
	DATABASE_URL=url pgrebase [-w] <sql_directory>

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
	`;

	fmt.Println( usage )
}

/*
 * Start the actual work
 */
func Process() {
}

func main() {
	ParseConfig()
	CheckSanity()
	Process()
}
