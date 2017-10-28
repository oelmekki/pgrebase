package core

// startWatching fires a watcher, will die as soon something changed.
func startWatching(errorChan chan error, doneChan chan bool) (err error) {
	w := watcher{Done: doneChan, Error: errorChan}
	go w.Start()

	return
}
