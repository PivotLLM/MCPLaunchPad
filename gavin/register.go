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
	}
}
