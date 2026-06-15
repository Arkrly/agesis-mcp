# Architecture Research

**Domain:** Zero-trust security gateway for AI agents
**Researched:** Mon Jun 15 2026
**Confidence:** HIGH

## Standard Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Entry/Ingress Layer                      │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐     │
│  │ Transport   │  │ Request      │  │ Initial        │     │
│  │ Layer       │  │ Validation   │  │ Authentication │     │
│  └──────┬──────┘  └──────┬───────┘  └──────┬──────────┘     │
│         │                │                 │                 │
├─────────┴────────────────┴────────────────┴─────────────────┤
│                  Policy Decision Engine                     │
├─────────────────────────────────────────────────────────────┤
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐  │
│  │ Semantic       │  │ Policy         │  │ Decision     │  │
│  │ Inspection     │  │ Evaluation     │  │ Integration  │  │
│  │ (LLM-based)    │  │ (OPA/Rego)     │  │ (Allow/Block)│  │
│  └──────┬────────┘  └──────┬──────────┘  └──────┬────────┘  │
│         │                  │                    │           │
├─────────┴────────────────┴────────────────┴─────────────────┤
│                    Egress/Proxy Layer                       │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌────────────────┐  ┌──────────────┐    │
│  │ MCP          │  │ Response       │  │ Observability  │    │
│  │ Proxy        │  │ Transformation │  │ & Logging      │    │
│  └──────┬───────┘  └──────┬──────────┘  └──────┬────────┘    │
│         │                 │                    │             │
└─────────┴─────────────────┴────────────────────┴────────────┘
                                                    ↓
                                         ┌──────────────────┐
                                         │   Storage Layer  │
                                         │ (Audit Logs,     │
                                         │  Config, Metrics)│
                                         └──────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| Transport Layer | Handles MCP protocol details (JSON-RPC 2.0 over HTTP/Stdio), connection management, message framing | Custom HTTP server, stdio pipes for local MCP servers |
| Request Validation | Validates MCP message format, size limits, basic schema compliance | JSON schema validation, message size checks |
| Initial Authentication | Validates JWT tokens, extracts agent identity and roles | JWT validation libraries (e.g., github.com/golang-jwt/jwt/v5) |
| Semantic Inspection | Analyzes prompt safety and intent using local LLMs to detect prompt injection, jailbreak attempts, harmful content | Local LLM inference (e.g., Llama.cpp, Ollama) with safety scoring |
| Policy Evaluation | Evaluates requests against RBAC, tool authorization, and resource access policies using Rego | Open Policy Agent (OPA) with custom Rego policies |
| Decision Integration | Combines semantic inspection and policy evaluation results to make final allow/block decisions | Custom logic combining safety scores and policy decisions |
| MCP Proxy | Forwards allowed requests to upstream MCP servers, returns responses to clients | HTTP client for remote MCP servers, stdio for local |
| Response Transformation | Processes MCP responses, applies any necessary transformations or filtering | Similar to request validation but for responses |
| Observability & Logging | Captures audit logs, metrics, traces for monitoring and debugging | Structured logging (JSON), Prometheus metrics, OpenTelemetry tracing |
| Storage Layer | Persists audit logs, configuration, metrics for compliance and analysis | Database (PostgreSQL) for audit logs, Redis for caching/config, ClickHouse for analytics |

## Recommended Project Structure

```
src/
├── transport/              # MCP protocol handling (HTTP/Stdio)
│   ├── http/               # HTTP transport implementation
│   │   ├── server.go       # MCP HTTP server
│   │   └── client.go       # MCP HTTP client for proxying
│   └── stdio/              # STDIO transport for local MCP servers
│       ├── server.go       # STDIO MCP server
│       └── client.go       # STDIO MCP client
├── validation/             # Request/response validation
│   ├── request_validator.go
│   └── response_validator.go
├── auth/                   # Authentication (JWT validation)
│   ├── jwt_validator.go
│   └── claims_extractor.go
├── inspection/             # Semantic inspection using LLMs
│   ├── llm_inspector.go
│   ├── safety_scorer.go
│   └── intent_classifier.go
├── policy/                 # Policy evaluation (OPA integration)
│   ├── opa_client.go
│   ├── policy_evaluator.go
│   └── rego_policies/      # Rego policy files
├── decision/               # Decision integration logic
│   ├── decision_engine.go
│   └── decision_types.go
├── proxy/                  # MCP request/response proxying
│   ├── mcp_proxy.go
│   ├── upstream_client.go
│   └── response_handler.go
├── observability/          # Logging, metrics, tracing
│   ├── logger.go
│   ├── metrics.go
│   └── tracer.go
├── storage/                # Data persistence
│   ├── audit_log_repo.go
│   ├── config_repo.go
│   └── metrics_repo.go
├── config/                 # Configuration management
│   ├── config.go
│   └── dynamic_config.go
└── main.go                 # Application entrypoint
```

