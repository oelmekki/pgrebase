package main

import (
	"flag"
	"fmt"
	"github.com/oelmekki/pgrebase/core"
	"os"
)

// Usage shows user how they're supposed to use application.
func Usage() {
	usage := `
PgRebase-1.2.0

USAGE:
	DATABASE_URL=url pgrebase [-w] <sql_directory>

PgRebase is a tool that allows you to easily handle your postgres codebase for
functions, triggers, types and views.

Your expected to provide a postgresql connection url as DATABASE_URL and
a sql directory as <sql_directory>.

<sql_directory> should be structured that way:
	<sql_directory>/
	├── functions/
	├── triggers/
	├── types/
	└── views/

At least one of functions/triggers/types/views should exist.

OPTIONS:
	-w: enter watch mode.
		In watch mode, pgrebase will keep watching for file changes and will
		automatically reload your sql code when it happens.
	`

	fmt.Println(usage)
	os.Exit(1)
}

// parseDatabaseUrl retrieves database connection info.
func parseDatabaseUrl() (url string, err error) {
	url = os.Getenv("DATABASE_URL")

	if len(url) == 0 {
		err = fmt.Errorf("You must provide database connection information as DATABASE_URL")
		return
	}

	return
}

// parseFlags parses options.
func parseFlags() (watch bool) {
	flag.BoolVar(&watch, "w", false, "Keep watching for filesystem change")
	flag.Parse()

	return
}

func main() {
	url, err := parseDatabaseUrl()
	if err != nil {
		fmt.Println(err)
		Usage()
	}

	if len(os.Args) == 1 {
		fmt.Println("You must provide a sql directory")
		Usage()
	}

	watch := parseFlags()

	sqlDir := os.Args[len(os.Args)-1]

	err = core.Init(url, sqlDir)
	if err != nil {
		fmt.Printf("Can't initialize pgrebase: %v\n\n", err)
		Usage()
	}

	if err := core.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if watch {
		core.Watch()
	}
}
