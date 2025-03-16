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
	"log"
	"os"

	"github.com/notedownorg/nd/pkg/fsnotify"
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

func (r *Reader) fileWatcher() {
	defer r.watcher.Close()
	for {
		select {
		case event := <-r.watcher.Events():
			switch event.Op {
			case fsnotify.Create:
				r.handleCreateEvent(event)
			case fsnotify.Remove:
				r.handleRemoveEvent(event)
			case fsnotify.Rename:
				r.handleRenameEvent(event)
			case fsnotify.Write:
				r.handleWriteEvent(event)
			}
		case err := <-r.watcher.Errors():
			log.Printf("error: %s", err)
		}
	}
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (r *Reader) handleCreateEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		r.log.Debug("ignoring directory create event", "dir", event.Name)
		return
	}
	r.log.Debug("handling file create event", "file", event.Name)
	r.processFile(event.Name, false)
}

func (r *Reader) handleRemoveEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		r.log.Debug("ignoring directory remove event", "dir", event.Name)
		return
	}
	r.log.Debug("handling file remove event", "file", event.Name)
	rel, err := r.relative(event.Name)
	if err != nil {
		r.log.Error("failed to get relative path", "file", event.Name, "error", err)
		r.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return
	}
	r.docMutex.Lock()
	defer r.docMutex.Unlock()
	delete(r.documents, rel)
	r.events <- reader.Event{Op: reader.Delete, Id: rel}
}

func (r *Reader) handleRenameEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		r.log.Debug("ignoring directory rename event", "dir", event.Name)
		return
	}
	r.log.Debug("handling file rename event", "file", event.Name)
	r.handleRemoveEvent(event) // rename sends the name of the old file, presumably it sends a create event for the new file
}

func (r *Reader) handleWriteEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		r.log.Debug("ignoring directory write event", "dir", event.Name)
		return
	}
	r.log.Debug("handling file write event", "file", event.Name)
	r.processFile(event.Name, false)
}
