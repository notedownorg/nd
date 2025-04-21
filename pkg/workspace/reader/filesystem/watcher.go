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
	"os"

	ev "github.com/notedownorg/nd/pkg/fsnotify"
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

func (r *Reader) fileWatcher() {
	defer r.watcher.Close()
	for {
		clock := r.clock.Add(1)
		select {
		case event := <-r.watcher.Events():
			switch event.Op {
			case ev.Change:
				r.handleChangeEvent(event, clock)
			case ev.Remove:
				r.handleRemoveEvent(event, clock)
			}
		case err := <-r.watcher.Errors():
			r.log.Warn("received error from filewatcher", "error", err)
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

func (r *Reader) handleChangeEvent(event ev.Event, clock uint64) {
	if isDir(event.Path) {
		r.log.Debug("ignoring directory change event", "dir", event.Path)
		return
	}
	r.log.Debug("handling file change event", "file", event.Path)
	r.processFile(event.Path, false, clock)
}

func (r *Reader) handleRemoveEvent(ev ev.Event, clock uint64) {
	if isDir(ev.Path) {
		r.log.Debug("ignoring directory remove event", "dir", ev.Path)
		return
	}
	r.log.Debug("handling file remove event", "file", ev.Path)
	rel, err := r.relative(ev.Path)
	if err != nil {
		r.log.Error("failed to get relative path", "file", ev.Path, "error", err)
		r.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return
	}
	r.docMutex.Lock()
	defer r.docMutex.Unlock()
	delete(r.documents, rel)
	r.events <- event{Event: reader.Event{Op: reader.Delete, Id: rel}, clock: clock}
}
