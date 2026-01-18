# MCPServer Package Refactoring - High-Level Design

## Overview
Transform MCPLaunchPad from a demo project into a reusable MCP server library with example implementations. The mcpserver package will become a standalone, production-ready library with no dependencies on this project's specific code.

---

## Design Principles

1. **Library-First Architecture** - mcpserver is the primary product; examples demonstrate usage
2. **Zero Breaking Changes** - Design for extensibility from day 1 using JSON Schema alignment
3. **Progressive Complexity** - Simple cases are trivial; advanced features available when needed
4. **Loose Coupling** - Interfaces for all integration points
5. **Transport Flexibility** - Support stdio, SSE, or HTTP (one mode per runtime instance)
6. **Proper MCP Compliance** - Full hint system and parameter type support

---

## Package Structure

```
MCPLaunchPad/
├── mcptypes/               # NEW: Shared MCP type definitions (no dependencies)
│   ├── logger.go           # Logger interface
│   ├── providers.go        # ToolProvider, ResourceProvider, PromptProvider
│   ├── parameters.go       # Parameter struct + helper constructors
│   ├── hints.go            # ToolHints struct
│   └── auth.go             # Future: BearerTokenValidator, OAuth2Provider (stubs)
│
├── mcpserver/              # REFACTORED: Standalone MCP server library
│   ├── mcpserver.go        # Main server with transport selection
│   ├── options.go          # Configuration options (With* functions)
│   ├── noop_logger.go      # No-op logger implementation
│   ├── hints.go            # Hint builder/helper functions
│   ├── tools.go            # Tool registration logic
│   ├── resources.go        # Resource registration logic
│   ├── prompts.go          # Prompt registration logic
│   ├── hooks.go            # MCP hooks
│   ├── transport.go        # Transport abstraction
│   └── README.md           # Library documentation
│
├── mlogger/                # UPDATED: Implements interfaces.Logger
│   └── mlogger.go
│
├── examples/               # NEW: Example implementations
│   ├── basic/              # Simple stdio example
│   │   ├── main.go
│   │   └── README.md
│   ├── http-api/           # HTTP/SSE example (current main.go + example1/2)
│   │   ├── main.go
│   │   ├── example1/
│   │   └── example2/
│   └── README.md           # Examples overview
│
├── global/                 # DEPRECATED: To be removed after refactor
│
├── go.mod
├── go.sum
├── README.md              # Updated: Focus on library usage
├── CLAUDE.md              # Updated: New architecture
└── LICENSE
```

---

## 1. Transport Architecture

### Design
- **Single transport mode** selected at initialization via option
- Three supported modes: stdio, SSE, HTTP
- Transport abstraction allows future extensions

### Implementation Approach

```go
type TransportMode int

const (
    TransportStdio TransportMode = iota
    TransportSSE
    TransportHTTP
)

// Options
func WithTransportStdio() Option
func WithTransportSSE(listen string) Option
func WithTransportHTTP(listen string) Option
```

### Mode Selection
- Exactly one transport option must be provided
- Returns error if none or multiple transports specified
- Default: None (require explicit selection)

---

## 2. Tool Hint System

### MCP Hints (4 total)
From mcp-go library:
1. **ReadOnlyHint** - Tool is read-only
2. **DestructiveHint** - Tool may perform destructive updates
3. **IdempotentHint** - Repeated calls with same args have no additional effect
4. **OpenWorldHint** - Tool interacts with external entities

### Three-Level Configuration

**Level 1: Package Defaults** (hardcoded fallbacks)
```go
// Internal defaults when nothing specified
defaultReadOnly    = false
defaultDestructive = false
defaultIdempotent  = false
defaultOpenWorld   = false
```

**Level 2: Server-Wide Configuration** (via options)
```go
// Override package defaults for all tools
func WithDefaultReadOnlyHint(value bool) Option
func WithDefaultDestructiveHint(value bool) Option
func WithDefaultIdempotentHint(value bool) Option
func WithDefaultOpenWorldHint(value bool) Option
```

**Level 3: Tool-Level Overrides** (via ToolDefinition)
```go
type ToolHints struct {
    ReadOnlyHint    *bool  // nil = inherit from Level 2/1
    DestructiveHint *bool
    IdempotentHint  *bool
    OpenWorldHint   *bool
}

type ToolDefinition struct {
    Name        string
    Description string
    Parameters  []Parameter
    Handler     ToolHandler
    Hints       *ToolHints  // nil = use all defaults
}
```

### Builder Pattern (Both Styles)

