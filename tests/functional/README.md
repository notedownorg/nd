# Functional Tests

This directory contains end-to-end functional tests for the nd streaming server.

## Overview

These tests run the actual `cmd/nd` server binary as an external process and test the complete gRPC API from a client perspective. The tests focus on:

- Document subscription functionality
- Event streaming and eventual consistency
- Concurrent client scenarios
- Error handling
- Nested directory support

## Architecture

The functional tests:

1. **Start Server Process**: Launch the `cmd/nd` binary with test configuration
2. **Connect gRPC Client**: Establish connection to the running server
3. **Test Scenarios**: Execute various streaming scenarios
4. **Verify Results**: Check for eventual consistency within 1 second
5. **Cleanup**: Gracefully stop server and clean up resources

## Test Structure

- `harness_test.go` - Test harness and client utilities
- `stream_test.go` - Main streaming functionality tests

## Key Features Tested

### Document Subscriptions
- Basic document subscription with confirmation
- Multiple operations (create, update, delete)
- Concurrent subscriptions from multiple clients
- **Frontmatter parsing** - YAML metadata extraction and streaming

### Eventual Consistency
- System convergence within 1 second timeout
- Handling of rapid file operations
- Verification that all clients see the same final state

### Directory Support
- Files in nested directory structures
- Deep path handling (up to 5 levels)
- Realistic filename generation using word combinations

## Running Tests

```bash
# Build and run functional tests (recommended)
make test-functional

# Or run all tests (unit + functional)
make test-all

# Manual approach (requires building binary first)
make build
go test ./tests/functional/...

# Note: These tests start actual server processes on port 9080
# Make sure the port is available before running
```

## Test Configuration

- **Server Port**: 9080 (to avoid conflicts with development servers)
- **Startup Timeout**: 10 seconds
- **Event Consistency Timeout**: 5 seconds  
- **Event Propagation Delay**: 500ms

## Implementation Notes

### Server Binary Path
Tests expect the server binary at `../../cmd/nd/main.go` and run it with `go run`.

### Workspace Management
Each test creates a temporary directory as the workspace and populates it with test files.

### Event Verification
Tests verify eventual consistency rather than strict event ordering, allowing for system debouncing and performance optimizations.

### Graceful Shutdown
Server processes are terminated with SIGTERM for graceful shutdown, with fallback to SIGKILL if needed.

## Features Tested

### Frontmatter Support
The functional tests verify complete frontmatter parsing functionality:
- **Basic metadata**: Strings, numbers, booleans in YAML frontmatter
- **Complex structures**: Arrays, nested objects (converted to JSON strings)
- **Type preservation**: Proper protobuf Any type conversion
- **Empty handling**: Documents without frontmatter return empty metadata

### Document Streaming
- Document ID verification with `Document|filename` format
- Event type validation (load, change, delete)
- Subscription management and error handling
- System stability under concurrent load

### Test Coverage
- Simple frontmatter (title, author, tags, priority, published)
- Complex frontmatter (nested objects, arrays, mixed types)
- Documents without frontmatter
- Real-time metadata updates during file changes
