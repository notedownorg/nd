---
title: "Initial Streaming API Implementation"
author: ["Liam White"]
date: "2025-06-17"
status: "Draft"
prd: "proposals/prds/2025-06-initial-streaming-api.md"
---

# Initial Streaming API Implementation

## Problem Statement

Applications working with markdown workspaces need to track and respond to changes in real-time. For example a document search will need to keep it's embeddings up to date as documents are edited in real-time or in the future a task manager might want to update tasks as they are completed in the source documents.

## Requirements

### Functional Requirements

**FR1**: Document-level subscriptions
- Subscribe to documents within a workspace
- Receive real-time updates when subscribed documents change
- Support for frontmatter as part of the document message

### Non-Functional Requirements

**NFR1**: Performance
- Document subscriptions perform efficiently for real-time updates
- Support reasonable number of concurrent document subscriptions per workspace

**NFR2**: Multiple Readers and Writers
- Support multiple concurrent subscribers to the same documents
- Handle updates from external clients (e.g. Obsidian, file system changes)

**NFR3**: Reliability
- All subscribers should be eventually consistent to the current workspace state within a second

## Scope

### In Scope

- Protobuf schema with complete document data
- Document subscription implementation
- Document content extraction and streaming including:
  - Document frontmatter metadata
- Real-time change notifications
- Core test suite for document streaming
- Basic documentation and examples

### Not in Scope

- Section-level subscriptions (future enhancement)
- Mutation/write capabilities
- Query interface for one-time requests
- Authentication and authorization
- Advanced filtering capabilities

## Proposal

### Approach: Document-Level Streaming API

**Strategy**: Implement a streaming API focused on document-level subscriptions:
1. Design protobuf schema with complete document data
2. Implement `DocumentSubscription` for workspace documents
3. Add document content extraction and streaming capabilities
4. Build server infrastructure to handle real-time document updates

**Pros**:
- Clean, focused streaming interface
- Efficient real-time document updates
- Foundation for future section-level enhancements
- Simpler initial implementation

**Cons**:
- Clients must process full documents for section-level changes
- Higher bandwidth usage for large documents

## Detailed Design

### Example Usage

**Subscribe to all documents in workspace:**
```go
subscription := &pb.DocumentSubscription{
    WorkspaceName: "project-workspace",
}
```

**gRPC service interface:**
```go
service NodeService {
    rpc Stream(stream StreamRequest) returns (stream StreamEvent);
}
```

### Protobuf Schema Design

```protobuf
// Document with complete content and metadata
message Document {
    string id = 1;  
    string workspace = 2;
    repeated string children = 3;  // Child node IDs (sections and placeholders)
    
    // Document metadata from frontmatter
    // Use arrays to maintain order
    message MetadataEntry {
        string key = 1;
        google.protobuf.Any value = 2;
    }
    repeated MetadataEntry metadata = 4;
}

message Unsubscribe { }

// Document-level subscription  
message DocumentSubscription {
    string workspace_name = 1;
}

// Subscription interface
message SubscriptionRequest {
    string subscription_id = 1;
    oneof msg {
        Unsubscribe unsubscribe = 10; 
        DocumentSubscription document_subscription = 11;
    }
}

// Streaming request and events
message StreamRequest {
    oneof request {
        SubscriptionRequest subscription = 1;
    }
}

message StreamEvent {
    oneof event {
        SubscriptionConfirmation subscription_confirmation = 1;
        InitializationComplete initialization_complete = 2;
        Load load = 3;
        Change change = 4;
        Delete delete = 5;
        Error error = 6;
    }
}

// Error types for streaming operations
enum ErrorCode {
    UNKNOWN_ERROR = 0;
    SUBSCRIPTION_ID_CONFLICT = 1;
    WORKSPACE_NOT_FOUND = 2;
    INVALID_REQUEST = 3;
    INTERNAL_ERROR = 4;
}

// Indicates an error with a request
message Error {
    string request_id = 1;  // Links to subscription_id or other request ID
    ErrorCode error_code = 2;
    string error_message = 3;
}

message SubscriptionConfirmation {
    string subscription_id = 1;
}

message InitializationComplete {
    string subscription_id = 1;
}

message Load {
    string subscription_id = 1;
    Node node = 2;
}

message Change {
    string subscription_id = 1;
    Node node = 2;
}

message Delete {
    string subscription_id = 1;
    Node node = 2;
}

message Node {
    oneof node {
        Document document = 1;
    }
}
```

