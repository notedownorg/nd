# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development
- `make test` - Run all tests with 10s timeout
- `make hygiene` - Full hygiene check (tidy, generate, format, licenser, dirty check)
- `make generate` - Generate protobuf code and run go generate
- `make format` - Format Go code with gofmt

### Single Test Execution
- `nix develop --command go test -timeout=10s ./pkg/workspace/` - Test specific package
- `nix develop --command go test -run TestSpecificFunction ./pkg/workspace/` - Run specific test

### Protobuf Management
- `nix develop --command buf generate` - Generate Go code from protobuf definitions
- `nix develop --command buf lint` - Lint protobuf files

## Architecture Overview

This is a gRPC server for real-time Markdown workspace graphs using a streaming-first API design.

### Core Components

**Workspace System** (`pkg/workspace/`):
- Central orchestrator managing document graphs with thread-safe caching
- Uses pub/sub pattern for real-time updates across multiple subscribers
- Processes markdown files into hierarchical nodes (Documents → Sections)

**Node Architecture** (`pkg/workspace/node/`):
- Type-safe IDs with kind prefixes (e.g., "Document|filename.md")
- Hierarchical structure: Documents contain Sections
- Bidirectional markdown serialization preserving frontmatter

**Reader Interface** (`pkg/workspace/reader/`):
- Abstraction layer for data sources with filesystem implementation
- Event-driven (Load/Change/Delete) with subscription model
- Threadsafe subscribers use buffered channels to prevent blocking

**gRPC Streaming API**:
- Single bidirectional `Stream` RPC in `api/proto/nodes/v1alpha1/`
- Subscription-based real-time workspace synchronization
- Event buffering prevents slow clients from blocking others

### Key Patterns

1. **Event Buffering Strategy**: Per-subscriber buffering with graceful degradation
2. **Load-Completion Ordering**: Events buffered until workspace initialization completes
3. **Filesystem Integration**: Custom fsnotify wrapper with reorder buffer for handling out-of-order events
4. **Markdown Processing**: Goldmark with custom frontmatter extension for YAML metadata

### Test Architecture

- Mock readers for isolated workspace testing (`pkg/workspace/reader/mock/`)
- Filesystem integration tests with real file operations
- Randomized stress testing (10,000+ event sequences)
- Race condition detection for concurrent operations

### Development Environment

Uses Nix flakes for reproducible development - all commands should be run via `nix develop --command` or within a nix develop shell.