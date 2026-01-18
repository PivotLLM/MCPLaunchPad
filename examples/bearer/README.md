# Bearer Token Authentication Example

This example demonstrates bearer token authentication for HTTP-based MCP servers.

## What It Does

Provides a single tool `get_greeting` that returns a personalized greeting message, protected by optional bearer token authentication.

## Features Demonstrated

- **HTTP Transport**: Communication over HTTP protocol
- **Bearer Token Authentication**: Optional token-based authentication
- **Tool Registration**: Simple tool with parameter
- **Hint System**: Marked as read-only tool
- **Parameter Helpers**: Using `StringParam()` for easy parameter definition
- **Logger Integration**: Using mlogger package
- **Graceful Shutdown**: Signal handling for clean server shutdown

## Building

```bash
go build
```

## Running

### Without Authentication

```bash
# Start server without authentication
./bearer

# In another terminal, test with probe
probe -transport http -url http://localhost:8080/mcp -list-only
probe -transport http -url http://localhost:8080/mcp -call get_greeting -params '{"name":"World"}'
```

### With Bearer Token Authentication

```bash
# Start server with token authentication
./bearer --token "mytoken123"

# Test without token (should fail)
probe -transport http -url http://localhost:8080/mcp -list-only
# Error: 401 Authorization required

# Test with correct token (should succeed)
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer mytoken123" -list-only
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer mytoken123" -call get_greeting -params '{"name":"World"}'
```

### Custom Listen Address

```bash
# Listen on a different address/port
./bearer --listen "localhost:9000"
./bearer --listen "localhost:9000" --token "mytoken123"
```

## Command Line Flags

- `--token` - Bearer token for authentication (if empty, no auth required)
- `--listen` - Address to listen on (default: "localhost:8080")

## Code Overview

The example demonstrates:

1. **Provider Implementation**: `SimpleProvider` implements `mcptypes.ToolProvider`
2. **Tool Definition**: Single tool with name, description, parameters, handler, and hints
3. **Bearer Token Validator**: Simple function-based token validation
4. **Server Creation**: Using `mcpserver.New()` with HTTP transport
5. **Conditional Authentication**: Bearer token auth enabled only if token is provided
6. **Background Execution**: HTTP server runs in background with graceful shutdown

## Authentication Flow

When bearer token is configured:

1. Client sends request with `Authorization: Bearer <token>` header
2. Server middleware extracts and validates the token
3. If valid, request proceeds to MCP handler
4. If invalid or missing, returns 401 Unauthorized

## Implementation Notes

- Token validation is intentionally simple for demonstration
- Production systems should use secure token storage (database, Redis, etc.)
- Always use HTTPS in production
- Consider adding rate limiting and logging for security monitoring