**Style A: Method Chaining**
```go
hints := mcpserver.NewHints().
    ReadOnly(true).
    Destructive(false).
    OpenWorld(true)
```

**Style B: Variadic Constructor**
```go
hints := mcpserver.NewHints(
    mcpserver.ReadOnly(true),
    mcpserver.Destructive(false),
    mcpserver.OpenWorld(true),
)
```

Both return `*ToolHints` ready to assign to `ToolDefinition.Hints`

### Resolution Logic
For each hint:
1. Check tool-level `ToolDefinition.Hints.{Hint}` - if non-nil, use it
2. Else check server-wide config option - if set, use it
3. Else use package default

---

## 3. Logger Design

### Interface Location
**mcptypes/logger.go** - Standalone, no dependencies

```go
package mcptypes

type Logger interface {
    Debug(string)
    Info(string)
    Notice(string)
    Warning(string)
    Error(string)
    Fatal(string)
    Debugf(string, ...any)
    Infof(string, ...any)
    Noticef(string, ...any)
    Warningf(string, ...any)
    Errorf(string, ...any)
    Fatalf(string, ...any)
    Close()
}
```

### No-Op Implementation
**mcpserver/noop_logger.go** - Used when no logger provided

```go
type noopLogger struct{}

func (n *noopLogger) Debug(string) {}
func (n *noopLogger) Info(string) {}
// ... all methods are no-ops
```

### Integration
```go
func WithLogger(logger mcptypes.Logger) Option

// In New():
if m.logger == nil {
    m.logger = &noopLogger{}  // Silent operation
}
```

### mlogger Update
Update `mlogger/mlogger.go` to implement `mcptypes.Logger` (currently implements `global.Logger`)

---

## 4. Provider Interfaces

### Location
**mcptypes/providers.go** - Replaces `global/interfaces.go`

### Updated Interfaces

```go
type ToolProvider interface {
    RegisterTools() []ToolDefinition
}

type ResourceProvider interface {
    RegisterResources() []ResourceDefinition
    RegisterResourceTemplates() []ResourceTemplateDefinition
}

type PromptProvider interface {
    RegisterPrompts() []PromptDefinition
}
```

### Tool Definition

```go
type ToolDefinition struct {
    Name        string
    Description string
    Parameters  []Parameter
    Handler     ToolHandler
    Hints       *ToolHints  // Optional hint overrides
}

type ToolHandler func(options map[string]any) (string, error)
```

---

## 5. Parameter System

### Full JSON Schema Support
**mcptypes/parameters.go**

```go
type Parameter struct {
    Name        string
    Description string
    Required    bool
    Type        string  // "string", "number", "integer", "boolean", "array", "object", "null"

    // String validation
    Pattern   *string
    MinLength *int
    MaxLength *int
    Format    *string  // "date-time", "email", "uri", etc.

    // Numeric validation
    Minimum          *float64
    Maximum          *float64
    ExclusiveMinimum *bool
    ExclusiveMaximum *bool
    MultipleOf       *float64

    // Array validation
    Items       *Parameter
    MinItems    *int
    MaxItems    *int
    UniqueItems *bool

    // Object validation
    Properties           map[string]*Parameter
    AdditionalProperties *bool

    // Enum constraint
    Enum []interface{}

    // Default value
    Default interface{}
}
```

### Helper Constructors (Simple Cases)

```go
// One-liners for common cases
func StringParam(name, description string, required bool) *Parameter
func NumberParam(name, description string, required bool) *Parameter
func IntegerParam(name, description string, required bool) *Parameter
func BoolParam(name, description string, required bool) *Parameter
func ArrayParam(name, description string, required bool, itemType *Parameter) *Parameter
func ObjectParam(name, description string, required bool, properties map[string]*Parameter) *Parameter
```

### Fluent API for Validation

```go
// Method chaining for validation rules
func (p *Parameter) WithPattern(pattern string) *Parameter
func (p *Parameter) WithFormat(format string) *Parameter
func (p *Parameter) WithMinLength(min int) *Parameter
func (p *Parameter) WithMaxLength(max int) *Parameter
func (p *Parameter) WithMinimum(min float64) *Parameter
func (p *Parameter) WithMaximum(max float64) *Parameter
func (p *Parameter) WithEnum(values ...interface{}) *Parameter
func (p *Parameter) WithDefault(value interface{}) *Parameter
// ... etc
```

### Usage Examples

