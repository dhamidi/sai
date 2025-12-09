package codebase

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileWatcher struct {
	codebase    *Codebase
	stopCh      chan struct{}
	pollInteval time.Duration
	modTimes    map[string]time.Time
}

func NewFileWatcher(c *Codebase) *FileWatcher {
	return &FileWatcher{
		codebase:    c,
		stopCh:      make(chan struct{}),
		pollInteval: 1 * time.Second,
		modTimes:    make(map[string]time.Time),
	}
}

func (w *FileWatcher) Start() {
	go w.run()
}

func (w *FileWatcher) Stop() {
	close(w.stopCh)
}

func (w *FileWatcher) run() {
	ticker := time.NewTicker(w.pollInteval)
	defer ticker.Stop()

	w.scan()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.scan()
		}
	}
}

func (w *FileWatcher) scan() {
	currentFiles := make(map[string]bool)

	filepath.Walk(w.codebase.RootDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".java" {
			return nil
		}

		currentFiles[path] = true

		lastMod, known := w.modTimes[path]
		if !known || info.ModTime().After(lastMod) {
			w.modTimes[path] = info.ModTime()
			w.codebase.ScanFile(path)
		}
		return nil
	})

	for path := range w.modTimes {
		if !currentFiles[path] {
			delete(w.modTimes, path)
			w.codebase.RemoveFile(path)
		}
	}
}
