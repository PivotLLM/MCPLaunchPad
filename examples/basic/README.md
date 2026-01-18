# Basic MCP Server Example

This is a minimal stdio-based MCP server demonstrating the core concepts.

## What It Does

Provides a single tool `get_greeting` that returns a personalized greeting message.

## Features Demonstrated

- **Stdio Transport**: Communication via stdin/stdout
- **Tool Registration**: Simple tool with parameter
- **Hint System**: Marked as read-only tool
- **Parameter Helpers**: Using `StringParam()` for easy parameter definition
- **Logger Integration**: Using mlogger package

## Building

```bash
go build
```

## Running

### With MCP Probe (Testing)

```bash
# List available tools
probe -stdio ./basic -list-only

# Call the greeting tool
probe -stdio ./basic -call get_greeting -params '{"name":"World"}'
```

### With Claude Desktop (Production)

Add to your Claude Desktop config:

```json
{
  "mcpServers": {
    "basic": {
      "command": "/path/to/MCPLaunchPad/examples/basic/basic"
    }
  }
}
```

## Code Overview

The example demonstrates:

1. **Provider Implementation**: `SimpleProvider` implements `mcptypes.ToolProvider`
2. **Tool Definition**: Single tool with name, description, parameters, handler, and hints
3. **Server Creation**: Using `mcpserver.New()` with stdio transport
4. **Blocking Execution**: `srv.Start()` blocks until stdin closes in stdio mode
