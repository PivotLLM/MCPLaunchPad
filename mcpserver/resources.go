/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// AddResources registers all resources from resource providers
func (m *MCPServer) AddResources() {

	// Iterate over resource providers and register their resources
	for _, provider := range m.resourceProviders {

		// Call the Register function of the provider to get resource definitions
		resourceDefinitions := provider.RegisterResources()

		// Iterate over the resource definitions and register each resource
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

					// Execute the resource handler, passing the options
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

// AddResourceTemplates registers all resource templates from resource providers
func (m *MCPServer) AddResourceTemplates() {

	// Iterate over resource providers and register their templates
	for _, provider := range m.resourceProviders {

		// Call the Register function of the provider to get resource template definitions
		resourceTemplates := provider.RegisterResourceTemplates()

		// Iterate over the resource template definitions and register each template
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

					// Execute the resource template handler, passing the options
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
