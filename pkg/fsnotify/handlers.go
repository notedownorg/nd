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
	"errors"
	"os"

	"github.com/fsnotify/fsnotify"
)

func (rw *RecursiveWatcher) handleCreate(event fsnotify.Event) {
	// If the event is a directory, add it to the watcher
	if isDir(event.Name) {
		rw.add(event.Name)
	}

	// Check if the file hasn't already been removed before sending the event.
	// As we have no order guarantee from fsnotify, there's a non-zero chance a client
	// may get this event after the file has been removed.
	if _, err := os.Stat(event.Name); errors.Is(err, os.ErrNotExist) {
		return
	}

	// Send the event to the events channel
	rw.events <- event
}

func (rw *RecursiveWatcher) handleRemove(event fsnotify.Event) {
	// If the event is a directory, remove it from the watcher
	if isDir(event.Name) {
		if err := rw.w.Remove(event.Name); err != nil {
			rw.errors <- err
		}
		delete(rw.watchers, event.Name)
	}

	// Send the event to the events channel
	rw.events <- event
}
