// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import "fmt"

//
// NOTE: All functions in this file must use the same function signature:
// func (c *Config) <FunctionName>(options map[string]any) (string, error)
//
// The httpPost, httpGet, and httpDelete functions return (string, error) and
// can therefore be returned directly to the MCP server.
//

// CreateWidget creates a new widget
func (c *Config) CreateWidget(options map[string]any) (string, error) {

	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("create_widget", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/widget", postParams)
}

// GetWidgets retrieves a list of all widgets
func (c *Config) GetWidgets(options map[string]any) (string, error) {

	// Validate and build query parameters using the helper function
	queryParams, err := c.validateURLParams("get_widgets", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpGet function
	return c.httpGet("/widget", queryParams)
}

// GetWidgetByID retrieves a widget by their ID
func (c *Config) GetWidgetByID(options map[string]any) (string, error) {

	// Safely get the id from options
	id, ok := options["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/widget/"+id, params)
}

// DeleteWidgetByID deletes a widget by their ID
func (c *Config) DeleteWidgetByID(options map[string]any) (string, error) {

	// Safely get the id from options
	id, ok := options["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpDelete function
	return c.httpDelete("/widget/"+id, params)
}
