/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package main

import (
	"fmt"
	"os"

	"github.com/PivotLLM/MCPLaunchPad/mcpserver"
	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
	"github.com/PivotLLM/MCPLaunchPad/mlogger"
)

// SimpleProvider provides a simple time tool
type SimpleProvider struct{}

// Ensure SimpleProvider implements ToolProvider
var _ mcptypes.ToolProvider = (*SimpleProvider)(nil)

// RegisterTools returns the list of tools this provider offers
func (s *SimpleProvider) RegisterTools() []mcptypes.ToolDefinition {
	return []mcptypes.ToolDefinition{
		{
			Name:        "get_greeting",
			Description: "Get a greeting message",
			Parameters: []*mcptypes.Parameter{
				mcptypes.StringParam("name", "Name to greet", false),
			},
			Handler: s.GetGreeting,
			Hints:   mcptypes.NewHints().ReadOnly(true),
		},
	}
}

// GetGreeting returns a greeting message
func (s *SimpleProvider) GetGreeting(options map[string]any) (string, error) {
	name := "World"
	if n, ok := options["name"].(string); ok && n != "" {
		name = n
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}

func main() {
	// Create logger
	logger, err := mlogger.New(
		mlogger.WithPrefix("BasicMCP"),
		mlogger.WithLogFile("basic-mcp.log"),
		mlogger.WithLogStdout(true),
		mlogger.WithDebug(true),
	)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Create provider
	provider := &SimpleProvider{}

	// Create MCP server with stdio transport
	srv, err := mcpserver.New(
		mcpserver.WithTransportStdio(),
		mcpserver.WithLogger(logger),
		mcpserver.WithName("BasicMCP"),
		mcpserver.WithVersion("1.0.0"),
		mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
		mcpserver.WithDefaultReadOnlyHint(false),
	)
	if err != nil {
		logger.Fatalf("Failed to create MCP server: %v", err)
	}

	// Start server (blocks until EOF in stdio mode)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
