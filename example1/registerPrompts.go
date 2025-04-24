// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

func (c *Config) RegisterPrompts() []global.PromptDefinition {
	return []global.PromptDefinition{
		{
			Name:        "greeting",
			Description: "A friendly greeting prompt",
			Parameters: []global.Parameter{
				{
					Name:        "name",
					Description: "Name of the person to greet",
					Required:    true,
				},
			},
			Handler: c.GreetingPrompt,
		},
	}

}
