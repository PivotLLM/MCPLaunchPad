# OAuth2 Google Authentication Example

This example demonstrates Google OAuth2 authentication for HTTP-based MCP servers using the device flow.

## What It Does

Provides tools protected by Google OAuth2 authentication:
- `get_greeting` - Returns a personalized greeting message
- `get_user_info` - Returns authenticated user information

## Features Demonstrated

- **Google OAuth2 Device Flow**: User authentication via Google
- **HTTP Transport**: Communication over HTTP protocol
- **Bearer Token Validation**: OAuth2 access tokens as bearer tokens
- **Automatic Token Validation**: Built-in token validation with Google
- **User Information**: Access to authenticated user's profile
- **Graceful Shutdown**: Signal handling for clean server shutdown

## Prerequisites

### 1. Google Cloud Project Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" → "Credentials"
4. Click "Create Credentials" → "OAuth client ID"
5. Select "Desktop app" as the application type
6. Note down your Client ID and Client Secret

### 2. OAuth Consent Screen

1. In Google Cloud Console, go to "APIs & Services" → "OAuth consent screen"
2. Select "External" user type (or "Internal" if you have Google Workspace)
3. Fill in the required information:
   - App name
   - User support email
   - Developer contact information
4. Add the following scopes:
   - `email`
   - `profile`
5. Add test users (yourself) under "Test users" section
6. Save and continue

### 3. Set Environment Variables

```bash
export GOOGLE_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="your-client-secret"
```

Or create a `.env` file:
```bash
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
```

## Building

```bash
go build
```

## Running

### With OAuth2 Authentication

```bash
# Set environment variables
export GOOGLE_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="your-client-secret"

# Start server
./oauth

# Follow the on-screen instructions:
# 1. Open the verification URL in your browser
# 2. Enter the user code displayed
# 3. Authorize the application
# 4. Server will automatically receive the access token
```

The server will display:
```
=== OAuth2 Device Flow ===
1. Open this URL in your browser: https://www.google.com/device
2. Enter this code: ABCD-EFGH
3. Waiting for authorization...

OAuth2 authentication configured successfully!

To use this server, include the access token in the Authorization header:
  Authorization: Bearer ya29.a0AfB_by...
```

### Testing with Probe

After authentication, use the access token with probe:

```bash
# List tools (requires access token)
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer YOUR_ACCESS_TOKEN" -list-only

# Call greeting tool
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer YOUR_ACCESS_TOKEN" -call get_greeting -params '{"name":"World"}'

# Get user info
probe -transport http -url http://localhost:8080/mcp -headers "Authorization:Bearer YOUR_ACCESS_TOKEN" -call get_user_info
```

### Skip Authentication (Testing Only)

For testing the server without OAuth2:

```bash
./oauth --skip-auth
```

This starts the server without requiring authentication.

## Command Line Flags

- `--listen` - Address to listen on (default: "localhost:8080")
- `--skip-auth` - Skip OAuth2 authentication for testing

## Authentication Flow

1. **Device Code Request**: Application requests a device code from Google
2. **User Authorization**: User opens verification URL and enters code on their device
3. **Token Polling**: Application automatically polls Google for access token
4. **Token Received**: Once user authorizes, application receives access token
5. **Bearer Token Validation**: Access token is used for subsequent API requests
6. **User Info**: Application can retrieve user profile information

## Code Overview

The example demonstrates:

1. **OAuth2 Provider Setup**: Creating Google OAuth2 provider with credentials
2. **Device Flow Execution**: Performing complete device flow with automatic polling
3. **Bearer Token Validator**: Converting OAuth2 provider to bearer token validator
4. **Server Configuration**: Integrating OAuth2 with MCP server
5. **Tool Implementation**: Simple tools that require authentication

## Security Notes

- **Never commit credentials**: Use environment variables for client ID/secret
- **Token Storage**: This example doesn't persist tokens (they expire after ~1 hour)
- **Production Use**: Implement token storage (database, Redis) for production
- **HTTPS Required**: Always use HTTPS in production
- **Scope Minimization**: Only request scopes you need (email, profile)
- **Token Refresh**: Implement refresh token logic for long-running services

## Troubleshooting

### "GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set"

Make sure you've exported the environment variables before running:
```bash
export GOOGLE_CLIENT_ID="your-client-id"
export GOOGLE_CLIENT_SECRET="your-client-secret"
```

### "OAuth2 device flow failed"

- Check that your credentials are correct
- Verify that the OAuth consent screen is properly configured
- Ensure you're added as a test user if the app isn't published
- Check that the device flow authorization URL is accessible

### "Token validation failed"

- Access tokens expire after ~1 hour
- Re-run the authentication flow to get a new token
- For production, implement refresh token logic

### "401 Unauthorized"

- Ensure you're including the `Authorization: Bearer <token>` header
- Verify the token hasn't expired
- Check that the token is valid with Google's tokeninfo endpoint

## Production Considerations

For production deployment:

1. **Token Storage**: Store access and refresh tokens in a secure database
2. **Token Refresh**: Implement automatic token refresh before expiration
3. **User Session**: Link tokens to user sessions
4. **Multi-User**: Support multiple authenticated users simultaneously
5. **Rate Limiting**: Implement rate limiting for API endpoints
6. **Monitoring**: Log authentication events for security monitoring
7. **HTTPS**: Always use HTTPS for token transmission
8. **App Verification**: Complete Google OAuth app verification process

## Next Steps

- Implement token storage (database, Redis)
- Add refresh token handling
- Support multiple concurrent users
- Add user session management
- Implement token revocation
- Add OAuth2 scope management

## References

- [Google OAuth2 Device Flow](https://developers.google.com/identity/protocols/oauth2/limited-input-device)
- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth2 Package Documentation](../../oauth2/README.md)
- [AUTHENTICATION.md](../../AUTHENTICATION.md)
