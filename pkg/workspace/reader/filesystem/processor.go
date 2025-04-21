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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notedownorg/nd/pkg/workspace/reader"
)

func (r *Reader) processFile(path string, load bool, clock uint64) {
	// If we have already processed this file and it is up to date, we can skip it
	if !load && r.isUpToDate(path) {
		r.log.Debug("file is up to date, stopping short", "file", path)
		return
	}
	// Do the rest in a goroutine so we can continue doing other things
	r.log.Debug("processing file", "file", path)
	go func() {
		rel, content := r.loadDocument(path)

		r.log.Debug("updating document in cache", "file", path, "relative", rel)
		r.docMutex.Lock()
		r.documents[rel] = document{id: rel, lastModified: time.Now()}
		r.docMutex.Unlock()

		r.log.Debug("emitting event", "file", path, "relative", rel)
		op := func() reader.Operation {
			if load {
				return reader.Load
			} else {
				return reader.Change
			}
		}()
		r.events <- event{Event: reader.Event{Op: op, Id: rel, Content: content}, clock: clock}
	}()
}

// Takes absolute or relative path and returns relative key + content
func (r *Reader) loadDocument(path string) (string, []byte) {
	rel, err := r.relative(path)
	if err != nil {
		r.log.Error("failed to get relative path", "file", path, "error", err)
		r.errors <- fmt.Errorf("failed to get relative path: %w", err)
	}

	// Acquire semaphore as we make a blocking syscall to read file
	// This ensures we don't exhaust the thread limit when reading files in parallel
	r.threadLimit.Acquire(context.Background(), 1)
	defer r.threadLimit.Release(1)

	content, err := os.ReadFile(r.absolute(path))
	if err != nil {
		r.log.Error("failed to read file", "path", path, "error", err)
		r.errors <- fmt.Errorf("failed to read file: %w", err)
		return rel, nil
	}
	return rel, content
}

func (r *Reader) isUpToDate(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		r.log.Error("failed to get file info", "file", file, "error", err)
		r.errors <- fmt.Errorf("failed to get file info: %w", err)
		return false
	}
	rel, err := r.relative(file)
	if err != nil {
		r.log.Error("failed to get relative path", "file", file, "error", err)
		r.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return false
	}
	r.docMutex.RLock()
	doc, ok := r.documents[rel]
	r.docMutex.RUnlock()
	return ok && doc.Modified(info.ModTime())
}
