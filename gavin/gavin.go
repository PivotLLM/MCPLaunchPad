// Package gavin provides utility functions for interacting with APIs.
package gavin

import (
	"github.com/PivotLLM/MCPLaunchPad/global"
)

// Ensure Config implements the global.APIClient interface.
var _ global.APIClient = (*Config)(nil)

// Config holds the configuration for the Gavin package.
type Config struct {
	BaseURL string
	Logger  global.Logger
}

// Option defines a function type for configuring the Gavin package.
type Option func(*Config)

// WithBaseURL sets the BaseURL for the Gavin package.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithLogger sets the logger for the Gavin package.
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
