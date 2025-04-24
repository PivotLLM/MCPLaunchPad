// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package global

// Parameter represents a parameter for a tool, resource, or prompt
type Parameter struct {
	Name        string
	Description string
	Required    bool
}

//
// Tools
//

// ToolDefinition represents the structure of a tool
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []Parameter
	Handler     ToolHandler
}

// ToolHandler defines the function signature for our tool handler
type ToolHandler func(options map[string]any) (string, error)

// ToolProvider defines an interface for providing tools
type ToolProvider interface {
	RegisterTools() []ToolDefinition
}

// NewTools is a helper function that returns an empty slice of ToolDefinition
//
//goland:noinspection GoUnusedExportedFunction
func NewTools() []ToolDefinition {
	return []ToolDefinition{}
}

//
// Resources
//

// ResourceDefinition represents the structure of a resource
type ResourceDefinition struct {
	Name        string
	Description string
	MIMEType    string
	URI         string
	Handler     ResourceHandler
}

// ResourceTemplateDefinition represents the structure of a resource template
type ResourceTemplateDefinition struct {
	Name        string
	Description string
	MIMEType    string
	URITemplate string
	Handler     ResourceHandler
}

// ResourceResponse represents the structure of a resource response
type ResourceResponse struct {
	URI      string
	MIMEType string
	Content  string
}

// ResourceHandler defines the function signature for our resource handler
type ResourceHandler func(uri string, options map[string]any) (ResourceResponse, error)

// ResourceProvider defines an interface for providing resources
type ResourceProvider interface {
	RegisterResources() []ResourceDefinition
	RegisterResourceTemplates() []ResourceTemplateDefinition
}

// NewResources is a helper function that returns an empty slice of ResourceDefinition
//
//goland:noinspection GoUnusedExportedFunction
func NewResources() []ResourceDefinition {
	return []ResourceDefinition{}
}

// NewResourceTemplates is a helper function that returns an empty slice of ResourceTemplateDefinition
//
//goland:noinspection GoUnusedExportedFunction
func NewResourceTemplates() []ResourceTemplateDefinition {
	return []ResourceTemplateDefinition{}
}

//
// Prompts
//

// PromptDefinition represents the structure of a prompt
type PromptDefinition struct {
	Name        string
	Description string
	Parameters  []Parameter
	Handler     PromptHandler
}

// Messages represents a collection of messages
type Messages []Message
type Message struct {
	Role    string
	Content string
}

// PromptHandler defines the function signature for our prompt handler
type PromptHandler func(options map[string]any) (string, Messages, error)

// PromptProvider defines an interface for providing prompts
type PromptProvider interface {
	RegisterPrompts() []PromptDefinition
}

// NewPrompts is a helper function that returns an empty slice of PromptDefinition
//
//goland:noinspection GoUnusedExportedFunction
func NewPrompts() []PromptDefinition {
	return []PromptDefinition{}
}
