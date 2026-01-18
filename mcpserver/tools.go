/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// Package-level hint defaults (Level 1)
const (
	defaultReadOnly    = false
	defaultDestructive = false
	defaultIdempotent  = false
	defaultOpenWorld   = false
)

// AddTools registers all tools from tool providers
func (m *MCPServer) AddTools() {

	// Iterate over tool providers and register their tools
	for _, provider := range m.toolProviders {

		// Call the Register function of the provider to get tool definitions
		toolDefinitions := provider.RegisterTools()

		// Iterate over the tool definitions and register each tool
		for _, toolDef := range toolDefinitions {

			// Start with description
			toolOptions := []mcp.ToolOption{
				mcp.WithDescription(toolDef.Description),
			}

			// Add parameters
			for _, param := range toolDef.Parameters {
				toolOptions = append(toolOptions, m.parameterToToolOption(param)...)
			}

			// Add hints with three-level resolution
			hints := m.resolveHints(&toolDef)
			toolOptions = append(toolOptions, mcp.WithToolAnnotation(hints))

			// Create the tool with all options
			tool := mcp.NewTool(toolDef.Name, toolOptions...)

			// Register the tool with the MCP server
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

// resolveHints implements three-level hint resolution:
// Level 3 (tool-level) > Level 2 (server-wide config) > Level 1 (package defaults)
func (m *MCPServer) resolveHints(toolDef *mcptypes.ToolDefinition) mcp.ToolAnnotation {
	hints := mcp.ToolAnnotation{}

	// ReadOnlyHint
	if toolDef.Hints != nil && toolDef.Hints.ReadOnlyHint != nil {
		// Level 3: Tool-level override
		hints.ReadOnlyHint = toolDef.Hints.ReadOnlyHint
	} else if m.defaultReadOnlyHint != nil {
		// Level 2: Server-wide config
		hints.ReadOnlyHint = m.defaultReadOnlyHint
	} else {
		// Level 1: Package default
		val := defaultReadOnly
		hints.ReadOnlyHint = &val
	}

	// DestructiveHint
	if toolDef.Hints != nil && toolDef.Hints.DestructiveHint != nil {
		hints.DestructiveHint = toolDef.Hints.DestructiveHint
	} else if m.defaultDestructiveHint != nil {
		hints.DestructiveHint = m.defaultDestructiveHint
	} else {
		val := defaultDestructive
		hints.DestructiveHint = &val
	}

	// IdempotentHint
	if toolDef.Hints != nil && toolDef.Hints.IdempotentHint != nil {
		hints.IdempotentHint = toolDef.Hints.IdempotentHint
	} else if m.defaultIdempotentHint != nil {
		hints.IdempotentHint = m.defaultIdempotentHint
	} else {
		val := defaultIdempotent
		hints.IdempotentHint = &val
	}

	// OpenWorldHint
	if toolDef.Hints != nil && toolDef.Hints.OpenWorldHint != nil {
		hints.OpenWorldHint = toolDef.Hints.OpenWorldHint
	} else if m.defaultOpenWorldHint != nil {
		hints.OpenWorldHint = m.defaultOpenWorldHint
	} else {
		val := defaultOpenWorld
		hints.OpenWorldHint = &val
	}

	return hints
}

// parameterToToolOption converts an mcptypes.Parameter to mcp.ToolOption slice
// Supporting full JSON Schema validation
func (m *MCPServer) parameterToToolOption(param *mcptypes.Parameter) []mcp.ToolOption {
	var options []mcp.ToolOption

	// Build property options based on parameter type
	switch param.Type {
	case "string":
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		// String-specific validations
		if param.Pattern != nil {
			propOpts = append(propOpts, mcp.Pattern(*param.Pattern))
		}
		if param.MinLength != nil {
			propOpts = append(propOpts, mcp.MinLength(*param.MinLength))
		}
		if param.MaxLength != nil {
			propOpts = append(propOpts, mcp.MaxLength(*param.MaxLength))
		}
		// Note: mcp.Format may not be available in current version
		// if param.Format != nil {
		// 	propOpts = append(propOpts, mcp.Format(*param.Format))
		// }
		if param.Enum != nil {
			// Convert []any to []string for enum
			strEnum := make([]string, len(param.Enum))
			for i, v := range param.Enum {
				if s, ok := v.(string); ok {
					strEnum[i] = s
				}
			}
			propOpts = append(propOpts, mcp.Enum(strEnum...))
		}
		options = append(options, mcp.WithString(param.Name, propOpts...))

	case "number", "integer":
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		// Note: Numeric validations may not be fully supported in current mcp-go version
		// if param.Minimum != nil {
		// 	propOpts = append(propOpts, mcp.Minimum(*param.Minimum))
		// }
		// if param.Maximum != nil {
		// 	propOpts = append(propOpts, mcp.Maximum(*param.Maximum))
		// }
		// Note: mcp-go may not support all numeric validations like MultipleOf, ExclusiveMin/Max
		// Add them if the library supports them
		options = append(options, mcp.WithNumber(param.Name, propOpts...))

	case "boolean":
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		options = append(options, mcp.WithBoolean(param.Name, propOpts...))

	case "array":
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		// Array validations
		if param.MinItems != nil {
			propOpts = append(propOpts, mcp.MinItems(*param.MinItems))
		}
		if param.MaxItems != nil {
			propOpts = append(propOpts, mcp.MaxItems(*param.MaxItems))
		}
		// Note: mcp-go may require item schema - check library documentation
		options = append(options, mcp.WithArray(param.Name, propOpts...))

	case "object":
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		// Note: Object properties may need special handling depending on mcp-go library
		options = append(options, mcp.WithObject(param.Name, propOpts...))

	default:
		// Fallback to string for unknown types
		propOpts := []mcp.PropertyOption{mcp.Description(param.Description)}
		if param.Required {
			propOpts = append(propOpts, mcp.Required())
		}
		options = append(options, mcp.WithString(param.Name, propOpts...))
	}

	return options
}
