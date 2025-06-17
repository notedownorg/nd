# nd streaming-server

A high-performance gRPC server for real-time Markdown workspace graphs using a streaming-first API design.

## Overview

This server provides real-time synchronization of Markdown document workspaces through a bidirectional gRPC streaming API. It processes Markdown files into hierarchical node structures and delivers live updates to connected clients as files change on disk.

## Key Features

- **Real-time Updates**: Bidirectional gRPC streaming with pub/sub pattern for instant workspace synchronization
- **Hierarchical Processing**: Documents organized into sections with type-safe IDs (e.g., "Document|filename.md")  
- **Filesystem Integration**: Automatic file watching with intelligent event buffering and reordering
- **Frontmatter Support**: Full YAML metadata preservation in Markdown files using custom Goldmark extension
- **Thread-safe Caching**: Concurrent workspace access with graceful degradation for slow clients
- **Event Buffering**: Per-subscriber buffering prevents any single client from blocking others

## Architecture

### Workspace System
The core workspace (`pkg/workspace/`) acts as a central orchestrator that:
- Manages document graphs with thread-safe caching
- Uses pub/sub pattern for real-time updates across multiple subscribers  
- Processes markdown files into hierarchical nodes (Documents → Sections)
- Handles load-completion ordering to buffer events until workspace initialization completes

### Node Architecture
Type-safe hierarchical structure (`pkg/workspace/node/`) with:
- Documents containing Sections in a tree structure
- Bidirectional markdown serialization preserving frontmatter
- Kind-prefixed IDs for type safety and debugging

### Reader Interface
Abstraction layer (`pkg/workspace/reader/`) providing:
- Pluggable data sources (filesystem implementation included)
- Event-driven model (Load/Change/Delete) with subscription support
- Threadsafe subscribers using buffered channels

### gRPC Streaming API
Single bidirectional `Stream` RPC (`api/proto/nodes/v1alpha1/`) enabling:
- Subscription-based real-time workspace synchronization
- Client-server communication for workspace operations
- Efficient event delivery with buffering strategies

## Use Cases

- **Live Markdown Editors**: Real-time collaborative editing with instant updates
- **Documentation Systems**: Dynamic content management with file system integration  
- **Knowledge Graphs**: Hierarchical document relationships with live synchronization
- **Content Management**: Automated processing of markdown file changes

## API

The server exposes a single bidirectional `Stream` RPC defined in `api/proto/nodes/v1alpha1/stream.proto`. Clients can:
- Subscribe to workspace changes
- Receive real-time document and section updates  
- Send commands for workspace operations

