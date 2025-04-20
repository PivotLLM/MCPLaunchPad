package global

// ToolDefinition represents the definition of a tool to be registered.
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []ToolParameter
	Handler     ToolHandler
}

// ToolParameter represents a parameter for a tool.
type ToolParameter struct {
	Name        string
	Description string
	Required    bool
}

// ToolProvider defines an interface for interacting with APIs and registering tools.
type ToolProvider interface {
	Register() []ToolDefinition
}

// ToolHandler defines a function type for our tool handler. This avoids the API package
// having to import definitions from the MCP server, etc.
type ToolHandler func(options map[string]any) (string, error)
