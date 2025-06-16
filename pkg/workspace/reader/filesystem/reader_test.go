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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Subscribe(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	events := make(chan reader.Event, 10)
	subscriptionID := r.Subscribe(events, false)

	assert.Greater(t, subscriptionID, 0)

	r.Unsubscribe(subscriptionID)
}

func TestReader_LoadInitialDocuments(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.md")
	testContent := "# Test Document\nThis is a test."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Add a small delay to ensure file is fully written before starting watcher
	time.Sleep(100 * time.Millisecond)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	events := make(chan reader.Event, 10)
	r.Subscribe(events, true)

	var loadEvent, completeEvent reader.Event
	select {
	case loadEvent = <-events:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for load event")
	}

	select {
	case completeEvent = <-events:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for load complete event")
	}

	assert.Equal(t, reader.Load, loadEvent.Op)
	assert.Equal(t, "test.md", loadEvent.Id)
	assert.Equal(t, []byte(testContent), loadEvent.Content)

	assert.Equal(t, reader.SubscriberLoadComplete, completeEvent.Op)
}

func TestReader_FileChanges(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	events := make(chan reader.Event, 10)
	r.Subscribe(events, false)

	time.Sleep(100 * time.Millisecond)

	testFile := filepath.Join(tmpDir, "new.md")
	testContent := "# New Document"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	var changeEvent reader.Event
	select {
	case changeEvent = <-events:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for change event")
	}

	assert.Equal(t, reader.Change, changeEvent.Op)
	assert.Equal(t, "new.md", changeEvent.Id)
	assert.Equal(t, []byte(testContent), changeEvent.Content)
}

func TestReader_OnlyMarkdownFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	events := make(chan reader.Event, 10)
	r.Subscribe(events, false)

	time.Sleep(100 * time.Millisecond)

	txtFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(txtFile, []byte("not markdown"), 0644)
	require.NoError(t, err)

	select {
	case <-events:
		t.Fatal("should not receive event for non-markdown file")
	case <-time.After(500 * time.Millisecond):
	}
}
