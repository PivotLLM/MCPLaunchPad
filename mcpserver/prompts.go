/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (m *MCPServer) AddPrompts() {

	// Iterate over prompt providers and register their prompts
	for _, provider := range m.promptProviders {

		// Call the Register function of the provider to get tool definitions
		promptDefinitions := provider.RegisterPrompts()

		// Iterate over the tool definitions and register each tool
		for _, prompt := range promptDefinitions {

			// Combine description and parameters into a slice of options
			options := []mcp.PromptOption{
				mcp.WithPromptDescription(prompt.Description),
			}

			for _, param := range prompt.Parameters {
				argOptions := []mcp.ArgumentOption{mcp.ArgumentDescription(param.Description)}
				if param.Required {
					argOptions = append(argOptions, mcp.RequiredArgument())
				}
				options = append(options, mcp.WithArgument(param.Name, argOptions...))
			}

			// Create the tool with all options
			newPrompt := mcp.NewPrompt(prompt.Name, options...)

			// Register the tool with the MCP server, creating a handler compatible with the MCP server
			// that wraps the tool's handler function with the provided options
			m.srv.AddPrompt(newPrompt, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {

				// Copy the MCP arguments to a map
				options := make(map[string]any)
				for key, value := range req.Params.Arguments {
					options[key] = value
				}

				// Execute the tool's handler, passing the options
				str, messages, err := prompt.Handler(options)
				if err != nil {
					return nil, err
				}

				// Convert the results to a PromptMessage struct
				var promptMessages []mcp.PromptMessage
				for _, message := range messages {

					var role mcp.Role
					role = mcp.RoleAssistant
					if message.Role == "user" {
						role = mcp.RoleUser
					}

					promptMessages = append(promptMessages,
						mcp.NewPromptMessage(role, mcp.NewTextContent(message.Content)))
				}
				return mcp.NewGetPromptResult(str, promptMessages), nil
			})
		}
	}
}
