# Pitfalls Research

**Domain:** Zero-trust security gateway for AI agents  
**Researched:** Mon Jun 15 2026  
**Confidence:** MEDIUM (based on authoritative sources for zero trust and MCP, with logical extension to AI agents)

## Critical Pitfalls

### Pitfall 1: Over-permissioning AI Agent Identities

**What goes wrong:**
AI agents are granted excessive permissions (e.g., broad API access, unrestricted tool usage) under the assumption that they need flexibility, leading to potential data breaches or misuse if compromised.

**Why it happens:**
Development teams prioritize ease of integration and fear blocking legitimate agent operations, following the principle of least privilege inadequately for non-human identities.

**How to avoid:**
Implement just-in-time (JIT) access and fine-grained permissions scoped to specific tasks. Use dynamic policy evaluation based on agent context (task, user, data sensitivity).

**Warning signs:**
- Agents using service accounts with domain-wide privileges
- Static API keys with no rotation or expiration
- Lack of audit trails showing fine-grained permission usage

**Phase to address:**
Phase 2: Identity & Access Management (after foundational zero trust architecture)

### Pitfall 2: Neglecting Prompt Injection and Data Leakage Risks

**What goes wrong:**
Failure to validate and sanitize AI agent prompts, allowing malicious inputs to exfiltrate sensitive data or execute unauthorized actions via connected tools/resources.

**Why it happens:**
Focus on network-level security overlooks the application-layer threat unique to AI agents processing natural language inputs.

**How to avoid:**
Implement prompt inspection, data loss prevention (DLP) for agent outputs, and enforce strict input validation schemas for all agent-tool interactions.

**Warning signs:**
- No content filtering on agent prompts/responses
- Agents able to access tools without context-aware authorization
- Absence of DLP policies for agent-generated content

**Phase to address:**
Phase 3: Application Security (post-identity foundation)

### Pitfall 3: Treating AI Agents Like Human Users in Identity Systems

**What goes wrong:**
Using human-centric identity practices (e.g., password-based auth, long-lived sessions) for AI agents, creating credential management overhead and security gaps.

**Why it happens:**
Identity systems designed for humans are repurposed for agents without adapting to non-human, programmatic identity requirements.

**How to avoid:**
Use machine-identity standards (e.g., SPIFFE/SPIRE, OAuth 2.0 Client Credentials) with short-lived certificates and automated rotation tailored for agents.

**Warning signs:**
- Agents using username/password or human-style MFA
- Long-lived tokens (>24 hours) for agent authentication
- Manual credential provisioning for agent workloads

**Phase to address:**
Phase 1: Foundation (during initial identity system design)

### Pitfall 4: Inadequate Monitoring of Agent Behavior Anomalies

**What goes wrong:**
Lack of behavioral baselines and anomaly detection for AI agent activities, enabling compromised agents to operate undetected for extended periods.

**Why it happens:**
Monitoring focuses on network traffic or known malware signatures, not the unique behavioral patterns of AI agent-tool interactions.

**How to avoid:**
Deploy user and entity behavior analytics (UEBA) specifically tuned for agent actions, tracking tool usage patterns, data access volumes, and request frequencies.

**Warning signs:**
- Alerts based solely on IP/reputation, not agent behavior
- No baselining of normal agent tool invocation patterns
- Missing correlation between agent requests and data egress

**Phase to address:**
Phase 4: Monitoring & Response (after core controls are in place)

### Pitfall 5: Ignoring Supply Chain Risks in Agent Tooling

**What goes wrong:**
Failure to verify and monitor the security of third-party tools/data sources that AI agents connect to, creating backdoor vectors via compromised MCP servers.

**Why it happens:**
Trust is implicitly extended to all registered MCP servers without ongoing validation of their security posture or code integrity.

**How to avoid:**
Implement tool/server attestation, signed MCP server registries, and runtime sandboxing for agent-tool interactions.

