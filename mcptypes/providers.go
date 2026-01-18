/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcptypes

//
// Tools
//

// ToolDefinition represents the structure of a tool
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []*Parameter
	Handler     ToolHandler
	Hints       *ToolHints // Optional hint overrides
}

// ToolHandler defines the function signature for tool handlers
type ToolHandler func(options map[string]any) (string, error)

// ToolProvider defines an interface for providing tools
type ToolProvider interface {
	RegisterTools() []ToolDefinition
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

// ResourceHandler defines the function signature for resource handlers
type ResourceHandler func(uri string, options map[string]any) (ResourceResponse, error)

// ResourceProvider defines an interface for providing resources
type ResourceProvider interface {
	RegisterResources() []ResourceDefinition
	RegisterResourceTemplates() []ResourceTemplateDefinition
}

//
// Prompts
//

// PromptDefinition represents the structure of a prompt
type PromptDefinition struct {
	Name        string
	Description string
	Parameters  []*Parameter
	Handler     PromptHandler
}

// Messages represents a collection of messages
type Messages []Message

// Message represents a single message
type Message struct {
	Role    string
	Content string
}

// PromptHandler defines the function signature for prompt handlers
type PromptHandler func(options map[string]any) (string, Messages, error)

// PromptProvider defines an interface for providing prompts
type PromptProvider interface {
	RegisterPrompts() []PromptDefinition
}
