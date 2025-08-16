/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

// Package example2 provides a simple time service to an MCP client
package example2

import (
	"fmt"
	"strings"
	"time"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

// Ensure Config implements the global.ToolProvider interface.
var _ global.ToolProvider = (*Config)(nil)

// Config serves as the package's object and holds configuration information
type Config struct {
	Logger global.Logger
}

// Option defines a function type for configuration options
type Option func(*Config)

// WithLogger sets the logger
func WithLogger(logger global.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// New creates a new Config instance with the provided options.
func New(options ...Option) *Config {
	config := &Config{}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// RegisterTools will be called by the MCP server to obtain information and register the tools
func (c *Config) RegisterTools() []global.ToolDefinition {
	return []global.ToolDefinition{
		{
			Name:        "get_time",
			Description: "Get the current time. Optionally set 'time_format' to '12' for 12-hour format or '24' for 24-hour format.",
			Parameters: []global.Parameter{
				{
					Name:        "time_format",
					Description: "Time format (12 or 24).",
					Required:    false,
				},
			},
			Handler: c.GetTime,
		},
	}
}

func (c *Config) GetTime(options map[string]any) (string, error) {

	// Assume 12-hour format by default
	h24 := false

	// Safely check for the 'time_format' key in the options map
	timeFormatInterface, exists := options["time_format"]
	if exists {

		// Convert it to a string
		timeFormatStr, ok := timeFormatInterface.(string)
		if !ok {
			return "", fmt.Errorf("'time_format' must be a string")
		}

		// Trim whitespace
		timeFormatStr = strings.TrimSpace(timeFormatStr)

		switch timeFormatStr {
		case "12":
			h24 = false
		case "24":
			h24 = true
		default:
			return "", fmt.Errorf("invalid 'time_format' value; must be '12' or '24'")
		}

	}

	var formattedTime string
	if h24 {
		formattedTime = time.Now().Format("15:04:05")
	} else {
		formattedTime = time.Now().Format("03:04:05 PM")
	}
	return formattedTime, nil
}
