/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// GoogleOAuth2Provider implements OAuth2Provider for Google
type GoogleOAuth2Provider struct {
	clientID     string
	clientSecret string
	scopes       []string
	httpClient   *http.Client

	// Google OAuth2 endpoints
	authURL       string
	tokenURL      string
	deviceAuthURL string
}

// Ensure GoogleOAuth2Provider implements OAuth2Provider
var _ mcptypes.OAuth2Provider = (*GoogleOAuth2Provider)(nil)

// NewGoogleProvider creates a new Google OAuth2 provider
func NewGoogleProvider(clientID, clientSecret string, scopes []string) *GoogleOAuth2Provider {
	return &GoogleOAuth2Provider{
		clientID:      clientID,
		clientSecret:  clientSecret,
		scopes:        scopes,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		authURL:       "https://accounts.google.com/o/oauth2/v2/auth",
		tokenURL:      "https://oauth2.googleapis.com/token",
		deviceAuthURL: "https://oauth2.googleapis.com/device/code",
	}
}

// GetDeviceCode initiates the OAuth2 device flow
func (g *GoogleOAuth2Provider) GetDeviceCode(ctx context.Context) (mcptypes.DeviceCodeResponse, error) {
	// Prepare request body
	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("scope", strings.Join(g.scopes, " "))

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", g.deviceAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return mcptypes.DeviceCodeResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return mcptypes.DeviceCodeResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcptypes.DeviceCodeResponse{}, fmt.Errorf("device code request failed: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var result struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURL         string `json:"verification_url"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
		VerificationURLComplete string `json:"verification_url_complete,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return mcptypes.DeviceCodeResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return mcptypes.DeviceCodeResponse{
		DeviceCode:      result.DeviceCode,
		UserCode:        result.UserCode,
		VerificationURI: result.VerificationURL,
		ExpiresIn:       result.ExpiresIn,
	}, nil
}

// ExchangeDeviceCode polls for token using device code
func (g *GoogleOAuth2Provider) ExchangeDeviceCode(ctx context.Context, deviceCode string) (mcptypes.TokenResponse, error) {
	// Prepare request body
	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("client_secret", g.clientSecret)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", g.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for pending authorization
	var errorResp struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		// These errors are expected during polling
		if errorResp.Error == "authorization_pending" || errorResp.Error == "slow_down" {
			return mcptypes.TokenResponse{}, fmt.Errorf("%s", errorResp.Error)
		}
		return mcptypes.TokenResponse{}, fmt.Errorf("token exchange failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
	}

	// Parse success response
	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return mcptypes.TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (g *GoogleOAuth2Provider) RefreshToken(ctx context.Context, refreshToken string) (mcptypes.TokenResponse, error) {
	// Prepare request body
	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("client_secret", g.clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", g.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return mcptypes.TokenResponse{}, fmt.Errorf("token refresh failed: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return mcptypes.TokenResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// Google may not return a new refresh token, keep the old one
	if result.RefreshToken == "" {
		result.RefreshToken = refreshToken
	}

	return mcptypes.TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// ValidateToken checks if an access token is valid
func (g *GoogleOAuth2Provider) ValidateToken(ctx context.Context, accessToken string) (bool, error) {
	// Use Google's token info endpoint to validate
	tokenInfoURL := "https://oauth2.googleapis.com/tokeninfo?access_token=" + url.QueryEscape(accessToken)

	req, err := http.NewRequestWithContext(ctx, "GET", tokenInfoURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Token is valid if status is 200
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// GetUserInfo retrieves user information from Google
func (g *GoogleOAuth2Provider) GetUserInfo(ctx context.Context, accessToken string) (map[string]any, error) {
	// Use Google's userinfo endpoint
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s - %s", resp.Status, string(body))
	}

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return userInfo, nil
}

// DeviceFlowWithPolling performs the complete device flow with automatic polling
func (g *GoogleOAuth2Provider) DeviceFlowWithPolling(ctx context.Context, interval time.Duration) (mcptypes.TokenResponse, mcptypes.DeviceCodeResponse, error) {
	// Get device code
	deviceResp, err := g.GetDeviceCode(ctx)
	if err != nil {
		return mcptypes.TokenResponse{}, mcptypes.DeviceCodeResponse{}, err
	}

	// Poll for token
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeout := time.After(time.Duration(deviceResp.ExpiresIn) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return mcptypes.TokenResponse{}, deviceResp, ctx.Err()
		case <-timeout:
			return mcptypes.TokenResponse{}, deviceResp, fmt.Errorf("device code expired")
		case <-ticker.C:
			tokenResp, err := g.ExchangeDeviceCode(ctx, deviceResp.DeviceCode)
			if err != nil {
				// Continue polling for expected errors
				if strings.Contains(err.Error(), "authorization_pending") {
					continue
				}
				if strings.Contains(err.Error(), "slow_down") {
					ticker.Reset(interval * 2) // Slow down polling
					continue
				}
				return mcptypes.TokenResponse{}, deviceResp, err
			}
			return tokenResp, deviceResp, nil
		}
	}
}

// CreateBearerTokenValidator creates a bearer token validator from the OAuth2 provider
func (g *GoogleOAuth2Provider) CreateBearerTokenValidator() mcptypes.BearerTokenValidator {
	return func(token string) (map[string]any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		valid, err := g.ValidateToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}
		if !valid {
			return nil, fmt.Errorf("invalid token")
		}

		// Get user info
		userInfo, err := g.GetUserInfo(ctx, token)
		if err != nil {
			// Token is valid but we couldn't get user info
			// Return minimal context
			return map[string]any{
				"authenticated": true,
			}, nil
		}

		return userInfo, nil
	}
}
