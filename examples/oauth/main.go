/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PivotLLM/MCPLaunchPad/mcpserver"
	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
	"github.com/PivotLLM/MCPLaunchPad/mlogger"
	"github.com/PivotLLM/MCPLaunchPad/oauth2"
)

// SimpleProvider provides a simple greeting tool
type SimpleProvider struct{}

// Ensure SimpleProvider implements ToolProvider
var _ mcptypes.ToolProvider = (*SimpleProvider)(nil)

// RegisterTools returns the list of tools this provider offers
func (s *SimpleProvider) RegisterTools() []mcptypes.ToolDefinition {
	return []mcptypes.ToolDefinition{
		{
			Name:        "get_greeting",
			Description: "Get a personalized greeting message",
			Parameters: []*mcptypes.Parameter{
				mcptypes.StringParam("name", "Name to greet", false),
			},
			Handler: s.GetGreeting,
			Hints:   mcptypes.NewHints().ReadOnly(true),
		},
		{
			Name:        "get_user_info",
			Description: "Get authenticated user information",
			Parameters:  []*mcptypes.Parameter{},
			Handler:     s.GetUserInfo,
			Hints:       mcptypes.NewHints().ReadOnly(true),
		},
	}
}

// GetGreeting returns a greeting message
func (s *SimpleProvider) GetGreeting(options map[string]any) (string, error) {
	name := "World"
	if n, ok := options["name"].(string); ok && n != "" {
		name = n
	}
	return fmt.Sprintf("Hello, %s! You are authenticated via OAuth2.", name), nil
}

// GetUserInfo returns authenticated user information
func (s *SimpleProvider) GetUserInfo(options map[string]any) (string, error) {
	// In a real implementation, you would extract user info from the request context
	// For this example, we just return a placeholder message
	return "User info would be available from the OAuth2 token context", nil
}

func main() {
	// Parse command line flags
	listen := flag.String("listen", "localhost:8080", "Address to listen on for HTTP mode")
	skipAuth := flag.Bool("skip-auth", false, "Skip OAuth2 authentication (for testing)")
	flag.Parse()

	// Create logger
	logger, err := mlogger.New(
		mlogger.WithPrefix("OAuth2MCP"),
		mlogger.WithLogFile("oauth-mcp.log"),
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
		mcpserver.WithName("OAuth2MCP"),
		mcpserver.WithVersion("1.0.0"),
		mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
		mcpserver.WithDefaultReadOnlyHint(false),
	}

	// Configure OAuth2 authentication unless skipped
	if !*skipAuth {
		// Get Google OAuth2 credentials from environment
		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

		if clientID == "" || clientSecret == "" {
			logger.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set")
		}

		logger.Info("Initializing Google OAuth2 authentication...")

		// Create Google OAuth2 provider
		oauth2Provider := oauth2.NewGoogleProvider(
			clientID,
			clientSecret,
			[]string{"email", "profile"},
		)

		// Perform device flow
		logger.Info("Starting OAuth2 device flow...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		tokenResp, deviceResp, err := oauth2Provider.DeviceFlowWithPolling(ctx, 5*time.Second)
		if err != nil {
			logger.Fatalf("OAuth2 device flow failed: %v", err)
		}

		// Display device code to user
		fmt.Println("\n=== OAuth2 Device Flow ===")
		fmt.Printf("1. Open this URL in your browser: %s\n", deviceResp.VerificationURI)
		fmt.Printf("2. Enter this code: %s\n", deviceResp.UserCode)
		fmt.Printf("3. Waiting for authorization...\n\n")

		logger.Infof("Authentication successful! Access token received (expires in %d seconds)", tokenResp.ExpiresIn)

		// Create bearer token validator from OAuth2 provider
		validator := oauth2Provider.CreateBearerTokenValidator()
		opts = append(opts, mcpserver.WithBearerTokenAuth(validator))

		logger.Info("OAuth2 bearer token authentication enabled")
		fmt.Println("OAuth2 authentication configured successfully!")
		fmt.Printf("\nTo use this server, include the access token in the Authorization header:\n")
		fmt.Printf("  Authorization: Bearer %s\n\n", tokenResp.AccessToken[:20]+"...")
	} else {
		logger.Info("Authentication skipped (--skip-auth flag set)")
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
