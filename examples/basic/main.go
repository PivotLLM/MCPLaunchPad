/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

// createBearerTokenValidator creates a simple bearer token validator
func createBearerTokenValidator(expectedToken string) mcptypes.BearerTokenValidator {
	return func(token string) (map[string]any, error) {
		if token != expectedToken {
			return nil, fmt.Errorf("invalid token")
		}
		// Return context data that will be available to handlers
		return map[string]any{
			"authenticated": true,
			"token":         token,
		}, nil
	}
}

func main() {
	// Parse command line flags
	token := flag.String("token", "", "Bearer token for authentication (if empty, no auth required)")
	listen := flag.String("listen", "localhost:8080", "Address to listen on for HTTP mode")
	flag.Parse()

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

	// Build server options
	opts := []mcpserver.Option{
		mcpserver.WithTransportHTTP(*listen),
		mcpserver.WithLogger(logger),
		mcpserver.WithName("BasicMCP"),
		mcpserver.WithVersion("1.0.0"),
		mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
		mcpserver.WithDefaultReadOnlyHint(false),
	}

	// Add bearer token authentication if token is provided
	if *token != "" {
		logger.Infof("Bearer token authentication enabled")
		opts = append(opts, mcpserver.WithBearerTokenAuth(createBearerTokenValidator(*token)))
	} else {
		logger.Info("No authentication configured")
	}

	// Create MCP server
	srv, err := mcpserver.New(opts...)
	if err != nil {
		logger.Fatalf("Failed to create MCP server: %v", err)
	}

	// Start server (runs in background for HTTP mode)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
	if err := srv.Stop(); err != nil {
		logger.Errorf("Error during shutdown: %v", err)
	}
}
