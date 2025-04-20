package gavin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

// httpGet is a generic function to make HTTP GET requests.
func (c *Config) httpGet(path string, queryParams map[string]string) (string, error) {

	// Build the full URL
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Append the path to the base URL
	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Add query parameters
	if len(queryParams) > 0 {
		query := fullURL.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		fullURL.RawQuery = query.Encode()
	}

	// Make the HTTP GET request
	resp, err := http.Get(fullURL.String())
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Parse the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Return the raw JSON response
	return string(responseBody), nil
}

// ValidateURLParams validates the options for a GET request.
func (c *Config) validateURLParams(toolName string, options map[string]any) (map[string]string, error) {

	// Find the tool definition from the registration
	var toolDef *global.ToolDefinition
	for _, def := range c.Register() {
		if def.Name == toolName {
			toolDef = &def
			break
		}
	}

	if toolDef == nil {
		return nil, fmt.Errorf("tool '%s' not found in registration", toolName)
	}

	// Validate and build query parameters
	queryParams := make(map[string]string)
	for _, param := range toolDef.Parameters {
		value, exists := options[param.Name]
		if !exists {
			if param.Required {
				return nil, fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		// Convert the value to a string or handle numbers
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int, int8, int16, int32, int64, float32, float64:
			strValue = fmt.Sprintf("%v", v)
		default:
			return nil, fmt.Errorf("parameter '%s' must be a string or a number", param.Name)
		}

		queryParams[param.Name] = strValue
	}

	return queryParams, nil
}
