package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	WatchMode     bool
	DatabaseUrl   string
	SqlDirPath    string
	FunctionFiles []string
	TriggerFiles  []string
	TypeFiles     []string
	ViewFiles     []string
}

/*
 * Retrieve configuration from command line options
 */
func (config *Config) Parse() (err error) {
	if err = config.parseDatabaseUrl(); err != nil {
		return err
	}
	if err = config.parseFlags(); err != nil {
		return err
	}
	if err = config.parseSqlDir(); err != nil {
		return err
	}

	return
}

/*
 * Scan sql directory for sql files
 */
func (config *Config) ScanFiles() {
	config.FunctionFiles = make([]string, 0)
	config.TriggerFiles = make([]string, 0)
	config.TypeFiles = make([]string, 0)
	config.ViewFiles = make([]string, 0)

	sourceWalker := SourceWalker{Config: config}
	sourceWalker.Process()
}

/*
 * Parse options
 */
func (config *Config) parseFlags() (err error) {
	flag.BoolVar(&config.WatchMode, "w", false, "Keep watching for filesystem change")
	flag.Parse()

	return
}

/*
 * Retrieve database connection info
 */
func (config *Config) parseDatabaseUrl() (err error) {
	config.DatabaseUrl = os.Getenv("DATABASE_URL")

	if len(config.DatabaseUrl) == 0 {
		return fmt.Errorf("You must provide database connection information as DATABASE_URL")
	}

	return
}

/*
 * Retrieve sql source directory
 */
func (config *Config) parseSqlDir() (err error) {
	if len(os.Args) == 1 {
		return fmt.Errorf("You must provide a sql directory")
	}

	config.SqlDirPath = os.Args[len(os.Args)-1]
	if string(config.SqlDirPath[len(config.SqlDirPath)-1]) != "/" {
		config.SqlDirPath += "/"
	}

	return
}
