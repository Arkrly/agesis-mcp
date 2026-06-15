# Getting Started

Follow this guide to get Aegis-MCP up and running in your environment.

## 1. Prerequisites

- **Go**: v1.26+ (for the backend)
- **Node.js**: v18+ & **npm** (for the dashboard)
- **Docker**: (Optional, for containerized deployment)

## 2. Fast Track (Monorepo)

The easiest way to start both the backend and the dashboard:

```bash
# Clone the repo
git clone https://github.com/yourorg/aegis-mcp.git
cd aegis-mcp

# Install orchestration dependencies
npm install

# Start everything
npm run dev
```

The Dashboard will be available at `http://localhost:5173`.

## 3. Backend Setup

If you want to run the backend individually:

```bash
# Copy example configuration
cp config/aegis.env.example .env

# Build and run
go build -o aegis-mcp ./cmd/aegis-mcp
set -a && source .env && set +a
./aegis-mcp
```

### Configuration
Key environment variables in `.env`:
- `AEGIS_LISTEN_ADDR`: Port to listen on (default `:8080`).
- `AEGIS_UPSTREAM_URL`: The destination MCP server.
- `AEGIS_JWT_SHARED_SECRET`: Key for verifying agent tokens.

## 4. Dashboard Setup

The dashboard provides a real-time view of security events and health status.

```bash
cd frontend
npm install
npm run dev
```

## 5. First Request

To test the proxy, send a signed MCP request:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Authorization: Bearer <YOUR_JWT>" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

Refer to the [API Reference](./API_REFERENCE.md) for detailed status codes and endpoints.
