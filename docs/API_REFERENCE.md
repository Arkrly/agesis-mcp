# API Reference

Aegis-MCP exposes a transparent proxy endpoint for MCP traffic and several management/monitoring endpoints.

## MCP Proxy Endpoint

### `POST /mcp`

Proxies JSON-RPC 2.0 messages to the upstream MCP server after security inspection.

**Headers:**
- `Authorization`: `Bearer <JWT_TOKEN>` (Required. Must be HS256 signed).
- `Mcp-Session-Id`: `<string>` (Optional. Forwarded for correlation).
- `Content-Type`: `application/json`

**Request Body:**
Standard MCP JSON-RPC 2.0 Request.

**Response Codes:**
- `200 OK`: Request allowed and upstream response returned.
- `400 Bad Request`: Invalid JSON or MCP format.
- `401 Unauthorized`: Missing or invalid JWT.
- `403 Forbidden`: Blocked by RBAC/Policy.
- `413 Payload Too Large`: Request exceeds configured body limit.
- `422 Unprocessable Entity`: Blocked by Semantic Inspection (LLM).
- `429 Too Many Requests`: Rate limit exceeded for the Agent ID.
- `502 Bad Gateway`: Upstream MCP server unreachable.

---

## Management & Monitoring

### `GET /api/summary`

Returns a high-level health and status summary for the dashboard.

**Response (200 OK):**
```json
{
  "status": "ok",
  "version": "v0.1.0",
  "timestamp": "2026-06-15T12:00:00Z",
  "components": {
    "proxy": "healthy",
    "policy": "healthy",
    "audit_log": "healthy"
  }
}
```

### `GET /api/audit`

Retrieves the most recent security decisions from the persistent BoltDB audit log.

**Query Params:**
- `limit`: (Optional) Max entries to return. Default: 100.

**Response (200 OK):**
```json
[
  {
    "ts": "2026-06-15T12:30:45Z",
    "agent_id": "agent-007",
    "method": "tools/call",
    "tool": "read_file",
    "allowed": false,
    "reason": "policy_denied: restricted_directory",
    "inspection": {
      "safety_score": 0.9,
      "intent_categories": ["file_access"]
    }
  }
]
```

### `GET /metrics`

Standard Prometheus metrics (Go defaults + Aegis specific metrics).

**Key Metrics:**
- `aegis_mcp_requests_total`: Counter of all requests.
- `aegis_mcp_request_duration_seconds`: Histogram of processing latency.
- `aegis_mcp_denials_total`: Counter of blocked requests categorized by reason.

### `GET /ready` & `GET /live`

Standard health probes for Kubernetes/Container orchestration.