```go
// Simple
params := []mcptypes.Parameter{
    mcptypes.StringParam("name", "User name", true),
    mcptypes.NumberParam("age", "User age", false),
}

// With validation
mcptypes.StringParam("email", "Email address", true).
    WithFormat("email").
    WithPattern(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

mcptypes.NumberParam("port", "Port number", true).
    WithMinimum(1).
    WithMaximum(65535)

// Complex
&mcptypes.Parameter{
    Name: "config",
    Type: "object",
    Properties: map[string]*mcptypes.Parameter{
        "host": mcptypes.StringParam("host", "Host", true),
        "port": mcptypes.IntegerParam("port", "Port", true).WithMinimum(1),
    },
}
```

### Conversion to mcp-go
**mcpserver/tools.go** converts `mcptypes.Parameter` to `mcp.ToolOption` using appropriate type methods:
- `mcp.WithString()` for string types
- `mcp.WithNumber()` for number types
- `mcp.WithBoolean()` for boolean types
- `mcp.WithArray()` for array types
- `mcp.WithObject()` for object types

Validation properties map to `mcp.PropertyOption` equivalents.

---

## 6. Authentication Hooks (Design Only)

### Bearer Token Authentication

**mcptypes/auth.go** (stub for future)
```go
// BearerTokenValidator validates bearer tokens and returns tenant/user context
type BearerTokenValidator func(token string) (context map[string]interface{}, err error)
```

**mcpserver/options.go**
```go
func WithBearerTokenAuth(validator BearerTokenValidator) Option
```

**Integration Point**
- HTTP/SSE transports: Middleware extracts `Authorization: Bearer <token>` header
- Calls validator function (user-provided via option)
- Returns 401 if validation fails
- Injects context into request for handler access

### OAuth2 Authentication

