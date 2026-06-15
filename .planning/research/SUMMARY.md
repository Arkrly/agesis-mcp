# Project Research Summary

**Project:** Aegis-MCP  
**Domain:** Zero-trust security gateway for AI agents  
**Researched:** Mon Jun 15 2026  
**Confidence:** MEDIUM  

## Executive Summary

Aegis-MCP is a zero-trust security gateway designed to secure AI agent communications over the Model Context Protocol (MCP). Based on research, experts build such systems using a pipeline architecture where each request flows through sequential security stages: transport handling, validation, authentication, semantic inspection, policy evaluation, decision integration, and proxying to upstream MCP servers. The recommended approach combines Go for performance, Docker for deployment, JWT for authentication, Open Policy Agent (OPA) for authorization, and local LLMs for semantic inspection.

Key risks include over-permissioning of AI agents, prompt injection vulnerabilities, treating AI agents like human users in identity systems, inadequate behavioral monitoring, and supply chain risks in agent tooling. These risks are mitigated through just-in-time access, prompt inspection and data loss prevention, machine-identity standards, user and entity behavior analytics (UEBA) for agents, and tool/server attestation with runtime sandboxing.

## Key Findings

### Recommended Stack

The technology stack emphasizes performance, security, and observability. Core components include Go 1.22+ for efficient concurrency, Docker for consistent deployment, and support for both HTTP and STDIO MCP transports. Authentication relies on industry-standard JWT tokens, while authorization uses Open Policy Agent (OPA) with Rego policies. For semantic inspection, local LLM inference (via Llama.cpp, Ollama, or Hugging Face) ensures privacy and predictable latency. Observability combines structured logging (Zap/Zerolog), Prometheus metrics, and OpenTelemetry tracing. Storage uses PostgreSQL for audit logs, Redis for caching/configuration, and optionally ClickHouse for queryable analytics, with BoltDB/SQLite as embedded alternatives for MVP.

**Core technologies:**
- **Go 1.22+**: High-performance networking and concurrency — Selected for cloud-native service development
- **Docker**: Container platform for deployment — Matches infrastructure requirements and ensures consistency
- **JWT (github.com/golang-jwt/jwt/v5)**: Stateless service-to-service authentication — Industry standard, integrates with existing identity systems
- **Open Policy Agent (OPA)**: Cloud-native authorization engine — Separates policy from code, strong ecosystem
- **Local LLM Inference**: Privacy-preserving prompt analysis — Prevents leakage, enables offline operation, predictable latency
- **Observability Stack (Zap/Zerolog, Prometheus, OpenTelemetry)**: Production monitoring — Wide adoption, comprehensive tracing/metrics/logging
- **PostgreSQL + Redis + ClickHouse**: Scalable storage architecture — Flexible from embedded MVP to production scale

### Expected Features

Research identifies clear distinctions between table stakes (expected by users), differentiators (competitive advantages), and features suitable for post-launch versions.

**Must have (table stakes):**
- Zero-trust request inspection — Every MCP request inspected and authorized
- Semantic prompt inspection using local LLMs — Understands intent beyond pattern matching
- Fine-grained RBAC for agents and tools — Enforces least-privilege access
- Real-time allow/block decisions — Security decisions without unacceptable latency
- Support for standard MCP methods — Basic protocol compliance for ecosystem compatibility
- JWT authentication — Standard auth mechanism integrating with common identity systems
- Audit logging of all requests — Required for compliance and security monitoring
- Configurable message size limits — Prevents resource exhaustion attacks
- Stateless horizontal scalability — Enables cloud-native deployment and load handling
- Security headers (X-Content-Type-Options, etc.) — Basic web security hygiene
- Open-source core — Enables community trust, auditing, and extensibility

**Should have (competitive):**
- Multi-format AI API proxy — Drop-in replacement for multiple AI SDKs
- Virtual API keys with scoped access — Enhances security over static keys
- Composable rate limits and budgets — Prevents abuse and enables cost control
- Real-time cost tracking and attribution — Critical for enterprise chargeback/showback
- Per-user upstream identity propagation (MCP) — Key differentiator for MCP security
- One-paste OAuth onboarding for MCP servers — Reduces friction in connecting MCP servers
- MCP Store with curated templates — Accelerates MCP server adoption
- Tool-level RBAC (MCP) — Fine-grained tool access control
- Namespace isolation for MCP tools — Prevents naming collisions across servers
- Full audit trail with queryable analytics — Enables deep investigations and compliance reporting
- Audit log forwarding — Integrates with existing enterprise security infrastructure
- Usage and cost analytics dashboard — Self-service visibility improves operational efficiency
- Health dashboard with dependency status — Rapid incident detection and response
- Unified log explorer — Correlates events for faster troubleshooting
- Dynamic configuration via UI — Reduces operational overhead, enables self-service
- Multi-instance configuration sync — Supports high availability and scaling
- Data retention policies with automatic purge — Manages storage costs and compliance
- SSO/OIDC integration — Leverages enterprise identity, improves security and UX
- Advanced security controls (session IP binding, distroless containers) — Defense-in-depth for high-security

