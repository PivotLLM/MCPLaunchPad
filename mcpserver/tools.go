// Copyright (c) 2025 Tenebris Technologies Inc.
// This software is licensed under the MIT License (see LICENSE for details).

package mcpserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (m *MCPServer) AddTools() {
	// Create a get_time tool
	// This tool returns the current time in a specified format (12-hour or 24-hour)
	setAllowTool := mcp.NewTool("get_time",
		mcp.WithDescription("Get the current time. Optionally set 'time_format' to '12' for 12-hour format or '24' for 24-hour format."),
		mcp.WithString("time_format",
			mcp.Description("Time format (12 or 24)"),
		),
	)

	// Register the "get_time" tool
	m.srv.AddTool(setAllowTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		// Log tool invocation
		m.logger.Debugf("Tool 'get_time' invoked with arguments: %+v", req.Params.Arguments)

		timeFormatInterface, ok := req.Params.Arguments["time_format"]
		if !ok {
			return nil, fmt.Errorf("missing 'time_format' argument")
		}

		// Convert the interface to a string
		timeFormatStr, ok := timeFormatInterface.(string)
		if !ok {
			return nil, fmt.Errorf("'time_format' must be a string")
		}

		// Trim whitespace and check if it's empty
		timeFormatStr = strings.TrimSpace(timeFormatStr)
		if timeFormatStr == "" {
			return nil, fmt.Errorf("'time_format' cannot be empty")
		}

		var formattedTime string
		switch timeFormatStr {
		case "12":
			formattedTime = time.Now().Format("03:04:05 PM")
		case "24":
			formattedTime = time.Now().Format("15:04:05")
		default:
			return nil, fmt.Errorf("invalid 'time_format' value; must be '12' or '24'")
		}

		// Return the formatted time
		return mcp.NewToolResultText(formattedTime), nil
	})
}