**Warning signs:**
- Agents connecting to unsigned or unverified MCP servers
- No mechanism to revoke trust in compromised tools
- Lack of sandboxing or resource limits for tool execution

**Phase to address:**
Phase 2: Identity & Access Management (when defining agent-tool trust framework)

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Using static API keys for agent auth | Simpler initial setup | Key leakage, no rotation | Never for production agents |
| Broad tool permissions ("allow all") | Faster integration | Increased blast radius | Only in isolated dev environments |
| Skipping prompt validation | Reduced latency | Vulnerable to injection attacks | Never |

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| MCP Servers | Trusting all registered servers | Implement server attestation and allowlisting |
| Identity Providers | Using human user directories for agents | Deploy dedicated machine identity infrastructure |
| Logging Systems | Treating agent logs like user logs | Schema-specific parsing for agent-structured data |

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Real-time prompt inspection | Increased latency >500ms | Asynchronous scanning with caching | >100 agents/sec throughput |
| Per-request policy evaluation | Policy decision point (PDP) overload | Policy caching with TTL based on agent context | >1K requests/sec |
| Centralized agent monitoring | Log ingestion bottlenecks | Distributed tracing with agent-specific sampling | >10K daily active agents |

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| Agents inheriting user context | Privilege escalation via compromised user | Strict context separation; agent-only identities |
| No encryption for agent-tool traffic | Eavesdropping on sensitive tool interactions | Mutual TLS (mTLS) for all agent-server channels |
| Static trust decisions | Failure to adapt to changing risk | Continuous trust evaluation based on behavior & threat intel |

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Overwhelming consent prompts | Agent usability degradation | Context-aware, just-in-time consent flows |
| Cryptic error messages | Difficult troubleshooting | Agent-specific error codes with remediation guidance |
| No visibility into agent activities | Blind spots in oversight | Dedicated agent activity dashboard with audit trails |

## "Looks Done But Isn't" Checklist

- [ ] **Agent Identity Lifecycle:** Automated provisioning/deprovisioning — verify [hooks trigger on agent deployment/teardown]
- [ ] **Prompt Security:** Input validation and sanitization — verify [tested with known injection patterns]
- [ ] **Least Privilege:** Permissions scoped to specific tasks — verify [no agent has wildcard permissions]
- [ ] **Behavioral Monitoring:** Anomaly detection baseline established — verify [alerts on deviations from normal patterns]
- [ ] **Tool Trust:** Server attestation for MCP connections — verify [all connections use verified servers]

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Over-permissioning | MEDIUM | 1. Immediately revoke excess permissions 2. Implement JIT access 3. Audit recent agent activity for misuse |
| Prompt injection vulnerability | HIGH | 1. Deploy emergency prompt filtering 2. Patch input validation 3. Rotate all agent credentials 4. Notify affected data owners |
| Compromised MCP server | HIGH | 1. Isolate server 2. Revoke trust 3. Rotate credentials for connected agents 4. Rebuild from clean image |

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Over-permissioning | Phase 2 | Permission audit shows least privilege compliance |
| Neglecting prompt injection | Phase 3 | Penetration test fails to exfiltrate data via prompts |
| Treating agents like human users | Phase 1 | Identity logs show machine-auth flows only |
| Inadequate monitoring | Phase 4 | UEBA alerts trigger on simulated anomalous behavior |
| Ignoring tool supply chain | Phase 2 | All MCP connections use attested servers |

## Sources

- NIST SP 800-207: Zero Trust Architecture (authoritative zero trust foundation)
- Model Context Protocol Specification Section 9: Security and Trust & Safety (MCP-specific agent risks)
- Cloudflare One documentation: AI agent governance patterns (industrial implementation)
- BeyondCorp: A New Approach to Enterprise Security (Google's zero trust lessons)
- CISA Zero Trust Maturity Model V2.0 (federal implementation guidance)