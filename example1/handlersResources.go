/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package example1

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

// ResourceHandler is a simple handler that returns a readme file
func (c *Config) ResourceHandler(uri string, options map[string]any) (global.ResourceResponse, error) {

	// Check if the URI is valid
	if uri == "file:///home/readme.txt" {
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
		return global.ResourceResponse{URI: uri, MIMEType: "text/plain", Content: msg}, nil
	}

	if strings.HasPrefix(uri, "abc:///info/") {
		// Get the letter or number from the URI
		letterOrNumber := strings.TrimPrefix(uri, "abc:///info/")

		// Check if the letter or number is valid
		if len(letterOrNumber) != 1 {
			return global.ResourceResponse{}, errors.New("invalid letter or number, one character only")
		}

		// Build some content because this is an example
		msg := fmt.Sprintf("This is information about '%s' because you asked and you're important to me", letterOrNumber)
		return global.ResourceResponse{URI: uri, MIMEType: "text/plain", Content: msg}, nil
	}

	// If the URI is not valid, return an error
	return global.ResourceResponse{}, errors.New("invalid URI")
}