**Defer (v2+):**
- AI-powered policy optimization — Uses ML to suggest policy improvements
- Blockchain-based audit trail — Immutable logs for regulated industries (long-term)
- Predictive threat detection — ML to identify emerging threats from prompt patterns
- Policy marketplace — Community sharing and rating of Rego policies
- AI-generated MCP server templates — LLMs to generate configurations from natural language
- Edge deployment modes — Lightweight versions for IoT or on-premise isolation
- Federated policy management — Central policy distribution with edge enforcement
- Quantum-resistant cryptography — Future-proofing for long-term secrets
- Homomorphic encryption for processing — Inspection of encrypted prompts (research stage)
- GDPR/CCPA compliance automation — Automated data subject request fulfillment

### Architecture Approach

The recommended architecture follows a pipeline (chain of responsibility) pattern where each request flows through sequential processing stages, with each stage capable of processing, modifying, or terminating the request. This is complemented by sidecar and strategy patterns for observability and policy evaluation flexibility. The architecture cleanly separates concerns into transport, validation, authentication, semantic inspection, policy evaluation, decision integration, proxying, response transformation, observability, and storage layers.

**Major components:**
1. **Transport Layer**: Handles MCP protocol details (JSON-RPC 2.0 over HTTP/Stdio), connection management, message framing
2. **Request Validation**: Validates MCP message format, size limits, basic schema compliance
3. **Initial Authentication**: Validates JWT tokens, extracts agent identity and roles
4. **Semantic Inspection**: Analyzes prompt safety and intent using local LLMs to detect prompt injection, jailbreak attempts, harmful content
5. **Policy Evaluation**: Evaluates requests against RBAC, tool authorization, and resource access policies using Rego
6. **Decision Integration**: Combines semantic inspection and policy evaluation results for final allow/block decisions
7. **MCP Proxy**: Forwards allowed requests to upstream MCP servers, returns responses to clients
8. **Observability & Logging**: Captures audit logs, metrics, traces for monitoring and debugging
9. **Storage Layer**: Persists audit logs, configuration, metrics for compliance and analysis

### Critical Pitfalls

Research identifies five critical pitfalls with specific prevention strategies mapped to implementation phases.

1. **Over-permissioning AI Agent Identities** — Implement just-in-time (JIT) access and fine-grained permissions scoped to specific tasks using dynamic policy evaluation based on agent context
2. **Neglecting Prompt Injection and Data Leakage Risks** — Implement prompt inspection, data loss prevention (DLP) for agent outputs, and enforce strict input validation schemas for all agent-tool interactions
3. **Treating AI Agents Like Human Users in Identity Systems** — Use machine-identity standards (e.g., SPIFFE/SPIRE, OAuth 2.0 Client Credentials) with short-lived certificates and automated rotation tailored for agents
4. **Inadequate Monitoring of Agent Behavior Anomalies** — Deploy user and entity behavior analytics (UEBA) specifically tuned for agent actions, tracking tool usage patterns, data access volumes, and request frequencies
5. **Ignoring Supply Chain Risks in Agent Tooling** — Implement tool/server attestation, signed MCP server registries, and runtime sandboxing for agent-tool interactions

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Foundation & Core Security
**Rationale:** Establishes non-human identity handling and basic gateway functionality before adding advanced features; addresses the most fundamental security risks first.
**Delivers:** Basic zero-trust gateway with JWT authentication, OPA policy engine, basic MCP proxy, and core transport/validation layers
**Addresses:** JWT authentication, Support for standard MCP methods, Basic fine-grained RBAC for agents and tools, Basic real-time allow/block decisions, Configurable message size limits, Security headers, Open-source core, Stateless horizontal scalability (foundational)
**Avoids:** Treating AI Agents Like Human Users in Identity Systems (Phase 1 focus)

### Phase 2: Identity & Access Management Enhancements
**Rationale:** Builds on Phase 1 to strengthen identity controls and supply chain security; addresses over-permissioning and tool trust risks.
**Delivers:** Advanced identity management, fine-grained RBAC, and supply chain security for MCP servers
**Addresses:** Per-user upstream identity propagation (MCP), Virtual API keys with scoped access, Tool-level RBAC (MCP), Namespace isolation for MCP tools, MCP Store with curated templates, Advanced security controls (session IP binding, distroless containers), SSO/OIDC integration
**Avoids:** Over-permissioning AI Agent Identities, Ignoring Supply Chain Risks in Agent Tooling

