package core

import (
	"fmt"
	"log"
	"time"
)

// Watch listens to FS change in sql dir and processes them.
func Watch() {
	fmt.Println("Watching filesystem for changes...")

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
