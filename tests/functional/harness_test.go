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

// Package functional contains end-to-end tests for the nd streaming server.
// These tests run the actual server binary as an external process and test
// the complete gRPC API from a client perspective.
package functional

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/test/words"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ServerStartupTimeout       = 10 * time.Second
	ServerPort                 = "9080"
	ServerAddr                 = "localhost:" + ServerPort
	EventualConsistencyTimeout = 5 * time.Second
	EventPropagationDelay      = 500 * time.Millisecond
)

// testHarness manages the server process and gRPC client for functional tests
type testHarness struct {
	t         *testing.T
	tmpDir    string
	serverCmd *exec.Cmd
	client    pb.NodeServiceClient
	conn      *grpc.ClientConn
}

// newTestHarness creates a new functional test environment
func newTestHarness(t *testing.T) *testHarness {
	tmpDir, err := os.MkdirTemp("", "nd-functional-test")
	require.NoError(t, err)

	h := &testHarness{
		t:      t,
		tmpDir: tmpDir,
	}

	// Start server process
	h.startServer()

	// Connect gRPC client
	h.connectClient()

	return h
}

// startServer launches the nd server binary as an external process
func (h *testHarness) startServer() {
	// Build path to server binary
	binaryPath := filepath.Join("..", "..", "bin", "nd")

	// Start server with test configuration
	h.serverCmd = exec.Command(binaryPath,
		"-port", ServerPort,
		"-workspace", h.tmpDir,
		"-name", "test-workspace",
		"-log", "debug")

	// Capture server output for debugging
	h.serverCmd.Stdout = os.Stdout
	h.serverCmd.Stderr = os.Stderr

	err := h.serverCmd.Start()
	require.NoError(h.t, err, "Failed to start server process")

	h.t.Logf("Started server process with PID %d", h.serverCmd.Process.Pid)

	// Wait for server to be ready
	h.waitForServerReady()
}

// waitForServerReady polls the server until it's accepting connections
func (h *testHarness) waitForServerReady() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerStartupTimeout)
	defer cancel()

	// Initial delay to let server start
	time.Sleep(200 * time.Millisecond)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.t.Fatal("Server failed to start within timeout")
		case <-ticker.C:
			// Try to connect
			conn, err := grpc.NewClient(ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err == nil {
				conn.Close()
				h.t.Log("Server is ready")
				return
			}
		}
	}
}

// connectClient establishes a gRPC connection to the server
func (h *testHarness) connectClient() {
	var err error
	h.conn, err = grpc.NewClient(ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to server")

	h.client = pb.NewNodeServiceClient(h.conn)
	h.t.Log("Connected gRPC client")
}

// cleanup stops the server and cleans up resources
func (h *testHarness) cleanup() {
	if h.conn != nil {
		h.conn.Close()
	}

	if h.serverCmd != nil && h.serverCmd.Process != nil {
		h.t.Logf("Stopping server process %d", h.serverCmd.Process.Pid)

		// Send SIGTERM for graceful shutdown
		err := h.serverCmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			h.t.Logf("Failed to send SIGTERM: %v", err)
		}

		// Wait for graceful shutdown with timeout
		done := make(chan error, 1)
		go func() {
			done <- h.serverCmd.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			h.t.Log("Server didn't stop gracefully, force killing")
			h.serverCmd.Process.Kill()
			h.serverCmd.Wait()
		case err := <-done:
			if err != nil {
				h.t.Logf("Server process exited with error: %v", err)
			} else {
				h.t.Log("Server process stopped gracefully")
			}
		}
	}

	os.RemoveAll(h.tmpDir)
}

// createTestFile creates a markdown file in the test workspace
func (h *testHarness) createTestFile(relativePath string, content string) error {
	fullPath := filepath.Join(h.tmpDir, relativePath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// updateTestFile updates an existing file
func (h *testHarness) updateTestFile(relativePath string, content string) error {
	fullPath := filepath.Join(h.tmpDir, relativePath)
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// deleteTestFile removes a file
func (h *testHarness) deleteTestFile(relativePath string) error {
	fullPath := filepath.Join(h.tmpDir, relativePath)
	return os.Remove(fullPath)
}

// generateRealisticContent creates realistic markdown content
func (h *testHarness) generateRealisticContent(title string) string {
	return fmt.Sprintf("# %s\n\n%s %s content created at %s",
		title, words.Random(), words.Random(), time.Now().Format("15:04:05"))
}

// streamWithTimeout creates a stream with automatic timeout
func (h *testHarness) streamWithTimeout(timeout time.Duration) (pb.NodeService_StreamClient, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	stream, err := h.client.Stream(ctx)
	require.NoError(h.t, err, "Failed to create stream")
	return stream, cancel
}

// expectStreamEvent waits for a specific event type with timeout
func (h *testHarness) expectStreamEvent(stream pb.NodeService_StreamClient, expectedEventType string, timeout time.Duration) *pb.StreamEvent {
	done := make(chan *pb.StreamEvent, 1)
	errorChan := make(chan error, 1)

	go func() {
		event, err := stream.Recv()
		if err != nil {
			errorChan <- err
			return
		}
		done <- event
	}()

	select {
	case event := <-done:
		switch expectedEventType {
		case "subscription_confirmation":
			assert.NotNil(h.t, event.GetSubscriptionConfirmation(), "Expected subscription confirmation event")
		case "initialization_complete":
			assert.NotNil(h.t, event.GetInitializationComplete(), "Expected initialization complete event")
		case "load":
			assert.NotNil(h.t, event.GetLoad(), "Expected load event")
		case "change":
			assert.NotNil(h.t, event.GetChange(), "Expected change event")
		case "delete":
			assert.NotNil(h.t, event.GetDelete(), "Expected delete event")
		case "error":
			assert.NotNil(h.t, event.GetError(), "Expected error event")
		}
		return event
	case err := <-errorChan:
		h.t.Fatalf("Stream error while waiting for %s event: %v", expectedEventType, err)
	case <-time.After(timeout):
		h.t.Fatalf("Timeout waiting for %s event", expectedEventType)
	}
	return nil
}

// drainStreamEvents collects all pending events from the stream
func (h *testHarness) drainStreamEvents(stream pb.NodeService_StreamClient, maxEvents int, timeout time.Duration) []*pb.StreamEvent {
	events := make([]*pb.StreamEvent, 0, maxEvents)
	deadline := time.Now().Add(timeout)

	for len(events) < maxEvents && time.Now().Before(deadline) {
		select {
		default:
			// Try to receive an event with a short timeout
			done := make(chan *pb.StreamEvent, 1)
			errorChan := make(chan error, 1)

			go func() {
				event, err := stream.Recv()
				if err != nil {
					errorChan <- err
					return
				}
				done <- event
			}()

			select {
			case event := <-done:
				events = append(events, event)
			case <-errorChan:
				// Stream closed or error, stop collecting
				return events
			case <-time.After(100 * time.Millisecond):
				// No more events immediately available
				return events
			}
		}
	}

	return events
}
