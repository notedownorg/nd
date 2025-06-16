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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/nd/pkg/fsnotify"
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

type Reader struct {
	root    string
	watcher *fsnotify.RecursiveWatcher

	subscribersMutex *sync.RWMutex
	subscribers      map[int]*subscriber
	nextID           int

	errors chan error
}

type subscriber struct {
	id       int
	events   chan reader.Event
	buffered []reader.Event
	mu       sync.Mutex
}

func NewReader(root string) (*Reader, error) {
	watcher, err := fsnotify.NewRecursiveWatcher(
		root,
		fsnotify.WithDebounce(50*time.Millisecond),
		fsnotify.WithMaxWait(150*time.Millisecond),
		fsnotify.WithIgnoredDirs([]string{".git", "node_modules", ".DS_Store"}),
	)
	if err != nil {
		return nil, err
	}

	r := &Reader{
		root:             root,
		watcher:          watcher,
		subscribersMutex: &sync.RWMutex{},
		subscribers:      make(map[int]*subscriber),
		nextID:           1,
		errors:           make(chan error, 100),
	}

	go r.handleEvents()
	go r.forwardErrors()

	return r, nil
}

func (r *Reader) Subscribe(ch chan reader.Event, loadInitialDocuments bool) int {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()

	id := r.nextID
	r.nextID++

	sub := &subscriber{
		id:     id,
		events: ch,
	}

	r.subscribers[id] = sub

	if loadInitialDocuments {
		go r.loadInitialDocuments(sub)
	}

	return id
}

func (r *Reader) Unsubscribe(id int) {
	r.subscribersMutex.Lock()
	defer r.subscribersMutex.Unlock()

	delete(r.subscribers, id)
}

func (r *Reader) Errors() <-chan error {
	return r.errors
}

func (r *Reader) Close() error {
	return r.watcher.Close()
}

func (r *Reader) handleEvents() {
	for event := range r.watcher.Events() {
		if !r.isMarkdownFile(event.Path) {
			continue
		}

		var readerEvent reader.Event
		switch event.Op {
		case fsnotify.Change:
			content, err := os.ReadFile(event.Path)
			if err != nil {
				r.errors <- err
				continue
			}
			readerEvent = reader.Event{
				Op:      reader.Change,
				Id:      r.pathToID(event.Path),
				Content: content,
			}
		case fsnotify.Remove:
			readerEvent = reader.Event{
				Op: reader.Delete,
				Id: r.pathToID(event.Path),
			}
		}

		r.broadcastEvent(readerEvent)
	}
}

func (r *Reader) forwardErrors() {
	for err := range r.watcher.Errors() {
		r.errors <- err
	}
}

func (r *Reader) loadInitialDocuments(sub *subscriber) {
	err := filepath.WalkDir(r.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !r.isMarkdownFile(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			r.errors <- err
			return nil
		}

		event := reader.Event{
			Op:      reader.Load,
			Id:      r.pathToID(path),
			Content: content,
		}

		sub.mu.Lock()
		select {
		case sub.events <- event:
		default:
			sub.buffered = append(sub.buffered, event)
		}
		sub.mu.Unlock()

		return nil
	})

	if err != nil {
		r.errors <- err
	}

	loadCompleteEvent := reader.Event{
		Op: reader.SubscriberLoadComplete,
	}

	sub.mu.Lock()
	select {
	case sub.events <- loadCompleteEvent:
	default:
		sub.buffered = append(sub.buffered, loadCompleteEvent)
	}
	sub.mu.Unlock()

	go r.flushBufferedEvents(sub)
}

func (r *Reader) flushBufferedEvents(sub *subscriber) {
	for {
		sub.mu.Lock()
		if len(sub.buffered) == 0 {
			sub.mu.Unlock()
			return
		}

		event := sub.buffered[0]
		sub.buffered = sub.buffered[1:]
		sub.mu.Unlock()

		select {
		case sub.events <- event:
		default:
			sub.mu.Lock()
			sub.buffered = append([]reader.Event{event}, sub.buffered...)
			sub.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (r *Reader) broadcastEvent(event reader.Event) {
	r.subscribersMutex.RLock()
	defer r.subscribersMutex.RUnlock()

	for _, sub := range r.subscribers {
		sub.mu.Lock()
		select {
		case sub.events <- event:
		default:
			sub.buffered = append(sub.buffered, event)
			go r.flushBufferedEvents(sub)
		}
		sub.mu.Unlock()
	}
}

func (r *Reader) isMarkdownFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".md")
}

func (r *Reader) pathToID(path string) string {
	relPath, err := filepath.Rel(r.root, path)
	if err != nil {
		return path
	}
	return relPath
}
