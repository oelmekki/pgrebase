package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/oelmekki/pgrebase/core/utils"
	"os"
	"path/filepath"
)

// watcher contains data for watching fs for code change.
type watcher struct {
	Done    chan bool         // pinged once watcher found changes
	Error   chan error        // pinged if an error occured
	notify  *fsnotify.Watcher // the fsnotify low level watcher
	watches []string          // the list of file paths being watched
}

// Start starts the watch loop.
func (w *watcher) Start() {
	var err error
	w.notify, err = fsnotify.NewWatcher()
	if err != nil {
		w.Error <- err
		return
	}
	defer w.notify.Close()

	w.build()
	w.loop()

	return
}

// build finds all directories and watch them.
func (w *watcher) build() (err error) {
	if err = w.notify.Add(conf.SqlDirPath); err != nil {
		return err
	}

	err = filepath.Walk(conf.SqlDirPath, func(path string, info os.FileInfo, err error) error {
		if utils.IsDir(path) {
			if err = w.notify.Add(path); err != nil {
				return err
			}
			w.watches = append(w.watches, path)
		}

		return nil
	})

	return
}

// loop starts the watcher event processing loop.
func (w *watcher) loop() {
	for {
		select {
		case event := <-w.notify.Events:
			if !utils.IsHiddenFile(event.Name) {
				fmt.Printf("\nFS changed. Building.\n")
				w.Done <- true
				return
			}

		case err := <-w.notify.Errors:
			w.Error <- err
			return
		}
	}
}
