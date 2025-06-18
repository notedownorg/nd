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

package server

import (
	"context"
	"testing"
	"time"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/workspace/reader/mock"
)

func TestDocumentSubscription(t *testing.T) {
	// Create a server with a mock workspace
	server := NewServer()
	mockReader := mock.NewReader(t)

	err := server.RegisterWorkspace("test-workspace", mockReader)
	if err != nil {
		t.Fatalf("Failed to register workspace: %v", err)
	}

	// Create a mock stream
	mockStream := &mockNodeServiceStreamServer{
		events: make(chan *pb.StreamEvent, 10),
		ctx:    context.Background(),
	}

	// Create a document subscription request
	subscriptionReq := &pb.SubscriptionRequest{
		SubscriptionId: "test-subscription-1",
		Msg: &pb.SubscriptionRequest_DocumentSubscription{
			DocumentSubscription: &pb.DocumentSubscription{
				WorkspaceName: "test-workspace",
			},
		},
	}

	// Handle the subscription in a goroutine
	go server.handleSubscriptionRequest(mockStream, subscriptionReq)

	// Wait for subscription confirmation using helper
	waitForSubscriptionConfirmation(t, mockStream, "test-subscription-1")

	// Wait a bit for the subscription to be fully set up
	time.Sleep(10 * time.Millisecond)

	// Verify subscription is active
	server.subsMu.RLock()
	activeSub, exists := server.subscriptions["test-subscription-1"]
	server.subsMu.RUnlock()

	if !exists {
		t.Fatal("Expected active subscription to exist")
	}

	if activeSub.subscriptionID != "test-subscription-1" {
		t.Errorf("Expected subscription ID 'test-subscription-1', got '%s'", activeSub.subscriptionID)
	}

	if activeSub.workspaceName != "test-workspace" {
		t.Errorf("Expected workspace name 'test-workspace', got '%s'", activeSub.workspaceName)
	}
}

func TestUnsubscribe(t *testing.T) {
	server := NewServer()
	mockReader := mock.NewReader(t)

	err := server.RegisterWorkspace("test-workspace", mockReader)
	if err != nil {
		t.Fatalf("Failed to register workspace: %v", err)
	}

	mockStream := &mockNodeServiceStreamServer{
		events: make(chan *pb.StreamEvent, 10),
		ctx:    context.Background(),
	}

	// First, create a subscription
	subscriptionReq := &pb.SubscriptionRequest{
		SubscriptionId: "test-subscription-1",
		Msg: &pb.SubscriptionRequest_DocumentSubscription{
			DocumentSubscription: &pb.DocumentSubscription{
				WorkspaceName: "test-workspace",
			},
		},
	}

	go server.handleSubscriptionRequest(mockStream, subscriptionReq)

	// Wait for confirmation using helper
	waitForSubscriptionConfirmation(t, mockStream, "test-subscription-1")

	// Wait for subscription to be set up
	time.Sleep(50 * time.Millisecond)

	// Now unsubscribe
	unsubscribeReq := &pb.SubscriptionRequest{
		SubscriptionId: "test-subscription-1",
		Msg: &pb.SubscriptionRequest_Unsubscribe{
			Unsubscribe: &pb.Unsubscribe{},
		},
	}

	server.handleSubscriptionRequest(mockStream, unsubscribeReq)

	// Give it a moment to process
	time.Sleep(10 * time.Millisecond)

	// Test functionality: try to subscribe with the same ID again
	// If unsubscribe worked, this should succeed without conflict
	go server.handleSubscriptionRequest(mockStream, subscriptionReq)

	// If unsubscribe worked properly, we should get a confirmation, not an error
	waitForSubscriptionConfirmation(t, mockStream, "test-subscription-1")
}

