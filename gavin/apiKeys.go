package gavin

import "fmt"

// CreateAPIKey creates a new API key
func (c *Config) CreateAPIKey(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("create_api_key", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/admin/api-keys", postParams)
}

// ListAPIKeys retrieves a list of all API keys
func (c *Config) ListAPIKeys(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	queryParams, err := c.validateURLParams("list_api_keys", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpGet function
	return c.httpGet("/admin/api-keys", queryParams)
}

// GetAPIKeyByID retrieves an API key by its ID
func (c *Config) GetAPIKeyByID(options map[string]any) (string, error) {
	// Safely get the api_key_id from options
	apiKeyID, ok := options["api_key_id"].(string)
	if !ok || apiKeyID == "" {
		return "", fmt.Errorf("api_key_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/admin/api-keys/"+apiKeyID, params)
}

// DeleteAPIKeyByID deletes an API key by its ID
func (c *Config) DeleteAPIKeyByID(options map[string]any) (string, error) {
	// Safely get the api_key_id from options
	apiKeyID, ok := options["api_key_id"].(string)
	if !ok || apiKeyID == "" {
		return "", fmt.Errorf("api_key_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpDelete function
	return c.httpDelete("/admin/api-keys/"+apiKeyID, params)
}
