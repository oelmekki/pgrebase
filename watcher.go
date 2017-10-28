package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

// Watcher contains data for watching fs for code change.
type Watcher struct {
	Done    chan bool         // pinged once watcher found changes
	Error   chan error        // pinged if an error occured
	notify  *fsnotify.Watcher // the fsnotify low level watcher
	watches []string          // the list of file paths being watched
}

// Start starts the watch loop.
func (watcher *Watcher) Start() {
	var err error
	watcher.notify, err = fsnotify.NewWatcher()
	if err != nil {
		watcher.Error <- err
		return
	}
	defer watcher.notify.Close()

	watcher.build()
	watcher.loop()

	return
}

// build finds all directories and watch them.
func (watcher *Watcher) build() (err error) {
	if err = watcher.notify.Add(Cfg.SqlDirPath); err != nil {
		return err
	}

	err = filepath.Walk(Cfg.SqlDirPath, func(path string, info os.FileInfo, err error) error {
		if IsDir(path) {
			if err = watcher.notify.Add(path); err != nil {
				return err
			}
			watcher.watches = append(watcher.watches, path)
		}

		return nil
	})

	return
}

// loop starts the watcher event processing loop.
func (watcher *Watcher) loop() {
	for {
		select {
		case event := <-watcher.notify.Events:
			if !IsHiddenFile(event.Name) {
				fmt.Printf("\nFS changed. Building.\n")
				watcher.Done <- true
				return
			}

		case err := <-watcher.notify.Errors:
			watcher.Error <- err
			return
		}
	}
}
