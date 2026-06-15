# Getting Started

## Run Locally

1. Copy [config/aegis.env.example](/home/arkrly/Projects/agesis-mcp/config/aegis.env.example) into your shell environment.
2. Ensure `AEGIS_UPSTREAM_URL` points to an MCP-compatible upstream endpoint.
3. Start the server:

```bash
set -a
source config/aegis.env.example
set +a
go run ./cmd/aegis-mcp
```

## Endpoints

- `POST /mcp`: validates JSON-RPC 2.0 MCP requests, authenticates the caller, runs semantic inspection, evaluates policy, and proxies allowed requests upstream.
- `GET /live`: process liveness.
- `GET /ready`: dependency readiness, including policy load health.
- `GET /metrics`: Prometheus-style counters for requests and denials.

## Default Behavior

- JWT authentication is mandatory and uses `HS256`.
- Requests are denied by default unless a matching policy rule allows them.
- Semantic inspection fails closed by default.
- Destructive tools such as `delete_file` and `exec_command` require the `admin` role in the sample policy.
