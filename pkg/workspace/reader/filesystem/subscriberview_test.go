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

// Package filesystem contains test utilities for invariant testing.
//
// subscriberView tracks a subscriber's view of the filesystem based on events received.
// This is the core of invariant testing - we build up what each subscriber thinks
// the filesystem looks like, then compare it to reality.

package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/notedownorg/nd/pkg/workspace/reader"
)

// subscriberView tracks a subscriber's view of the filesystem based on events received.
// This is the core of invariant testing - we build up what each subscriber thinks
// the filesystem looks like, then compare it to reality.
type subscriberView struct {
	files map[string][]byte // Maps file ID to content
	mu    sync.RWMutex      // Protects concurrent access during event processing
}

func newSubscriberView() *subscriberView {
	return &subscriberView{
		files: make(map[string][]byte),
	}
}

// applyEvent updates the subscriber's view based on a received event.
// This simulates how a real client would maintain state from events.
func (sv *subscriberView) applyEvent(event reader.Event) {
	sv.mu.Lock()
	defer sv.mu.Unlock()

	switch event.Op {
	case reader.Load, reader.Change:
		// File created or updated - store the content
		sv.files[event.Id] = make([]byte, len(event.Content))
		copy(sv.files[event.Id], event.Content)
	case reader.Delete:
		// File deleted - remove from view
		delete(sv.files, event.Id)
	case reader.SubscriberLoadComplete:
		// Marker event - no state change needed
	}
}

// getFiles returns a thread-safe copy of all files in this subscriber's view.
func (sv *subscriberView) getFiles() map[string][]byte {
	sv.mu.RLock()
	defer sv.mu.RUnlock()

	result := make(map[string][]byte)
	for id, content := range sv.files {
		result[id] = make([]byte, len(content))
		copy(result[id], content)
	}
	return result
}

// getFileContent returns the content of a specific file, if it exists.
func (sv *subscriberView) getFileContent(id string) ([]byte, bool) {
	sv.mu.RLock()
	defer sv.mu.RUnlock()
	content, exists := sv.files[id]
	if !exists {
		return nil, false
	}
	result := make([]byte, len(content))
	copy(result, content)
	return result, true
}

// getActualFilesystemState reads the current state of markdown files from disk.
// This is our "ground truth" that subscriber views should match.
func getActualFilesystemState(rootDir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Convert to relative path (same format as event IDs)
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		files[relPath] = content
		return nil
	})

	return files, err
}
