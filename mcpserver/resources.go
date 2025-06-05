// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (m *MCPServer) AddResources() {

	// Iterate over resource providers and register their resources
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
			m.srv.AddResource(newResource,
				func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {

					// Copy the MCP arguments to a map
					options := request.Params.Arguments
					if options == nil {
						options = make(map[string]any)
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
				},
			)
		}
	}
}

func (m *MCPServer) AddResourceTemplates() {

	// Iterate over resource providers and register their templates
	for _, provider := range m.resourceProviders {

		// Call the Register function of the provider to get tool definitions
		resourceTemplates := provider.RegisterResourceTemplates()

		// Iterate over the tool definitions and register each tool
		for _, resourceTemplate := range resourceTemplates {

			template := mcp.NewResourceTemplate(
				resourceTemplate.URITemplate,
				resourceTemplate.Name,
				mcp.WithTemplateDescription(resourceTemplate.Description),
				mcp.WithTemplateMIMEType(resourceTemplate.MIMEType))

			// Add resource template with its handler
			m.srv.AddResourceTemplate(template,
				func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
					// Copy the MCP arguments to a map
					options := request.Params.Arguments
					if options == nil {
						options = make(map[string]any)
					}

					// Execute the tool's handler, passing the options
					resp, err := resourceTemplate.Handler(request.Params.URI, options)
					if err != nil {
						return nil, err
					}

					return []mcp.ResourceContents{
						mcp.TextResourceContents{
							URI:      resp.URI,
							MIMEType: resp.MIMEType,
							Text:     resp.Content}}, nil
				},
			)
		}
	}
}
