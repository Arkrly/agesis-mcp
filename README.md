# Aegis-MCP: Zero-Trust Security Gateway for AI Agents

Aegis-MCP is a secure proxy for AI agents communicating via the **Model Context Protocol (MCP)**. It provides a zero-trust security layer that intercepts, inspects, and authorizes every request before it reaches your MCP servers.

## 🚀 Quick Start (One Command)

To run the entire stack (Proxy + Dashboard + Mock MCP Server):

```bash
git clone git@github.com:Arkrly/agesis-mcp.git
cd agesis-mcp
npm install
npm run dev
```

- **Gateway/Proxy**: `http://localhost:8080`
- **Management Dashboard**: `http://localhost:5173`
- **Mock MCP Server**: `http://localhost:9090` (for testing)

---

## 🏗️ Architecture

Aegis-MCP sits between your AI Agents (clients) and your MCP Servers.

```text
AI Agent  ──>  Aegis-MCP Gateway  ──>  Policy Engine (OPA)  ──>  Real MCP Server
                (Port 8080)             (Rego Rules)              (Upstream)
```

1.  **Intercept**: Validates MCP JSON-RPC structure.
2.  **Authenticate**: Verifies caller identity via HS256 JWT tokens.
3.  **Inspect**: Performs semantic analysis on prompts to detect injections/jailbreaks.
4.  **Authorize**: Evaluates Open Policy Agent (OPA) rules for RBAC and tool-level access.
5.  **Audit**: Logs every decision to a persistent BoltDB database.

---

## 📊 Management Dashboard

The project includes a modern React 19 dashboard ("HackCulture" theme) to monitor your gateway:
- **Real-time Status**: Monitor proxy and component health.
- **Security Audit**: View a detailed table of allowed and blocked requests.
- **Policy Preview**: Visualize and draft new security rules.

---

## 🛠️ Components

| Component | Path | Description |
| --- | --- | --- |
| **Backend** | `/internal` | Go core implementing the proxy and security logic. |
| **Dashboard** | `/frontend` | React + Tailwind v4 management console. |
| **Mock Server**| `/test/mock-mcp` | Lightweight MCP server for integration testing. |
| **Policies** | `/config` | Rego (OPA) rules and JSON data for authorization. |

---

## 🚀 Deployment

### Railway (Recommended)
Aegis-MCP is optimized for one-click deployment on **Railway**.
1. Connect this GitHub repository to [Railway](https://railway.app/).
2. Railway will automatically detect the `Dockerfile` and build a unified production image.
3. The dashboard will be served directly at your assigned Railway domain.

---

## 📖 Documentation

- **[Getting Started](docs/GETTING_STARTED.md)**: Detailed setup and first-request guide.
- **[API Reference](docs/API_REFERENCE.md)**: Endpoint documentation and status codes.
- **[Policy Writing](docs/POLICY_WRITING.md)**: How to write custom Rego security rules.
- **[Configuration](CONFIG.md)**: Full list of environment variables and settings.

---

## 🛡️ Security Features

- **Semantic Guardrails**: Heuristic detection of malicious intent (Phase 4 LLM integration ready).
- **Fine-Grained RBAC**: Restrict tools (e.g., `read_file` vs `delete_file`) based on agent roles.
- **Rate Limiting**: Protect upstream servers from token exhaustion/abuse.
- **Fail-Closed Design**: If security components fail, the gateway denies all traffic.
