# Development Guide

This guide covers how to set up and work with the nd streaming-server codebase.

## Environment Setup

This project uses Nix flakes for reproducible development environments. All development commands should be run within the Nix environment.

### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled

### Setup

```bash
# Enter development environment
nix develop

# Or run individual commands
nix develop --command <command>
```

## Development Commands

### Testing

```bash
# Run all tests with 30s timeout
make test

# Test specific package
nix develop --command go test -timeout=30s ./pkg/workspace/

# Run specific test function
nix develop --command go test -run TestSpecificFunction ./pkg/workspace/

# Run tests with race detection (already enabled by default)
nix develop --command go test -race ./...
```

### Code Quality

```bash
# Full hygiene check (tidy, generate, format, licenser, dirty check)
make hygiene

# Format Go code
make format

# Check for formatting issues
gofmt -l .
```

### Code Generation

```bash
# Generate all code (protobuf + go generate)
make generate

# Generate only protobuf code
nix develop --command buf generate

# Lint protobuf files
nix develop --command buf lint
```

## Project Structure

```
├── api/                    # Generated protobuf code
│   ├── go/                # Go bindings
│   └── proto/             # Protobuf definitions
├── pkg/                   # Core packages
│   ├── fsnotify/         # Filesystem watching
│   ├── goldmark/         # Markdown processing extensions
│   ├── reorderbuffer/    # Event ordering utilities
│   ├── server/           # gRPC server implementation
│   ├── test/             # Test utilities
│   └── workspace/        # Core workspace logic
└── flake.nix             # Nix development environment
```

## Testing Architecture

### Test Organization

The project uses comprehensive testing with several patterns:

- **Unit Tests**: Individual component testing with mocks
- **Integration Tests**: Real filesystem operations
- **Stress Tests**: Randomized event sequences (10,000+ operations)
- **Race Detection**: Concurrent operation validation

### Key Test Files

- `*_test.go` - Standard unit and integration tests
- `*_perf_test.go` - Performance and stress testing
- `testharness_test.go` - Shared test infrastructure
- `mock/` directories - Mock implementations for testing

### Test Helpers

The project includes centralized test infrastructure in `testharness_test.go`:

```go
// Create test harness for single subscriber
harness := newSingleSubscriberHarness(t, tempDir)

// Wait for specific events
harness.waitForEvent(t, "expected event")

// File operations with automatic cleanup
harness.createTestFile("test.md", "content")
harness.updateTestFile("test.md", "new content")
harness.deleteTestFile("test.md")
```

### Running Specific Tests

```bash
# Filesystem reader tests
nix develop --command go test ./pkg/workspace/reader/filesystem/

# Performance tests
nix develop --command go test -run Perf ./pkg/workspace/reader/filesystem/

# Workspace core tests
nix develop --command go test ./pkg/workspace/
```

## Protobuf Development

### File Structure

```
api/proto/nodes/v1alpha1/
├── stream.proto    # Main streaming API
└── types.proto     # Shared types
```

### Workflow

1. Edit `.proto` files in `api/proto/`
2. Run `buf lint` to check syntax
3. Run `buf generate` to update Go bindings
4. Update Go code to use new protobuf types

### Buffer Configuration

- `buf.yaml` - Protobuf linting and breaking change detection
- `buf.gen.yaml` - Code generation configuration

## Architecture Guidelines

### Event Handling

- **Buffering**: Each subscriber gets its own event buffer
- **Ordering**: Events are reordered to handle filesystem race conditions
- **Graceful Degradation**: Slow clients don't block fast ones

### Error Handling

- Use structured logging for debugging
- Graceful handling of filesystem errors
- Client disconnection handling in streaming API

### Concurrency

- All workspace operations are thread-safe
- Use channels for event communication
- Prefer composition over inheritance for testability

## Debugging

### Logging

The project uses structured logging. Key areas to monitor:

- Workspace initialization events
- File system event processing
- gRPC connection lifecycle
- Buffer overflow conditions

### Common Issues

1. **Race Conditions**: Use `go test -race` to detect
2. **Event Ordering**: Check reorder buffer configuration
3. **Memory Leaks**: Monitor subscriber cleanup
4. **Filesystem Events**: Verify fsnotify integration

## Performance Considerations

### Optimization Areas

- Event buffering strategies
- Markdown parsing efficiency  
- Concurrent subscriber handling
- Memory usage in large workspaces

### Benchmarking

```bash
# Run performance tests
nix develop --command go test -bench=. ./pkg/workspace/reader/filesystem/

# Profile memory usage
nix develop --command go test -memprofile=mem.prof ./pkg/workspace/

# Profile CPU usage  
nix develop --command go test -cpuprofile=cpu.prof ./pkg/workspace/
```

## Contributing

1. Run `make hygiene` before committing
2. Ensure all tests pass with `make test`
3. Follow existing code patterns and naming conventions
4. Add tests for new functionality
5. Update documentation for API changes