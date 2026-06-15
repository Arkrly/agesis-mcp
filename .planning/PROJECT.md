# Aegis-MCP

## What This Is

A zero-trust security gateway for AI agents communicating over the Model Context Protocol (MCP). It intercepts MCP traffic via an HTTP proxy, semantically inspects prompts using a local LLM, enforces RBAC and tool authorization policies via OPA (Open Policy Agent), and allows or blocks requests in real-time. Built in Go, targeting developer-friendly, production-grade, open-source core.

## Core Value

Never trust, always verify - every MCP request is inspected and authorized before being forwarded to ensure AI agent communications remain secure and compliant.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Intercept MCP traffic via HTTP proxy
- [ ] Perform semantic inspection of prompts using local LLM
- [ ] Enforce RBAC policies via OPA for agent authorization
- [ ] Enforce tool authorization policies via OPA
- [ ] Make real-time allow/block decisions
- [ ] Provide JWT-based authentication for gateway access
- [ ] Deploy as Docker container
- [ ] Maintain audit trail using embedded database (BoltDB)
- [ ] Implement CI/CD pipeline with GitHub Actions
- [ ] Ensure fully open-source implementation

### Out of Scope

- [MCP protocol implementation] — Aegis-MCP assumes MCP is implemented by clients/servers; it only proxies and secures existing MCP traffic
- [LLM training/fine-tuning] — Uses pre-trained LLMs for inspection but doesn't include model training capabilities
- [Identity provider replacement] — Integrates with existing JWT-based auth rather than replacing identity systems
- [Network-level security] — Focuses on application-layer security; assumes transport security (TLS) is handled separately
- [Multi-tenant SaaS platform] — Designed for self-hosted deployment; not a hosted service offering

## Context

- Target deployment: Docker containers for ease of deployment and isolation
- Authentication: JWT tokens for securing the gateway itself
- Authorization: Role-based (RBAC) with predefined roles in OPA
- Semantic inspection: Local/open-source LLMs (e.g., Llama 3 8B or Mistral 7B) for privacy and predictable latency
- Persistence: Embedded database (BoltDB) for policy metadata and audit logs
- CI/CD: GitHub Actions for automated testing, building, and deployment
- Team size: Small core team (2-3 developers) with community contributions
- Timeline: Targeting MVP in 3 months, production-ready in 6 months
- Open-source commitment: Fully open-source implementation to foster community adoption

## Constraints

- **[Deployment]**: Docker containers — Matches infrastructure requirements and ensures consistent deployment
- **[Auth]**: JWT Tokens — Integrates with existing identity systems and is stateless for scalability
- **[OPA]**: Role-based (RBAC) with predefined roles — Provides clear authorization model that's easy to audit and manage
- **[LLM]**: Local/open-source — Ensures privacy (no prompt leakage), predictable latency/cost, and offline operation
- **[Persistence]**: Embedded database (BoltDB) — Zero-configuration, ACID transactions, and excellent read performance for audit trails
- **[CI/CD]**: GitHub Actions — Leverages existing GitHub repository for streamlined DevOps
- **[License]**: Fully open-source — Encourages community contributions and transparency
- **[Timeline]**: 3-6 months for MVP to production — Balances thoroughness with market readiness
- **[Team]**: Small core team (2-3 developers) — Focuses on core functionality before expanding scope

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Zero-trust security model for MCP | Addresses growing security concerns with AI agent communications | — Pending |
| HTTP proxy architecture | Enables easy deployment without requiring changes to existing MCP clients/servers | — Pending |
| Local LLM for semantic inspection | Protects sensitive prompts from external exposure and provides predictable performance | — Pending |
| OPA for policy enforcement | Industry-standard for cloud-native authorization with strong ecosystem and performance | — Pending |
| Embedded database (BoltDB) | Zero-configuration storage for audit trails and policy metadata | — Pending |
| JWT authentication | Stateless, widely adopted, and easy to integrate with existing identity providers | — Pending |
| GitHub Actions CI/CD | Tight integration with repository hosting for automated testing and deployment | — Pending |
| Fully open-source implementation | Fosters community trust, adoption, and contributions | — Pending |

---
*Last updated: Mon Jun 15 2026 after initialization*