# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development
- `make test` - Run all tests with 30s timeout
- `make hygiene` - Full hygiene check (tidy, generate, format, licenser, dirty check)
- `make generate` - Generate protobuf code and run go generate
- `make format` - Format Go code with gofmt

### Single Test Execution
- `nix develop --command go test -timeout=30s ./pkg/workspace/` - Test specific package
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

#### Filesystem Reader Test Structure

The filesystem reader tests follow a well-organized structure with several key components:

**Test Files:**
- `reader_test.go` - Core functionality tests (basic operations, initial load, file changes)
- `reader_perf_test.go` - Performance and stress tests with invariant verification  
- `testharness_test.go` - Reusable test infrastructure for invariant testing
- `subscriberview_test.go` - Test utilities for tracking subscriber state

**Test Patterns:**
- **Invariant Testing**: Tests verify that subscriber views match actual filesystem state
- **Eventual Consistency**: All stress tests include eventual consistency verification
- **Realistic Test Data**: Uses `pkg/test/words` for word-based filenames and nested directories
- **Helper Functions**: Common setup/teardown and assertion patterns abstracted into helpers
- **Constants**: Magic numbers extracted into named constants for maintainability

**Refactoring Opportunities Identified:**
1. **Code Duplication**: Test setup patterns could be further abstracted
2. **Magic Numbers**: Some timeout values could be made configurable constants
3. **Complex Functions**: Large test functions could benefit from breaking into smaller focused tests
4. **Missing Helpers**: Could add more specialized helper functions for common operations
5. **Pattern Inconsistencies**: Some tests use different assertion patterns that could be standardized

**Recent Improvements:**
- **Consolidated Test Infrastructure**: Moved common constants, helper functions, and patterns into centralized test harness
- **Enhanced Test Harness**: Added single-subscriber helpers, file operation utilities, and realistic data generators
- **Reduced Code Duplication**: Eliminated ~150 lines of repeated setup/teardown code across test files  
- **Standardized Event Handling**: All tests now use consistent event waiting and validation patterns
- **Improved Maintainability**: Test infrastructure changes now happen in one place (testharness_test.go)
- **Better Test Readability**: Individual tests focus on their logic rather than boilerplate setup
- **Shared Constants**: All timeout values and buffer sizes centralized for consistency

**Test Harness Features:**
- `newSingleSubscriberHarness()` - Simplified setup for single subscriber tests
- `harness.waitForEvent()`, `harness.assertEventMatch()` - Standardized event validation
- `harness.createTestFile()`, `harness.updateTestFile()`, `harness.deleteTestFile()` - File operation helpers
- `harness.generateRealisticFilename()`, `harness.generateMarkdownContent()` - Realistic test data generation
- `harness.expectNoEvent()`, `harness.drainEvents()` - Additional event utilities
- Shared constants for timeouts, buffer sizes, and operation durations

### Development Environment

Uses Nix flakes for reproducible development - all commands should be run via `nix develop --command` or within a nix develop shell.