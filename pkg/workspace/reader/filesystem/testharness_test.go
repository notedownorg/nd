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
// testHarness provides a reusable test environment for filesystem reader tests.
// It encapsulates common setup, subscriber management, and invariant verification.

package filesystem

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHarness encapsulates common test setup and teardown.
type testHarness struct {
	t               *testing.T
	tmpDir          string
	reader          *Reader
	subscribers     []chan reader.Event
	subscriberViews []*subscriberView
	stopChan        chan struct{}
	subscriberWg    sync.WaitGroup
}

// newTestHarness creates a new test environment with the specified number of subscribers.
func newTestHarness(t *testing.T, numSubscribers int, loadInitial bool) *testHarness {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-test")
	require.NoError(t, err)

	r, err := NewReader(tmpDir)
	require.NoError(t, err)

	harness := &testHarness{
		t:               t,
		tmpDir:          tmpDir,
		reader:          r,
		subscribers:     make([]chan reader.Event, numSubscribers),
		subscriberViews: make([]*subscriberView, numSubscribers),
		stopChan:        make(chan struct{}),
	}

	// Set up subscribers and their event processors
	for i := 0; i < numSubscribers; i++ {
		harness.subscribers[i] = make(chan reader.Event, 2000) // Large buffer to prevent blocking
		harness.subscriberViews[i] = newSubscriberView()
		r.Subscribe(harness.subscribers[i], loadInitial)

		// Start goroutine to process events for this subscriber
		harness.subscriberWg.Add(1)
		go func(subIdx int) {
			defer harness.subscriberWg.Done()

			for {
				select {
				case <-harness.stopChan:
					return
				case event := <-harness.subscribers[subIdx]:
					harness.subscriberViews[subIdx].applyEvent(event)
				}
			}
		}(i)
	}

	return harness
}

// cleanup shuts down the test harness and cleans up resources.
func (h *testHarness) cleanup() {
	close(h.stopChan)
	h.subscriberWg.Wait()
	h.reader.Close()
	os.RemoveAll(h.tmpDir)
}

// waitForEventPropagation allows time for all events to be processed by subscribers.
func (h *testHarness) waitForEventPropagation() {
	time.Sleep(500 * time.Millisecond) // Allow events to propagate and be processed
}

// verifyEventualConsistency checks that all subscribers eventually achieve consistency
// with the filesystem state within the specified timeout. This is a fundamental
// requirement for correctness under concurrent operations.
func (h *testHarness) verifyEventualConsistency(maxWaitTime time.Duration) {
	h.t.Logf("Verifying eventual consistency within %v...", maxWaitTime)

	checkInterval := 200 * time.Millisecond
	if maxWaitTime < time.Second {
		checkInterval = maxWaitTime / 10 // Use smaller intervals for short timeouts
	}

	startTime := time.Now()
	var finalActualFiles map[string][]byte
	var finalErr error
	consistencyAchieved := false

	for time.Since(startTime) < maxWaitTime && !consistencyAchieved {
		// Get current filesystem state
		finalActualFiles, finalErr = getActualFilesystemState(h.tmpDir)
		require.NoError(h.t, finalErr)

		// Check if all subscribers have converged to the same state as the filesystem
		allConsistent := true
		subscriberConsistent := true

		for i, view := range h.subscriberViews {
			subscriberFiles := view.getFiles()

			// Check consistency with filesystem
			filesystemConsistent := len(subscriberFiles) == len(finalActualFiles)
			if filesystemConsistent {
				for fileId, expectedContent := range finalActualFiles {
					actualContent, exists := view.getFileContent(fileId)
					if !exists || !bytes.Equal(expectedContent, actualContent) {
						filesystemConsistent = false
						break
					}
				}
			}

			if !filesystemConsistent {
				allConsistent = false
				h.t.Logf("Retry %v: Subscriber %d not yet consistent with filesystem (%d vs %d files)",
					time.Since(startTime).Truncate(time.Millisecond), i, len(subscriberFiles), len(finalActualFiles))
			}
		}

		// Check consistency between subscribers
		if allConsistent && len(h.subscriberViews) > 1 {
			baseFiles := h.subscriberViews[0].getFiles()
			for i := 1; i < len(h.subscriberViews); i++ {
				currentFiles := h.subscriberViews[i].getFiles()
				if len(baseFiles) != len(currentFiles) {
					subscriberConsistent = false
					break
				}
				for fileId, baseContent := range baseFiles {
					currentContent, exists := h.subscriberViews[i].getFileContent(fileId)
					if !exists || !bytes.Equal(baseContent, currentContent) {
						subscriberConsistent = false
						break
					}
				}
				if !subscriberConsistent {
					break
				}
			}
		}

		consistencyAchieved = allConsistent && subscriberConsistent

		if !consistencyAchieved {
			time.Sleep(checkInterval)
		}
	}

	// Report results
	elapsed := time.Since(startTime).Truncate(time.Millisecond)
	h.t.Logf("Eventual consistency check completed after %v", elapsed)
	h.t.Logf("Final filesystem state: %d files", len(finalActualFiles))

	for i, view := range h.subscriberViews {
		subscriberFiles := view.getFiles()
		h.t.Logf("Subscriber %d final state: %d files", i, len(subscriberFiles))
	}

	if consistencyAchieved {
		h.t.Logf("SUCCESS: Eventual consistency achieved after %v", elapsed)
	} else {
		// Eventual consistency failure - provide detailed diagnostics
		h.t.Errorf("FAILURE: Eventual consistency was not achieved within %v timeout", maxWaitTime)
		h.logConsistencyDiagnostics(finalActualFiles)
		h.t.Fatal("Eventual consistency failure: The filesystem reader must achieve eventual consistency")
	}
}

