/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// bearerTokenHTTPMiddleware wraps an HTTP handler with bearer token authentication
type bearerTokenHTTPMiddleware struct {
	handler   http.Handler
	validator mcptypes.BearerTokenValidator
	logger    mcptypes.Logger
}

// newBearerTokenHTTPMiddleware creates a new bearer token HTTP middleware
func newBearerTokenHTTPMiddleware(handler http.Handler, validator mcptypes.BearerTokenValidator, logger mcptypes.Logger) http.Handler {
	return &bearerTokenHTTPMiddleware{
		handler:   handler,
		validator: validator,
		logger:    logger,
	}
}

// ServeHTTP implements http.Handler
func (m *bearerTokenHTTPMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		m.logger.Warning("Missing Authorization header")
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	// Check Bearer prefix
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		m.logger.Warning("Invalid Authorization format")
		http.Error(w, "Invalid Authorization format - expected Bearer token", http.StatusUnauthorized)
		return
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	// Validate token
	contextData, err := m.validator(token)
	if err != nil {
		m.logger.Warningf("Bearer token validation failed: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Store auth context in request context for downstream handlers
	ctx := r.Context()
	for key, value := range contextData {
		ctx = context.WithValue(ctx, key, value)
	}

	// Create new request with enriched context and pass to handler
	r = r.WithContext(ctx)
	m.handler.ServeHTTP(w, r)
}

// authenticatedHTTPServer wraps an HTTP server with authentication
type authenticatedHTTPServer struct {
	server    *http.Server
	validator mcptypes.BearerTokenValidator
	logger    mcptypes.Logger
}

// Start starts the authenticated HTTP server
func (a *authenticatedHTTPServer) Start(handler http.Handler) error {
	// Wrap handler with authentication middleware
	authHandler := newBearerTokenHTTPMiddleware(handler, a.validator, a.logger)

	a.server = &http.Server{
		Handler: authHandler,
	}

	return a.server.ListenAndServe()
}

// Shutdown shuts down the server
func (a *authenticatedHTTPServer) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}

// wrapHTTPHandlerWithAuth wraps an http.Handler with bearer token authentication
func wrapHTTPHandlerWithAuth(handler http.Handler, validator mcptypes.BearerTokenValidator, logger mcptypes.Logger) http.Handler {
	return newBearerTokenHTTPMiddleware(handler, validator, logger)
}

// startHTTPServerWithHandler starts an HTTP server with the provided handler
func startHTTPServerWithHandler(listen string, handler http.Handler) error {
	srv := &http.Server{
		Addr:    listen,
		Handler: handler,
	}
	return srv.ListenAndServe()
}
