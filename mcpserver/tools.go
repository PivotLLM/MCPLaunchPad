// Copyright (c) 2025 Tenebris Technologies Inc.
// This software is licensed under the MIT License (see LICENSE for details).

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (m *MCPServer) AddTools() {

	// Iterate over tool providers and register their tools
	for _, provider := range m.toolProviders {

		// Call the Register function of the provider to get tool definitions
		toolDefinitions := provider.RegisterTools()

		// Iterate over the tool definitions and register each tool
		for _, toolDef := range toolDefinitions {

			// Combine description and parameters into a slice of options
			toolOptions := []mcp.ToolOption{
				mcp.WithDescription(toolDef.Description),
			}
			for _, param := range toolDef.Parameters {
				toolOptions = append(toolOptions, mcp.WithString(param.Name, mcp.Description(param.Description)))
			}

			// Create the tool with all options
			tool := mcp.NewTool(toolDef.Name, toolOptions...)

			// Register the tool with the MCP server, creating a handler compatible with the MCP server
			// that wraps the tool's handler function with the provided options
			m.srv.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

				// Copy the MCP arguments to a map
				options := req.GetArguments()

				// Execute the tool's handler, passing the options
				result, err := toolDef.Handler(options)
				if err != nil {
					return mcp.NewToolResultError(err.Error()), err
				}
				return mcp.NewToolResultText(result), nil
			})
		}
	}
}
