// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

// RegisterResources will be called by the MCP server. This is a very simplistic
// example and in practice would likely provide a list of resources instead of just one.
// The MCP server is also capable of returning more than plain text, so users may wish
// to expand this, or bypass this wrapper entirely and either import
// "github.com/mark3labs/mcp-go/mcp" or implement resources within the mcpserver package.
func (c *Config) RegisterResources() []global.ResourceDefinition {
	return []global.ResourceDefinition{
		{
			Name:        "readme.txt",
			Description: "A readme file",
			MIMEType:    "text/plain",
			URI:         "file:///home/readme.txt",
			Handler:     c.ResourceHandler,
		},
	}
}

// RegisterResourceTemplates will be called by the MCP server.
func (c *Config) RegisterResourceTemplates() []global.ResourceTemplateDefinition {
	return []global.ResourceTemplateDefinition{
		{
			Name:        "ABC Data",
			Description: "ABC data provides all the information you need on the alphabet",
			MIMEType:    "text/plain",
			URITemplate: "abc:///info/{letter_or_number}",
			Handler:     c.ResourceHandler,
		},
	}
}
