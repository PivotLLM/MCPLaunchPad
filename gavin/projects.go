package gavin

import "fmt"

// CreateProject creates a new project and decomposes it into tasks using the LLM planner
func (c *Config) CreateProject(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("create_project", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/project/create", postParams)
}

// DecomposeTasks decomposes a project description into tasks using the LLM planner without creating a project
func (c *Config) DecomposeTasks(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("decompose_tasks", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/project/decompose", postParams)
}

// ListProjects retrieves a list of all projects
func (c *Config) ListProjects(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	queryParams, err := c.validateURLParams("list_projects", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpGet function
	return c.httpGet("/project/list", queryParams)
}

// GetProjectByID retrieves a project by its ID
func (c *Config) GetProjectByID(options map[string]any) (string, error) {
	// Safely get the project_id from options
	projectID, ok := options["project_id"].(string)
	if !ok || projectID == "" {
		return "", fmt.Errorf("project_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/project/"+projectID, params)
}

// DeleteProjectByID deletes a project by its ID
func (c *Config) DeleteProjectByID(options map[string]any) (string, error) {
	// Safely get the project_id from options
	projectID, ok := options["project_id"].(string)
	if !ok || projectID == "" {
		return "", fmt.Errorf("project_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpDelete function
	return c.httpDelete("/project/"+projectID, params)
}
