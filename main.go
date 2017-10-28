package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Cfg is the global level configuration data structure.
var Cfg Config

// ParseConfig makes sure user provided needed configuration.
func ParseConfig() {
	if err := Cfg.Parse(); err != nil {
		fmt.Printf("Error: %v\n", err)
		Usage()
	}
}

// CheckSanity makes sure all requirements are satisfied.
func CheckSanity() {
	sanity := Sanity{}
	if err := sanity.Check(); err != nil {
		fmt.Printf("Error: %v\n", err)
		Usage()
	}
}

// Usage shows user how they're supposed to use application.
func Usage() {
	usage := `
PgRebase-1.1.2

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

// Process starts the actual work.
func Process() (err error) {
	if err = LoadTypes(); err != nil {
		return err
	}
	if err = LoadViews(); err != nil {
		return err
	}
	if err = LoadFunctions(); err != nil {
		return err
	}
	if err = LoadTriggers(); err != nil {
		return err
	}

	return
}

// StartWatching fires a watcher, will die as soon something changed.
func StartWatching(errorChan chan error, doneChan chan bool) (err error) {
	watcher := Watcher{Done: doneChan, Error: errorChan}
	go watcher.Start()

	return
}

// WatchTheWatcher processes events from watchers.
func WatchTheWatcher() {
	fmt.Println("Watching filesystem for changes...")

	errorChan := make(chan error)
	doneChan := make(chan bool)

	if err := StartWatching(errorChan, doneChan); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-doneChan:
			time.Sleep(300 * time.Millisecond) // without this, new file watcher is started faster than file writing has ended
			Cfg.ScanFiles()

			if err := Process(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			if err := StartWatching(errorChan, doneChan); err != nil {
				log.Fatal(err)
			}

		case err := <-errorChan:
			fmt.Printf("Error: %v\n", err)
			if err := StartWatching(errorChan, doneChan); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	ParseConfig()
	CheckSanity()

	if err := Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if Cfg.WatchMode {
		WatchTheWatcher()
	}
}