**mcptypes/auth.go** (stub for future)
```go
type OAuth2Provider interface {
    // GetDeviceCode initiates device flow
    GetDeviceCode(ctx context.Context) (DeviceCodeResponse, error)

    // ExchangeDeviceCode polls for token
    ExchangeDeviceCode(ctx context.Context, deviceCode string) (TokenResponse, error)

    // RefreshToken refreshes an access token
    RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)

    // ValidateToken checks token validity
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

**mcpserver/options.go**
```go
func WithOAuth2(provider OAuth2Provider, endpoints bool) Option
// endpoints: if true, add OAuth2 API endpoints to server
```

**Integration Point**
- If enabled, add `/api/v1/oauth/device-code` and `/api/v1/oauth/token` endpoints
- Store tokens per tenant (requires database/state management - user responsibility)
- Middleware validates OAuth tokens for MCP requests
- Can coexist with bearer token auth (check OAuth first, fall back to bearer)

### Design Notes
- Authentication is **optional** - server works without it
- **Interfaces for loose coupling** - users provide implementations
- **Separate files** when implemented: `mcpserver/bearer_auth.go`, `mcpserver/oauth2_auth.go`
- **Separate packages** also viable: `mcpserver/auth/bearer`, `mcpserver/auth/oauth2`
- Both auth types can be enabled simultaneously for different use cases

---

## 7. Context-Aware Handlers

### Current Handler Signature
```go
type ToolHandler func(options map[string]any) (string, error)
```

### Future Enhancement
```go
type ContextAwareToolHandler func(ctx context.Context, options map[string]any) (string, error)
```

### Implementation Strategy
- Keep current signature for backwards compatibility
- Add optional context-aware handler to `ToolDefinition`:

```go
type ToolDefinition struct {
    Name        string
    Description string
    Parameters  []Parameter
    Handler     ToolHandler              // Simple handler
    ContextHandler ContextAwareToolHandler  // Optional context-aware handler
    Hints       *ToolHints
}
```

- Registration logic checks `ContextHandler` first, falls back to `Handler`
- MCP context (tenant info, auth context) passed through Go context

### Migration Path
- Existing code continues working (no breaking change)
- New code can use `ContextHandler` when needed
- Context can carry:
  - Tenant ID (from auth)
  - User info (from OAuth2)
  - Request correlation ID
  - Timeout/cancellation signals

---

## 8. Error Handling & Bug Fixes

### util.go Bug Fix
**Current code (mcpserver/util.go:11-14):**
```go
func (s *MCPServer) logInJSON(data any) {
    b, err := json.MarshalIndent(data, "", "  ")
    if err == nil {  // BUG: Should be err != nil
        s.logger.Debugf("Failed to marshal type %T to JSON: %v", data, data)
        return
    }
    s.logger.Debugf("JSON DATA:\n%s", string(b))
}
```

**Fix:** Change `if err == nil` to `if err != nil`

---

## 9. Configuration Options Summary

### Transport Selection (required, exactly one)
```go
WithTransportStdio()
WithTransportSSE(listen string)
WithTransportHTTP(listen string)
```

### Basic Configuration
```go
WithLogger(logger interfaces.Logger)
WithDebug(debug bool)
WithName(name string)
WithVersion(version string)
```

### Provider Registration
```go
WithToolProviders(providers []mcptypes.ToolProvider)
WithResourceProviders(providers []mcptypes.ResourceProvider)
WithPromptProviders(providers []mcptypes.PromptProvider)
```

### Hint Defaults
```go
WithDefaultReadOnlyHint(value bool)
WithDefaultDestructiveHint(value bool)
WithDefaultIdempotentHint(value bool)
WithDefaultOpenWorldHint(value bool)
```

### Future Authentication (not implemented initially)
```go
WithBearerTokenAuth(validator BearerTokenValidator)
WithOAuth2(provider OAuth2Provider, endpoints bool)
```

---

## 10. Migration Strategy

### Phase 1: New mcptypes Package
1. Create `mcptypes/` package with all interface definitions
2. No breaking changes yet - parallel to `global/`

### Phase 2: Update mcpserver
1. Refactor `mcpserver/` to use `mcptypes/` instead of `global/`
2. Add hint system, parameter helpers, noop logger
3. Add transport selection logic
4. Fix util.go bug

### Phase 3: Update mlogger
1. Change import from `global` to `mcptypes`
2. Verify interface compliance

### Phase 4: Create Examples
1. Move `main.go` to `examples/http-api/main.go`
2. Move `example1/`, `example2/` to `examples/http-api/`
3. Create `examples/basic/` with simple stdio example
4. Update imports in example code

### Phase 5: Documentation
1. Update README.md to focus on library usage
2. Add mcpserver/README.md with comprehensive documentation
3. Add examples/README.md with example descriptions
4. Update CLAUDE.md with new architecture

### Phase 6: Deprecation
1. Add deprecation notice to `global/` package
2. Eventually remove `global/` after confirming no external usage

---

## 11. Documentation Requirements

### mcpserver/README.md
- Installation instructions
- Quick start guide
- Transport mode selection guide
- Hint system explanation
- Parameter system with examples
- Provider interface implementation guide
- Full API reference

### examples/README.md
- Overview of all examples
- When to use each example
- Running instructions

### Individual Example READMEs
- Purpose of the example
- How to run it
- What it demonstrates
- Code walkthrough

---

## 12. Testing Strategy

### Unit Tests
- Hint resolution logic (all three levels)
- Parameter helper constructors
- Transport selection validation
- Logger no-op implementation
- Provider registration

### Integration Tests
- End-to-end stdio server
- End-to-end SSE server
- End-to-end HTTP server
- Hint propagation to mcp-go
- Parameter type conversion

### Example Tests
- Each example should have a test script
- Verify examples compile and run
- Basic smoke tests

---

## 13. Backwards Compatibility

### Breaking Changes
- Move from `global` to `interfaces` package (import path change)
- Transport mode must be explicitly specified (was defaulted before)

### Non-Breaking Enhancements
- Parameter system is additive (existing code works with simple params)
- Hint system is optional (nil = use defaults)
- Logger is optional (nil = no-op)
- Context handlers are optional (regular handlers still work)

### Migration Guide
Document for existing users:
1. Change imports from `github.com/PivotLLM/MCPLaunchPad/global` to `github.com/PivotLLM/MCPLaunchPad/mcptypes`
2. Add explicit transport option to `mcpserver.New()`
3. Optional: Adopt new parameter helpers for cleaner code
4. Optional: Add hints to tool definitions

---

## Success Criteria

✅ mcpserver package has zero dependencies on project-specific code
✅ All four MCP hints supported with three-level configuration
✅ Full JSON Schema parameter support with simple helper API
✅ Three transport modes (stdio, SSE, HTTP) work correctly
✅ Logger interface abstraction with no-op fallback
✅ Clean example code demonstrating library usage
✅ Comprehensive documentation
✅ Authentication hooks designed (stubs in place)
✅ No breaking changes to parameter/hint system in future
✅ mlogger implements new Logger interface
✅ All existing functionality preserved

---

## Next Steps

Once this high-level design is approved:
1. Create detailed implementation plan with file-by-file changes
2. Implement Phase 1 (mcptypes package)
3. Implement Phase 2 (mcpserver refactor)
4. Continue through remaining phases
5. Test thoroughly
6. Update documentation
7. Release as v1.0.0 of the library
