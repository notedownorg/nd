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

// Note: Test constants and helper functions have been moved to testharness_test.go

func TestReader_Subscribe(t *testing.T) {
	harness := newSingleSubscriberHarness(t, false)
	defer harness.cleanup()

	events := harness.getSingleSubscriber()
	subscriptionID := harness.reader.Subscribe(events, false)
	assert.Greater(t, subscriptionID, 0)

	harness.reader.Unsubscribe(subscriptionID)
}

func TestReader_LoadInitialDocuments(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.md")
	testContent := "# Test Document\nThis is a test."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Add delay to ensure file is fully written before starting watcher
	time.Sleep(FileWriteDelay)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	events := make(chan reader.Event, DefaultEventBufferSize)
	r.Subscribe(events, true)

	// Use test harness for event validation but with external reader
	harness := &testHarness{t: t}

	select {
	case loadEvent := <-events:
		harness.assertEventMatch(loadEvent, reader.Load, "test.md", []byte(testContent))
	case <-time.After(DefaultTimeout):
		t.Fatal("timeout waiting for load event")
	}

	select {
	case completeEvent := <-events:
		harness.assertEventMatch(completeEvent, reader.SubscriberLoadComplete, "", nil)
	case <-time.After(DefaultTimeout):
		t.Fatal("timeout waiting for load complete event")
	}
}

func TestReader_FileChanges(t *testing.T) {
	harness := newSingleSubscriberHarness(t, false)
	defer harness.cleanup()

	events := harness.getSingleSubscriber()
	harness.reader.Subscribe(events, false)

	testContent := "# New Document"
	err := harness.createTestFile("new.md", testContent)
	require.NoError(t, err)

	changeEvent := harness.waitForEvent(reader.Change, ChangeEventTimeout)
	harness.assertEventMatch(changeEvent, reader.Change, "new.md", []byte(testContent))
}

func TestReader_OnlyMarkdownFiles(t *testing.T) {
	harness := newSingleSubscriberHarness(t, false)
	defer harness.cleanup()

	events := harness.getSingleSubscriber()
	harness.reader.Subscribe(events, false)

	// Create non-markdown file (should be ignored)
	txtFile := filepath.Join(harness.tmpDir, "test.txt")
	err := os.WriteFile(txtFile, []byte("not markdown"), 0644)
	require.NoError(t, err)

	// Should not receive event for non-markdown file
	harness.expectNoEvent(NonMarkdownTimeout)
}
