# Feature Research

**Domain:** Zero-trust security gateway for AI agents
**Researched:** Mon Jun 15 2026
**Confidence:** MEDIUM

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Zero-trust request inspection | Every MCP request must be inspected and authorized; core to zero-trust model | HIGH | Requires semantic inspection and policy evaluation for each request |
| Semantic prompt inspection using local LLMs | Understand intent and safety beyond pattern matching; prevents prompt injection and jailbreak attempts | HIGH | Requires local LLM integration (e.g., Llama 3) and safety scoring |
| Fine-grained RBAC for agents and tools | Enforce least-privilege access; different agents/tools require different permissions | MEDIUM | Typically implemented via OPA or similar policy engine |
| Real-time allow/block decisions | Security decisions must not introduce unacceptable latency; affects user experience | MEDIUM | Requires low-latency policy evaluation and caching |
| Support for standard MCP methods | Must proxy core MCP functionality: tools/list, tools/call, resources/list, resources/read, prompts/list, prompts/get | LOW | Standard MCP compliance is basic expectation |
| JWT authentication | Industry standard for service-to-service authentication; integrates with existing identity systems | LOW | Well-established pattern; libraries available |
| Audit logging of all requests | Required for compliance, forensics, and detecting policy violations | MEDIUM | Must capture request/response metadata, decisions, and timestamps |
| Configurable message size limits | Prevent resource exhaustion attacks; aligns with MCP specification | LOW | Simple to implement; common in proxies |
| Stateless horizontal scalability | Enables deployment in cloud-native environments; handles variable load | HIGH | Requires externalizing session state (e.g., Redis) |
| Security headers (X-Content-Type-Options, etc.) | Basic web security hygiene; prevents common vulnerabilities | LOW | Standard set of defensive headers |
| Open-source core | Enables community trust, auditing, and extensibility; aligns with project goals | LOW | Licensing and community management overhead |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valuable.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Multi-format AI API proxy (OpenAI, Anthropic, Gemini, etc.) | Drop-in replacement for multiple AI SDKs; reduces integration complexity | HIGH | Requires protocol translation and format conversion |
| Virtual API keys with scoped access | Enables key rotation, least-privilege access, and usage tracking per key | MEDIUM | Requires key management system with hashing and lifecycle |
| Composable rate limits and budgets (sliding window + natural period) | Flexible quota management prevents abuse and controls costs | HIGH | Requires Redis-backed token/bucket algorithms and UI configuration |
| Real-time cost tracking and attribution | Enables chargeback/showback and cost optimization; critical for enterprise adoption | MEDIUM | Requires pricing metadata and token-weighted calculations |
| Per-user upstream identity propagation (MCP) | Ensures upstream systems see real user identity, not shared service accounts | HIGH | Requires OAuth/PAT vault with encryption and per-request token injection |
| One-paste OAuth onboarding for MCP servers | Simplifies MCP server connection; reduces configuration errors and support burden | MEDIUM | Requires implementing RFC 9728/8414/7591 discovery flow |
| MCP Store with curated templates | Accelerates onboarding for popular services (GitHub, Notion, Slack, etc.) | LOW | Community-maintained registry of MCP server configurations |
| Tool-level RBAC (MCP) | Granular control over which tools each agent/key can access | MEDIUM | Extends RBAC to individual tools within MCP servers |
| Namespace isolation for MCP tools | Prevents tool name collisions across different MCP servers | LOW | Automatic prefixing (e.g., github__create_issue) |
| Full audit trail with queryable analytics (ClickHouse) | Enables deep investigations, compliance reporting, and usage insights | HIGH | Columnar storage optimized for analytical queries |
| Audit log forwarding (Syslog, Kafka, HTTP webhooks) | Integrates with existing SIEM and monitoring pipelines | MEDIUM | Standard protocols for enterprise security operations |
| Usage and cost analytics dashboard | Self-service visibility into consumption patterns and spend | MEDIUM | Requires frontend visualization of aggregated audit data |
| Health dashboard with dependency status | Rapid incident response and operational awareness | LOW | Real-time view of PostgreSQL, Redis, MCP server health |
| Unified log explorer (gateway, MCP, audit, access) | Correlates events across system boundaries for troubleshooting | MEDIUM | Requires structured logging and search interface |
| Dynamic configuration via UI (no restart) | Reduces operational overhead and enables self-service adjustments | MEDIUM | Stores configuration in database with change propagation |
| Multi-instance configuration sync (Redis Pub/Sub) | Supports high-availability deployments with consistent state | MEDIUM | Eventually consistent; handles network partitions gracefully |
| Data retention policies with automatic purge | Manages storage costs and complies with data protection regulations | LOW | Configurable TTL for different data types |
| SSO/OIDC integration (Zitadel, Okta, Azure AD) | Leverages existing enterprise identity infrastructure; improves security | MEDIUM | Standard protocols; reduces password fatigue |
| Advanced security controls (session IP binding, distroless containers) | Defense-in-depth against credential theft and container escapes | LOW | Specifically valuable for high-security environments |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Store plaintext API keys/secrets | Simplicity; avoids encryption key management | Critical security breach if database leaked; violates compliance | Encrypt at rest with AES-256-GCM; show plaintext only once upon creation |
| Shared admin token for upstream MCP servers | Simplifies configuration; avoids per-user credential management | Violates zero-trust; upstream sees all actions as same user; no per-user audit | Per-user OAuth/PAT vault with encryption; upstream sees real identity |
| Hard-coded configuration (no UI/API) | Simplicity; avoids configuration complexity | Inflexible; requires code changes and redeployment for adjustments | Database-backed dynamic configuration with UI and API endpoints |
| No rate limiting or budgeting | Maximizes throughput; avoids complexity | Vulnerable to abuse, DoS attacks, and uncontrolled costs | Composable rate limits and budgets with soft/hard options |
| Synchronous blocking on all operations | Simplicity; avoids async complexity | Poor performance under latency; poor resource utilization | Async/non-blocking I/O with connection pooling |
| Single-process, single-instance deployment | Simplicity; avoids distributed systems complexity | No fault tolerance; limited scalability; single point of failure | Design for horizontal scaling; stateless where possible |
| Exposing admin console to internet | Convenience for remote administration | Large attack surface; risks credential theft and configuration changes | Internal-only console port; restrict to admin network/VPN |
| Complex setup without wizard | Assumes technical expertise; avoids UX investment | High barrier to entry; increases support burden; slows adoption | Guided setup wizard for initial configuration |
| Missing SSRF protection | Simplicity; avoids URL validation logic | Vulnerable to internal network probing and metadata service attacks | Validate URLs against allowlists; block private/link-local/metadata IPs |
| Vulnerable to OAuth replay attacks | Simplicity; avoids session binding | Stolen tokens usable from any location/network | Bind sessions to client IP; require re-authentication on network change |
| Non-standard MCP method handling | Simplicity; avoids protocol complexity | Breaks MCP clients; incomplete protocol proxy | Full RFC compliance for all standard MCP methods |
| Monolithic policy engine (no OPA/pluggable) | Simplicity; avoids external dependency | Inflexible policy updates; no community policy sharing | OPA integration with hot-reload; support for Rego policies |
| Opaque error messages (no trace IDs) | Simplicity; avoids logging complexity | Difficult troubleshooting; increases mean time to resolution | Structured errors with trace IDs (internal) and user-safe messages |
| Mandatory TLS termination at gateway | Simplicity; avoids certificate management complexity | Limits deployment flexibility; may double-encrypt in service mesh | Support for both TLS termination and passthrough modes |

