package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

// conf is the global level configuration data structure.
var conf config

// Init stores the global config object.
//
// databaseUrl should be a connection string to the database (eg: postgres://postgres:@localhost/database).
//
// sqlDir is the path to the directory where sql source files live.
//
// watch should be true if you want to keep watching for changes in source files rather
// than just loading them once.
func Init(databaseUrl, sqlDir string) (err error) {
	conf = newConfig(databaseUrl, sqlDir)

	checker := sanity{}
	err = checker.check()

	return
}

// Process loads sql code, just once.
func Process() (err error) {
	if err = loadViews(); err != nil {
		return err
	}
	if err = loadFunctions(); err != nil {
		return err
	}
	if err = loadTriggers(); err != nil {
		return err
	}

	return
}

// startWatching fires a watcher, will die as soon something changed.
func startWatching(errorChan chan error, doneChan chan bool) (err error) {
	w := watcher{Done: doneChan, Error: errorChan}
	go w.Start()

	return
}

// Watch listens to FS change in sql dir and processes them.
func Watch() {
	if os.Getenv("QUIET") != "true" {
		fmt.Println("Watching filesystem for changes...")
	}

	errorChan := make(chan error)
	doneChan := make(chan bool)

	if err := startWatching(errorChan, doneChan); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-doneChan:
			time.Sleep(300 * time.Millisecond) // without this, new file watcher is started faster than file writing has ended
			scanFiles(&conf)

			if err := Process(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			if err := startWatching(errorChan, doneChan); err != nil {
				log.Fatal(err)
			}

		case err := <-errorChan:
			fmt.Printf("Error: %v\n", err)
			if err := startWatching(errorChan, doneChan); err != nil {
				log.Fatal(err)
			}
		}
	}
}
