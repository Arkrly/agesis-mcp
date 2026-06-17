# REQUIREMENTS.md

## v1 Requirements

### Authentication & Identity
- [ ] **AUTH-01**: User can create account with email/password (Deferred to management UI/external IDP)
- [ ] **AUTH-02**: User can log in and stay logged in across sessions (Deferred)
- [ ] **AUTH-03**: User can log out from any page (Deferred)
- [x] **AUTH-04**: System validates JWT tokens for service-to-service authentication
- [x] **AUTH-05**: System extracts agent identity and roles from JWT claims
- [ ] **AUTH-06**: System supports machine-identity standards (SPIFFE/SPIRE equivalent) for AI agents (In progress)
- [ ] **AUTH-07**: System issues short-lived credentials for AI agent identities (Deferred)

### MCP Protocol Handling
- [x] **MCP-01**: System intercepts MCP traffic via HTTP proxy
- [x] **MCP-02**: System validates MCP message format (JSON-RPC 2.0)
- [x] **MCP-03**: System enforces configurable message size limits
- [x] **MCP-04**: System supports standard MCP methods (tools/list, tools/call, resources/list, resources/read, prompts/list, prompts/get)
- [x] **MCP-05**: System maintains stateless horizontal scalability
- [x] **MCP-06**: System adds security headers (X-Content-Type-Options, X-Frame-Options, etc.)

### Semantic Inspection
- [x] **SEM-01**: System performs semantic inspection of prompts using local LLM (Heuristic V1 implemented, Real LLM in progress)
- [x] **SEM-02**: System detects prompt injection attempts (Heuristic)
- [x] **SEM-03**: System identifies jailbreak attempts (Heuristic)
- [x] **SEM-04**: System detects harmful content in prompts (Heuristic)
- [x] **SEM-05**: System provides safety scoring for prompts (0.0-1.0 scale)
- [x] **SEM-06**: System categorizes intent of prompts (Heuristic)

### Policy Enforcement
- [x] **POL-01**: System enforces RBAC policies via OPA for agent authorization
- [x] **POL-02**: System enforces tool authorization policies via OPA
- [x] **POL-03**: System evaluates requests against resource access policies
- [x] **POL-04**: System integrates semantic inspection results with policy evaluation
- [x] **POL-05**: System makes real-time allow/block decisions (<100ms latency target)
- [x] **POL-06**: System provides detailed denial reasons with codes and metadata

### Proxy Functionality
- [x] **PROX-01**: System forwards allowed requests to upstream MCP servers
- [x] **PROX-02**: System returns MCP responses to clients unchanged
- [x] **PROX-03**: System blocks requests that fail security checks with appropriate error responses
- [x] **PROX-04**: System handles MCP notifications (messages without ID field)
- [x] **PROX-05**: System preserves MCP request/response semantics

### Observability & Logging
- [x] **OBS-01**: System maintains audit trail of all requests (allowed and blocked)
- [x] **OBS-02**: System logs agent identity, action, timestamp, and decision for each request
- [x] **OBS-03**: System exports structured logs compatible with log aggregation systems
- [x] **OBS-04**: System provides Prometheus metrics endpoint
- [ ] **OBS-05**: System instruments code with OpenTelemetry for distributed tracing (Pending)
- [x] **OBS-06**: System implements health check endpoints (/ready, /live)

### Storage & Configuration
- [x] **STO-01**: System uses embedded database (BoltDB) for policy metadata and audit logs (MVP)
- [ ] **STO-02**: System supports external database (PostgreSQL) for production scale (v1.1)
- [ ] **STO-03**: System uses Redis for caching and temporary storage (v1.1)
- [x] **STO-04**: System loads configuration from environment variables and config file
- [x] **STO-05**: System supports hot-reloading of non-security-critical configuration
- [ ] **STO-06**: System implements data retention policies with automatic purge (In progress)

### Security & Compliance
- [x] **SEC-01**: System implements defense-in-depth with multiple security layers
- [x] **SEC-02**: System prevents over-permissioning of AI agent identities
- [x] **SEC-03**: System mitigates prompt injection and data leakage risks
- [x] **SEC-04**: System treats AI agents as machine identities (not human users)
- [ ] **SEC-05**: System addresses supply chain risks in agent tooling (In progress)
- [x] **SEC-06**: System ensures fully open-source implementation

### DevOps & Deployment
- [x] **DEV-01**: System provides Dockerfile for containerized deployment
- [x] **DEV-02**: System implements multi-stage Docker build for minimal image size
- [x] **DEV-03**: System configures GitHub Actions CI/CD pipeline
- [x] **DEV-04**: System runs unit tests, integration tests, and security scans in CI
- [ ] **DEV-05**: System builds and pushes Docker image to GitHub Container Registry (Pending)
- [ ] **DEV-06**: System provides Helm chart for Kubernetes deployment (v1.1)
- [ ] **DEV-07**: System provides Kubernetes manifests for deployment (v1.1)

### Out of Scope (v1)

- [Blockchain-based audit trail] — Immutable logs for regulated industries (long-term)
- [AI-powered policy optimization] — Uses ML to suggest policy improvements
- [Predictive threat detection] — ML to identify emerging threats from prompt patterns
- [Policy marketplace] — Community sharing and rating of Rego policies
- [AI-generated MCP server templates] — LLMs to generate configurations from natural language
- [Edge deployment modes] — Lightweight versions for IoT or on-premise isolation
- [Federated policy management] — Central policy distribution with edge enforcement
- [Quantum-resistant cryptography] — (long-term)
