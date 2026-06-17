# Aegis-MCP Backend Architecture

## System Overview
Aegis-MCP functions as a security-first HTTP proxy that sits between MCP clients and servers. It enforces dual-layer security: **protocol-level RBAC** (via OPA) and **content-level semantic guardrails** (via heuristic/LLM inspection). The system follows a modular pipeline architecture where each stage processes the request sequentially to ensure zero-trust verification of every interaction.

## Component Diagram (ASCII)
```
MCP Client ────────────┐
                       ▼
           ┌─────────────────┐
           │   TLS Terminator│◀───────┐
           │ (optional)      │        │
           └─────────────────┘        │
                       ▼              │
           ┌─────────────────┐        │
           │   Request Parser │        │
           │  (MCP → HTTP)    │        │
           └─────────────────┘        │
                       ▼              │
           ┌─────────────────┐        │
           │   JSON Validator│        │
           │   (Schema)      │        │
           └─────────────────┘        │
                       ▼              │
           ┌─────────────────┐        │
           │ Semantic        │        │
           │   Inspector     │◀───────┘
           │   (Local LLM)   │
           └─────────────────┘
                       ▼
           ┌─────────────────┐
           │    OPA Engine   │
           │   (Policy Eval) │
           └─────────────────┘
                       ▼
           ┌─────────────────┐
           │  Decision Cache │
           │   (LRU, 5min TTL)│
           └─────────────────┘
                       ▼
           ┌─────────────────┐
           │   Response      │
           │   Formatter     │
           │  (Allow/Block)  │
           └─────────────────┘
                       ▼
           ┌─────────────────┐
           │   MCP Encoder   │
           │  (HTTP → MCP)   │
           └─────────────────┘
                       ▼
               MCP Server
```

## Data Flow
1. **Ingress**: MCP client connects via HTTP (with optional TLS termination)
2. **Parsing**: Raw MCP messages parsed from HTTP payload
3. **Validation**: JSON schema validation against MCP specification
4. **Semantic Inspection**: Prompts sent to local LLM for safety/intent analysis
5. **Policy Evaluation**: OPA evaluates request against RBAC policies
6. **Caching**: Recent decisions cached to reduce LLM/OPA calls
7. **Decision**: Request allowed, blocked, or modified based on evaluation
8. **Egress**: Approved requests forwarded to MCP server; blocked requests return error

## Technology Decisions & Rationale

### Language: Go
- **Why**: Excellent concurrency model for proxy workloads, strong standard library for HTTP/JSON, static binaries for easy deployment, growing ecosystem for security tools
- **Alternatives Considered**: Rust (steeper learning curve), Java (higher memory overhead), Node.js (single-threaded limitations)

### Local LLM for Semantic Inspection
- **Why**: Privacy-sensitive (no prompt leakage), predictable latency/cost, offline operation, compliance-friendly
- **Model Choice**: Quantized Llama 3 8B or Mistral 7B via llama.cpp or similar
- **Alternatives Considered**: API-based LLMs (rejected due to privacy concerns and variable latency)

### OPA for Policy Engine
- **Why**: Industry-standard for cloud-native authorization, declarative policy language (Rego), excellent performance, strong ecosystem
- **Alternatives Considered**: Custom policy engine (rejected due to complexity and maintenance burden)

### Embedded Database: BoltDB
- **Why**: Zero-configuration, ACID transactions, excellent read performance, embedded in Go binary
- **Use Cases**: Policy metadata, audit logs, decision cache persistence
- **Alternatives Considered**: SQLite (similar but BoltDB has better Go integration), Redis (adds external dependency)

### JWT Authentication
- **Why**: Stateless, widely adopted, easy to integrate with existing identity providers
- **Implementation**: Validate JWT signatures, extract claims for OPA policy input
- **Alternatives Considered**: API keys (less secure), OAuth 2.0 (overkill for service-to-service)

### Docker Containers
- **Why**: Consistent deployment, easy orchestration, isolation, matches user requirements
- **Image Base**: Distroless or Alpine for minimal attack surface
- **Multi-Stage Build**: Separate build and runtime images for security

### Observability Stack
- **Logging**: Structured JSON logs to stdout (compatible with Docker logging drivers)
- **Metrics**: Prometheus endpoint (/metrics) for request rates, latency, error rates
- **Tracing**: OpenTelemetry instrumentation for end-to-end tracing

## Security Considerations
1. **Privilege Separation**: Run as non-root user in container
2. **Secrets Management**: JWT verification keys mounted as read-only volumes
3. **Network Policies**: Container only exposes HTTP port; no outbound connectivity required
4. **Resource Limits**: CPU/memory limits to prevent noisy neighbor issues
5. **Image Scanning**: Regular vulnerability scanning of base images

## Assumptions
- [ASSUMPTION] TLS termination handled externally (load balancer or sidecar) for simplicity
- [ASSUMPTION] Initial implementation focuses on text-based MCP prompts; multimodal support (images, audio) in v1.2
- [ASSUMPTION] OPA policies loaded at startup with hot-reload capability via file watcher
- [ASSUMPTION] Decision cache size bounded at 10,000 entries to prevent memory exhaustion