## Feature Dependencies

```
Zero-trust request inspection
    └──requires──> Semantic prompt inspection
    └──requires──> Local LLM integration
    └──requires──> Safety scoring and intent detection
    └──requires──> Policy evaluation engine (OPA)
                        └──requires──> JWT authentication
                                └──requires→> User identity extraction
                                        └──requires→> Agent/role mapping
Real-time allow/block decisions
    └──requires→> Low-latency policy evaluation
                        └──requires→> Caching of frequent policies
                        └──requires→> Efficient OPA integration
Support for standard MCP methods
    └──requires→> MCP protocol comprehension
                        └──requires→> JSON-RPC 2.0 handling
Audit logging
    └──requires→> Structured logging
                        └──requires→> Request/response capture
                        └──requires→> Decision audit (allow/block reasons)
                        └──requires→> Storage backend (database/ClickHouse)
Per-user upstream identity propagation (MCP)
    └──requires→> OAuth/PAT vault
                        └──requires→> Encryption at rest (AES-256-GCM)
                        └──requires→> Per-request token injection
                        └──requires→> Upstream identity resolution (JWT/userinfo)
Composable rate limits and budgets
    └──requires→> Redis-backed counters/buckets
                        └──requires→> Sliding window algorithm
                        └──requires→> Natural period reset mechanism
                        └──requires→> Token weighting model
                        └──requires→> UI for configuration
Virtual API keys
    └──requires→> Key generation with entropy
                        └──requires→> SHA-256 hashing for storage
                        └──requires→> Lifecycle management (rotation/expiry)
                        └──requires→> Scoping to services (AI/MCP)
MCP Store with curated templates
    └──requires→> Template registry (GitHub or custom)
                        └──requires→> One-click install workflow
                        └──requires→> OAuth scope pre-configuration
Full audit trail with queryable analytics
    └──requires→> Columnar storage (ClickHouse or similar)
                        └──requires→> Structured audit event schema
                        └──requires→> Efficient ingestion pipeline
                        └──requires→> SQL query interface
Audit log forwarding
    └──requires→> Multiple output channels (Syslog/Kafka/Webhook)
                        └──requires→> Reliable delivery with retries
                        └──requires→> Format translation (JSON to target)
Health dashboard
    └──requires→> Dependency health checks (PostgreSQL, Redis, MCP servers)
                        └──requires→> Aggregated status API
                        └──requires→> Simple UI for status display
Unified log explorer
    └──requires→> Centralized log collection (all services)
                        └──requires→> Structured logging format
                        └──requires→> Full-text search with filters
                        └──requires→> Time-range and severity navigation
Dynamic configuration via UI
    └──requires→> Database-backed configuration store
                        └──requires→> Change detection mechanism
                        └──requires→> Configuration propagation (Redis Pub/Sub)
                        └──requires→> UI for editing all settings
Multi-instance configuration sync
    └──requires→> Redis Pub/Sub for change notifications
                        └──requires→> Conflict resolution strategy
                        └──requires→> Eventually consistent model
Data retention policies
    └──requires→> TTL per data type (audit logs, usage metrics, etc.)
                        └──requires→> Automatic purge mechanism
                        └──requires→> Configurable retention periods
SSO/OIDC integration
    └──requires→> OIDC client implementation
                        └──requires→> User provisioning/just-in-time
                        └──requires→> Role mapping from SSO groups
                        └──requires→> Session management
Advanced security controls
    └──requires→> Session IP binding and validation
                        └──requires→> Distroless container base image
                        └──requires→> Minimal runtime permissions
```

