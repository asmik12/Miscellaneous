package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	logFile, err := os.OpenFile("downloads_activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "FSLOG ", log.LstdFlags)

	// Watch ~/Downloads
	home, _ := os.UserHomeDir()
	print(home, "\n")
	watchPath := filepath.Join(home, "/Desktop/TestWatcj")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(watchPath)
	if err != nil {
		log.Fatal(err)
	}

	logger.Printf("Watching folder: %s\n", watchPath)
	fmt.Printf("Monitoring %s ...\n", watchPath)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			logger.Printf("EVENT: %s %s\n", event.Op, event.Name)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Printf("ERROR: %s\n", err)
		}
	}
}
