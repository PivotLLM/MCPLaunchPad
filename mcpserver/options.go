/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import "github.com/PivotLLM/MCPLaunchPad/mcptypes"

// Option defines a function type for configuring the MCPServer.
type Option func(*MCPServer)

// TransportMode represents the transport mode for the MCP server
type TransportMode int

const (
	// TransportStdio uses stdin/stdout for communication
	TransportStdio TransportMode = iota
	// TransportSSE uses Server-Sent Events over HTTP
	TransportSSE
	// TransportHTTP uses plain HTTP (non-streaming)
	TransportHTTP
)

// Transport selection options (exactly one required)

// WithTransportStdio configures the server to use stdio transport
func WithTransportStdio() Option {
	return func(m *MCPServer) {
		m.transportMode = TransportStdio
		m.transportConfigured = true
	}
}

// WithTransportSSE configures the server to use SSE transport
func WithTransportSSE(listen string) Option {
	return func(m *MCPServer) {
		m.transportMode = TransportSSE
		m.listen = listen
		m.transportConfigured = true
	}
}

// WithTransportHTTP configures the server to use HTTP transport
func WithTransportHTTP(listen string) Option {
	return func(m *MCPServer) {
		m.transportMode = TransportHTTP
		m.listen = listen
		m.transportConfigured = true
	}
}

// Basic configuration options

// WithLogger sets the logger for the server
func WithLogger(logger mcptypes.Logger) Option {
	return func(m *MCPServer) {
		m.logger = logger
	}
}

// WithDebug enables debug mode
func WithDebug(debug bool) Option {
	return func(m *MCPServer) {
		m.debug = debug
	}
}

// WithName sets the server name
func WithName(name string) Option {
	return func(m *MCPServer) {
		m.name = name
	}
}

// WithVersion sets the server version
func WithVersion(version string) Option {
	return func(m *MCPServer) {
		m.version = version
	}
}

// Provider registration options

// WithToolProviders sets the tool providers
func WithToolProviders(providers []mcptypes.ToolProvider) Option {
	return func(s *MCPServer) {
		s.toolProviders = providers
	}
}

// WithResourceProviders sets the resource providers
func WithResourceProviders(providers []mcptypes.ResourceProvider) Option {
	return func(s *MCPServer) {
		s.resourceProviders = providers
	}
}

// WithPromptProviders sets the prompt providers
func WithPromptProviders(providers []mcptypes.PromptProvider) Option {
	return func(s *MCPServer) {
		s.promptProviders = providers
	}
}

// Authentication options

// WithBearerTokenAuth enables bearer token authentication
func WithBearerTokenAuth(validator mcptypes.BearerTokenValidator) Option {
	return func(m *MCPServer) {
		m.bearerTokenValidator = validator
	}
}

// Hint default configuration options

// WithDefaultReadOnlyHint sets the default ReadOnlyHint for all tools
func WithDefaultReadOnlyHint(value bool) Option {
	return func(m *MCPServer) {
		m.defaultReadOnlyHint = &value
	}
}

// WithDefaultDestructiveHint sets the default DestructiveHint for all tools
func WithDefaultDestructiveHint(value bool) Option {
	return func(m *MCPServer) {
		m.defaultDestructiveHint = &value
	}
}

// WithDefaultIdempotentHint sets the default IdempotentHint for all tools
func WithDefaultIdempotentHint(value bool) Option {
	return func(m *MCPServer) {
		m.defaultIdempotentHint = &value
	}
}

// WithDefaultOpenWorldHint sets the default OpenWorldHint for all tools
func WithDefaultOpenWorldHint(value bool) Option {
	return func(m *MCPServer) {
		m.defaultOpenWorldHint = &value
	}
}

// Deprecated option (for backwards compatibility)

// WithNoStreaming is deprecated. Use WithTransportHTTP instead.
func WithNoStreaming(noStreaming bool) Option {
	return func(m *MCPServer) {
		if noStreaming {
			m.transportMode = TransportHTTP
			if m.listen == "" {
				m.listen = "localhost:8080"
			}
			m.transportConfigured = true
		}
	}
}

// WithListen is deprecated. Use WithTransportSSE or WithTransportHTTP instead.
func WithListen(listen string) Option {
	return func(m *MCPServer) {
		m.listen = listen
	}
}