### Dependency Notes

- **[Zero-trust request inspection] requires [Semantic prompt inspection]:** Without understanding prompt intent, authorization decisions are based solely on superficial patterns, missing sophisticated attacks.
- **[Semantic prompt inspection] requires [Local LLM integration]:** Local LLMs provide offline, low-latency inspection without exposing prompts to third parties.
- **[Real-time allow/block decisions] requires [Low-latency policy evaluation]:** High-latency policy checks defeat the purpose of a gateway by introducing unacceptable delays.
- **[Support for standard MCP methods] requires [MCP protocol comprehension]:** Incomplete MCP method support breaks interoperability with standard clients and servers.
- **[Audit logging] requires [Structured logging]:** Unstructured logs are difficult to parse and analyze systematically for security monitoring.
- **[Per-user upstream identity propagation (MCP)] requires [OAuth/PAT vault]:** Securely storing and retrieving per-user credentials is foundational to identity propagation.
- **[Composable rate limits and budgets] requires [Redis-backed counters/buckets]:** Redis provides atomic operations and TTL needed for sliding window and natural period algorithms.
- **[Virtual API keys] requires [Key generation with entropy]:** Cryptographically secure random keys prevent guessing attacks; entropy ensures uniqueness.
- **[MCP Store with curated templates] requires [One-click install workflow]:** Reducing MCP server connection complexity drives adoption of the gateway.
- **[Full audit trail with queryable analytics] requires [Columnar storage]:** Columnar stores like ClickHouse provide fast aggregation for usage and cost analytics.
- **[Audit log forwarding] requires [Multiple output channels]:** Enterprises have diverse SIEM requirements; supporting multiple protocols increases adoption.
- **[Health dashboard] requires [Dependency health checks]:** Operators need immediate visibility into infrastructure health to prevent cascading failures.
- **[Unified log explorer] requires [Centralized log collection]:** Correlating events across gateway, MCP, and audit logs is essential for incident investigation.
- **[Dynamic configuration via UI] requires [Database-backed configuration store]:** Runtime configuration changes without restart require persistent, shared state.
- **[Multi-instance configuration sync] requires [Redis Pub/Sub for change notifications]:** Ensures configuration consistency across horizontally scaled instances.
- **[Data retention policies] requires [TTL per data type]:** Automatic data purge manages storage costs and complies with regulations like GDPR.
- **[SSO/OIDC integration] requires [OIDC client implementation]:** Leveraging existing enterprise identity reduces password fatigue and improves security posture.
- **[Advanced security controls] requires [Session IP binding]:** Binding sessions to client IP prevents token replay attacks from compromised devices.

