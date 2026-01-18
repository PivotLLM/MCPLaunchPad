/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// withRequestLogging is a middleware function that logs request details.
func withRequestLogging(logger mcptypes.Logger) server.ServerOption {
	return server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Log the request details
			logger.Debugf("Request: %+v", request)

			// Call the next handler in the chain
			return next(ctx, request)
		}
	})
}

// withBearerTokenAuth is a middleware function that validates bearer tokens.
func withBearerTokenAuth(validator mcptypes.BearerTokenValidator, logger mcptypes.Logger) server.ServerOption {
	return server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Call validator
			contextData, err := validator("")
			if err != nil {
				logger.Warningf("Bearer token validation failed: %v", err)
				return mcp.NewToolResultError("Authentication required"), err
			}

			// Add auth context data to context
			for key, value := range contextData {
				ctx = context.WithValue(ctx, key, value)
			}

			// Call next handler with enriched context
			return next(ctx, request)
		}
	})
}
