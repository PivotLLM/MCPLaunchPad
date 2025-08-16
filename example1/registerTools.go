/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package example1

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

// RegisterTools provides the details of all available tools. This function is called by the
// MCP server, and it registers each of them. This function is also called by helper functions in
// this package to validate options and build query parameters. That should ensure that the
// information provided to the LLM via the MCP server and the implementations remain consistent.
// If you change the Name: fields here, you must also update them in handlers.go because
// handers pass the name field to validation functions.
func (c *Config) RegisterTools() []global.ToolDefinition {
	return []global.ToolDefinition{
		{
			Name:        "list_widgets",
			Description: "Fetch a list of widgets with optional pagination. Use 'offset' and 'limit' for pagination.",
			Parameters: []global.Parameter{
				{
					Name:        "offset",
					Description: "Starting record offset.",
					Required:    false,
				},
				{
					Name:        "limit",
					Description: "Maximum number of records to return.",
					Required:    false,
				},
			},
			Handler: c.GetWidgets,
		},
		{
			Name:        "create_widget",
			Description: "Create a new widget with a name and description, and an optional radius.",
			Parameters: []global.Parameter{
				{
					Name:        "name",
					Description: "The name of the widget.",
					Required:    true,
				},
				{
					Name:        "description",
					Description: "A description of the widget.",
					Required:    true,
				},
				{
					Name:        "radius",
					Description: "The radius of the widget (optional).",
					Required:    false,
				},
			},
			Handler: c.CreateWidget,
		},
		{
			Name:        "get_widget",
			Description: "Get details of a specific widget by id.",
			Parameters: []global.Parameter{
				{
					Name:        "id",
					Description: "The ID of the widget to get.",
					Required:    true,
				},
			},
			Handler: c.GetWidgetByID,
		},
		{
			Name:        "delete_widget",
			Description: "Delete a widget by id.",
			Parameters: []global.Parameter{
				{
					Name:        "id",
					Description: "The ID of the widget to delete.",
					Required:    true,
				},
			},
			Handler: c.DeleteWidgetByID,
		},
	}
}
