package global

// ToolDefinition represents the definition of a tool to be registered.
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []ToolParameter
	Handler     APIHandler
}

// ToolParameter represents a parameter for a tool.
type ToolParameter struct {
	Name        string
	Description string
	Required    bool
}

// APIClient defines an interface for interacting with APIs and registering tools.
type APIClient interface {
	Register() []ToolDefinition
}

// APIHandler defines a function type for our tool handler. This avoids the API package
// having to import definitions from the MCP server, etc.
type APIHandler func(options map[string]any) (string, error)
