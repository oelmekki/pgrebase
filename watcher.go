package main

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"os"
	"fmt"
)

type Watcher struct {
	Done    chan bool
	Error   chan error
	notify  *fsnotify.Watcher
	watches []string
}

/*
 * Start the watch loop
 */
func ( watcher *Watcher ) Start() {
	var err error
	watcher.notify, err = fsnotify.NewWatcher()
	if err != nil { watcher.Error <- err ; return }
	defer watcher.notify.Close()

	watcher.build()
	watcher.loop()

	return
}

/*
 * Find all directories and watch them
 */
func ( watcher *Watcher ) build() ( err error ) {
	if err = watcher.notify.Add( Cfg.SqlDirPath ) ; err != nil { return err }

	err = filepath.Walk( Cfg.SqlDirPath, func( path string, info os.FileInfo, err error ) error {
		if IsDir( path ) {
			if err = watcher.notify.Add( path ) ; err != nil { return err }
			watcher.watches = append( watcher.watches, path )
		}

		return nil
	})

	return
}

/*
 * Watcher event loop
 */
func ( watcher *Watcher ) loop() {
	for {
		select {
			case event := <-watcher.notify.Events:
				if ! IsHiddenFile( event.Name ) {
					fmt.Printf( "\nFS changed. Building.\n" )
					watcher.Done <- true
					return
				}

			case err := <-watcher.notify.Errors:
				watcher.Error <- err
				return
		}
	}
}
