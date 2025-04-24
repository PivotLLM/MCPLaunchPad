// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package example1

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// httpDelete is a generic function to make HTTP DELETE requests.
func (c *Config) httpDelete(path string, queryParams map[string]string) (string, error) {

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

	// Create a new HTTP DELETE request
	req, err := http.NewRequest("DELETE", fullURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create DELETE request: %w", err)
	}

	// Add authentication header
	req.Header.Set(c.AuthHeader, c.AuthKey)

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute DELETE request: %w", err)
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
