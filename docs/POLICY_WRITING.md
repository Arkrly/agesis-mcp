# Policy Writing Guide

Aegis-MCP uses [Open Policy Agent (OPA)](https://www.openpolicyagent.org/) and the Rego language to define security policies.

## Input Schema

The following JSON object is provided to your Rego policies for every request:

```json
{
  "request": {
    "method": "tools/call",
    "params": {
      "name": "read_file",
      "arguments": { "path": "/etc/passwd" }
    },
    "id": 1,
    "jsonrpc": "2.0"
  },
  "agent": {
    "id": "agent-123",
    "roles": ["developer"],
    "metadata": {}
  },
  "semantic_inspection": {
    "safety_score": 0.85,
    "intent_categories": ["file_system", "sensitive_data"]
  }
}
```

## Policy Structure

Policies are located in `config/policy.rego`. Aegis-MCP watches this file and hot-reloads changes instantly.

### Package Names

Aegis-MCP looks for decisions in specific packages:
- `aegis.mcp.auth`: Primary RBAC rules.
- `aegis.mcp.semantic`: Thresholds for LLM safety scores.
- `aegis.mcp.decision`: The final combiner (usually standard).

### Example: Basic RBAC

```rego
package aegis.mcp.auth

default allow = false

# Allow all agents to list tools
allow {
    input.request.method == "tools/list"
}

# Only 'admin' role can call 'delete_file'
allow {
    input.request.method == "tools/call"
    input.request.params.name == "delete_file"
    input.agent.roles[_] == "admin"
}
```

### Example: Semantic Guardrails

```rego
package aegis.mcp.semantic

default allow = true

# Block if the LLM safety score is too low
allow = false {
    input.semantic_inspection.safety_score < 0.4
}

# Block specific dangerous intents
allow = false {
    dangerous := {"prompt_injection", "jailbreak"}
    intent := input.semantic_inspection.intent_categories[_]
    dangerous[intent]
}
```

## Testing Policies

We recommend using the `opa` CLI to test your rules locally before deploying:

```bash
opa test config/policy.rego
```