### Structure Rationale

- **transport/:** Separated by transport type (HTTP vs STDIO) to isolate protocol-specific concerns; follows MCP specification which defines these two transports
- **validation/:** Contains request and response validation logic; keeps validation concerns separate from business logic
- **auth/:** Focused solely on JWT authentication and claims extraction; follows single responsibility principle
- **inspection/:** Houses all LLM-based semantic inspection components; allows for easy swapping of different LLM backends
- **policy/:** Encapsulates OPA integration and policy evaluation; makes it easy to update policies without changing core logic
- **decision/:** Combines results from inspection and policy evaluation; centralizes the allow/block decision logic
- **proxy/:** Handles the actual MCP protocol proxying to upstream servers; separates concern of forwarding from decision-making
- **observability/:** Centralized logging, metrics, and tracing; follows observability best practices for production systems
- **storage/:** Data persistence layer for audit logs, configuration, and metrics; enables separation of concerns
- **config/:** Dynamic configuration management; supports runtime configuration changes without restart
- **main.go:** Standard application entrypoint that wires all components together

## Architectural Patterns

### Pattern 1: Pipeline (Chain of Responsibility)

**What:** Each request flows through a series of processing stages (transport → validation → auth → inspection → policy → decision → proxy) where each stage can process, modify, or terminate the request.

**When to use:** When requests need to undergo multiple sequential checks or transformations, and any step can short-circuit the flow (e.g., reject invalid requests early).

**Trade-offs:**
- Pros: Clear separation of concerns, easy to add/remove stages, explicit failure points
- Cons: Potential latency from multiple stages, complexity in error propagation

**Example:**
```go
func HandleMCPRequest(req *MCPRequest) (*MCPResponse, error) {
    // Transport layer - already handled by HTTP server
    
    // Validation stage
    if err := v.ValidateRequest(req); err != nil {
        return &MCPResponse{Error: validationError(err)}, nil
    }
    
    // Authentication stage
    claims, err := a.Authenticate(req.Headers["Authorization"])
    if err != nil {
        return &MCPResponse{Error: authError(err)}, nil
    }
    
    // Semantic inspection
    safetyScore, intents, err := i.InspectPrompt(req)
    if err != nil {
        return &MCPResponse{Error: inspectionError(err)}, nil
    }
    
    // Policy evaluation
    allowed, reason, err := p.Evaluate(req, claims, safetyScore, intents)
    if err != nil {
        return &MCPResponse{Error: policyError(err)}, nil
    }
    
    // Decision integration
    if !d.ShouldAllow(allowed, safetyScore) {
        return &MCPResponse{Error: decisionError(reason)}, nil
    }
    
    // MCP proxying
    return p.ProxyRequest(req)
}
```

### Pattern 2: Sidecar Pattern for Observability

**What:** Observability concerns (logging, metrics, tracing) are implemented as cross-cutting concerns that wrap or decorate core business logic without modifying it.

**When to use:** When you want to add monitoring capabilities to existing services without changing their core functionality.

**Trade-offs:**
- Pros: Non-invasive, easy to enable/disable, follows separation of concerns
- Cons: Can add complexity if not implemented carefully, potential performance overhead

