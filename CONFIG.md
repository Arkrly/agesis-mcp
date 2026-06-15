# Aegis-MCP Configuration

## Environment Variables

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `AEGIS_LISTEN_ADDR` | No | `:8080` | HTTP listen address |
| `AEGIS_UPSTREAM_URL` | Yes | None | Upstream MCP endpoint URL |
| `AEGIS_UPSTREAM_TIMEOUT` | No | `30s` | Upstream request timeout |
| `AEGIS_READ_TIMEOUT` | No | `10s` | Max time for inbound request reads |
| `AEGIS_READ_HEADER_TIMEOUT` | No | `5s` | Max time for inbound header reads |
| `AEGIS_WRITE_TIMEOUT` | No | `15s` | Max time for response writes |
| `AEGIS_IDLE_TIMEOUT` | No | `60s` | Idle keepalive timeout |
| `AEGIS_DECISION_TIMEOUT` | No | `5s` | Budget for auth, semantic inspection, and policy evaluation |
| `AEGIS_MAX_BODY_BYTES` | No | `1048576` | Max accepted request size |
| `AEGIS_METRICS_NAMESPACE` | No | `aegis_mcp` | Prefix for emitted metrics |
| `AEGIS_JWT_SHARED_SECRET` | Yes | None | HS256 secret for incoming bearer tokens |
| `AEGIS_JWT_ISSUER` | No | None | Expected JWT issuer |
| `AEGIS_JWT_AUDIENCE` | No | None | Comma-separated acceptable audiences |
| `AEGIS_JWT_REQUIRE_AUDIENCE` | No | `false` | Require an audience match |
| `AEGIS_JWT_ALLOWED_CLOCK_SKEW` | No | `30s` | Allowed JWT timestamp skew |
| `AEGIS_POLICY_FILE` | No | `config/policy.rego` | Rego policy document path |
| `AEGIS_POLICY_RELOAD_INTERVAL` | No | `5s` | Minimum interval between policy reload checks |
| `AEGIS_SEMANTIC_PROVIDER` | No | `heuristic` | Semantic analysis provider (`heuristic`) |
| `AEGIS_SEMANTIC_FAIL_CLOSED` | No | `true` | Deny requests when inspection fails |
| `AEGIS_SEMANTIC_MINIMUM_SCORE` | No | `0.35` | Local semantic minimum score threshold |
| `AEGIS_SEMANTIC_BLOCKED_INTENTS` | No | built-in defaults | Comma-separated blocked intent categories |
| `AEGIS_AUDIT_FILE` | No | `config/audit.db` | BoltDB audit log path |

## Notes

- `AEGIS_UPSTREAM_URL` should point to the MCP server endpoint that Aegis-MCP proxies to.
- The current backend validates MCP JSON-RPC request shape before forwarding, preserves the `Authorization` and `Mcp-Session-Id` headers, and exposes `/live`, `/ready`, `/metrics`.
- Policies are enforced using **Open Policy Agent (OPA)**. It loads `policy.rego` (rules) and optionally `policy.json` (data) from the same directory as the configured policy file.
- Semantic inspection supports a `heuristic` provider initially. The architecture is ready for local LLM integration (Phase 4).
- A durable audit log is maintained in **BoltDB** at the configured `AEGIS_AUDIT_FILE`.
- Per-agent **Rate Limiting** is enforced (default 10 req/s, burst 20).
