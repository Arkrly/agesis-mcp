# Agent Behavior Rules for Aegis-MCP Development

## Core Principles
When contributing to Aegis-MCP, the AI coding agent must prioritize security, correctness, and maintainability over speed or convenience. This is a security-critical project where vulnerabilities could have severe consequences.

## Mandatory Practices (MUST Do)

### Security-First Development
1. **Threat Model Every Change**: Before implementing any feature, identify potential attack vectors it might introduce
2. **Input Validation First**: All external inputs (MCP requests, policy updates, config changes) must be validated before processing
3. **Principle of Least Privilege**: Code should run with minimal permissions; never grant unnecessary access
4. **Secure Defaults**: Configurations should default to most restrictive secure state
5. **Dependency Vigilance**: Regularly update dependencies; never introduce dependencies without security review

### Code Quality & Maintainability
1. **Go Idiomatic Code**: Follow [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. **Comprehensive Error Handling**: Every function must handle errors explicitly; never ignore errors with `_`
3. **Clear Documentation**: Export functions and complex structs must have godoc comments
4. **Unit Test Coverage**: Aim for >80% unit test coverage; security-critical paths >95%
5. **Dependency Injection**: Use interfaces for external services (LLM clients, OPA, databases) to enable testing

### Performance Awareness
1. **Benchmark Sensitive Paths**: Latency-critical paths (request processing) must have benchmarks
2. **Memory Efficiency**: Avoid unnecessary allocations in hot paths; use sync.Pool where appropriate
3. **Concurrency Safety**: All shared state must be properly synchronized; prefer immutable data structures
4. **Resource Bounding**: Implement limits on concurrent requests, queue sizes, and cache entries

### Testing Rigor
1. **Security Testing**: Include fuzzing for input parsers and property-based testing for security logic
2. **Integration Tests**: Test full request flows with mock LLM and OPA servers
3. **Chaos Testing**: Simulate network failures, slow dependencies, and resource exhaustion
4. **Golden Request Tests**: Maintain a corpus of known good/bad MCP requests for regression testing

### Open-Source Collaboration
1. **Clear Commit Messages**: Follow conventional commits format (feat:, fix:, docs:, etc.)
2. **Backward Compatibility**: Maintain API compatibility unless absolutely necessary to break
3. **Upgrade Paths**: Provide clear migration guides for breaking changes
4. **Third-Party Contributions**: Review all contributions with same rigor as internal work

## Prohibited Practices (MUST NOT Do)

### Security Violations
1. **Never Log Sensitive Data**: Never log MCP prompts, JWT tokens, or policy decisions in plaintext
2. **No Hardcoded Secrets**: Never embed API keys, passwords, or cryptographic keys in source code
3. **Avoid Eval/Dynamic Code**: Never use `go:generate` with user input or similar dynamic execution patterns
4. **No Weak Cryptography**: Never use MD5, SHA1, DES, or other deprecated algorithms
5. **No Timing Attacks**: Use constant-time comparison for cryptographic operations (e.g., JWT signature validation)

### Quality Anti-Patterns
1. **No Global State**: Avoid package-level variables that create hidden dependencies
2. **No Magic Numbers**: Replace unexplained constants with named variables
3. **No Empty Catch-Alls**: Never use `recover()` to ignore panics without logging and alerting
4. **No Ignored Linters**: All code must pass `golangci-lint` with project-configured rules
5. **No Vendor Lock-in**: Avoid platform-specific code unless abstracted behind interfaces

### Performance Pitfalls
1. **No Unbounded Growth**: Never allow queues, caches, or slices to grow without limits
2. **No Blocking in Hot Paths**: Avoid long-running operations in request processing goroutines
3. **No Excessive Logging**: Debug-level logging should not impact performance in production
4. **No Premature Optimization**: Optimize only after profiling identifies actual bottlenecks

### Process Violations
1. **No Direct Main Pushes**: All changes must go through pull requests with review
2. **No Skipping Tests**: Never skip tests to make CI pass; fix failing tests instead
3. **No Undocumented Config**: All configuration options must be documented in CONFIG.md
4. **No Breaking Changes in Patch Versions**: Follow semantic versioning strictly

## Special Considerations for Security Components

### LLM Integration
- **Must**: Validate LLM output format before parsing; implement timeouts for LLM calls
- **Must**: Implement fallback rules when LLM is unavailable (fail-closed by default)
- **Must Not**: Send raw MCP prompts to external LLM APIs without encryption and audit
- **Must Not**: Trust LLM outputs without validation; treat as untrusted input

### OPA Policy Engine
- **Must**: Cache policy decisions with appropriate TTL to reduce OPA calls
- **Must**: Validate OPA bundle integrity if using bundle distribution
- **Must Not**: Allow dynamic policy loading from untrusted sources without signature verification
- **Must Not**: Expose OPA admin interface without authentication

### Cryptography
- **Must**: Use standard library crypto packages only (`crypto/`, `golang.org/x/crypto`)
- **Must**: Use TLS 1.2+ for any external communications
- **Must Not**: Implement custom cryptographic protocols
- **Must Not**: Use RSA keys < 2048-bit or ECC curves < P-256

## Review Checklist for AI Agent
Before marking a task complete, verify:
- [ ] All inputs validated and sanitized
- [ ] Errors handled explicitly (no `_ = err`)
- [ ] Security implications considered and documented
- [ ] Unit tests cover normal, edge, and error cases
- [ ] Benchmarks exist for performance-critical paths
- [ ] Code follows Go idioms and project conventions
- [ ] Documentation updated if APIs changed
- [ ] No secrets or sensitive data in logs/commits
- [ ] Dependencies updated and vetted
- [ ] CHANGELOG entry added if user-facing change

## Violations & Remediation
Any violation of these rules requires:
1. Immediate rework to bring code into compliance
2. Documentation of why violation occurred and how to prevent recurrence
3. Potential additional review for related code
4. In severe cases, temporary restriction from security-sensitive areas