## MVP Definition

### Launch With (v1)

Minimum viable product — what's needed to validate the concept.

- [ ] Zero-trust request inspection — Core security function; without it, product is not a zero-trust gateway
- [ ] Semantic prompt inspection using local LLMs — Differentiates from pattern-matching gateways; essential for AI-specific threats
- [ ] Fine-grained RBAC for agents and tools — Enables least-privilege access; table stake for enterprise security
- [ ] Real-time allow/block decisions — Performance requirement; excessive latency makes product unusable
- [ ] Support for standard MCP methods — Basic protocol compliance; required to work with existing MCP ecosystem
- [ ] JWT authentication — Standard auth mechanism; integrates with common identity systems
- [ ] Audit logging of all requests — Required for compliance and security monitoring
- [ ] Configurable message size limits — Basic security hygiene; prevents resource exhaustion
- [ ] Stateless horizontal scalability — Enables production deployment; single instance limits adoption
- [ ] Security headers (X-Content-Type-Options, etc.) — Basic web security; expected in any internet-facing service
- [ ] Open-source core — Project goal; enables community trust and contributions

### Add After Validation (v1.x)

Features to add once core is working.

- [ ] Multi-format AI API proxy — Expands beyond MCP to direct AI API use cases; increases TAM
- [ ] Virtual API keys with scoped access — Enhances security and usability over static keys
- [ ] Composable rate limits and budgets — Prevents abuse and enables cost control; expected in production systems
- [ ] Real-time cost tracking and attribution — Critical for enterprise chargeback/showback
- [ ] Per-user upstream identity propagation (MCP) — Key differentiator for MCP security; moves beyond shared tokens
- [ ] One-paste OAuth onboarding for MCP servers — Reduces friction in connecting MCP servers
- [ ] MCP Store with curated templates — Accelerates MCP server adoption
- [ ] Tool-level RBAC (MCP) — Fine-grained tool access control; expected in mature MCP gateways
- [ ] Namespace isolation for MCP tools — Prevents naming collisions; improves usability
- [ ] Full audit trail with queryable analytics — Enables deep investigations and compliance reporting
- [ ] Audit log forwarding — Integrates with existing enterprise security infrastructure
- [ ] Usage and cost analytics dashboard — Self-service visibility improves operational efficiency
- [ ] Health dashboard with dependency status — Rapid incident detection and response
- [ ] Unified log explorer — Correlates events for faster troubleshooting
- [ ] Dynamic configuration via UI — Reduces operational overhead; enables self-service
- [ ] Multi-instance configuration sync — Supports high availability and scaling
- [ ] Data retention policies with automatic purge — Manages storage costs and compliance
- [ ] SSO/OIDC integration — Leverages enterprise identity; improves security and UX
- [ ] Advanced security controls (session IP binding, distroless containers) — Defense-in-depth for high-security deployments

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] AI-powered policy optimization — Uses ML to suggest policy improvements based on usage patterns
- [ ] Blockchain-based audit trail — Immutable logs for regulated industries (long-term)
- [ ] Predictive threat detection — Uses ML to identify emerging threats from prompt patterns
- [ ] Policy marketplace — Community sharing and rating of Rego policies
- [ ] AI-generated MCP server templates — Uses LLMs to generate MCP configurations from natural language
- [ ] Edge deployment modes — Lightweight versions for IoT or on-premise isolation
- [ ] Federated policy management — Central policy distribution with edge enforcement
- [ ] Quantum-resistant cryptography — Future-proofing for long-term secrets
- [ ] Homomorphic encryption for processing — Enables inspection of encrypted prompts (research stage)
- [ ] GDPR/CCPA compliance automation — Automated data subject request fulfillment

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Zero-trust request inspection | HIGH | HIGH | P1 |
| Semantic prompt inspection using local LLMs | HIGH | HIGH | P1 |
| Fine-grained RBAC for agents and tools | HIGH | MEDIUM | P1 |
| Real-time allow/block decisions | HIGH | MEDIUM | P1 |
| Support for standard MCP methods | HIGH | LOW | P1 |
| JWT authentication | HIGH | LOW | P1 |
| Audit logging of all requests | HIGH | MEDIUM | P1 |
| Configurable message size limits | MEDIUM | LOW | P1 |
| Stateless horizontal scalability | HIGH | HIGH | P1 |
| Security headers (X-Content-Type-Options, etc.) | MEDIUM | LOW | P1 |
| Open-source core | MEDIUM | LOW | P1 |
| Multi-format AI API proxy | HIGH | HIGH | P2 |
| Virtual API keys with scoped access | HIGH | MEDIUM | P2 |
| Composable rate limits and budgets | HIGH | HIGH | P2 |
| Real-time cost tracking and attribution | MEDIUM | MEDIUM | P2 |
| Per-user upstream identity propagation (MCP) | HIGH | HIGH | P2 |
| One-paste OAuth onboarding for MCP servers | MEDIUM | MEDIUM | P2 |
| MCP Store with curated templates | MEDIUM | LOW | P2 |
| Tool-level RBAC (MCP) | MEDIUM | MEDIUM | P2 |
| Namespace isolation for MCP tools | MEDIUM | LOW | P2 |
| Full audit trail with queryable analytics | HIGH | HIGH | P2 |
| Audit log forwarding | MEDIUM | MEDIUM | P2 |
| Usage and cost analytics dashboard | MEDIUM | MEDIUM | P2 |
| Health dashboard with dependency status | MEDIUM | LOW | P2 |
| Unified log explorer | MEDIUM | MEDIUM | P2 |
| Dynamic configuration via UI | MEDIUM | MEDIUM | P2 |
| Multi-instance configuration sync | MEDIUM | MEDIUM | P2 |
| Data retention policies with automatic purge | LOW | LOW | P2 |
| SSO/OIDC integration | MEDIUM | MEDIUM | P2 |
| Advanced security controls (session IP binding, distilless containers) | LOW | LOW | P2 |
| AI-powered policy optimization | MEDIUM | HIGH | P3 |
| Blockchain-based audit trail | LOW | HIGH | P3 |
| Predictive threat detection | MEDIUM | HIGH | P3 |
| Policy marketplace | LOW | MEDIUM | P3 |
| AI-generated MCP server templates | MEDIUM | HIGH | P3 |
| Edge deployment modes | MEDIUM | HIGH | P3 |
| Federated policy management | MEDIUM | HIGH | P3 |
| Quantum-resistant cryptography | LOW | HIGH | P3 |
| Homomorphic encryption for processing | LOW | HIGH | P3 |
| GDPR/CCPA compliance automation | LOW | MEDIUM | P3 |

