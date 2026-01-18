/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcptypes

// Parameter represents a tool parameter with full JSON Schema support
type Parameter struct {
	Name        string
	Description string
	Required    bool
	Type        string // "string", "number", "integer", "boolean", "array", "object", "null"

	// String validation
	Pattern   *string
	MinLength *int
	MaxLength *int
	Format    *string // "date-time", "email", "uri", etc.

	// Numeric validation
	Minimum          *float64
	Maximum          *float64
	ExclusiveMinimum *bool
	ExclusiveMaximum *bool
	MultipleOf       *float64

	// Array validation
	Items       *Parameter // Schema for array items
	MinItems    *int
	MaxItems    *int
	UniqueItems *bool

	// Object validation
	Properties           map[string]*Parameter
	AdditionalProperties *bool

	// Enum constraint (works with any type)
	Enum []any

	// Default value
	Default any
}

// Helper constructors for common parameter types

// StringParam creates a string parameter
func StringParam(name, description string, required bool) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "string",
	}
}

// NumberParam creates a number parameter
func NumberParam(name, description string, required bool) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "number",
	}
}

// IntegerParam creates an integer parameter
func IntegerParam(name, description string, required bool) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "integer",
	}
}

// BoolParam creates a boolean parameter
func BoolParam(name, description string, required bool) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "boolean",
	}
}

// ArrayParam creates an array parameter
func ArrayParam(name, description string, required bool, itemType *Parameter) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "array",
		Items:       itemType,
	}
}

// ObjectParam creates an object parameter
func ObjectParam(name, description string, required bool, properties map[string]*Parameter) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		Type:        "object",
		Properties:  properties,
	}
}

// Fluent API methods for validation

// WithPattern sets a regex pattern for string validation
func (p *Parameter) WithPattern(pattern string) *Parameter {
	p.Pattern = &pattern
	return p
}

// WithFormat sets a format for string validation
func (p *Parameter) WithFormat(format string) *Parameter {
	p.Format = &format
	return p
}

// WithMinLength sets minimum string length
func (p *Parameter) WithMinLength(min int) *Parameter {
	p.MinLength = &min
	return p
}

// WithMaxLength sets maximum string length
func (p *Parameter) WithMaxLength(max int) *Parameter {
	p.MaxLength = &max
	return p
}

// WithMinimum sets minimum numeric value
func (p *Parameter) WithMinimum(min float64) *Parameter {
	p.Minimum = &min
	return p
}

// WithMaximum sets maximum numeric value
func (p *Parameter) WithMaximum(max float64) *Parameter {
	p.Maximum = &max
	return p
}

// WithExclusiveMinimum sets whether minimum is exclusive
func (p *Parameter) WithExclusiveMinimum(exclusive bool) *Parameter {
	p.ExclusiveMinimum = &exclusive
	return p
}

// WithExclusiveMaximum sets whether maximum is exclusive
func (p *Parameter) WithExclusiveMaximum(exclusive bool) *Parameter {
	p.ExclusiveMaximum = &exclusive
	return p
}

// WithMultipleOf sets the multipleOf constraint for numbers
func (p *Parameter) WithMultipleOf(value float64) *Parameter {
	p.MultipleOf = &value
	return p
}

// WithMinItems sets minimum array length
func (p *Parameter) WithMinItems(min int) *Parameter {
	p.MinItems = &min
	return p
}

// WithMaxItems sets maximum array length
func (p *Parameter) WithMaxItems(max int) *Parameter {
	p.MaxItems = &max
	return p
}

// WithUniqueItems sets whether array items must be unique
func (p *Parameter) WithUniqueItems(unique bool) *Parameter {
	p.UniqueItems = &unique
	return p
}

// WithAdditionalProperties sets whether additional object properties are allowed
func (p *Parameter) WithAdditionalProperties(allowed bool) *Parameter {
	p.AdditionalProperties = &allowed
	return p
}

// WithEnum sets enum constraint
func (p *Parameter) WithEnum(values ...any) *Parameter {
	p.Enum = values
	return p
}

// WithDefault sets a default value
func (p *Parameter) WithDefault(value any) *Parameter {
	p.Default = value
	return p
}