**Example:**
```go
type ObservingLLMInspector struct {
    Inspector   LLMSemanticInspector
    Logger      *zap.Logger
    Metrics     *ObservabilityMetrics
    Tracer      *oteltrace.Tracer
}

func (o *ObservingLLMInspector) InspectPrompt(ctx context.Context, prompt string) (float64, []string, error) {
    ctx, span := o.Tracer.Start(ctx, "LLMInspection")
    defer span.End()
    
    safetyScore, intents, err := o.Inspector.InspectPrompt(ctx, prompt)
    
    if err != nil {
        o.Logger.Error("LLM inspection failed", zap.Error(err))
        o.Metrics.InspectionErrors.Inc()
        return 0, nil, err
    }
    
    o.Logger.Info("LLM inspection completed",
        zap.Float64("safety_score", safetyScore),
        zap.Strings("intents", intents))
    o.Metrics.Inspections.Inc()
    o.Metrics.AverageSafetyScore.Observe(safetyScore)
    
    return safetyScore, intents, nil
}
```

### Pattern 3: Strategy Pattern for Policy Evaluation

**What:** Different policy evaluation strategies can be swapped in and out (e.g., OPA vs local Rego evaluation vs external policy service) without changing the core decision logic.

**When to use:** When you need to support multiple policy engines or want to easily switch between different policy evaluation approaches.

**Trade-offs:**
- Pros: Flexibility to change policy engines, easier testing with mocks, separation of policy concerns
- Cons: Additional abstraction layer, slight performance overhead

**Example:**
```go
type PolicyEvaluator interface {
    EvaluateRequest(ctx context.Context, req *MCPRequest, claims jwt.MapClaims) (bool, string, error)
    ReloadPolicies() error
}

type OpaPolicyEvaluator struct {
    client *opa.Client
    module string
}

func (e *OpaPolicyEvaluator) EvaluateRequest(ctx context.Context, req *MCPRequest, claims jwt.MapClaims) (bool, string, error) {
    // OPA evaluation logic
}

type LocalRegoPolicyEvaluator struct {
    rego *rego.Rego
}

func (e *LocalRegoPolicyEvaluator) EvaluateRequest(ctx context.Context, req *MCPRequest, claims jwt.MapClaims) (bool, string, error) {
    // Local Rego evaluation logic
}
```

## Data Flow

### Request Flow

```
[AI Agent/MCP Client]
        ↓ (MCP Request over HTTP/Stdio)
[Transport Layer] → [Request Validator] → [Authenticator] → [Semantic Inspector]
        ↓                   ↓                   ↓                   ↓
[Policy Evaluator] → [Decision Engine] → [MCP Proxy] → [Upstream MCP Server]
        ↓                   ↓                   ↓                   ↓
[Response Handler] ← [Response Validator] ← [Observability] ← [Storage]
        ↓
[AI Agent/MCP Client] ← (MCP Response)
```

### State Management

```
[Configuration Store]
        ↓ (watch for changes)
[Dynamic Config Manager] ←→ [All Components] (via config updates)
        ↓
[Audit Log Store] ←→ [Observability Layer] (async writes)
        ↓
[Metrics Store] ←→ [Observability Layer] (async aggregation)
```

### Key Data Flows

1. **Security Inspection Flow:** MCP Request → Transport → Validation → Auth → Semantic Inspection → Policy Evaluation → Decision → (Allow: Proxy to MCP Server | Block: Return Error Response)
2. **Observability Flow:** All components emit structured logs and metrics → Observability Layer → Async writes to Storage Layer (database/ClickHouse) + real-time streaming to external systems (Kafka/Syslog/webhooks)
3. **Configuration Flow:** External config changes → Config Watcher → Dynamic Config Manager → Push updates to all components via Redis Pub/Sub or direct function calls

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|--------------------------|
| 0-1k users | Single instance deployment is fine; use SQLite for storage if needed |
| 1k-100k users | Horizontal scaling of stateless components; separate Redis for caching/sessions; PostgreSQL for audit logs |
| 100k+ users | Microservices separation; dedicated instances for high-load components (inspection, proxy); read replicas for databases; CDN for static assets |

### Scaling Priorities

1. **First bottleneck:** Semantic inspection (LLM inference) - typically the most computationally expensive operation
   - How to fix: LLM batching, GPU acceleration, quantization, caching frequent prompts, LLM pooling

2. **Second bottleneck:** Policy evaluation (OPA) - can become complex with many policies
   - How to fix: Policy caching, OPA optimization, policy simplification, incremental policy evaluation

3. **Third bottleneck:** MCP proxying - network I/O to upstream servers
   - How to fix: Connection pooling, keep-alive connections, geographic load balancing of upstream MCP servers

