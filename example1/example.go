// Package example1 demonstrates interacting with APIs
package example1

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

// Ensure Config implements the global.ToolProvider interface.
var _ global.ToolProvider = (*Config)(nil)

// Config serves as the package's object and holds configuration information
type Config struct {
	BaseURL string
	Logger  global.Logger
}

// Option defines a function type for configuration options
type Option func(*Config)

// WithBaseURL sets the BaseURL
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithLogger sets the logger
func WithLogger(logger global.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// New creates a new Config instance with the provided options.
func New(options ...Option) *Config {
	config := &Config{}
	for _, opt := range options {
		opt(config)
	}
	return config
}
