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
