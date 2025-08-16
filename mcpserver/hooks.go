/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

//goland:noinspection GoUnusedParameter
func (m *MCPServer) hookAfterListPrompts(ctx context.Context, id any, request *mcp.ListPromptsRequest, result *mcp.ListPromptsResult) {
	if m.debug {
		m.logger.Debugf("%s: %v", request.Request.Method, result.Prompts)
	} else {
		m.logger.Infof("%s: %s items returned", request.Request.Method, len(result.Prompts))
	}
}

//goland:noinspection GoUnusedParameter
func (m *MCPServer) hookAfterListResources(ctx context.Context, id any, request *mcp.ListResourcesRequest, result *mcp.ListResourcesResult) {
	if m.debug {
		m.logger.Debugf("%s: %v", request.Request.Method, result.Resources)
	} else {
		m.logger.Infof("%s: %s items returned", request.Request.Method, len(result.Resources))
	}
}

//goland:noinspection GoUnusedParameter
func (m *MCPServer) hookAfterListResourceTemplates(ctx context.Context, id any, request *mcp.ListResourceTemplatesRequest, result *mcp.ListResourceTemplatesResult) {
	if m.debug {
		m.logger.Debugf("%s: %v", request.Request.Method, result.ResourceTemplates)
	} else {
		m.logger.Infof("%s: %s items returned", request.Request.Method, len(result.ResourceTemplates))
	}
}

//goland:noinspection GoUnusedParameter
func (m *MCPServer) hookAfterListTools(ctx context.Context, id any, request *mcp.ListToolsRequest, result *mcp.ListToolsResult) {
	if m.debug {
		m.logger.Debugf("%s: %v", request.Request.Method, result.Tools)
	} else {
		m.logger.Infof("%s: %s items returned", request.Request.Method, len(result.Tools))
	}
}
