// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"errors"
	"fmt"

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

// ResourceHandler is a simple handler that returns a readme file
func (c *Config) ResourceHandler(uri string, options map[string]any) (global.ResourceResponse, error) {

	// Check if the URI is valid
	if uri != "file:///home/readme.txt" {
		return global.ResourceResponse{}, errors.New("invalid URI")
	}

	// Build some content because this is an example
	msg := "This is a simple readme file.\nIf it was a real file, it would hopefully have meaningful content.\nHave a great day!"

	// If the client sent any options, add them
	if len(options) > 0 {
		msg += "\n\nI noticed some options:\n"
		for k, y := range options {
			msg += fmt.Sprintf("%s: %v\n", k, y)
		}
	}

	// Return the readme content
	return global.ResourceResponse{
		URI:      uri,
		MIMEType: "text/plain",
		Content:  "This is a simple readme file.\nIf it was a real file, it would hopefully have meaningful content.\nHave a great day!"}, nil
}
