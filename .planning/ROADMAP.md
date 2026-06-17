# ROADMAP.md

## Phases

- [x] **Phase 0: Project Foundations** - Establish development environment, tooling, and basic project structure.
- [x] **Phase 1: Core Proxy Infrastructure** - Build a functional HTTP-to-MCP proxy with basic request/response handling.
- [x] **Phase 2: Authentication Layer** - Implement JWT-based authentication and authorization context extraction.
- [x] **Phase 3: OPA Policy Engine Integration** - Integrate OPA for policy decision evaluation.
- [/] **Phase 4: LLM Semantic Inspection** - Integrate local LLM for prompt safety and intent analysis (Heuristic V1 completed, llama.cpp bindings pending).
- [x] **Phase 5: Policy Decision Engine & Audit Logging** - Combine auth, semantic, and OPA results into final allow/block decisions; implement persistent audit logging.
- [x] **Phase 6: Observability & Production Readiness** - Make the system observable, configurable, and production-ready.
- [/] **Phase 7: Testing & Security Hardening** - Ensure security, reliability, and readiness for external review.
- [ ] **Phase 8: Documentation & Release Preparation** - Prepare for public release and community adoption.
- [ ] **Phase 9: Launch & Community Building** - Release v0.1.0 and begin community engagement.

## Progress Details

### Completed Phases
- **Phase 0-3, 5-6**: Fully implemented and tested.
- **Phase 1**: Core proxy with JSON-RPC 2.0 framing, transparent proxy mode, and graceful shutdown.
- **Phase 2**: JWT extraction (HS256) and claim validation.
- **Phase 3**: Embedded OPA, rego loading, and decision caching.
- **Phase 5**: Decision combiner, BoltDB audit logging, and per-agent rate limiting.
- **Phase 6**: Prometheus metrics, health checks, and structured logging.

### Active Phases
- **Phase 4 (Semantic Inspection)**: Heuristic safety inspection is active. Real local LLM integration via llama.cpp is the next major step.
- **Phase 7 (Security Hardening)**: Fuzzing, property-based testing, and race detection are implemented. Threat modeling and final security audit are pending.

### Upcoming Phases
- **Phase 8**: Focusing on comprehensive documentation (GETTING_STARTED.md, API_REFERENCE.md, etc.) and release automation.
- **Phase 9**: Official v0.1.0 release.

## Progress Table

| Phase | Status | Completed |
|-------|--------|-----------|
| 0. Foundations | COMPLETED | Yes |
| 1. Core Proxy | COMPLETED | Yes |
| 2. Authentication | COMPLETED | Yes |
| 3. OPA Integration | COMPLETED | Yes |
| 4. Semantic Inspection | IN PROGRESS | Partial |
| 5. Policy Decision | COMPLETED | Yes |
| 6. Observability | COMPLETED | Yes |
| 7. Security Hardening | IN PROGRESS | Partial |
| 8. Documentation | NOT STARTED | - |
| 9. Launch | NOT STARTED | - |
