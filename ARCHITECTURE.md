# Architecture Overview

This document describes the layered architecture of the nd server, a gRPC streaming server for real-time Markdown workspace graphs. The system is designed around a streaming-first API with event-driven architecture and careful attention to thread safety and performance.

## Layer Hierarchy

```
┌─────────────────────────────────────────┐
│                Server                   │ ← gRPC API & protobuf conversion
├─────────────────────────────────────────┤
│               Workspace                 │ ← Document graph management & pub/sub
├─────────────────────────────────────────┤
│                Reader                   │ ← Data source abstraction
├─────────────────────────────────────────┤
│               fsnotify                  │ ← Filesystem event monitoring
└─────────────────────────────────────────┘
                    ↑
               Filesystem
```

## Layer Responsibilities

### fsnotify Layer (`pkg/fsnotify/`)

**Purpose**: Low-level filesystem event monitoring with intelligent debouncing and deduplication.

**Key Components**:
- `RecursiveWatcher`: Monitors directory trees for file changes
- Event debouncing with configurable timers (50ms default, 150ms max wait)
- Automatic recursive directory watching with ignored patterns (`.git`, `node_modules`, etc.)

**Key Behaviors**:
- **Immediate CREATE events**: New files are reported instantly to prevent race conditions
- **Debounced WRITE events**: Multiple rapid writes are consolidated into single events
- **Race condition mitigation**: When adding new directories, scans for existing files to catch missed events
- **Timer management**: Dual-timer system (debounce + max wait) ensures events are never lost
- **Directory lifecycle**: Automatically adds/removes watchers as directories are created/deleted

**Event Types**: `Change` (create/modify), `Remove` (delete/rename)

### Reader Layer (`pkg/workspace/reader/`)

**Purpose**: Document-level abstraction that converts filesystem events into content-aware document events.

**Key Components**:
- `Reader` interface: Pluggable data source abstraction
- `filesystem.Reader`: Concrete implementation for local filesystem
- Per-subscriber event buffering with graceful degradation

**Key Behaviors**:
- **Document-level operations**: Maps individual files to document entities
- **Content loading**: Reads file contents into memory when changes occur
- **Read-only operations**: Cannot mutate documents, only observes and reports changes
- **Event translation**: Converts fsnotify path events to document events with full content
- **Initial load**: Walks directory tree to load existing markdown files with content
- **Buffering strategy**: Per-subscriber buffered channels prevent slow clients from blocking others
- **Markdown filtering**: Only processes `.md` files, treating each as a document
- **Path-to-ID mapping**: Converts absolute filesystem paths to relative document IDs

**Event Types**: `Load`, `Change`, `Delete`, `SubscriberLoadComplete`

### Workspace Layer (`pkg/workspace/`)

**Purpose**: Central orchestrator managing document graphs with thread-safe caching and pub/sub coordination.

**Key Components**:
- `Workspace`: Main coordinator with cached document storage
- `Cache[T]`: Thread-safe generic cache with automatic event emission
- Subscriber management with load-completion ordering
- Document parsing with goldmark and frontmatter support

**Key Behaviors**:
- **Thread-safe caching**: RWMutex-protected node storage with automatic event emission
- **Load-completion ordering**: Events are buffered during initial load, then flushed in order
- **Document processing**: Parses markdown into hierarchical node structures (Documents → Sections)
- **Frontmatter handling**: Extracts and streams YAML metadata from document headers
- **Deletion strategy**: Uses depth-first traversal to clean up orphaned nodes
- **Event buffering**: 1000-event buffers per subscriber with drop-on-full behavior

**Node Types**: Documents (with frontmatter metadata), Sections (heading-based), Placeholders (content blocks)

### Server Layer (`pkg/server/`)

**Purpose**: gRPC API endpoint with protobuf conversion and client session management.

**Key Components**:
- `Server`: gRPC service implementation with workspace registry
- Subscription management with per-client isolation
- Protobuf conversion for network serialization
- Error handling with typed error codes

**Key Behaviors**:
- **Bidirectional streaming**: Single `Stream` RPC handles all client communication
- **Subscription lifecycle**: Manages client subscriptions with automatic cleanup on disconnect
- **Workspace selection**: Supports multiple workspaces with default fallback
- **Metadata conversion**: Transforms YAML frontmatter to protobuf Any types
- **Event forwarding**: Converts workspace events to protobuf stream events
- **Error propagation**: Structured error responses with request correlation

**Protocol**: Single bidirectional stream with subscription-based event delivery

## Data Flow

```
File Change → fsnotify → Reader → Workspace → Server → gRPC Client
     ↓           ↓         ↓         ↓          ↓        ↓
  debounce → filter+load → parse → cache → convert → stream
```

**Detailed Flow**:
1. **fsnotify**: Detects filesystem changes, debounces rapid writes
2. **Reader**: Filters for markdown files, loads content into memory, emits document-level events
3. **Workspace**: Parses document content into hierarchical nodes, caches in thread-safe storage
4. **Server**: Converts workspace events to protobuf, manages client subscriptions
5. **gRPC Client**: Receives streaming document updates

## Key Design Patterns

### Event Buffering Strategy
Each layer implements buffering to prevent blocking:
- **fsnotify**: Timer-based debouncing to reduce event volume
- **Reader**: Per-subscriber buffering with overflow handling
- **Workspace**: Load-completion ordering with 1000-event buffers
- **Server**: Context-based cancellation for clean disconnection

### Thread Safety
- **Cache operations**: RWMutex protection with atomic event emission
- **Subscriber management**: Separate mutexes for subscription state
- **Event channels**: Buffered channels with non-blocking sends
- **Graceful degradation**: Events dropped rather than blocking on slow consumers

### Clear Separation of Concerns
- **Reader immutability**: Reader layer cannot mutate documents, ensuring predictable cache invalidation
- **Single source of truth**: Filesystem is the authoritative source, workspace cache reflects current state
- **Event-driven updates**: All changes flow through the event pipeline, no direct cache manipulation
- **Layer isolation**: Each layer has distinct responsibilities without overlap

### Error Handling
- **Layer isolation**: Errors propagated up through structured channels
- **Client feedback**: Typed error codes sent via gRPC stream
- **Panic recovery**: Subscriber goroutines recover from panics
- **Resource cleanup**: Automatic unsubscription on errors

## Performance Characteristics

### Scalability
- **O(1) node lookup**: Cache uses map-based storage
- **Bounded memory**: Event buffers have fixed sizes
- **Non-blocking writes**: Slow clients don't affect fast ones
- **Efficient parsing**: Single-pass markdown processing

### Latency
- **Immediate creates**: New files reported without debouncing
- **Sub-200ms writes**: Debounced events delivered within 150ms max
- **Streaming delivery**: Events sent as soon as processed
- **Load balancing**: Multiple subscribers process independently

This architecture enables real-time markdown workspace synchronization with eventual consistency guarantees.
