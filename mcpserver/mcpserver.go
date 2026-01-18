/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// MCPServer represents the server instance.
type MCPServer struct {
	// Transport configuration
	listen              string
	transportMode       TransportMode
	transportConfigured bool

	// mcp-go server and transports
	srv        *server.MCPServer
	sseServer  *server.SSEServer
	httpServer *server.StreamableHTTPServer

	// Lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Configuration
	logger  mcptypes.Logger
	debug   bool
	name    string
	version string

	// Providers
	toolProviders     []mcptypes.ToolProvider
	resourceProviders []mcptypes.ResourceProvider
	promptProviders   []mcptypes.PromptProvider

	// Authentication
	bearerTokenValidator mcptypes.BearerTokenValidator

	// Default hint values (Level 2 configuration)
	defaultReadOnlyHint    *bool
	defaultDestructiveHint *bool
	defaultIdempotentHint  *bool
	defaultOpenWorldHint   *bool
}

// New creates a new MCPServer instance with the provided options.
func New(options ...Option) (*MCPServer, error) {

	// Create a new MCPServer instance with default values
	m := &MCPServer{
		listen:              "localhost:8080",
		transportMode:       TransportSSE, // Default if not specified
		transportConfigured: false,
		srv:                 nil,
		sseServer:           nil,
		httpServer:          nil,
		ctx:                 nil,
		cancel:              nil,
		logger:              nil,
		debug:               false,
		name:                "Generic-MCP",
		version:             "0.0.1",
		wg:                  sync.WaitGroup{},
		// Hint defaults are nil (will use package defaults)
	}

	// Apply options
	for _, opt := range options {
		opt(m)
	}

	// Validate transport configuration
	if !m.transportConfigured {
		return nil, fmt.Errorf("no transport mode specified; use WithTransportStdio(), WithTransportSSE(), or WithTransportHTTP()")
	}

	// If there is no logger, use no-op logger
	if m.logger == nil {
		m.logger = &noopLogger{}
	}

	// Create hooks
	hooks := &server.Hooks{}
	hooks.AddAfterListPrompts(m.hookAfterListPrompts)
	hooks.AddAfterListResources(m.hookAfterListResources)
	hooks.AddAfterListResourceTemplates(m.hookAfterListResourceTemplates)
	hooks.AddAfterListTools(m.hookAfterListTools)

	// Create an MCP server using the mcp-go library
	m.srv = server.NewMCPServer(
		m.name,
		m.version,
		server.WithLogging(),
		server.WithRecovery(),
		server.WithHooks(hooks),
		withRequestLogging(m.logger), // Our custom request logging middleware
	)

	// Register tools, resources, and prompts
	m.AddTools()
	m.AddResources()
	m.AddResourceTemplates()
	m.AddPrompts()

	// Return the MCPServer instance
	return m, nil
}

// Start runs the MCP server.
// For stdio mode, this blocks until EOF or signal.
// For HTTP/SSE modes, this runs in a background goroutine.
func (m *MCPServer) Start() error {
	if m.logger == nil {
		return fmt.Errorf("logger not set")
	}

	switch m.transportMode {
	case TransportStdio:
		// Stdio mode blocks until EOF
		m.logger.Info("MCP server starting in stdio mode")
		err := server.ServeStdio(m.srv)
		if err != nil {
			m.logger.Errorf("Stdio server error: %v", err)
			return err
		}
		m.logger.Info("Stdio server closed")
		return nil

	case TransportSSE:
		// SSE mode runs in background
		m.ctx, m.cancel = context.WithCancel(context.Background())
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.logger.Infof("MCP server listening on %s (SSE mode)", m.listen)
			m.sseServer = server.NewSSEServer(m.srv)

			// Wrap with authentication if configured
			if m.bearerTokenValidator != nil {
				handler := wrapHTTPHandlerWithAuth(m.sseServer, m.bearerTokenValidator, m.logger)
				err := startHTTPServerWithHandler(m.listen, handler)
				_ = err
			} else {
				err := m.sseServer.Start(m.listen)
				_ = err
			}
		}()
		return nil

	case TransportHTTP:
		// HTTP mode runs in background
		m.ctx, m.cancel = context.WithCancel(context.Background())
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.logger.Infof("MCP server listening on %s (HTTP mode)", m.listen)
			m.httpServer = server.NewStreamableHTTPServer(m.srv)

			// Wrap with authentication if configured
			if m.bearerTokenValidator != nil {
				handler := wrapHTTPHandlerWithAuth(m.httpServer, m.bearerTokenValidator, m.logger)
				err := startHTTPServerWithHandler(m.listen, handler)
				_ = err
			} else {
				err := m.httpServer.Start(m.listen)
				_ = err
			}
		}()
		return nil

	default:
		return fmt.Errorf("unknown transport mode: %d", m.transportMode)
	}
}

// Stop signals the MCP server to shut down and waits for goroutines to exit.
// For stdio mode, this is a no-op (use signal handling).
func (m *MCPServer) Stop() error {
	// Stdio mode doesn't need explicit stop
	if m.transportMode == TransportStdio {
		return nil
	}

	// Cancel context to signal shutdown
	if m.cancel != nil {
		m.cancel()
	}

	// Shutdown the appropriate transport
	var transport interface {
		Shutdown(context.Context) error
	}

	switch m.transportMode {
	case TransportSSE:
		transport = m.sseServer
	case TransportHTTP:
		transport = m.httpServer
	}

	if transport != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = transport.Shutdown(ctx)
	}

	// Wait for server goroutine to exit with timeout
	waitCh := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		return nil
	case <-time.After(1 * time.Second):
		return nil
	}
}