### Server Architecture

#### Core Components

**1. Document Content Management**
- Implement document nodes with no content (for now).
- Extract and include document metadata from frontmatter
- Build document-to-protobuf conversion

**2. Subscription Management**
- Implement `handleDocumentSubscription` for document-level subscriptions
- Add subscription routing and event distribution
- Manage subscription lifecycle and cleanup

**3. Workspace Integration**
- Build workspace management to track document changes
- Implement event system for real-time document updates
- Add efficient file system monitoring and change detection

#### Implementation Approach

1. **Design protobuf schema** with complete document data including child nodes
2. **Implement document parsing** that extracts the frontmatter metadata
3. **Build subscription handlers** for document subscriptions
4. **Create server infrastructure** for streaming events and managing subscriptions
5. **Leverage existing goldmark parsing** to maintain compatibility with current markdown processing

## Test Plan

### Unit Tests

**Document Content Extraction**
- Test full document content extraction including frontmatter
- Verify metadata parsing from YAML frontmatter
- Test content extraction preserves markdown formatting
- Test edge cases (malformed frontmatter, empty documents)

**Document Subscription Logic**
- Test document subscription to workspace
- Verify document events are properly delivered to subscribers
- Test subscription cleanup and lifecycle management

**Protobuf Conversion**
- Test document-to-protobuf conversion with metadata
- Verify all document fields are properly populated
- Test conversion edge cases and error handling

### Integration Tests

**End-to-End Document Streaming**
- Test document subscription with real workspace
- Verify document events reach all subscribers
- Test subscription cleanup on client disconnect
- Test file system change detection and propagation

**Real-world Scenario**
- Create workspace with multiple markdown documents
- Subscribe to workspace and verify initial document load
- Modify documents externally (via file system)
- Verify all subscribers receive change notifications

### Performance Tests

**Document Subscription Performance**
- Benchmark document streaming with various document sizes
- Test multiple concurrent document subscriptions
- Measure memory usage and event propagation latency

## Implementation Plan

### Phase 1: Core Infrastructure

**Protobuf Schema and Server Foundation**
- [ ] Design protobuf schema with Document message and metadata support
- [ ] Add DocumentSubscription and related message types
- [ ] Implement basic gRPC server with Stream method
- [ ] Generate Go code and set up project structure

**Document Content Management**
- [ ] Implement document node with frontmatter metadata extraction
- [ ] Build document-to-protobuf conversion
- [ ] Add unit tests for content extraction and conversion
- [ ] Integrate with existing goldmark parsing infrastructure

### Phase 2: Streaming Implementation

**Document Subscriptions**
- [ ] Implement handleDocumentSubscription method
- [ ] Add workspace management and document tracking
- [ ] Implement file system monitoring for change detection
- [ ] Add subscription lifecycle management and cleanup

**Integration and Testing**
- [ ] End-to-end document streaming tests
- [ ] Test external file changes (file system simulation)
- [ ] Performance testing and optimization
- [ ] Documentation and usage examples
- [ ] Basic integration tests for document streaming

### Success Metrics

**Functionality**:
- [ ] Clients can subscribe to all documents in a workspace
- [ ] Document events include complete markdown content and metadata
- [ ] Real-time change detection from external file modifications
- [ ] Multiple concurrent subscribers supported
- [ ] Proper subscription cleanup and lifecycle management

**Performance**:
- [ ] Document subscriptions handle large workspaces efficiently
- [ ] Change notifications delivered within 1 second
- [ ] No significant memory leaks or performance degradation

**Quality**:
- [ ] Comprehensive test coverage for document streaming functionality
- [ ] Clean, extensible API design ready for future section-level enhancements
- [ ] Complete documentation and usage examples
