# Authentication Guide for MCPLaunchPad

This guide explains how to add Bearer Token and/or OAuth2 authentication to your MCP server built with MCPLaunchPad.

## Overview

The mcpserver package provides hooks for authentication via:

1. **Bearer Token Authentication** - Simple function-based validation (dependency injection)
2. **OAuth2 Device Flow** - Interface-based provider pattern

Both are designed but not yet implemented. The interfaces are defined in `mcptypes/auth.go` as stubs.

---

## Bearer Token Authentication

### Architecture

Bearer tokens are validated via a user-provided function that you inject during server creation.

### Interface

```go
// mcptypes/auth.go
type BearerTokenValidator func(token string) (contextData map[string]any, err error)
```

### Implementation Steps

#### Step 1: Create Your Validator Function

Create a function that validates tokens and returns context data:

```go
package main

import (
    "fmt"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

func validateBearerToken(token string) (map[string]any, error) {
    // Example: Check against database or cache
    // This is YOUR implementation

    // Simple example - check against hardcoded tokens
    validTokens := map[string]string{
        "secret-token-123": "user-alice",
        "secret-token-456": "user-bob",
    }

    userID, ok := validTokens[token]
    if !ok {
        return nil, fmt.Errorf("invalid token")
    }

    // Return context that will be available to handlers
    return map[string]any{
        "user_id": userID,
        "tenant":  "acme-corp",
        "role":    "admin",
    }, nil
}
```

#### Step 2: Add Authentication Middleware to mcpserver

You'll need to modify `mcpserver/middleware.go` to add bearer token middleware:

```go
// mcpserver/middleware.go (add this function)

func withBearerTokenAuth(validator mcptypes.BearerTokenValidator) server.ServerOption {
    return server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
        return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

            // Extract Authorization header from context
            // Note: You'll need to store headers in context during HTTP handling
            authHeader, ok := ctx.Value("Authorization").(string)
            if !ok || authHeader == "" {
                return mcp.NewToolResultError("Missing Authorization header"),
                    fmt.Errorf("unauthorized")
            }

            // Extract Bearer token
            const prefix = "Bearer "
            if !strings.HasPrefix(authHeader, prefix) {
                return mcp.NewToolResultError("Invalid Authorization format"),
                    fmt.Errorf("unauthorized")
            }

            token := strings.TrimPrefix(authHeader, prefix)

            // Validate token
            contextData, err := validator(token)
            if err != nil {
                return mcp.NewToolResultError("Invalid token"),
                    fmt.Errorf("unauthorized: %w", err)
            }

            // Add context data to request context
            ctx = context.WithValue(ctx, "auth_context", contextData)

            // Call next handler with enriched context
            return next(ctx, request)
        }
    })
}
```

#### Step 3: Add Option to mcpserver

Add to `mcpserver/options.go`:

```go
// WithBearerTokenAuth enables bearer token authentication
func WithBearerTokenAuth(validator mcptypes.BearerTokenValidator) Option {
    return func(m *MCPServer) {
        m.bearerTokenValidator = validator
    }
}
```

And add field to `MCPServer` struct in `mcpserver/mcpserver.go`:

```go
type MCPServer struct {
    // ... existing fields ...
    bearerTokenValidator mcptypes.BearerTokenValidator
}
```

Then in `New()`, conditionally add the middleware:

```go
func New(options ...Option) (*MCPServer, error) {
    // ... existing code ...

    // Build server options
    serverOpts := []server.ServerOption{
        server.WithLogging(),
        server.WithRecovery(),
        server.WithHooks(hooks),
        withRequestLogging(m.logger),
    }

    // Add bearer token auth if configured
    if m.bearerTokenValidator != nil {
        serverOpts = append(serverOpts, withBearerTokenAuth(m.bearerTokenValidator))
    }

    // Create MCP server
    m.srv = server.NewMCPServer(m.name, m.version, serverOpts...)

    // ... rest of code ...
}
```

#### Step 4: Use in Your Application

```go
package main

import (
    "github.com/PivotLLM/MCPLaunchPad/mcpserver"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

func main() {
    srv, _ := mcpserver.New(
        mcpserver.WithTransportHTTP("localhost:8080"),
        mcpserver.WithBearerTokenAuth(validateBearerToken), // Your validator
        mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
    )

    srv.Start()
    // ... handle shutdown ...
}
```

