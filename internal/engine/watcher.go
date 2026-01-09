package engine

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/ismailtsdln/DevTestrider/internal/config"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	config  config.WatchConfig
	Events  chan string
	done    chan bool
	mu      sync.Mutex
}

func NewWatcher(cfg config.WatchConfig) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: w,
		config:  cfg,
		Events:  make(chan string),
		done:    make(chan bool),
	}, nil
}

func (w *Watcher) Start() {
	defer w.watcher.Close()

	if err := w.addPaths(w.config.Paths); err != nil {
		log.Printf("Error adding paths: %v", err)
	}

	var debounceTimer *time.Timer
	const debounceDuration = 500 * time.Millisecond

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Handle file operations
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {

				// Handle new directory creation regardless of name (unless ignored)
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err == nil && fi.IsDir() && !w.shouldIgnore(event.Name) {
						if err := w.watcher.Add(event.Name); err == nil {
							log.Printf("Watching new directory: %s", event.Name)
						}
						// Don't return, as we might want to trigger a test run? Maybe not for just a dir create.
					}
				}

				// Check if it's a Go file or relevant file for triggering tests
				if !w.shouldIgnore(event.Name) && strings.HasSuffix(event.Name, ".go") {
					// Debounce logic
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(debounceDuration, func() {
						w.Events <- event.Name
					})
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		case <-w.done:
			return
		}
	}
}

func (w *Watcher) Stop() {
	w.done <- true
}

func (w *Watcher) addPaths(paths []string) error {
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if w.shouldIgnore(path) {
					return filepath.SkipDir
				}
				return w.watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) shouldIgnore(path string) bool {
	for _, ignore := range w.config.Ignore {
		if strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}
