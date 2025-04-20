// Copyright (c) 2025 Tenebris Technologies Inc.
// This software is licensed under the MIT License (see LICENSE for details).

package mcpserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PivotLLM/MCPLaunchPad/global"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Option defines a function type for configuring the MCPServer.
type Option func(*MCPServer)

// MCPServer represents the server instance.
type MCPServer struct {
	listen        string
	srv           *server.MCPServer
	sseServer     *server.SSEServer
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	logger        global.Logger
	debug         bool
	name          string
	version       string
	toolProviders []global.ToolProvider
}

func WithListen(listen string) Option {
	return func(m *MCPServer) {
		m.listen = listen
	}
}

func WithLogger(logger global.Logger) Option {
	return func(m *MCPServer) {
		m.logger = logger
	}
}

func WithDebug(debug bool) Option {
	return func(m *MCPServer) {
		m.debug = debug
	}
}

func WithName(name string) Option {
	return func(m *MCPServer) {
		m.name = name
	}
}

func WithVersion(version string) Option {
	return func(m *MCPServer) {
		m.version = version
	}
}

func WithToolProviders(providers []global.ToolProvider) Option {
	return func(s *MCPServer) {
		s.toolProviders = providers
	}
}

// New creates a new MCPServer instance with the provided options.
func New(options ...Option) (*MCPServer, error) {

	// Create a new MCPServer instance with default values
	// This is a wrapper around the mcp-go server
	m := &MCPServer{
		listen:    "localhost:8080",
		srv:       nil,
		sseServer: nil,
		ctx:       nil,
		cancel:    nil,
		logger:    nil,
		debug:     false,
		name:      "Generic-MCP",
		version:   "0.0.1",
		wg:        sync.WaitGroup{},
	}

	// Apply options
	for _, opt := range options {
		opt(m)
	}

	// If there is no logger, create one
	if m.logger == nil {
		return nil, fmt.Errorf("logger not set")
	}

	// Create an MCP server using the mcp-go library
	m.srv = server.NewMCPServer(
		m.name,
		m.version,
		server.WithLogging(),
		server.WithRecovery(),
		WithRequestLogging(m.logger), // Our custom request logging middleware
	)

	// Tools are in a separate file for better organization
	m.AddTools()

	// Return the MCPServer instance
	return m, nil
}

// Start runs the MCP server in a background goroutine and checks for a logger.
func (m *MCPServer) Start() error {
	if m.logger == nil {
		return fmt.Errorf("logger not set")
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		// Log the start
		m.logger.Infof("MCP server listening on TCP port %s", m.listen)

		// Create SSE server
		m.sseServer = server.NewSSEServer(m.srv)

		// Start the SSE server
		err := m.sseServer.Start(m.listen)
		// We don't need to log anything here - if the server is shutting down,
		// this is expected behavior and not an error condition
		_ = err
		return
	}()
	return nil
}

// Stop signals the MCP server to shut down and waits for the goroutine to exit.
func (m *MCPServer) Stop() error {
	// First cancel the context to signal all operations to stop
	if m.cancel != nil {
		m.cancel()
	}

	if m.sseServer != nil {
		// Attempt graceful shutdown with a timeout
		// Use a shorter timeout to avoid the context deadline exceeded error
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Shutdown the server and ignore all errors during shutdown
		// This prevents both the ErrServerClosed and context deadline exceeded errors
		_ = m.sseServer.Shutdown(ctx)
	}

	// Wait for the server goroutine to exit with a timeout
	waitCh := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(waitCh)
	}()

	// Wait for either the waitgroup to finish or a timeout
	select {
	case <-waitCh:
		// Goroutine completed successfully
		return nil
	case <-time.After(1 * time.Second):
		// If we're still waiting after 1 second, continue anyway
		// This prevents the context deadline exceeded error
		return nil
	}
}

func WithRequestLogging(logger global.Logger) server.ServerOption {
	return server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Log the request details
			logger.Debugf("Received request: %+v\n", request)

			// Call the next handler in the chain
			return next(ctx, request)
		}
	})
}
