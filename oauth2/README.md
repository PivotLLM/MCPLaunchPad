# OAuth2 Package

Reusable OAuth2 implementation for MCP servers, with built-in Google OAuth2 support.

## Features

- **Google OAuth2**: Full implementation of Google OAuth2 device flow
- **Bearer Token Validation**: Convert OAuth2 tokens to bearer token validators
- **User Info**: Automatic user information retrieval
- **Automatic Polling**: Built-in device code polling with timeout handling

## Quick Start

### Basic Usage with Google OAuth2

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/PivotLLM/MCPLaunchPad/oauth2"
)

func main() {
    // Create Google OAuth2 provider
    provider := oauth2.NewGoogleProvider(
        "your-client-id.apps.googleusercontent.com",
        "your-client-secret",
        []string{"email", "profile"},
    )

    // Perform device flow with automatic polling
    tokenResp, deviceResp, err := provider.DeviceFlowWithPolling(
        context.Background(),
        5*time.Second, // Poll every 5 seconds
    )
    if err != nil {
        panic(err)
    }

    // Display device code to user
    fmt.Printf("Go to: %s\n", deviceResp.VerificationURI)
    fmt.Printf("Enter code: %s\n", deviceResp.UserCode)

    // Use the access token
    fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
}
```

### Integration with MCPServer

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/PivotLLM/MCPLaunchPad/mcpserver"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
    "github.com/PivotLLM/MCPLaunchPad/oauth2"
)

func main() {
    // Create OAuth2 provider
    provider := oauth2.NewGoogleProvider(
        "your-client-id.apps.googleusercontent.com",
        "your-client-secret",
        []string{"email", "profile"},
    )

    // Perform device flow
    tokenResp, deviceResp, err := provider.DeviceFlowWithPolling(
        context.Background(),
        5*time.Second,
    )
    if err != nil {
        panic(err)
    }

    fmt.Printf("Visit: %s and enter: %s\n", deviceResp.VerificationURI, deviceResp.UserCode)

    // Create bearer token validator from OAuth2
    validator := provider.CreateBearerTokenValidator()

    // Create MCP server with bearer token auth
    srv, err := mcpserver.New(
        mcpserver.WithTransportHTTP("localhost:8080"),
        mcpserver.WithBearerTokenAuth(validator),
        mcpserver.WithToolProviders([]mcptypes.ToolProvider{myProvider}),
    )
    if err != nil {
        panic(err)
    }

    srv.Start()
}
```

## Google OAuth2 Setup

### 1. Create OAuth2 Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Navigate to "APIs & Services" → "Credentials"
4. Click "Create Credentials" → "OAuth client ID"
5. Select "Desktop app" as application type
6. Download the credentials JSON or copy client ID and secret

### 2. Configure Scopes

Common Google OAuth2 scopes:
- `email` - User's email address
- `profile` - Basic profile information
- `openid` - OpenID Connect support
- `https://www.googleapis.com/auth/userinfo.email` - Email scope (full URI)
- `https://www.googleapis.com/auth/userinfo.profile` - Profile scope (full URI)

### 3. Test Mode vs Production

For development:
- Add test users in Google Cloud Console under "OAuth consent screen"
- Test users can authenticate without app verification

For production:
- Complete OAuth consent screen verification process
- App must be approved by Google

## API Reference

### GoogleOAuth2Provider

#### NewGoogleProvider(clientID, clientSecret string, scopes []string)

Creates a new Google OAuth2 provider.

#### GetDeviceCode(ctx context.Context) (DeviceCodeResponse, error)

Initiates the OAuth2 device flow and returns device code information.

#### ExchangeDeviceCode(ctx context.Context, deviceCode string) (TokenResponse, error)

Exchanges device code for access token. Returns special errors:
- `authorization_pending` - User hasn't authorized yet
- `slow_down` - Should increase polling interval

#### RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)

Refreshes an access token using a refresh token.

#### ValidateToken(ctx context.Context, accessToken string) (bool, error)

Validates an access token.

#### GetUserInfo(ctx context.Context, accessToken string) (map[string]any, error)

Retrieves user information from Google.

#### DeviceFlowWithPolling(ctx context.Context, interval time.Duration) (TokenResponse, DeviceCodeResponse, error)

Performs complete device flow with automatic polling. Recommended interval: 5 seconds.

#### CreateBearerTokenValidator() BearerTokenValidator

Creates a bearer token validator for use with MCPServer.

## Device Flow Explained

The OAuth2 device flow allows users to authenticate on a separate device (like their phone or computer browser):

1. **Request Device Code**: App requests a device code from Google
2. **Display to User**: App shows verification URL and user code
3. **User Authenticates**: User opens URL and enters code on their device
4. **Poll for Token**: App polls Google's server for access token
5. **Receive Token**: Once user authorizes, app receives access token

This is perfect for CLI applications and devices without browsers.

## Security Notes

- **Never commit credentials**: Store client ID/secret in environment variables or secure vaults
- **Use HTTPS**: Always use HTTPS in production for token exchange
- **Token Storage**: Store tokens securely (encrypted database, Redis with encryption)
- **Token Rotation**: Implement refresh token rotation for security
- **Scope Minimization**: Only request scopes you actually need

## Example Integration

See `examples/oauth/` for a complete working example with:
- Device flow authentication
- Token storage
- Bearer token validation
- MCP server integration
