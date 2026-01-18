/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcptypes

import "context"

// This file contains authentication-related interfaces for future implementation.
// These are stubs/placeholders to establish the API contract.

// BearerTokenValidator validates bearer tokens and returns tenant/user context.
// Implement this interface and pass it via WithBearerTokenAuth() option.
type BearerTokenValidator func(token string) (contextData map[string]any, err error)

// OAuth2Provider defines the interface for OAuth2 authentication providers.
// Implement this interface to add OAuth2 support to your MCP server.
type OAuth2Provider interface {
	// GetDeviceCode initiates the OAuth2 device flow
	GetDeviceCode(ctx context.Context) (DeviceCodeResponse, error)

	// ExchangeDeviceCode polls for token using device code
	ExchangeDeviceCode(ctx context.Context, deviceCode string) (TokenResponse, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)

	// ValidateToken checks if an access token is valid
	ValidateToken(ctx context.Context, accessToken string) (bool, error)
}

// DeviceCodeResponse represents the response from GetDeviceCode
type DeviceCodeResponse struct {
	DeviceCode      string
	UserCode        string
	VerificationURI string
	ExpiresIn       int
}

// TokenResponse represents an OAuth2 token response
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
