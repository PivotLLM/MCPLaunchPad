// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (m *MCPServer) AddResources() {

	// Iterate over prompt providers and register their prompts
	for _, provider := range m.resourceProviders {

		// Call the Register function of the provider to get tool definitions
		resourceDefinitions := provider.RegisterResources()

		// Iterate over the tool definitions and register each tool
		for _, resource := range resourceDefinitions {

			newResource := mcp.NewResource(
				resource.URI,
				resource.Name,
				mcp.WithResourceDescription(resource.Description),
				mcp.WithMIMEType(resource.MIMEType),
			)

			// Add resource with its handler
			m.srv.AddResource(newResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {

				// Copy the MCP arguments to a map
				options := make(map[string]any)
				for key, value := range request.Params.Arguments {
					options[key] = value
				}

				// Execute the tool's handler, passing the options
				resp, err := resource.Handler(request.Params.URI, options)
				if err != nil {
					return nil, err
				}

				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      resp.URI,
						MIMEType: resp.MIMEType,
						Text:     resp.Content}}, nil
			})
		}
	}
}
