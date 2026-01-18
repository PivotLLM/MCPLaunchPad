/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcptypes

// ToolHints represents MCP tool annotation hints.
// All fields are pointers to allow nil (inherit from defaults) vs explicit true/false.
type ToolHints struct {
	ReadOnlyHint    *bool // If true, the tool is read-only
	DestructiveHint *bool // If true, the tool may perform destructive updates
	IdempotentHint  *bool // If true, repeated calls with same args have no additional effect
	OpenWorldHint   *bool // If true, tool interacts with external entities
}

// HintOption is a function that sets a specific hint value
type HintOption func(*ToolHints)

// NewHints creates a new ToolHints with optional hint setters
// Supports variadic constructor pattern: NewHints(ReadOnly(true), Destructive(false))
func NewHints(options ...HintOption) *ToolHints {
	hints := &ToolHints{}
	for _, opt := range options {
		opt(hints)
	}
	return hints
}

// Method chaining pattern methods

// ReadOnly sets the ReadOnlyHint and returns the hints for chaining
func (h *ToolHints) ReadOnly(value bool) *ToolHints {
	h.ReadOnlyHint = &value
	return h
}

// Destructive sets the DestructiveHint and returns the hints for chaining
func (h *ToolHints) Destructive(value bool) *ToolHints {
	h.DestructiveHint = &value
	return h
}

// Idempotent sets the IdempotentHint and returns the hints for chaining
func (h *ToolHints) Idempotent(value bool) *ToolHints {
	h.IdempotentHint = &value
	return h
}

// OpenWorld sets the OpenWorldHint and returns the hints for chaining
func (h *ToolHints) OpenWorld(value bool) *ToolHints {
	h.OpenWorldHint = &value
	return h
}

// Variadic constructor option functions

// ReadOnly creates a HintOption that sets ReadOnlyHint
func ReadOnly(value bool) HintOption {
	return func(h *ToolHints) {
		h.ReadOnlyHint = &value
	}
}

// Destructive creates a HintOption that sets DestructiveHint
func Destructive(value bool) HintOption {
	return func(h *ToolHints) {
		h.DestructiveHint = &value
	}
}

// Idempotent creates a HintOption that sets IdempotentHint
func Idempotent(value bool) HintOption {
	return func(h *ToolHints) {
		h.IdempotentHint = &value
	}
}

// OpenWorld creates a HintOption that sets OpenWorldHint
func OpenWorld(value bool) HintOption {
	return func(h *ToolHints) {
		h.OpenWorldHint = &value
	}
}
