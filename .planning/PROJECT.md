# Aegis-MCP

## What This Is

A zero-trust security gateway for AI agents communicating over the Model Context Protocol (MCP). It intercepts MCP traffic via an HTTP proxy, semantically inspects prompts using a local LLM, enforces RBAC and tool authorization policies via OPA (Open Policy Agent), and allows or blocks requests in real-time. Built in Go, targeting developer-friendly, production-grade, open-source core.

## Core Value

Never trust, always verify - every MCP request is inspected and authorized before being forwarded to ensure AI agent communications remain secure and compliant.

## Requirements

### Validated

- [x] Intercept MCP traffic via HTTP proxy
- [x] Perform semantic inspection (Heuristic V1)
- [x] Enforce RBAC and tool authorization via OPA
- [x] Real-time allow/block decisions
- [x] JWT-based authentication
- [x] Audit trail with BoltDB
- [x] Prometheus metrics and structured logging
- [x] Docker-based deployment

### Active

- [ ] Real local LLM integration (llama.cpp)
- [ ] Comprehensive documentation (v0.1.0)
- [ ] Final security hardening and threat modeling
- [ ] Release automation (GitHub Actions)

### Out of Scope

- [MCP protocol implementation] — Aegis-MCP assumes MCP is implemented by clients/servers.
- [LLM training/fine-tuning] — Uses pre-trained LLMs.
- [Identity provider replacement] — Integrates with existing JWT-based auth.
- [Multi-tenant SaaS platform] — Designed for self-hosted deployment.

## Context

- **Progress**: ~85% complete. Core infrastructure, auth, policy engine, and observability are fully implemented.
- **Next Steps**: Phase 4 (Real LLM) and Phase 8 (Documentation).
- **Timeline**: Targeting v0.1.0 release soon.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Zero-trust security model | Addresses growing security concerns with AI agent communications | COMPLETED |
| HTTP proxy architecture | Enables easy deployment without requiring changes to existing MCP clients/servers | COMPLETED |
| Local LLM for semantic inspection | Protects sensitive prompts and provides predictable performance | PARTIAL (Heuristic V1) |
| OPA for policy enforcement | Industry-standard for cloud-native authorization | COMPLETED |
| Embedded database (BoltDB) | Zero-configuration storage for audit trails | COMPLETED |
| JWT authentication | Stateless, widely adopted, and easy to integrate | COMPLETED |
| GitHub Actions CI/CD | Tight integration with repository hosting | COMPLETED |
| Fully open-source | Fosters community trust and adoption | COMPLETED |

---
*Last updated: Wed Jun 17 2026*