**Priority key:**
- P1: Must have for launch
- P2: Should have, add when possible
- P3: Nice to have, future consideration

## Competitor Feature Analysis

| Feature | ThinkWatch | Hoop | Casbin Gateway | Our Approach |
|---------|------------|------|----------------|--------------|
| Zero-trust request inspection | ✅ (RBAC + semantic) | ✅ (policy-based) | ✅ (OPA + Casbin) | ✅ (OPA + LLM inspection) |
| Semantic prompt inspection using local LLMs | ❌ (focus on API/MCP, not deep semantic) | ❌ | ❌ | ✅ (core differentiator) |
| Fine-grained RBAC for agents and tools | ✅ (5-tier) | ✅ (policy-based) | ✅ (Casbin RBAC) | ✅ (OPA-based RBAC) |
| Real-time allow/block decisions | ✅ (low latency claimed) | ✅ (wire-level <5ms) | ✅ | ✅ (optimized OPA + caching) |
| Support for standard MCP methods | ✅ (full MCP proxy) | ✅ (multi-protocol) | ✅ (MCP focus) | ✅ (standard MCP methods) |
| JWT authentication | ✅ | ✅ (implied) | ✅ | ✅ (standard JWT) |
| Audit logging of all requests | ✅ (ClickHouse) | ❓ (not specified) | ❓ | ✅ (structured audit trail) |
| Configurable message size limits | ✅ | ✅ | ✅ | ✅ |
| Stateless horizontal scalability | ✅ (Redis-backed) | ✅ | ✅ | ✅ |
| Security headers | ✅ | ❓ | ❓ | ✅ (standard set) |
| Open-source core | ❌ (BSL-1.1) | ✅ (MIT) | ✅ (Apache 2.0) | ✅ (OSI-approved license) |
| Multi-format AI API proxy | ✅ (OpenAI/Anthropic/Gemini) | ❌ (MCP-focused) | ❌ | ✅ (planned for v1.x) |
| Virtual API keys with scoped access | ✅ (tw- keys) | ❓ | ❓ | ✅ (planned for v1.x) |
| Composable rate limits and budgets | ✅ (sliding window + natural) | ❓ | ❓ | ✅ (planned for v1.x) |
| Real-time cost tracking and attribution | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |
| Per-user upstream identity propagation (MCP) | ✅ (OAuth/PAT per user) | ❓ | ❓ | ✅ (planned for v1.x) |
| One-paste OAuth onboarding for MCP servers | ✅ (RFC flow) | ❓ | ❓ | ✅ (planned for v1.x) |
| MCP Store with curated templates | ✅ (23+ templates) | ❓ | ❓ | ✅ (planned for v1.x) |
| Tool-level RBAC (MCP) | ✅ (per-role grants) | ❓ | ❓ | ✅ (planned for v1.x) |
| Namespace isolation for MCP tools | ✅ (github__create_issue) | ❓ | ❓ | ✅ (planned for v1.x) |
| Full audit trail with queryable analytics | ✅ (ClickHouse) | ❓ | ❓ | ✅ (planned for v1.x) |
| Audit log forwarding | ✅ (Syslog/Kafka/Webhook) | ❓ | ❓ | ✅ (planned for v1.x) |
| Usage and cost analytics dashboard | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |
| Health dashboard with dependency status | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |
| Unified log explorer | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |
| Dynamic configuration via UI | ✅ (database + Redis Pub/Sub) | ❓ | ❓ | ✅ (planned for v1.x) |
| Multi-instance configuration sync | ✅ (Redis Pub/Sub) | ❓ | ❓ | ✅ (planned for v1.x) |
| Data retention policies with automatic purge | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |
| SSO/OIDC integration | ✅ (Zitadel/OIDC) | ❓ | ❓ | ✅ (planned for v1.x) |
| Advanced security controls (session IP binding, distroless containers) | ✅ | ❓ | ❓ | ✅ (planned for v1.x) |

## Sources

- Aegis-MCP PRD.md: Product requirements and goals
- Aegis-MCP ARCH.md: Technical architecture and component design
- ThinkWatch GitHub Repository: Enterprise AI bastion host features (https://github.com/ThinkWatchProject/ThinkWatch)
- Hoop GitHub Repository: Multi-protocol gateway (https://github.com/hoophq/hoop)
- Casbin Gateway GitHub Repository: AI & MCP security gateway (https://github.com/apasdin/casbin-gateway)
- MCP Specification: https://modelcontextprotocol.io
- Zero Trust Architecture (NIST SP 800-207): https://csrc.nist.gov/publications/detail/sp/800-207/final
- OAuth 2.0 Threat Model and Security Considerations (RFC 6819): https://tools.ietf.org/html/rfc6819
- OPA Documentation: https://www.openpolicyagent.org
- Redis Rate Limiting Patterns: https://redis.io/topics/rate-limiting