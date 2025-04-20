package gavin

import "fmt"

//
// NOTE: All functions in this file must use the same function signature:
// func (c *Config) <FunctionName>(options map[string]any) (string, error)
//
// The httpPost, httpGet, and httpDelete functions return (string, error) and
// can therefore be passed directly to the MCP server.
//

// CreateUser creates a new user
func (c *Config) CreateUser(options map[string]any) (string, error) {

	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("create_user", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/admin/users", postParams)
}

// GetUsers retrieves a list of all users
func (c *Config) GetUsers(options map[string]any) (string, error) {

	// Validate and build query parameters using the helper function
	queryParams, err := c.validateURLParams("get_users", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpGet function
	return c.httpGet("/admin/users", queryParams)
}

// GetUserByID retrieves a user by their ID
func (c *Config) GetUserByID(options map[string]any) (string, error) {

	// Safely get the user_id from options
	userID, ok := options["user_id"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/admin/users/"+userID, params)
}

// DeleteUserByID deletes a user by their ID
func (c *Config) DeleteUserByID(options map[string]any) (string, error) {

	// Safely get the user_id from options
	userID, ok := options["user_id"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpDelete function
	return c.httpDelete("/admin/users/"+userID, params)
}
