package gavin

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

	// Create the HTTP DELETE request
	req, err := http.NewRequest(http.MethodDelete, fullURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create DELETE request: %w", err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

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
