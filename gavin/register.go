package gavin

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

// RegisterTools registers all tools from the gavin package with the MCP server.
func (c *Config) Register() []global.ToolDefinition {
	return []global.ToolDefinition{
		{
			Name:        "get_users",
			Description: "Fetch a list of users with optional pagination. Use 'skip' and 'limit' to control pagination.",
			Parameters: []global.ToolParameter{
				{
					Name:        "skip",
					Description: "Number of records to skip.",
					Required:    false,
				},
				{
					Name:        "limit",
					Description: "Maximum number of records to return.",
					Required:    false,
				},
			},
			Handler: c.GetUsers,
		},
		{
			Name:        "create_user",
			Description: "Create a new user with username, email, and password.",
			Parameters: []global.ToolParameter{
				{
					Name:        "username",
					Description: "The username of the new user.",
					Required:    true,
				},
				{
					Name:        "email",
					Description: "The email address of the new user.",
					Required:    true,
				},
				{
					Name:        "password",
					Description: "The password for the new user.",
					Required:    true,
				},
			},
			Handler: c.CreateUser,
		},
		{
			Name:        "get_user",
			Description: "Get details of a specific user by user_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "user_id",
					Description: "The ID of the user to get.",
					Required:    true,
				},
			},
			Handler: c.GetUserByID,
		},
		{
			Name:        "delete_user",
			Description: "Delete a user by user_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "user_id",
					Description: "The ID of the user to delete.",
					Required:    true,
				},
			},
			Handler: c.DeleteUserByID,
		},
		{
			Name:        "create_api_key",
			Description: "Create a new API key with a name and optionally a user_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "name",
					Description: "The name of the API key.",
					Required:    true,
				},
				{
					Name:        "user_id",
					Description: "The ID of the user associated with the API key.",
					Required:    false,
				},
			},
			Handler: c.CreateAPIKey,
		},
		{
			Name:        "list_api_keys",
			Description: "List all API keys with optional pagination.",
			Parameters: []global.ToolParameter{
				{
					Name:        "skip",
					Description: "Number of records to skip.",
					Required:    false,
				},
				{
					Name:        "limit",
					Description: "Maximum number of records to return.",
					Required:    false,
				},
			},
			Handler: c.ListAPIKeys,
		},
		{
			Name:        "get_api_key",
			Description: "Fetch details of a specific API key by api_key_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "api_key_id",
					Description: "The ID of the API key to fetch.",
					Required:    true,
				},
			},
			Handler: c.GetAPIKeyByID,
		},
		{
			Name:        "delete_api_key",
			Description: "Delete an API key by api_key_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "api_key_id",
					Description: "The ID of the API key to delete.",
					Required:    true,
				},
			},
			Handler: c.DeleteAPIKeyByID,
		},
		{
			Name:        "create_project",
			Description: "Create a new project and decompose it into tasks using the LLM planner.",
			Parameters: []global.ToolParameter{
				{
					Name:        "title",
					Description: "The title of the project.",
					Required:    true,
				},
				{
					Name:        "description",
					Description: "The description of the project.",
					Required:    true,
				},
			},
			Handler: c.CreateProject,
		},
		{
			Name:        "decompose_tasks",
			Description: "Decompose a project description into tasks using the LLM planner without creating a project.",
			Parameters: []global.ToolParameter{
				{
					Name:        "description",
					Description: "The description of the project to decompose.",
					Required:    true,
				},
			},
			Handler: c.DecomposeTasks,
		},
		{
			Name:        "list_projects",
			Description: "List all projects with optional pagination.",
			Parameters: []global.ToolParameter{
				{
					Name:        "skip",
					Description: "Number of records to skip.",
					Required:    false,
				},
				{
					Name:        "limit",
					Description: "Maximum number of records to return.",
					Required:    false,
				},
			},
			Handler: c.ListProjects,
		},
		{
			Name:        "get_project",
			Description: "Fetch details of a specific project by project_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "project_id",
					Description: "The ID of the project to fetch.",
					Required:    true,
				},
			},
			Handler: c.GetProjectByID,
		},
		{
			Name:        "delete_project",
			Description: "Delete a project by project_id.",
			Parameters: []global.ToolParameter{
				{
					Name:        "project_id",
					Description: "The ID of the project to delete.",
					Required:    true,
				},
			},
			Handler: c.DeleteProjectByID,
		},
		{
			Name:        "send_task",
			Description: "Send a new task for processing or update an existing one.",
			Parameters: []global.ToolParameter{
				{
					Name:        "id",
					Description: "The external ID of the task.",
					Required:    true,
				},
				{
					Name:        "sessionId",
					Description: "The session ID associated with the task.",
					Required:    true,
				},
				{
					Name:        "message",
					Description: "The message payload for the task.",
					Required:    true,
				},
			},
			Handler: c.SendTask,
		},
		{
			Name:        "get_task",
			Description: "Get the status and details of a task by its external ID.",
			Parameters: []global.ToolParameter{
				{
					Name:        "id",
					Description: "The external ID of the task.",
					Required:    true,
				},
				{
					Name:        "historyLength",
					Description: "Optional parameter to specify the length of task history to retrieve.",
					Required:    false,
				},
			},
			Handler: c.GetTask,
		},
		{
			Name:        "cancel_task",
			Description: "Cancel a task by its external ID.",
			Parameters: []global.ToolParameter{
				{
					Name:        "id",
					Description: "The external ID of the task to cancel.",
					Required:    true,
				},
			},
			Handler: c.CancelTask,
		},
		{
			Name:        "manually_process_task",
			Description: "Manually process a task that is stuck in submitted status.",
			Parameters: []global.ToolParameter{
				{
					Name:        "task_id",
					Description: "The internal task ID to process.",
					Required:    true,
				},
			},
			Handler: c.ManuallyProcessTask,
		},
		{
			Name:        "get_tasks_by_project",
			Description: "Get all tasks for a specific project.",
			Parameters: []global.ToolParameter{
				{
					Name:        "project_id",
					Description: "The ID of the project to fetch tasks for.",
					Required:    true,
				},
			},
			Handler: c.GetTasksByProject,
		},
		{
			Name:        "get_task_details",
			Description: "Get detailed information about a task by its internal or external ID.",
			Parameters: []global.ToolParameter{
				{
					Name:        "task_id",
					Description: "The internal or external ID of the task to fetch details for.",
					Required:    true,
				},
			},
			Handler: c.GetTaskDetails,
		},
		{
			Name:        "send_task_with_subscription",
			Description: "Send a new task for processing or update an existing one and subscribe to SSE for task updates.",
			Parameters: []global.ToolParameter{
				{
					Name:        "id",
					Description: "The external ID of the task.",
					Required:    true,
				},
				{
					Name:        "sessionId",
					Description: "The session ID associated with the task.",
					Required:    true,
				},
				{
					Name:        "message",
					Description: "The message payload for the task.",
					Required:    true,
				},
			},
			Handler: c.SendTaskWithSubscription,
		},
	}
}
