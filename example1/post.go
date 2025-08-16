/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package example1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

// httpPost is a generic function to make HTTP POST requests.
func (c *Config) httpPost(path string, data map[string]any) (string, error) {

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Build the full URL
	url := c.BaseURL + path

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create POST request: %w", err)
	}

	// Add authentication header
	req.Header.Set(c.AuthHeader, c.AuthKey)

	// Set Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check the status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body) // Read the response body for error details
		return "", fmt.Errorf("received non-OK HTTP status: %s, body: %s", resp.Status, string(responseBody))
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Return the response body as a string
	return string(responseBody), nil
}

// ValidatePostParams validates the options for a POST request.
func (c *Config) validatePostParams(toolName string, options map[string]any) (map[string]any, error) {
	// Find the tool definition from the registration
	var toolDef *global.ToolDefinition
	for _, def := range c.RegisterTools() {
		if def.Name == toolName {
			toolDef = &def
			break
		}
	}

	if toolDef == nil {
		return nil, fmt.Errorf("tool '%s' not found in registration", toolName)
	}

	// Validate and build the parameters
	validatedParams := make(map[string]any)
	for _, param := range toolDef.Parameters {
		value, exists := options[param.Name]
		if !exists {
			if param.Required {
				return nil, fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		// Add the parameter as-is to the validatedParams map
		validatedParams[param.Name] = value
	}

	return validatedParams, nil
}