func TestSubscriptionIDConflict(t *testing.T) {
	server := NewServer()
	mockReader := mock.NewReader(t)

	err := server.RegisterWorkspace("test-workspace", mockReader)
	if err != nil {
		t.Fatalf("Failed to register workspace: %v", err)
	}

	mockStream := &mockNodeServiceStreamServer{
		events: make(chan *pb.StreamEvent, 10),
		ctx:    context.Background(),
	}

	// Create first subscription
	subscriptionReq1 := &pb.SubscriptionRequest{
		SubscriptionId: "duplicate-id",
		Msg: &pb.SubscriptionRequest_DocumentSubscription{
			DocumentSubscription: &pb.DocumentSubscription{
				WorkspaceName: "test-workspace",
			},
		},
	}

	go server.handleSubscriptionRequest(mockStream, subscriptionReq1)

	// Wait for confirmation
	confirmEvent := <-mockStream.events
	if confirmEvent.GetSubscriptionConfirmation() == nil {
		t.Fatalf("Expected subscription confirmation, got %T", confirmEvent.GetEvent())
	}
	time.Sleep(50 * time.Millisecond)

	// Try to create second subscription with same ID
	go server.handleSubscriptionRequest(mockStream, subscriptionReq1)

	// Wait for error event (might need to drain other events first)
	var errorReceived bool
	for i := 0; i < 3; i++ { // Try up to 3 events
		select {
		case event := <-mockStream.events:
			if errorEvent := event.GetError(); errorEvent != nil {
				if errorEvent.RequestId != "duplicate-id" {
					t.Errorf("Expected request ID 'duplicate-id', got '%s'", errorEvent.RequestId)
				}
				if errorEvent.ErrorCode != pb.ErrorCode_SUBSCRIPTION_ID_CONFLICT {
					t.Errorf("Expected error code SUBSCRIPTION_ID_CONFLICT, got %v", errorEvent.ErrorCode)
				}
				if errorEvent.ErrorMessage != "subscription ID already exists" {
					t.Errorf("Expected error message 'subscription ID already exists', got '%s'", errorEvent.ErrorMessage)
				}
				errorReceived = true
				break
			}
			// Continue to next event if this wasn't an error
		case <-time.After(500 * time.Millisecond):
			break // Try next iteration
		}
		if errorReceived {
			break
		}
	}

	if !errorReceived {
		t.Fatal("Did not receive expected error event")
	}
}

func TestWorkspaceNotFound(t *testing.T) {
	server := NewServer()

	mockStream := &mockNodeServiceStreamServer{
		events: make(chan *pb.StreamEvent, 10),
		ctx:    context.Background(),
	}

	// Try to subscribe to non-existent workspace
	subscriptionReq := &pb.SubscriptionRequest{
		SubscriptionId: "test-subscription",
		Msg: &pb.SubscriptionRequest_DocumentSubscription{
			DocumentSubscription: &pb.DocumentSubscription{
				WorkspaceName: "non-existent-workspace",
			},
		},
	}

	go server.handleSubscriptionRequest(mockStream, subscriptionReq)

	// Wait for error event using helper
	waitForErrorEvent(t, mockStream, "test-subscription", pb.ErrorCode_WORKSPACE_NOT_FOUND, "workspace not found")
}

// Test helper functions

// waitForSubscriptionConfirmation waits for and validates a subscription confirmation event
func waitForSubscriptionConfirmation(t *testing.T, mockStream *mockNodeServiceStreamServer, expectedSubscriptionID string) {
	t.Helper()
	
	// We might get either confirmation first or initialization complete first
	// Try to get the confirmation within the first two events
	for i := 0; i < 2; i++ {
		select {
		case event := <-mockStream.events:
			if confirmation := event.GetSubscriptionConfirmation(); confirmation != nil {
				if confirmation.SubscriptionId != expectedSubscriptionID {
					t.Errorf("Expected subscription ID '%s', got '%s'", expectedSubscriptionID, confirmation.SubscriptionId)
				}
				return // Found confirmation, test passed
			}
			// If it's not a confirmation, continue to next event
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for subscription confirmation")
		}
	}
	
	t.Fatal("Did not receive subscription confirmation within expected events")
}

// waitForErrorEvent waits for and validates an error event
func waitForErrorEvent(t *testing.T, mockStream *mockNodeServiceStreamServer, expectedRequestID string, expectedErrorCode pb.ErrorCode, expectedErrorMessage string) {
	t.Helper()
	select {
	case event := <-mockStream.events:
		if errorEvent := event.GetError(); errorEvent != nil {
			if errorEvent.RequestId != expectedRequestID {
				t.Errorf("Expected request ID '%s', got '%s'", expectedRequestID, errorEvent.RequestId)
			}
			if errorEvent.ErrorCode != expectedErrorCode {
				t.Errorf("Expected error code %v, got %v", expectedErrorCode, errorEvent.ErrorCode)
			}
			if errorEvent.ErrorMessage != expectedErrorMessage {
				t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorEvent.ErrorMessage)
			}
		} else {
			t.Errorf("Expected error event, got %T", event.GetEvent())
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for error event")
	}
}
