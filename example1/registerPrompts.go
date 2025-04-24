// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"fmt"

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

// GreetingPrompt is an example
func (c *Config) GreetingPrompt(options map[string]any) (string, global.Messages, error) {
	name, ok := options["name"].(string)
	if !ok || name == "" {
		name = "friend"
	}
	return "A friendly greeting", []global.Message{
		{
			Role:    "assistant",
			Content: fmt.Sprintf("Hello, %s! How can I help you today?", name),
		},
	}, nil
}
