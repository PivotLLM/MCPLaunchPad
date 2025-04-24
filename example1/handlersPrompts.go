// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"fmt"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

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
