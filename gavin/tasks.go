package gavin

import "fmt"

// SendTask sends a new task for processing or updates an existing one
func (c *Config) SendTask(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("send_task", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/tasks/send", postParams)
}

// GetTask retrieves the status and details of a task by its external ID
func (c *Config) GetTask(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("get_task", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/tasks/get", postParams)
}

// CancelTask cancels a task by its external ID
func (c *Config) CancelTask(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("cancel_task", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/tasks/cancel", postParams)
}

// ManuallyProcessTask manually processes a task stuck in submitted status
func (c *Config) ManuallyProcessTask(options map[string]any) (string, error) {
	// Safely get the task_id from options
	taskID, ok := options["task_id"].(string)
	if !ok || taskID == "" {
		return "", fmt.Errorf("task_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]any{}

	// Use the generic httpPost function
	return c.httpPost("/tasks/"+taskID+"/process", params)
}

// GetTasksByProject retrieves all tasks for a specific project
func (c *Config) GetTasksByProject(options map[string]any) (string, error) {
	// Safely get the project_id from options
	projectID, ok := options["project_id"].(string)
	if !ok || projectID == "" {
		return "", fmt.Errorf("project_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/tasks/by-project/"+projectID, params)
}

// GetTaskDetails retrieves detailed information about a task by its internal or external ID
func (c *Config) GetTaskDetails(options map[string]any) (string, error) {
	// Safely get the task_id from options
	taskID, ok := options["task_id"].(string)
	if !ok || taskID == "" {
		return "", fmt.Errorf("task_id is required and must be a non-empty string")
	}

	// No additional params
	params := map[string]string{}

	// Use the generic httpGet function
	return c.httpGet("/tasks/"+taskID+"/details", params)
}

// SendTaskWithSubscription sends a new task for processing or updates an existing one and subscribes to SSE for task updates
func (c *Config) SendTaskWithSubscription(options map[string]any) (string, error) {
	// Validate and build query parameters using the helper function
	postParams, err := c.validatePostParams("send_task_with_subscription", options)
	if err != nil {
		return "", err
	}

	// Use the generic httpPost function
	return c.httpPost("/tasks/sendSubscribe", postParams)
}
