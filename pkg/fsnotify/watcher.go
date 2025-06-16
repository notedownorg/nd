// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsnotify

import (
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	root        string
	ignoredDirs []string

	w      *fsnotify.Watcher
	events chan Event
	errors chan error

	watchersMutex *sync.RWMutex
	watchers      map[string]struct{}

	timersMutex   *sync.RWMutex
	timers        map[string]*time.Timer // used to debounce/dedupe events
	maxWaitTimers map[string]*time.Timer // ensures events are sent within max wait time
	debounce      time.Duration
	maxWait       time.Duration
}

type Option func(*RecursiveWatcher)

func WithIgnoredDirs(dirs []string) Option {
	return func(rw *RecursiveWatcher) {
		rw.ignoredDirs = dirs
	}
}

func WithDebounce(d time.Duration) Option {
	return func(rw *RecursiveWatcher) {
		rw.debounce = d
	}
}

func WithMaxWait(d time.Duration) Option {
	return func(rw *RecursiveWatcher) {
		rw.maxWait = d
	}
}

func NewRecursiveWatcher(root string, opts ...Option) (*RecursiveWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rw := &RecursiveWatcher{
		root:        root,
		ignoredDirs: []string{},
		w:           w,
		events:      make(chan Event),
		errors:      make(chan error),

		watchersMutex: &sync.RWMutex{},
		watchers:      make(map[string]struct{}),

		timersMutex:   &sync.RWMutex{},
		timers:        make(map[string]*time.Timer),
		maxWaitTimers: make(map[string]*time.Timer),
		maxWait:       500 * time.Millisecond, // default max wait
	}

	for _, opt := range opts {
		opt(rw)
	}

	go rw.eventLoop()
	go rw.add(root) // run in a goroutine so we can return before events are being read

	return rw, nil
}

func (rw RecursiveWatcher) Events() <-chan Event {
	return rw.events
}

func (rw RecursiveWatcher) Errors() <-chan error {
	return rw.errors
}

func (rw *RecursiveWatcher) Close() error {
	return rw.w.Close()
}

func (rw *RecursiveWatcher) eventLoop() {
	for {
		select {
		case event, ok := <-rw.w.Events:
			if !ok {
				return
			}

			// Debounce events to dedupe writes
			// It has the nice side effect of reducing the likelihood of receiving a write event after a remove event

			// If there is a create event send it immediately
			if event.Op.Has(fsnotify.Create) {
				slog.Debug("handling create event", "path", event.Name)
				rw.handleCreate(event)
				continue
			}

			// If there is a write event, debounce it.
			if event.Op.Has(fsnotify.Write) {
				rw.timersMutex.Lock()

				// Cancel existing debounce timer
				if timer, ok := rw.timers[event.Name]; ok {
					timer.Stop()
				}

				// Start/restart debounce timer
				rw.timers[event.Name] = time.AfterFunc(rw.debounce, func() {
					rw.sendWriteEvent(event.Name)
				})

				// Start max wait timer if it doesn't exist
				if _, ok := rw.maxWaitTimers[event.Name]; !ok {
					rw.maxWaitTimers[event.Name] = time.AfterFunc(rw.maxWait, func() {
						rw.sendWriteEvent(event.Name)
					})
				}

				rw.timersMutex.Unlock()
			}

			// If the event is a remove or rename event, tidy up the timer and handle the event
			// Rename only tells us information about the old path, so we treat it like a remove event
			if event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) {
				rw.timersMutex.Lock()
				timer, ok := rw.timers[event.Name]
				if ok {
					timer.Stop()
					delete(rw.timers, event.Name)
				}
				maxWaitTimer, ok := rw.maxWaitTimers[event.Name]
				if ok {
					maxWaitTimer.Stop()
					delete(rw.maxWaitTimers, event.Name)
				}
				rw.timersMutex.Unlock()
				slog.Debug("handling remove event", "path", event.Name)
				rw.handleRemove(event)
			}
		case err, ok := <-rw.w.Errors:
			if !ok {
				return
			}
			rw.errors <- err
		}
	}
}

func (rw *RecursiveWatcher) add(path string) {
	slog.Debug("processing path", "path", path)

	// Check if the path is already being watched
	rw.watchersMutex.RLock()
	_, ok := rw.watchers[path]
	rw.watchersMutex.RUnlock()
	if ok {
		slog.Debug("path already being watched", "path", path)
		return
	}

	// Check the path does not match any ignored directories
	for _, ignoredDir := range rw.ignoredDirs {
		if strings.Contains(path, ignoredDir) {
			slog.Debug("path matches ignored directory", "path", path, "ignoredDir", ignoredDir)
			return
		}
	}

	// If the path is not a directory, return
	// See fsnotify docs as to why we dont bother watching individual files
	if !isDir(path) {
		slog.Debug("skipping path, not a directory", "path", path)
		return
	}

	// Add the path to the watcher
	if err := rw.w.Add(path); err != nil {
		rw.errors <- err
		return
	}

	// Add the path to the watchers map
	rw.watchersMutex.Lock()
	rw.watchers[path] = struct{}{}
	rw.watchersMutex.Unlock()
	slog.Debug("watching path", "path", path)

	// Iterate over all the entries in the directory (subdirectories) and recurse
	entries, err := os.ReadDir(path)
	if err != nil {
		rw.errors <- err
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			rw.add(path + "/" + entry.Name())
		}
	}

	// This is a bit of a hack to ensure we get events when the following race condition occurs:
	// 1. A new directory is created
	// 2. A file is created in the new directory
	// 3. The directory is added to the watcher
	// 4. No event is sent for the file == sad face

	// Reload the entries
	entries, err = os.ReadDir(path)
	if err != nil {
		rw.errors <- err
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			slog.Debug("sending create event for file in new directory", "path", path+"/"+entry.Name())
			rw.events <- Event{Path: path + "/" + entry.Name(), Op: Change}
		}
	}
}

func (rw *RecursiveWatcher) sendWriteEvent(path string) {
	rw.timersMutex.Lock()
	defer rw.timersMutex.Unlock()

	// Clean up both timers
	if timer, ok := rw.timers[path]; ok {
		timer.Stop()
		delete(rw.timers, path)
	}
	if maxWaitTimer, ok := rw.maxWaitTimers[path]; ok {
		maxWaitTimer.Stop()
		delete(rw.maxWaitTimers, path)
	}

	slog.Debug("handling write event", "path", path)
	rw.events <- Event{Path: path, Op: Change}
}

func isDir(path string) bool {
	fi, error := os.Stat(path)
	if error != nil {
		return false
	}
	return fi.IsDir()
}