## Anti-Patterns

### Anti-Pattern 1: Monolithic LLM Inspection

**What people do:** Using a single, large LLM for all semantic inspection tasks without considering performance or specialization.

**Why it's wrong:** LLMs are resource-intensive; using one large model for everything creates unnecessary latency and cost; different inspection tasks may benefit from different model sizes/specializations.

**Do this instead:** Use a pipeline of models - lightweight models for initial screening (e.g., classifying obvious safe/unsafe prompts), reserving larger models for ambiguous cases; consider specialized models for specific threat types (prompt injection vs jailbreak vs harmful content).

### Anti-Pattern 2: Tight Coupling Between Inspection and Policy

**What people do:** Making policy decisions directly within the LLM inspection logic or vice versa, creating dependencies that make it hard to change one without affecting the other.

**Why it's wrong:** Violates separation of concerns; makes testing difficult; prevents independent scaling of inspection and policy components; complicates policy updates.

**Do this instead:** Keep inspection and policy evaluation as separate stages with well-defined interfaces; use a decision integration layer that combines their outputs; this allows each to be developed, tested, and scaled independently.

### Anti-Pattern 3: Synchronous Blocking on All Operations

**What people do:** Using synchronous, blocking I/O for all operations (LLM inference, network calls, disk writes) which underutilizes system resources.

**Why it's wrong:** Creates poor performance under load; limits throughput; causes resource exhaustion under concurrent requests; doesn't take advantage of async capabilities in modern languages/runtimes.

**Do this instead:** Use async/non-blocking I/O where possible; implement connection pooling for external services; use worker queues for expensive operations like LLM inference; leverage Go's goroutines or Rust's async runtime for concurrency.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|-------------------|-------|
| Upstream MCP Servers | HTTP client (for remote) or stdio pipes (for local) | Must support both MCP transports; handle connection pooling and timeouts |
| Identity Providers | JWT validation with JWKS endpoint support | Should support common OAuth 2.0/OIDC providers; validate token signatures and claims |
| External Policy OPA Server | HTTP client to remote OPA | For centralized policy management; fallback to local OPA for resilience |
| SIEM/Logging Systems | Async log forwarders (Kafka, Syslog, webhook) | Support multiple formats; handle network partitions with buffering and retries |
| Monitoring Systems | Prometheus metrics endpoint + OpenTelemetry tracing | Standard endpoints; ensure cardinality control for high-scale deployments |
| Storage Systems | Database/sqlx for PostgreSQL, redis.go for Redis | Connection pooling; proper error handling and retry logic |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Transport ↔ Validation | Function calls with validated MCPRequest objects | Validation should return sanitized/validated request or error |
| Validation ↔ Auth | Function calls extracting JWT from validated request | Auth should return identity/claims or error |
| Auth ↔ Inspection | Function calls passing claims and request to inspector | Inspector should not need to re-validate JWT |
| Inspection ↔ Policy | Function calls passing inspection results and claims | Policy evaluator should receive pre-computed safety scores and intents |
| Policy ↔ Decision | Function calls returning boolean allowance and reason | Decision layer combines policy outcome with inspection scores |
| Decision ↔ Proxy | Function calls forwarding allowed requests or returning errors | Proxy should only receive requests that have passed all security checks |
| All Components ↔ Observability | Structured logging via dependency-injected logger | Logs should include trace IDs for correlation across components |
| All Components ↔ Config | Dynamic config updates via pub/sub or callback | Configuration changes should be propagated without requiring restart |

## Sources

- NIST SP 800-207: Zero Trust Architecture - https://csrc.nist.gov/publications/detail/sp/800-207/final
- OAuth 2.0 Threat Model and Security Considerations (RFC 6819) - https://tools.ietf.org/html/rfc6819
- Open Policy Agent Documentation - https://www.openpolicyagent.org
- Model Context Protocol Specification - https://modelcontextprotocol.io/specification/latest
- Zero Trust Architecture for AI/ML Systems - IBM Research (internal research)
- API Gateway Patterns - Kong Gateway Documentation
- LLM Security: Prompt Injection, Jailbreaking, and Defense Techniques - arXiv:2305.14965
- Cloud Native Security Patterns - CNCF White Papers