# ROADMAP.md

## Phases

- [ ] **Phase 1: Foundation & Basic Proxy** - Securely proxy MCP traffic with transport validation, basic JWT auth, and core proxy functionality
- [ ] **Phase 2: Identity & Access Management** - Manage identities and enforce fine-grained access control for human admins and AI agents
- [ ] **Phase 3: Semantic Inspection & AI Security** - Intelligently inspect and filter MCP prompts for AI-specific threats using local LLMs
- [ ] **Phase 4: Observability, Analytics & Operations** - Provide comprehensive monitoring, logging, and operational capabilities

## Phase Details

### Phase 1: Foundation & Basic Proxy
**Goal**: Users can securely proxy MCP traffic with basic security controls
**Depends on**: Nothing (first phase)
**Requirements**: MCP-01, MCP-02, MCP-03, MCP-04, MCP-05, MCP-06, AUTH-04, AUTH-05, PROX-01, PROX-02, PROX-03, PROX-04, PROX-05, STO-01, DEV-01, DEV-02, DEV-03
**Success Criteria** (what must be TRUE):
  1. User can start the gateway and it intercepts MCP traffic on configured port
  2. Gateway validates incoming MCP messages for proper JSON-RPC 2.0 format and size limits
  3. Gateway validates JWT tokens and extracts agent identity/roles
  4. Gateway forwards allowed requests to upstream MCP servers and returns responses unchanged
  5. Gateway blocks requests that fail security checks with appropriate error responses
**Plans**: TBD

### Phase 2: Identity & Access Management
**Goal**: Users can manage identities and enforce fine-grained access control
**Depends on**: Phase 1
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-06, AUTH-07, POL-01, POL-02, POL-03, POL-05, POL-06, STO-02, STO-03, STO-04, SEC-01, SEC-02, SEC-04, SEC-05
**Success Criteria** (what must be TRUE):
  1. Admin user can create account, log in, and log out of the gateway management interface
  2. Gateway supports machine-identity standards for AI agent authentication and issues short-lived credentials
  3. Gateway enforces RBAC and tool authorization policies via OPA
  4. Gateway evaluates requests against resource access policies and makes real-time allow/block decisions
  5. Gateway provides detailed denial reasons when blocking requests
**Plans**: TBD

### Phase 3: Semantic Inspection & AI Security
**Goal**: Gateway intelligently inspects and filters MCP prompts for AI-specific threats
**Depends on**: Phase 2
**Requirements**: SEM-01, SEM-02, SEM-03, SEM-04, SEM-05, SEM-06, POL-04, SEC-03, DEV-04, STO-05, SEC-06
**Success Criteria** (what must be TRUE):
  1. Gateway performs semantic inspection of prompts using local LLM to understand intent and categorize them
  2. Gateway detects and blocks prompt injection and jailbreak attempts in MCP requests
  3. Gateway identifies harmful content in prompts and provides safety scoring (0.0-1.0)
  4. Gateway integrates semantic inspection results with policy evaluation for final allow/block decisions
  5. Gateway mitigates prompt injection and data leakage risks through defense-in-depth approach
**Plans**: TBD

### Phase 4: Observability, Analytics & Operations
**Goal**: Gateway provides comprehensive monitoring, logging, and operational capabilities
**Depends on**: Phase 3
**Requirements**: OBS-01, OBS-02, OBS-03, OBS-04, OBS-05, OBS-06, STO-06, DEV-05, DEV-06, DEV-07
**Success Criteria** (what must be TRUE):
  1. Gateway maintains audit trail of all MCP requests with agent identity, action, timestamp, and decision
  2. Gateway exports structured logs and provides Prometheus metrics for monitoring and alerting
  3. Gateway implements health check endpoints and OpenTelemetry tracing for observability
  4. Gateway implements data retention policies with automatic purge of old audit logs
  5. Gateway provides multiple deployment options (Docker, Helm, K8s manifests)
**Plans**: TBD

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Basic Proxy | 0/0 | Not started | - |
| 2. Identity & Access Management | 0/0 | Not started | - |
| 3. Semantic Inspection & AI Security | 0/0 | Not started | - |
| 4. Observability, Analytics & Operations | 0/0 | Not started | - |