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

// Define a function type for the tool handler
type APIHandler func(options map[string]any) (string, error)
