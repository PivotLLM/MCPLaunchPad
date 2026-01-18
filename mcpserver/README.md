# mcpserver Package

A production-ready Go library for building MCP (Model Context Protocol) servers.

## Features

- **Multiple Transport Modes**: stdio, SSE, or HTTP (one per server instance)
- **Full Hint System**: Four-level MCP hint support (ReadOnly, Destructive, Idempotent, OpenWorld)
- **Rich Parameters**: Full JSON Schema validation with helper constructors
- **Provider Pattern**: Clean separation via ToolProvider, ResourceProvider, PromptProvider interfaces
- **Flexible Logging**: Optional logger with no-op fallback
- **Three-Level Hint Configuration**: Package defaults, server-wide config, tool-level overrides

## Quick Start

### Stdio Server

```go
package main

import (
    "github.com/PivotLLM/MCPLaunchPad/mcpserver"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
    "github.com/PivotLLM/MCPLaunchPad/mlogger"
)

func main() {
    logger, _ := mlogger.New(mlogger.WithLogStdout(true))
    defer logger.Close()

    provider := &MyProvider{} // Implements mcptypes.ToolProvider

    srv, _ := mcpserver.New(
        mcpserver.WithTransportStdio(),
        mcpserver.WithLogger(logger),
        mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
    )

    srv.Start() // Blocks until EOF
}
```

### SSE Server

```go
srv, _ := mcpserver.New(
    mcpserver.WithTransportSSE("localhost:8080"),
    mcpserver.WithLogger(logger),
    mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
)

srv.Start() // Runs in background
// ... do other work ...
srv.Stop()  // Graceful shutdown
```

## Provider Interface

Implement one or more provider interfaces:

```go
type ToolProvider interface {
    RegisterTools() []ToolDefinition
}

type ResourceProvider interface {
    RegisterResources() []ResourceDefinition
    RegisterResourceTemplates() []ResourceTemplateDefinition
}

type PromptProvider interface {
    RegisterPrompts() []PromptDefinition
}
```

## Parameter Helpers

```go
// Simple parameters
params := []*mcptypes.Parameter{
    mcptypes.StringParam("name", "User name", true),
    mcptypes.NumberParam("age", "User age", false),
    mcptypes.BoolParam("active", "Is active", false),
}

// With validation
mcptypes.StringParam("email", "Email address", true).
    WithPattern(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

mcptypes.NumberParam("port", "Port number", true).
    WithMinimum(1).
    WithMaximum(65535)
```

## Hint System

### Three Levels of Configuration

**Level 1: Package Defaults** (hardcoded)
- ReadOnly: false
- Destructive: false
- Idempotent: false
- OpenWorld: false

**Level 2: Server-Wide Config** (via options)
```go
mcpserver.New(
    mcpserver.WithDefaultReadOnlyHint(true),
    mcpserver.WithDefaultDestructiveHint(false),
    // ...
)
```

**Level 3: Tool-Level Overrides** (in ToolDefinition)
```go
// Method chaining
hints := mcptypes.NewHints().ReadOnly(true).Destructive(false)

// Variadic constructor
hints := mcptypes.NewHints(
    mcptypes.ReadOnly(true),
    mcptypes.Destructive(false),
)

// In tool definition
ToolDefinition{
    Name: "my_tool",
    // ...
    Hints: hints,
}
```

## Configuration Options

### Transport (required - exactly one)
- `WithTransportStdio()` - stdin/stdout communication
- `WithTransportSSE(listen string)` - Server-Sent Events
- `WithTransportHTTP(listen string)` - Plain HTTP

### Basic
- `WithLogger(logger mcptypes.Logger)` - Optional, defaults to no-op
- `WithDebug(bool)` - Enable debug mode
- `WithName(string)` - Server name
- `WithVersion(string)` - Server version

### Providers
- `WithToolProviders([]mcptypes.ToolProvider)`
- `WithResourceProviders([]mcptypes.ResourceProvider)`
- `WithPromptProviders([]mcptypes.PromptProvider)`

### Hint Defaults
- `WithDefaultReadOnlyHint(bool)`
- `WithDefaultDestructiveHint(bool)`
- `WithDefaultIdempotentHint(bool)`
- `WithDefaultOpenWorldHint(bool)`

## Logger Interface

If no logger provided, uses silent no-op logger. Implement `mcptypes.Logger`:

```go
type Logger interface {
    Debug(string)
    Info(string)
    Notice(string)
    Warning(string)
    Error(string)
    Fatal(string)
    Debugf(string, ...any)
    Infof(string, ...any)
    Noticef(string, ...any)
    Warningf(string, ...any)
    Errorf(string, ...any)
    Fatalf(string, ...any)
    Close()
}
```

## Architecture

The package wraps `github.com/mark3labs/mcp-go` and provides:

1. **Transport abstraction** - Unified API for stdio/SSE/HTTP
2. **Provider registration** - Automatic conversion to mcp-go types
3. **Hint resolution** - Three-level configuration system
4. **Parameter conversion** - JSON Schema to mcp-go mapping
5. **Lifecycle management** - Start/Stop with graceful shutdown

## Examples

See the `examples/` directory for complete working examples:
- `examples/basic/` - Simple stdio server
- (More examples to be added)
