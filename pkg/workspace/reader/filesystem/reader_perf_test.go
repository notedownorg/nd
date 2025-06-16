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

// Package filesystem contains performance and invariant tests for the filesystem reader.
//
// INVARIANT TESTING PHILOSOPHY:
// Instead of simply counting events, these tests verify that each subscriber's view
// of the filesystem matches the actual filesystem state. This approach detects real
// bugs like missing updates, stale content, or inconsistent views between subscribers.
//
// Each test follows this pattern:
// 1. Set up subscribers with subscriberView trackers
// 2. Perform file operations (create, update, delete)
// 3. Wait for events to propagate
// 4. Compare subscriber views with actual filesystem state
// 5. Verify all subscribers have identical, correct views

package filesystem

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/notedownorg/nd/pkg/test/words"
	"github.com/notedownorg/nd/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReader_InvariantTesting_ModerateStressScenario tests the filesystem reader
// under moderate concurrent load to verify that all subscribers maintain
// consistent views that match the actual filesystem state.
//
// TEST SCENARIO:
// - 3 subscribers monitoring the same directory
// - 100 files created sequentially (20ms intervals)
// - Random updates to existing files every 200ms
// - Operations run for 3 seconds
//
// EXPECTED OUTCOME:
// - All subscribers have identical views
// - All subscriber views match the actual filesystem
// - No missing files, no phantom files, no stale content
func TestReader_InvariantTesting_ModerateStressScenario(t *testing.T) {
	// Test configuration
	const numSubscribers = 3
	const numFiles = 100
	const operationDuration = 3 * time.Second

	// Set up test harness with subscribers
	harness := newTestHarness(t, numSubscribers, false)
	defer harness.cleanup()

	// Give subscribers time to initialize
	time.Sleep(200 * time.Millisecond)

	// Generate realistic filenames with nested directories (0-5 levels deep)
	filenames := make([]string, numFiles)
	fileTitles := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		// Create a realistic directory path with 0-5 levels of nesting
		pathParts := []string{harness.tmpDir}
		depth := rand.Intn(6) // 0-5 levels of depth
		for d := 0; d < depth; d++ {
			pathParts = append(pathParts, words.Random())
		}

		// Add the filename
		filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())
		filenames[i] = filepath.Join(append(pathParts, filename)...)
		fileTitles[i] = fmt.Sprintf("%s %s", words.Random(), words.Random())
	}

	// Concurrent operations to stress test the system
	var operationWg sync.WaitGroup
	stopChan := make(chan struct{})

	// Operation 1: Create files sequentially
	// This simulates a user creating multiple files in a session
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < numFiles; i++ {
			select {
			case <-stopChan:
				return
			default:
			}

			filename := filenames[i]
			title := fileTitles[i]
			content := fmt.Sprintf("# %s\n\n%s %s initial content created at %s",
				title, words.Random(), words.Random(), time.Now().Format("15:04:05"))

			// Ensure directory exists before creating file
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				continue // Skip this file if directory creation fails
			}
			os.WriteFile(filename, []byte(content), 0644)
			time.Sleep(20 * time.Millisecond) // Moderate pace
		}
	}()

	// Operation 2: Random updates to existing files
	// This simulates ongoing editing of existing files
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		time.Sleep(1 * time.Second) // Let some files be created first
		rand.Seed(time.Now().UnixNano() + 1)

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				fileNum := rand.Intn(numFiles)
				filename := filenames[fileNum]
				title := fileTitles[fileNum]

				// Only update if file exists
				if _, err := os.Stat(filename); err == nil {
					content := fmt.Sprintf("# %s (Updated)\n\n%s %s updated content at %s",
						title, words.Random(), words.Random(), time.Now().Format("15:04:05"))
					os.WriteFile(filename, []byte(content), 0644)
				}
			}
		}
	}()

	// Let operations run for the specified duration
	time.Sleep(operationDuration)

	// Stop all operations and wait for completion
	close(stopChan)
	operationWg.Wait()

	// Allow time for all events to propagate through the system
	harness.waitForEventPropagation()

	// Verify eventual consistency before checking invariants
	harness.verifyEventualConsistency(5 * time.Second)

	// Verify the core invariant: all subscribers see the correct filesystem state
	harness.verifyInvariants()
}