### Phase 3: Application Security & Intelligence
**Rationale:** Adds intelligent inspection capabilities after identity foundation is secure; focuses on AI-specific threats like prompt injection.
**Delivers:** Semantic prompt inspection, prompt injection protection, and enhanced request inspection capabilities
**Addresses:** Semantic prompt inspection using local LLMs, Enhanced zero-trust request inspection, Real-time cost tracking and attribution, Composable rate limits and budgets
**Avoids:** Neglecting Prompt Injection and Data Leakage Risks

### Phase 4: Observability, Analytics & Operations
**Rationale:** Provides visibility and operational capabilities once core security functions are working; enables production monitoring and management.
**Delivers:** Comprehensive audit trail, analytics dashboards, operational tooling, and configuration management
**Addresses:** Full audit trail with queryable analytics (ClickHouse), Audit log forwarding, Usage and cost analytics dashboard, Health dashboard with dependency status, Unified log explorer, Dynamic configuration via UI, Multi-instance configuration sync, Data retention policies with automatic purge
**Avoids:** Inadequate Monitoring of Agent Behavior Anomalies

### Phase Ordering Rationale
- **Dependency-based ordering**: Identity foundation (Phase 1) must precede advanced identity features (Phase 2); semantic inspection (Phase 3) requires the basic request flow established in Phases 1-2; observability (Phase 4) builds on working core functions
- **Risk mitigation ordering**: Addresses most critical identity risks first (Phase 1), then access control and supply chain risks (Phase 2), then AI-specific application risks (Phase 3), finally monitoring and operational risks (Phase 4)
- **Architecture alignment**: Follows the natural pipeline flow from transport → validation → auth → inspection → policy → decision → proxy → observability/storage

### Research Flags
Phases likely needing deeper research during planning:
- **Phase 3:** Semantic prompt inspection requires LLM model selection, performance benchmarking, and safety scoring tuning — complex integration with multiple LLM backend options
- **Phase 4:** Unified log explorer and advanced analytics need research into log correlation strategies and query performance optimization

Phases with standard patterns (skip research-phase):
- **Phase 1:** Foundation authentication and basic proxy — well-documented patterns with abundant examples (JWT validation, HTTP proxying)
- **Phase 2:** Identity and access management enhancements — OPA and RBAC patterns are well-established in cloud-native security

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Based on explicit constraints in PROJECT.md and detailed technology choices in ARCHITECTURE.md Component Responsibilities table |
| Features | MEDIUM | FEATURES.md reports MEDIUM confidence; includes both table stakes (well-established) and speculative differentiators |
| Architecture | HIGH | Detailed patterns, data flow diagrams, scaling considerations, and anti-patterns provide strong guidance |
| Pitfalls | MEDIUM | PITFALLS.md reports MEDIUM confidence; based on authoritative sources but with logical extension to AI agent domain |

**Overall confidence:** MEDIUM

### Gaps to Address
- **Specific LLM model selection and performance benchmarks**: Need to evaluate Llama 3 8B vs Mistral 7B vs other models for inspection accuracy and latency; plan for performance testing during implementation
- **Exact OPA policy structure and examples**: While OPA is selected, specific Rego policies for agent authorization need definition; will develop during Phase 1 implementation
- **Detailed API specifications for MCP proxy**: Need to define exact MCP method handling and transformation rules; will refine during Phases 1-2
- **Performance targets and latency budgets**: Need to establish acceptable latency thresholds for different operations; will measure and optimize during implementation
- **Specific encryption algorithms and key management details**: Need to define AES-GCM parameters, key rotation schedules, and vault implementation; will detail during Phase 2
- **User interface designs for dashboards and configuration**: Need to prototype UI/UX for admin dashboards and configuration screens; will address during Phase 4

## Sources

### Primary (HIGH confidence)
- PROJECT.md — Project constraints, key decisions, and context
- ARCHITECTURE.md — Component responsibilities, architectural patterns, data flow, and scaling considerations
- Model Context Protocol Specification — https://modelcontextprotocol.io/specification/latest
- NIST SP 800-207: Zero Trust Architecture — https://csrc.nist.gov/publications/detail/sp/800-207/final

### Secondary (MEDIUM confidence)
- FEATURES.md — Feature landscape, differentiators, anti-features, and dependency analysis
- PITFALLS.md — Critical pitfalls, technical debt patterns, integration gotchas, and recovery strategies
- Open Policy Agent Documentation — https://www.openpolicyagent.org
- OAuth 2.0 Threat Model and Security Considerations (RFC 6819) — https://tools.ietf.org/html/rfc6819

### Tertiary (LOW confidence)
- ThinkWatch/Hoop/Casbin Gateway repositories — Competitor feature analysis (implementation details may vary)
- Blog posts and articles on LLM security — Rapidly evolving field, needs validation
- Emerging UEBA tools for agent behavior — New domain, limited proven solutions

---
*Research completed: Mon Jun 15 2026*
*Ready for roadmap: yes*