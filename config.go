package main

import (
	"fmt"
	"flag"
	"os"
)

type Config struct {
	WatchMode       bool
	DatabaseUrl     string
	SqlDirPath      string
	FunctionFiles   []string
	TriggerFiles    []string
	ViewFiles       []string
}

/*
 * Retrieve configuration from command line options
 */
func ( config *Config ) Parse() ( err error ) {
	if err = config.checkWatchMode() ; err != nil { return err }
	if err = config.parseDatabaseUrl() ; err != nil { return err }
	if err = config.parseSqlDir() ; err != nil { return err }

	return
}

/*
 * Scan sql directory for sql files
 */
func ( config *Config ) ScanFiles() {
	sourceWalker := SourceWalker{ Config: config }
	sourceWalker.Process()
}

/*
 * Check if we should exit or keep watching for fs change
 */
func ( config *Config ) checkWatchMode() ( err error ) {
	flag.BoolVar( &config.WatchMode, "w", false, "Keep watching for filesystem change" )
	flag.Parse()

	if config.WatchMode {
	 return fmt.Errorf( "Watch mode is not supported for now" )
	}

	return
}

/*
 * Retrieve database connection info
 */
func ( config *Config ) parseDatabaseUrl() ( err error ) {
	config.DatabaseUrl = os.Getenv( "DATABASE_URL" )

	if  len( config.DatabaseUrl ) == 0 {
		return fmt.Errorf( "You must provide database connection information as DATABASE_URL" )
	}

	return
}

/*
 * Retrieve sql source directory
 */
func ( config *Config ) parseSqlDir() ( err error ) {
	if len( os.Args ) == 1 || ( config.WatchMode && len( os.Args ) == 2 ) {
		return fmt.Errorf( "You must provide a sql directory" )
	}

	config.SqlDirPath = os.Args[ len( os.Args ) - 1 ]

	return
}
