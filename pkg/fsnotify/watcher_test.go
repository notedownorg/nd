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

package fsnotify_test

import (
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/fsnotify"
	"github.com/notedownorg/nd/pkg/test/words"
	"github.com/stretchr/testify/assert"
)

func randomFile(root string) string {
	var b strings.Builder
	b.WriteString(root)
	b.WriteString("/")

	// Increase chances of top-level directories being unique as there are issues with case sensitive conflicts on some filesystems
	for range rand.Intn(3) {
		b.WriteString(words.Random())
		b.WriteString("-")
	}
	b.WriteString(words.Random())
	b.WriteString("/")

	// Random depth
	for range rand.Intn(4) {
		b.WriteString(words.Random())
		b.WriteString("/")
	}

	// Random file name
	b.WriteString(words.Random())
	b.WriteString(".file")
	return b.String()
}

func TestRecursiveWatcher(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo.Level()})))

	dir, _ := os.MkdirTemp("", "testrecursivewatcher")
	w, _ := fsnotify.NewRecursiveWatcher(dir, fsnotify.WithDebounce(time.Millisecond*100))

	// What we want to test is that we get an accurate view of the filesystem based on the events we receive
	// This because events are non-deteministic even if you dont take ordering into account
	got := &fileview{files: make(map[string]string)}
	go tracker(t, w, got)

	// Do a bunch of things
	want := &fileview{files: make(map[string]string)}
	for range 1000 {
		path := randomFile(dir)
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			t.Fatal(err)
		}

		// Create
		content := words.Random()
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		want.add(path)
	}
	for path := range want.files {
		// Randomly update, rename or remove
		switch rand.Intn(3) {
		case 0: // Update
			slog.Debug("updating file", "path", path)
			content := words.Random()
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			want.add(path)
		case 1: // Rename
			newpath := randomFile(dir)
			slog.Debug("renaming file", "path", path, "newpath", newpath)
			if err := os.MkdirAll(filepath.Dir(newpath), 0777); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(path, newpath); err != nil {
				t.Fatal(err)
			}
			want.add(newpath)
			want.delete(path)
		case 2: // Remove
			slog.Debug("removing file", "path", path)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			want.delete(path)
		}
	}

	// Wait for the tracker to catch up
	// Theres not a good way to do this deterministically
	time.Sleep(time.Second * 5)
	defer want.mutex.RUnlock()
	defer got.mutex.RUnlock()
	want.mutex.RLock()
	got.mutex.RLock()
	assert.Equal(t, want.files, got.files)
}

// Map of file paths to their content
type fileview struct {
	mutex sync.RWMutex
	files map[string]string
}

func (f *fileview) add(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	f.mutex.Lock()
	(*f).files[path] = string(data)
	f.mutex.Unlock()
}

func (f *fileview) exists(path string) bool {
	f.mutex.RLock()
	_, ok := (*f).files[path]
	f.mutex.RUnlock()
	return ok
}

func (f *fileview) delete(path string) {
	f.mutex.Lock()
	delete((*f).files, path)
	f.mutex.Unlock()
}

func tracker(t *testing.T, w *fsnotify.RecursiveWatcher, view *fileview) {
	for {
		select {
		case event := <-w.Events():
			switch event.Op {
			case fsnotify.Change:
				view.add(event.Path)
			case fsnotify.Remove:
				if !view.exists(event.Path) {
					t.Logf("file %s does not exist during remove", event.Path)
				}
				view.delete(event.Path)
			}
		case err := <-w.Errors():
			t.Log(err)
		}
	}
}
