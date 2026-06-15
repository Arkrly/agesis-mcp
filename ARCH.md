# Aegis-MCP Design Specification

## API Design

### External HTTP API
Aegis-MCP exposes a single HTTP endpoint that acts as a transparent MCP proxy.

#### Endpoint
```
POST /mcp
```

#### Request Format
- **Method**: POST
- **Content-Type**: application/json (MCP messages are JSON-RPC 2.0)
- **Body**: Valid MCP JSON-RPC request object
- **Headers**: 
  - `Authorization`: Bearer <JWT> (required)
  - `Content-Type`: application/json
  - `Mcp-Session-Id`: <optional session identifier>

#### Response Format
- **Success (200)**: MCP JSON-RPC response object (if request allowed)
- **Client Error (4xx)**: 
  - 400: Invalid MCP request format
  - 401: Missing or invalid JWT token
  - 403: Request blocked by policy (RBAC or semantic inspection)
  - 422: MCP request passed validation but failed semantic/policy checks
- **Server Error (5xx)**: Internal processing errors

#### Example Request
```http
POST /mcp HTTP/1.1
Host: aegis-mcp.local:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
Mcp-Session-Id: sess_abc123

{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 1
}
```

#### Example Success Response
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "read_file",
        "description": "Read contents of a file",
        "inputSchema": {
          "type": "object",
          "properties": {
            "path": { "type": "string" }
          },
          "required": ["path"]
        }
      }
    ]
  },
  "id": 1
}
```

#### Example Blocked Response
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json
X-Aegis-Reason: policy_violation
X-Aegis-Policy: rbac:agent_role

{
  "error": {
    "code": -32003,
    "message": "Request blocked by Aegis-MCP security policy",
    "data": {
      "reason": "policy_violation",
      "policy": "rbac:agent_role",
      "agent_id": "agent_789",
      "required_role": "admin",
      "actual_role": "developer"
    }
  }
}
```

### Internal Component Interfaces

#### LLM Inspector Interface
```go
type SemanticInspector interface {
    // InspectPrompt analyzes a prompt for safety and intent
    // Returns: safety score (0.0-1.0), detected intent categories, error
    InspectPrompt(ctx context.Context, prompt string) (float64, []string, error)
    
    // Close releases any resources
    Close() error
}
```

#### OPA Policy Evaluator Interface
```go
type PolicyEvaluator interface {
    // EvaluateRequest checks if an MCP request is allowed
    // Returns: allowed (bool), reason string if denied, error
    EvaluateRequest(ctx context.Context, req *MCPRequest, claims jwt.MapClaims) (bool, string, error)
    
    // ReloadPolicies hot-loads updated policy definitions
    ReloadPolicies() error
}
```

#### MCP Request/Response Structures
```go
type MCPRequest struct {
    JSONRPC string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
    ID      *json.RawMessage `json:"id,omitempty"` // null for notifications
}

type MCPResponse struct {
    JSONRPC string          `json:"jsonrpc"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *RPCError       `json:"error,omitempty"`
    ID      *json.RawMessage `json:"id,omitempty"`
}

type RPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}
```

## OPA Policy Schema

### Input Schema
OPA receives a JSON input with the following structure:
```json
{
  "request": {
    "method": "string",
    "params": {},
    "id": null,
    "jsonrpc": "2.0"
  },
  "agent": {
    "id": "string",
    "roles": ["string", ...],
    "metadata": {}
  },
  "session": {
    "id": "string",
    "created_at": "timestamp",
    "metadata": {}
  },
  "semantic_inspection": {
    "safety_score": 0.0,
    "intent_categories": ["string", ...],
    "timestamp": "timestamp"
  }
}
```

### Policy Structure (Rego)
#### RBAC Policies
```rego
package aegis.mcp.auth

# Default deny
default allow = false

# Allow if agent has required role for method
allow {
    input.request.method == "tools/list"
    agent_has_role("developer")
}

allow {
    input.request.method == "tools/call"
    input.request.params.tool == "read_file"
    agent_has_role("developer")
}

allow {
    input.request.method == "tools/call"
    input.request.params.tool == "delete_file"
    agent_has_role("admin")
}

# Helper function
agent_has_role(role) {
    r := input.agent.roles[_]
    r == role
}
```

#### Semantic Inspection Policies
```rego
package aegis.mcp.semantic

# Block prompts with low safety scores
violation["low_safety_score"] {
    input.semantic_inspection.safety_score < 0.3
}

# Block specific dangerous intent categories
violation["dangerous_intent"] {
    intent := input.semantic_inspection.intent_categories[_]
    dangerous_intents := ["prompt_injection", "jailbreak_attempt", "malicious_code_generation"]
    dangerous_intents[_] == intent
}

# Allow if no violations
allow {
    not violation[_]
}
```

#### Combined Policy Decision
```rego
package aegis.mcp.decision

import data.aegis.mcp.auth
import data.aegis.mcp.semantic

# Final allow decision requires both auth and semantic to pass
allow {
    auth.allow
    semantic.allow
}

# Provide detailed denial reason
denial_reason[reason] {
    not auth.allow
    reason := concat("; ", auth.denial_reasons)
} else {
    not semantic.allow
    reason := concat("; ", semantic.violation_reasons)
}
```

## MCP Proxy Interface Contracts

### Supported MCP Methods
Aegis-MCP must correctly proxy all standard MCP methods, but may apply different policies based on method type.

#### Standard Methods Requiring Inspection
1. `tools/list` - List available tools (auth only, typically allowed)
2. `tools/call` - Execute a tool (full inspection: auth + semantic + tool-specific)
3. `resources/list` - List resources (auth only)
4. `resources/read` - Read resource content (auth + semantic on params)
5. `prompts/list` - List prompts (auth only)
6. `prompts/get` - Get prompt template (auth + semantic on params)

### Transport Requirements
1. **HTTP/1.1 Only**: Initial implementation targets HTTP/1.1 for simplicity
2. **Keep-Alive Connections**: Support HTTP persistent connections for performance
3. **No HTTP/2**: Deferred to v1.2 to reduce complexity
4. **WebSocket Not Supported**: MCP over HTTP only; WebSocket transport in v2.0

### Message Size Limits
- **Maximum Request Size**: 4MB (configurable)
- **Maximum Response Size**: 16MB (configurable)
- **Oversize Messages**: Return 413 Payload Too Large

### Session Handling
- **Stateless by Design**: No server-side session storage required
- **Optional Session ID**: Mcp-Session-Id header forwarded unchanged for correlation
- **Timeouts**: 
  - Read timeout: 30s
  - Write timeout: 30s
  - Idle timeout: 5m (connection level)

### Error Handling Contracts
1. **MCP-Level Errors**: Returned as valid JSON-RPC error objects
2. **Aegis-MCP Errors**: Returned as HTTP status codes with JSON error bodies
3. **Never Mix**: HTTP errors never contain JSON-RPC bodies and vice versa
4. **Error Chaining**: Internal errors include trace IDs for debugging but never expose stack traces to clients

### Security Headers
All responses include:
```
X-Aegis-Version: v0.1.0
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
```

## Assumptions
- [ASSUMPTION] Initial implementation supports only JSON-RPC 2.0 over HTTP (no batch requests)
- [ASSUMPTION] OPA policies loaded from local filesystem; bundle distribution in v1.1
- [ASSUMPTION] Semantic inspector uses a single LLM instance; pooling in v1.2
- [ASSUMPTION] JWT validation uses HS256 algorithm only; RS256 support in v1.1
- [ASSUMPTION] Rate limiting implemented at infrastructure level (reverse proxy), not in Aegis-MCP