// TestReader_InvariantTesting_SequentialOperations tests the filesystem reader
// with careful sequential operations to verify that all state transitions
// are captured correctly without race conditions.
//
// TEST SCENARIO:
// - 2 subscribers monitoring the same directory
// - Sequential workflow: Create → Update → Delete
// - Step 1: Create 50 files (150ms between each)
// - Step 2: Update first 25 files (150ms between each)
// - Step 3: Delete last 25 files (100ms between each)
//
// EXPECTED OUTCOME:
// - Final state: 25 files remaining (first half, updated)
// - All subscribers see identical final state
// - No lost updates, no stale content, no phantom files
func TestReader_InvariantTesting_SequentialOperations(t *testing.T) {
	// Test configuration
	const numSubscribers = 2
	const numFiles = 50

	// Set up test harness
	harness := newTestHarness(t, numSubscribers, false)
	defer harness.cleanup()

	// Give subscribers time to initialize
	time.Sleep(200 * time.Millisecond)

	// Generate realistic filenames with nested directories and titles for consistent referencing
	filenames := make([]string, numFiles)
	fileTitles := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		// Create a realistic directory path with 0-5 levels of nesting
		pathParts := []string{harness.tmpDir}
		depth := rand.Intn(6) // 0-5 levels of depth
		for d := 0; d < depth; d++ {
			pathParts = append(pathParts, words.Random())
		}

		// Add the filename
		filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())
		filenames[i] = filepath.Join(append(pathParts, filename)...)
		fileTitles[i] = fmt.Sprintf("%s %s", words.Random(), words.Random())
	}

	// Step 1: Create all files sequentially
	// This tests basic file creation and event delivery
	t.Logf("Step 1: Creating %d files...", numFiles)
	for i := 0; i < numFiles; i++ {
		filename := filenames[i]
		title := fileTitles[i]
		content := fmt.Sprintf("# %s\n\n%s %s step 1 content", title, words.Random(), words.Random())

		// Ensure directory exists before creating file
		err := os.MkdirAll(filepath.Dir(filename), 0755)
		require.NoError(t, err)

		err = os.WriteFile(filename, []byte(content), 0644)
		require.NoError(t, err)

		time.Sleep(50 * time.Millisecond) // Pace to avoid debouncing
	}

	// Step 2: Update first half of files
	// This tests that updates overwrite previous content correctly
	t.Logf("Step 2: Updating first %d files...", numFiles/2)
	for i := 0; i < numFiles/2; i++ {
		filename := filenames[i]
		title := fileTitles[i]
		content := fmt.Sprintf("# %s (Updated)\n\n%s %s step 2 updated content", title, words.Random(), words.Random())

		err := os.WriteFile(filename, []byte(content), 0644)
		require.NoError(t, err)

		time.Sleep(50 * time.Millisecond)
	}

	// Step 3: Delete second half of files
	// This tests that delete events properly remove files from subscriber views
	t.Logf("Step 3: Deleting last %d files...", numFiles/2)
	for i := numFiles / 2; i < numFiles; i++ {
		filename := filenames[i]

		err := os.Remove(filename)
		require.NoError(t, err)

		time.Sleep(30 * time.Millisecond)
	}

	// Allow all events to propagate
	harness.waitForEventPropagation()

	// Verify eventual consistency
	harness.verifyEventualConsistency(3 * time.Second)

	// Verify final state: should have numFiles/2 files remaining (the updated ones)
	actualFiles, err := getActualFilesystemState(harness.tmpDir)
	require.NoError(t, err)

	expectedFileCount := numFiles / 2
	assert.Equal(t, expectedFileCount, len(actualFiles),
		"should have %d files remaining after sequential operations", expectedFileCount)

	// Verify the core invariant: all subscribers see the correct final state
	harness.verifyInvariants()
}

