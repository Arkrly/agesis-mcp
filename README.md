# Aegis-MCP

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![CI](https://github.com/yourorg/aegis-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/yourorg/aegis-mcp/actions)

**Aegis-MCP** is a zero-trust security gateway for AI agents communicating over the Model Context Protocol (MCP). 

It acts as an intelligent HTTP proxy that intercepts MCP traffic, validates JSON-RPC payloads, authenticates callers via JWT, enforces fine-grained RBAC using Open Policy Agent (OPA), and performs semantic inspection to block malicious prompts (like prompt injection, jailbreaks, and secret exfiltration).

> [!IMPORTANT]
> **Never trust, always verify.** Every single MCP request passing through Aegis-MCP is inspected, evaluated, and explicitly authorized before being forwarded to your backend MCP servers.

## Features

- **Zero-Trust Architecture**: Deny-by-default routing ensures only explicitly authorized traffic reaches your MCP servers.
- **Fine-Grained Authorization**: Integrated Open Policy Agent (OPA) evaluates Rego policies against the request context, JWT claims, and semantic intents.
- **Semantic Prompt Inspection**: Built-in heuristic engine identifies and blocks prompt injection, malicious code generation, and secret exfiltration attempts. (Future support planned for local quantized LLMs).
- **Strong Authentication**: Requires standard JWTs (HS256) with strict validation of claims (`exp`, `nbf`, `aud`, `iss`, and custom `agent_id` or `roles`).
- **Resilience & Protection**: Protects upstream services with per-agent rate limiting (token bucket) and configurable request body size limits.
- **Observability**: Exposes Prometheus metrics, structured JSON logging, and maintains a durable, persistent audit log of all decisions via BoltDB.

## Request Lifecycle

1. **Ingress**: Receives HTTP POST to `/mcp`. Enforces content-type and payload size limits.
2. **Decoding**: Validates the JSON-RPC 2.0 structure and extracts the targeted MCP method and tool.
3. **Authentication**: Extracts and cryptographically verifies the bearer token (JWT).
4. **Semantic Inspection**: Analyzes the prompt text to score safety and categorize intent.
5. **Policy Evaluation**: OPA evaluates the `policy.rego` rules using the agent's identity, roles, tool requested, and semantic inspection results.
6. **Decision & Audit**: The decision is durably logged to a local BoltDB database.
7. **Egress**: If fully authorized, the proxy forwards the payload to the configured upstream MCP server, preserving necessary session headers.

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.26.4 or higher.

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/yourorg/aegis-mcp.git
cd aegis-mcp
go build -o aegis-mcp ./cmd/aegis-mcp
```

### Configuration

Aegis-MCP is primarily configured via environment variables.

1. Copy the example environment file:
   ```bash
   cp config/aegis.env.example .env
   ```

2. Edit `.env` to set your upstream URL and JWT secret:
   ```ini
   AEGIS_LISTEN_ADDR=:8080
   AEGIS_UPSTREAM_URL=http://127.0.0.1:9090/mcp
   AEGIS_JWT_SHARED_SECRET=your-super-secret-key
   AEGIS_POLICY_FILE=config/policy.rego
   ```

> [!NOTE]
> See [CONFIG.md](./CONFIG.md) for a complete list of all supported environment variables, timeouts, and metrics configurations.

### Running the Gateway

Load your environment variables and start the server:

```bash
set -a
source .env
set +a

./aegis-mcp
```

The gateway will start and bind to the configured `AEGIS_LISTEN_ADDR`.

## Usage

Once running, point your AI Agents to Aegis-MCP instead of directly to your MCP servers.

### Endpoints

- `POST /mcp`
  The primary proxy endpoint. Requires a valid `Authorization: Bearer <token>` header. Evaluates policies and forwards traffic to the upstream URL.
  
- `GET /metrics`
  Prometheus-compatible metrics exposing request counters, latencies, and detailed denial reasons.

- `GET /live`
  Liveness probe. Returns `200 OK` if the HTTP server is responsive.

- `GET /ready`
  Readiness probe. Returns `200 OK` if all dependencies (like the OPA policy engine and semantic inspector) have loaded successfully. Returns `503 Service Unavailable` otherwise.

## Policies

Aegis-MCP uses standard Rego files for authorization. A default `policy.rego` is provided in the `config/` directory.

By default, the engine operates in **fail-closed** mode. You must explicitly define `allow` rules based on methods, tools, roles, and agent IDs. The policy engine automatically hot-reloads when the `policy.rego` file is modified on disk.

## Development

To run both the backend (Go) and the React dashboard together in development mode:

1. **Install Dependencies**:
   ```bash
   npm install
   ```

2. **Start Services**:
   ```bash
   npm run dev
   ```

This launches:
- **Backend API**: `http://localhost:8080` (requires environment variables from `.env`)
- **Frontend Dashboard**: `http://localhost:5173`

Individual services can also be run via:
- `npm run backend`
- `npm run frontend`