#### Step 5: Access Auth Context in Handlers

Modify your tool handlers to accept context-aware signatures:

```go
// Current signature
type ToolHandler func(options map[string]any) (string, error)

// Context-aware signature (future enhancement)
type ContextAwareToolHandler func(ctx context.Context, options map[string]any) (string, error)

// Example handler
func (p *MyProvider) MyTool(ctx context.Context, options map[string]any) (string, error) {
    // Extract auth context
    authData, ok := ctx.Value("auth_context").(map[string]any)
    if !ok {
        return "", fmt.Errorf("no auth context")
    }

    userID := authData["user_id"].(string)

    // Use userID for authorization, logging, etc.
    return fmt.Sprintf("Hello, %s!", userID), nil
}
```

---

## OAuth2 Device Flow Authentication

OAuth2 is more complex and requires a provider pattern.

### Architecture

OAuth2 support requires:
1. A provider implementing the `OAuth2Provider` interface
2. Additional HTTP endpoints for the OAuth flow
3. Token storage (database or cache)

### Interface

```go
// mcptypes/auth.go
type OAuth2Provider interface {
    GetDeviceCode(ctx context.Context) (DeviceCodeResponse, error)
    ExchangeDeviceCode(ctx context.Context, deviceCode string) (TokenResponse, error)
    RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)
    ValidateToken(ctx context.Context, accessToken string) (bool, error)
}

type DeviceCodeResponse struct {
    DeviceCode      string
    UserCode        string
    VerificationURI string
    ExpiresIn       int
}

type TokenResponse struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int
}
```

### Implementation Steps

#### Step 1: Implement OAuth2Provider

```go
package myauth

import (
    "context"
    "fmt"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
    "database/sql"
)

type MyOAuth2Provider struct {
    db              *sql.DB
    clientID        string
    clientSecret    string
    authServerURL   string
}

func NewOAuth2Provider(db *sql.DB, clientID, clientSecret, authServerURL string) *MyOAuth2Provider {
    return &MyOAuth2Provider{
        db:            db,
        clientID:      clientID,
        clientSecret:  clientSecret,
        authServerURL: authServerURL,
    }
}

func (p *MyOAuth2Provider) GetDeviceCode(ctx context.Context) (mcptypes.DeviceCodeResponse, error) {
    // Call your OAuth2 authorization server
    // POST /oauth/device/code

    // Example implementation
    resp, err := http.Post(
        p.authServerURL+"/oauth/device/code",
        "application/json",
        strings.NewReader(fmt.Sprintf(`{"client_id":"%s"}`, p.clientID)),
    )
    if err != nil {
        return mcptypes.DeviceCodeResponse{}, err
    }
    defer resp.Body.Close()

    var result mcptypes.DeviceCodeResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return mcptypes.DeviceCodeResponse{}, err
    }

    return result, nil
}

func (p *MyOAuth2Provider) ExchangeDeviceCode(ctx context.Context, deviceCode string) (mcptypes.TokenResponse, error) {
    // Poll authorization server
    // POST /oauth/token with grant_type=urn:ietf:params:oauth:grant-type:device_code

    // Your implementation here
    // ...

    return mcptypes.TokenResponse{}, nil
}

func (p *MyOAuth2Provider) RefreshToken(ctx context.Context, refreshToken string) (mcptypes.TokenResponse, error) {
    // POST /oauth/token with grant_type=refresh_token

    // Your implementation here
    // ...

    return mcptypes.TokenResponse{}, nil
}

func (p *MyOAuth2Provider) ValidateToken(ctx context.Context, accessToken string) (bool, error) {
    // Call token introspection endpoint or check local cache/DB

    // Your implementation here
    // ...

    return true, nil
}
```

#### Step 2: Add OAuth2 API Endpoints to mcpserver

Create `mcpserver/oauth_handlers.go`:

