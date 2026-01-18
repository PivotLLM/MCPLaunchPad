# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MCPLaunchPad is a production-ready Go library for building MCP (Model Context Protocol) servers. It provides a clean, extensible architecture with support for multiple transport modes, full hint system, and rich parameter validation.

**Release Status**: This code is unreleased and under active development. Backwards compatibility is not required for refactoring or API changes.

## Build and Run Commands

```bash
# Build library packages
go build ./mcptypes
go build ./mcpserver
go build ./mlogger

# Build and test basic example
cd examples/basic
go build
probe -stdio ./basic -list-only
probe -stdio ./basic -call get_greeting -params '{"name":"World"}'
cd ../..

# Run all tests
go test ./...

# Format code
go fmt ./...

# Verify builds
go build ./...
```

## Development Commands

```bash
# Install dependencies
go mod tidy

# Check for issues (requires golangci-lint)
golangci-lint run

# Compile without running
go build ./...
```

## Architecture Overview

### Core Philosophy

MCPLaunchPad is a **library-first project**. The `mcpserver` package is the primary product; examples demonstrate usage.

### Package Structure

**mcptypes/** - Shared MCP type definitions (no dependencies)
- `logger.go` - Logger interface
- `providers.go` - ToolProvider, ResourceProvider, PromptProvider interfaces
- `parameters.go` - Parameter struct with JSON Schema support + helper constructors
- `hints.go` - ToolHints struct with builder patterns
- `auth.go` - Future authentication interfaces (stubs)

**mcpserver/** - MCP server implementation
- `mcpserver.go` - Main server with transport selection logic
- `options.go` - Configuration options (With* functions)
- `noop_logger.go` - No-op logger fallback
- `hints.go` - Three-level hint resolution
- `tools.go` - Tool registration with parameter conversion
- `resources.go` - Resource registration
- `prompts.go` - Prompt registration
- `hooks.go` - MCP hooks for logging
- `middleware.go` - Request logging middleware

**mlogger/** - Logger implementation
- `mlogger.go` - File-based logger implementing mcptypes.Logger

**examples/** - Example implementations
- `basic/` - Minimal stdio server

**global/** - DEPRECATED - To be removed

### Key Design Patterns

**Provider Pattern**: Services implement interfaces (ToolProvider, ResourceProvider, PromptProvider) to register their capabilities

**Functional Options Pattern**: Configuration via `WithXxx()` functions passed to constructors

**Three-Level Hint Configuration**:
1. Package defaults (hardcoded in mcpserver/tools.go)
2. Server-wide config (via WithDefaultXxxHint options)
3. Tool-level overrides (via ToolDefinition.Hints)

**Transport Abstraction**: Single transport mode selected at initialization (stdio, SSE, or HTTP)

## Transport Modes

### Stdio Mode
- Uses stdin/stdout for communication
- `Start()` blocks until EOF
- For Claude Desktop and similar MCP clients
- Example: `mcpserver.WithTransportStdio()`

### SSE Mode (Server-Sent Events)
- HTTP server with streaming
- `Start()` runs in background goroutine
- `Stop()` for graceful shutdown
- Example: `mcpserver.WithTransportSSE("localhost:8080")`

### HTTP Mode (Non-Streaming)
- Plain HTTP request/response
- `Start()` runs in background goroutine
- `Stop()` for graceful shutdown
- Example: `mcpserver.WithTransportHTTP("localhost:8080")`

## Hint System

### Four MCP Hints
1. **ReadOnlyHint** - Tool is read-only
2. **DestructiveHint** - Tool may perform destructive updates
3. **IdempotentHint** - Repeated calls with same args have no additional effect
4. **OpenWorldHint** - Tool interacts with external entities

### Resolution Logic (mcpserver/tools.go:resolveHints)
For each hint:
1. Check tool-level `ToolDefinition.Hints.{Hint}` - if non-nil, use it
2. Else check server-wide config (m.default{Hint}) - if non-nil, use it
3. Else use package default constant (default{Hint})

### Builder Patterns
```go
// Method chaining
hints := mcptypes.NewHints().ReadOnly(true).Destructive(false)

// Variadic constructor
hints := mcptypes.NewHints(
    mcptypes.ReadOnly(true),
    mcptypes.Destructive(false),
)
```

## Parameter System

### Full JSON Schema Support
`mcptypes.Parameter` struct supports:
- Types: string, number, integer, boolean, array, object
- String validation: Pattern, MinLength, MaxLength, Format, Enum
- Numeric validation: Minimum, Maximum, ExclusiveMin/Max, MultipleOf
- Array validation: Items schema, MinItems, MaxItems, UniqueItems
- Object validation: Properties, AdditionalProperties
- Default values

### Helper Constructors (mcptypes/parameters.go)
```go
// Simple types
mcptypes.StringParam(name, description, required)
mcptypes.NumberParam(name, description, required)
mcptypes.IntegerParam(name, description, required)
mcptypes.BoolParam(name, description, required)
mcptypes.ArrayParam(name, description, required, itemType)
mcptypes.ObjectParam(name, description, required, properties)

// Fluent API
.WithPattern(pattern)
.WithMinLength(min)
.WithMaxLength(max)
.WithMinimum(min)
.WithMaximum(max)
.WithEnum(values...)
.WithDefault(value)
```

### Conversion to mcp-go (mcpserver/tools.go:parameterToToolOption)
Maps `mcptypes.Parameter` to `mcp.ToolOption` based on Type field:
- "string" → `mcp.WithString()`
- "number", "integer" → `mcp.WithNumber()`
- "boolean" → `mcp.WithBoolean()`
- "array" → `mcp.WithArray()`
- "object" → `mcp.WithObject()`

Note: Some validation options (Format, Minimum, Maximum) may not be supported by current mcp-go version and are commented out.

## Adding New Providers

1. Create a package/struct that implements one or more provider interfaces
2. Use functional options pattern for configuration
3. Implement `Register*()` methods returning definition slices
4. In your main.go, create provider instance
5. Pass to `mcpserver.New()` via `WithToolProviders()`, `WithResourceProviders()`, or `WithPromptProviders()`

## Configuration

### Required Options
Exactly one transport option:
- `WithTransportStdio()`
- `WithTransportSSE(listen string)`
- `WithTransportHTTP(listen string)`

### Optional Options
- `WithLogger(mcptypes.Logger)` - Defaults to no-op logger
- `WithDebug(bool)`
- `WithName(string)`
- `WithVersion(string)`
- `WithToolProviders([]mcptypes.ToolProvider)`
- `WithResourceProviders([]mcptypes.ResourceProvider)`
- `WithPromptProviders([]mcptypes.PromptProvider)`
- `WithDefaultReadOnlyHint(bool)`
- `WithDefaultDestructiveHint(bool)`
- `WithDefaultIdempotentHint(bool)`
- `WithDefaultOpenWorldHint(bool)`

## Testing

### Unit Tests
Currently no unit tests (library is new). Should add:
- Hint resolution logic tests
- Parameter conversion tests
- Transport mode selection tests

### Integration Tests
Manual testing with probe tool:
```bash
cd examples/basic
go build
probe -stdio ./basic -list-only
probe -stdio ./basic -call tool_name -params '{}'
```

## Code Style

- Use functional options pattern for configuration
- Implement provider interfaces for clean separation
- Use helper constructors for common parameter types
- Follow Go naming conventions
- Add godoc comments to all exported symbols

## Security Notes

- Server binds to `localhost` only in examples
- No authentication mechanism exists for HTTP/SSE modes
- Stdio mode relies on OS-level process isolation
- Do not expose HTTP/SSE servers to untrusted networks without authentication

## Future Enhancements (Designed but Not Implemented)

**Authentication** (mcptypes/auth.go has stub interfaces):
- Bearer token validation via function injection
- OAuth2 device flow via provider interface
- Middleware-based implementation

**Context-Aware Handlers**:
- Optional `ContextAwareToolHandler` alongside regular `ToolHandler`
- Pass tenant info, auth context via Go context

## Common Tasks

**Add a new tool**:
1. Add method to your provider's `RegisterTools()` return slice
2. Define name, description, parameters, handler, and optional hints
3. Rebuild and test

**Change hint defaults**:
1. Modify package constants in `mcpserver/tools.go` (Level 1)
2. Or use `WithDefaultXxxHint()` options (Level 2)
3. Or set in `ToolDefinition.Hints` (Level 3)

**Add validation to parameter**:
1. Use fluent API methods like `.WithPattern()`, `.WithMinimum()`, etc.
2. See `mcptypes/parameters.go` for full list

**Change transport mode**:
1. Replace `WithTransportXxx()` option in server creation
2. Exactly one transport option must be specified

## Dependencies

- `github.com/mark3labs/mcp-go` v0.43.2 - Core MCP protocol implementation
- `github.com/joho/godotenv` v1.5.1 - Environment variable loading (examples only)

## File Organization

Keep the library clean:
- Core types in `mcptypes/`
- Server implementation in `mcpserver/`
- Examples in `examples/`
- Each example has its own README.md
