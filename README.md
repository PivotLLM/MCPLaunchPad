# MCPLaunchPad

**OATH2 functionality has not been fully tested**

Go library for building MCP (Model Context Protocol) servers with a clean, extensible architecture.

## Features

- **Multiple Transport Modes**: stdio, SSE, or HTTP (one mode per server instance)
- **Full MCP Hint System**: Complete support for ReadOnly, Destructive, Idempotent, and OpenWorld hints
- **Rich Parameter Types**: Full JSON Schema validation with easy-to-use helper constructors
- **Provider Pattern**: Clean separation of concerns via interfaces
- **Three-Level Hint Configuration**: Package defaults → server-wide config → tool-level overrides
- **Flexible Logging**: Optional logger with silent no-op fallback
- **Zero Dependencies on Project Code**: Reusable across your projects

## Quick Start

### Installation

```bash
go get github.com/PivotLLM/MCPLaunchPad
```

### Basic Stdio Server

```go
package main

import (
    "github.com/PivotLLM/MCPLaunchPad/mcpserver"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
    "github.com/PivotLLM/MCPLaunchPad/mlogger"
)

// Implement ToolProvider interface
type MyProvider struct{}

func (m *MyProvider) RegisterTools() []mcptypes.ToolDefinition {
    return []mcptypes.ToolDefinition{
        {
            Name:        "greet",
            Description: "Say hello",
            Parameters: []*mcptypes.Parameter{
                mcptypes.StringParam("name", "Name to greet", true),
            },
            Handler: func(opts map[string]any) (string, error) {
                name := opts["name"].(string)
                return "Hello, " + name + "!", nil
            },
            Hints: mcptypes.NewHints().ReadOnly(true),
        },
    }
}

func main() {
    logger, _ := mlogger.New(mlogger.WithLogStdout(true))
    defer logger.Close()

    srv, _ := mcpserver.New(
        mcpserver.WithTransportStdio(),
        mcpserver.WithLogger(logger),
        mcpserver.WithName("MyMCP"),
        mcpserver.WithVersion("1.0.0"),
        mcpserver.WithToolProviders([]mcptypes.ToolProvider{&MyProvider{}}),
    )

    srv.Start() // Blocks until EOF in stdio mode
}
```

## Architecture

### Core Packages

- **mcptypes/** - Shared interfaces and types (Logger, ToolProvider, ResourceProvider, PromptProvider, Parameter, ToolHints)
- **mcpserver/** - MCP server implementation with transport abstraction
- **mlogger/** - Simple file-based logger implementing mcptypes.Logger

### Provider Interfaces

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

## Transport Modes

### Stdio (for Claude Desktop, etc.)

```go
mcpserver.WithTransportStdio()
```

Communicates via stdin/stdout. `Start()` blocks until EOF. Use for Claude Desktop integration.

### SSE (Server-Sent Events)

```go
mcpserver.WithTransportSSE("localhost:8080")
```

HTTP server with streaming. `Start()` runs in background, `Stop()` for graceful shutdown.

### HTTP (Non-Streaming)

```go
mcpserver.WithTransportHTTP("localhost:8080")
```

Plain HTTP server. `Start()` runs in background, `Stop()` for graceful shutdown.

## Hint System

MCP supports four tool hints:

- **ReadOnlyHint**: Tool only reads data
- **DestructiveHint**: Tool may modify/delete data
- **IdempotentHint**: Repeated calls have no additional effect
- **OpenWorldHint**: Tool interacts with external entities

### Three-Level Configuration

**Level 1: Package Defaults**
```go
// Hardcoded in mcpserver package
ReadOnly: false, Destructive: false, Idempotent: false, OpenWorld: false
```

**Level 2: Server-Wide Overrides**
```go
mcpserver.New(
    mcpserver.WithDefaultReadOnlyHint(true),
    mcpserver.WithDefaultDestructiveHint(false),
    // ... affects all tools unless overridden
)
```

**Level 3: Tool-Level Overrides**
```go
// Method chaining
hints := mcptypes.NewHints().ReadOnly(true).OpenWorld(false)

// Or variadic constructor
hints := mcptypes.NewHints(
    mcptypes.ReadOnly(true),
    mcptypes.OpenWorld(false),
)

// In tool definition
ToolDefinition{
    Hints: hints,  // Overrides server-wide and package defaults
}
```

## Parameter Helpers

### Simple Types

```go
params := []*mcptypes.Parameter{
    mcptypes.StringParam("name", "User name", true),
    mcptypes.NumberParam("age", "User age", false),
    mcptypes.IntegerParam("count", "Item count", true),
    mcptypes.BoolParam("active", "Is active", false),
}
```

### With Validation

```go
mcptypes.StringParam("email", "Email address", true).
    WithPattern(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

mcptypes.NumberParam("port", "Port number", true).
    WithMinimum(1).
    WithMaximum(65535)

mcptypes.StringParam("status", "Status", true).
    WithEnum("pending", "active", "closed")
```

### Complex Types

```go
mcptypes.ArrayParam("tags", "Tags", false,
    mcptypes.StringParam("tag", "Tag", false))

mcptypes.ObjectParam("config", "Configuration", true,
    map[string]*mcptypes.Parameter{
        "host": mcptypes.StringParam("host", "Hostname", true),
        "port": mcptypes.IntegerParam("port", "Port", true),
    })
```

## Examples

See the `examples/` directory:

- **examples/bearer/** - HTTP server with bearer token authentication

## Documentation

- [mcpserver Package README](mcpserver/README.md) - Detailed API documentation
- [REFACTOR_DESIGN.md](REFACTOR_DESIGN.md) - Architecture and design decisions
- [CLAUDE.md](CLAUDE.md) - Development guide for AI assistants

## Prerequisites

- Go 1.24 or later

## Testing

```bash
# Build example
cd examples/bearer
go build

# Test without authentication
./bearer &
probe -transport http -url http://localhost:8080/mcp -list-only

# Test with bearer token authentication
./bearer --token "mytoken123" &
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer mytoken123" -list-only
```

## Security Note

Authentication support:
- **Bearer Token**: Built-in support via `WithBearerTokenAuth()` option
- **OAuth2**: Implement using `OAuth2Provider` interface (see [AUTHENTICATION.md](AUTHENTICATION.md))
- **Stdio mode**: Relies on OS-level process isolation

Do not expose HTTP/SSE servers to untrusted networks without proper authentication.

## Copyright and License

Copyright (c) 2025 by Tenebris Technologies Inc.
Licensed under the MIT License. See LICENSE for details.

## No Warranty

THIS SOFTWARE IS PROVIDED "AS IS," WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT. IN NO EVENT SHALL THE COPYRIGHT HOLDERS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Made in Canada 🍁