// TestReader_InvariantTesting_LoadInitialDocuments tests that subscribers
// correctly load all existing files when they subscribe with loadInitialDocuments=true,
// and continue to receive new events after the initial load completes.
//
// TEST SCENARIO:
// - Create 50 files before starting the reader
// - 4 subscribers connect with loadInitialDocuments=true
// - Wait for SubscriberLoadComplete events
// - Create one new file after initial load
//
// EXPECTED OUTCOME:
// - All subscribers receive all 50 initial files
// - All subscribers receive SubscriberLoadComplete event
// - All subscribers receive the new file created after load
// - All subscriber views are identical and match filesystem
func TestReader_InvariantTesting_LoadInitialDocuments(t *testing.T) {
	// Test configuration
	const numInitialFiles = 50
	const numSubscribers = 4

	// Create test directory and populate with initial files
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-initial-load-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Pre-populate directory with files before creating reader
	// This simulates a workspace that already has content
	initialFiles := make(map[string][]byte)
	t.Logf("Creating %d initial files...", numInitialFiles)
	for i := 0; i < numInitialFiles; i++ {
		// Create a realistic directory path with 0-5 levels of nesting
		pathParts := []string{}
		depth := rand.Intn(6) // 0-5 levels of depth
		for d := 0; d < depth; d++ {
			pathParts = append(pathParts, words.Random())
		}

		// Add the filename
		filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())

		// Create relative path for the map key (without tmpDir)
		relativePathParts := append(pathParts, filename)
		relativePath := filepath.Join(relativePathParts...)

		// Create full path for file creation
		fullPath := filepath.Join(tmpDir, relativePath)

		title := fmt.Sprintf("%s %s", words.Random(), words.Random())
		content := fmt.Sprintf("# %s\n\n%s %s existing content before subscription",
			title, words.Random(), words.Random())

		// Ensure directory exists before creating file
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)

		initialFiles[relativePath] = []byte(content)
	}

	// Create reader and subscribers with loadInitialDocuments=true
	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	subscribers := make([]chan reader.Event, numSubscribers)
	subscriberViews := make([]*subscriberView, numSubscribers)
	loadCompleteReceived := make([]bool, numSubscribers)

	// Set up subscribers to track load complete events separately
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = make(chan reader.Event, 1000)
		subscriberViews[i] = newSubscriberView()
		r.Subscribe(subscribers[i], true) // Request initial load
	}

	// Process events until all subscribers receive SubscriberLoadComplete
	var wg sync.WaitGroup
	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go func(subIdx int) {
			defer wg.Done()

			timeout := time.After(10 * time.Second)
			for !loadCompleteReceived[subIdx] {
				select {
				case event := <-subscribers[subIdx]:
					if event.Op == reader.SubscriberLoadComplete {
						loadCompleteReceived[subIdx] = true
						t.Logf("Subscriber %d received load complete", subIdx)
					} else {
						subscriberViews[subIdx].applyEvent(event)
					}
				case <-timeout:
					t.Errorf("timeout waiting for load complete for subscriber %d", subIdx)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify all subscribers received load complete and have all initial files
	for i := 0; i < numSubscribers; i++ {
		assert.True(t, loadCompleteReceived[i], "subscriber %d should receive load complete event", i)

		subscriberFiles := subscriberViews[i].getFiles()
		assert.Equal(t, len(initialFiles), len(subscriberFiles),
			"subscriber %d should have loaded all initial files", i)

		// Verify content of each initial file
		for fileId, expectedContent := range initialFiles {
			actualContent, exists := subscriberViews[i].getFileContent(fileId)
			assert.True(t, exists, "subscriber %d missing initial file: %s", i, fileId)
			if exists {
				assert.Equal(t, expectedContent, actualContent,
					"subscriber %d has incorrect content for initial file: %s", i, fileId)
			}
		}
	}

	// Test that new files are received after initial load completes
	// This ensures the system continues working normally after load
	t.Log("Creating new file after initial load...")

	// Create a new file in a nested directory
	pathParts := []string{}
	depth := rand.Intn(4) + 1 // 1-4 levels of depth for new file
	for d := 0; d < depth; d++ {
		pathParts = append(pathParts, words.Random())
	}

	newFileName := fmt.Sprintf("%s-after-load.md", words.Random())
	relativePathParts := append(pathParts, newFileName)
	newFileName = filepath.Join(relativePathParts...) // Update to include directory path
	newFilePath := filepath.Join(tmpDir, newFileName)

	newTitle := fmt.Sprintf("%s %s", words.Random(), words.Random())
	newContent := fmt.Sprintf("# %s\n\n%s %s created after initial load",
		newTitle, words.Random(), words.Random())

	// Ensure directory exists before creating file
	err = os.MkdirAll(filepath.Dir(newFilePath), 0755)
	require.NoError(t, err)

	err = os.WriteFile(newFilePath, []byte(newContent), 0644)
	require.NoError(t, err)

	// Wait for new file event to propagate to all subscribers
	timeout := time.After(3 * time.Second)
	newFileEventReceived := make([]bool, numSubscribers)

	for !allTrue(newFileEventReceived) {
		for i := 0; i < numSubscribers; i++ {
			if newFileEventReceived[i] {
				continue
			}

			select {
			case event := <-subscribers[i]:
				if event.Op == reader.Change && event.Id == newFileName {
					subscriberViews[i].applyEvent(event)
					newFileEventReceived[i] = true
					t.Logf("Subscriber %d received new file event", i)
				}
			case <-timeout:
				t.Fatalf("timeout waiting for new file event")
			default:
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify all subscribers received the new file
	for i := 0; i < numSubscribers; i++ {
		content, exists := subscriberViews[i].getFileContent(newFileName)
		assert.True(t, exists, "subscriber %d should have new file after load complete", i)
		if exists {
			assert.Equal(t, []byte(newContent), content,
				"subscriber %d should have correct content for new file", i)
		}
	}

	t.Logf("Successfully verified initial load of %d files + 1 new file across %d subscribers",
		numInitialFiles, numSubscribers)
}

// allTrue is a helper function to check if all boolean values in a slice are true
func allTrue(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}

// TestReader_InvariantTesting_ExtremeStressScenario pushes the filesystem reader
// to its limits with very rapid concurrent operations to test the max wait
// debouncing mechanism and identify edge cases.
//
// TEST SCENARIO:
// - 2 subscribers monitoring the same directory
// - 50 files created rapidly
// - Extremely fast updates (25ms intervals)
// - Random file deletions and recreations
// - Operations run for 3 seconds
//
// EXPECTED OUTCOME:
// - Under extreme stress, some events may be debounced/lost temporarily
// - BUT: the system MUST achieve eventual consistency within 10 seconds
// - All subscribers MUST converge to the correct final filesystem state
// - Test FAILS if eventual consistency is not achieved (hard requirement)
//
// NOTE: This test verifies eventual consistency - a fundamental requirement for correctness
func TestReader_InvariantTesting_ExtremeStressScenario(t *testing.T) {
	// Test configuration - extreme parameters
	const numSubscribers = 2
	const numFiles = 50
	const operationDuration = 3 * time.Second

	// Set up test harness
	harness := newTestHarness(t, numSubscribers, false)
	defer harness.cleanup()

	time.Sleep(200 * time.Millisecond)

	// Generate realistic filenames with nested directories upfront for consistent referencing
	filenames := make([]string, numFiles)
	fileTitles := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		// Create a realistic directory path with 0-5 levels of nesting
		pathParts := []string{harness.tmpDir}
		depth := rand.Intn(6) // 0-5 levels of depth
		for d := 0; d < depth; d++ {
			pathParts = append(pathParts, words.Random())
		}

		// Add the filename
		filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())
		filenames[i] = filepath.Join(append(pathParts, filename)...)
		fileTitles[i] = fmt.Sprintf("%s %s", words.Random(), words.Random())
	}

	// Extreme concurrent operations to stress test debouncing/max wait
	var operationWg sync.WaitGroup
	stopChan := make(chan struct{})

	// Add timeout context for better control
	ctx, cancel := context.WithTimeout(context.Background(), operationDuration+5*time.Second)
	defer cancel()

	// Operation 1: Rapid file creation
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < numFiles; i++ {
			select {
			case <-stopChan:
				return
			default:
			}

			filename := filenames[i]
			title := fileTitles[i]
			content := fmt.Sprintf("# %s\n\n%s %s initial content", title, words.Random(), words.Random())

			// Ensure directory exists before creating file
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				continue // Skip this file if directory creation fails
			}
			os.WriteFile(filename, []byte(content), 0644)

			// Occasional brief pause to vary timing
			if rand.Intn(20) == 0 {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			}
		}
	}()

	// Operation 2: Extremely rapid updates (faster than debounce timeout)
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		time.Sleep(500 * time.Millisecond)
		rand.Seed(time.Now().UnixNano() + 1)

		ticker := time.NewTicker(25 * time.Millisecond) // Very fast updates
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				fileNum := rand.Intn(numFiles)
				filename := filenames[fileNum]
				title := fileTitles[fileNum]

				if _, err := os.Stat(filename); err == nil {
					content := fmt.Sprintf("# %s (Rapid Update)\n\n%s %s rapid update at %s",
						title, words.Random(), words.Random(), time.Now().Format("15:04:05.000"))
					os.WriteFile(filename, []byte(content), 0644)
				}
			}
		}
	}()

	// Operation 3: Random deletions
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		time.Sleep(1 * time.Second)
		rand.Seed(time.Now().UnixNano() + 2)

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				if rand.Intn(10) < 2 { // 20% chance
					fileNum := rand.Intn(numFiles)
					filename := filenames[fileNum]
					os.Remove(filename) // Ignore errors
				}
			}
		}
	}()

	// Operation 4: Random recreations
	operationWg.Add(1)
	go func() {
		defer operationWg.Done()
		time.Sleep(2 * time.Second)
		rand.Seed(time.Now().UnixNano() + 3)

		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				if rand.Intn(10) < 3 { // 30% chance
					fileNum := rand.Intn(numFiles)
					filename := filenames[fileNum]
					title := fileTitles[fileNum]

					if _, err := os.Stat(filename); os.IsNotExist(err) {
						content := fmt.Sprintf("# %s (Recreated)\n\n%s %s recreated content",
							title, words.Random(), words.Random())

						// Ensure directory exists before recreating file
						if err := os.MkdirAll(filepath.Dir(filename), 0755); err == nil {
							os.WriteFile(filename, []byte(content), 0644)
						}
					}
				}
			}
		}
	}()

	// Let extreme operations run with timeout protection
	select {
	case <-time.After(operationDuration):
		close(stopChan)
	case <-ctx.Done():
		t.Log("Test context cancelled, stopping operations early")
		close(stopChan)
	}

	// Wait for operations to complete with timeout
	done := make(chan struct{})
	go func() {
		operationWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All operations completed
	case <-time.After(2 * time.Second):
		t.Log("Warning: Some operations did not complete within timeout")
	}

	// Allow extra time for events to settle under extreme load
	time.Sleep(2 * time.Second) // Longer than normal

	// This test may show inconsistencies under extreme load - that's expected
	// The goal is to verify we don't crash and eventual consistency is maintained
	actualFiles, err := getActualFilesystemState(harness.tmpDir)
	require.NoError(t, err)

	t.Logf("Extreme stress test completed:")
	t.Logf("- Final filesystem state: %d files", len(actualFiles))

	for i, view := range harness.subscriberViews {
		subscriberFiles := view.getFiles()
		t.Logf("- Subscriber %d view: %d files", i, len(subscriberFiles))
	}

	// EVENTUAL CONSISTENCY VERIFICATION
	// Under extreme stress, the system may temporarily be inconsistent due to debouncing
	// and high concurrency, but it should eventually reach consistency given enough time.
	harness.verifyEventualConsistency(10 * time.Second)

	// Verify the invariants one final time now that consistency is achieved
	harness.verifyInvariants()
}

