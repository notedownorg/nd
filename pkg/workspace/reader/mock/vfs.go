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

package mock

import (
	"math/rand/v2"

	"github.com/notedownorg/nd/pkg/workspace/reader"
)

func (r *Reader) Add(path string, content []byte) {
	r.vfsMutex.Lock()
	r.vfs[path] = content
	r.vfsMutex.Unlock()
	r.sendEvent(reader.Event{Op: reader.Change, Id: path, Content: content})
}

func (r *Reader) Update(path string, content []byte) {
	r.vfsMutex.Lock()
	r.vfs[path] = content
	r.vfsMutex.Unlock()
	r.sendEvent(reader.Event{Op: reader.Change, Id: path, Content: content})
}

func (r *Reader) Rename(oldPath, newPath string) {
	r.vfsMutex.Lock()
	content := r.vfs[oldPath]
	r.vfs[newPath] = content
	delete(r.vfs, oldPath)
	r.vfsMutex.Unlock()
	r.sendEvent(reader.Event{Op: reader.Delete, Id: oldPath})
	r.sendEvent(reader.Event{Op: reader.Change, Id: newPath, Content: content})
}

func (r *Reader) Remove(path string) {
	r.vfsMutex.Lock()
	delete(r.vfs, path)
	r.vfsMutex.Unlock()
	r.sendEvent(reader.Event{Op: reader.Delete, Id: path})
}

func (r *Reader) ListFiles() []string {
	r.vfsMutex.RLock()
	defer r.vfsMutex.RUnlock()

	paths := make([]string, 0, len(r.vfs))
	for path := range r.vfs {
		paths = append(paths, path)
	}
	return paths
}

func (r *Reader) RandomFile() string {
	r.vfsMutex.RLock()
	defer r.vfsMutex.RUnlock()

	files := r.ListFiles()
	if len(files) == 0 {
		return ""
	}
	return files[rand.IntN(len(files))]
}

func (r *Reader) GetContent(path string) []byte {
	r.vfsMutex.RLock()
	defer r.vfsMutex.RUnlock()

	return r.vfs[path]
}
