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

package reader

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/nd/pkg/fsnotify"
	"golang.org/x/sync/semaphore"
)

// The reader is responsible for processing files on disk and emitting events when changes
// to the graph are detected (i.e. a file is added, removed or modified)
// It is not responsible for parsing the documents, maintaining the graph or handling mutations to the graph/files.
type Reader struct {
	log  *slog.Logger
	root string

	// Map of relative location (id) to document
	documents map[string]document
	docMutex  sync.RWMutex

	watcher *fsnotify.RecursiveWatcher

	subscribers []chan Event

	// Everytime a goroutine makes a blocking syscall (in our case usually file i/o) it uses a new thread so to avoid
	// large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted

	errors chan error
	events chan Event
}

func NewClient(name string, location string) (*Reader, error) {
	if !filepath.IsAbs(location) {
		return nil, fmt.Errorf("location must be an absolute path got %s", location)
	}
	ignoredDirs := []string{".git", ".vscode", ".debug", ".stversions", ".stfolder"}
	watcher, err := fsnotify.NewRecursiveWatcher(location, fsnotify.WithIgnoredDirs(ignoredDirs))
	if err != nil {
		return nil, err
	}

	client := &Reader{
		log:         slog.Default().With("name", name),
		root:        location,
		documents:   make(map[string]document),
		docMutex:    sync.RWMutex{},
		watcher:     watcher,
		subscribers: make([]chan Event, 0),
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
		errors:      make(chan error),
		events:      make(chan Event),
	}

	// Create a subscription so we can listen for the initial load events
	sub := make(chan Event)
	subscriberIndex := client.Subscribe(sub)

	// For each file we process on intial load, a load event is emitted
	// Therefore if our subscriber has received a load event for each file we have finished the initial load
	var wg sync.WaitGroup
	go func() {
		for ev := range sub {
			if ev.Op == Load {
				wg.Done()
			}
		}
	}()

	go client.fileWatcher()
	go client.eventDispatcher()

	// Recurse through the root directory and process all the files to build the initial state
	client.log.Debug("walking workspace to build initial state")
	err = filepath.Walk(client.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ignoredDir := range ignoredDirs {
			if strings.Contains(path, ignoredDir) {
				return nil
			}
		}
		if strings.HasSuffix(path, ".md") {
			wg.Add(1) // Increment the wait group for each file we process
			client.processFile(path, true)
		}
		return nil
	})

	// Wait for all initial loads to finish, unsubscribe and close the channel
	client.log.Debug("waiting for initial load to complete")
	wg.Wait()
	client.Unsubscribe(subscriberIndex)
	close(sub)

	return client, nil
}

func (r *Reader) absolute(relative string) string {
	if filepath.IsAbs(relative) {
		return relative
	}
	return filepath.Join(r.root, relative)
}

func (r *Reader) relative(absolute string) (string, error) {
	if !filepath.IsAbs(absolute) {
		return absolute, nil
	}
	return filepath.Rel(r.root, absolute)
}

// We dont need to store the actual document content in the reader, just the minimal information required
// to determine if the document has been modified and how to inform subscribers (i.e. the last modified time and ID)
type document struct {
	id           string
	lastModified time.Time
}

func (d document) Modified(lastModified time.Time) bool {
	return lastModified.UnixNano() > d.lastModified.UnixNano()
}