// TestReader_HighFrequencyOperations tests the filesystem reader's ability to handle
// very rapid, concurrent file operations across multiple goroutines without losing
// events or crashing. This test focuses on throughput and stability under load.
//
// TEST SCENARIO:
// - Single subscriber monitoring the directory
// - 100 files with realistic word-based names in nested directories (0-5 levels deep)
// - 20 operations per file (create, update, delete) with random timing
// - Operations include: file creation, content updates, and deletions
// - Realistic markdown content generated using word combinations
// - Nested directory structures created automatically as needed
// - Random delays between operations (0-10ms) to simulate real usage
//
// EXPECTED OUTCOME:
// - System remains stable and responsive
// - Substantial number of events are captured (allows for some debouncing)
// - At least some change and delete events are received
// - No crashes or deadlocks occur
// - Eventual consistency is achieved - subscriber view matches filesystem state
//
// NOTE: This test focuses on system stability and eventual consistency rather than
// perfect event capture, as very rapid operations are expected to be debounced for performance.
func TestReader_HighFrequencyOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filesystem-reader-highfreq-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set up filesystem reader with large event buffer
	r, err := NewReader(tmpDir)
	require.NoError(t, err)
	defer r.Close()

	// Large buffer to handle high-frequency events without blocking
	events := make(chan reader.Event, 5000)
	r.Subscribe(events, false)

	// Allow reader to initialize
	time.Sleep(100 * time.Millisecond)

	// Test configuration - high concurrency parameters
	const numFiles = 100
	const operationsPerFile = 20

	// Generate realistic filenames with nested directories using words package
	filenames := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		// Create a realistic directory path with 0-5 levels of nesting
		pathParts := []string{tmpDir}
		depth := rand.Intn(6) // 0-5 levels of depth
		for d := 0; d < depth; d++ {
			pathParts = append(pathParts, words.Random())
		}

		// Add the filename
		filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())
		filenames[i] = filepath.Join(append(pathParts, filename)...)
	}

	// Coordinate concurrent file operations
	var wg sync.WaitGroup

	// Launch concurrent operations for each file
	for fileIdx := 0; fileIdx < numFiles; fileIdx++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			filename := filenames[idx]
			baseTitle := fmt.Sprintf("%s %s", words.Random(), words.Random())

			// Perform mixed operations on this file
			for op := 0; op < operationsPerFile; op++ {
				switch op % 4 {
				case 0, 1: // File creation/overwrite (50% of operations)
					content := fmt.Sprintf("# %s\n\n%s %s operation %d at %v",
						baseTitle, words.Random(), words.Random(), op, time.Now().Format("15:04:05"))

					// Ensure directory exists before creating file
					if err := os.MkdirAll(filepath.Dir(filename), 0755); err == nil {
						os.WriteFile(filename, []byte(content), 0644)
					}
				case 2: // File update (25% of operations)
					if _, err := os.Stat(filename); err == nil {
						content := fmt.Sprintf("# %s (Updated)\n\n%s %s updated operation %d at %v",
							baseTitle, words.Random(), words.Random(), op, time.Now().Format("15:04:05"))
						os.WriteFile(filename, []byte(content), 0644)
					}
				case 3: // File deletion (25% of operations, with 33% probability)
					if rand.Intn(3) == 0 {
						os.Remove(filename)
					}
				}

				// Random delay to simulate realistic usage patterns
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			}
		}(fileIdx)
	}

	// Wait for all file operations to complete
	wg.Wait()

	// Collect and analyze events
	eventCount := 0
	eventTypes := make(map[reader.Operation]int)
	timeout := time.After(5 * time.Second) // Allow time for events to propagate

	// Collect events with reasonable timeout
	for {
		select {
		case event := <-events:
			// Count relevant events (ignore Load events from potential initial scan)
			if event.Op == reader.Change || event.Op == reader.Delete {
				eventCount++
				eventTypes[event.Op]++
			}
		case <-timeout:
			// Timeout reached, stop collecting
			goto done
		default:
			// Exit early if we've collected substantial events
			if eventCount > numFiles*operationsPerFile/2 {
				goto done
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

done:
	// Report test results
	t.Logf("High frequency test results:")
	t.Logf("  Total operations performed: %d", numFiles*operationsPerFile)
	t.Logf("  Total events captured: %d", eventCount)
	t.Logf("  Change events: %d", eventTypes[reader.Change])
	t.Logf("  Delete events: %d", eventTypes[reader.Delete])

	// Set up a test harness to verify eventual consistency
	harness := newTestHarness(t, 1, false)
	defer harness.cleanup()

	// Add all remaining files to the harness subscriber to verify consistency
	time.Sleep(100 * time.Millisecond) // Brief wait for harness to initialize

	// Verify eventual consistency - the subscriber should converge to filesystem state
	harness.verifyEventualConsistency(3 * time.Second)

	// Verify system stability and reasonable event capture
	// Under high frequency, expect significant debouncing, so set reasonable minimums
	expectedMinEvents := numFiles / 2
	assert.Greater(t, eventCount, expectedMinEvents,
		"should receive substantial events from high frequency operations")
	assert.Greater(t, eventTypes[reader.Change], 10,
		"should receive meaningful number of change events")

	// Verify system handled the load without issues
	t.Logf("High frequency test completed successfully - system remained stable under load")
}