```go
package mcpserver

import (
    "encoding/json"
    "net/http"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
)

// OAuth2Handler wraps an OAuth2Provider and provides HTTP endpoints
type OAuth2Handler struct {
    provider mcptypes.OAuth2Provider
    logger   mcptypes.Logger
}

func newOAuth2Handler(provider mcptypes.OAuth2Provider, logger mcptypes.Logger) *OAuth2Handler {
    return &OAuth2Handler{
        provider: provider,
        logger:   logger,
    }
}

func (h *OAuth2Handler) handleDeviceCode(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    resp, err := h.provider.GetDeviceCode(r.Context())
    if err != nil {
        h.logger.Errorf("GetDeviceCode failed: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *OAuth2Handler) handleTokenExchange(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        DeviceCode string `json:"device_code"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    resp, err := h.provider.ExchangeDeviceCode(r.Context(), req.DeviceCode)
    if err != nil {
        h.logger.Errorf("ExchangeDeviceCode failed: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

#### Step 3: Integrate OAuth2 Endpoints into HTTP Server

Modify transport setup to add OAuth2 routes when OAuth2 is enabled:

```go
// In mcpserver.go, modify Start() for HTTP/SSE modes

case TransportHTTP:
    m.ctx, m.cancel = context.WithCancel(context.Background())
    m.wg.Add(1)
    go func() {
        defer m.wg.Done()

        m.httpServer = server.NewStreamableHTTPServer(m.srv)

        // If OAuth2 is enabled, wrap with OAuth endpoints
        if m.oauth2Provider != nil {
            handler := m.wrapWithOAuthEndpoints(m.httpServer)
            // Start with custom handler
            // ... custom HTTP server setup ...
        } else {
            err := m.httpServer.Start(m.listen)
            _ = err
        }
    }()
```

#### Step 4: Use OAuth2 in Your Application

```go
package main

import (
    "database/sql"
    "github.com/PivotLLM/MCPLaunchPad/mcpserver"
    "github.com/PivotLLM/MCPLaunchPad/mcptypes"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "...")
    oauth2Provider := myauth.NewOAuth2Provider(db, "client-id", "secret", "https://auth.example.com")

    srv, _ := mcpserver.New(
        mcpserver.WithTransportHTTP("localhost:8080"),
        mcpserver.WithOAuth2(oauth2Provider, true), // true = add OAuth endpoints
        mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
    )

    srv.Start()
    // ... shutdown handling ...
}
```

---

## Combining Both Authentication Methods

You can use both Bearer tokens and OAuth2 simultaneously:

```go
srv, _ := mcpserver.New(
    mcpserver.WithTransportHTTP("localhost:8080"),
    mcpserver.WithBearerTokenAuth(validateBearerToken),  // For API keys
    mcpserver.WithOAuth2(oauth2Provider, true),          // For user authentication
    mcpserver.WithToolProviders([]mcptypes.ToolProvider{provider}),
)
```

**Middleware Order** (authentication chain):
1. Try OAuth2 token validation first
2. If no OAuth2 token, try Bearer token
3. If both fail, return 401 Unauthorized

---

## Token Storage Recommendations

### For Bearer Tokens
- **Simple**: In-memory map (development only)
- **Production**: Redis or database table with indexes on token column
- **Format**: Securely hashed tokens, never store plain text

### For OAuth2
- **Required Storage**:
  - Access tokens (short-lived, 1-hour typical)
  - Refresh tokens (long-lived, days/months)
  - Token metadata (user_id, scopes, expiry)

- **Recommended Schema**:
```sql
CREATE TABLE oauth_tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL UNIQUE,
    refresh_token TEXT UNIQUE,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    expires_at TIMESTAMP NOT NULL,
    scopes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_access_token (access_token),
    INDEX idx_user (user_id)
);
```

---

## Security Best Practices

1. **Always use HTTPS** for HTTP/SSE transports with authentication
2. **Token expiration**: Implement short-lived access tokens (1 hour) with refresh tokens
3. **Rate limiting**: Add rate limiting to authentication endpoints
4. **Token rotation**: Rotate refresh tokens on each use
5. **Secure storage**: Use secure vaults (HashiCorp Vault, AWS Secrets Manager) for client secrets
6. **Logging**: Log authentication failures for security monitoring
7. **CORS**: Configure appropriate CORS policies for web clients

---

## Example: Complete Bearer Token Implementation

See the full example in `/examples/bearer-auth/` (to be created as a future example).

## Example: Complete OAuth2 Implementation

See the full example in `/examples/oauth2-auth/` (to be created as a future example).

---

## Summary

**Bearer Token**: Simpler, function-based, good for API keys and internal services

**OAuth2**: More complex, provider-based, good for user authentication and third-party integrations

Both integrate seamlessly with the mcpserver middleware architecture. The key is implementing the validator function (Bearer) or provider interface (OAuth2) according to your specific authentication infrastructure.
