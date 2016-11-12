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
	MaxConnection   int
	FunctionFiles   []string
	TriggerFiles    []string
	ViewFiles       []string
}

/*
 * Retrieve configuration from command line options
 */
func ( config *Config ) Parse() ( err error ) {
	if err = config.parseDatabaseUrl() ; err != nil { return err }
	if err = config.parseFlags() ; err != nil { return err }
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
 * Parse options
 */
func ( config *Config ) parseFlags() ( err error ) {
	flag.BoolVar( &config.WatchMode, "w", false, "Keep watching for filesystem change" )
	flag.IntVar( &config.MaxConnection, "n", 0, "Maximum number of connection (default: pg max_conn / 2)" )
	flag.Parse()

	if config.MaxConnection == 0 {
		config.MaxConnection, err = FindMaxConnection()
		if err != nil { return err }
		if config.MaxConnection == 0 { return fmt.Errorf( "Can't find max connection for postgres" ) }
	}

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
	if len( os.Args ) == 1 {
		return fmt.Errorf( "You must provide a sql directory" )
	}

	config.SqlDirPath = os.Args[ len( os.Args ) - 1 ]

	return
}