// logConsistencyDiagnostics provides detailed information about consistency failures
// to help with debugging when eventual consistency is not achieved.
func (h *testHarness) logConsistencyDiagnostics(actualFiles map[string][]byte) {
	h.t.Log("Diagnostic information about the consistency failure:")

	for i, view := range h.subscriberViews {
		subscriberFiles := view.getFiles()
		h.t.Logf("Subscriber %d detailed state: %d files", i, len(subscriberFiles))

		if len(subscriberFiles) > 0 || len(actualFiles) > 0 {
			// Calculate and report consistency metrics for debugging
			commonFiles := 0
			correctContent := 0
			missingFiles := make([]string, 0)
			incorrectContent := make([]string, 0)

			for fileId, expectedContent := range actualFiles {
				if actualContent, exists := view.getFileContent(fileId); exists {
					commonFiles++
					if bytes.Equal(expectedContent, actualContent) {
						correctContent++
					} else {
						incorrectContent = append(incorrectContent, fileId)
					}
				} else {
					missingFiles = append(missingFiles, fileId)
				}
			}

			if len(actualFiles) > 0 {
				presenceRatio := float64(commonFiles) / float64(len(actualFiles))
				contentRatio := float64(correctContent) / float64(len(actualFiles))

				h.t.Logf("Subscriber %d consistency: %.1f%% files present (%d/%d), %.1f%% content correct (%d/%d)",
					i, presenceRatio*100, commonFiles, len(actualFiles),
					contentRatio*100, correctContent, len(actualFiles))

				// Report specific inconsistencies for debugging (limit output)
				if len(missingFiles) > 0 {
					if len(missingFiles) <= 10 {
						h.t.Logf("Subscriber %d missing files: %v", i, missingFiles)
					} else {
						h.t.Logf("Subscriber %d missing %d files (showing first 10): %v",
							i, len(missingFiles), missingFiles[:10])
					}
				}
				if len(incorrectContent) > 0 {
					if len(incorrectContent) <= 10 {
						h.t.Logf("Subscriber %d files with incorrect content: %v", i, incorrectContent)
					} else {
						h.t.Logf("Subscriber %d has %d files with incorrect content (showing first 10): %v",
							i, len(incorrectContent), incorrectContent[:10])
					}
				}
			}
		}
	}
}

// verifyInvariants checks that all subscribers have views matching the filesystem.
// This is the core invariant verification that ensures correctness.
// It assumes eventual consistency has already been achieved.
func (h *testHarness) verifyInvariants() {
	// Get the ground truth from the filesystem
	actualFiles, err := getActualFilesystemState(h.tmpDir)
	require.NoError(h.t, err)

	h.t.Logf("Filesystem contains %d files", len(actualFiles))

	for i, view := range h.subscriberViews {
		subscriberFiles := view.getFiles()
		h.t.Logf("Subscriber %d sees %d files", i, len(subscriberFiles))

		// Verify subscriber has same number of files as filesystem
		assert.Equal(h.t, len(actualFiles), len(subscriberFiles),
			"subscriber %d should have same number of files as filesystem", i)

		// Verify each file in filesystem exists in subscriber view with correct content
		for fileId, actualContent := range actualFiles {
			subscriberContent, exists := view.getFileContent(fileId)
			assert.True(h.t, exists, "subscriber %d missing file: %s", i, fileId)
			if exists {
				assert.Equal(h.t, actualContent, subscriberContent,
					"subscriber %d has incorrect content for file: %s", i, fileId)
			}
		}

		// Verify subscriber doesn't have phantom files (files that don't exist on disk)
		for fileId := range subscriberFiles {
			_, exists := actualFiles[fileId]
			assert.True(h.t, exists, "subscriber %d has phantom file: %s", i, fileId)
		}
	}

	// Verify all subscribers have identical views
	h.verifySubscriberConsistency()
}

// verifySubscriberConsistency ensures all subscribers have exactly the same view.
func (h *testHarness) verifySubscriberConsistency() {
	if len(h.subscriberViews) < 2 {
		return // Nothing to compare
	}

	baseFiles := h.subscriberViews[0].getFiles()

	for i := 1; i < len(h.subscriberViews); i++ {
		currentFiles := h.subscriberViews[i].getFiles()

		assert.Equal(h.t, len(baseFiles), len(currentFiles),
			"subscriber %d should have same number of files as subscriber 0", i)

		for fileId, baseContent := range baseFiles {
			currentContent, exists := h.subscriberViews[i].getFileContent(fileId)
			assert.True(h.t, exists, "subscriber %d missing file that subscriber 0 has: %s", i, fileId)
			if exists {
				assert.Equal(h.t, baseContent, currentContent,
					"subscribers 0 and %d have different content for file: %s", i, fileId)
			}
		}
	}
}
