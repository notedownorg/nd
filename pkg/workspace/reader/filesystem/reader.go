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

package filesystem

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/notedownorg/nd/pkg/fsnotify"
	"github.com/notedownorg/nd/pkg/workspace/reader"
	"golang.org/x/sync/semaphore"
)

var _ reader.Reader = &Reader{}

// The filesystem reader is responsible for processing files on disk and emitting events when changes
// to the graph are detected (i.e. a file is added, removed or modified)
// It is not responsible for parsing the documents, maintaining the graph or handling mutations to the graph/files.
type Reader struct {
	log   *slog.Logger
	root  string
	clock atomic.Uint64

	// Map of relative location (id) to document
	documents map[string]document
	docMutex  sync.RWMutex

	watcher *fsnotify.RecursiveWatcher

	subscriberMutex sync.RWMutex
	subscribers     map[int]*subscriber

	// Everytime a goroutine makes a blocking syscall (in our case usually file i/o) it uses a new thread so to avoid
	// large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted

	errors chan error
	events chan event
}

func NewReader(name string, location string) (*Reader, error) {
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
		subscribers: make(map[int]*subscriber),
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
		errors:      make(chan error),
		events:      make(chan event),
	}

	// Create a subscription so we can listen for the initial load events
	sub := make(chan reader.Event)
	defer client.Unsubscribe(client.Subscribe(sub, true))

	// Wait for the subscriber to complete its initial load
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for ev := range sub {
			if ev.Op == reader.SubscriberLoadComplete {
				wg.Done()
				return
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
			client.processFile(path, true, 0) // ignore clock for initial load
		}
		return nil
	})

	// Wait for all initial loads to finish, unsubscribe and close the channel
	client.log.Debug("waiting for initial load to complete")
	wg.Wait()
	client.log.Debug("initial load complete")

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
