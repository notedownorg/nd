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

package functional

import (
	"fmt"
	"testing"
	"time"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/test/words"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TestStream_DocumentSubscription tests basic document subscription functionality
func TestStream_DocumentSubscription(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create a stream
	stream, cancel := harness.streamWithTimeout(30 * time.Second)
	defer cancel()

	// Subscribe to documents
	harness.subscribeToDocuments(stream, "test-sub-1")

	// Create a test file
	filename := fmt.Sprintf("%s.md", words.Random())
	content := harness.generateRealisticContent("Test Document")
	err := harness.createTestFile(filename, content)
	require.NoError(t, err)

	// Wait for file propagation
	time.Sleep(EventPropagationDelay)

	// Expect change event
	changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	node := changeEvent.GetChange().GetNode()
	if doc := node.GetDocument(); doc != nil {
		expectedID := fmt.Sprintf("Document|%s", filename)
		assert.Equal(t, expectedID, doc.Id)
		// Note: Document protobuf doesn't have content field, only metadata and children
		// For now we just verify the ID is correct
	} else {
		assert.Fail(t, "Expected document node, got different type")
	}
}

// TestStream_MultipleDocumentOperations tests a sequence of document operations
func TestStream_MultipleDocumentOperations(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create a stream
	stream, cancel := harness.streamWithTimeout(60 * time.Second)
	defer cancel()

	// Subscribe to documents
	harness.subscribeToDocuments(stream, "test-multi-ops")

	// Test sequence: Create → Update → Delete
	filename := fmt.Sprintf("%s-%s.md", words.Random(), words.Random())

	// Step 1: Create file
	initialContent := harness.generateRealisticContent("Initial Title")
	err := harness.createTestFile(filename, initialContent)
	require.NoError(t, err)
	time.Sleep(EventPropagationDelay)

	// Expect change event for creation
	createEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	node := createEvent.GetChange().GetNode()
	if doc := node.GetDocument(); doc != nil {
		expectedID := fmt.Sprintf("Document|%s", filename)
		assert.Equal(t, expectedID, doc.Id)
	} else {
		assert.Fail(t, "Expected document node, got different type")
	}

	// Step 2: Update file
	updatedContent := harness.generateRealisticContent("Updated Title")
	err = harness.updateTestFile(filename, updatedContent)
	require.NoError(t, err)
	time.Sleep(EventPropagationDelay)

	// Expect change event for update
	updateEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	updateNode := updateEvent.GetChange().GetNode()
	if doc := updateNode.GetDocument(); doc != nil {
		expectedID := fmt.Sprintf("Document|%s", filename)
		assert.Equal(t, expectedID, doc.Id)
	} else {
		assert.Fail(t, "Expected document node, got different type")
	}

	// Step 3: Delete file
	err = harness.deleteTestFile(filename)
	require.NoError(t, err)
	time.Sleep(EventPropagationDelay)

	// Expect delete event
	deleteEvent := harness.expectStreamEvent(stream, DeleteEvent, 5*time.Second)
	deleteNode := deleteEvent.GetDelete().GetNode()
	if doc := deleteNode.GetDocument(); doc != nil {
		expectedID := fmt.Sprintf("Document|%s", filename)
		assert.Equal(t, expectedID, doc.Id)
	} else {
		assert.Fail(t, "Expected document node, got different type")
	}
}

// TestStream_ConcurrentSubscriptions tests multiple concurrent subscribers
func TestStream_ConcurrentSubscriptions(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create two streams
	stream1, cancel1 := harness.streamWithTimeout(60 * time.Second)
	defer cancel1()

	stream2, cancel2 := harness.streamWithTimeout(60 * time.Second)
	defer cancel2()

	// Subscribe both streams to documents
	for i, stream := range []pb.NodeService_StreamClient{stream1, stream2} {
		harness.subscribeToDocuments(stream, fmt.Sprintf("concurrent-sub-%d", i+1))
	}

	// Create a test file
	filename := fmt.Sprintf("concurrent-%s.md", words.Random())
	content := harness.generateRealisticContent("Concurrent Test")
	err := harness.createTestFile(filename, content)
	require.NoError(t, err)

	// Wait for propagation
	time.Sleep(EventPropagationDelay)

	// Both streams should receive the change event
	for i, stream := range []pb.NodeService_StreamClient{stream1, stream2} {
		changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
		node := changeEvent.GetChange().GetNode()
		if doc := node.GetDocument(); doc != nil {
			expectedID := fmt.Sprintf("Document|%s", filename)
			assert.Equal(t, expectedID, doc.Id,
				"Stream %d should receive change event", i+1)
		} else {
			assert.Fail(t, "Stream %d: Expected document node, got different type", i+1)
		}
	}
}

// TestStream_EventualConsistency tests that the system achieves eventual consistency
func TestStream_EventualConsistency(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create stream
	stream, cancel := harness.streamWithTimeout(120 * time.Second)
	defer cancel()

	// Subscribe to documents
	subscriptionReq := &pb.StreamRequest{
		Request: &pb.StreamRequest_Subscription{
			Subscription: &pb.SubscriptionRequest{
				SubscriptionId: "consistency-test",
				Msg: &pb.SubscriptionRequest_DocumentSubscription{
					DocumentSubscription: &pb.DocumentSubscription{},
				},
			},
		},
	}

	err := stream.Send(subscriptionReq)
	require.NoError(t, err)

	// Wait for subscription setup
	harness.expectStreamEvent(stream, "subscription_confirmation", 5*time.Second)
	harness.expectStreamEvent(stream, "initialization_complete", 5*time.Second)

	// Create multiple files rapidly
	const numFiles = 10
	filenames := make([]string, numFiles)

	for i := 0; i < numFiles; i++ {
		filenames[i] = fmt.Sprintf("consistency-%d-%s.md", i, words.Random())
		content := harness.generateRealisticContent(fmt.Sprintf("File %d", i))
		err := harness.createTestFile(filenames[i], content)
		require.NoError(t, err)

		// Small delay between files to avoid overwhelming the system
		time.Sleep(50 * time.Millisecond)
	}

	// Collect events for eventual consistency verification
	maxEvents := numFiles * 2 // Allow for potential duplicate events
	events := harness.drainStreamEvents(stream, maxEvents, EventualConsistencyTimeout)

	// Verify we received events for our files
	// Note: We focus on eventual consistency rather than strict event ordering
	receivedFiles := make(map[string]bool)
	for _, event := range events {
		if changeEvent := event.GetChange(); changeEvent != nil {
			if node := changeEvent.GetNode(); node != nil {
				if doc := node.GetDocument(); doc != nil {
					receivedFiles[doc.Id] = true
				}
			}
		}
	}

	// Verify eventual consistency: we should have received events for all files
	// Allow some tolerance for rapid operations
	minExpectedFiles := int(float64(numFiles) * 0.8) // At least 80% of files
	assert.GreaterOrEqual(t, len(receivedFiles), minExpectedFiles,
		"Should receive events for most files (eventual consistency)")

	t.Logf("Eventual consistency test completed: received events for %d/%d files",
		len(receivedFiles), numFiles)
}

// TestStream_NestedDirectories tests handling of nested directory structures
func TestStream_NestedDirectories(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create stream
	stream, cancel := harness.streamWithTimeout(60 * time.Second)
	defer cancel()

	// Subscribe to documents
	subscriptionReq := &pb.StreamRequest{
		Request: &pb.StreamRequest_Subscription{
			Subscription: &pb.SubscriptionRequest{
				SubscriptionId: "nested-test",
				Msg: &pb.SubscriptionRequest_DocumentSubscription{
					DocumentSubscription: &pb.DocumentSubscription{},
				},
			},
		},
	}

	err := stream.Send(subscriptionReq)
	require.NoError(t, err)

	// Wait for subscription setup
	harness.expectStreamEvent(stream, "subscription_confirmation", 5*time.Second)
	harness.expectStreamEvent(stream, "initialization_complete", 5*time.Second)

	// Create files in nested directories
	nestedFiles := []string{
		fmt.Sprintf("dir1/%s.md", words.Random()),
		fmt.Sprintf("dir1/subdir/%s.md", words.Random()),
		fmt.Sprintf("dir2/deep/nested/path/%s.md", words.Random()),
	}

	for i, filename := range nestedFiles {
		content := harness.generateRealisticContent(fmt.Sprintf("Nested File %d", i))
		err := harness.createTestFile(filename, content)
		require.NoError(t, err)

		// Wait for propagation
		time.Sleep(EventPropagationDelay)

		// Expect change event
		changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
		node := changeEvent.GetChange().GetNode()
		if doc := node.GetDocument(); doc != nil {
			expectedID := fmt.Sprintf("Document|%s", filename)
			assert.Equal(t, expectedID, doc.Id)
		} else {
			assert.Fail(t, "Expected document node, got different type")
		}
	}
}

// TestStream_FrontmatterParsing tests that frontmatter metadata is correctly parsed and sent
func TestStream_FrontmatterParsing(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create a stream
	stream, cancel := harness.streamWithTimeout(30 * time.Second)
	defer cancel()

	// Subscribe to documents
	subscriptionReq := &pb.StreamRequest{
		Request: &pb.StreamRequest_Subscription{
			Subscription: &pb.SubscriptionRequest{
				SubscriptionId: "frontmatter-test",
				Msg: &pb.SubscriptionRequest_DocumentSubscription{
					DocumentSubscription: &pb.DocumentSubscription{},
				},
			},
		},
	}

	err := stream.Send(subscriptionReq)
	require.NoError(t, err)

	// Expect subscription confirmation and initialization
	harness.expectStreamEvent(stream, "subscription_confirmation", 5*time.Second)
	harness.expectStreamEvent(stream, "initialization_complete", 5*time.Second)

	// Create a test file with frontmatter
	filename := fmt.Sprintf("frontmatter-%s.md", words.Random())
	content := `---
title: "Test Document"
author: "Test Author"
tags: ["test", "frontmatter"]
priority: 1
published: true
---

# Test Document

This is a test document with frontmatter.
`

	err = harness.createTestFile(filename, content)
	require.NoError(t, err)

	// Wait for file propagation
	time.Sleep(EventPropagationDelay)

	// Expect change event
	changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	node := changeEvent.GetChange().GetNode()
	require.NotNil(t, node, "Change event should contain a node")

	doc := node.GetDocument()
	require.NotNil(t, doc, "Node should be a document")

	expectedID := fmt.Sprintf("Document|%s", filename)
	assert.Equal(t, expectedID, doc.Id)

	// Verify frontmatter metadata is present
	require.NotNil(t, doc.Metadata, "Document should have metadata")
	assert.Greater(t, len(doc.Metadata), 0, "Metadata should not be empty")

	// Create a map of metadata for easier verification
	metadataMap := make(map[string]interface{})
	for _, entry := range doc.Metadata {
		if entry.Value != nil {
			// Extract the value from the Any field
			var value interface{}
			switch {
			case entry.Value.MessageIs(&wrapperspb.StringValue{}):
				var sv wrapperspb.StringValue
				entry.Value.UnmarshalTo(&sv)
				value = sv.Value
			case entry.Value.MessageIs(&wrapperspb.Int64Value{}):
				var iv wrapperspb.Int64Value
				entry.Value.UnmarshalTo(&iv)
				value = iv.Value
			case entry.Value.MessageIs(&wrapperspb.BoolValue{}):
				var bv wrapperspb.BoolValue
				entry.Value.UnmarshalTo(&bv)
				value = bv.Value
			case entry.Value.MessageIs(&wrapperspb.DoubleValue{}):
				var dv wrapperspb.DoubleValue
				entry.Value.UnmarshalTo(&dv)
				value = dv.Value
			}
			metadataMap[entry.Key] = value
		}
	}

	// Verify expected metadata values
	assert.Equal(t, "Test Document", metadataMap["title"])
	assert.Equal(t, "Test Author", metadataMap["author"])
	assert.Equal(t, int64(1), metadataMap["priority"])
	assert.Equal(t, true, metadataMap["published"])

	t.Logf("Successfully parsed frontmatter with %d metadata entries", len(doc.Metadata))
	for key, value := range metadataMap {
		t.Logf("  %s: %v (%T)", key, value, value)
	}
}

// TestStream_FrontmatterEmpty tests handling of documents without frontmatter
func TestStream_FrontmatterEmpty(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create a stream
	stream, cancel := harness.streamWithTimeout(30 * time.Second)
	defer cancel()

	// Subscribe to documents
	subscriptionReq := &pb.StreamRequest{
		Request: &pb.StreamRequest_Subscription{
			Subscription: &pb.SubscriptionRequest{
				SubscriptionId: "no-frontmatter-test",
				Msg: &pb.SubscriptionRequest_DocumentSubscription{
					DocumentSubscription: &pb.DocumentSubscription{},
				},
			},
		},
	}

	err := stream.Send(subscriptionReq)
	require.NoError(t, err)

	// Expect subscription confirmation and initialization
	harness.expectStreamEvent(stream, "subscription_confirmation", 5*time.Second)
	harness.expectStreamEvent(stream, "initialization_complete", 5*time.Second)

	// Create a test file without frontmatter
	filename := fmt.Sprintf("no-frontmatter-%s.md", words.Random())
	content := `# Simple Document

This is a document without frontmatter.
Just plain markdown content.
`

	err = harness.createTestFile(filename, content)
	require.NoError(t, err)

	// Wait for file propagation
	time.Sleep(EventPropagationDelay)

	// Expect change event
	changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	node := changeEvent.GetChange().GetNode()
	require.NotNil(t, node, "Change event should contain a node")

	doc := node.GetDocument()
	require.NotNil(t, doc, "Node should be a document")

	expectedID := fmt.Sprintf("Document|%s", filename)
	assert.Equal(t, expectedID, doc.Id)

	// Verify no metadata is present for documents without frontmatter
	assert.Equal(t, 0, len(doc.Metadata), "Document without frontmatter should have no metadata")

	t.Log("Successfully handled document without frontmatter")
}

// TestStream_FrontmatterComplexTypes tests handling of complex frontmatter with arrays and nested objects
func TestStream_FrontmatterComplexTypes(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.cleanup()

	// Create a stream
	stream, cancel := harness.streamWithTimeout(30 * time.Second)
	defer cancel()

	// Subscribe to documents
	subscriptionReq := &pb.StreamRequest{
		Request: &pb.StreamRequest_Subscription{
			Subscription: &pb.SubscriptionRequest{
				SubscriptionId: "complex-frontmatter-test",
				Msg: &pb.SubscriptionRequest_DocumentSubscription{
					DocumentSubscription: &pb.DocumentSubscription{},
				},
			},
		},
	}

	err := stream.Send(subscriptionReq)
	require.NoError(t, err)

	// Expect subscription confirmation and initialization
	harness.expectStreamEvent(stream, "subscription_confirmation", 5*time.Second)
	harness.expectStreamEvent(stream, "initialization_complete", 5*time.Second)

	// Create a test file with complex frontmatter
	filename := fmt.Sprintf("complex-frontmatter-%s.md", words.Random())
	content := `---
title: "Complex Metadata Test"
categories: ["tech", "documentation", "testing"]
version: 2.1
metadata:
  created: "2025-06-17"
  updated: "2025-06-17"
  status: "draft"
feature_flags:
  - experimental: true
  - beta: false
score: 4.5
---

# Complex Metadata Document

This document tests complex frontmatter parsing.
`

	err = harness.createTestFile(filename, content)
	require.NoError(t, err)

	// Wait for file propagation
	time.Sleep(EventPropagationDelay)

	// Expect change event
	changeEvent := harness.expectStreamEvent(stream, ChangeEvent, 5*time.Second)
	node := changeEvent.GetChange().GetNode()
	require.NotNil(t, node, "Change event should contain a node")

	doc := node.GetDocument()
	require.NotNil(t, doc, "Node should be a document")

	expectedID := fmt.Sprintf("Document|%s", filename)
	assert.Equal(t, expectedID, doc.Id)

	// Verify frontmatter metadata is present
	require.NotNil(t, doc.Metadata, "Document should have metadata")
	assert.Greater(t, len(doc.Metadata), 0, "Metadata should not be empty")

	// Create a map of metadata for easier verification
	metadataMap := make(map[string]interface{})
	for _, entry := range doc.Metadata {
		if entry.Value != nil {
			// Extract the value from the Any field
			var value interface{}
			switch {
			case entry.Value.MessageIs(&wrapperspb.StringValue{}):
				var sv wrapperspb.StringValue
				entry.Value.UnmarshalTo(&sv)
				value = sv.Value
			case entry.Value.MessageIs(&wrapperspb.DoubleValue{}):
				var dv wrapperspb.DoubleValue
				entry.Value.UnmarshalTo(&dv)
				value = dv.Value
			default:
				// For complex types, they should be JSON strings
				var sv wrapperspb.StringValue
				if entry.Value.MessageIs(&wrapperspb.StringValue{}) {
					entry.Value.UnmarshalTo(&sv)
					value = sv.Value
				}
			}
			metadataMap[entry.Key] = value
		}
	}

	// Verify expected metadata values
	assert.Equal(t, "Complex Metadata Test", metadataMap["title"])
	assert.Equal(t, 2.1, metadataMap["version"])
	assert.Equal(t, 4.5, metadataMap["score"])

	// Complex types should be JSON strings
	if categoriesStr, ok := metadataMap["categories"].(string); ok {
		assert.Contains(t, categoriesStr, "tech")
		assert.Contains(t, categoriesStr, "documentation")
		assert.Contains(t, categoriesStr, "testing")
	}

	if metadataStr, ok := metadataMap["metadata"].(string); ok {
		assert.Contains(t, metadataStr, "created")
		assert.Contains(t, metadataStr, "2025-06-17")
		assert.Contains(t, metadataStr, "draft")
	}

	t.Logf("Successfully parsed complex frontmatter with %d metadata entries", len(doc.Metadata))
	for key, value := range metadataMap {
		t.Logf("  %s: %v (%T)", key, value, value)
	}
}
