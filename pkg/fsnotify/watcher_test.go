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
	"fmt"
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

func TestMaxWaitDebouncing(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsnotify-maxwait-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	w, err := fsnotify.NewRecursiveWatcher(
		tmpDir,
		fsnotify.WithDebounce(100*time.Millisecond),
		fsnotify.WithMaxWait(200*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	testFile := filepath.Join(tmpDir, "test.file")

	events := make([]fsnotify.Event, 0)
	eventTimes := make([]time.Time, 0)

	go func() {
		for {
			select {
			case event := <-w.Events():
				events = append(events, event)
				eventTimes = append(eventTimes, time.Now())
			case err := <-w.Errors():
				t.Log("Error:", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	startTime := time.Now()

	for i := 0; i < 10; i++ {
		content := []byte(strings.Repeat("x", i+1))
		err := os.WriteFile(testFile, content, 0644)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(500 * time.Millisecond)

	if len(events) == 0 {
		t.Fatal("Expected at least one event")
	}

	changeEvents := 0
	for _, event := range events {
		if event.Op == fsnotify.Change && event.Path == testFile {
			changeEvents++
		}
	}

	assert.Greater(t, changeEvents, 0, "Should receive at least one change event")

	if len(eventTimes) > 0 {
		firstEventTime := eventTimes[0]
		timeSinceStart := firstEventTime.Sub(startTime)

		assert.True(t, timeSinceStart <= 250*time.Millisecond,
			"First event should arrive within max wait time (got %v)", timeSinceStart)
	}

	t.Logf("Received %d change events from %d rapid writes", changeEvents, 10)
}

func TestMaxWaitEnsuresEventDelivery(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsnotify-delivery-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	w, err := fsnotify.NewRecursiveWatcher(
		tmpDir,
		fsnotify.WithDebounce(200*time.Millisecond),
		fsnotify.WithMaxWait(300*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	testFile := filepath.Join(tmpDir, "delivery.file")

	eventReceived := make(chan time.Time, 1)

	go func() {
		for {
			select {
			case event := <-w.Events():
				if event.Op == fsnotify.Change && event.Path == testFile {
					eventReceived <- time.Now()
					return
				}
			case err := <-w.Errors():
				t.Log("Error:", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	startTime := time.Now()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	updateCount := 0
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				updateCount++
				content := []byte(fmt.Sprintf("content-%d", updateCount))
				os.WriteFile(testFile, content, 0644)

				if updateCount >= 5 {
					done <- struct{}{}
					return
				}
			case <-done:
				return
			}
		}
	}()

	select {
	case eventTime := <-eventReceived:
		timeSinceStart := eventTime.Sub(startTime)
		assert.True(t, timeSinceStart <= 400*time.Millisecond,
			"Event should be delivered within max wait time despite rapid updates (got %v)", timeSinceStart)
		t.Logf("Event delivered after %v with %d rapid updates", timeSinceStart, updateCount)
	case <-time.After(1 * time.Second):
		t.Fatal("Event was not delivered within reasonable time")
	}
}

func TestDebounceWithoutMaxWait(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsnotify-debounce-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	w, err := fsnotify.NewRecursiveWatcher(
		tmpDir,
		fsnotify.WithDebounce(200*time.Millisecond),
		fsnotify.WithMaxWait(1*time.Second), // Long max wait to ensure debounce works
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	testFile := filepath.Join(tmpDir, "debounce.file")

	events := make([]fsnotify.Event, 0)

	go func() {
		for {
			select {
			case event := <-w.Events():
				if event.Op == fsnotify.Change && event.Path == testFile {
					events = append(events, event)
				}
			case err := <-w.Errors():
				t.Log("Error:", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Create file first and wait for initial event
	os.WriteFile(testFile, []byte("initial"), 0644)
	time.Sleep(300 * time.Millisecond)

	// Clear events from file creation
	events = make([]fsnotify.Event, 0)

	// Now do rapid updates
	os.WriteFile(testFile, []byte("first"), 0644)
	time.Sleep(50 * time.Millisecond)
	os.WriteFile(testFile, []byte("second"), 0644)
	time.Sleep(50 * time.Millisecond)
	os.WriteFile(testFile, []byte("third"), 0644)

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, 1, len(events), "Should receive exactly one debounced event")
	t.Logf("Received %d events from 3 rapid writes (debounce working)", len(events))
}

func TestTimerCleanupOnRemove(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsnotify-cleanup-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	w, err := fsnotify.NewRecursiveWatcher(
		tmpDir,
		fsnotify.WithDebounce(200*time.Millisecond),
		fsnotify.WithMaxWait(300*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	testFile := filepath.Join(tmpDir, "cleanup.file")

	events := make([]fsnotify.Event, 0)

	go func() {
		for {
			select {
			case event := <-w.Events():
				events = append(events, event)
			case err := <-w.Errors():
				t.Log("Error:", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	os.WriteFile(testFile, []byte("content"), 0644)
	time.Sleep(50 * time.Millisecond)
	os.WriteFile(testFile, []byte("more content"), 0644)
	time.Sleep(50 * time.Millisecond)

	os.Remove(testFile)

	time.Sleep(500 * time.Millisecond)

	changeEvents := 0
	removeEvents := 0
	for _, event := range events {
		if event.Path == testFile {
			switch event.Op {
			case fsnotify.Change:
				changeEvents++
			case fsnotify.Remove:
				removeEvents++
			}
		}
	}

	assert.Equal(t, 1, removeEvents, "Should receive exactly one remove event")
	assert.True(t, changeEvents <= 1, "Should receive at most one change event before removal")

	t.Logf("Received %d change events and %d remove events", changeEvents, removeEvents)
}

func TestNoDuplicateEvents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsnotify-duplicate-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	w, err := fsnotify.NewRecursiveWatcher(
		tmpDir,
		fsnotify.WithDebounce(50*time.Millisecond),
		fsnotify.WithMaxWait(100*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	testFile := filepath.Join(tmpDir, "duplicate.file")

	events := make([]fsnotify.Event, 0)
	eventMutex := sync.Mutex{}

	go func() {
		for {
			select {
			case event := <-w.Events():
				if event.Op == fsnotify.Change && event.Path == testFile {
					eventMutex.Lock()
					events = append(events, event)
					eventMutex.Unlock()
				}
			case err := <-w.Errors():
				t.Log("Error:", err)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	os.WriteFile(testFile, []byte("trigger"), 0644)
	time.Sleep(25 * time.Millisecond)
	os.WriteFile(testFile, []byte("multiple"), 0644)
	time.Sleep(25 * time.Millisecond)
	os.WriteFile(testFile, []byte("rapid"), 0644)
	time.Sleep(25 * time.Millisecond)
	os.WriteFile(testFile, []byte("updates"), 0644)

	time.Sleep(300 * time.Millisecond)

	eventMutex.Lock()
	eventCount := len(events)
	eventMutex.Unlock()

	assert.Greater(t, eventCount, 0, "Should receive at least one event")
	assert.LessOrEqual(t, eventCount, 2, "Should not receive more than 2 events (debounce + max wait)")

	t.Logf("Received %d events from 4 rapid updates (no duplicates)", eventCount)
}
